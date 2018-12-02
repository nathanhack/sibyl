package rest

type Histories struct {
	Histories  []History  `json:"History"`
	ErrorState ErrorState `json:"ErrorState"`
}

type History struct {
	ClosePrice string `json:"ClosePrice"`
	HighPrice  string `json:"HighPrice"`
	LowPrice   string `json:"LowPrice"`
	OpenPrice  string `json:"OpenPrice"`
	Symbol     string `json:"Symbol"`
	Timestamp  int64  `json:"Timestamp"`
	Volume     string `json:"Volume"`
}
