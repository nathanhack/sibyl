package scanners

import (
	"database/sql"
	"fmt"
	"github.com/nathanhack/sibyl/core"
)

func ScanSibylOptionQuoteRecordRow(rows *sql.Rows) (string,*core.SibylOptionQuoteRecord, error) {
	var id string
	var ask sql.NullFloat64
	var askTime sql.NullInt64
	var askSize sql.NullInt64
	var bid sql.NullFloat64
	var bidTime sql.NullInt64
	var bidSize sql.NullInt64
	var change sql.NullFloat64
	var delta sql.NullFloat64
	var equityType core.EquityType
	var expiration core.DateType
	var gamma sql.NullFloat64
	var highPrice sql.NullFloat64
	var impliedVolatility sql.NullFloat64
	var lastTradePrice sql.NullFloat64
	var lastTradeVolume sql.NullInt64
	var lastTradeTimestamp sql.NullInt64
	var lowPrice sql.NullFloat64
	var openInterest sql.NullInt64
	var rho sql.NullFloat64
	var symbol core.StockSymbolType
	var strikePrice float64
	var theta sql.NullFloat64
	var timestamp core.TimestampType
	var vega sql.NullFloat64

	err := rows.Scan(
		&id,
		&ask,
		&askTime,
		&askSize,
		&bid,
		&bidTime,
		&bidSize,
		&change,
		&delta,
		&equityType,
		&expiration,
		&gamma,
		&highPrice,
		&impliedVolatility,
		&lastTradePrice,
		&lastTradeVolume,
		&lastTradeTimestamp,
		&lowPrice,
		&openInterest,
		&rho,
		&strikePrice,
		&symbol,
		&theta,
		&timestamp,
		&vega,
	)

	if err != nil {
		return "",nil, fmt.Errorf("ScanSibylQuoteRow: had an error while scanning: %v", err)
	}

	quote := core.SibylOptionQuoteRecord{
		Ask:                ask,
		AskTime:            askTime,
		AskSize:            askSize,
		Bid:                bid,
		BidTime:            bidTime,
		BidSize:            bidSize,
		Change:             change,
		Delta:              delta,
		EquityType:         equityType,
		Expiration:         expiration,
		Gamma:              gamma,
		HighPrice:          highPrice,
		ImpliedVolatility:  impliedVolatility,
		LastTradePrice:     lastTradePrice,
		LastTradeVolume:    lastTradeVolume,
		LastTradeTimestamp: lastTradeTimestamp,
		LowPrice:           lowPrice,
		OpenInterest:       openInterest,
		Rho:                rho,
		Symbol:             symbol,
		StrikePrice:        strikePrice,
		Theta:              theta,
		Timestamp:          timestamp,
		Vega:               vega,
	}
	return id, &quote, nil
}
