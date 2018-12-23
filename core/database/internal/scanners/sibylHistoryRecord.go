package scanners

import (
	"database/sql"
	"fmt"
	"github.com/nathanhack/sibyl/core"
)

func ScanSibylHistoryRecordRow(rows *sql.Rows) (string,*core.SibylHistoryRecord, error) {
	var id string
	var closePrice sql.NullFloat64
	var highPrice sql.NullFloat64
	var lowPrice sql.NullFloat64
	var openPrice sql.NullFloat64
	var symbol core.StockSymbolType
	var timestamp core.DateType // millis
	var volume sql.NullInt64

	err := rows.Scan(
		&id,
		&closePrice,
		&highPrice,
		&lowPrice,
		&openPrice,
		&symbol,
		&timestamp,
		&volume,
	)

	if err != nil {
		return "",nil, fmt.Errorf("ScanSibylHistoryRecordRow: had an error while scanning: %v", err)
	}

	return id,&core.SibylHistoryRecord{
		ClosePrice: closePrice,
		HighPrice:  highPrice,
		LowPrice:   lowPrice,
		OpenPrice:  openPrice,
		Symbol:     symbol,
		Timestamp:  timestamp,
		Volume:     volume,
	}, nil
}
