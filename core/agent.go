package core

import (
	"context"
)

type HistoryInterval string

type IntradayInterval string

const (
	DailyInterval   HistoryInterval = "daily"
	WeeklyInterval  HistoryInterval = "weekly"
	MonthlyInterval HistoryInterval = "monthly"
	YearlyInterval  HistoryInterval = "yearly"

	OneMinInterval  IntradayInterval = "1min"
	FiveMinInterval IntradayInterval = "5min"
	TickInterval    IntradayInterval = "tick"
)

type SibylAgent interface {
	GetHistory(ctx context.Context, symbol StockSymbolType, interval HistoryInterval, startDate, endDate DateType) ([]*SibylHistoryRecord, error)
	GetIntraday(ctx context.Context, symbol StockSymbolType, interval IntradayInterval, startDate, endDate TimestampType) ([]*SibylIntradayRecord, error)
	GetQuotes(ctx context.Context, stockSymbols map[StockSymbolType]bool, optionSymbols map[OptionSymbolType]bool) ([]*SibylStockQuoteRecord, []*SibylOptionQuoteRecord, error)
	GetStableQuotes(ctx context.Context, stockSymbols map[StockSymbolType]bool, optionSymbols map[OptionSymbolType]bool) ([]*SibylStableStockQuoteRecord, []*SibylStableOptionQuoteRecord, error)
	GetStockOptionSymbols(ctx context.Context, symbol StockSymbolType) ([]*OptionSymbolType, error)
	VerifyStockSymbol(ctx context.Context, symbol StockSymbolType) (good, hasOptions bool, exchange, exchangeName, name string, err error)
}
