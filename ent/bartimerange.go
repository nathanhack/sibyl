// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/nathanhack/sibyl/ent/bartimerange"
	"github.com/nathanhack/sibyl/ent/interval"
)

// BarTimeRange is the model entity for the BarTimeRange schema.
type BarTimeRange struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Start holds the value of the "start" field.
	Start time.Time `json:"start,omitempty"`
	// End holds the value of the "end" field.
	End time.Time `json:"end,omitempty"`
	// The number of BarGroups
	Count int `json:"count,omitempty"`
	// IntervalID holds the value of the "interval_id" field.
	IntervalID int `json:"interval_id,omitempty"`
	// Status holds the value of the "status" field.
	Status bartimerange.Status `json:"status,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the BarTimeRangeQuery when eager-loading is set.
	Edges        BarTimeRangeEdges `json:"edges"`
	selectValues sql.SelectValues
}

// BarTimeRangeEdges holds the relations/edges for other nodes in the graph.
type BarTimeRangeEdges struct {
	// Interval holds the value of the interval edge.
	Interval *Interval `json:"interval,omitempty"`
	// Groups holds the value of the groups edge.
	Groups []*BarGroup `json:"groups,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// IntervalOrErr returns the Interval value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e BarTimeRangeEdges) IntervalOrErr() (*Interval, error) {
	if e.Interval != nil {
		return e.Interval, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: interval.Label}
	}
	return nil, &NotLoadedError{edge: "interval"}
}

// GroupsOrErr returns the Groups value or an error if the edge
// was not loaded in eager-loading.
func (e BarTimeRangeEdges) GroupsOrErr() ([]*BarGroup, error) {
	if e.loadedTypes[1] {
		return e.Groups, nil
	}
	return nil, &NotLoadedError{edge: "groups"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*BarTimeRange) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case bartimerange.FieldID, bartimerange.FieldCount, bartimerange.FieldIntervalID:
			values[i] = new(sql.NullInt64)
		case bartimerange.FieldStatus:
			values[i] = new(sql.NullString)
		case bartimerange.FieldStart, bartimerange.FieldEnd, bartimerange.FieldUpdateTime:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the BarTimeRange fields.
func (btr *BarTimeRange) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case bartimerange.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			btr.ID = int(value.Int64)
		case bartimerange.FieldStart:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field start", values[i])
			} else if value.Valid {
				btr.Start = value.Time
			}
		case bartimerange.FieldEnd:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field end", values[i])
			} else if value.Valid {
				btr.End = value.Time
			}
		case bartimerange.FieldCount:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field count", values[i])
			} else if value.Valid {
				btr.Count = int(value.Int64)
			}
		case bartimerange.FieldIntervalID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field interval_id", values[i])
			} else if value.Valid {
				btr.IntervalID = int(value.Int64)
			}
		case bartimerange.FieldStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value.Valid {
				btr.Status = bartimerange.Status(value.String)
			}
		case bartimerange.FieldUpdateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field update_time", values[i])
			} else if value.Valid {
				btr.UpdateTime = value.Time
			}
		default:
			btr.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the BarTimeRange.
// This includes values selected through modifiers, order, etc.
func (btr *BarTimeRange) Value(name string) (ent.Value, error) {
	return btr.selectValues.Get(name)
}

// QueryInterval queries the "interval" edge of the BarTimeRange entity.
func (btr *BarTimeRange) QueryInterval() *IntervalQuery {
	return NewBarTimeRangeClient(btr.config).QueryInterval(btr)
}

// QueryGroups queries the "groups" edge of the BarTimeRange entity.
func (btr *BarTimeRange) QueryGroups() *BarGroupQuery {
	return NewBarTimeRangeClient(btr.config).QueryGroups(btr)
}

// Update returns a builder for updating this BarTimeRange.
// Note that you need to call BarTimeRange.Unwrap() before calling this method if this BarTimeRange
// was returned from a transaction, and the transaction was committed or rolled back.
func (btr *BarTimeRange) Update() *BarTimeRangeUpdateOne {
	return NewBarTimeRangeClient(btr.config).UpdateOne(btr)
}

// Unwrap unwraps the BarTimeRange entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (btr *BarTimeRange) Unwrap() *BarTimeRange {
	_tx, ok := btr.config.driver.(*txDriver)
	if !ok {
		panic("ent: BarTimeRange is not a transactional entity")
	}
	btr.config.driver = _tx.drv
	return btr
}

// String implements the fmt.Stringer.
func (btr *BarTimeRange) String() string {
	var builder strings.Builder
	builder.WriteString("BarTimeRange(")
	builder.WriteString(fmt.Sprintf("id=%v, ", btr.ID))
	builder.WriteString("start=")
	builder.WriteString(btr.Start.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("end=")
	builder.WriteString(btr.End.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("count=")
	builder.WriteString(fmt.Sprintf("%v", btr.Count))
	builder.WriteString(", ")
	builder.WriteString("interval_id=")
	builder.WriteString(fmt.Sprintf("%v", btr.IntervalID))
	builder.WriteString(", ")
	builder.WriteString("status=")
	builder.WriteString(fmt.Sprintf("%v", btr.Status))
	builder.WriteString(", ")
	builder.WriteString("update_time=")
	builder.WriteString(btr.UpdateTime.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// BarTimeRanges is a parsable slice of BarTimeRange.
type BarTimeRanges []*BarTimeRange
