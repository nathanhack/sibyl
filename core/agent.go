package core

import (
	"context"
)

type HistoryTick string

const (
	MinuteTicks     HistoryTick = "1min"
	FiveMinuteTicks HistoryTick = "5min"
	DailyTicks      HistoryTick = "day"
)

type SibylAgent interface {
	GetHistory(ctx context.Context, symbol StockSymbolType, tickSize HistoryTick, startDate, endDate DateType) ([]*SibylHistoryRecord, error)
	GetIntraday(ctx context.Context, symbol StockSymbolType, tickSize HistoryTick, startDate, endDate TimestampType) ([]*SibylIntradayRecord, error)
	GetQuotes(ctx context.Context, stockSymbols map[StockSymbolType]bool, optionSymbols map[OptionSymbolType]bool) ([]*SibylStockQuoteRecord, []*SibylOptionQuoteRecord, error)
	GetStableQuotes(ctx context.Context, stockSymbols map[StockSymbolType]bool, optionSymbols map[OptionSymbolType]bool) ([]*SibylStableStockQuoteRecord, []*SibylStableOptionQuoteRecord, error)
	GetStockOptionSymbols(ctx context.Context, symbol StockSymbolType) ([]*OptionSymbolType, error)
	VerifyStockSymbol(ctx context.Context, symbol StockSymbolType) (good, hasOptions bool, exchange, exchangeName, name string, err error)
}
