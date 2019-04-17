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
	var intradayHistoryState core.ScanStateType
	var exchange string
	var exchangeName string
	var hasOptions string
	var name string
	var quotesStatus core.ActivityStatusType
	var stableQuotesStatus core.ActivityStatusType
	var stockSymbol core.StockSymbolType
	var validationStatus core.ValidationStatusType
	var validationTimestamp core.DateType

	err := rows.Scan(&id,
		&downloadStatus,
		&exchange,
		&exchangeName,
		&hasOptions,
		&historyStatus,
		&intradayHistoryStatus,
		&intradayHistoryState,
		&name,
		&quotesStatus,
		&stableQuotesStatus,
		&stockSymbol,
		&validationStatus,
		&validationTimestamp,
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
		IntradayHistoryState:  intradayHistoryState,
		Exchange:              exchange,
		ExchangeDescription:   exchangeName,
		HasOptions:            options,
		Name:                  name,
		QuotesStatus:          quotesStatus,
		StableQuotesStatus:    stableQuotesStatus,
		Symbol:                stockSymbol,
		ValidationStatus:      validationStatus,
		ValidationTimestamp:   validationTimestamp,
	}
	return toReturn, nil
}
