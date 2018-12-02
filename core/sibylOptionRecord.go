package core

import (
	"fmt"
)

type SibylOptionRecord struct {
	Expiration  DateType
	OptionType  EquityType
	StrikePrice float64
	Symbol      StockSymbolType
}

func (so *SibylOptionRecord) String() string {
	return "{" + so.StringBlindWithDelimiter(",", "", true) + "}"
}

func (so *SibylOptionRecord) StringBlindWithDelimiter(delimiter string, nullString string, stringEscapes bool) string {
	esc := ""
	if stringEscapes {
		esc = "'"
	}
	return fmt.Sprintf("%v%v%v%v%v%v%v%v%v%v%v",
		so.Expiration,
		delimiter,
		esc,
		so.OptionType,
		esc,
		delimiter,
		so.StrikePrice,
		delimiter,
		esc,
		so.Symbol,
		esc,
	)

}
