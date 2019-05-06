package core

import (
	"fmt"
)

type SibylStockRecord struct {
	DownloadStatus        ActivityStatusType
	Exchange              string
	ExchangeDescription   string
	HistoryStatus         HistoryStatusType
	HistoryTimestamp      DateType
	IntradayState         IntradayStateType
	IntradayStatus        IntradayStatusType
	IntradayTimestamp1Min TimestampType
	IntradayTimestamp5Min TimestampType
	IntradayTimestampTick TimestampType
	Name                  string
	OptionListTimestamp   DateType
	OptionStatus          OptionStatusType
	QuotesStatus          ActivityStatusType
	StableQuotesStatus    ActivityStatusType
	Symbol                StockSymbolType
	ValidationStatus      ValidationStatusType
	ValidationTimestamp   DateType
}

func (ss *SibylStockRecord) String() string {
	return fmt.Sprintf("{%v}", ss.StringBlindWithDelimiter(",", "", true))
}

func (ss *SibylStockRecord) StringBlindWithDelimiter(delimiter string, nullString string, stringEscapes bool) string {
	esc := ""
	if stringEscapes {
		esc = "'"
	}
	return fmt.Sprintf("%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%%v%v%v%v%v%v%v%v%v%v%vv%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v",
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
		ss.HistoryStatus,
		esc,
		delimiter,
		esc,
		ss.HistoryTimestamp,
		esc,
		delimiter,
		esc,
		ss.IntradayState,
		esc,
		delimiter,
		esc,
		ss.IntradayStatus,
		esc,
		delimiter,
		esc,
		ss.IntradayTimestamp1Min,
		esc,
		delimiter,
		esc,
		ss.IntradayTimestamp5Min,
		esc,
		delimiter,
		esc,
		ss.IntradayTimestampTick,
		esc,
		delimiter,
		esc,
		ss.Name,
		esc,
		delimiter,
		esc,
		ss.OptionListTimestamp,
		esc,
		delimiter,
		esc,
		ss.OptionStatus,
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
		delimiter,
		esc,
		ss.ValidationTimestamp,
		esc,
	)
}
