package internal

import (
	"context"
	"fmt"
	"github.com/nathanhack/sibyl/core"
	"github.com/nathanhack/sibyl/core/database"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

//StockCache is meant to be a database backed in memory cache.  Upon startup it loads from the database.
// Then it will periodical update the database based on what is in the cache.  Methods to get a copy of the
// current list of stocks will be available and methods to update specific fields based
// on topic (history, intraday, etc).
type StockCache struct {
	db                              *database.SibylDatabase
	done                            context.CancelFunc
	doneCtx                         context.Context
	kill                            context.CancelFunc
	killCtx                         context.Context
	running                         bool
	stateChange                     bool
	stocks                          map[core.StockSymbolType]core.SibylStockRecord
	stocksMu                        sync.RWMutex
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
}

func NewStockCache(db *database.SibylDatabase) *StockCache {
	killCtx, kill := context.WithCancel(context.Background())
	doneCtx, done := context.WithCancel(context.Background())
	return &StockCache{
		db:                              db,
		done:                            done,
		doneCtx:                         doneCtx,
		kill:                            kill,
		killCtx:                         killCtx,
		stocks:                          make(map[core.StockSymbolType]core.SibylStockRecord),
		stocksMu:                        sync.RWMutex{},
		QuoteOptionsSymbols:             make(map[core.OptionSymbolType]bool),
		QuoteOptionsSymbolsMu:           sync.RWMutex{},
		QuoteStockSymbols:               make(map[core.StockSymbolType]bool),
		QuoteStockSymbolsMu:             sync.RWMutex{},
		StableQuoteOptionSymbolsChanged: make(chan bool, 1),
		StableQuoteStockSymbols:         make(map[core.StockSymbolType]bool),
		StableQuoteStockSymbolsMu:       sync.RWMutex{},
		StableQuoteStockSymbolsChanged:  make(chan bool, 1),
	}
}

//RunUpdater: creates and runs a go routine that will periodically (and if requested) update the symbols in the struct
func (sc *StockCache) Run() error {
	//we want to prime the update requester becuase we want to update the cache as soon as we start running
	if sc.running {
		return fmt.Errorf("StockCache is already running")
	}
	sc.running = true

	go func(sc *StockCache) {
		//on startup the cache is empty so we need to get data from the DB
		// then update the DB ever so often
	getData:
		for {
			select {
			case <-sc.killCtx.Done():
				break getData
			case <-time.After(5 * time.Second):
				if err := sc.updatedSymbolsList(); err != nil {
					logrus.Errorf("StockCache: had an error while getting data from DB: %v", err)
				} else {
					break getData
				}
			}
		}
		logrus.Infof("StockCache: loaded Data from Database")
		// now the only thing left to do is update the DB every so often
		updateDBTicker := time.NewTicker(1 * time.Second)
	mainLoop:
		for {
			select {
			case <-sc.killCtx.Done():
				break mainLoop
			case <-updateDBTicker.C:
				startTime := time.Now()
				if !sc.stateChange {
					continue
				}
				sc.stateChange = false

				stocks, err := sc.db.GetAllStockRecords(sc.killCtx)
				if err != nil {
					logrus.Infof("StockCache: unable to retrieve stocks for update")
					continue
				}

				sc.stocksMu.RLock()
				//first we remove anything that's no longer in the cache
				for _, stock := range stocks {
					if _, has := sc.stocks[stock.Symbol]; !has {
						sc.db.StockDelete(sc.killCtx, stock.Symbol)
					}
				}

				//finally we take what's in the cache and update the DB
				for _, stock := range sc.stocks {
					err := sc.db.InsertOrUpdateStock(sc.killCtx, &stock)
					if err != nil {
						logrus.Errorf("StockCache:  %v", err)
					}
				}
				sc.stocksMu.RUnlock()
				logrus.Infof("StockCache: updated Database in %v", time.Since(startTime))
			}
		}
		sc.done()
	}(sc)
	return nil
}

func (sc *StockCache) updatedSymbolsList() error {
	quoteOptionSymbols := make(map[core.OptionSymbolType]bool)
	stableQuoteOptionSymbols := make(map[core.OptionSymbolType]bool)
	quoteStockSymbols := make(map[core.StockSymbolType]bool)
	stableQuoteStockSymbols := make(map[core.StockSymbolType]bool)

	allStocks, err := sc.db.GetAllStockRecords(sc.killCtx)
	if err != nil {
		return fmt.Errorf("getUpdatedSymbolsList: had a problem getting list of stocks: %v", err)
	}

	sc.stocksMu.Lock()
	//now we have all the stocks so pick on the ones that are good with downloading
	for _, stock := range allStocks {
		sc.stocks[stock.Symbol] = *stock
		if stock.ValidationStatus == core.ValidationValid &&
			stock.DownloadStatus == core.ActivityEnabled {

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
	sc.stocksMu.Unlock()

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

	//we do a quick context check
	select {
	case <-sc.killCtx.Done():
		return fmt.Errorf("getUpdatedSymbolsList: error :context canceled")
	default:
	}

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

func (sc *StockCache) AddStockSymbol(symbol core.StockSymbolType) {
	sc.stocksMu.RLock()
	if _, has := sc.stocks[symbol]; !has {
		sc.stocks[symbol] = core.SibylStockRecord{
			Symbol: core.StockSymbolType(strings.ToUpper(string(symbol))),
		}
		sc.stateChange = true
	}
	sc.stocksMu.RUnlock()
}

func (sc *StockCache) GetStock(symbol core.StockSymbolType) *core.SibylStockRecord {
	sc.stocksMu.RLock()
	defer sc.stocksMu.RUnlock()
	if s, has := sc.stocks[symbol]; has {
		return &s
	}
	return nil
}

func (sc *StockCache) RemoveStockSymbol(symbol core.StockSymbolType) {

	sc.stocksMu.RLock()
	if _, has := sc.stocks[symbol]; has {
		sc.stocks[symbol] = core.SibylStockRecord{
			Symbol: symbol,
		}
		sc.stateChange = true
	}
	sc.stocksMu.RUnlock()
}

func (sc *StockCache) GetAllStocks() []core.SibylStockRecord {
	toReturn := make([]core.SibylStockRecord, 0, len(sc.stocks))
	sc.stocksMu.RLock()
	for _, stock := range sc.stocks {
		toReturn = append(toReturn, stock)
	}
	sc.stocksMu.RUnlock()
	return toReturn
}

//GetValidationStatus returns a slice containing the stock records with a validationStatus equal to status
func (sc *StockCache) GetValidationStatus(status core.ValidationStatusType) []core.SibylStockRecord {
	toReturn := make([]core.SibylStockRecord, 0, len(sc.stocks))
	sc.stocksMu.RLock()
	for _, stock := range sc.stocks {
		if stock.ValidationStatus == status {
			toReturn = append(toReturn, stock)
		}
	}
	sc.stocksMu.RUnlock()
	return toReturn
}

func (sc *StockCache) HistoryStocks(hasDaily, hasWeekly, hasMonthly, hasYearly bool) []core.SibylStockRecord {
	toReturn := make([]core.SibylStockRecord, 0, len(sc.stocks))
	sc.stocksMu.RLock()
	for _, stock := range sc.stocks {
		if stock.DownloadStatus == core.ActivityEnabled &&
			stock.ValidationStatus == core.ValidationValid &&
			((hasDaily && hasDaily == stock.HistoryStatus.HasDaily()) ||
				(hasWeekly && hasWeekly == stock.HistoryStatus.HasWeekly()) ||
				(hasMonthly && hasMonthly == stock.HistoryStatus.HasMonthly()) ||
				(hasYearly && hasYearly == stock.HistoryStatus.HasYearly())) {
			toReturn = append(toReturn, stock)
		}
	}
	sc.stocksMu.RUnlock()
	return toReturn
}

func (sc *StockCache) IntradayStocks(hasTick, has1Min, has5Min, mustBeActive bool) []core.SibylStockRecord {
	toReturn := make([]core.SibylStockRecord, 0, len(sc.stocks))
	sc.stocksMu.RLock()
	for _, stock := range sc.stocks {
		if stock.DownloadStatus == core.ActivityEnabled &&
			stock.ValidationStatus == core.ValidationValid &&
			((hasTick && hasTick == stock.IntradayStatus.HasTicks()) ||
				(has1Min && has1Min == stock.IntradayStatus.Has1Min()) ||
				(has5Min && has5Min == stock.IntradayStatus.Has5Min())) &&
			((mustBeActive && mustBeActive == stock.IntradayState.IsActive()) ||
				mustBeActive == false) {
			tmp := stock
			toReturn = append(toReturn, tmp)
		}
	}
	sc.stocksMu.RUnlock()
	return toReturn
}

func (sc *StockCache) UpdateDownloadStatus(symbol core.StockSymbolType, status core.ActivityStatusType) error {
	sc.stocksMu.Lock()
	x, has := sc.stocks[symbol]
	if !has {
		return fmt.Errorf("StockCache could not find %v to update DownloadStatus", symbol)
	}
	x.DownloadStatus = status
	sc.stocks[symbol] = x
	sc.stateChange = true
	sc.stocksMu.Unlock()
	return nil
}

func (sc *StockCache) UpdateExchange(symbol core.StockSymbolType, exchange string) error {
	sc.stocksMu.Lock()
	x, has := sc.stocks[symbol]
	if !has {
		return fmt.Errorf("StockCache could not find %v to update Exchange", symbol)
	}
	x.Exchange = exchange
	sc.stocks[symbol] = x
	sc.stateChange = true
	sc.stocksMu.Unlock()
	return nil
}

func (sc *StockCache) UpdateExchangeDescription(symbol core.StockSymbolType, exchangeDescription string) error {
	sc.stocksMu.Lock()
	x, has := sc.stocks[symbol]
	if !has {
		return fmt.Errorf("StockCache could not find %v to update ExchangeDescription", symbol)
	}
	x.ExchangeDescription = exchangeDescription
	sc.stocks[symbol] = x
	sc.stateChange = true
	sc.stocksMu.Unlock()
	return nil
}

func (sc *StockCache) UpdateOptionStatus(symbol core.StockSymbolType, status core.OptionStatusType) error {
	sc.stocksMu.Lock()
	x, has := sc.stocks[symbol]
	if !has {
		return fmt.Errorf("StockCache could not find %v to update OptionStatus", symbol)
	}
	x.OptionStatus = status
	sc.stocks[symbol] = x
	sc.stateChange = true
	sc.stocksMu.Unlock()
	return nil
}

func (sc *StockCache) UpdateHistoryStatus(symbol core.StockSymbolType, status core.HistoryStatusType) error {
	sc.stocksMu.Lock()
	x, has := sc.stocks[symbol]
	if !has {
		return fmt.Errorf("StockCache could not find %v to update HistoryStatus", symbol)
	}
	x.HistoryStatus = status
	sc.stocks[symbol] = x
	sc.stateChange = true
	sc.stocksMu.Unlock()
	return nil
}

func (sc *StockCache) UpdateHistoryTimestamp(symbol core.StockSymbolType, date core.DateType) error {
	sc.stocksMu.Lock()
	x, has := sc.stocks[symbol]
	if !has {
		return fmt.Errorf("StockCache could not find %v to update HistoryTimestamp", symbol)
	}
	x.HistoryTimestamp = date
	sc.stocks[symbol] = x
	sc.stateChange = true
	sc.stocksMu.Unlock()
	return nil
}

func (sc *StockCache) UpdateIntradayState(symbol core.StockSymbolType, state core.IntradayStateType) error {
	sc.stocksMu.Lock()
	x, has := sc.stocks[symbol]
	if !has {
		return fmt.Errorf("StockCache could not find %v to update IntradayState", symbol)
	}
	x.IntradayState = state
	sc.stocks[symbol] = x
	sc.stateChange = true
	sc.stocksMu.Unlock()
	return nil
}

func (sc *StockCache) UpdateIntradayStatus(symbol core.StockSymbolType, status core.IntradayStatusType) error {
	sc.stocksMu.Lock()
	x, has := sc.stocks[symbol]
	if !has {
		return fmt.Errorf("StockCache could not find %v to update IntradayStatus", symbol)
	}
	x.IntradayStatus = status
	sc.stocks[symbol] = x
	sc.stateChange = true
	sc.stocksMu.Unlock()
	return nil
}

func (sc *StockCache) UpdateIntradayTimestamp1Min(symbol core.StockSymbolType, timestamp core.TimestampType) error {
	sc.stocksMu.Lock()
	x, has := sc.stocks[symbol]
	if !has {
		return fmt.Errorf("StockCache could not find %v to update IntradayTimestamp1Min", symbol)
	}
	x.IntradayTimestamp1Min = timestamp
	sc.stocks[symbol] = x
	sc.stateChange = true
	sc.stocksMu.Unlock()
	return nil
}

func (sc *StockCache) UpdateIntradayTimestamp5Min(symbol core.StockSymbolType, timestamp core.TimestampType) error {
	sc.stocksMu.Lock()
	x, has := sc.stocks[symbol]
	if !has {
		return fmt.Errorf("StockCache could not find %v to update IntradayTimestamp5Min", symbol)
	}
	x.IntradayTimestamp5Min = timestamp
	sc.stocks[symbol] = x
	sc.stateChange = true
	sc.stocksMu.Unlock()
	return nil
}

func (sc *StockCache) UpdateIntradayTimestampTick(symbol core.StockSymbolType, timestamp core.TimestampType) error {
	sc.stocksMu.Lock()
	x, has := sc.stocks[symbol]
	if !has {
		return fmt.Errorf("StockCache could not find %v to update IntradayTimestampTick", symbol)
	}
	x.IntradayTimestampTick = timestamp
	sc.stocks[symbol] = x
	sc.stateChange = true
	sc.stocksMu.Unlock()
	return nil
}

func (sc *StockCache) UpdateName(symbol core.StockSymbolType, name string) error {
	sc.stocksMu.Lock()
	x, has := sc.stocks[symbol]
	if !has {
		return fmt.Errorf("StockCache could not find %v to update Name", symbol)
	}
	x.Name = name
	sc.stocks[symbol] = x
	sc.stateChange = true
	sc.stocksMu.Unlock()
	return nil
}

func (sc *StockCache) UpdateOptionListTimestamp(symbol core.StockSymbolType, date core.DateType) error {
	sc.stocksMu.Lock()
	x, has := sc.stocks[symbol]
	if !has {
		return fmt.Errorf("StockCache could not find %v to update OptionListTimestamp", symbol)
	}
	x.OptionListTimestamp = date
	sc.stocks[symbol] = x
	sc.stateChange = true
	sc.stocksMu.Unlock()
	return nil
}

func (sc *StockCache) UpdateQuoteStatus(symbol core.StockSymbolType, status core.ActivityStatusType) error {
	sc.stocksMu.Lock()
	x, has := sc.stocks[symbol]
	if !has {
		return fmt.Errorf("StockCache could not find %v to update QuotesStatus", symbol)
	}
	x.QuotesStatus = status
	sc.stocks[symbol] = x
	sc.stateChange = true
	sc.stocksMu.Unlock()
	return nil
}

func (sc *StockCache) UpdateStableQuoteStatus(symbol core.StockSymbolType, status core.ActivityStatusType) error {
	sc.stocksMu.Lock()
	x, has := sc.stocks[symbol]
	if !has {
		return fmt.Errorf("StockCache could not find %v to update StableQuotesStatus", symbol)
	}
	x.StableQuotesStatus = status
	sc.stocks[symbol] = x
	sc.stateChange = true
	sc.stocksMu.Unlock()
	return nil
}

func (sc *StockCache) UpdateValidationStatus(symbol core.StockSymbolType, status core.ValidationStatusType) error {
	sc.stocksMu.Lock()
	x, has := sc.stocks[symbol]
	if !has {
		return fmt.Errorf("StockCache could not find %v to update ValidationStatus", symbol)
	}
	x.ValidationStatus = status
	sc.stocks[symbol] = x
	sc.stateChange = true
	sc.stocksMu.Unlock()
	return nil
}

func (sc *StockCache) UpdateValidationTimestamp(symbol core.StockSymbolType, date core.DateType) error {
	sc.stocksMu.Lock()
	x, has := sc.stocks[symbol]
	if !has {
		return fmt.Errorf("StockCache could not find %v to update ValidationTimestamp", symbol)
	}
	x.ValidationTimestamp = date
	sc.stocks[symbol] = x
	sc.stateChange = true
	sc.stocksMu.Unlock()
	return nil
}

func (sc *StockCache) Stop(waitUpTo time.Duration) {
	//next stop the StockCache
	logrus.Infof("Waiting for StockCache to finish")
	startTime := time.Now()
	sc.kill()
	select {
	case <-sc.doneCtx.Done():
		logrus.Infof("StockCache finished in %v", time.Since(startTime))
	case <-time.After(waitUpTo):
		logrus.Errorf("StockCache failed to gracefully finish in %v", time.Since(startTime))
	}
}
