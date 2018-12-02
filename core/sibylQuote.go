package core

//type SibylQuote struct {
//	ask                 sql.NullFloat64 // both   //
//	askTime             sql.NullInt64   // both   //
//	askSize             sql.NullInt64   // both   //
//	beta                sql.NullFloat64 // stock  //
//	bid                 sql.NullFloat64 // both   //
//	bidTime             sql.NullInt64   // both   // millis
//	bidSize             sql.NullInt64   // both   //
//	change              sql.NullFloat64 // both   //  (+/-) number
//	delta               sql.NullFloat64 // option //  "idelta"
//	equityType          NullEquityType  // both   //  "put_call"     //CALL or PUT
//	expirationTimestamp sql.NullInt64   // option //  millis
//	gamma               sql.NullFloat64 // option //  "igamma"
//	highPrice           sql.NullFloat64 // both   //  "hi"`
//	impliedVolatility   sql.NullFloat64 // option //
//	lastTradePrice      sql.NullFloat64 // both   //  ally LastTrade ... TD mark or regularMarketLastPrice
//	lastTradeVolume     sql.NullInt64   // both   //  ally incr_vl   TD : lastSize
//	lastTradeTimestamp  sql.NullInt64   // both   //  ally DateTime   TD :tradeTimeInLong
//	lowPrice            sql.NullFloat64 // both   //  "lo"
//	openInterest        sql.NullInt64   // option //  "openinterest"
//	rho                 sql.NullFloat64 // option //  "irho"
//	stockSymbol         sql.NullString  // both   //  either the stock symbol or root symbol ex. CAT
//	strikePrice         sql.NullFloat64 // option //
//	theta               sql.NullFloat64 // option //  "itheta"
//	timestamp           sql.NullInt64   // both   //  millis // this is the GMT timestamp (ms from epoch) of the this quote
//	vega                sql.NullFloat64 // option //  "ivega"
//	volume              sql.NullInt64   // stock  //
//	volWeightedAvgPrice sql.NullFloat64 // stock  //
//}
//
//func NewSibylQuote(
//	ask sql.NullFloat64,
//	askTime sql.NullInt64,
//	askSize sql.NullInt64,
//	beta sql.NullFloat64,
//	bid sql.NullFloat64,
//	bidTime sql.NullInt64,
//	bidSize sql.NullInt64,
//	change sql.NullFloat64,
//	delta sql.NullFloat64,
//	equityType NullEquityType,
//	expirationTimestamp sql.NullInt64,
//	gamma sql.NullFloat64,
//	highPrice sql.NullFloat64,
//	impliedVolatility sql.NullFloat64,
//	lastTradePrice sql.NullFloat64,
//	lastTradeVolume sql.NullInt64,
//	lastTradeTimestamp sql.NullInt64,
//	lowPrice sql.NullFloat64,
//	openInterest sql.NullInt64,
//	rho sql.NullFloat64,
//	stockSymbol sql.NullString,
//	strikePrice sql.NullFloat64,
//	theta sql.NullFloat64,
//	timestamp sql.NullInt64,
//	vega sql.NullFloat64,
//	volume sql.NullInt64,
//	volWeightedAvgPrice sql.NullFloat64,
//) *SibylQuote {
//	return &SibylQuote{
//		ask,
//		askTime,
//		askSize,
//		beta,
//		bid,
//		bidTime,
//		bidSize,
//		change,
//		delta,
//		equityType,
//		expirationTimestamp,
//		gamma,
//		highPrice,
//		impliedVolatility,
//		lastTradePrice,
//		lastTradeVolume,
//		lastTradeTimestamp,
//		lowPrice,
//		openInterest,
//		rho,
//		stockSymbol,
//		strikePrice,
//		theta,
//		timestamp,
//		vega,
//		volume,
//		volWeightedAvgPrice,
//	}
//}
//
//func (sbq *SibylQuote) Ask() sql.NullFloat64 {
//	return sbq.ask
//}
//
//func (sbq *SibylQuote) AskTime() sql.NullInt64 {
//	return sbq.askTime
//}
//
//func (sbq *SibylQuote) AskSize() sql.NullInt64 {
//	return sbq.askSize
//}
//
//func (sbq *SibylQuote) Beta() sql.NullFloat64 {
//	return sbq.beta
//}
//
//func (sbq *SibylQuote) Bid() sql.NullFloat64 {
//	return sbq.bid
//}
//
//func (sbq *SibylQuote) BidTime() sql.NullInt64 {
//	return sbq.bidTime
//}
//
//func (sbq *SibylQuote) BidSize() sql.NullInt64 {
//	return sbq.bidSize
//}
//
//func (sbq *SibylQuote) Change() sql.NullFloat64 {
//	return sbq.change
//}
//
//func (sbq *SibylQuote) Delta() sql.NullFloat64 {
//	return sbq.delta
//}
//
//func (sbq *SibylQuote) EquityType() NullEquityType {
//	return sbq.equityType
//}
//
//func (sbq *SibylQuote) ExpirationTimestamp() sql.NullInt64 {
//	return sbq.expirationTimestamp
//}
//
//func (sbq *SibylQuote) Gamma() sql.NullFloat64 {
//	return sbq.gamma
//}
//
//func (sbq *SibylQuote) HighPrice() sql.NullFloat64 {
//	return sbq.highPrice
//}
//
//func (sbq *SibylQuote) ImpliedVolatility() sql.NullFloat64 {
//	return sbq.impliedVolatility
//}
//
//func (sbq *SibylQuote) LastTradePrice() sql.NullFloat64 {
//	return sbq.lastTradePrice
//}
//
//func (sbq *SibylQuote) LastTradeVolume() sql.NullInt64 {
//	return sbq.lastTradeVolume
//}
//
//func (sbq *SibylQuote) LastTradeTimestamp() sql.NullInt64 {
//	return sbq.lastTradeTimestamp
//}
//
//func (sbq *SibylQuote) LowPrice() sql.NullFloat64 {
//	return sbq.lowPrice
//}
//
//func (sbq *SibylQuote) OpenInterest() sql.NullInt64 {
//	return sbq.openInterest
//}
//
//func (sbq *SibylQuote) Rho() sql.NullFloat64 {
//	return sbq.rho
//}
//
//func (sbq *SibylQuote) Symbol() sql.NullString {
//	return sbq.stockSymbol
//}
//
//func (sbq *SibylQuote) StrikePrice() sql.NullFloat64 {
//	return sbq.strikePrice
//}
//
//func (sbq *SibylQuote) Theta() sql.NullFloat64 {
//	return sbq.theta
//}
//
//func (sbq *SibylQuote) Timestamp() sql.NullInt64 {
//	return sbq.timestamp
//}
//
//func (sbq *SibylQuote) Vega() sql.NullFloat64 {
//	return sbq.vega
//}
//
//func (sbq *SibylQuote) Volume() sql.NullInt64 {
//	return sbq.volume
//}
//
//func (sbq *SibylQuote) VolWeightedAvgPrice() sql.NullFloat64 {
//	return sbq.volWeightedAvgPrice
//}
//
//func (sbq *SibylQuote) String() string {
//	return fmt.Sprintf("{%v}", sbq.StringBlindWithDelimiter(",", "", true))
//}
//
//func (sbq *SibylQuote) StringBlindWithDelimiter(delimiter string, nullString string, stringEscapes bool) string {
//	//it's faster fmt. or strings.Join
//	builder := strings.Builder{}
//	builder.WriteString(nullFloat64ToString(sbq.ask, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullInt64ToString(sbq.askTime, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullInt64ToString(sbq.askSize, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullFloat64ToString(sbq.beta, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullFloat64ToString(sbq.bid, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullInt64ToString(sbq.bidTime, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullInt64ToString(sbq.bidSize, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullFloat64ToString(sbq.change, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullFloat64ToString(sbq.delta, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullEquityTypeToString(sbq.equityType, nullString, stringEscapes))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullInt64ToString(sbq.expirationTimestamp, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullFloat64ToString(sbq.gamma, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullFloat64ToString(sbq.highPrice, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullFloat64ToString(sbq.impliedVolatility, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullFloat64ToString(sbq.lastTradePrice, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullInt64ToString(sbq.lastTradeVolume, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullInt64ToString(sbq.lastTradeTimestamp, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullFloat64ToString(sbq.lowPrice, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullInt64ToString(sbq.openInterest, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullFloat64ToString(sbq.rho, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullStringToString(sbq.stockSymbol, nullString, stringEscapes))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullFloat64ToString(sbq.strikePrice, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullFloat64ToString(sbq.theta, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullInt64ToString(sbq.timestamp, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullFloat64ToString(sbq.vega, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullInt64ToString(sbq.volume, nullString))
//	builder.WriteString(delimiter)
//	builder.WriteString(nullFloat64ToString(sbq.volWeightedAvgPrice, nullString))
//	return builder.String()
//}
