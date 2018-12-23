package scanners

import (
	"database/sql"
	"fmt"
	"github.com/nathanhack/sibyl/core"
)

func ScanSibylStableOptionQuoteRecordRow(rows *sql.Rows) (string, *core.SibylStableOptionQuoteRecord, error) {
	var id string
	var closePrice sql.NullFloat64
	var contractSize sql.NullInt64
	var equityType core.EquityType
	var expirationTimestamp core.DateType
	var highPrice52Wk sql.NullFloat64
	var highPrice52WkTimestamp sql.NullInt64
	var lowPrice52Wk sql.NullFloat64
	var lowPrice52WkTimestamp sql.NullInt64
	var multiplier sql.NullInt64
	var openPrice sql.NullFloat64
	var strikePrice float64
	var symbol core.StockSymbolType
	var timestamp core.DateType

	err := rows.Scan(
		&id,
		&closePrice,
		&contractSize,
		&equityType,
		&expirationTimestamp,
		&highPrice52Wk,
		&highPrice52WkTimestamp,
		&lowPrice52Wk,
		&lowPrice52WkTimestamp,
		&multiplier,
		&openPrice,
		&strikePrice,
		&symbol,
		&timestamp,
	)

	if err != nil {
		return "", nil, fmt.Errorf("ScanSibylStableOptionQuoteRecordRow: had a error while scanning: %v", err)
	}

	quote := core.SibylStableOptionQuoteRecord{
		ClosePrice:             closePrice,
		ContractSize:           contractSize,
		EquityType:             equityType,
		Expiration:             expirationTimestamp,
		HighPrice52Wk:          highPrice52Wk,
		HighPrice52WkTimestamp: highPrice52WkTimestamp,
		LowPrice52Wk:           lowPrice52Wk,
		LowPrice52WkTimestamp:  lowPrice52WkTimestamp,
		Multiplier:             multiplier,
		OpenPrice:              openPrice,
		Symbol:                 symbol,
		StrikePrice:            strikePrice,
		Timestamp:              timestamp,
	}
	return id, &quote, nil
}
