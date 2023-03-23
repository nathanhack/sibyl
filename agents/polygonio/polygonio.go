package polygonio

import (
	"context"
	"fmt"
	"time"

	"github.com/go-faster/errors"
	"github.com/nathanhack/sibyl/agents"
	"github.com/nathanhack/sibyl/agents/internal/barresult"
	"github.com/nathanhack/sibyl/ent"
	"github.com/nathanhack/sibyl/ent/datasource"
	"github.com/nathanhack/sibyl/ent/dividend"
	"github.com/nathanhack/sibyl/ent/interval"
	"github.com/sirupsen/logrus"

	polygon "github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/rest/models"
	"go.uber.org/ratelimit"
)

type Plan string

const (
	companyName    = "Polygon.io"
	companyAddress = "331 Elizabeth St NE Suite D, Atlanta, Georgia, 30307 US"

	Basic     = Plan("basic")
	Starter   = Plan("starter")
	Developer = Plan("developer")
	Advanced  = Plan("advanced")

	maxLimit = 50_000 //directly from polygon.io website
)

type Polygonio struct {
	datasourceID int
	plan         Plan
	ent          *ent.Client
	client       *polygon.Client
	rateLimit    ratelimit.Limiter
}

func New(ctx context.Context, client *ent.Client, apikey string, plan Plan) (*Polygonio, error) {
	id, err := client.DataSource.Query().Where(datasource.Name(companyName)).OnlyID(ctx)
	if err != nil {
		// must not exist so we create it!
		ds, err := client.DataSource.Create().
			SetName(companyName).
			SetAddress(companyAddress).
			Save(ctx)

		if err != nil {
			return nil, errors.Wrap(err, "polygonio.New")
		}
		id = ds.ID
	}

	var r ratelimit.Limiter
	if plan == Basic {
		// the website says 5/min but using this ratelimiter it seems it needs to be
		// but experiments showed it needs to be just less than that
		r = ratelimit.New(7, ratelimit.Per(2*time.Minute), ratelimit.WithoutSlack)
	} else {
		r = ratelimit.NewUnlimited()
	}

	return &Polygonio{
		datasourceID: id,
		plan:         plan,
		ent:          client,
		client:       polygon.New(apikey),
		rateLimit:    r,
	}, nil
}

// DataSourceId provides the database DataSourceID for Polygon.io
func (pio *Polygonio) DataSourceId() int {
	return pio.datasourceID
}

func (pio *Polygonio) Name() string {
	return companyName
}

func (pio *Polygonio) CompanyName() string {
	return companyName
}

func (pio *Polygonio) CompanyAddress() string {
	return companyAddress
}

// MaxTimeRange: return the valid query time range based on the particular Polygon.io plan
func (pio *Polygonio) MaxTimeRange(intervalValue interval.Interval) (time.Time, time.Time) {
	// from https://polygon.io/pricing
	now := time.Now()
	switch pio.plan {
	case Basic:
		// 2 years, and end of day yesterday
		return time.Now().AddDate(-2, 0, 0),
			time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local).AddDate(0, 0, -1)
	case Starter:
		// 5 years, and 15 min delayed
		return time.Now().AddDate(-5, 0, 0), now.Add(-15 * time.Minute)
	case Developer:
		// 10 years, and 15 min delayed
		return time.Now().AddDate(-10, 0, 0), now.Add(-15 * time.Minute)
	case Advanced:
		// 15+ years, and now
		return time.Time{}, now
	}
	panic(fmt.Errorf("invalid plan type: '%v'", pio.plan))
}

func toTimespan(val interval.Interval) models.Timespan {
	switch val {
	case interval.Interval1min:
		return models.Minute
	case interval.IntervalDaily:
		return models.Day
	case interval.IntervalMonthly:
		return models.Month
	case interval.IntervalYearly:
		return models.Year
	}
	panic(fmt.Errorf("wrong type: %v", val))
}

func (pio *Polygonio) BarRequest(ctx context.Context, ticker string, intervalValue interval.Interval, startInterval, endInterval time.Time) (*agents.BarHistoryResults, error) {
	logrus.Infof("Polygon.io processing %v request for %v {%v,%v}", intervalValue, ticker, startInterval.Local(), endInterval.Local())

	rb := barresult.Builder{
		Ent:           pio.ent,
		DataSourceID:  pio.datasourceID,
		Ticker:        ticker,
		Interval:      intervalValue,
		IntervalStart: startInterval,
		IntervalEnd:   endInterval,
	}

	// for polygonio we need to subtract one sec from the endInterval
	endInterval = endInterval.Add(-time.Second)
	for start := startInterval; start.Before(endInterval); {
		select {
		case <-ctx.Done():
			return nil, errors.Errorf("context ended")
		default:
		}

		end := endInterval

		switch intervalValue {
		case interval.Interval1min:
			if end.Sub(start) > (maxLimit-120)*time.Minute {
				end = start.Add((maxLimit - 120) * time.Minute) //we subtract 120 to give a little buffer for daylight savings differences
				end = end.Add(-1 * time.Second)                 // end is meant to be noninclusive
			}
		case interval.IntervalDaily:
			if end.After(start.AddDate(0, 0, maxLimit)) {
				end = start.AddDate(0, 0, maxLimit)
				end = end.Add(-1 * time.Second) // end is meant to be noninclusive
			}
		case interval.IntervalMonthly:
			if end.After(start.AddDate(0, maxLimit, 0)) {
				end = start.AddDate(0, maxLimit, 0)
				end = end.Add(-1 * time.Second) // end is meant to be noninclusive
			}
		case interval.IntervalYearly:
			if end.After(start.AddDate(maxLimit, 0, 0)) {
				end = start.AddDate(maxLimit, 0, 0)
				end = end.Add(-1 * time.Second) // end is meant to be noninclusive
			}
		}
		params := models.dfsParams{
			Ticker:     ticker,
			Multiplier: 1,
			Timespan:   toTimespan(intervalValue),
			From:       models.Millis(start),
			To:         models.Millis(end),
		}.WithOrder(models.Asc).WithAdjusted(false)

		pio.rateLimit.Take()
		logrus.Debugf("Polygon.io GetAggs: {%v,%v,%v,%v,%v}", params.Ticker, params.Multiplier, params.Timespan, time.Time(params.From).Local(), time.Time(params.To).Local())
		res, err := pio.client.GetAggs(ctx, params)
		if err != nil {
			//at the first sign of an error were do the call back
			return nil, errors.Wrap(err, fmt.Sprintf("Polygon.io GetAggs(%v,%v,%v,%v) Client error", ticker, intervalValue, startInterval, endInterval))
		}

		if res != nil && len(res.Results) > 0 {
			for _, r := range res.Results {
				rb.AddBar(r.Open, r.High, r.Low, r.Close, r.Volume, int32(r.Transactions), time.Time(r.Timestamp))
			}
		}

		start = end.Add(time.Second) // add back the second and set start to it
	}

	return rb.Result(), nil
}

func (pio *Polygonio) EntitySearch(ctx context.Context, ticker string, limit int) ([]agents.EntitySearchResults, error) {
	params := models.ListTickersParams{
		TickerGTE: &ticker,
	}.WithOrder(models.Asc).WithLimit(limit)

	results := make([]agents.EntitySearchResults, 0)
	pio.rateLimit.Take()
	iter := pio.client.ListTickers(ctx, params)
	for i := 0; i < limit && iter.Next(); i++ {
		results = append(results, agents.EntitySearchResults{
			Ticker: iter.Item().Ticker,
			Name:   iter.Item().Name,
		})
	}
	if iter.Err() != nil {
		logrus.Debugf("EntitySearch had an error and %v results", len(results))
		return results, iter.Err()
	}

	return results, nil
}

func (pio *Polygonio) Entity(ctx context.Context, ticker string) (*ent.EntityCreate, error) {
	params := models.GetTickerDetailsParams{
		Ticker: ticker,
	}

	pio.rateLimit.Take()
	res, err := pio.client.GetTickerDetails(ctx, &params)
	if err != nil {
		return nil, err
	}

	d := time.Time(res.Results.ListDate)
	return pio.ent.Entity.Create().
		SetActive(res.Results.Active).
		SetTicker(res.Results.Ticker).
		SetName(res.Results.Name).
		SetDescription(res.Results.Description).
		SetListDate(time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.Local)), nil
}

func (pio *Polygonio) DividendRequest(ctx context.Context, ticker string, start, end time.Time) (results []*ent.DividendCreate, payDates []time.Time, err error) {
	defer func() {
		if err != nil {
			err = errors.Wrap(err, fmt.Sprintf("Polygon.io ListDividends(%v,%v,%v) Client error", ticker, start.Local(), end.Local()))
		}
	}()
	results = make([]*ent.DividendCreate, 0)
	payDates = make([]time.Time, 0)
	startDate := models.Date(start)
	endDate := models.Date(end)

	params := models.ListDividendsParams{
		TickerEQ:   &ticker,
		PayDateGTE: &startDate,
		PayDateLT:  &endDate,
	}.WithOrder(models.Asc).WithLimit(1000)

	type dividendResult struct {
		data models.Dividend
		err  error
	}

	data := make(chan dividendResult, 1000)
	done, doneCancel := context.WithCancel(context.Background())

	go func() {
		pio.rateLimit.Take()
		logrus.Debugf("Polygon.io ListDividends: {%v,%v,%v}", *params.TickerEQ, time.Time(*params.PayDateGTE).Local(), time.Time(*params.PayDateLT).Local())
		iter := pio.client.ListDividends(ctx, params)

		for iter.Next() {
			data <- dividendResult{
				data: iter.Item(),
			}
		}
		if iter.Err() != nil {
			data <- dividendResult{
				err: iter.Err(),
			}
		} else {
			doneCancel()
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil, nil, fmt.Errorf("context ended")
		case d := <-data:
			if d.err != nil {
				return nil, nil, err
			}

			declarationDate, err := time.ParseInLocation("2006-01-02", d.data.DeclarationDate, time.Local)
			if err != nil {
				return nil, nil, err
			}
			exDividendDate, err := time.ParseInLocation("2006-01-02", d.data.ExDividendDate, time.Local)
			if err != nil {
				return nil, nil, err
			}
			payDate, err := time.ParseInLocation("2006-01-02", d.data.PayDate, time.Local)
			if err != nil {
				return nil, nil, err
			}
			recordDate, err := time.ParseInLocation("2006-01-02", d.data.RecordDate, time.Local)
			if err != nil {
				return nil, nil, err
			}

			results = append(results, pio.ent.Dividend.Create().
				SetCashAmount(d.data.CashAmount).
				SetDeclarationDate(declarationDate).
				SetExDividendDate(exDividendDate).
				SetFrequency(int(d.data.Frequency)).
				SetPayDate(payDate).
				SetRecordDate(recordDate).
				SetDividendType(dividend.DividendType(d.data.DividendType)), //the types came from Polygon.io so a straight case is fine
			)
			payDates = append(payDates, payDate)
		case <-done.Done():
			return results, payDates, nil
		}
	}
}

func (pio *Polygonio) SplitRequest(ctx context.Context, ticker string, start, end time.Time) (results []*ent.SplitCreate, executionDates []time.Time, err error) {
	defer func() {
		if err != nil {
			err = errors.Wrap(err, fmt.Sprintf("Polygon.io SplitRequest(%v,%v,%v) Client error", ticker, start.Local(), end.Local()))
		}
	}()
	results = make([]*ent.SplitCreate, 0)
	executionDates = make([]time.Time, 0)
	startDate := models.Date(start)
	endDate := models.Date(end)

	params := models.ListSplitsParams{
		TickerEQ:         &ticker,
		ExecutionDateGTE: &startDate,
		ExecutionDateLT:  &endDate,
	}.WithOrder(models.Asc).WithLimit(1000)

	type splitResult struct {
		data models.Split
		err  error
	}
	data := make(chan splitResult, 1000)
	done, doneCancel := context.WithCancel(context.Background())
	go func() {
		pio.rateLimit.Take()
		logrus.Debugf("Polygon.io SplitRequest: {%v,%v,%v}", *params.TickerEQ, time.Time(*params.ExecutionDateGTE).Local(), time.Time(*params.ExecutionDateLT).Local())
		iter := pio.client.ListSplits(ctx, params)
		for iter.Next() {
			data <- splitResult{
				data: iter.Item(),
			}
		}
		if iter.Err() != nil {
			data <- splitResult{
				err: iter.Err(),
			}
		} else {
			doneCancel()
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil, nil, fmt.Errorf("context ended")
		case d := <-data:
			if d.err != nil {
				return nil, nil, err
			}
			executionDate, err := time.Parse("2006-01-02", d.data.ExecutionDate)
			if err != nil {
				return nil, nil, err
			}

			results = append(results, pio.ent.Split.Create().
				SetExecutionDate(executionDate).
				SetFrom(d.data.SplitFrom).
				SetTo(d.data.SplitTo),
			)
			executionDates = append(executionDates, executionDate)
		case <-done.Done():
			return results, executionDates, nil
		}
	}
}
