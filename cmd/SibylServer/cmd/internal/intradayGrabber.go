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

type IntradayGrabber struct {
	killCtx     context.Context
	kill        context.CancelFunc
	doneCtx     context.Context
	done        context.CancelFunc
	db          *database.SibylDatabase
	symbolCache *SymbolCache
	running     bool
}

func NewIntradayGrabber(db *database.SibylDatabase, symbolCache *SymbolCache) *IntradayGrabber {
	killCtx, kill := context.WithCancel(context.Background())
	doneCtx, done := context.WithCancel(context.Background())
	return &IntradayGrabber{
		killCtx:     killCtx,
		kill:        kill,
		doneCtx:     doneCtx,
		done:        done,
		db:          db,
		symbolCache: symbolCache,
	}
}

func (ig *IntradayGrabber) Run() error {
	if ig.running {
		return fmt.Errorf("IntradayGrabber is already running")
	}

	ig.running = true
	go func(ig *IntradayGrabber) {
		//we'll start this up with a short delay
		durationToWait := 15 * time.Second
		failedStocksAndDays := make(map[core.StockSymbolType]core.TimestampType)
		runIntradyGrabber := make(chan bool, 100)
	mainLoop:
		for {
			select {
			case <-ig.killCtx.Done():
				break mainLoop
			case <-ig.symbolCache.IntradaySymbolsChanged:
				runIntradyGrabber <- true
			case <-time.After(durationToWait):
				runIntradyGrabber <- true
			case <-runIntradyGrabber:
				//first we drain the chan
			drainLoop:
				for {
					select {
					case <-runIntradyGrabber:
						continue
					default:
						break drainLoop
					}
				}

				currentTime := time.Now()
				agent, err := ig.db.GetAgent(ig.killCtx)
				if err != nil {
					logrus.Errorf("IntradayGrabber: could not retrieve agent: %v", err)
					//we want to run again but not immediately
					durationToWait = 5 * time.Minute
					continue
				}

				//schedule the next update to be tomorrow morning at 6am
				year, month, day := currentTime.Date()
				durationToWait = time.Date(year, month, day+1, 6, 0, 0, 0, time.Local).Sub(currentTime)

				ig.symbolCache.IntradaySymbolsMu.RLock()
				intradaySymbolsToDownload := ig.symbolCache.IntradaySymbols
				ig.symbolCache.IntradaySymbolsMu.RUnlock()

				emptyTime := core.NewTimestampTypeFromUnix(0)
				endtime := core.NewTimestampTypeFromTime(time.Now())
				stockAndDays := make(map[core.StockSymbolType]core.TimestampType)
				//first we make a map of stocks and the number of days we need to get for each one upto yesterdayOrLastWeekday
				// intraday should also be upto 1600 however, sometimes the there are missing values so we'll be an earlier time just in case
				yesterdayOrLastWeekday := core.NewTimestampTypeFromTime(time.Date(year, month, day-1, 12, 0, 0, 0, time.Local))
				for yesterdayOrLastWeekday.IsWeekDay() {
					yesterdayOrLastWeekday = yesterdayOrLastWeekday.AddDate(0, 0, -1)
				}
				for stock := range intradaySymbolsToDownload {
					//TODO make the intradayhistory time a configuration
					startTime, err := ig.db.LastIntradayHistoryDate(ig.killCtx, stock)
					if err != nil ||
						startTime == emptyTime ||
						(len(failedStocksAndDays) == 0 && rand.Intn(2) == 0) {
						// if we couldn't find the last date or the date was zero or if we randomly picked it we'll
						// want the max number of days of intradays data which is either 5 or 10 days (as defined by the services web sites)
						// HOWEVER, depending on discount broker .. it can b e 45 - 60 days depending on the stock
						// so we'll start high and over time roll down on error
						if _, has := failedStocksAndDays[stock]; has {
							//if it was in the failedStockAndDays it means we want to move the date up by one day
							nTime := failedStocksAndDays[stock].AddDate(0, 0, 1)
							if nTime.Before(yesterdayOrLastWeekday) {
								stockAndDays[stock] = nTime
							} else {
								stockAndDays[stock] = yesterdayOrLastWeekday.AddDate(0, 0, -60)
							}
						} else {
							stockAndDays[stock] = yesterdayOrLastWeekday.AddDate(0, 0, -60)
						}

					} else {
						if startTime.Before(yesterdayOrLastWeekday) {
							stockAndDays[stock] = startTime.AddDate(0, 0, -1)
						}
					}
					if startTime, has := stockAndDays[stock]; has {
						logrus.Debugf("IntradayGrabber: Grabbing Intraday for %v startTime: %v and endTime: %v", stock, startTime, yesterdayOrLastWeekday)
					}
				}

				//we clear out the failed stocks from the previous run and so we can get
				failedStocksAndDays = make(map[core.StockSymbolType]core.TimestampType)

				//now we have a set of stock that need to be updated
				for stock, startTime := range stockAndDays {
					// so we'll wait in a select watching for killCtx
					select {
					case <-ig.killCtx.Done():
						logrus.Errorf("IntradayGrabber: context canceled")
						break mainLoop
					default:
					}

					intradayRecords, err := agent.GetIntraday(ig.killCtx, stock, core.MinuteTicks, startTime, endtime)
					if err != nil {
						logrus.Errorf("IntradayGrabber: had a problem getting Intraday data on %v: %v", stock, err)
						failedStocksAndDays[stock] = startTime
					} else {
						if err = ig.db.LoadIntradayRecords(ig.killCtx, intradayRecords); err != nil {
							logrus.Errorf("IntradayGrabber: had a problem saving Intraday Record data: %v", err)
							failedStocksAndDays[stock] = startTime
						}
					}
				}

				if len(failedStocksAndDays) > 0 {
					//if we had failures we start up right away
					runIntradyGrabber <- true
					logrus.Infof("IntradayGrabber: finished a round in %v some stocks were not updated, next round to being in %v", time.Since(currentTime), durationToWait)
				} else {
					//else everything was successful and we'll schedule the to do this again in 4 hrs
					// there is the chance we didn't get everything so this will give us 6 trys per day
					logrus.Infof("IntradayGrabber: finished a round in %v next round to begin in %v", time.Since(currentTime), durationToWait)
				}
			}
		}
		ig.done() //signal this is finished
	}(ig)
	return nil
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