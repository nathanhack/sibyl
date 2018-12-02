package core

import (
	"database/sql"
	"fmt"
	"strings"
)

type SibylStableStockQuoteRecord struct {
	AnnualDividend         sql.NullFloat64  // stock  // ally : iad   TD   divAmount
	BookValue              sql.NullFloat64  // stock  //  "prbook"
	ClosePrice             sql.NullFloat64  // both   // previous day close
	Div                    sql.NullFloat64  // stock  // ally only //(= AnnualDiv/ divFreq) Latest announced cash dividend
	DivExTimestamp         sql.NullInt64    // stock  // millis  // date of last dividend
	DivFreq                NullDivFrequency // stock  // Ally
	DivPayTimestamp        sql.NullInt64    // stock  // millis // Ally
	Eps                    sql.NullFloat64  // stock  // "eps"
	HighPrice52Wk          sql.NullFloat64  // both   //
	HighPrice52WkTimestamp sql.NullInt64    // both   // ally only //
	LowPrice52Wk           sql.NullFloat64  // both   //
	LowPrice52WkTimestamp  sql.NullInt64    // both   // ally only //millis -- ally only
	OpenPrice              sql.NullFloat64  // both   //  "opn"
	PriceEarnings          sql.NullFloat64  // stock  //  "pe"
	SharesOutstanding      sql.NullInt64    // stock  //
	Symbol                 StockSymbolType  // both   //  either the stock symbol or root symbol ex. CAT // this is (and must be) guaranteed to be a NON - Null value
	Timestamp              DateType         // both   // millis // Since these values should be stable the date is just the month/day/year // this is (and must be) guaranteed to be a NON - Null value
	Volatility             sql.NullFloat64  // stock  // one year volatility measure
	Yield                  sql.NullFloat64  // stock  //
}

func (ssq *SibylStableStockQuoteRecord) String() string {
	return fmt.Sprintf("{%v}", ssq.StringBlindWithDelimiter(",", "", true))
}

func (ssq *SibylStableStockQuoteRecord) StringBlindWithDelimiter(delimiter string, nullString string, stringEscapes bool) string {
	//strings.Builder is faster fmt. or strings.Join
	esc := ""
	if stringEscapes {
		esc = "'"
	}
	builder := strings.Builder{}
	builder.WriteString(nullFloat64ToString(ssq.AnnualDividend, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(ssq.BookValue, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(ssq.ClosePrice, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(ssq.Div, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullInt64ToString(ssq.DivExTimestamp, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullDivFrequencyToString(ssq.DivFreq, nullString, stringEscapes))
	builder.WriteString(delimiter)
	builder.WriteString(nullInt64ToString(ssq.DivPayTimestamp, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(ssq.Eps, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(ssq.HighPrice52Wk, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullInt64ToString(ssq.HighPrice52WkTimestamp, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(ssq.LowPrice52Wk, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullInt64ToString(ssq.LowPrice52WkTimestamp, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(ssq.OpenPrice, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(ssq.PriceEarnings, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullInt64ToString(ssq.SharesOutstanding, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(fmt.Sprintf("%v%v%v", esc, ssq.Symbol, esc))
	builder.WriteString(delimiter)
	builder.WriteString(fmt.Sprintf("%v", ssq.Timestamp))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(ssq.Volatility, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(ssq.Yield, nullString))
	return builder.String()
}
