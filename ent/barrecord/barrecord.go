// Code generated by ent, DO NOT EDIT.

package barrecord

const (
	// Label holds the string label denoting the barrecord type in the database.
	Label = "bar_record"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldClose holds the string denoting the close field in the database.
	FieldClose = "close"
	// FieldHigh holds the string denoting the high field in the database.
	FieldHigh = "high"
	// FieldLow holds the string denoting the low field in the database.
	FieldLow = "low"
	// FieldOpen holds the string denoting the open field in the database.
	FieldOpen = "open"
	// FieldTimestamp holds the string denoting the timestamp field in the database.
	FieldTimestamp = "timestamp"
	// FieldVolume holds the string denoting the volume field in the database.
	FieldVolume = "volume"
	// FieldTransactions holds the string denoting the transactions field in the database.
	FieldTransactions = "transactions"
	// EdgeGroup holds the string denoting the group edge name in mutations.
	EdgeGroup = "group"
	// Table holds the table name of the barrecord in the database.
	Table = "bar_records"
	// GroupTable is the table that holds the group relation/edge.
	GroupTable = "bar_records"
	// GroupInverseTable is the table name for the BarGroup entity.
	// It exists in this package in order to avoid circular dependency with the "bargroup" package.
	GroupInverseTable = "bar_groups"
	// GroupColumn is the table column denoting the group relation/edge.
	GroupColumn = "bar_group_records"
)

// Columns holds all SQL columns for barrecord fields.
var Columns = []string{
	FieldID,
	FieldClose,
	FieldHigh,
	FieldLow,
	FieldOpen,
	FieldTimestamp,
	FieldVolume,
	FieldTransactions,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "bar_records"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"bar_group_records",
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
