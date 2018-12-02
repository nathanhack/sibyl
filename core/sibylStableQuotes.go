package core

//type SibylStableQuote struct {
//	annualDividend         sql.NullFloat64  // stock  // ally : iad   TD   divAmount
//	bookValue              sql.NullFloat64  // stock  //  "prbook"
//	closePrice             sql.NullFloat64  // both   // previous day close
//	contractSize           sql.NullInt64    // option // contract size for option
//	div                    sql.NullFloat64  // stock  // ally only //(= AnnualDiv/ divFreq) Latest announced cash dividend
//	divFreq                NullDivFrequency // stock  // Ally
//	divExTimestamp         sql.NullInt64    // stock  // millis  // date of last dividend
//	divPayTimestamp        sql.NullInt64    // stock  // millis // Ally
//	eps                    sql.NullFloat64  // stock  // "eps"
//	equityType             NullEquityType   // both   //  "put_call"    //CALL or PUT
//	exchange               sql.NullString   // both   // "exch"
//	exchangeDescription    sql.NullString   // both   // "exch_desc"
//	expirationTimestamp    sql.NullInt64    // option // millis
//	highPrice52Wk          sql.NullFloat64  // both   //
//	highPrice52WkTimestamp sql.NullInt64    // both   // ally only //
//	lowPrice52Wk           sql.NullFloat64  // both   //
//	lowPrice52WkTimestamp  sql.NullInt64    // both   // ally only //millis -- ally only
//	multiplier             sql.NullInt64    // option // "prem_mult"
//	name                   sql.NullString   // both   // "name"
//	openPrice              sql.NullFloat64  // both   //  "opn"
//	priceEarnings          sql.NullFloat64  // stock  //  "pe"
//	sharesOutstanding      sql.NullInt64    // stock  //
//	stockSymbol            sql.NullString   // both   //  either the stock symbol or root symbol ex. CAT
//	strikePrice            sql.NullFloat64  // option //
//	timestamp              sql.NullInt64    // both   // millis // Since these values should be stable the date is just the month/day/year
//	volatility             sql.NullFloat64  // stock  // one year volatility measure
//	yield                  sql.NullFloat64  // stock  //
//}
//
//
//
//
//func NewSibylStableQuote(
//	annualDividend sql.NullFloat64,
//	bookValue sql.NullFloat64,
//	closePrice sql.NullFloat64,
//	contractSize sql.NullInt64,
//	div sql.NullFloat64,
//	divFreq NullDivFrequency,
//	divExTimestamp sql.NullInt64,
//	divPayTimestamp sql.NullInt64,
//	eps sql.NullFloat64,
//	equityType NullEquityType,
//	exchange sql.NullString,
//	exchangeDescription sql.NullString,
//	expirationTimestamp sql.NullInt64,
//	highPrice52Wk sql.NullFloat64,
//	highPrice52WkTimestamp sql.NullInt64,
//	lowPrice52Wk sql.NullFloat64,
//	lowPrice52WkTimestamp sql.NullInt64,
//	multiplier sql.NullInt64,
//	name sql.NullString,
//	openPrice sql.NullFloat64,
//	priceEarnings sql.NullFloat64,
//	sharesOutstanding sql.NullInt64,
//	stockSymbol sql.NullString,
//	strikePrice sql.NullFloat64,
//	timestamp sql.NullInt64,
//	volatility sql.NullFloat64,
//	yield sql.NullFloat64,
//) *SibylStableQuote {
//	return &SibylStableQuote{
//		annualDividend,
//		bookValue,
//		closePrice,
//		contractSize,
//		div,
//		divFreq,
//		divExTimestamp,
//		divPayTimestamp,
//		eps,
//		equityType,
//		exchange,
//		exchangeDescription,
//		expirationTimestamp,
//		highPrice52Wk,
//		highPrice52WkTimestamp,
//		lowPrice52Wk,
//		lowPrice52WkTimestamp,
//		multiplier,
//		name,
//		openPrice,
//		priceEarnings,
//		sharesOutstanding,
//		stockSymbol,
//		strikePrice,
//		timestamp,
//		volatility,
//		yield,
//	}
//}
//func (ssq *SibylStableQuote) AnnualDividend() sql.NullFloat64 {
//	return ssq.annualDividend
//}
//
//func (ssq *SibylStableQuote) BookValue() sql.NullFloat64 {
//	return ssq.bookValue
//}
//
//func (ssq *SibylStableQuote) ClosePrice() sql.NullFloat64 {
//	return ssq.closePrice
//}
//
//func (ssq *SibylStableQuote) ContractSize() sql.NullInt64 {
//	return ssq.contractSize
//}
//
//func (ssq *SibylStableQuote) Div() sql.NullFloat64 {
//	return ssq.div
//}
//
//func (ssq *SibylStableQuote) DivFreq() NullDivFrequency {
//	return ssq.divFreq
//}
//
//func (ssq *SibylStableQuote) DivExTimestamp() sql.NullInt64 {
//	return ssq.divExTimestamp
//}
//
//func (ssq *SibylStableQuote) DivPayTimestamp() sql.NullInt64 {
//	return ssq.divPayTimestamp
//}
//
//func (ssq *SibylStableQuote) Eps() sql.NullFloat64 {
//	return ssq.eps
//}
//
//func (ssq *SibylStableQuote) EquityType() NullEquityType {
//	return ssq.equityType
//}
//
//func (ssq *SibylStableQuote) Exchange() sql.NullString {
//	return ssq.exchange
//}
//
//func (ssq *SibylStableQuote) ExchangeDescription() sql.NullString {
//	return ssq.exchangeDescription
//}
//
//func (ssq *SibylStableQuote) ExpirationTimestamp() sql.NullInt64 {
//	return ssq.expirationTimestamp
//}
//
//func (ssq *SibylStableQuote) HighPrice52Wk() sql.NullFloat64 {
//	return ssq.highPrice52Wk
//}
//
//func (ssq *SibylStableQuote) HighPrice52WkTimestamp() sql.NullInt64 {
//	return ssq.highPrice52WkTimestamp
//}
//
//func (ssq *SibylStableQuote) LowPrice52Wk() sql.NullFloat64 {
//	return ssq.lowPrice52Wk
//}
//
//func (ssq *SibylStableQuote) LowPrice52WkTimestamp() sql.NullInt64 {
//	return ssq.lowPrice52WkTimestamp
//}
//
//func (ssq *SibylStableQuote) Multiplier() sql.NullInt64 {
//	return ssq.multiplier
//}
//
//func (ssq *SibylStableQuote) Name() sql.NullString {
//	return ssq.name
//}
//
//func (ssq *SibylStableQuote) OpenPrice() sql.NullFloat64 {
//	return ssq.openPrice
//}
//
//func (ssq *SibylStableQuote) PriceEarnings() sql.NullFloat64 {
//	return ssq.priceEarnings
//}
//
//func (ssq *SibylStableQuote) SharesOutstanding() sql.NullInt64 {
//	return ssq.sharesOutstanding
//}
//
//func (ssq *SibylStableQuote) Symbol() sql.NullString {
//	return ssq.stockSymbol
//}
//
//func (ssq *SibylStableQuote) StrikePrice() sql.NullFloat64 {
//	return ssq.strikePrice
//}
//
//func (ssq *SibylStableQuote) Timestamp() sql.NullInt64 {
//	return ssq.timestamp
//}
//
//func (ssq *SibylStableQuote) Volatility() sql.NullFloat64 {
//	return ssq.volatility
//}
//
//func (ssq *SibylStableQuote) Yield() sql.NullFloat64 {
//	return ssq.yield
//}
//
//func (ssq *SibylStableQuote) String() string {
//	return fmt.Sprintf("{%v}", ssq.StringBlindWithDelimiter(",", "", true))
//}
//
//func (ssq *SibylStableQuote) StringBlindWithDelimiter(delimiter string, nullString string, stringEscapes bool) string {
//	builder := strings.Builder{}
//	builder.WriteString(nullFloat64ToString(ssq.annualDividend, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullFloat64ToString(ssq.bookValue, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullFloat64ToString(ssq.closePrice, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullInt64ToString(ssq.contractSize, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullFloat64ToString(ssq.div, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullDivFrequencyToString(ssq.divFreq, nullString, stringEscapes))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullInt64ToString(ssq.divExTimestamp, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullInt64ToString(ssq.divPayTimestamp, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullFloat64ToString(ssq.eps, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullEquityTypeToString(ssq.equityType, nullString, stringEscapes))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullStringToString(ssq.exchange, nullString, stringEscapes))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullStringToString(ssq.exchangeDescription, nullString, stringEscapes))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullInt64ToString(ssq.expirationTimestamp, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullFloat64ToString(ssq.highPrice52Wk, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullInt64ToString(ssq.highPrice52WkTimestamp, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullFloat64ToString(ssq.lowPrice52Wk, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullInt64ToString(ssq.lowPrice52WkTimestamp, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullInt64ToString(ssq.multiplier, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullStringToString(ssq.name, nullString, stringEscapes))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullFloat64ToString(ssq.openPrice, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullFloat64ToString(ssq.priceEarnings, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullInt64ToString(ssq.sharesOutstanding, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullStringToString(ssq.stockSymbol, nullString, stringEscapes))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullFloat64ToString(ssq.strikePrice, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullInt64ToString(ssq.timestamp, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullFloat64ToString(ssq.volatility, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullFloat64ToString(ssq.yield, nullString))
//	return builder.String()
//}
