package barrequester

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-faster/errors"
	"github.com/nathanhack/sibyl/agents"
	"github.com/nathanhack/sibyl/cmd/server/cmd/internal"
	"github.com/nathanhack/sibyl/ent"
	"github.com/nathanhack/sibyl/ent/bartimerange"
	"github.com/nathanhack/sibyl/ent/datasource"
	"github.com/nathanhack/sibyl/ent/entity"
	"github.com/nathanhack/sibyl/ent/interval"
	"github.com/sirupsen/logrus"
)

func Grabber(ctx context.Context, client *ent.Client, agent agents.BarRequester, wg *sync.WaitGroup) {
	logrus.Infof("Bars.Grabber(%v): Running", agent.Name())
	startupTimer := time.NewTimer(5 * time.Second)
	wg.Add(1)
	defer func() {
		logrus.Infof("Bars.Grabber(%v): Stopped", agent.Name())
		wg.Done()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-startupTimer.C:
		case <-internal.AtMidnight():
		}
		//let's clean up any BarTimeRange's that are still pending older than 30 minutes
		removeBadBarTimeRangesOlderThan(ctx, client, time.Now().Add(-30*time.Minute))

		//we first get all the stocks
		stocks, err := client.Entity.Query().Where(
			entity.Active(true),
		).All(ctx)

		if err != nil {
			logrus.Errorf("Bars.Grabber(%v): failed to get stocks: %v", agent.Name(), err)
			continue
		}

		for _, stock := range stocks {
			//we start off with a timerange that is the intersection of the stocks' listed date to today and what agent
			// can give us for the plan

			for _, intervalValue := range []interval.Interval{
				interval.Interval1min,
				interval.IntervalDaily,
				interval.IntervalMonthly,
				interval.IntervalYearly} {
				select {
				case <-ctx.Done():
					return
				default:
				}

				stockInterval, err := stock.QueryIntervals().Where(
					interval.ActiveEQ(true),
					interval.IntervalEQ(intervalValue),
					interval.HasDataSourceWith(datasource.ID(agent.DataSourceId())),
				).Only(ctx)
				if err != nil {
					switch err.(type) {
					case *ent.NotFoundError:
					default:
						logrus.Warnf("Bars.Grabber(%v):%v", agent.Name(), err)
					}
					continue
				}

				//we clean the interval because we want to make accurate requests
				err = clean(ctx, client, stockInterval)
				if err != nil {
					//we need to log this.. but not let it keep us from requesting new data
					logrus.Errorf("Bars.Grabber(%v): clean error for %v (%v): %v", agent.Name(), stock.Ticker, intervalValue, err)
				}

				s, e := agent.MaxTimeRange(intervalValue)
				agentInterval := internal.TimeInterval{Start: truncateTime(s, intervalValue), End: truncateTime(e, intervalValue)}
				historyInterval := createInterval(stock.ListDate, intervalValue)
				mainTimeRange := internal.IntervalIntersection(historyInterval, agentInterval)
				timeranges, err := stockInterval.QueryBars().All(ctx)
				if err != nil {
					switch err.(type) {
					case *ent.NotFoundError:
					default:
						logrus.Warnf("Bars.Grabber(%v): QueryBars error: %v", agent.Name(), err)
					}
					continue
				}

				//now get all the timeranges from the interval and
				// remove them from the mainTimerange that will requested
				for _, t := range timeranges {
					mainTimeRange = internal.IntervalDifferenceSlice(mainTimeRange, internal.TimeInterval{
						Start: t.Start,
						End:   t.End,
					})
				}

				//now we have some set of timeranges in mainTimerange
				// now put in the request to download the data
				for _, t := range mainTimeRange {
					results, err := agent.BarRequest(ctx, stock.Ticker, intervalValue, t.Start, t.End)
					if err != nil {
						logrus.Errorf("Bars.Grabber(%v): BarHistoryRequest req: %v", agent.Name(), err)
						continue
					}

					limitedCtx, limitedCancel := context.WithTimeout(ctx, 10*time.Minute)
					err = processResults(limitedCtx, agent.Name(), client, stockInterval, results)
					limitedCancel()
					if err != nil {
						logrus.Errorf("Bars.Grabber(%v): Process: %v", agent.Name(), err)
					}
				}

				needsConsolidating, err := stockInterval.QueryBars().
					Where(bartimerange.StatusIn(bartimerange.StatusClean)).
					Exist(ctx)
				if err != nil {
					logrus.Errorf("Bars.Grabber(%v): query before consolidating: %v", agent.Name(), err)
					continue
				}

				if needsConsolidating {
					//We want to get rid of duplication by consolidating BarTimeRange and BarGroups and making sure we
					// have unique BarRecords so we consolidate
					logrus.Debugf("Bars.Grabber(%v): consolidating", agent.Name())
					limitedCtx, limitedCancel := context.WithTimeout(ctx, 10*time.Minute)
					err = consolidate(limitedCtx, client, stockInterval)
					limitedCancel()
					if err != nil {
						err = errors.Wrapf(err, "consolidate error for %v (%v) : %v", stock.Ticker, stockInterval.Interval, err)
						logrus.Errorf("Bars.Grabber(%v): consolidate: %v", agent.Name(), err)
					}
				}
				logrus.Infof("Bars.Grabber(%v): %v (%v): done", agent.Name(), stock.Ticker, intervalValue)
			}
		}
	}
}

func pickStartingInterval(originalInterval, firstBar time.Time) time.Time {
	//TODO this should be smarter?, if the original interval was bigger and included a
	// weekend or holiday for now we truncate based on the data
	// side note: could use take Ticker Event (ticker_change) into account
	// once implemented, BUT maybe not, it's only a problem with Polygon.io that
	// truncate history when tickers change names (plus they have said they're working on fixing it)

	return firstBar
}

func addTimeRange(ctx context.Context, client *ent.Client, stockIntervalID int, result *agents.BarHistoryResults) (*ent.BarTimeRange, error) {
	// with the interval we can create the BarTimeRange and BarGroup
	barTimeRange, err := client.BarTimeRange.Create().
		SetIntervalID(stockIntervalID).
		SetStart(pickStartingInterval(result.IntervalStart, result.FirstBarTimestamp)).
		SetEnd(result.IntervalEnd).
		SetCount(len(result.BarGroups)).
		Save(ctx)

	if err != nil {
		return nil, errors.Wrap(err, "addTimeRange failed Create")
	}

	//update the BarGroup barTimerangeId
	for _, g := range result.BarGroups {
		g.SetTimeRangeID(barTimeRange.ID)
	}
	delta := 65535/9 - 1
	//make the BarGroups
	barGroups := make([]*ent.BarGroup, 0, len(result.BarGroups))
	for i := 0; i < len(result.BarGroups); i += delta {
		end := i + delta
		if end > len(result.BarGroups) {
			end = len(result.BarGroups)
		}
		tmp, err := client.BarGroup.CreateBulk(result.BarGroups[i:end]...).Save(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "addTimeRange failed CreateBulk1")
		}
		barGroups = append(barGroups, tmp...)
	}

	//now update the all the BarRecord groupid's
	records := make([]*ent.BarRecordCreate, 0, result.BarCount)
	for i, g := range barGroups {
		for _, b := range result.Bars[i] {
			b.SetGroupID(g.ID)

		}
		records = append(records, result.Bars[i]...)
	}

	for i := 0; i < len(records); i += delta {
		end := i + delta
		if end > len(records) {
			end = len(records)
		}

		err = client.BarRecord.CreateBulk(records[i:end]...).Exec(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "addTimeRange failed CreateBulk2")
		}
	}

	barTimeRange, err = barTimeRange.Update().SetStatus(bartimerange.StatusCreated).Save(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "addTimeRange failed Update")
	}

	return barTimeRange, nil
}

func processResults(ctx context.Context, agentName string, client *ent.Client, stockInterval *ent.Interval, result *agents.BarHistoryResults) error {
	if len(result.BarGroups) != len(result.Bars) {
		return fmt.Errorf("number of BarGroups(%v) should equal groups of Bars(%v)", len(result.BarGroups), len(result.Bars))
	}

	// every set of bars is connected to its entity by this chain
	// entity->interval->BarTimeRange->BarGroup->BarRecord
	// The interval will be unique for (stock,result.Interval,datasource)
	// Each interval may have multiple BarTimeRanges (ideally there would be just 1 but it allows for more than 1)
	// The BarTimeRanges should have times [start,end) that span a given time range
	// Each BarTimeRange must have >0 BarGroups
	// The BarGroup in the exact time range [firstTimestamp, lastTimestamp] from the BarRecord

	//the goal of this function is to insert new data

	logrus.Debugf("Bars.Grabber(%v): initial add %v bars {%v,%v}", agentName, result.BarCount, result.IntervalStart.Local(), result.IntervalEnd.Local())
	bar, err := addTimeRange(ctx, client, stockInterval.ID, result)
	if err != nil {
		return errors.Wrapf(err, "addTimeRange error for %v (%v)", result.Ticker, result.Interval)
	}

	logrus.Debugf("Bars.Grabber(%v): cleaning bar:%v", agentName, bar)
	err = cleanBarTimeRange(ctx, client, bar)
	if err != nil {
		return errors.Wrapf(err, "cleanBarTimeRange error for %v (%v)", result.Ticker, result.Interval)
	}

	logrus.Debugf("Bars.Grabber(%v): processResults done %v (%v)", agentName, result.Ticker, result.Interval)
	return nil
}

func truncateTime(val time.Time, delta interval.Interval) time.Time {
	switch delta {
	case interval.IntervalTrades:
		return val.Truncate(time.Second)
	case interval.Interval1min:
		return val.Truncate(time.Minute)
	case interval.IntervalDaily:
		return time.Date(val.Year(), val.Month(), val.Day(), 0, 0, 0, 0, time.Local)
	case interval.IntervalMonthly:
		return time.Date(val.Year(), val.Month(), 1, 0, 0, 0, 0, time.Local)
	case interval.IntervalYearly:
		return time.Date(val.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	}

	panic(fmt.Errorf("unsupported interval: %v", val))
}

// createInterval create a [start, now] both truncated to the interval
func createInterval(start time.Time, val interval.Interval) internal.TimeInterval {
	now := time.Now()
	switch val {
	case interval.IntervalTrades:
		return internal.TimeInterval{Start: start.Truncate(time.Second), End: now.Truncate(time.Second)}
	case interval.Interval1min:
		return internal.TimeInterval{Start: start.Truncate(time.Minute), End: now.Truncate(time.Minute)}
	case interval.IntervalDaily:
		return internal.TimeInterval{
			Start: time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.Local),
			End:   time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local),
		}
	case interval.IntervalMonthly:
		return internal.TimeInterval{
			Start: time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, time.Local),
			End:   time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local),
		}
	case interval.IntervalYearly:
		return internal.TimeInterval{
			Start: time.Date(start.Year(), 1, 1, 0, 0, 0, 0, time.Local),
			End:   time.Date(now.Year()+1, 1, 1, 0, 0, 0, 0, time.Local),
		}
	}

	panic(fmt.Errorf("unsupported interval: %v", val))
}
