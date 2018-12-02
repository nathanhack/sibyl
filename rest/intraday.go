package rest

type Intradays struct {
	Intradays  []Intraday `json:"Intraday"`
	ErrorState ErrorState `json:"ErrorState"`
}

type Intraday struct {
	HighPrice string `json:"HighPrice"`
	LastPrice string `json:"LastPrice"`
	LowPrice  string `json:"LowPrice"`
	OpenPrice string `json:"OpenPrice"`
	Symbol    string `json:"Symbol"`
	Timestamp int64  `json:"Timestamp"`
	Volume    string `json:"Volume"`
}
