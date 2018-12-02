package core

////Symbol
//type Symbol interface {
//	ExpirationTimestamp() int64
//	IsOption() bool
//	IsStock() bool
//	Symbol() string
//	StrikePrice() float64
//	String() string
//	Type() EquityType
//}
//
//type Symbol struct {
//	Symbol string
//}
//
//func (ss *Symbol) ExpirationTimestamp() int64 {
//	return 0
//}
//
//func (ss *Symbol) IsOption() bool {
//	return false
//}
//
//func (ss *Symbol) IsStock() bool {
//	return true
//}
//
//func (ss *Symbol) Symbol() string {
//	return ss.Symbol
//}
//
//func (ss *Symbol) StrikePrice() float64 {
//	return 0
//}
//
//func (ss *Symbol) String() string {
//	return fmt.Sprintf("{%v}", ss.Symbol)
//}
//
//func (ss *Symbol) Type() EquityType {
//	return ""
//}
