package rest

type DatabaseRecords struct {
	Histories          string     `json:"Histories"`
	Intradays          string     `json:"Intradays"`
	StockQuotes        string     `json:"StockQuotes"`
	StockStableQuotes  string     `json:"StockStableQuotes"`
	OptionQuotes       string     `json:"OptionQuotes"`
	OptionStableQuotes string     `json:"OptionStableQuotes"`
	ErrorState         ErrorState `json:"ErrorState"`
}
