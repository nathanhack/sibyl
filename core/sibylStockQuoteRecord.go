package core

import (
	"database/sql"
	"fmt"
	"strings"
)

type SibylStockQuoteRecord struct {
	Ask                 sql.NullFloat64 // both   //
	AskTime             sql.NullInt64   // both   //
	AskSize             sql.NullInt64   // both   //
	Beta                sql.NullFloat64 // stock  //
	Bid                 sql.NullFloat64 // both   //
	BidTime             sql.NullInt64   // both   // millis
	BidSize             sql.NullInt64   // both   //
	Change              sql.NullFloat64 // both   //  (+/-) number
	HighPrice           sql.NullFloat64 // both   //  "hi"`
	LastTradePrice      sql.NullFloat64 // both   //  ally LastTrade ... TD mark or regularMarketLastPrice
	LastTradeTimestamp  sql.NullInt64   // both   //  ally DateTime   TD :tradeTimeInLong
	LastTradeVolume     sql.NullInt64   // both   //  ally incr_vl   TD : lastSize
	LowPrice            sql.NullFloat64 // both   //  "lo"
	Symbol              StockSymbolType // both   //  either the stock symbol or root symbol ex. CAT // this is (and must be) guaranteed to be a NON - Null value
	Timestamp           TimestampType   // both   //  millis // this is the GMT timestamp (ms from epoch) of the this quote // this is (and must be) guaranteed to be a NON - Null value
	Volume              sql.NullInt64   // stock  //
	VolWeightedAvgPrice sql.NullFloat64 // stock  //
}

func (sbq *SibylStockQuoteRecord) String() string {
	return fmt.Sprintf("{%v}", sbq.StringBlindWithDelimiter(",", "", true))
}

func (sbq *SibylStockQuoteRecord) StringBlindWithDelimiter(delimiter string, nullString string, stringEscapes bool) string {
	//strings.Builder is faster fmt. or strings.Join
	esc := ""
	if stringEscapes {
		esc = "'"
	}
	builder := strings.Builder{}
	builder.WriteString(nullFloat64ToString(sbq.Ask, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullInt64ToString(sbq.AskTime, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullInt64ToString(sbq.AskSize, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(sbq.Beta, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(sbq.Bid, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullInt64ToString(sbq.BidTime, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullInt64ToString(sbq.BidSize, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(sbq.Change, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(sbq.HighPrice, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(sbq.LastTradePrice, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullInt64ToString(sbq.LastTradeTimestamp, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullInt64ToString(sbq.LastTradeVolume, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(sbq.LowPrice, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(fmt.Sprintf("%v%v%v", esc, sbq.Symbol, esc))
	builder.WriteString(delimiter)
	builder.WriteString(fmt.Sprintf("%v", sbq.Timestamp))
	builder.WriteString(delimiter)
	builder.WriteString(nullInt64ToString(sbq.Volume, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(sbq.VolWeightedAvgPrice, nullString))
	return builder.String()
}
