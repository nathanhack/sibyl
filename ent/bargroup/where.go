// Code generated by ent, DO NOT EDIT.

package bargroup

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/nathanhack/sibyl/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldLTE(FieldID, id))
}

// First applies equality check predicate on the "first" field. It's identical to FirstEQ.
func First(v time.Time) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldEQ(FieldFirst, v))
}

// Last applies equality check predicate on the "last" field. It's identical to LastEQ.
func Last(v time.Time) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldEQ(FieldLast, v))
}

// Count applies equality check predicate on the "count" field. It's identical to CountEQ.
func Count(v int) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldEQ(FieldCount, v))
}

// TimeRangeID applies equality check predicate on the "time_range_id" field. It's identical to TimeRangeIDEQ.
func TimeRangeID(v int) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldEQ(FieldTimeRangeID, v))
}

// FirstEQ applies the EQ predicate on the "first" field.
func FirstEQ(v time.Time) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldEQ(FieldFirst, v))
}

// FirstNEQ applies the NEQ predicate on the "first" field.
func FirstNEQ(v time.Time) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldNEQ(FieldFirst, v))
}

// FirstIn applies the In predicate on the "first" field.
func FirstIn(vs ...time.Time) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldIn(FieldFirst, vs...))
}

// FirstNotIn applies the NotIn predicate on the "first" field.
func FirstNotIn(vs ...time.Time) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldNotIn(FieldFirst, vs...))
}

// FirstGT applies the GT predicate on the "first" field.
func FirstGT(v time.Time) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldGT(FieldFirst, v))
}

// FirstGTE applies the GTE predicate on the "first" field.
func FirstGTE(v time.Time) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldGTE(FieldFirst, v))
}

// FirstLT applies the LT predicate on the "first" field.
func FirstLT(v time.Time) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldLT(FieldFirst, v))
}

// FirstLTE applies the LTE predicate on the "first" field.
func FirstLTE(v time.Time) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldLTE(FieldFirst, v))
}

// LastEQ applies the EQ predicate on the "last" field.
func LastEQ(v time.Time) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldEQ(FieldLast, v))
}

// LastNEQ applies the NEQ predicate on the "last" field.
func LastNEQ(v time.Time) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldNEQ(FieldLast, v))
}

// LastIn applies the In predicate on the "last" field.
func LastIn(vs ...time.Time) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldIn(FieldLast, vs...))
}

// LastNotIn applies the NotIn predicate on the "last" field.
func LastNotIn(vs ...time.Time) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldNotIn(FieldLast, vs...))
}

// LastGT applies the GT predicate on the "last" field.
func LastGT(v time.Time) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldGT(FieldLast, v))
}

// LastGTE applies the GTE predicate on the "last" field.
func LastGTE(v time.Time) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldGTE(FieldLast, v))
}

// LastLT applies the LT predicate on the "last" field.
func LastLT(v time.Time) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldLT(FieldLast, v))
}

// LastLTE applies the LTE predicate on the "last" field.
func LastLTE(v time.Time) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldLTE(FieldLast, v))
}

// CountEQ applies the EQ predicate on the "count" field.
func CountEQ(v int) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldEQ(FieldCount, v))
}

// CountNEQ applies the NEQ predicate on the "count" field.
func CountNEQ(v int) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldNEQ(FieldCount, v))
}

// CountIn applies the In predicate on the "count" field.
func CountIn(vs ...int) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldIn(FieldCount, vs...))
}

// CountNotIn applies the NotIn predicate on the "count" field.
func CountNotIn(vs ...int) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldNotIn(FieldCount, vs...))
}

// CountGT applies the GT predicate on the "count" field.
func CountGT(v int) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldGT(FieldCount, v))
}

// CountGTE applies the GTE predicate on the "count" field.
func CountGTE(v int) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldGTE(FieldCount, v))
}

// CountLT applies the LT predicate on the "count" field.
func CountLT(v int) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldLT(FieldCount, v))
}

// CountLTE applies the LTE predicate on the "count" field.
func CountLTE(v int) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldLTE(FieldCount, v))
}

// TimeRangeIDEQ applies the EQ predicate on the "time_range_id" field.
func TimeRangeIDEQ(v int) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldEQ(FieldTimeRangeID, v))
}

// TimeRangeIDNEQ applies the NEQ predicate on the "time_range_id" field.
func TimeRangeIDNEQ(v int) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldNEQ(FieldTimeRangeID, v))
}

// TimeRangeIDIn applies the In predicate on the "time_range_id" field.
func TimeRangeIDIn(vs ...int) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldIn(FieldTimeRangeID, vs...))
}

// TimeRangeIDNotIn applies the NotIn predicate on the "time_range_id" field.
func TimeRangeIDNotIn(vs ...int) predicate.BarGroup {
	return predicate.BarGroup(sql.FieldNotIn(FieldTimeRangeID, vs...))
}

// HasTimeRange applies the HasEdge predicate on the "time_range" edge.
func HasTimeRange() predicate.BarGroup {
	return predicate.BarGroup(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, TimeRangeTable, TimeRangeColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasTimeRangeWith applies the HasEdge predicate on the "time_range" edge with a given conditions (other predicates).
func HasTimeRangeWith(preds ...predicate.BarTimeRange) predicate.BarGroup {
	return predicate.BarGroup(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(TimeRangeInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, TimeRangeTable, TimeRangeColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasRecords applies the HasEdge predicate on the "records" edge.
func HasRecords() predicate.BarGroup {
	return predicate.BarGroup(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, RecordsTable, RecordsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasRecordsWith applies the HasEdge predicate on the "records" edge with a given conditions (other predicates).
func HasRecordsWith(preds ...predicate.BarRecord) predicate.BarGroup {
	return predicate.BarGroup(func(s *sql.Selector) {
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
func And(predicates ...predicate.BarGroup) predicate.BarGroup {
	return predicate.BarGroup(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.BarGroup) predicate.BarGroup {
	return predicate.BarGroup(func(s *sql.Selector) {
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
func Not(p predicate.BarGroup) predicate.BarGroup {
	return predicate.BarGroup(func(s *sql.Selector) {
		p(s.Not())
	})
}
