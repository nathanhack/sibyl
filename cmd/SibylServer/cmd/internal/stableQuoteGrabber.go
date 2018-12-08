package internal

import (
	"context"
	"fmt"
	"github.com/nathanhack/sibyl/core/database"
	"github.com/sirupsen/logrus"
	"time"
)

//TODO Change stablequote to allow for multiple updates taking the newest data but deal with NULLs from the agents in a smart way

type StableQuoteGrabber struct {
	killCtx     context.Context
	kill        context.CancelFunc
	doneCtx     context.Context
	done        context.CancelFunc
	db          *database.SibylDatabase
	symbolCache *SymbolCache
	running     bool
}

func NewStableQuoteGrabber(db *database.SibylDatabase, symbolCache *SymbolCache) *StableQuoteGrabber {
	killCtx, kill := context.WithCancel(context.Background())
	doneCtx, done := context.WithCancel(context.Background())
	return &StableQuoteGrabber{
		killCtx:     killCtx,
		kill:        kill,
		doneCtx:     doneCtx,
		done:        done,
		db:          db,
		symbolCache: symbolCache,
	}
}

func (sqg *StableQuoteGrabber) Run() error {
	if sqg.running {
		return fmt.Errorf("StableQuoteGrabber is already running")
	}

	sqg.running = true
	go func(sqg *StableQuoteGrabber) {
		durationToWait := 16 * time.Second
		runGrabber := make(chan bool, 100)
	mainLoop:
		for {
			select {
			case <-sqg.killCtx.Done():
				break mainLoop
			case <-time.After(durationToWait):
				select {
				//non-blocking add
				case runGrabber <- true:
				default:
				}
			case <-sqg.symbolCache.StableQuoteStockSymbolsChanged:
				//this is a signal from the cache that we've had an update
				select {
				//non-blocking add
				case runGrabber <- true:
				default:
				}
			case <-sqg.symbolCache.StableQuoteOptionSymbolsChanged:
				select {
				//non-blocking add
				case runGrabber <- true:
				default:
				}
			case <-runGrabber:
				//now schedule the to do this again in 4 hrs there is the chance we didn't get everything so
				// this will give us 6 trys per day
				durationToWait = 4 * time.Hour
				startTime := time.Now()
				if startTime.Weekday() != time.Saturday && startTime.Weekday() != time.Sunday {
					//update the options before continuing
					agent, err := sqg.db.GetAgent(sqg.killCtx)
					if err != nil {
						logrus.Errorf("StableQuoteGrabber: could not retrieve agent: %v", err)
						continue
					}

					//get the current list of symbols from the cache
					sqg.symbolCache.StableQuoteStockSymbolsMu.RLock()
					stableStockQuoteSymbolsToDownLoad := sqg.symbolCache.StableQuoteStockSymbols
					sqg.symbolCache.StableQuoteStockSymbolsMu.RUnlock()

					sqg.symbolCache.StableQuoteOptionSymbolsMu.RLock()
					stableOptionQuoteSymbolsToDownLoad := sqg.symbolCache.StableQuoteOptionSymbols
					sqg.symbolCache.StableQuoteOptionSymbolsMu.RUnlock()

					if len(stableStockQuoteSymbolsToDownLoad) > 0 || len(stableOptionQuoteSymbolsToDownLoad) > 0 {
						stocks, options, err := agent.GetStableQuotes(sqg.killCtx, stableStockQuoteSymbolsToDownLoad, stableOptionQuoteSymbolsToDownLoad)
						errorsFound := false
						if err != nil {
							logrus.Errorf("StableQuoteGrabber: had a problem getting Stable Quotes: %v", err)
							errorsFound = true
						}

						if err = sqg.db.LoadStableStockQuoteRecords(sqg.killCtx, stocks); err != nil {
							logrus.Errorf("StableQuoteGrabber: had a problem saving Stable Stock Quotes: %v", err)
							errorsFound = true
						}

						if err = sqg.db.LoadStableOptionQuoteRecords(sqg.killCtx, options); err != nil {
							logrus.Errorf("StableQuoteGrabber: had a problem saving Stable Option Quotes: %v", err)
							errorsFound = true
						}

						if errorsFound {
							//then we want to try again until we get a clean run
							// so we'll reduce the durationToWait
							durationToWait = 1 * time.Minute
						}
					}
				}
				logrus.Infof("StableQuoteGrabber: finished a round in %v, next round in %v", time.Since(startTime), durationToWait)
			}
		}
		sqg.done() //signal this is finished
	}(sqg)
	return nil
}

func (sqg *StableQuoteGrabber) Stop(waitUpTo time.Duration) {
	//next stop the quoter
	logrus.Infof("Waiting for StableQuoteGrabber to finish")
	startTime := time.Now()
	sqg.kill()
	select {
	case <-sqg.doneCtx.Done():
		logrus.Infof("StableQuoteGrabber finished in %v", time.Since(startTime))
	case <-time.After(waitUpTo):
		logrus.Errorf("StableQuoteGrabber failed to gracefully finish in %v", time.Since(startTime))
	}
}
