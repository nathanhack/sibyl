package scanners

import (
	"database/sql"
	"fmt"
	"github.com/nathanhack/sibyl/core"
)

func ScanSibylStockRecordRow(rows *sql.Rows) (*core.SibylStockRecord, error) {
	var downloadStatus core.ActivityStatusType
	var historyStatus core.HistoryStatusType
	var historyTimestamp core.DateType
	var intradayState core.IntradayStateType
	var intradayStatus core.IntradayStatusType
	var intradayTimestamp1Min core.TimestampType
	var intradayTimestamp5Min core.TimestampType
	var intradayTimestampTick core.TimestampType
	var exchange string
	var exchangeName string
	var hasOptions core.OptionStatusType
	var name string
	var optionListTimestamp core.DateType
	var quotesStatus core.ActivityStatusType
	var stableQuotesStatus core.ActivityStatusType
	var stockSymbol core.StockSymbolType
	var validationStatus core.ValidationStatusType
	var validationTimestamp core.DateType

	err := rows.Scan(
		&downloadStatus,
		&exchange,
		&exchangeName,
		&hasOptions,
		&historyStatus,
		&historyTimestamp,
		&intradayState,
		&intradayStatus,
		&intradayTimestamp1Min,
		&intradayTimestamp5Min,
		&intradayTimestampTick,
		&name,
		&optionListTimestamp,
		&quotesStatus,
		&stableQuotesStatus,
		&stockSymbol,
		&validationStatus,
		&validationTimestamp,
	)

	if err != nil {
		return nil, fmt.Errorf("ScanSibylStockRow: had a problem while scanning: %v", err)
	}

	toReturn := &core.SibylStockRecord{
		DownloadStatus:        downloadStatus,
		HistoryStatus:         historyStatus,
		HistoryTimestamp:      historyTimestamp,
		IntradayState:         intradayState,
		IntradayStatus:        intradayStatus,
		IntradayTimestamp1Min: intradayTimestamp1Min,
		IntradayTimestamp5Min: intradayTimestamp5Min,
		IntradayTimestampTick: intradayTimestampTick,
		Exchange:              exchange,
		ExchangeDescription:   exchangeName,
		OptionStatus:          hasOptions,
		Name:                  name,
		OptionListTimestamp:   optionListTimestamp,
		QuotesStatus:          quotesStatus,
		StableQuotesStatus:    stableQuotesStatus,
		Symbol:                stockSymbol,
		ValidationStatus:      validationStatus,
		ValidationTimestamp:   validationTimestamp,
	}
	return toReturn, nil
}
