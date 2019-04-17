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
	optionSymbolGrabber *OptionSymbolGrabber
	running             bool
	RequestUpdate       chan bool
}

func NewStockValidator(db *database.SibylDatabase, optionSymbolGrabber *OptionSymbolGrabber) *StockValidator {
	killCtx, kill := context.WithCancel(context.Background())
	doneCtx, done := context.WithCancel(context.Background())
	return &StockValidator{
		killCtx:             killCtx,
		kill:                kill,
		doneCtx:             doneCtx,
		done:                done,
		db:                  db,
		optionSymbolGrabber: optionSymbolGrabber,
		RequestUpdate:       make(chan bool, 100),
	}
}

func (sv *StockValidator) Run() error {
	if sv.running {
		return fmt.Errorf("StockValidator is already running")
	}

	sv.running = true
	go func(sqg *StockValidator) {
		ticker := time.NewTicker(1 * time.Minute)
	mainLoop:
		for {
			select {
			case <-sv.killCtx.Done():
				break mainLoop
			case <-ticker.C:
				startTime := time.Now()
				agent, err := sv.db.GetAgent(sv.killCtx)
				if err != nil {
					logrus.Errorf("StockValidator: could not retrieve agent: %v", err)
					continue
				}

				if stocks, err := sv.db.GetAllStockRecords(sv.killCtx); err != nil {
					logrus.Errorf("StockValidator: had a problem executing GetAllStocks: %v", err)
				} else {
					stocksWereUpdated := false
					updatedChan := make(chan bool, len(stocks))
					runningCount := 0
					ctx, cancel := context.WithCancel(sv.killCtx)
					for _, stock := range stocks {
						if stock.ValidationStatus == core.ValidationPending ||
							(stock.ValidationStatus != core.ValidationInvalid && stock.ValidationTimestamp.IsZero()) ||
							(stock.ValidationStatus == core.ValidationValid && stock.ValidationTimestamp.Before(core.NewDateTypeFromTime(time.Now()))) {
							runningCount++
							go validateStock(ctx, agent, sv.db, stock.Symbol, updatedChan)
						}
					}
					// now we will issue an updated if we get a true from the updateChan
				updateLoop:
					for runningCount > 0 {
						select {
						case <-sv.killCtx.Done():
							break updateLoop
						case update := <-updatedChan:
							stocksWereUpdated = stocksWereUpdated || update
							runningCount--
						case <-time.After(6 * time.Hour):
							//a fail safe this should take 6 hours to complete
							logrus.Errorf("StockValidator: had an issue getting all the result for stock validation in a timely manor")
							cancel()
							break updateLoop
						}
					}

					if stocksWereUpdated {
						//TODO consider moving this action into the SymbolsCache
						sv.optionSymbolGrabber.RequestUpdate <- true
					}
				}

				logrus.Infof("StockValidator: finished a round in %v", time.Since(startTime))
			}
		}
		sqg.done() //signal this is finished
	}(sv)
	return nil
}

func validateStock(ctx context.Context, agent core.SibylAgent, db *database.SibylDatabase, stock core.StockSymbolType, updated chan bool) {
	good, hasOptions, exchange, exchangeName, name, err := agent.VerifyStockSymbol(ctx, stock)
	if err != nil {
		logrus.Errorf("StockValidator: had the following error: %v", err)
	} else if good {
		if err := db.StockSetExchangeInfoAndName(ctx, stock, hasOptions, exchange, exchangeName, name); err != nil {
			logrus.Errorf("StockValidator: had the following error while updating stock: %v: %v", stock, err)
		} else {
			if err := db.StockValidate(ctx, stock); err != nil {
				logrus.Errorf("StockValidator: had the following error while updating stock: %v to valid: %v", stock, err)
			} else {
				logrus.Infof("StockValidator: the stock %v has been validated.", stock)
				//we'll let the cache and OptionSymbolGrabber know
				// since the OptionSymbolGrabber will update the cache we'll
				// just let it know and it will take care of the rest
				updated <- true
				db.StockSetValidateTimestamp(ctx, stock, core.NewDateTypeFromTime(time.Now()))
				return
			}
		}
	} else {
		if err := db.StockInvalidate(ctx, stock); err != nil {
			logrus.Errorf("StockValidator: had the following error while updating stock:%v to invalid :%v", stock, err)
			updated <- true
			db.StockSetValidateTimestamp(ctx, stock, core.NewDateTypeFromTime(time.Now()))
			return
		} else {
			logrus.Infof("StockValidator: the stock %v has NOT been validated.", stock)
		}
	}

	//if we made it here we didn't update anything
	updated <- false
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
