// Code generated by ent, DO NOT EDIT.

package tradetimerange

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/nathanhack/sibyl/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(sql.FieldLTE(FieldID, id))
}

// Start applies equality check predicate on the "start" field. It's identical to StartEQ.
func Start(v time.Time) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(sql.FieldEQ(FieldStart, v))
}

// End applies equality check predicate on the "end" field. It's identical to EndEQ.
func End(v time.Time) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(sql.FieldEQ(FieldEnd, v))
}

// IntervalID applies equality check predicate on the "interval_id" field. It's identical to IntervalIDEQ.
func IntervalID(v int) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(sql.FieldEQ(FieldIntervalID, v))
}

// StartEQ applies the EQ predicate on the "start" field.
func StartEQ(v time.Time) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(sql.FieldEQ(FieldStart, v))
}

// StartNEQ applies the NEQ predicate on the "start" field.
func StartNEQ(v time.Time) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(sql.FieldNEQ(FieldStart, v))
}

// StartIn applies the In predicate on the "start" field.
func StartIn(vs ...time.Time) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(sql.FieldIn(FieldStart, vs...))
}

// StartNotIn applies the NotIn predicate on the "start" field.
func StartNotIn(vs ...time.Time) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(sql.FieldNotIn(FieldStart, vs...))
}

// StartGT applies the GT predicate on the "start" field.
func StartGT(v time.Time) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(sql.FieldGT(FieldStart, v))
}

// StartGTE applies the GTE predicate on the "start" field.
func StartGTE(v time.Time) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(sql.FieldGTE(FieldStart, v))
}

// StartLT applies the LT predicate on the "start" field.
func StartLT(v time.Time) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(sql.FieldLT(FieldStart, v))
}

// StartLTE applies the LTE predicate on the "start" field.
func StartLTE(v time.Time) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(sql.FieldLTE(FieldStart, v))
}

// EndEQ applies the EQ predicate on the "end" field.
func EndEQ(v time.Time) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(sql.FieldEQ(FieldEnd, v))
}

// EndNEQ applies the NEQ predicate on the "end" field.
func EndNEQ(v time.Time) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(sql.FieldNEQ(FieldEnd, v))
}

// EndIn applies the In predicate on the "end" field.
func EndIn(vs ...time.Time) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(sql.FieldIn(FieldEnd, vs...))
}

// EndNotIn applies the NotIn predicate on the "end" field.
func EndNotIn(vs ...time.Time) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(sql.FieldNotIn(FieldEnd, vs...))
}

// EndGT applies the GT predicate on the "end" field.
func EndGT(v time.Time) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(sql.FieldGT(FieldEnd, v))
}

// EndGTE applies the GTE predicate on the "end" field.
func EndGTE(v time.Time) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(sql.FieldGTE(FieldEnd, v))
}

// EndLT applies the LT predicate on the "end" field.
func EndLT(v time.Time) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(sql.FieldLT(FieldEnd, v))
}

// EndLTE applies the LTE predicate on the "end" field.
func EndLTE(v time.Time) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(sql.FieldLTE(FieldEnd, v))
}

// IntervalIDEQ applies the EQ predicate on the "interval_id" field.
func IntervalIDEQ(v int) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(sql.FieldEQ(FieldIntervalID, v))
}

// IntervalIDNEQ applies the NEQ predicate on the "interval_id" field.
func IntervalIDNEQ(v int) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(sql.FieldNEQ(FieldIntervalID, v))
}

// IntervalIDIn applies the In predicate on the "interval_id" field.
func IntervalIDIn(vs ...int) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(sql.FieldIn(FieldIntervalID, vs...))
}

// IntervalIDNotIn applies the NotIn predicate on the "interval_id" field.
func IntervalIDNotIn(vs ...int) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(sql.FieldNotIn(FieldIntervalID, vs...))
}

// HasInterval applies the HasEdge predicate on the "interval" edge.
func HasInterval() predicate.TradeTimeRange {
	return predicate.TradeTimeRange(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, IntervalTable, IntervalColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasIntervalWith applies the HasEdge predicate on the "interval" edge with a given conditions (other predicates).
func HasIntervalWith(preds ...predicate.Interval) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(IntervalInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, IntervalTable, IntervalColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasRecords applies the HasEdge predicate on the "records" edge.
func HasRecords() predicate.TradeTimeRange {
	return predicate.TradeTimeRange(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, RecordsTable, RecordsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasRecordsWith applies the HasEdge predicate on the "records" edge with a given conditions (other predicates).
func HasRecordsWith(preds ...predicate.TradeRecord) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(RecordsInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, RecordsTable, RecordsColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.TradeTimeRange) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.TradeTimeRange) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.TradeTimeRange) predicate.TradeTimeRange {
	return predicate.TradeTimeRange(func(s *sql.Selector) {
		p(s.Not())
	})
}
