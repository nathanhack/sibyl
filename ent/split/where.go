// Code generated by ent, DO NOT EDIT.

package split

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/nathanhack/sibyl/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Split {
	return predicate.Split(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Split {
	return predicate.Split(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Split {
	return predicate.Split(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Split {
	return predicate.Split(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Split {
	return predicate.Split(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Split {
	return predicate.Split(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Split {
	return predicate.Split(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Split {
	return predicate.Split(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Split {
	return predicate.Split(sql.FieldLTE(FieldID, id))
}

// ExecutionDate applies equality check predicate on the "execution_date" field. It's identical to ExecutionDateEQ.
func ExecutionDate(v time.Time) predicate.Split {
	return predicate.Split(sql.FieldEQ(FieldExecutionDate, v))
}

// From applies equality check predicate on the "from" field. It's identical to FromEQ.
func From(v float64) predicate.Split {
	return predicate.Split(sql.FieldEQ(FieldFrom, v))
}

// To applies equality check predicate on the "to" field. It's identical to ToEQ.
func To(v float64) predicate.Split {
	return predicate.Split(sql.FieldEQ(FieldTo, v))
}

// ExecutionDateEQ applies the EQ predicate on the "execution_date" field.
func ExecutionDateEQ(v time.Time) predicate.Split {
	return predicate.Split(sql.FieldEQ(FieldExecutionDate, v))
}

// ExecutionDateNEQ applies the NEQ predicate on the "execution_date" field.
func ExecutionDateNEQ(v time.Time) predicate.Split {
	return predicate.Split(sql.FieldNEQ(FieldExecutionDate, v))
}

// ExecutionDateIn applies the In predicate on the "execution_date" field.
func ExecutionDateIn(vs ...time.Time) predicate.Split {
	return predicate.Split(sql.FieldIn(FieldExecutionDate, vs...))
}

// ExecutionDateNotIn applies the NotIn predicate on the "execution_date" field.
func ExecutionDateNotIn(vs ...time.Time) predicate.Split {
	return predicate.Split(sql.FieldNotIn(FieldExecutionDate, vs...))
}

// ExecutionDateGT applies the GT predicate on the "execution_date" field.
func ExecutionDateGT(v time.Time) predicate.Split {
	return predicate.Split(sql.FieldGT(FieldExecutionDate, v))
}

// ExecutionDateGTE applies the GTE predicate on the "execution_date" field.
func ExecutionDateGTE(v time.Time) predicate.Split {
	return predicate.Split(sql.FieldGTE(FieldExecutionDate, v))
}

// ExecutionDateLT applies the LT predicate on the "execution_date" field.
func ExecutionDateLT(v time.Time) predicate.Split {
	return predicate.Split(sql.FieldLT(FieldExecutionDate, v))
}

// ExecutionDateLTE applies the LTE predicate on the "execution_date" field.
func ExecutionDateLTE(v time.Time) predicate.Split {
	return predicate.Split(sql.FieldLTE(FieldExecutionDate, v))
}

// FromEQ applies the EQ predicate on the "from" field.
func FromEQ(v float64) predicate.Split {
	return predicate.Split(sql.FieldEQ(FieldFrom, v))
}

// FromNEQ applies the NEQ predicate on the "from" field.
func FromNEQ(v float64) predicate.Split {
	return predicate.Split(sql.FieldNEQ(FieldFrom, v))
}

// FromIn applies the In predicate on the "from" field.
func FromIn(vs ...float64) predicate.Split {
	return predicate.Split(sql.FieldIn(FieldFrom, vs...))
}

// FromNotIn applies the NotIn predicate on the "from" field.
func FromNotIn(vs ...float64) predicate.Split {
	return predicate.Split(sql.FieldNotIn(FieldFrom, vs...))
}

// FromGT applies the GT predicate on the "from" field.
func FromGT(v float64) predicate.Split {
	return predicate.Split(sql.FieldGT(FieldFrom, v))
}

// FromGTE applies the GTE predicate on the "from" field.
func FromGTE(v float64) predicate.Split {
	return predicate.Split(sql.FieldGTE(FieldFrom, v))
}

// FromLT applies the LT predicate on the "from" field.
func FromLT(v float64) predicate.Split {
	return predicate.Split(sql.FieldLT(FieldFrom, v))
}

// FromLTE applies the LTE predicate on the "from" field.
func FromLTE(v float64) predicate.Split {
	return predicate.Split(sql.FieldLTE(FieldFrom, v))
}

// ToEQ applies the EQ predicate on the "to" field.
func ToEQ(v float64) predicate.Split {
	return predicate.Split(sql.FieldEQ(FieldTo, v))
}

// ToNEQ applies the NEQ predicate on the "to" field.
func ToNEQ(v float64) predicate.Split {
	return predicate.Split(sql.FieldNEQ(FieldTo, v))
}

// ToIn applies the In predicate on the "to" field.
func ToIn(vs ...float64) predicate.Split {
	return predicate.Split(sql.FieldIn(FieldTo, vs...))
}

// ToNotIn applies the NotIn predicate on the "to" field.
func ToNotIn(vs ...float64) predicate.Split {
	return predicate.Split(sql.FieldNotIn(FieldTo, vs...))
}

// ToGT applies the GT predicate on the "to" field.
func ToGT(v float64) predicate.Split {
	return predicate.Split(sql.FieldGT(FieldTo, v))
}

// ToGTE applies the GTE predicate on the "to" field.
func ToGTE(v float64) predicate.Split {
	return predicate.Split(sql.FieldGTE(FieldTo, v))
}

// ToLT applies the LT predicate on the "to" field.
func ToLT(v float64) predicate.Split {
	return predicate.Split(sql.FieldLT(FieldTo, v))
}

// ToLTE applies the LTE predicate on the "to" field.
func ToLTE(v float64) predicate.Split {
	return predicate.Split(sql.FieldLTE(FieldTo, v))
}

// HasStock applies the HasEdge predicate on the "stock" edge.
func HasStock() predicate.Split {
	return predicate.Split(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, StockTable, StockColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasStockWith applies the HasEdge predicate on the "stock" edge with a given conditions (other predicates).
func HasStockWith(preds ...predicate.Entity) predicate.Split {
	return predicate.Split(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(StockInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, StockTable, StockColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Split) predicate.Split {
	return predicate.Split(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Split) predicate.Split {
	return predicate.Split(func(s *sql.Selector) {
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
func Not(p predicate.Split) predicate.Split {
	return predicate.Split(func(s *sql.Selector) {
		p(s.Not())
	})
}
