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
	maxISRRetry             = 3
	defaultISRStartingDelta = 14 * 24 * time.Hour  // 14 days
	defaultIntradyStart     = 5 * 8760 * time.Hour //5 years - tends to only be 1 year available from brokers //TODO consider making this a configuration
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
		runIntradyGrabber := make(chan bool, 1)
	mainLoop:
		for {
			select {
			case <-ig.killCtx.Done():
				break mainLoop
			case <-ig.symbolCache.IntradaySymbolsChanged:
				select {
				//non-blocking add if fails it's not a problem
				case runIntradyGrabber <- true:
				default:
				}
			case <-time.After(durationToWait):
				select {
				//non-blocking add if fails it's not a problem
				case runIntradyGrabber <- true:
				default:
				}
			case <-runIntradyGrabber:
				currentTime := time.Now()
				agent, err := ig.db.GetAgent(ig.killCtx)
				if err != nil {
					logrus.Errorf("IntradayGrabber: could not retrieve agent: %v", err)
					//we want to run again but not immediately
					durationToWait = 5 * time.Minute
					continue
				}

				//basically we take all current stocks with intraday history enable
				// and the intraday values for it
				// we focus on getting values from the last intraday values until today
				// for a small fraction we'll try and get all the history available
				// if there isn't any history then we try and get all available history

				//schedule the next update to be tomorrow morning at 6am
				year, month, day := currentTime.Date()
				durationToWait = time.Date(year, month, day+1, 6, 0, 0, 0, time.Local).Sub(currentTime)

				//get the current list of stocks that have intraday history enabled
				ig.symbolCache.IntradaySymbolsMu.RLock()
				intradaySymbolsToDownload := ig.symbolCache.IntradaySymbols
				ig.symbolCache.IntradaySymbolsMu.RUnlock()

				if len(intradaySymbolsToDownload) == 0 {
					//if there wasn't anything in there to look at
					// we just skip this round
				}

				emptyTime := core.NewTimestampTypeFromUnix(0)

				//for each stock send out an process to get the intraday values
				finishedChan := make(chan bool, 10)
				runningCount := 0
				ctx, cancel := context.WithCancel(ig.killCtx)
				for stock := range intradaySymbolsToDownload {
					//first we do a quick check if we've stopped the application
					select {
					case <-ig.killCtx.Done():
						break
					default:
					}
					startTime, err := ig.db.LastIntradayHistoryDate(ig.killCtx, stock)
					if err != nil || startTime == emptyTime {
						// if there was an error or we couldn't find the last date
						isr := intradayStockRange{
							Stock:     stock,
							Delta:     defaultISRStartingDelta,
							StartDate: time.Now().Truncate(24 * time.Hour).Add(-defaultIntradyStart),
							EndDate:   time.Now().Truncate(24 * time.Hour),
							Retry:     maxISRRetry,
							Finished:  finishedChan,
						}
						runningCount++
						go processISR(ctx, &isr, agent, ig.db)

					} else if rand.Intn(25) == 0 {
						//if we randomly picked it we'll try and get all the data
						// we assume what ever data is already there is contiguous
						// there we find the first and last days
						// we take the last day and today and search that
						// and if the defaultIntradyStart is before the first day then
						// we'll search that too

						precededStartTime, _ := ig.db.FirstIntradayHistoryDate(ig.killCtx, stock)

						//here the error doesn't matter startTime will be zero it happened
						if time.Now().Truncate(24 * time.Hour).Add(-defaultIntradyStart).Before(precededStartTime.Time()) {
							isr := intradayStockRange{
								Stock:     stock,
								Delta:     defaultISRStartingDelta,
								StartDate: time.Now().Truncate(24 * time.Hour).Add(-defaultIntradyStart),
								EndDate:   precededStartTime.Time().Truncate(24*time.Hour).AddDate(0, 0, 7), //plus one weeks
								Retry:     maxISRRetry,
								Finished:  finishedChan,
							}
							runningCount++
							go processISR(ctx, &isr, agent, ig.db)
						}
						isr := intradayStockRange{
							Stock:     stock,
							Delta:     defaultISRStartingDelta,
							StartDate: startTime.Time().Truncate(24 * time.Hour),
							EndDate:   time.Now().Truncate(24 * time.Hour),
							Retry:     maxISRRetry,
							Finished:  finishedChan,
						}
						runningCount++
						go processISR(ctx, &isr, agent, ig.db)

					} else {
						isr := intradayStockRange{
							Stock:     stock,
							Delta:     defaultISRStartingDelta,
							StartDate: startTime.Time().Truncate(24 * time.Hour),
							EndDate:   time.Now().Truncate(24 * time.Hour),
							Retry:     maxISRRetry,
							Finished:  finishedChan,
						}
						runningCount++
						go processISR(ctx, &isr, agent, ig.db)
					}
				}

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
						logrus.Errorf("IntradayGrabber: had an issue getting all the result in a timely manor")
						cancel() //kill any of the go routines still running for this round
						break finishedLoop
					}
				}

				logrus.Infof("IntradayGrabber: finished a round in %v next round to begin in %v", time.Since(currentTime), durationToWait)
			}
		}
		ig.done() //signal this is finished
	}(ig)
	return nil
}

type intradayStockRange struct {
	Stock     core.StockSymbolType
	StartDate time.Time
	EndDate   time.Time
	Delta     time.Duration
	Retry     int
	Finished  chan bool
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

func processISR(ctx context.Context, isr *intradayStockRange, agent core.SibylAgent, db *database.SibylDatabase) {
	//for each isr we try and to take a
	// chunk of time starting at the EndDate-Delta
	// each failed attempt will reduce retry
	// when retry is equal to zero, reduceDelta() is called
	// if reduceDelta() returns true then the attempts restart.
	// if it returns false the stock is marked to "scanned" and move to the next isr

	//but first we update the database to indicate that we're scanning this stock
	if err := db.StockIntradayHistorySetScanning(ctx, isr.Stock); err != nil {
		//if this failed that's bad BUT we'll continue and hope everything was sorted out
		logrus.Errorf("IntradayGrabber: had a problem setting stock %v to \"scanning\": %v", isr.Stock, err)
	}

	for isr.StartDate.Before(isr.EndDate) {
		//there is the chance we're in this loop and the program is wanting to exit
		// so we need to check the ctx in case it's done
		select {
		case <-ctx.Done():
			return
		default:
		}

		startDate := isr.EndDate.Add(-isr.Delta)

		//we asked for everything between startDate and endDate
		// let's make sure we stay in inside those bounds
		if startDate.Before(isr.StartDate) {
			startDate = isr.StartDate
		}
		intradayRecords, err := agent.GetIntraday(ctx, isr.Stock, core.MinuteTicks, core.NewTimestampTypeFromTime(startDate), core.NewTimestampTypeFromTime(isr.EndDate))
		if err != nil {
			logrus.Errorf("IntradayGrabber: had a problem getting Intraday data on stock %v for range (%v-%v): %v", isr.Stock, startDate, isr.EndDate, err)
			isr.Retry--
			if isr.Retry == 0 {
				isr.ResetRetry()
				if !isr.ReduceDelta() {
					//well we can't reduce any more
					// set the stock value to scanned
					if err1 := db.StockIntradayHistorySetScanned(ctx, isr.Stock); err1 != nil {
						logrus.Errorf("IntradayGrabber: had a problem setting stock %v to \"scanned\": %v : after this error: %v", isr.Stock, err1, err)
					}
					isr.Finished <- false
					logrus.Infof("IntradayGrabber: unable to reduce delta (%v-%v) finished getting Intraday history for %v, after error: %v", startDate, isr.EndDate, isr.Stock, err)
					return
				}
			}
			continue
		} else {

			if len(intradayRecords) == 0 {
				//we either had an error or there wasn't anything to download
				// either way we're done no need to go any further back in time
				// there be another random chance to try again later
				logrus.Debugf("IntradayGrabber: received zero records, assuming the last of the history for %v is found at: %v - %v", isr.Stock, startDate, isr.EndDate)
				break
			}

			if err = db.LoadIntradayRecords(ctx, intradayRecords); err != nil {
				logrus.Errorf("IntradayGrabber: had a problem saving Intraday Record data: %v", err)
				isr.Retry--
				if isr.Retry == 0 {
					isr.ResetRetry()
					if !isr.ReduceDelta() {
						//well we can't reduce any more
						// set the stock value to scanned
						if err1 := db.StockIntradayHistorySetScanned(ctx, isr.Stock); err1 != nil {
							logrus.Errorf("IntradayGrabber: had a problem setting stock %v to \"scanned\": %v : after this error: %v", isr.Stock, err1, err)
						}
						isr.Finished <- false
						logrus.Infof("IntradayGrabber: loading data failed on dates(%v - %v) finished getting Intraday history for %v, after error: %v", startDate, isr.EndDate, isr.Stock, err)
						return
					}
				}
				continue
			}
		}
		isr.EndDate = startDate
	}

	if err := db.StockIntradayHistorySetScanned(ctx, isr.Stock); err != nil {
		logrus.Errorf("IntradayGrabber: had a problem setting stock %v to \"scanned\": %v", isr.Stock, err)
	}
	isr.Finished <- true
	logrus.Infof("IntradayGrabber: finished getting Intraday history for %v", isr.Stock)
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
