// Code generated by ent, DO NOT EDIT.

package datasource

const (
	// Label holds the string label denoting the datasource type in the database.
	Label = "data_source"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldAddress holds the string denoting the address field in the database.
	FieldAddress = "address"
	// EdgeIntervals holds the string denoting the intervals edge name in mutations.
	EdgeIntervals = "intervals"
	// Table holds the table name of the datasource in the database.
	Table = "data_sources"
	// IntervalsTable is the table that holds the intervals relation/edge.
	IntervalsTable = "intervals"
	// IntervalsInverseTable is the table name for the Interval entity.
	// It exists in this package in order to avoid circular dependency with the "interval" package.
	IntervalsInverseTable = "intervals"
	// IntervalsColumn is the table column denoting the intervals relation/edge.
	IntervalsColumn = "data_source_id"
)

// Columns holds all SQL columns for datasource fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldAddress,
}

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
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// DefaultAddress holds the default value on creation for the "address" field.
	DefaultAddress string
)
