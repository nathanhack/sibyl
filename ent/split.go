// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/nathanhack/sibyl/ent/entity"
	"github.com/nathanhack/sibyl/ent/split"
)

// Split is the model entity for the Split schema.
type Split struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// ExecutionDate holds the value of the "execution_date" field.
	ExecutionDate time.Time `json:"execution_date,omitempty"`
	// From holds the value of the "from" field.
	From float64 `json:"from,omitempty"`
	// To holds the value of the "to" field.
	To float64 `json:"to,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the SplitQuery when eager-loading is set.
	Edges         SplitEdges `json:"edges"`
	entity_splits *int
	selectValues  sql.SelectValues
}

// SplitEdges holds the relations/edges for other nodes in the graph.
type SplitEdges struct {
	// Stock holds the value of the stock edge.
	Stock *Entity `json:"stock,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// StockOrErr returns the Stock value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e SplitEdges) StockOrErr() (*Entity, error) {
	if e.Stock != nil {
		return e.Stock, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: entity.Label}
	}
	return nil, &NotLoadedError{edge: "stock"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Split) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case split.FieldFrom, split.FieldTo:
			values[i] = new(sql.NullFloat64)
		case split.FieldID:
			values[i] = new(sql.NullInt64)
		case split.FieldExecutionDate:
			values[i] = new(sql.NullTime)
		case split.ForeignKeys[0]: // entity_splits
			values[i] = new(sql.NullInt64)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Split fields.
func (s *Split) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case split.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			s.ID = int(value.Int64)
		case split.FieldExecutionDate:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field execution_date", values[i])
			} else if value.Valid {
				s.ExecutionDate = value.Time
			}
		case split.FieldFrom:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field from", values[i])
			} else if value.Valid {
				s.From = value.Float64
			}
		case split.FieldTo:
			if value, ok := values[i].(*sql.NullFloat64); !ok {
				return fmt.Errorf("unexpected type %T for field to", values[i])
			} else if value.Valid {
				s.To = value.Float64
			}
		case split.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field entity_splits", value)
			} else if value.Valid {
				s.entity_splits = new(int)
				*s.entity_splits = int(value.Int64)
			}
		default:
			s.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Split.
// This includes values selected through modifiers, order, etc.
func (s *Split) Value(name string) (ent.Value, error) {
	return s.selectValues.Get(name)
}

// QueryStock queries the "stock" edge of the Split entity.
func (s *Split) QueryStock() *EntityQuery {
	return NewSplitClient(s.config).QueryStock(s)
}

// Update returns a builder for updating this Split.
// Note that you need to call Split.Unwrap() before calling this method if this Split
// was returned from a transaction, and the transaction was committed or rolled back.
func (s *Split) Update() *SplitUpdateOne {
	return NewSplitClient(s.config).UpdateOne(s)
}

// Unwrap unwraps the Split entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (s *Split) Unwrap() *Split {
	_tx, ok := s.config.driver.(*txDriver)
	if !ok {
		panic("ent: Split is not a transactional entity")
	}
	s.config.driver = _tx.drv
	return s
}

// String implements the fmt.Stringer.
func (s *Split) String() string {
	var builder strings.Builder
	builder.WriteString("Split(")
	builder.WriteString(fmt.Sprintf("id=%v, ", s.ID))
	builder.WriteString("execution_date=")
	builder.WriteString(s.ExecutionDate.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("from=")
	builder.WriteString(fmt.Sprintf("%v", s.From))
	builder.WriteString(", ")
	builder.WriteString("to=")
	builder.WriteString(fmt.Sprintf("%v", s.To))
	builder.WriteByte(')')
	return builder.String()
}

// Splits is a parsable slice of Split.
type Splits []*Split
