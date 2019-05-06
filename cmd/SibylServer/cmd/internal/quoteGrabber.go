package internal

import (
	"context"
	"fmt"
	"github.com/nathanhack/sibyl/core"
	"github.com/nathanhack/sibyl/core/database"
	"github.com/sirupsen/logrus"
	"time"
)

type QuoteGrabber struct {
	killCtx    context.Context
	kill       context.CancelFunc
	doneCtx    context.Context
	done       context.CancelFunc
	db         *database.SibylDatabase
	stockCache *StockCache
	running    bool
}

func NewQuoteGrabber(db *database.SibylDatabase, symbolCache *StockCache) *QuoteGrabber {
	killCtx, kill := context.WithCancel(context.Background())
	doneCtx, done := context.WithCancel(context.Background())
	return &QuoteGrabber{
		killCtx:    killCtx,
		kill:       kill,
		doneCtx:    doneCtx,
		done:       done,
		db:         db,
		stockCache: symbolCache,
	}
}

func (qg *QuoteGrabber) Run() error {
	if qg.running {
		return fmt.Errorf("QuoteGrabber is already running")
	}

	qg.running = true
	go func(qg *QuoteGrabber) {

		// we end up in weird cases when it takes a long time
		// to help with this we start two go routines to hand the database
		// communications for storing the stock and options
		// each channel can buffer 800 sets of quotes, 12 hours
		// should only be 720 so plenty of room
		stockQuotesChannel := make(chan []*core.SibylStockQuoteRecord, 800)
		optionQuotesChannel := make(chan []*core.SibylOptionQuoteRecord, 800)

		//the first one for stocks
		go func(qg *QuoteGrabber) {
		loaderLoop:
			for {
				select {
				case <-qg.killCtx.Done():
					break loaderLoop
				case quotes := <-stockQuotesChannel:
					//we'll see if there is more than one in the channel
					startTime := time.Now()
				drainLoop:
					for {
						select {
						case more := <-stockQuotesChannel:
							quotes = append(quotes, more...)
							continue
						default:
							break drainLoop
						}
					}
					//then add them
					if err := qg.db.LoadStockQuoteRecords(qg.killCtx, quotes); err != nil {
						logrus.Errorf("QuoteGrabber: failed to load %v stock quotes into database: %v", len(quotes), err)
					}
					logrus.Infof("QuoteGrabber: finished loading %v stock quotes into database in %v", len(quotes), time.Since(startTime))
				}
			}
			logrus.Infof("QuoteGrabber: database stocks loader thread stopped.")
		}(qg)

		//the second one for options
		go func(qg *QuoteGrabber) {
		loaderLoop:
			for {
				select {
				case <-qg.killCtx.Done():
					break loaderLoop
				case quotes := <-optionQuotesChannel:
					//we'll see if there is more than one in the channel
					startTime := time.Now()
				drainLoop:
					for {
						select {
						case more := <-optionQuotesChannel:
							quotes = append(quotes, more...)
							continue
						default:
							break drainLoop
						}
					}
					//then add them
					if err := qg.db.LoadOptionQuoteRecords(qg.killCtx, quotes); err != nil {
						logrus.Errorf("QuoteGrabber: failed to load %v option quotes into database: %v", len(quotes), err)
					}
					logrus.Infof("QuoteGrabber: finished loading %v option quotes into database in %v", len(quotes), time.Since(startTime))
				}
			}
			logrus.Infof("QuoteGrabber: database options loader thread stopped.")
		}(qg)

		///Quotes are currently set to happen every minutes during the week from the hours of 0700-1800
		currentTime := time.Now()
		var startTime = currentTime.Truncate(time.Minute).Add(time.Minute)
		select {
		case <-qg.killCtx.Done():
		case <-time.After(startTime.Sub(currentTime)):
		}
		ticker := time.NewTicker(1 * time.Minute)
	mainLoop:
		for {
			select {
			case <-qg.killCtx.Done():
				break mainLoop
			case <-ticker.C:
				go func(qg *QuoteGrabber) {
					qg.executeOneRound(stockQuotesChannel, optionQuotesChannel)
				}(qg)
			}
		}

		qg.done() //signal this is finished
	}(qg)
	return nil
}

func (qg *QuoteGrabber) executeOneRound(stockQuotesChannel chan []*core.SibylStockQuoteRecord, optionQuotesChannel chan []*core.SibylOptionQuoteRecord) {
	qg.stockCache.QuoteStockSymbolsMu.RLock()
	quoteStockSymbolsToDownLoad := qg.stockCache.QuoteStockSymbols
	qg.stockCache.QuoteStockSymbolsMu.RUnlock()

	qg.stockCache.QuoteOptionsSymbolsMu.RLock()
	quoteOptionSymbolsToDownLoad := qg.stockCache.QuoteOptionsSymbols
	qg.stockCache.QuoteOptionsSymbolsMu.RUnlock()

	t := time.Now()
	year, month, day := t.Date()
	if (len(quoteStockSymbolsToDownLoad) > 0 || len(quoteOptionSymbolsToDownLoad) > 0) &&
		t.Weekday() != time.Saturday &&
		t.Weekday() != time.Sunday &&
		t.After(time.Date(year, month, day, 7, 0, 0, 0, time.Local)) &&
		t.Before(time.Date(year, month, day, 18, 0, 0, 0, time.Local)) {

		startTime := time.Now()
		agent, err := qg.db.GetAgent(qg.killCtx)
		if err != nil {
			logrus.Errorf("QuoteGrabber: could not retrieve agent: %v", err)
			return
		}
		//this could take a VERY long time so after 15 mins we kill it
		ctx, _ := context.WithTimeout(qg.killCtx, 1*time.Minute+30*time.Second)
		stocks, options, err := agent.GetQuotes(ctx, quoteStockSymbolsToDownLoad, quoteOptionSymbolsToDownLoad)
		if err != nil {
			logrus.Errorf("QuoteGrabber: had error while getting quotes: %v", err)
		}

		if len(stocks) > 0 {
			stockQuotesChannel <- stocks
		}

		if len(options) > 0 {
			optionQuotesChannel <- options
		}

		logrus.Infof("QuoteGrabber: finished a round of getting %v of %v stock quotes and %v of %v option quotes in %v", len(stocks), len(quoteStockSymbolsToDownLoad), len(options), len(quoteOptionSymbolsToDownLoad), time.Since(startTime))
	}
}

func (qg *QuoteGrabber) Stop(waitUpTo time.Duration) {
	//next stop the quoter
	logrus.Infof("Waiting for QuoteGrabber to finish")
	startTime := time.Now()
	qg.kill()
	select {
	case <-qg.doneCtx.Done():
		logrus.Infof("QuoteGrabber finished in %v", time.Since(startTime))
	case <-time.After(waitUpTo):
		logrus.Errorf("QuoteGrabber failed to gracefully finish in %v", time.Since(startTime))
	}
}
