package internal

import (
	"context"
	"fmt"
	"github.com/nathanhack/sibyl/core"
	"github.com/nathanhack/sibyl/core/database"
	"github.com/sirupsen/logrus"
	"time"
)

type OptionSymbolGrabber struct {
	killCtx       context.Context
	kill          context.CancelFunc
	doneCtx       context.Context
	done          context.CancelFunc
	db            *database.SibylDatabase
	symbolCache   *SymbolCache
	running       bool
	RequestUpdate chan bool
}

func NewOptionSymbolGrabber(db *database.SibylDatabase, symbolCache *SymbolCache) *OptionSymbolGrabber {
	killCtx, kill := context.WithCancel(context.Background())
	doneCtx, done := context.WithCancel(context.Background())
	return &OptionSymbolGrabber{
		killCtx:       killCtx,
		kill:          kill,
		doneCtx:       doneCtx,
		done:          done,
		db:            db,
		symbolCache:   symbolCache,
		RequestUpdate: make(chan bool, 100),
	}
}

func (osg *OptionSymbolGrabber) Run() error {
	if osg.running {
		return fmt.Errorf("OptionSymbolGrabber is already running")
	}
	//prime this to pull down the latest options symbols from the agent
	osg.running = true
	go func(osg *OptionSymbolGrabber) {
		// initially we wait a few seconds then start up
		durationToWait := 5 * time.Second

		///we create a retry timeout for when we need to retry
		waitForForever := 1000 * time.Hour
		wait1Min := 1 * time.Minute
		tryAgainIn := waitForForever
		failedSymbols := make(map[core.StockSymbolType]int)
		runGrabber := make(chan bool, 1)
	mainLoop:
		for {
			select {
			case <-osg.killCtx.Done():
				break mainLoop
			case <-time.After(durationToWait):
				//clean out any failed ones
				failedSymbols = make(map[core.StockSymbolType]int)
				tryAgainIn = waitForForever
				select {
				//non-blocking add
				case runGrabber <- true:
				default:
				}
			case <-time.After(tryAgainIn):
				tryAgainIn = waitForForever
				select {
				//non-blocking add
				case runGrabber <- true:
				default:
				}
			case <-osg.RequestUpdate:
				//first we drain the chan
			drainRequestUpdateLoop:
				for {
					select {
					case <-osg.RequestUpdate:
						continue
					default:
						break drainRequestUpdateLoop
					}
				}
				//clear all history
				failedSymbols = make(map[core.StockSymbolType]int)
				tryAgainIn = waitForForever
				select {
				//non-blocking add
				case runGrabber <- true:
				default:
				}
			case <-runGrabber:
				//now we start the actual task
				startTime := time.Now()
				agent, err := osg.db.GetAgent(osg.killCtx)
				if err != nil {
					logrus.Errorf("OptionSymbolGrabber: could not retrieve agent: %v", err)
					continue
				}

				//schedule the next update to be tomorrow morning at 6am
				year, month, day := startTime.Date()
				durationToWait = time.Date(year, month, day+1, 6, 0, 0, 0, time.Local).Sub(startTime)

				symbols := make([]core.StockSymbolType, 0)

				if len(failedSymbols) == 0 {
					stocks, err := osg.db.GetAllStockRecords(osg.killCtx)
					if err != nil {
						logrus.Errorf("OptionSymbolGrabber: had a problem getting list of stocks: %v", err)
						continue
					}
					for _, stock := range stocks {
						if stock.ValidationStatus == core.ValidationValid &&
							stock.DownloadStatus == core.ActivityEnabled &&
							stock.HasOptions {
							symbols = append(symbols, stock.Symbol)
						}
					}
				} else {
					for stock, count := range failedSymbols {
						if count < 4 {
							symbols = append(symbols, stock)
						}
					}
				}

				//for each stock in the list get updated options and put them in the database
				// we update the database so the quotes can use them (which may o
				for _, symbol := range symbols {
					select {
					case <-osg.killCtx.Done():
						logrus.Errorf("OptionSymbolGrabber: context canceled")
						break mainLoop
					default:
					}

					logrus.Infof("OptionSymbolGrabber: getting option symbols for %v", symbol)
					options, err := agent.GetStockOptionSymbols(osg.killCtx, symbol)
					if err != nil {
						//here we only log the error because there is a case of a weird options that are invalid
						logrus.Errorf("OptionSymbolGrabber: had an error while gather options for stock symbol:%v the error: %v", symbol, err)
					}
					if len(options) > 0 {
						if err := osg.db.SetOptionsForStock(osg.killCtx, symbol, options); err != nil {
							logrus.Errorf("OptionSymbolGrabber: failed during adding options to data, submitting to retry, error found: %v", err)
							//request an another update until we get a clean run
							if _, has := failedSymbols[symbol]; has {
								failedSymbols[symbol]++
							} else {
								failedSymbols[symbol] = 0
							}
							tryAgainIn = wait1Min
						}
					} else if len(options) == 0 {
						//there should have been something so we should do another request
						logrus.Errorf("OptionSymbolGrabber: found 0 option symbols for %v, submitting to retry", symbol)
						if _, has := failedSymbols[symbol]; has {
							failedSymbols[symbol]++
						} else {
							failedSymbols[symbol] = 0
						}
						tryAgainIn = wait1Min
					}
				}
				//now that we're done let the cache know
				osg.symbolCache.RequestUpdate <- true
				logrus.Infof("OptionSymbolGrabber: finished a round in %v", time.Since(startTime))
			}
		}
		osg.done() //signal this is finished
	}(osg)
	return nil
}

func (osg *OptionSymbolGrabber) Stop(waitUpTo time.Duration) {
	//next stop the quoter
	logrus.Infof("Waiting for OptionSymbolGrabber to finish")
	startTime := time.Now()
	osg.kill()
	select {
	case <-osg.doneCtx.Done():
		logrus.Infof("OptionSymbolGrabber finished in %v", time.Since(startTime))
	case <-time.After(waitUpTo):
		logrus.Errorf("OptionSymbolGrabber failed to gracefully finish in %v", time.Since(startTime))
	}
}
