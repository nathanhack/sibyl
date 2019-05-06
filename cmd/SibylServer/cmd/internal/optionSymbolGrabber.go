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
	stockCache    *StockCache
	running       bool
	RequestUpdate chan core.StockSymbolType //TODO consider moving this action into the SymbolsCache
}

func NewOptionSymbolGrabber(db *database.SibylDatabase, symbolCache *StockCache) *OptionSymbolGrabber {
	killCtx, kill := context.WithCancel(context.Background())
	doneCtx, done := context.WithCancel(context.Background())
	return &OptionSymbolGrabber{
		killCtx:       killCtx,
		kill:          kill,
		doneCtx:       doneCtx,
		done:          done,
		db:            db,
		stockCache:    symbolCache,
		RequestUpdate: make(chan core.StockSymbolType, 100),
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
		cache := make(map[core.StockSymbolType]bool)
		nextCheckDuration := 15 * time.Second

	mainLoop:
		for {
			select {
			case <-osg.killCtx.Done():
				break mainLoop
			case stockSymbol := <-osg.RequestUpdate:
				record := osg.stockCache.GetStock(stockSymbol)

				if record.ValidationStatus == core.ValidationValid &&
					record.DownloadStatus == core.ActivityEnabled &&
					record.OptionStatus == core.OptionsEnabled {
					cache[record.Symbol] = true
				}
			case <-time.After(nextCheckDuration):
				//now we start the actual task
				startTime := time.Now()
				agent, err := osg.db.GetAgent(osg.killCtx)
				if err != nil {
					logrus.Errorf("OptionSymbolGrabber: could not retrieve agent: %v", err)
					continue
				}

				//we every time we run we check all the stocks for those that should be checked
				today := core.NewDateTypeFromTime(time.Now())
				for _, stock := range osg.stockCache.GetAllStocks() {
					if stock.ValidationStatus == core.ValidationValid &&
						stock.DownloadStatus == core.ActivityEnabled &&
						stock.OptionStatus == core.OptionsEnabled &&
						stock.OptionListTimestamp.Before(today) {
						cache[stock.Symbol] = true
					}
				}

				//now we do a cursor check if there's any work to do
				// if not we bail this round
				if len(cache) == 0 {
					continue
				}

				//for each stock in the local cache get updated options and put them in the database
				// we update the database so the quotes can use them (which may o
				ctx, cancel := context.WithCancel(osg.killCtx)
				runningCount := 0
				finishedChan := make(chan bool, len(cache))
				for symbol := range cache {
					select {
					case <-osg.killCtx.Done():
						logrus.Errorf("OptionSymbolGrabber: context canceled")
						break mainLoop
					default:
					}
					runningCount++
					go processOS(ctx, symbol, agent, osg.db, osg.stockCache, finishedChan)
				}

				//clear the cache
				cache = make(map[core.StockSymbolType]bool)

				//next wait for results
			finishedLoop:
				for runningCount > 0 {
					select {
					case <-osg.killCtx.Done():
						break finishedLoop
					case <-finishedChan:
						runningCount--
					case <-time.After(30 * time.Minute):
						//a fail safe this should take 30 mins to complete
						logrus.Errorf("OptionSymbolGrabber: had an issue getting all the result in a timely manor")
						break finishedLoop
					}
				}
				cancel()

				logrus.Infof("OptionSymbolGrabber: finished a round in %v", time.Since(startTime))
			}
		}
		osg.done() //signal this is finished
	}(osg)
	return nil
}

const maxOSRetries = 3

func processOS(ctx context.Context, stock core.StockSymbolType, agent core.SibylAgent, db *database.SibylDatabase, stockCache *StockCache, finishedChan chan bool) {
	logrus.Infof("OptionSymbolGrabber: getting option symbols for %v", stock)

	retries := maxOSRetries
	var options []*core.OptionSymbolType
	var err error

	for retries > 0 {
		select {
		case <-ctx.Done():
			return
		default:
		}

		options, err = agent.GetStockOptionSymbols(ctx, stock)
		//this guy can return errors due to problem with some of the options but still return the valid ones
		if err != nil && len(options) == 0 {
			//here we have a full error, there were an error and zero options were returned
			logrus.Errorf("OptionSymbolGrabber: had an error while gather options for stock symbol:%v the error: %v", stock, err)
			retries--
			if retries == 0 {
				break
			}
			continue
		}

		if len(options) == 0 {
			logrus.Errorf("OptionSymbolGrabber: found 0 option symbols for %v, submitting to retry", stock)
			retries--
			if retries == 0 {
				break
			}
			continue
		}

		//this was that weird case where it can send back a partial error (aka the len(options) != 0)
		if err != nil {
			logrus.Errorf("OptionSymbolGrabber: had an error while gather options for stock symbol:%v the error: %v", stock, err)
		}

		if err := db.SetOptionsForStock(ctx, stock, options); err != nil {
			logrus.Errorf("OptionSymbolGrabber: failed during adding options to data, submitting to retry, error found: %v", err)
			//request an another update until we get a clean run
			retries--
			continue
		}

		//if we're here then we're as good as it gets so update the timestamp for the optionlist
		// we don't really care if this fails so we ignore the error if there is one
		stockCache.UpdateOptionListTimestamp(stock, core.NewDateTypeFromTime(time.Now()))

		break
	}
	finishedChan <- true
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
