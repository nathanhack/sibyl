package alpaca

import (
	"context"
	"fmt"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
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
	client       alpaca.Client
	dataClient   marketdata.Client
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
		ApiKey:    apikey,
		ApiSecret: apiSecret,
		BaseURL:   url,
	})

	dataClient := marketdata.NewClient(marketdata.ClientOpts{
		ApiKey:    apikey,
		ApiSecret: apiSecret,
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
func (pio *Alpaca) DataSourceId() int {
	return pio.datasourceID
}

func (pio *Alpaca) Name() string {
	return companyName
}

// MaxTimeRange returns the valid time range based on the particular Alpaca plan
func (pio *Alpaca) MaxTimeRange(intervalValue interval.Interval) (time.Time, time.Time) {
	//they have data up to 2016 but they said it's spotty at the beginning
	// so we limit to 2017.  they said they'll be expanding at some point
	// https://alpaca.markets/support/alpaca-data-timeline
	start := time.Date(2016, 12, 15, 0, 0, 0, 0, time.Local)
	if pio.plan == Free {
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

func (pio *Alpaca) BarRequest(ctx context.Context, ticker string, intervalValue interval.Interval, startInterval, endInterval time.Time) (*agents.BarHistoryResults, error) {
	logrus.Infof("Alpaca.BarRequest processing %v request for %v {%v,%v}", intervalValue, ticker, startInterval.Local(), endInterval.Local())
	rb := barresult.Builder{
		Ent:           pio.ent,
		DataSourceID:  pio.datasourceID,
		Ticker:        ticker,
		Interval:      intervalValue,
		IntervalStart: startInterval,
		IntervalEnd:   endInterval,
	}

	// The End filter equal to or before this time and we want less than the endInterval
	params := marketdata.GetBarsParams{
		TimeFrame:  toTimeFrame(intervalValue),
		Adjustment: marketdata.Raw,
		Start:      startInterval,
		End:        endInterval.Add(-time.Second),
	}
	pio.rateLimit.Take()
	logrus.Debugf("Alpaca.BarRequest GetBarsAsync: {%v,%v,%v,%v,%v}", ticker, params.TimeFrame, params.Adjustment, params.Start.Local(), params.End.Local())

	for barItem := range pio.dataClient.GetBarsAsync(ticker, params) {
		select {
		case <-ctx.Done():
			return nil, errors.Errorf("context ended")
		default:
		}

		if err := barItem.Error; err != nil {
			//at the first sign of an error were do the call back
			return nil, errors.Wrap(err, fmt.Sprintf("Alpaca GetBarsAsync(%v,%v,%v,%v) Client error", ticker, intervalValue, startInterval.Local(), endInterval.Local()))

		}

		rb.AddBar(barItem.Bar.Open, barItem.Bar.High, barItem.Bar.Low, barItem.Bar.Close, float64(barItem.Bar.Volume), int32(barItem.Bar.TradeCount), barItem.Bar.Timestamp)
	}

	return rb.Result(), nil
}

func (pio *Alpaca) MarketHoursRequest(ctx context.Context, start, end time.Time) ([]*ent.MarketHoursCreate, []time.Time, error) {
	logrus.Infof("Alpaca.MarketHoursRequest processing request for {%v,%v}", start.Local(), end.Local())

	format := "2006-01-02"
	startDate := start.Format(format)
	endDate := end.Format(format)

	pio.rateLimit.Take()
	logrus.Debugf("Alpaca.BarRequest GetCalendar: {%v,%v}", startDate, endDate)

	days, err := pio.client.GetCalendar(&startDate, &endDate)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "Alpaca.BarRequest GetCalendar: {%v,%v}", startDate, endDate)
	}

	results := make([]*ent.MarketHoursCreate, len(days))
	dates := make([]time.Time, len(days))
	dateFormat := "2006-01-02"
	format = "2006-01-02 15:04"
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
		results[i] = pio.ent.MarketHours.Create().
			SetDate(d).
			SetStartTime(s).
			SetEndTime(e)

		dates[i] = d
	}

	return results, dates, nil

}
