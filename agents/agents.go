package agents

import (
	"context"
	"time"

	"github.com/nathanhack/sibyl/ent"
	"github.com/nathanhack/sibyl/ent/interval"
)

type BarHistoryResults struct {
	DataSourceID      int
	Ticker            string
	Interval          interval.Interval
	IntervalStart     time.Time
	IntervalEnd       time.Time
	FirstBarTimestamp time.Time
	BarGroups         []*ent.BarGroupCreate
	Bars              [][]*ent.BarRecordCreate
	BarCount          int
}

type BarRequester interface {
	DataSourceId() int
	Name() string //Agent Name
	MaxTimeRange(intervalValue interval.Interval) (start time.Time, end time.Time)
	BarRequest(ctx context.Context, ticker string, intervalValue interval.Interval, start, end time.Time) (*BarHistoryResults, error)
}

type EntitySearchResults struct {
	Ticker string
	Name   string
}

type EntitySearcher interface {
	Name() string // Agent Name
	EntitySearch(ctx context.Context, ticker string, limit int) ([]EntitySearchResults, error)
}

type EntitySourcing interface {
	Name() string // Agent Name
	EntityCreate(ctx context.Context, ticker string) (*ent.EntityCreate, error)
	EntityCreateUpdate(ctx context.Context, current ...*ent.Entity) ([]*ent.EntityUpdateOne, error)
}

type DividendRequester interface {
	Name() string // Agent Name
	DividendRequest(ctx context.Context, ticker string, start, end time.Time) ([]*ent.DividendCreate, []time.Time, error)
}

type SplitRequester interface {
	Name() string // Agent Name
	SplitRequest(ctx context.Context, ticker string, start, end time.Time) ([]*ent.SplitCreate, []time.Time, error)
}

type MarketHoursRequester interface {
	Name() string // Agent Name
	MarketHoursRequest(ctx context.Context, start, end time.Time) ([]*ent.MarketHoursCreate, []time.Time, error)
}
