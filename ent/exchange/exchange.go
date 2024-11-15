// Code generated by ent, DO NOT EDIT.

package exchange

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the exchange type in the database.
	Label = "exchange"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCode holds the string denoting the code field in the database.
	FieldCode = "code"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// EdgeStocks holds the string denoting the stocks edge name in mutations.
	EdgeStocks = "stocks"
	// Table holds the table name of the exchange in the database.
	Table = "exchanges"
	// StocksTable is the table that holds the stocks relation/edge. The primary key declared below.
	StocksTable = "entity_exchanges"
	// StocksInverseTable is the table name for the Entity entity.
	// It exists in this package in order to avoid circular dependency with the "entity" package.
	StocksInverseTable = "entities"
)

// Columns holds all SQL columns for exchange fields.
var Columns = []string{
	FieldID,
	FieldCode,
	FieldName,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "exchanges"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"trade_record_exchange",
}

var (
	// StocksPrimaryKey and StocksColumn2 are the table columns denoting the
	// primary key for the stocks relation (M2M).
	StocksPrimaryKey = []string{"entity_id", "exchange_id"}
)

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

// OrderOption defines the ordering options for the Exchange queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByCode orders the results by the code field.
func ByCode(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCode, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByStocksCount orders the results by stocks count.
func ByStocksCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newStocksStep(), opts...)
	}
}

// ByStocks orders the results by stocks terms.
func ByStocks(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newStocksStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newStocksStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(StocksInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, true, StocksTable, StocksPrimaryKey...),
	)
}
