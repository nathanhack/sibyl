// Code generated by ent, DO NOT EDIT.

package financial

const (
	// Label holds the string label denoting the financial type in the database.
	Label = "financial"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// EdgeStock holds the string denoting the stock edge name in mutations.
	EdgeStock = "stock"
	// Table holds the table name of the financial in the database.
	Table = "financials"
	// StockTable is the table that holds the stock relation/edge. The primary key declared below.
	StockTable = "entity_financials"
	// StockInverseTable is the table name for the Entity entity.
	// It exists in this package in order to avoid circular dependency with the "entity" package.
	StockInverseTable = "entities"
)

// Columns holds all SQL columns for financial fields.
var Columns = []string{
	FieldID,
}

var (
	// StockPrimaryKey and StockColumn2 are the table columns denoting the
	// primary key for the stock relation (M2M).
	StockPrimaryKey = []string{"entity_id", "financial_id"}
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
