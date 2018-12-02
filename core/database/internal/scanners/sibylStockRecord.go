package scanners

import (
	"database/sql"
	"fmt"
	"github.com/nathanhack/sibyl/core"
)

func ScanSibylStockRecordRow(rows *sql.Rows) (*core.SibylStockRecord, error) {
	var id int
	var downloadStatus core.ActivityStatusType
	var historyStatus core.ActivityStatusType
	var intradayHistoryStatus core.ActivityStatusType
	var exchange string
	var exchangeName string
	var hasOptions string
	var name string
	var quotesStatus core.ActivityStatusType
	var stableQuotesStatus core.ActivityStatusType
	var stockSymbol core.StockSymbolType
	var validationStatus core.ValidationStatusType

	err := rows.Scan(&id,
		&downloadStatus,
		&exchange,
		&exchangeName,
		&hasOptions,
		&historyStatus,
		&intradayHistoryStatus,
		&name,
		&quotesStatus,
		&stableQuotesStatus,
		&stockSymbol,
		&validationStatus,
	)

	if err != nil {
		return nil, fmt.Errorf("ScanSibylStockRow: had a problem while scanning: %v", err)
	}
	options := false
	if hasOptions == "yes" {
		options = true
	}
	toReturn := &core.SibylStockRecord{
		DownloadStatus:        downloadStatus,
		HistoryStatus:         historyStatus,
		IntradayHistoryStatus: intradayHistoryStatus,
		Exchange:              exchange,
		ExchangeDescription:   exchangeName,
		HasOptions:            options,
		Name:                  name,
		QuotesStatus:          quotesStatus,
		StableQuotesStatus:    stableQuotesStatus,
		Symbol:                stockSymbol,
		ValidationStatus:      validationStatus,
	}
	return toReturn, nil
}
