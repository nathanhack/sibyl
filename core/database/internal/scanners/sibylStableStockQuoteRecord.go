package scanners

import (
	"database/sql"
	"fmt"
	"github.com/nathanhack/sibyl/core"
)

func ScanSibylStableStockQuoteRecordRow(rows *sql.Rows) (*core.SibylStableStockQuoteRecord, error) {
	var id int
	var annualDividend sql.NullFloat64
	var bookValue sql.NullFloat64
	var closePrice sql.NullFloat64
	var div sql.NullFloat64
	var divFreq core.NullDivFrequency
	var divExTimestamp sql.NullInt64
	var divPayTimestamp sql.NullInt64
	var eps sql.NullFloat64
	var highPrice52Wk sql.NullFloat64
	var highPrice52WkTimestamp sql.NullInt64
	var lowPrice52Wk sql.NullFloat64
	var lowPrice52WkTimestamp sql.NullInt64
	var openPrice sql.NullFloat64
	var priceEarnings sql.NullFloat64
	var sharesOutstanding sql.NullInt64
	var symbol core.StockSymbolType
	var timestamp core.DateType
	var volatility sql.NullFloat64
	var yield sql.NullFloat64

	err := rows.Scan(
		&id,
		&annualDividend,
		&bookValue,
		&closePrice,
		&div,
		&divFreq,
		&divExTimestamp,
		&divPayTimestamp,
		&eps,
		&highPrice52Wk,
		&highPrice52WkTimestamp,
		&lowPrice52Wk,
		&lowPrice52WkTimestamp,
		&openPrice,
		&priceEarnings,
		&sharesOutstanding,
		&symbol,
		&timestamp,
		&volatility,
		&yield,
	)

	if err != nil {
		return nil, fmt.Errorf("ScanSibylStableStockQuoteRecordRow: had a error while scanning: %v", err)
	}

	quote := core.SibylStableStockQuoteRecord{
		AnnualDividend:         annualDividend,
		BookValue:              bookValue,
		ClosePrice:             closePrice,
		Div:                    div,
		DivFreq:                divFreq,
		DivExTimestamp:         divExTimestamp,
		DivPayTimestamp:        divPayTimestamp,
		Eps:                    eps,
		HighPrice52Wk:          highPrice52Wk,
		HighPrice52WkTimestamp: highPrice52WkTimestamp,
		LowPrice52Wk:           lowPrice52Wk,
		LowPrice52WkTimestamp:  lowPrice52WkTimestamp,
		OpenPrice:              openPrice,
		PriceEarnings:          priceEarnings,
		SharesOutstanding:      sharesOutstanding,
		Symbol:                 symbol,
		Timestamp:              timestamp,
		Volatility:             volatility,
		Yield:                  yield,
	}
	return &quote, nil
}
