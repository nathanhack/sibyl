package core

import (
	"database/sql"
	"fmt"
	"strings"
)

type SibylOptionQuoteRecord struct {
	Ask                sql.NullFloat64 // both   //
	AskTime            sql.NullInt64   // both   //
	AskSize            sql.NullInt64   // both   //
	Bid                sql.NullFloat64 // both   //
	BidTime            sql.NullInt64   // both   // millis
	BidSize            sql.NullInt64   // both   //
	Change             sql.NullFloat64 // both   //  (+/-) number
	Delta              sql.NullFloat64 // option //  "idelta"
	EquityType         EquityType      // both   //  "put_call"     //CALL or PUT // this is (and must be) guaranteed to be a NON - Null value
	Expiration         DateType        // option //  millis // this is (and must be) guaranteed to be a NON - Null value
	Gamma              sql.NullFloat64 // option //  "igamma"
	HighPrice          sql.NullFloat64 // both   //  "hi"`
	ImpliedVolatility  sql.NullFloat64 // option //
	LastTradePrice     sql.NullFloat64 // both   //  ally LastTrade ... TD mark or regularMarketLastPrice
	LastTradeTimestamp sql.NullInt64   // both   //  ally DateTime   TD :tradeTimeInLong
	LastTradeVolume    sql.NullInt64   // both   //  ally incr_vl   TD : lastSize
	LowPrice           sql.NullFloat64 // both   //  "lo"
	OpenInterest       sql.NullInt64   // option //  "openinterest"
	Rho                sql.NullFloat64 // option //  "irho"
	StrikePrice        float64         // option // // this is (and must be) guaranteed to be a NON - Null value
	Symbol             StockSymbolType // both   //  either the stock symbol or root symbol ex. CAT // this is (and must be) guaranteed to be a NON - Null value
	Theta              sql.NullFloat64 // option //  "itheta"
	Timestamp          TimestampType   // both   //  millis // this is the GMT timestamp (ms from epoch) of the this quote // this is (and must be) guaranteed to be a NON - Null value
	Vega               sql.NullFloat64 // option //  "ivega"
}

func (sbq *SibylOptionQuoteRecord) String() string {
	return fmt.Sprintf("{%v}", sbq.StringBlindWithDelimiter(",", "", true))
}

func (sbq *SibylOptionQuoteRecord) StringBlindWithDelimiter(delimiter string, nullString string, stringEscapes bool) string {
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
	builder.WriteString(nullFloat64ToString(sbq.Bid, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullInt64ToString(sbq.BidTime, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullInt64ToString(sbq.BidSize, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(sbq.Change, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(sbq.Delta, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(fmt.Sprintf("%v%v%v", esc, sbq.EquityType, esc))
	builder.WriteString(delimiter)
	builder.WriteString(fmt.Sprintf("%v", sbq.Expiration))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(sbq.Gamma, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(sbq.HighPrice, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(sbq.ImpliedVolatility, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(sbq.LastTradePrice, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullInt64ToString(sbq.LastTradeTimestamp, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullInt64ToString(sbq.LastTradeVolume, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(sbq.LowPrice, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullInt64ToString(sbq.OpenInterest, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(sbq.Rho, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(fmt.Sprintf("%v", sbq.StrikePrice))
	builder.WriteString(delimiter)
	builder.WriteString(fmt.Sprintf("%v%v%v", esc, sbq.Symbol, esc))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(sbq.Theta, nullString))
	builder.WriteString(delimiter)
	builder.WriteString(fmt.Sprintf("%v", sbq.Timestamp))
	builder.WriteString(delimiter)
	builder.WriteString(nullFloat64ToString(sbq.Vega, nullString))
	return builder.String()
}
