// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/nathanhack/sibyl/ent/markethours"
	"github.com/nathanhack/sibyl/ent/marketinfo"
)

// MarketHours is the model entity for the MarketHours schema.
type MarketHours struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Date holds the value of the "date" field.
	Date time.Time `json:"date,omitempty"`
	// StartTime holds the value of the "start_time" field.
	StartTime time.Time `json:"start_time,omitempty"`
	// EndTime holds the value of the "end_time" field.
	EndTime time.Time `json:"end_time,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the MarketHoursQuery when eager-loading is set.
	Edges             MarketHoursEdges `json:"edges"`
	market_info_hours *int
}

// MarketHoursEdges holds the relations/edges for other nodes in the graph.
type MarketHoursEdges struct {
	// MarketInfo holds the value of the market_info edge.
	MarketInfo *MarketInfo `json:"market_info,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// MarketInfoOrErr returns the MarketInfo value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e MarketHoursEdges) MarketInfoOrErr() (*MarketInfo, error) {
	if e.loadedTypes[0] {
		if e.MarketInfo == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: marketinfo.Label}
		}
		return e.MarketInfo, nil
	}
	return nil, &NotLoadedError{edge: "market_info"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*MarketHours) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case markethours.FieldID:
			values[i] = new(sql.NullInt64)
		case markethours.FieldDate, markethours.FieldStartTime, markethours.FieldEndTime:
			values[i] = new(sql.NullTime)
		case markethours.ForeignKeys[0]: // market_info_hours
			values[i] = new(sql.NullInt64)
		default:
			return nil, fmt.Errorf("unexpected column %q for type MarketHours", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the MarketHours fields.
func (mh *MarketHours) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case markethours.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			mh.ID = int(value.Int64)
		case markethours.FieldDate:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field date", values[i])
			} else if value.Valid {
				mh.Date = value.Time
			}
		case markethours.FieldStartTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field start_time", values[i])
			} else if value.Valid {
				mh.StartTime = value.Time
			}
		case markethours.FieldEndTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field end_time", values[i])
			} else if value.Valid {
				mh.EndTime = value.Time
			}
		case markethours.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field market_info_hours", value)
			} else if value.Valid {
				mh.market_info_hours = new(int)
				*mh.market_info_hours = int(value.Int64)
			}
		}
	}
	return nil
}

// QueryMarketInfo queries the "market_info" edge of the MarketHours entity.
func (mh *MarketHours) QueryMarketInfo() *MarketInfoQuery {
	return (&MarketHoursClient{config: mh.config}).QueryMarketInfo(mh)
}

// Update returns a builder for updating this MarketHours.
// Note that you need to call MarketHours.Unwrap() before calling this method if this MarketHours
// was returned from a transaction, and the transaction was committed or rolled back.
func (mh *MarketHours) Update() *MarketHoursUpdateOne {
	return (&MarketHoursClient{config: mh.config}).UpdateOne(mh)
}

// Unwrap unwraps the MarketHours entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (mh *MarketHours) Unwrap() *MarketHours {
	_tx, ok := mh.config.driver.(*txDriver)
	if !ok {
		panic("ent: MarketHours is not a transactional entity")
	}
	mh.config.driver = _tx.drv
	return mh
}

// String implements the fmt.Stringer.
func (mh *MarketHours) String() string {
	var builder strings.Builder
	builder.WriteString("MarketHours(")
	builder.WriteString(fmt.Sprintf("id=%v, ", mh.ID))
	builder.WriteString("date=")
	builder.WriteString(mh.Date.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("start_time=")
	builder.WriteString(mh.StartTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("end_time=")
	builder.WriteString(mh.EndTime.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// MarketHoursSlice is a parsable slice of MarketHours.
type MarketHoursSlice []*MarketHours

func (mh MarketHoursSlice) config(cfg config) {
	for _i := range mh {
		mh[_i].config = cfg
	}
}