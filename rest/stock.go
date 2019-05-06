package rest

import "github.com/go-humble/rest"

type Stock struct {
	Exchange              string `json:"Exchange"`
	ExchangeDescription   string `json:"ExchangeDescription"`
	DownloadStatus        int    `json:"DownloadStatus"`
	HistoryStatus         int    `json:"HistoryStatus"`
	HistoryTimestamp      int64  `json:"HistoryTimstamp"`
	IntradayStatus        int    `json:"IntradayStatus"`
	IntradayState         int    `json:"IntradayState"`
	IntradayTimestamp1Min int64  `json:"IntradayTimestamp1Min"`
	IntradayTimestamp5Min int64  `json:"IntradayTimestamp5Min"`
	IntradayTimestampTick int64  `json:"IntradayTimestampTick"`
	Name                  string `json:"Name"`
	OptionListTimestamp   int64  `json:"OptionListTimestamp"`
	OptionStatus          int    `json:"OptionStatus"`
	QuotesStatus          int    `json:"QuotesStatus"`
	StableQuotesStatus    int    `json:"StableQuotesStatus"`
	Symbol                string `json:"Symbol"`
	Validation            int    `json:"Validation"`
	ValidationTimestamp   int64  `json:"ValidationTimestamp"`
}

type StocksResponse struct {
	rest.DefaultId
	Stocks     []Stock    `json:"Stocks"`
	ErrorState ErrorState `json:"ErrorState"`
}

func (sr *StocksResponse) RootURL() string {
	return "/stocks/get"
}
