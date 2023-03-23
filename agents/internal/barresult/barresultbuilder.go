package barresult

import (
	"time"

	"github.com/nathanhack/sibyl/agents"
	"github.com/nathanhack/sibyl/ent"
	"github.com/nathanhack/sibyl/ent/interval"
)

type Builder struct {
	Ent           *ent.Client
	DataSourceID  int
	Ticker        string
	Interval      interval.Interval
	IntervalStart time.Time
	IntervalEnd   time.Time
	result        *agents.BarHistoryResults
	lastBarTime   *time.Time
	currGroup     *ent.BarGroupCreate
	currBars      []*ent.BarRecordCreate
}

func (b *Builder) AddBar(open, high, low, close, volume float64, transactions int32, timestamp time.Time) {
	if b.lazyInit() {
		//when it's the very first bar
		b.result.FirstBarTimestamp = timestamp
	}

	newGroup := false
	if b.lastBarTime != nil {
		calcNext := *b.lastBarTime
		switch b.Interval {
		case interval.Interval1min:
			calcNext = calcNext.Add(time.Minute)
		case interval.IntervalDaily:
			calcNext = calcNext.AddDate(0, 0, 1)
		case interval.IntervalMonthly:
			calcNext = calcNext.AddDate(0, 1, 0)
		case interval.IntervalYearly:
			calcNext = calcNext.AddDate(1, 0, 0)
		}
		if calcNext.Before(timestamp) {
			newGroup = true
		}
	}

	if newGroup || b.lastBarTime == nil || b.currGroup == nil {
		b.currGroup = b.Ent.BarGroup.Create().
			SetFirst(time.Time(timestamp))

		b.currBars = make([]*ent.BarRecordCreate, 0)
		b.result.BarGroups = append(b.result.BarGroups, b.currGroup)
		b.result.Bars = append(b.result.Bars, b.currBars)
	}

	//update currBars
	b.currBars = append(b.currBars, b.createBarRecord(open, high, low, close, volume, transactions, timestamp))
	b.result.Bars[len(b.result.Bars)-1] = b.currBars // currBars could point to some place new so update it

	//update currGroup
	b.currGroup.SetLast(timestamp).
		SetCount(len(b.currBars))

	//update the rest
	b.lastBarTime = &timestamp
	b.result.BarCount++
}

func (b *Builder) Result() *agents.BarHistoryResults {
	b.lazyInit()
	return b.result
}

func (b *Builder) lazyInit() bool {
	if b.result != nil {
		return false
	}

	b.result = &agents.BarHistoryResults{
		DataSourceID:      b.DataSourceID,
		Ticker:            b.Ticker,
		Interval:          b.Interval,
		IntervalStart:     b.IntervalStart,
		IntervalEnd:       b.IntervalEnd,
		FirstBarTimestamp: b.IntervalStart, // it defaults to the IntervalStart
		BarGroups:         make([]*ent.BarGroupCreate, 0),
		Bars:              make([][]*ent.BarRecordCreate, 0),
	}
	return true
}

func (b *Builder) createBarRecord(open, high, low, close, volume float64, transactions int32, timestamp time.Time) *ent.BarRecordCreate {
	return b.Ent.BarRecord.Create().
		SetOpen(open).
		SetHigh(high).
		SetLow(low).
		SetClose(close).
		SetVolume(volume).
		SetTransactions(transactions).
		SetTimestamp(timestamp)
}
