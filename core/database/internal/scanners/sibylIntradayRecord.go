package scanners

import (
	"database/sql"
	"fmt"
	"github.com/nathanhack/sibyl/core"
)

func ScanSibylIntradayRecordRow(rows *sql.Rows) (string, *core.SibylIntradayRecord, error) {
	var id string
	var HighPrice sql.NullFloat64
	var LastPrice sql.NullFloat64
	var LowPrice sql.NullFloat64
	var OpenPrice sql.NullFloat64
	var Symbol core.StockSymbolType
	var Timestamp core.TimestampType
	var Volume sql.NullInt64

	err := rows.Scan(
		&id,
		&HighPrice,
		&LastPrice,
		&LowPrice,
		&OpenPrice,
		&Symbol,
		&Timestamp,
		&Volume,
	)

	if err != nil {
		return "", nil, fmt.Errorf("ScanSibylIntradayRecordRow: had an error while scanning: %v", err)
	}

	return id, &core.SibylIntradayRecord{
		HighPrice: HighPrice,
		LastPrice: LastPrice,
		LowPrice:  LowPrice,
		OpenPrice: OpenPrice,
		Symbol:    Symbol,
		Timestamp: Timestamp,
		Volume:    Volume,
	}, nil
}
