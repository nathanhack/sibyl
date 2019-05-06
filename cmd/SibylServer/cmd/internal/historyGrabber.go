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
	defaultHistoryStartYears     = 20 //years
	defaultHistoryDurationToWait = 1 * time.Second
	maxHSRRetry                  = 1
)

type HistoryGrabber struct {
	killCtx    context.Context
	kill       context.CancelFunc
	doneCtx    context.Context
	done       context.CancelFunc
	db         *database.SibylDatabase
	stockCache *StockCache
	running    bool
}

func NewHistoryGrabber(db *database.SibylDatabase, symbolCache *StockCache) *HistoryGrabber {
	killCtx, kill := context.WithCancel(context.Background())
	doneCtx, done := context.WithCancel(context.Background())
	return &HistoryGrabber{
		killCtx:    killCtx,
		kill:       kill,
		doneCtx:    doneCtx,
		done:       done,
		db:         db,
		stockCache: symbolCache,
	}
}

func (hg *HistoryGrabber) Run() error {
	if hg.running {
		return fmt.Errorf("HistoryGrabber is already running")
	}

	hg.running = true
	go func(hg *HistoryGrabber) {
		onceDailyDeadline, onceDailyDeadlineCancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
		durationToWait := defaultHistoryDurationToWait
		updateFullStockDailyCache := make(map[core.StockSymbolType]*historyStockRange)
		updateFullStockWeeklyCache := make(map[core.StockSymbolType]*historyStockRange)
		updateFullStockMonthlyCache := make(map[core.StockSymbolType]*historyStockRange)
		updateFullStockYearlyCache := make(map[core.StockSymbolType]*historyStockRange)
		finishedChan := make(chan bool, 100)
	mainLoop:
		for {
			select {
			case <-hg.killCtx.Done():
				break mainLoop
			case <-onceDailyDeadline.Done():
				onceDailyDeadlineCancel() // (REQUIRED) we call the cancel func to release resources associated with the context
				onceDailyDeadline, onceDailyDeadlineCancel = context.WithDeadline(context.Background(), nextValidHistoryTime())

				//here we go through all valid stocks and do a full update
				updateFullStockDailyCache = addToHistoryCache(hg.killCtx, hg.db, hg.stockCache, finishedChan, updateFullStockDailyCache, core.DailyInterval, false, true)
				updateFullStockWeeklyCache = addToHistoryCache(hg.killCtx, hg.db, hg.stockCache, finishedChan, updateFullStockWeeklyCache, core.WeeklyInterval, false, true)
				updateFullStockMonthlyCache = addToHistoryCache(hg.killCtx, hg.db, hg.stockCache, finishedChan, updateFullStockMonthlyCache, core.MonthlyInterval, false, true)
				updateFullStockYearlyCache = addToHistoryCache(hg.killCtx, hg.db, hg.stockCache, finishedChan, updateFullStockYearlyCache, core.YearlyInterval, false, true)

			case <-time.After(durationToWait):
				currentTime := time.Now()
				durationToWait = defaultHistoryDurationToWait
				agent, err := hg.db.GetAgent(hg.killCtx)
				if err != nil {
					logrus.Errorf("HistoryGrabber: could not retrieve agent: %v", err)
					//we extend the wait time for a round
					durationToWait = 5 * time.Minute
					continue
				}

				//first we do a check for anything new to the database
				updateFullStockDailyCache = addToHistoryCache(hg.killCtx, hg.db, hg.stockCache, finishedChan, updateFullStockDailyCache, core.DailyInterval, true, false)
				updateFullStockWeeklyCache = addToHistoryCache(hg.killCtx, hg.db, hg.stockCache, finishedChan, updateFullStockWeeklyCache, core.WeeklyInterval, true, false)
				updateFullStockMonthlyCache = addToHistoryCache(hg.killCtx, hg.db, hg.stockCache, finishedChan, updateFullStockMonthlyCache, core.MonthlyInterval, true, false)
				updateFullStockYearlyCache = addToHistoryCache(hg.killCtx, hg.db, hg.stockCache, finishedChan, updateFullStockYearlyCache, core.YearlyInterval, true, false)

				//and now we check if there is anything to do if not we bail on this round
				if len(updateFullStockDailyCache) == 0 && len(updateFullStockWeeklyCache) == 0 &&
					len(updateFullStockMonthlyCache) == 0 && len(updateFullStockYearlyCache) == 0 {
					continue
				}

				runningCount := 0
				ctx, cancel := context.WithCancel(hg.killCtx)
				//now we run each cache item in it's own routine
				// we start with the FullCaches first
			updateFullDailyCacheLoop:
				for _, hsr := range updateFullStockDailyCache {
					//do a quick check if we've stopped the application
					if areWeDone(hg.killCtx) {
						break updateFullDailyCacheLoop
					}

					runningCount++
					go processHSR(ctx, hg.db, hg.stockCache, agent, hsr)
				}
				updateFullStockDailyCache = make(map[core.StockSymbolType]*historyStockRange)

			updateFullWeeklyCacheLoop:
				for _, hsr := range updateFullStockWeeklyCache {
					//do a quick check if we've stopped the application
					if areWeDone(hg.killCtx) {
						break updateFullWeeklyCacheLoop
					}

					runningCount++
					go processHSR(ctx, hg.db, hg.stockCache, agent, hsr)
				}
				updateFullStockWeeklyCache = make(map[core.StockSymbolType]*historyStockRange)

			updateFullMonthlyCacheLoop:
				for _, hsr := range updateFullStockMonthlyCache {
					//do a quick check if we've stopped the application
					if areWeDone(hg.killCtx) {
						break updateFullMonthlyCacheLoop
					}

					runningCount++
					go processHSR(ctx, hg.db, hg.stockCache, agent, hsr)
				}
				updateFullStockMonthlyCache = make(map[core.StockSymbolType]*historyStockRange)

			updateFullYearlyCacheLoop:
				for _, hsr := range updateFullStockYearlyCache {
					//do a quick check if we've stopped the application
					if areWeDone(hg.killCtx) {
						break updateFullYearlyCacheLoop
					}

					runningCount++
					go processHSR(ctx, hg.db, hg.stockCache, agent, hsr)
				}
				updateFullStockYearlyCache = make(map[core.StockSymbolType]*historyStockRange)

				//now we loop through finishedChan and break once we have all the results
			finishedLoop:
				for runningCount > 0 {
					select {
					case <-hg.killCtx.Done():
						break finishedLoop
					case <-finishedChan:
						runningCount--
					case <-time.After(30 * time.Minute):
						//a fail safe this should take 30 minutes to complete
						logrus.Errorf("HistoryGrabber: had an issue getting all the result in a timely manner")
						break finishedLoop
					}
				}
				cancel() //kill any of the go routines still running for this round
				logrus.Infof("HistoryGrabber: finished a round in %v", time.Since(currentTime))
			}
		}
		hg.done() //signal this is finished
	}(hg)
	return nil
}

func nextValidHistoryTime() time.Time {
	tomorrowOrNextWeekday := core.NewDateTypeFromTime(time.Now()).AddDate(0, 0, 1)
	for !tomorrowOrNextWeekday.IsWeekDay() {
		//since we need at least two days for history quote
		tomorrowOrNextWeekday = tomorrowOrNextWeekday.AddDate(0, 0, 1)
	}
	return tomorrowOrNextWeekday.Time().Add(6 * time.Hour)
}

func addToHistoryCache(ctx context.Context, db *database.SibylDatabase,
	symbolCache *StockCache, finishedChan chan bool,
	oldCache map[core.StockSymbolType]*historyStockRange,
	interval core.HistoryInterval, onlyNewlyAdded bool, fullHistoryUpdate bool) map[core.StockSymbolType]*historyStockRange {

	toReturn := make(map[core.StockSymbolType]*historyStockRange)

	today := core.NewDateTypeFromTime(time.Now())

	for _, stock := range symbolCache.HistoryStocks(interval == core.DailyInterval,
		interval == core.WeeklyInterval, interval == core.MonthlyInterval, interval == core.YearlyInterval) {
		if areWeDone(ctx) {
			return toReturn //if we got the context done we end early
		}

		if fullHistoryUpdate {
			if x, has := oldCache[stock.Symbol]; !has {
				toReturn[stock.Symbol] = makeHSRFull(ctx, db, finishedChan, stock.Symbol, interval)
			} else {
				toReturn[stock.Symbol] = x
			}
		} else {
			// if either we're not current or if we're looking for onlyNewlyAdded
			if (!onlyNewlyAdded && stock.HistoryTimestamp.Before(today)) || (onlyNewlyAdded && stock.HistoryTimestamp.IsZero()) {
				if x, has := oldCache[stock.Symbol]; !has {
					toReturn[stock.Symbol] = makeHSRLatestOnly(ctx, db, finishedChan, stock.Symbol, interval)
				} else {
					toReturn[stock.Symbol] = x
				}
			}
		}
	}

	return toReturn
}

func makeHSRLatestOnly(ctx context.Context, db *database.SibylDatabase, finishedChan chan bool, stock core.StockSymbolType, interval core.HistoryInterval) *historyStockRange {
	now := core.NewDateTypeFromTime(time.Now())
	r := make([]stockDatestampRange, 0)
	postStartTime, err := db.NewestHistoryDate(ctx, stock, interval)
	if err == nil {
		//we grab the postStartTime to now. if it's not in the DB then it will
		//return a zero which will be get it all
		r = append(r, stockDatestampRange{StartDate: postStartTime, EndDate: now})
	}

	return &historyStockRange{
		Stock:    stock,
		Delta:    defaultHistoryStartYears * 365,
		Ranges:   r,
		Retry:    maxISRRetry,
		Interval: interval,
		Finished: finishedChan,
	}
}

func makeHSRFull(ctx context.Context, db *database.SibylDatabase, finishedChan chan bool, stock core.StockSymbolType, interval core.HistoryInterval) *historyStockRange {
	r := make([]stockDatestampRange, 0, 2)
	now := core.NewDateTypeFromTime(time.Now())
	postStartTime, err := db.NewestHistoryDate(ctx, stock, interval)
	if err != nil || postStartTime.IsZero() || rand.Intn(25) == 0 {
		// if there was an error or we couldn't find the last date
		r = append(r, stockDatestampRange{
			StartDate: core.DateType{},
			EndDate:   now,
		})
	} else {
		// we assume what ever data is already there is contiguous
		// there we find the first and last days
		// we take the last day and today and search that
		// and if the defaultIntradyStart is before the first day then
		// we'll search that too

		precededStartTime, _ := db.OldestHistoryDate(ctx, stock, interval)

		//we want to check the range before the first intraday info we have in the DB it
		// it's after the default intraday history goal
		if now.AddDate(-defaultHistoryStartYears, 0, 0).Before(precededStartTime) {
			//we make the range from way back to the start of the data(plus a week for good measure)
			r = append(r, stockDatestampRange{
				StartDate: now.AddDate(-defaultHistoryStartYears, 0, 0),
				EndDate:   precededStartTime.AddDate(0, 0, 7),
			})
		}
		// and we still want to check from the (we do -7 days for cushion)
		r = append(r, stockDatestampRange{
			StartDate: postStartTime.AddDate(0, 0, -7),
			EndDate:   now,
		})
	}

	return &historyStockRange{
		Stock:    stock,
		Delta:    defaultHistoryStartYears * 365,
		Ranges:   r,
		Interval: interval,
		Retry:    maxISRRetry,
		Finished: finishedChan,
	}
}

type historyStockRange struct {
	Stock    core.StockSymbolType
	Ranges   []stockDatestampRange
	Delta    int //current length of time download per query
	Retry    int
	Interval core.HistoryInterval
	Finished chan bool
}

func (hsr *historyStockRange) ReduceDelta() bool {
	hsr.Delta /= 2
	return hsr.Delta > 2
}

func (hsr *historyStockRange) ResetRetry() {
	hsr.Retry = maxHSRRetry
}

type stockDatestampRange struct {
	StartDate core.DateType
	EndDate   core.DateType
}

func processHSR(ctx context.Context, db *database.SibylDatabase, stockCache *StockCache, agent core.SibylAgent, hsr *historyStockRange) {
	//for each hsr we try and to take a
	// chunk of time starting at the EndDate-Delta
	// each failed attempt will reduce retry
	// when reduceDelta() is called
	// if reduceDelta() returns true then the attempts restart.
	// if it returns false the stock is marked to "scanned" and move to the next isr

	for _, r := range hsr.Ranges {
		for r.StartDate.Before(r.EndDate) {
			//there is the chance we're in this loop and the program is wanting to exit
			// so we need to check the ctx in case it's done
			if areWeDone(ctx) {
				return
			}

			startDate := r.EndDate.AddDate(0, 0, -hsr.Delta)

			//we asked for everything between startDate and endDate
			// let's make sure we stay in inside those bounds
			if startDate.Before(r.StartDate) {
				startDate = r.StartDate
			}
			historyRecords, err := agent.GetHistory(ctx, hsr.Stock, hsr.Interval, startDate, r.EndDate)
			if err != nil {
				logrus.Errorf("HistoryGrabber: had a problem getting History data on stock %v for range (%v - %v): %v", hsr.Stock, startDate, r.EndDate, err)
				hsr.Retry--
				if hsr.Retry == 0 {
					hsr.ResetRetry()
					if !hsr.ReduceDelta() {
						//well we can't reduce any more
						hsr.Finished <- false
						logrus.Infof("HistoryGrabber: unable to reduce delta (%v-%v) finished getting History for %v, after error: %v", startDate, r.EndDate, hsr.Stock, err)
						return
					}
				}
				continue
			} else {

				if len(historyRecords) == 0 {
					//we either had an error or there wasn't anything to download
					// either way we're done no need to go any further back in time
					// there be another random chance to try again later
					logrus.Debugf("HistoryGrabber: received zero records, assuming the last of the history for %v is found at: %v - %v", hsr.Stock, startDate, r.EndDate)
					break
				}

				if err = db.LoadHistoryRecords(ctx, historyRecords); err != nil {
					logrus.Errorf("HistoryGrabber: had a problem saving History Record data: %v", err)
					hsr.Retry--
					if hsr.Retry == 0 {
						hsr.ResetRetry()
						if !hsr.ReduceDelta() {
							//well we can't reduce any more
							hsr.Finished <- false
							logrus.Infof("HistoryGrabber: loading data failed on dates(%v - %v) finished getting History for %v, after error: %v", startDate, r.EndDate, hsr.Stock, err)
							return
						}
					}
					continue
				}
			}
			r.EndDate = startDate
			hsr.ResetRetry()
		}
	}

	//now we need to update the stock info but we don't care it has a problem
	stockCache.UpdateHistoryTimestamp(hsr.Stock, core.NewDateTypeFromTime(time.Now()))

	hsr.Finished <- true
	logrus.Infof("HistoryGrabber: finished getting %v History for %v", hsr.Interval, hsr.Stock)
}

func (hg *HistoryGrabber) Stop(waitUpTo time.Duration) {
	//next stop the quoter
	logrus.Infof("Waiting for HistoryGrabber to finish")
	startTime := time.Now()
	hg.kill()
	select {
	case <-hg.doneCtx.Done():
		logrus.Infof("HistoryGrabber finished in %v", time.Since(startTime))
	case <-time.After(waitUpTo):
		logrus.Errorf("HistoryGrabber failed to gracefully finish in %v", time.Since(startTime))
	}
}
