// Code generated by ent, DO NOT EDIT.

package entity

const (
	// Label holds the string label denoting the entity type in the database.
	Label = "entity"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldActive holds the string denoting the active field in the database.
	FieldActive = "active"
	// FieldTicker holds the string denoting the ticker field in the database.
	FieldTicker = "ticker"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldDescription holds the string denoting the description field in the database.
	FieldDescription = "description"
	// FieldListDate holds the string denoting the list_date field in the database.
	FieldListDate = "list_date"
	// FieldDelisted holds the string denoting the delisted field in the database.
	FieldDelisted = "delisted"
	// EdgeExchanges holds the string denoting the exchanges edge name in mutations.
	EdgeExchanges = "exchanges"
	// EdgeIntervals holds the string denoting the intervals edge name in mutations.
	EdgeIntervals = "intervals"
	// EdgeDividends holds the string denoting the dividends edge name in mutations.
	EdgeDividends = "dividends"
	// EdgeSplits holds the string denoting the splits edge name in mutations.
	EdgeSplits = "splits"
	// EdgeFinancials holds the string denoting the financials edge name in mutations.
	EdgeFinancials = "financials"
	// Table holds the table name of the entity in the database.
	Table = "entities"
	// ExchangesTable is the table that holds the exchanges relation/edge. The primary key declared below.
	ExchangesTable = "entity_exchanges"
	// ExchangesInverseTable is the table name for the Exchange entity.
	// It exists in this package in order to avoid circular dependency with the "exchange" package.
	ExchangesInverseTable = "exchanges"
	// IntervalsTable is the table that holds the intervals relation/edge.
	IntervalsTable = "intervals"
	// IntervalsInverseTable is the table name for the Interval entity.
	// It exists in this package in order to avoid circular dependency with the "interval" package.
	IntervalsInverseTable = "intervals"
	// IntervalsColumn is the table column denoting the intervals relation/edge.
	IntervalsColumn = "stock_id"
	// DividendsTable is the table that holds the dividends relation/edge. The primary key declared below.
	DividendsTable = "entity_dividends"
	// DividendsInverseTable is the table name for the Dividend entity.
	// It exists in this package in order to avoid circular dependency with the "dividend" package.
	DividendsInverseTable = "dividends"
	// SplitsTable is the table that holds the splits relation/edge.
	SplitsTable = "splits"
	// SplitsInverseTable is the table name for the Split entity.
	// It exists in this package in order to avoid circular dependency with the "split" package.
	SplitsInverseTable = "splits"
	// SplitsColumn is the table column denoting the splits relation/edge.
	SplitsColumn = "entity_splits"
	// FinancialsTable is the table that holds the financials relation/edge. The primary key declared below.
	FinancialsTable = "entity_financials"
	// FinancialsInverseTable is the table name for the Financial entity.
	// It exists in this package in order to avoid circular dependency with the "financial" package.
	FinancialsInverseTable = "financials"
)

// Columns holds all SQL columns for entity fields.
var Columns = []string{
	FieldID,
	FieldActive,
	FieldTicker,
	FieldName,
	FieldDescription,
	FieldListDate,
	FieldDelisted,
}

var (
	// ExchangesPrimaryKey and ExchangesColumn2 are the table columns denoting the
	// primary key for the exchanges relation (M2M).
	ExchangesPrimaryKey = []string{"entity_id", "exchange_id"}
	// DividendsPrimaryKey and DividendsColumn2 are the table columns denoting the
	// primary key for the dividends relation (M2M).
	DividendsPrimaryKey = []string{"entity_id", "dividend_id"}
	// FinancialsPrimaryKey and FinancialsColumn2 are the table columns denoting the
	// primary key for the financials relation (M2M).
	FinancialsPrimaryKey = []string{"entity_id", "financial_id"}
)

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// TickerValidator is a validator for the "ticker" field. It is called by the builders before save.
	TickerValidator func(string) error
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// DescriptionValidator is a validator for the "description" field. It is called by the builders before save.
	DescriptionValidator func(string) error
)