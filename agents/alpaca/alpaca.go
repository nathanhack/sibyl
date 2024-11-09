package alpaca

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/civil"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/nathanhack/sibyl/agents"
	"github.com/nathanhack/sibyl/agents/internal/barresult"
	"github.com/nathanhack/sibyl/ent"
	"github.com/nathanhack/sibyl/ent/datasource"
	"github.com/nathanhack/sibyl/ent/interval"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.uber.org/ratelimit"
)

type Plan string

const (
	companyName    = "Alpaca"
	companyAddress = ""

	Free      = Plan("free")
	Unlimited = Plan("unlimited")
)

type Alpaca struct {
	datasourceID int
	ent          *ent.Client
	plan         Plan
	client       *alpaca.Client
	dataClient   *marketdata.Client
	rateLimit    ratelimit.Limiter
}

func New(ctx context.Context, client *ent.Client, apikey, apiSecret, url string, plan Plan) (*Alpaca, error) {
	id, err := client.DataSource.Query().Where(datasource.Name(companyName)).OnlyID(ctx)
	if err != nil {
		// must not exist so we create it!
		ds, err := client.DataSource.Create().
			SetName(companyName).
			SetAddress(companyAddress).
			Save(ctx)

		if err != nil {
			return nil, errors.Wrap(err, "alpaca.New")
		}
		id = ds.ID
	}

	// for rate limit details https://alpaca.markets/docs/market-data/
	var r ratelimit.Limiter
	if plan == Free {
		r = ratelimit.New(200, ratelimit.Per(time.Minute))
	} else {
		r = ratelimit.NewUnlimited()
	}
	aplacaClient := alpaca.NewClient(alpaca.ClientOpts{
		APIKey:    apikey,
		APISecret: apiSecret,
		BaseURL:   url,
	})

	dataClient := marketdata.NewClient(marketdata.ClientOpts{
		APIKey:    apikey,
		APISecret: apiSecret,
	})

	return &Alpaca{
		datasourceID: id,
		ent:          client,
		plan:         plan,
		client:       aplacaClient,
		dataClient:   dataClient,
		rateLimit:    r,
	}, nil
}

// DataSourceId provides the database DataSourceID for Alpaca
func (am *Alpaca) DataSourceId() int {
	return am.datasourceID
}

func (am *Alpaca) Name() string {
	return companyName
}

// MaxTimeRange returns the valid time range based on the particular Alpaca plan
func (am *Alpaca) MaxTimeRange(intervalValue interval.Interval) (time.Time, time.Time) {
	//they have data up to 2016 but they said it's spotty at the beginning
	// so we limit to 2017.  they said they'll be expanding at some point
	// https://alpaca.markets/support/alpaca-data-timeline
	start := time.Date(2016, 12, 15, 0, 0, 0, 0, time.Local)
	if am.plan == Free {
		// 15 min delay for data
		return start, time.Now().Add(-15 * time.Minute)
	}
	return start, time.Now()
}

func toTimeFrame(val interval.Interval) marketdata.TimeFrame {
	switch val {
	case interval.Interval1min:
		return marketdata.OneMin
	case interval.IntervalDaily:
		return marketdata.OneDay
	case interval.IntervalMonthly:
		return marketdata.OneMonth
	case interval.IntervalYearly:
		return marketdata.NewTimeFrame(12, marketdata.Month)
	}
	panic(fmt.Errorf("wrong type: %v", val))
}

func (am *Alpaca) BarRequest(ctx context.Context, ticker string, intervalValue interval.Interval, startInterval, endInterval time.Time) (*agents.BarHistoryResults, error) {
	logrus.Infof("Alpaca.BarRequest processing %v request for %v {%v,%v}", intervalValue, ticker, startInterval.Local(), endInterval.Local())
	rb := barresult.Builder{
		Ent:           am.ent,
		DataSourceID:  am.datasourceID,
		Ticker:        ticker,
		Interval:      intervalValue,
		IntervalStart: startInterval,
		IntervalEnd:   endInterval,
	}

	// The End filter equal to or before this time and we want less than the endInterval
	barsRequest := marketdata.GetBarsRequest{
		TimeFrame:  toTimeFrame(intervalValue),
		Adjustment: marketdata.Raw,
		Start:      startInterval,
		End:        endInterval.Add(-time.Second),
	}
	am.rateLimit.Take()
	logrus.Debugf("Alpaca.BarRequest GetBarsAsync: {%v,%v,%v,%v,%v}", ticker, barsRequest.TimeFrame, barsRequest.Adjustment, barsRequest.Start.Local(), barsRequest.End.Local())

	type barResultsType struct {
		B []marketdata.Bar
		E error
	}

	bars := make(chan barResultsType)
	go func() {
		defer close(bars)		
		barItems, err := am.dataClient.GetBars(ticker, barsRequest)
		bars <- barResultsType{
			B: barItems,
			E: err,
		}
	}()

	var barResults barResultsType
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context ended")
	case barResults = <-bars:
	}
	barItems, err := barResults.B, barResults.E

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Alpaca GetBarsAsync(%v,%v,%v,%v) Client error", ticker, intervalValue, startInterval.Local(), endInterval.Local()))
	}
	for _, barItem := range barItems {
		select {
		case <-ctx.Done():
			return nil, errors.Errorf("context ended")
		default:
		}

		rb.AddBar(barItem.Open, barItem.High, barItem.Low, barItem.Close, float64(barItem.Volume), int32(barItem.TradeCount), barItem.Timestamp)
	}

	return rb.Result(), nil
}

func (am *Alpaca) MarketHoursRequest(ctx context.Context, start, end time.Time) ([]*ent.MarketHoursCreate, []time.Time, error) {
	logrus.Infof("Alpaca.MarketHoursRequest processing request for {%v,%v}", start.Local(), end.Local())

	am.rateLimit.Take()
	logrus.Debugf("Alpaca.BarRequest GetCalendar: {%v,%v}", start, end)

	calendarRequest := alpaca.GetCalendarRequest{
		Start: start,
		End:   end,
	}

	days, err := am.client.GetCalendar(calendarRequest)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "Alpaca.BarRequest GetCalendar: {%v,%v}", start, end)
	}

	results := make([]*ent.MarketHoursCreate, len(days))
	dates := make([]time.Time, len(days))
	dateFormat := "2006-01-02"
	format := "2006-01-02 15:04"
	for i, day := range days {
		d, err := time.ParseInLocation(dateFormat, day.Date, time.Local)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "Alpaca.BarRequest date: %v // %v", day, day.Date)
		}
		tmp := fmt.Sprintf("%v %v", day.Date, day.Open)
		s, err := time.ParseInLocation(format, tmp, time.Local)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "Alpaca.BarRequest open: %v // %v", day, tmp)
		}
		tmp = fmt.Sprintf("%v %v", day.Date, day.Close)
		e, err := time.ParseInLocation(format, tmp, time.Local)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "Alpaca.BarRequest close: %v // %v", day, tmp)
		}
		results[i] = am.ent.MarketHours.Create().
			SetDate(d).
			SetStartTime(s).
			SetEndTime(e)

		dates[i] = d
	}

	return results, dates, nil
}

func (am *Alpaca) DividendRequest(ctx context.Context, ticker string, start, end time.Time) ([]*ent.DividendCreate, []time.Time, error) {
	request := marketdata.GetCorporateActionsRequest{
		Symbols: []string{ticker},
		Types:   []string{"cash_dividend"},
		Start:   civil.DateOf(start),
		End:     civil.DateOf(end),
	}
	am.rateLimit.Take()
	action, err := am.dataClient.GetCorporateActions(request)
	if err != nil {
		return nil, nil, errors.Wrap(err, "alpaca.DividendRequest")
	}

	results := make([]*ent.DividendCreate, 0)
	paydates := make([]time.Time, 0)
	for _, cd := range action.CashDividends {
		if cd.PayableDate == nil || cd.RecordDate == nil {
			continue
		}
		results = append(results, am.ent.Dividend.Create().
			SetRate(cd.Rate).
			SetDeclarationDate(cd.ProcessDate.In(time.Local)).
			SetExDividendDate(cd.ExDate.In(time.Local)).
			SetPayDate(cd.PayableDate.In(time.Local)).
			SetRecordDate(cd.RecordDate.In(time.Local)),
		)
	}

	return results, paydates, nil
}

func (am *Alpaca) EntityCreate(ctx context.Context, ticker string) (*ent.EntityCreate, error) {

	am.rateLimit.Take()
	res, err := am.client.GetAsset(ticker)
	if err != nil {
		return nil, errors.Wrap(err, "Alpaca Entity")
	}

	hasOptions := false
	for _, attribute := range res.Attributes {
		if attribute == "has_options" {
			hasOptions = true
			break
		}
	}

	return am.ent.Entity.Create().
		SetActive(res.Status == alpaca.AssetActive).
		SetTicker(res.Symbol).
		SetName(res.Name).
		SetOptions(hasOptions).
		SetTradable(res.Tradable), nil
}

func (am *Alpaca) EntityCreateUpdate(ctx context.Context, current ...*ent.Entity) ([]*ent.EntityUpdateOne, error) {
	updates := make([]*ent.EntityUpdateOne, 0, len(current))
	for _, entity := range current {
		am.rateLimit.Take()
		res, err := am.client.GetAsset(entity.Ticker)
		if err != nil {
			return updates, err
		}
		hasOptions := false
		for _, attribute := range res.Attributes {
			if attribute == "has_options" {
				hasOptions = true
				break
			}
		}
		updates = append(updates,
			am.ent.Entity.UpdateOne(entity).
				SetActive(res.Status == alpaca.AssetActive).
				SetTicker(res.Symbol).
				SetName(res.Name).
				SetOptions(hasOptions).
				SetTradable(res.Tradable),
		)
	}

	return updates, nil
}

func (am *Alpaca) SplitRequest(ctx context.Context, ticker string, start, end time.Time) ([]*ent.SplitCreate, []time.Time, error) {

	request := marketdata.GetCorporateActionsRequest{
		Symbols: []string{ticker},
		Types:   []string{"cash_dividend"},
		Start:   civil.DateOf(start),
		End:     civil.DateOf(end),
	}
	am.rateLimit.Take()
	action, err := am.dataClient.GetCorporateActions(request)
	if err != nil {
		return nil, nil, errors.Wrap(err, "alpaca.DividendRequest")
	}

	results := make([]*ent.SplitCreate, 0)
	paydates := make([]time.Time, 0)
	for _, split := range action.ForwardSplits {

		results = append(results, am.ent.Split.Create().
			SetExecutionDate(split.ExDate.In(time.Local)).
			SetFrom(split.OldRate).
			SetTo(split.NewRate),
		)
	}
	for _, split := range action.ReverseSplits {

		results = append(results, am.ent.Split.Create().
			SetExecutionDate(split.ExDate.In(time.Local)).
			SetFrom(split.OldRate).
			SetTo(split.NewRate),
		)
	}

	for _, split := range action.UnitSplits {
		results = append(results, am.ent.Split.Create().
			SetExecutionDate(split.EffectiveDate.In(time.Local)).
			SetFrom(split.OldRate).
			SetTo(split.NewRate),
		)
	}

	return results, paydates, nil
}