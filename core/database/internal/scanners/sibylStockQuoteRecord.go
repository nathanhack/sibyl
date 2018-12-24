package scanners

import (
	"database/sql"
	"fmt"

	"github.com/nathanhack/sibyl/core"
)

func ScanSibylStockQuoteRecordRow(rows *sql.Rows) (string, *core.SibylStockQuoteRecord, error) {
	var id string
	var ask sql.NullFloat64
	var askTime sql.NullInt64
	var askSize sql.NullInt64
	var beta sql.NullFloat64
	var bid sql.NullFloat64
	var bidTime sql.NullInt64
	var bidSize sql.NullInt64
	var change sql.NullFloat64
	var highPrice sql.NullFloat64
	var lastTradePrice sql.NullFloat64
	var lastTradeTimestamp sql.NullInt64
	var lastTradeVolume sql.NullInt64
	var lowPrice sql.NullFloat64
	var stockSymbol core.StockSymbolType
	var timestamp core.TimestampType
	var volume sql.NullInt64
	var volWeightedAvgPrice sql.NullFloat64

	err := rows.Scan(
		&id,
		&ask,
		&askTime,
		&askSize,
		&beta,
		&bid,
		&bidTime,
		&bidSize,
		&change,
		&highPrice,
		&lastTradePrice,
		&lastTradeTimestamp,
		&lastTradeVolume,
		&lowPrice,
		&stockSymbol,
		&timestamp,
		&volume,
		&volWeightedAvgPrice,
	)

	if err != nil {
		return "", nil, fmt.Errorf("ScanSibylStockQuoteRecordRow: had an error while scanning: %v", err)
	}

	quote := core.SibylStockQuoteRecord{
		Ask:                 ask,
		AskTime:             askTime,
		AskSize:             askSize,
		Beta:                beta,
		Bid:                 bid,
		BidTime:             bidTime,
		BidSize:             bidSize,
		Change:              change,
		HighPrice:           highPrice,
		LastTradePrice:      lastTradePrice,
		LastTradeVolume:     lastTradeVolume,
		LastTradeTimestamp:  lastTradeTimestamp,
		LowPrice:            lowPrice,
		Symbol:              stockSymbol,
		Timestamp:           timestamp,
		Volume:              volume,
		VolWeightedAvgPrice: volWeightedAvgPrice,
	}
	return id, &quote, nil
}
