package rest

type DatabaseRecords struct {
	Histories          string     `json:"Histories"`
	Intradays          string     `json:"Intradays"`
	StockQuotes        string     `json:"StockQuotes"`
	StockStableQuotes  string     `json:"StockStableQuotes"`
	OptionQuotes       string     `json:"OptionQuotes"`
	OptionStableQuotes string     `json:"OptionStableQuotes"`
	LastID             string     `json:"LastID"`
	More               bool       `json:"More"`
	ErrorState         ErrorState `json:"ErrorState"`
}
