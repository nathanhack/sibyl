package core

import (
	"database/sql"
	"fmt"
	"strings"
)

type SibylStableOptionQuoteRecord struct {
	ClosePrice             sql.NullFloat64 // both   // previous day close
	ContractSize           sql.NullInt64   // option // contract size for option
	EquityType             EquityType      // both   //  "put_call"    //CALL or PUT
	Expiration             DateType        // option // millis // this is (and must be) guaranteed to be a NON - Null value
	HighPrice52Wk          sql.NullFloat64 // both   //
	HighPrice52WkTimestamp sql.NullInt64   // both   // ally only //
	LowPrice52Wk           sql.NullFloat64 // both   //
	LowPrice52WkTimestamp  sql.NullInt64   // both   // ally only //millis -- ally only
	Multiplier             sql.NullInt64   // option // "prem_mult"
	OpenPrice              sql.NullFloat64 // both   //  "opn"
	StrikePrice            float64         // option // // this is (and must be) guaranteed to be a NON - Null value
	Symbol                 StockSymbolType // both   //  either the stock symbol or root symbol ex. CAT // this is (and must be) guaranteed to be a NON - Null value
	Timestamp              DateType        // both   // millis // Since these values should be stable the date is just the month/day/year // this is (and must be) guaranteed to be a NON - Null value
}

func (ssq *SibylStableOptionQuoteRecord) String() string {
	return fmt.Sprintf("{%v}", ssq.StringBlindWithDelimiter(",", "", true))
}

func (ssq *SibylStableOptionQuoteRecord) StringBlindWithDelimiter(delimiter string, nullString string, stringEscapes bool) string {
	//strings.Builder is faster fmt. or strings.Join
	esc := ""
	if stringEscapes {
		esc = "'"
	}
	builder := strings.Builder{}
	builder.WriteString(nullFloat64ToString(ssq.ClosePrice, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullInt64ToString(ssq.ContractSize, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(fmt.Sprintf("%v%v%v", esc, ssq.EquityType, esc))
	builder.WriteString(delimiter)
	builder.WriteString(fmt.Sprintf("%v", ssq.Expiration))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(ssq.HighPrice52Wk, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullInt64ToString(ssq.HighPrice52WkTimestamp, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(ssq.LowPrice52Wk, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullInt64ToString(ssq.LowPrice52WkTimestamp, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullInt64ToString(ssq.Multiplier, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(ssq.OpenPrice, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(fmt.Sprintf("%v", ssq.StrikePrice))
	builder.WriteString(delimiter)
	builder.WriteString(fmt.Sprintf("%v%v%v", esc, ssq.Symbol, esc))
	builder.WriteString(delimiter)
	builder.WriteString(fmt.Sprintf("%v", ssq.Timestamp))
	return builder.String()
}
