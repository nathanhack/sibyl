package scanners

import (
	"database/sql"
	"fmt"
	"github.com/nathanhack/sibyl/core"
)

func ScanSibylOptionRecordRow(rows *sql.Rows) (*core.SibylOptionRecord, error) {
	var id int
	var expiration core.DateType
	var optionType core.EquityType
	var strikePrice float64
	var symbol core.StockSymbolType

	err := rows.Scan(
		&id,
		&expiration,
		&optionType,
		&strikePrice,
		&symbol,
	)
	if err != nil {
		return nil, fmt.Errorf("ScanSibylOptionRecordRow: had an error while scanning: %v", err)
	}

	return &core.SibylOptionRecord{
		Expiration:  expiration,
		OptionType:  optionType,
		StrikePrice: strikePrice,
		Symbol:      symbol,
	}, nil

}
