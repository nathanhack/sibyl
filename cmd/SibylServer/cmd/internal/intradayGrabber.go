package internal

import (
	"context"
	"fmt"
	"github.com/nathanhack/sibyl/core"
	"github.com/nathanhack/sibyl/core/database"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

const (
	maxISRRetry             = 1
	defaultISRStartingDelta = 14 * 24 * time.Hour  // 14 days
	defaultIntradyStart     = 5 * 8760 * time.Hour //5 years - tends to only be 1 year available from brokers  //TODO consider making this a configuration
	defaultStartTime        = 8 * time.Hour        //the amount of time after midnight to start minute updates //TODO consider making this a configuration
	defaultStopTime         = 17 * time.Hour       //the amount of time after midnight to stop minute updates  //TODO consider making this a configuration
)

type IntradayGrabber struct {
	killCtx    context.Context
	kill       context.CancelFunc
	doneCtx    context.Context
	done       context.CancelFunc
	db         *database.SibylDatabase
	stockCache *StockCache
	running    bool
}

//NewIntradayGrabber create a new intraday history grabber.  Two main goals first once a day do a round looking
// for all intraday history upto defaultIntradyStart back in time once a day. And secondly if the state is "active"
// then every 5 seconds do a look up of the latest during the times defaultStartTime - defaultStopTime (local)
func NewIntradayGrabber(db *database.SibylDatabase, symbolCache *StockCache) *IntradayGrabber {
	killCtx, kill := context.WithCancel(context.Background())
	doneCtx, done := context.WithCancel(context.Background())
	return &IntradayGrabber{
		killCtx:    killCtx,
		kill:       kill,
		doneCtx:    doneCtx,
		done:       done,
		db:         db,
		stockCache: symbolCache,
	}
}

const (
	deadlineTickTime = 5 * time.Second
	deadline1MinTime = 1 * time.Minute
	deadline5MinTime = 5 * time.Minute

	defaultIntradayDurationToWait = 1 * time.Second
)

func (ig *IntradayGrabber) Run() error {
	if ig.running {
		return fmt.Errorf("IntradayGrabber is already running")
	}

	ig.running = true
	go func(ig *IntradayGrabber) {
		//when we start up we want to do a FULL update to right now
		// then we'll let it update the stocks as it comes up in the future

		//we start off with a 15 second
		processDeadline, processDeadlineCancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
		onceDailyDeadline, onceDailyDeadlineCancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
		deadlineTick, deadlineTickCancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
		deadline1Min, deadline1MinCancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
		deadline5Min, deadline5MinCancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))

		updateFullTickCache := make(map[core.StockSymbolType]*intradayStockRange)
		updateFull1MinCache := make(map[core.StockSymbolType]*intradayStockRange)
		updateFull5MinCache := make(map[core.StockSymbolType]*intradayStockRange)
		updateActiveTickCache := make(map[core.StockSymbolType]*intradayStockRange)
		updateActive1MinCache := make(map[core.StockSymbolType]*intradayStockRange)
		updateActive5MinCache := make(map[core.StockSymbolType]*intradayStockRange)
		finishedChan := make(chan bool, 100)
	mainLoop:
		for {
			select {
			case <-ig.killCtx.Done():
				break mainLoop
			case <-deadlineTick.Done():
				deadlineTickCancel() // (REQUIRED) we call the cancel func to release resources associated with the context
				deadlineTick, deadlineTickCancel = nextValid(core.TickInterval)
				updateActiveTickCache = addToIntradayCache(ig.killCtx, ig.db, ig.stockCache, finishedChan, updateActiveTickCache, core.TickInterval, true, false, false)
				logrus.Infof("IntradayGrabber: deadlineTick.Done() Active count(%v)", len(updateActiveTickCache))
			case <-deadline1Min.Done():
				deadline1MinCancel() // (REQUIRED) we call the cancel func to release resources associated with the context
				deadline1Min, deadline1MinCancel = nextValid(core.OneMinInterval)
				updateActive1MinCache = addToIntradayCache(ig.killCtx, ig.db, ig.stockCache, finishedChan, updateActive1MinCache, core.OneMinInterval, true, false, false)
				logrus.Infof("IntradayGrabber: deadline1Min.Done() Active count(%v)", len(updateActive1MinCache))
			case <-deadline5Min.Done():
				deadline5MinCancel() // (REQUIRED) we call the cancel func to release resources associated with the context
				deadline5Min, deadline5MinCancel = nextValid(core.FiveMinInterval)
				updateActive5MinCache = addToIntradayCache(ig.killCtx, ig.db, ig.stockCache, finishedChan, updateActive5MinCache, core.FiveMinInterval, true, false, false)
				logrus.Infof("IntradayGrabber: deadline5Min.Done() Active count(%v)", len(updateActive5MinCache))
			case <-onceDailyDeadline.Done():
				onceDailyDeadlineCancel() // (REQUIRED) we call the cancel func to release resources associated with the context
				onceDailyDeadline, onceDailyDeadlineCancel = context.WithDeadline(context.Background(), tomorrowAt6AM())

				//when the once daily deadline expires we
				// want to run all stocks
				updateFullTickCache = addToIntradayCache(ig.killCtx, ig.db, ig.stockCache, finishedChan, updateFullTickCache, core.TickInterval, false, false, true)
				updateFull1MinCache = addToIntradayCache(ig.killCtx, ig.db, ig.stockCache, finishedChan, updateFull1MinCache, core.OneMinInterval, false, false, true)
				updateFull5MinCache = addToIntradayCache(ig.killCtx, ig.db, ig.stockCache, finishedChan, updateFull5MinCache, core.FiveMinInterval, false, false, true)
				logrus.Infof("IntradayGrabber: onceDailyDeadline.Done() Ticks(%v) 1Min(%v) 5min(%v)", len(updateFullTickCache), len(updateFull1MinCache), len(updateFull5MinCache))
			case <-processDeadline.Done():
				processDeadlineCancel() // (REQUIRED) we call the cancel func to release resources associated with the context
				processDeadline, processDeadlineCancel = context.WithDeadline(context.Background(), time.Now().Add(defaultIntradayDurationToWait))
				logrus.Infof("IntradayGrabber: processDeadline()")
				currentTime := time.Now()
				agent, err := ig.db.GetAgent(ig.killCtx)
				if err != nil {
					logrus.Errorf("IntradayGrabber: could not retrieve agent: %v", err)
					//we want to run again but not immediately
					processDeadline, processDeadlineCancel = context.WithDeadline(context.Background(), currentTime.Add(5*time.Minute))
					continue
				}

				//first we check if there are any new comers to the database they are the only thing
				// checked all the time
				updateFullTickCache = addToIntradayCache(ig.killCtx, ig.db, ig.stockCache, finishedChan, updateFullTickCache, core.TickInterval, false, true, false)
				updateFull1MinCache = addToIntradayCache(ig.killCtx, ig.db, ig.stockCache, finishedChan, updateFull1MinCache, core.OneMinInterval, false, true, false)
				updateFull5MinCache = addToIntradayCache(ig.killCtx, ig.db, ig.stockCache, finishedChan, updateFull5MinCache, core.FiveMinInterval, false, true, false)
				//now we do a quick check if there's any thing to work on, if not we pass
				if len(updateFullTickCache) == 0 && len(updateFull1MinCache) == 0 && len(updateFull5MinCache) == 0 &&
					len(updateActiveTickCache) == 0 && len(updateActive1MinCache) == 0 && len(updateActive5MinCache) == 0 {
					//if there wasn't anything in there to look at
					// we just skip this round
					continue
				}

				//we remove any dups in the caches
				// found in the full cache verse the non full caches
				updateActiveTickCache = removeActiveDups(updateFullTickCache, updateActiveTickCache)
				updateActive1MinCache = removeActiveDups(updateFull1MinCache, updateActive1MinCache)
				updateActive5MinCache = removeActiveDups(updateFull5MinCache, updateActive5MinCache)

				logrus.Infof("IntadayGrabber: processing Full: Tick(%v) 1Min(%v) 5Min(%v), Active: Tick(%v) 1Min(%v) 5Min(%v)", len(updateFullTickCache), len(updateFull1MinCache), len(updateFull5MinCache), len(updateActiveTickCache), len(updateActive1MinCache), len(updateActive5MinCache))
				runningCount := 0
				ctx, cancel := context.WithCancel(ig.killCtx)
				//now we run each cache item in it's own routine
				// we start with the FullCaches first

				//update the running count
				runningCount += len(updateFullTickCache)
				runningCount += len(updateFull1MinCache)
				runningCount += len(updateFull5MinCache)
				runningCount += len(updateActiveTickCache)
				runningCount += len(updateActive1MinCache)
				runningCount += len(updateActive5MinCache)

				//kick off the routines
				kickOffRoutines(ctx, updateFullTickCache, agent, ig.db, ig.stockCache)
				kickOffRoutines(ctx, updateFull1MinCache, agent, ig.db, ig.stockCache)
				kickOffRoutines(ctx, updateFull5MinCache, agent, ig.db, ig.stockCache)
				kickOffRoutines(ctx, updateActiveTickCache, agent, ig.db, ig.stockCache)
				kickOffRoutines(ctx, updateActive1MinCache, agent, ig.db, ig.stockCache)
				kickOffRoutines(ctx, updateActive5MinCache, agent, ig.db, ig.stockCache)

				//clear the caches
				updateFullTickCache = make(map[core.StockSymbolType]*intradayStockRange)
				updateFull1MinCache = make(map[core.StockSymbolType]*intradayStockRange)
				updateFull5MinCache = make(map[core.StockSymbolType]*intradayStockRange)
				updateActiveTickCache = make(map[core.StockSymbolType]*intradayStockRange)
				updateActive1MinCache = make(map[core.StockSymbolType]*intradayStockRange)
				updateActive5MinCache = make(map[core.StockSymbolType]*intradayStockRange)

				//now we loop through finishedChan and break once we have all the results
			finishedLoop:
				for runningCount > 0 {
					select {
					case <-ig.killCtx.Done():
						break finishedLoop
					case <-finishedChan:
						runningCount--
					case <-time.After(6 * time.Hour):
						//a fail safe this should take 6 hours to complete
						logrus.Errorf("IntradayGrabber: had an issue getting all the result in a timely manner")
						break finishedLoop
					}
				}
				cancel() //kill any of the go routines still running for this round

				logrus.Infof("IntradayGrabber: finished a round in %v", time.Since(currentTime))
			}
		}
		ig.done() //signal this is finished
	}(ig)
	return nil
}

func removeActiveDups(full, active map[core.StockSymbolType]*intradayStockRange) map[core.StockSymbolType]*intradayStockRange {
	toReturn := make(map[core.StockSymbolType]*intradayStockRange)
	for symbol, stock := range active {
		if _, has := full[symbol]; !has {
			toReturn[symbol] = stock
		}
	}
	return toReturn
}

func kickOffRoutines(ctx context.Context, stockToUpdate map[core.StockSymbolType]*intradayStockRange, agent core.SibylAgent, db *database.SibylDatabase, stockCache *StockCache) {
	for _, isr := range stockToUpdate {
		//do a quick check if we've stopped the application
		if areWeDone(ctx) {
			return
		}

		go processISR(ctx, isr, agent, db, stockCache)
	}
}

func nextValid(interval core.IntradayInterval) (context.Context, context.CancelFunc) {
	now := time.Now()
	midnight := core.NewDateTypeFromTime(time.Now())
	endOfToday := midnight.Time().Add(defaultStopTime)

	if now.After(endOfToday) {
		//if not try again tomorrow (or Monday if today is friday or saturday) at market start
		switch now.Weekday() {
		case time.Friday:
			return context.WithDeadline(context.Background(), midnight.AddDate(0, 0, 3).Time().Add(defaultStartTime))
		case time.Saturday:
			return context.WithDeadline(context.Background(), midnight.AddDate(0, 0, 2).Time().Add(defaultStartTime))
		}

		return context.WithDeadline(context.Background(), midnight.AddDate(0, 0, 1).Time().Add(defaultStartTime))
	}

	switch interval {
	case core.OneMinInterval:
		return context.WithDeadline(context.Background(), now.Truncate(deadline1MinTime).Add(deadline1MinTime))
	case core.FiveMinInterval:
		return context.WithDeadline(context.Background(), now.Truncate(deadline5MinTime).Add(deadline5MinTime))
	default: // core.TickInterval:
		return context.WithDeadline(context.Background(), now.Truncate(deadlineTickTime).Add(deadlineTickTime))
	}
}

func addToIntradayCache(ctx context.Context, db *database.SibylDatabase,
	symbolCache *StockCache, finishedChan chan bool,
	oldCache map[core.StockSymbolType]*intradayStockRange,
	interval core.IntradayInterval, active bool, onlyNewlyAdded bool, fullHistoryUpdate bool) map[core.StockSymbolType]*intradayStockRange {

	toReturn := make(map[core.StockSymbolType]*intradayStockRange)
	endOfToday := core.NewDateTypeFromTime(time.Now()).Add(defaultStopTime)
	hasTick, has1Min, has5Min := interval == core.TickInterval, interval == core.OneMinInterval, interval == core.FiveMinInterval

	//copy th old into the return
	for symbol, stock := range oldCache {
		toReturn[symbol] = stock
	}

	//next add what new stuff
	for _, stock := range symbolCache.IntradayStocks(hasTick, has1Min, has5Min, active) {
		if areWeDone(ctx) {
			return toReturn //if we got the context done we end early
		}

		//if the return already has this stock then we skip it
		if _, has := toReturn[stock.Symbol]; has {
			continue
		}

		var intervalTimestamp core.TimestampType
		var nextIntervalTimestamp core.TimestampType
		var now core.TimestampType
		switch interval {
		case core.OneMinInterval:
			intervalTimestamp = stock.IntradayTimestamp1Min
			nextIntervalTimestamp = intervalTimestamp.Add(1 * time.Minute)
			now = core.NewTimestampTypeFromTime(time.Now().Add(15 * time.Second)) //add a little buffer
		case core.FiveMinInterval:
			intervalTimestamp = stock.IntradayTimestamp5Min
			nextIntervalTimestamp = intervalTimestamp.Add(5 * time.Minute)
			now = core.NewTimestampTypeFromTime(time.Now().Add(1 * time.Minute)) //add a little buffer
		default:
			intervalTimestamp = stock.IntradayTimestampTick
			nextIntervalTimestamp = intervalTimestamp.Add(1 * time.Second)
			now = core.NewTimestampTypeFromTime(time.Now().Add(250 * time.Millisecond)) //add a little buffer
		}

		if fullHistoryUpdate {
			if x, has := oldCache[stock.Symbol]; !has {
				toReturn[stock.Symbol] = makeISRFull(ctx, stock.Symbol, interval, db, finishedChan)
			} else {
				toReturn[stock.Symbol] = x
			}
		} else {
			// if either we're not current or if we're looking for onlyNewlyAdded
			current := stock.IntradayState == core.IntradayStateActive &&
				intervalTimestamp.Before(endOfToday) &&
				(nextIntervalTimestamp.Before(now) || nextIntervalTimestamp.Equal(now)) &&
				!onlyNewlyAdded

			newly := onlyNewlyAdded && intervalTimestamp.IsZero()

			if current || newly {
				if x, has := oldCache[stock.Symbol]; !has {
					toReturn[stock.Symbol] = makeISRLatestOnly(ctx, stock.Symbol, interval, db, finishedChan)
				} else {
					toReturn[stock.Symbol] = x
				}
			}
		}
	}

	return toReturn
}

func makeISRLatestOnly(ctx context.Context, stock core.StockSymbolType, interval core.IntradayInterval, db *database.SibylDatabase, finishedChan chan bool) *intradayStockRange {
	now := core.NewTimestampTypeFromTime(time.Now())
	r := make([]stockTimestampRange, 0)
	postStartTime, err := db.NewestIntradayHistoryDate(ctx, stock, interval)
	if err == nil {
		//we grab the postStartTime to now. if it's not in the DB then it will
		//return a zero which will be get it all
		r = append(r, stockTimestampRange{StartDate: postStartTime, EndDate: now})
	}
	return &intradayStockRange{
		Stock:    stock,
		Delta:    defaultISRStartingDelta,
		Ranges:   r,
		Retry:    maxISRRetry,
		Interval: interval,
		Finished: finishedChan,
	}
}

func makeISRFull(ctx context.Context, stock core.StockSymbolType, interval core.IntradayInterval, db *database.SibylDatabase, finishedChan chan bool) *intradayStockRange {
	r := make([]stockTimestampRange, 0, 2)
	now := core.NewTimestampTypeFromTime(time.Now())
	postStartTime, err := db.NewestIntradayHistoryDate(ctx, stock, interval)
	if err != nil || postStartTime.IsZero() || rand.Intn(25) == 0 {
		// if there was an error or we couldn't find the last date
		r = append(r, stockTimestampRange{
			StartDate: core.TimestampType{},
			EndDate:   now,
		})
	} else {
		// we assume what ever data is already there is contiguous
		// there we find the first and last days
		// we take the last day and today and search that
		// and if the defaultIntradyStart is before the first day then
		// we'll search that too

		precededStartTime, _ := db.OldestIntradayHistoryDate(ctx, stock, interval)

		//we want to check the range before the first intraday info we have in the DB it
		// it's after the default intraday history goal
		if now.Add(-defaultIntradyStart).Before(precededStartTime) {
			//we make the range from way back to the start of the data(plus a week for good measure)
			r = append(r, stockTimestampRange{
				StartDate: now.Add(-defaultIntradyStart),
				EndDate:   core.NewTimestampTypeFromDate(precededStartTime.Date().AddDate(0, 0, 7)),
			})
		}
		// and we still want to check from the
		r = append(r, stockTimestampRange{
			StartDate: postStartTime.Truncate(24 * time.Hour),
			EndDate:   now,
		})
	}

	return &intradayStockRange{
		Stock:    stock,
		Delta:    defaultISRStartingDelta,
		Ranges:   r,
		Interval: interval,
		Retry:    maxISRRetry,
		Finished: finishedChan,
	}
}

type intradayStockRange struct {
	Stock    core.StockSymbolType
	Ranges   []stockTimestampRange
	Delta    time.Duration
	Retry    int
	Interval core.IntradayInterval
	Finished chan bool
}

type stockTimestampRange struct {
	StartDate core.TimestampType
	EndDate   core.TimestampType
}

func (isr *intradayStockRange) ReduceDelta() bool {
	//we round down to the nearest day
	isr.Delta = (isr.Delta / 2).Truncate(24 * time.Hour)
	//168h0m0s
	//72h0m0s
	//24h0m0s
	//0s
	return isr.Delta > 0
}

func (isr *intradayStockRange) ResetRetry() {
	isr.Retry = maxISRRetry
}

func processISR(ctx context.Context, isr *intradayStockRange, agent core.SibylAgent, db *database.SibylDatabase, stockCache *StockCache) {
	//for each isr we try and to take a
	// chunk of time starting at the EndDate-Delta
	// each failed attempt will reduce retry
	// when reduceDelta() is called
	// if reduceDelta() returns true then the attempts restart.
	// if it returns false

	for _, r := range isr.Ranges {
		for r.StartDate.Before(r.EndDate) {
			//there is the chance we're in this loop and the program is wanting to exit
			// so we need to check the ctx in case it's done
			select {
			case <-ctx.Done():
				return
			default:
			}

			startDate := r.EndDate.Add(-isr.Delta)

			//we asked for everything between startDate and endDate
			// let's make sure we stay in inside those bounds
			if startDate.Before(r.StartDate) {
				startDate = r.StartDate
			}
			intradayRecords, err := agent.GetIntraday(ctx, isr.Stock, isr.Interval, startDate, r.EndDate)
			if err != nil {
				logrus.Errorf("IntradayGrabber: had a problem getting %v Intraday data on stock %v for range (%v - %v): %v", isr.Interval, isr.Stock, startDate, r.EndDate, err)
				isr.Retry--
				if isr.Retry == 0 {
					isr.ResetRetry()
					if !isr.ReduceDelta() {
						//well we can't reduce any more
						isr.Finished <- false
						logrus.Errorf("IntradayGrabber: unable to reduce delta (%v - %v) finished getting %v Intraday history for %v, after error: %v", startDate, r.EndDate, isr.Interval, isr.Stock, err)
						return
					}
				}
				continue
			} else {

				if len(intradayRecords) == 0 {
					//we either had an error or there wasn't anything to download
					// either way we're done no need to go any further back in time
					// there be another random chance to try again later
					logrus.Debugf("IntradayGrabber: received zero records, assuming the last of the %v Intraday history for %v is found at: %v - %v", isr.Interval, isr.Stock, startDate, r.EndDate)
					break
				}

				if err = db.LoadIntradayRecords(ctx, intradayRecords); err != nil {
					logrus.Errorf("IntradayGrabber: had a problem saving Intraday Record data: %v", err)
					isr.Retry--
					if isr.Retry == 0 {
						isr.ResetRetry()
						if !isr.ReduceDelta() {
							//well we can't reduce any more
							isr.Finished <- false
							logrus.Errorf("IntradayGrabber: loading data failed on dates(%v - %v) finished getting %v Intraday history for %v, after error: %v", startDate, r.EndDate, isr.Interval, isr.Stock, err)
							return
						}
					}
					continue
				}
			}
			r.EndDate = startDate
			isr.ResetRetry()
		}
	}

	//now we need to update the stock info but we don't care it has a problem
	switch isr.Interval {
	case core.TickInterval:
		stockCache.UpdateIntradayTimestampTick(isr.Stock, core.NewTimestampTypeFromTime(time.Now()))
	case core.OneMinInterval:
		stockCache.UpdateIntradayTimestamp1Min(isr.Stock, core.NewTimestampTypeFromTime(time.Now()))
	case core.FiveMinInterval:
		stockCache.UpdateIntradayTimestamp5Min(isr.Stock, core.NewTimestampTypeFromTime(time.Now()))
	}

	isr.Finished <- true
	logrus.Infof("IntradayGrabber: finished getting %v Intraday history (%v) for %v", isr.Interval, isr.Ranges, isr.Stock)
}

func (ig *IntradayGrabber) Stop(waitUpTo time.Duration) {
	//next stop the quoter
	logrus.Infof("Waiting for IntradayGrabber to finish")
	startTime := time.Now()
	ig.kill()
	select {
	case <-ig.doneCtx.Done():
		logrus.Infof("IntradayGrabber finished in %v", time.Since(startTime))
	case <-time.After(waitUpTo):
		logrus.Errorf("IntradayGrabber failed to gracefully finish in %v", time.Since(startTime))
	}
}
