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
				agent, err := sv.db.GetAgent(sv.killCtx)
				if err != nil {
					logrus.Errorf("StockValidator: could not retrieve agent: %v", err)
					continue
				}

				if stocks, err := sv.db.GetAllStockRecords(sv.killCtx); err != nil {
					logrus.Errorf("StockValidator: had a problem executing GetAllStocks: %v", err)
				} else {
					stocksWhereUpdated := false
					for _, stock := range stocks {
						if stock.ValidationStatus == core.ValidationPending {
							good, hasOptions, exchange, exchangeName, name, err := agent.VerifyStockSymbol(sv.killCtx, stock.Symbol)
							if err != nil {
								logrus.Errorf("StockValidator: had the following error: %v", err)
								continue
							}
							if good {
								if err := sv.db.StockSetExchangeInfoAndName(sv.killCtx, stock.Symbol, hasOptions, exchange, exchangeName, name); err != nil {
									logrus.Errorf("StockValidator: had the following error while updating stock: %v: %v", stock.Symbol, err)
								} else {
									if err := sv.db.StockValidate(sv.killCtx, stock.Symbol); err != nil {
										logrus.Errorf("StockValidator: had the following error while updating stock: %v to valid: %v", stock.Symbol, err)
									} else {
										logrus.Infof("StockValidator: the stock %v has been validated.", stock.Symbol)
										//we'll let the cache and OptionSymbolGrabber know
										// since the OptionSymbolGrabber will update the cache we'll
										// just let it know and it will take care of the rest
										stocksWhereUpdated = true
									}
								}
							} else {
								if err := sv.db.StockInvalidate(sv.killCtx, stock.Symbol); err != nil {
									logrus.Errorf("StockValidator: had the following error while updating stock:%v to invalid :%v", stock.Symbol, err)
								} else {
									logrus.Infof("StockValidator: the stock %v has NOT been validated.", stock.Symbol)
								}
							}
						}
					}
					if stocksWhereUpdated {
						//TODO consider moving this action into the SymbolsCache
						sv.optionSymbolGrabber.RequestUpdate <- true
					}
				}
			}
		}
		sqg.done() //signal this is finished
	}(sv)
	return nil
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
