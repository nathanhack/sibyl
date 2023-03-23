// Code generated by ent, DO NOT EDIT.

package bargroup

const (
	// Label holds the string label denoting the bargroup type in the database.
	Label = "bar_group"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldFirst holds the string denoting the first field in the database.
	FieldFirst = "first"
	// FieldLast holds the string denoting the last field in the database.
	FieldLast = "last"
	// FieldCount holds the string denoting the count field in the database.
	FieldCount = "count"
	// FieldTimeRangeID holds the string denoting the time_range_id field in the database.
	FieldTimeRangeID = "time_range_id"
	// EdgeTimeRange holds the string denoting the time_range edge name in mutations.
	EdgeTimeRange = "time_range"
	// EdgeRecords holds the string denoting the records edge name in mutations.
	EdgeRecords = "records"
	// Table holds the table name of the bargroup in the database.
	Table = "bar_groups"
	// TimeRangeTable is the table that holds the time_range relation/edge.
	TimeRangeTable = "bar_groups"
	// TimeRangeInverseTable is the table name for the BarTimeRange entity.
	// It exists in this package in order to avoid circular dependency with the "bartimerange" package.
	TimeRangeInverseTable = "bar_time_ranges"
	// TimeRangeColumn is the table column denoting the time_range relation/edge.
	TimeRangeColumn = "time_range_id"
	// RecordsTable is the table that holds the records relation/edge.
	RecordsTable = "bar_records"
	// RecordsInverseTable is the table name for the BarRecord entity.
	// It exists in this package in order to avoid circular dependency with the "barrecord" package.
	RecordsInverseTable = "bar_records"
	// RecordsColumn is the table column denoting the records relation/edge.
	RecordsColumn = "bar_group_records"
)

// Columns holds all SQL columns for bargroup fields.
var Columns = []string{
	FieldID,
	FieldFirst,
	FieldLast,
	FieldCount,
	FieldTimeRangeID,
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
