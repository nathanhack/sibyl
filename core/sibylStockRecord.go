package core

import (
	"fmt"
)

type SibylStockRecord struct {
	DownloadStatus        ActivityStatusType
	Exchange              string
	ExchangeDescription   string
	HasOptions            bool
	HistoryStatus         ActivityStatusType
	IntradayHistoryStatus ActivityStatusType
	IntradayHistoryState  ScanStateType
	Name                  string
	QuotesStatus          ActivityStatusType
	StableQuotesStatus    ActivityStatusType
	Symbol                StockSymbolType
	ValidationStatus      ValidationStatusType
}

func (ss *SibylStockRecord) String() string {
	return fmt.Sprintf("{%v}", ss.StringBlindWithDelimiter(",", "", true))
}

func (ss *SibylStockRecord) StringBlindWithDelimiter(delimiter string, nullString string, stringEscapes bool) string {
	esc := ""
	if stringEscapes {
		esc = "'"
	}
	hasOptions := "no"
	if ss.HasOptions {
		hasOptions = "yes"
	}
	return fmt.Sprintf("%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%%v%v%v%v%v%v%v%v%v%v%vv%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v",
		esc,
		ss.DownloadStatus,
		esc,
		delimiter,
		esc,
		ss.Exchange,
		esc,
		delimiter,
		esc,
		ss.ExchangeDescription,
		esc,
		delimiter,
		esc,
		hasOptions,
		esc,
		delimiter,
		esc,
		ss.HistoryStatus,
		esc,
		delimiter,
		esc,
		ss.IntradayHistoryStatus,
		esc,
		delimiter,
		esc,
		ss.IntradayHistoryState,
		esc,
		delimiter,
		esc,
		ss.Name,
		esc,
		delimiter,
		esc,
		ss.QuotesStatus,
		esc,
		delimiter,
		esc,
		ss.StableQuotesStatus,
		esc,
		delimiter,
		esc,
		ss.Symbol,
		esc,
		delimiter,
		esc,
		ss.ValidationStatus,
		esc,
	)
}
