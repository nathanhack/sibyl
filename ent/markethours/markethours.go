// Code generated by ent, DO NOT EDIT.

package markethours

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the markethours type in the database.
	Label = "market_hours"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldDate holds the string denoting the date field in the database.
	FieldDate = "date"
	// FieldStartTime holds the string denoting the start_time field in the database.
	FieldStartTime = "start_time"
	// FieldEndTime holds the string denoting the end_time field in the database.
	FieldEndTime = "end_time"
	// EdgeMarketInfo holds the string denoting the market_info edge name in mutations.
	EdgeMarketInfo = "market_info"
	// Table holds the table name of the markethours in the database.
	Table = "market_hours"
	// MarketInfoTable is the table that holds the market_info relation/edge.
	MarketInfoTable = "market_hours"
	// MarketInfoInverseTable is the table name for the MarketInfo entity.
	// It exists in this package in order to avoid circular dependency with the "marketinfo" package.
	MarketInfoInverseTable = "market_infos"
	// MarketInfoColumn is the table column denoting the market_info relation/edge.
	MarketInfoColumn = "market_info_hours"
)

// Columns holds all SQL columns for markethours fields.
var Columns = []string{
	FieldID,
	FieldDate,
	FieldStartTime,
	FieldEndTime,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "market_hours"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"market_info_hours",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

// OrderOption defines the ordering options for the MarketHours queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByDate orders the results by the date field.
func ByDate(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDate, opts...).ToFunc()
}

// ByStartTime orders the results by the start_time field.
func ByStartTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldStartTime, opts...).ToFunc()
}

// ByEndTime orders the results by the end_time field.
func ByEndTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEndTime, opts...).ToFunc()
}

// ByMarketInfoField orders the results by market_info field.
func ByMarketInfoField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newMarketInfoStep(), sql.OrderByField(field, opts...))
	}
}
func newMarketInfoStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(MarketInfoInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, MarketInfoTable, MarketInfoColumn),
	)
}
