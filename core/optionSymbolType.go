package core

type OptionSymbolType struct {
	Expiration  DateType
	OptionType  EquityType
	StrikePrice float64
	Symbol      StockSymbolType
}
