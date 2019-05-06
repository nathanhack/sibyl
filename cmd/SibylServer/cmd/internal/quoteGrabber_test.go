package internal

import (
	"context"
	"fmt"
	"github.com/nathanhack/sibyl/core"
	"github.com/nathanhack/sibyl/core/database"
	"testing"
)

const mysqlAddress = "localhost:3306"

func Experiement_QuoteGrabberRates(t *testing.T) {
	fmt.Println("setup starting")
	db, err := database.ConnectAndEnsureSibylDatabase(context.Background(), mysqlAddress)
	if err != nil {
		t.Errorf("Could not establish connection to database: %v", err)
		return
	}

	symbolCache := NewStockCache(db)
	if err := symbolCache.updatedSymbolsList(); err != nil {
		t.Errorf("expected no error for symbol updated: %v", err)
		return
	}

	quoteGrabber := NewQuoteGrabber(db, symbolCache)
	fmt.Println("xml(7k):withfields")
	stockQuotesChannel := make(chan []*core.SibylStockQuoteRecord, 800)
	optionQuotesChannel := make(chan []*core.SibylOptionQuoteRecord, 800)
	fmt.Println("setup done now executing test function")
	quoteGrabber.executeOneRound(stockQuotesChannel, optionQuotesChannel)
	stockQuotesChannel = make(chan []*core.SibylStockQuoteRecord, 800)
	optionQuotesChannel = make(chan []*core.SibylOptionQuoteRecord, 800)
	quoteGrabber.executeOneRound(stockQuotesChannel, optionQuotesChannel)
	stockQuotesChannel = make(chan []*core.SibylStockQuoteRecord, 800)
	optionQuotesChannel = make(chan []*core.SibylOptionQuoteRecord, 800)
	quoteGrabber.executeOneRound(stockQuotesChannel, optionQuotesChannel)
	stockQuotesChannel = make(chan []*core.SibylStockQuoteRecord, 800)
	optionQuotesChannel = make(chan []*core.SibylOptionQuoteRecord, 800)
	quoteGrabber.executeOneRound(stockQuotesChannel, optionQuotesChannel)

}
