package core

import (
	"database/sql"
	"fmt"
	"strings"
)

type SibylHistoryRecord struct {
	ClosePrice sql.NullFloat64
	HighPrice  sql.NullFloat64
	Interval   HistoryStatusType
	LowPrice   sql.NullFloat64
	OpenPrice  sql.NullFloat64
	Symbol     StockSymbolType // this is (and must be) guaranteed to be a NON - Null value
	Timestamp  DateType        // millis // this is (and must be) guaranteed to be a NON - Null value
	Volume     sql.NullInt64
}

func (shr *SibylHistoryRecord) String() string {
	return shr.StringBlindWithDelimiter(",", "", true)
}

func (shr *SibylHistoryRecord) StringBlindWithDelimiter(delimiter string, nullString string, stringEscapes bool) string {
	//strings.Builder is faster fmt. or strings.Join
	esc := ""
	if stringEscapes {
		esc = "'"
	}
	builder := strings.Builder{}
	builder.WriteString(nullFloat64ToString(shr.ClosePrice, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(shr.HighPrice, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(fmt.Sprintf("%v", int(shr.Interval)))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(shr.LowPrice, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(shr.OpenPrice, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(fmt.Sprintf("%v%v%v", esc, shr.Symbol, esc))
	builder.WriteString(delimiter)
	builder.WriteString(fmt.Sprintf("%v", shr.Timestamp))
	builder.WriteString(delimiter)
	builder.WriteString(nullInt64ToString(shr.Volume, nullString))
	return builder.String()
}
