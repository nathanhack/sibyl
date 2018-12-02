package internal

import (
	"context"
	"fmt"
	"github.com/nathanhack/sibyl/core"
	"github.com/nathanhack/sibyl/core/database"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type SymbolCache struct {
	killCtx                         context.Context
	kill                            context.CancelFunc
	doneCtx                         context.Context
	done                            context.CancelFunc
	db                              *database.SibylDatabase
	HistorySymbols                  map[core.StockSymbolType]bool
	HistorySymbolsMu                sync.RWMutex
	HistorySymbolsChanged           chan bool
	IntradaySymbols                 map[core.StockSymbolType]bool
	IntradaySymbolsMu               sync.RWMutex
	IntradaySymbolsChanged          chan bool
	QuoteOptionsSymbols             map[core.OptionSymbolType]bool
	QuoteOptionsSymbolsMu           sync.RWMutex
	QuoteStockSymbols               map[core.StockSymbolType]bool
	QuoteStockSymbolsMu             sync.RWMutex
	StableQuoteOptionSymbols        map[core.OptionSymbolType]bool
	StableQuoteOptionSymbolsMu      sync.RWMutex
	StableQuoteOptionSymbolsChanged chan bool
	StableQuoteStockSymbols         map[core.StockSymbolType]bool
	StableQuoteStockSymbolsMu       sync.RWMutex
	StableQuoteStockSymbolsChanged  chan bool

	RequestUpdate chan bool
	running       bool
}

func NewSymbolsCache(db *database.SibylDatabase) *SymbolCache {
	killCtx, kill := context.WithCancel(context.Background())
	doneCtx, done := context.WithCancel(context.Background())
	return &SymbolCache{
		killCtx:                         killCtx,
		kill:                            kill,
		doneCtx:                         doneCtx,
		done:                            done,
		db:                              db,
		HistorySymbols:                  make(map[core.StockSymbolType]bool),
		HistorySymbolsMu:                sync.RWMutex{},
		HistorySymbolsChanged:           make(chan bool, 1),
		IntradaySymbols:                 make(map[core.StockSymbolType]bool),
		IntradaySymbolsMu:               sync.RWMutex{},
		IntradaySymbolsChanged:          make(chan bool, 1),
		QuoteOptionsSymbols:             make(map[core.OptionSymbolType]bool),
		QuoteOptionsSymbolsMu:           sync.RWMutex{},
		QuoteStockSymbols:               make(map[core.StockSymbolType]bool),
		QuoteStockSymbolsMu:             sync.RWMutex{},
		StableQuoteOptionSymbolsChanged: make(chan bool, 1),
		StableQuoteStockSymbols:         make(map[core.StockSymbolType]bool),
		StableQuoteStockSymbolsMu:       sync.RWMutex{},
		StableQuoteStockSymbolsChanged:  make(chan bool, 1),
		RequestUpdate:                   make(chan bool, 100),
	}
}

//RunUpdater: creates and runs a go routine that will periodically (and if requested) update the symbols in the struct
func (sc *SymbolCache) Run() error {
	//we want to prime the update requester becuase we want to update the cache as soon as we start running
	if sc.running {
		return fmt.Errorf("SymbolCache is already running")
	}

	sc.running = true
	sc.RequestUpdate <- true

	go func(sc *SymbolCache) {
		updateTicker := time.NewTicker(30 * time.Minute)
		regulateTicker := time.NewTicker(10 * time.Second)
	mainLoop:
		for {
			select {
			case <-sc.killCtx.Done():
				break mainLoop
			case <-updateTicker.C:
				//every 30 mins we want to update the state of the cache
				sc.RequestUpdate <- true
			case <-regulateTicker.C:
				//so every 10 seconds we check if there are any requests to update the cache
				select {
				case <-sc.RequestUpdate:
				default:
					continue
				}
				startTime := time.Now()
				//we start by draining the channel dry
			dryLoop:
				for {
					select {
					case <-sc.RequestUpdate:
						break
					default:
						break dryLoop
					}
				}

				if err := sc.updatedSymbolsList(); err != nil {
					logrus.Errorf("SymbolCache: had an error while gathering new symbols: %v", err)
					//since we had a failure we'll keep trying to update the cache
					sc.RequestUpdate <- true
				}

				logrus.Infof("SymbolCache: updated symbols from database in %v", time.Since(startTime))
			}
		}
		sc.done()
	}(sc)
	return nil
}

func (sc *SymbolCache) updatedSymbolsList() error {
	historySymbols := make(map[core.StockSymbolType]bool)
	intradaySymbols := make(map[core.StockSymbolType]bool)
	quoteOptionSymbols := make(map[core.OptionSymbolType]bool)
	stableQuoteOptionSymbols := make(map[core.OptionSymbolType]bool)
	quoteStockSymbols := make(map[core.StockSymbolType]bool)
	stableQuoteStockSymbols := make(map[core.StockSymbolType]bool)

	allStocks, err := sc.db.GetAllStockRecords(sc.killCtx)
	if err != nil {
		return fmt.Errorf("getUpdatedSymbolsList: had a problem getting list of stocks: %v", err)
	}

	//now we have all the stocks so pick on the ones that are good with downloading
	for _, stock := range allStocks {
		if stock.ValidationStatus == core.ValidationValid &&
			stock.DownloadStatus == core.ActivityEnabled {

			if stock.HistoryStatus == core.ActivityEnabled {
				historySymbols[stock.Symbol] = true
			}
			if stock.IntradayHistoryStatus == core.ActivityEnabled {
				intradaySymbols[stock.Symbol] = true
			}
			if stock.QuotesStatus == core.ActivityEnabled {
				quoteStockSymbols[stock.Symbol] = true
			}
			if stock.StableQuotesStatus == core.ActivityEnabled {
				stableQuoteStockSymbols[stock.Symbol] = true
			}
		}
		//we do a quick context check
		select {
		case <-sc.killCtx.Done():
			return fmt.Errorf("getUpdatedSymbolsList: error :context canceled")
		default:
		}

	}

	//now time for all the options
	errString := []string{}

	if options, err := sc.db.GetOptionsFor(sc.killCtx, quoteStockSymbols); err != nil {
		errString = append(errString, fmt.Sprintf("problem with quote options symbols: %v", err))
	} else {
		for _, option := range options {
			quoteOptionSymbols[*option] = true
		}
	}

	if options, err := sc.db.GetOptionsFor(sc.killCtx, stableQuoteStockSymbols); err != nil {
		errString = append(errString, fmt.Sprintf("problem with stable quote option symbols: %v", err))
	} else {
		for _, option := range options {
			stableQuoteOptionSymbols[*option] = true
		}
	}

	//we do a quick context check
	select {
	case <-sc.killCtx.Done():
		return fmt.Errorf("getUpdatedSymbolsList: error :context canceled")
	default:
	}

	//TODO consider adding some smarts for the case the update was bad aka results in missing symbols
	sc.HistorySymbolsMu.Lock()
	//we compare teh lists and if the new is different
	// we swap it and send notification of a change
	if len(sc.HistorySymbols) != len(historySymbols) {
		sc.HistorySymbols = historySymbols
		//now nonblock add to channel
		select {
		case sc.HistorySymbolsChanged <- true:
		default:
			//only here if it was full and didn't need this added
		}
	} else {
		for symbol := range historySymbols {
			if _, has := sc.HistorySymbols[symbol]; !has {
				sc.HistorySymbols = historySymbols
				//now nonblock add to channel
				select {
				case sc.HistorySymbolsChanged <- true:
				default:
					//only here if it was full and didn't need this added
				}
				break
			}
		}
	}
	sc.HistorySymbolsMu.Unlock()

	//we do a quick context check
	select {
	case <-sc.killCtx.Done():
		return fmt.Errorf("getUpdatedSymbolsList: error :context canceled")
	default:
	}

	sc.IntradaySymbolsMu.Lock()
	//we compare teh lists and if the new is different
	// we swap it and send notification of a change
	if len(sc.IntradaySymbols) != len(intradaySymbols) {
		sc.IntradaySymbols = intradaySymbols
		//now nonblock add to channel
		select {
		case sc.IntradaySymbolsChanged <- true:
		default:
			//only here if it was full and didn't need this added
		}
	} else {
		for symbol := range intradaySymbols {
			if _, has := sc.IntradaySymbols[symbol]; !has {
				sc.IntradaySymbols = intradaySymbols
				//now nonblock add to channel
				select {
				case sc.IntradaySymbolsChanged <- true:
				default:
					//only here if it was full and didn't need this added
				}
				break
			}
		}
	}
	sc.IntradaySymbolsMu.Unlock()

	sc.QuoteOptionsSymbolsMu.Lock()
	sc.QuoteOptionsSymbols = quoteOptionSymbols
	sc.QuoteOptionsSymbolsMu.Unlock()

	sc.QuoteStockSymbolsMu.Lock()
	sc.QuoteStockSymbols = quoteStockSymbols
	sc.QuoteStockSymbolsMu.Unlock()

	//we do a quick context check
	select {
	case <-sc.killCtx.Done():
		return fmt.Errorf("getUpdatedSymbolsList: error :context canceled")
	default:
	}

	sc.StableQuoteOptionSymbolsMu.Lock()
	//we compare teh lists and if the new is different
	// we swap it and send notification of a change
	if len(sc.StableQuoteOptionSymbols) != len(stableQuoteOptionSymbols) {
		sc.StableQuoteOptionSymbols = stableQuoteOptionSymbols
		//now nonblock add to channel
		select {
		case sc.StableQuoteOptionSymbolsChanged <- true:
		default:
			//only here if it was full and didn't need this added
		}
	} else {
		for symbol := range stableQuoteOptionSymbols {
			if _, has := sc.StableQuoteOptionSymbols[symbol]; !has {
				sc.StableQuoteOptionSymbols = stableQuoteOptionSymbols
				//now nonblock add to channel
				select {
				case sc.StableQuoteOptionSymbolsChanged <- true:
				default:
					//only here if it was full and didn't need this added
				}
				break
			}
		}
	}
	sc.StableQuoteOptionSymbolsMu.Unlock()

	//we do a quick context check
	select {
	case <-sc.killCtx.Done():
		return fmt.Errorf("getUpdatedSymbolsList: error :context canceled")
	default:
	}

	sc.StableQuoteStockSymbolsMu.Lock()
	//we compare teh lists and if the new is different
	// we swap it and send notification of a change
	if len(sc.StableQuoteStockSymbols) != len(stableQuoteStockSymbols) {
		sc.StableQuoteStockSymbols = stableQuoteStockSymbols

		//now nonblock add to channel
		select {
		case sc.StableQuoteStockSymbolsChanged <- true:
		default:
			//only here if it was full and didn't need this added
		}
	} else {
		for symbol := range stableQuoteStockSymbols {
			if _, has := sc.StableQuoteStockSymbols[symbol]; !has {
				sc.StableQuoteStockSymbols = stableQuoteStockSymbols
				//now nonblock add to channel
				select {
				case sc.StableQuoteStockSymbolsChanged <- true:
				default:
					//only here if it was full and didn't need this added
				}
				break
			}
		}
	}
	sc.StableQuoteStockSymbolsMu.Unlock()

	if len(errString) > 0 {
		return fmt.Errorf("getUpdatedSymbolsList: had an problem getting database options, had error: %v", errString)
	}

	return nil
}

func (sc *SymbolCache) Stop(waitUpTo time.Duration) {
	//next stop the quoter
	logrus.Infof("Waiting for SymbolCache to finish")
	startTime := time.Now()
	sc.kill()
	select {
	case <-sc.doneCtx.Done():
		logrus.Infof("SymbolCache finished in %v", time.Since(startTime))
	case <-time.After(waitUpTo):
		logrus.Errorf("SymbolCache failed to gracefully finish in %v", time.Since(startTime))
	}
}
