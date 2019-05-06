package internal

import (
	"context"
	"fmt"
	"github.com/nathanhack/sibyl/core"
	"github.com/nathanhack/sibyl/core/database"
	"github.com/sirupsen/logrus"
	"time"
)

type StockValidator struct {
	killCtx             context.Context
	kill                context.CancelFunc
	doneCtx             context.Context
	done                context.CancelFunc
	db                  *database.SibylDatabase
	stockCache          *StockCache
	optionSymbolGrabber *OptionSymbolGrabber
	running             bool
}

func NewStockValidator(db *database.SibylDatabase, stockCache *StockCache, optionSymbolGrabber *OptionSymbolGrabber) *StockValidator {
	killCtx, kill := context.WithCancel(context.Background())
	doneCtx, done := context.WithCancel(context.Background())
	return &StockValidator{
		killCtx:             killCtx,
		kill:                kill,
		doneCtx:             doneCtx,
		done:                done,
		db:                  db,
		stockCache:          stockCache,
		optionSymbolGrabber: optionSymbolGrabber,
	}
}

func (sv *StockValidator) Run() error {
	if sv.running {
		return fmt.Errorf("StockValidator is already running")
	}

	sv.running = true
	go func(sqg *StockValidator) {
		onceDailyDeadline, onceDailyDeadlineCancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))

		ticker := time.NewTicker(1 * time.Second)
		symbolCache := make(map[core.StockSymbolType]bool)
	mainLoop:
		for {
			select {
			case <-sv.killCtx.Done():
				break mainLoop
			case <-onceDailyDeadline.Done():
				onceDailyDeadlineCancel()
				onceDailyDeadline, onceDailyDeadlineCancel = context.WithDeadline(context.Background(), tomorrowAt6AM())
				now := core.NewDateTypeFromTime(time.Now())
				for _, stock := range sv.stockCache.GetAllStocks() {
					if stock.ValidationStatus == core.ValidationPending ||
						(stock.ValidationStatus != core.ValidationInvalid && stock.ValidationTimestamp.IsZero()) ||
						(stock.ValidationStatus == core.ValidationValid && stock.ValidationTimestamp.Before(now)) {
						symbolCache[stock.Symbol] = true
					}
				}
			case <-ticker.C:
				startTime := time.Now()
				agent, err := sv.db.GetAgent(sv.killCtx)
				if err != nil {
					logrus.Errorf("StockValidator: could not retrieve agent: %v", err)
					continue
				}

				for _, stock := range sv.stockCache.GetValidationStatus(core.ValidationPending) {
					symbolCache[stock.Symbol] = true
				}

				//now we do a cursor check if there's any work to do
				// if not we bail this round
				if len(symbolCache) == 0 {
					continue
				}

				updatedChan := make(chan updated, len(symbolCache))
				runningCount := 0
				ctx, cancel := context.WithCancel(sv.killCtx)
				for stock := range symbolCache {
					runningCount++
					go validateStock(ctx, agent, sv.stockCache, stock, updatedChan)
				}

				//we clear out the cache for the next round
				symbolCache = make(map[core.StockSymbolType]bool)

				// now we will issue an updated if we get a true from the updateChan
			updateLoop:
				for runningCount > 0 {
					select {
					case <-sv.killCtx.Done():
						break updateLoop
					case update := <-updatedChan:
						if update.Updated {
							sv.optionSymbolGrabber.RequestUpdate <- update.Stock
						}
						runningCount--
					case <-time.After(6 * time.Hour):
						//a fail safe this should take 6 hours to complete
						logrus.Errorf("StockValidator: had an issue getting all the result for stock validation in a timely manor")
						break updateLoop
					}
				}
				cancel()

				logrus.Infof("StockValidator: finished a round in %v", time.Since(startTime))
			}
		}
		sqg.done() //signal this is finished
	}(sv)
	return nil
}

type updated struct {
	Stock   core.StockSymbolType
	Updated bool
}

func validateStock(ctx context.Context, agent core.SibylAgent, stockCache *StockCache, stock core.StockSymbolType, updatedChan chan updated) {
	good, hasOptions, exchange, exchangeDescription, name, err := agent.VerifyStockSymbol(ctx, stock)
	if err != nil {
		logrus.Errorf("StockValidator: had the following error: %v", err)
		updatedChan <- updated{
			Stock:   stock,
			Updated: false,
		}
		return
	}

	if good {
		logrus.Infof("StockValidator: the stock %v was valid", stock)
		optionstatus := core.OptionsDisabled
		if hasOptions {
			optionstatus = core.OptionsEnabled

		}

		for _, err := range []error{
			stockCache.UpdateOptionStatus(stock, optionstatus),
			stockCache.UpdateExchange(stock, exchange),
			stockCache.UpdateExchangeDescription(stock, exchangeDescription),
			stockCache.UpdateName(stock, name),
			stockCache.UpdateValidationStatus(stock, core.ValidationValid),
			stockCache.UpdateValidationTimestamp(stock, core.NewDateTypeFromTime(time.Now())),
		} {
			if err != nil {
				logrus.Errorf("StockValidator: failed to update stock %v: %v", stock, err)
				updatedChan <- updated{
					Stock:   stock,
					Updated: false,
				}
				return
			}
		}
	} else {
		logrus.Infof("StockValidator: the stock %v was NOT valid", stock)
		for _, err := range []error{
			stockCache.UpdateValidationStatus(stock, core.ValidationInvalid),
			stockCache.UpdateValidationTimestamp(stock, core.NewDateTypeFromTime(time.Now())),
		} {
			if err != nil {
				logrus.Errorf("StockValidator: failed to update stock %v: %v", stock, err)
				updatedChan <- updated{
					Stock:   stock,
					Updated: false,
				}
				return
			}
		}
	}

	updatedChan <- updated{
		Stock:   stock,
		Updated: true,
	}
}

func (sv *StockValidator) Stop(waitUpTo time.Duration) {
	//next stop the quoter
	logrus.Infof("Waiting for StockValidator to finish")
	startTime := time.Now()
	sv.kill()
	select {
	case <-sv.doneCtx.Done():
		logrus.Infof("StockValidator finished in %v", time.Since(startTime))
	case <-time.After(waitUpTo):
		logrus.Errorf("StockValidator failed to gracefully finish in %v", time.Since(startTime))
	}
}
