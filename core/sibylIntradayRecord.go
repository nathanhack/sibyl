package core

import (
	"database/sql"
	"fmt"
	"strings"
)

type SibylIntradayRecord struct {
	HighPrice sql.NullFloat64
	LastPrice sql.NullFloat64
	LowPrice  sql.NullFloat64
	OpenPrice sql.NullFloat64
	Symbol    StockSymbolType // this is (and must be) guaranteed to be a NON - Null value
	Timestamp TimestampType   // this is (and must be) guaranteed to be a NON - Null value
	Volume    sql.NullInt64
}

func (si *SibylIntradayRecord) String() string {
	return fmt.Sprintf("{%v}", si.StringBlindWithDelimiter(",", "", true))
}

func (si *SibylIntradayRecord) StringBlindWithDelimiter(delimiter string, nullString string, stringEscapes bool) string {
	esc := ""
	if stringEscapes {
		esc = "'"
	}
	builder := strings.Builder{}
	builder.WriteString(nullFloat64ToString(si.HighPrice, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(si.LastPrice, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(si.LowPrice, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(si.OpenPrice, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(fmt.Sprintf("%v%v%v", esc, si.Symbol, esc))
	builder.WriteString(delimiter)
	builder.WriteString(fmt.Sprintf("%v", si.Timestamp))
	builder.WriteString(delimiter)
	builder.WriteString(nullInt64ToString(si.Volume, nullString))
	return builder.String()
}
