package rest

type Stock struct {
	Exchange              string `json:"Exchange"`
	ExchangeDescription   string `json:"ExchangeDescription"`
	DownloadStatus        string `json:"DownloadStatus"`
	HasOptions            bool   `json:"HasOptions"`
	HistoryStatus         string `json:"HistoryStatus"`
	IntradayHistoryStatus string `json:"IntradayHistoryStatus"`
	IntradayHistoryState  string `json:"IntradayHistoryState"`
	Name                  string `json:"Name"`
	QuotesStatus          string `json:"QuotesStatus"`
	StableQuotesStatus    string `json:"StableQuotesStatus"`
	Symbol                string `json:"Symbol"`
	Validation            string `json:"Validation"`
}

type StocksResponse struct {
	Stocks     []Stock    `json:"Stocks"`
	ErrorState ErrorState `json:"ErrorState"`
}
