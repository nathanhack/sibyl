// Code generated by ent, DO NOT EDIT.

package tradecondition

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/nathanhack/sibyl/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.TradeCondition {
	return predicate.TradeCondition(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.TradeCondition {
	return predicate.TradeCondition(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.TradeCondition {
	return predicate.TradeCondition(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.TradeCondition {
	return predicate.TradeCondition(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.TradeCondition {
	return predicate.TradeCondition(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.TradeCondition {
	return predicate.TradeCondition(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.TradeCondition {
	return predicate.TradeCondition(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.TradeCondition {
	return predicate.TradeCondition(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.TradeCondition {
	return predicate.TradeCondition(sql.FieldLTE(FieldID, id))
}

// Condition applies equality check predicate on the "condition" field. It's identical to ConditionEQ.
func Condition(v string) predicate.TradeCondition {
	return predicate.TradeCondition(sql.FieldEQ(FieldCondition, v))
}

// ConditionEQ applies the EQ predicate on the "condition" field.
func ConditionEQ(v string) predicate.TradeCondition {
	return predicate.TradeCondition(sql.FieldEQ(FieldCondition, v))
}

// ConditionNEQ applies the NEQ predicate on the "condition" field.
func ConditionNEQ(v string) predicate.TradeCondition {
	return predicate.TradeCondition(sql.FieldNEQ(FieldCondition, v))
}

// ConditionIn applies the In predicate on the "condition" field.
func ConditionIn(vs ...string) predicate.TradeCondition {
	return predicate.TradeCondition(sql.FieldIn(FieldCondition, vs...))
}

// ConditionNotIn applies the NotIn predicate on the "condition" field.
func ConditionNotIn(vs ...string) predicate.TradeCondition {
	return predicate.TradeCondition(sql.FieldNotIn(FieldCondition, vs...))
}

// ConditionGT applies the GT predicate on the "condition" field.
func ConditionGT(v string) predicate.TradeCondition {
	return predicate.TradeCondition(sql.FieldGT(FieldCondition, v))
}

// ConditionGTE applies the GTE predicate on the "condition" field.
func ConditionGTE(v string) predicate.TradeCondition {
	return predicate.TradeCondition(sql.FieldGTE(FieldCondition, v))
}

// ConditionLT applies the LT predicate on the "condition" field.
func ConditionLT(v string) predicate.TradeCondition {
	return predicate.TradeCondition(sql.FieldLT(FieldCondition, v))
}

// ConditionLTE applies the LTE predicate on the "condition" field.
func ConditionLTE(v string) predicate.TradeCondition {
	return predicate.TradeCondition(sql.FieldLTE(FieldCondition, v))
}

// ConditionContains applies the Contains predicate on the "condition" field.
func ConditionContains(v string) predicate.TradeCondition {
	return predicate.TradeCondition(sql.FieldContains(FieldCondition, v))
}

// ConditionHasPrefix applies the HasPrefix predicate on the "condition" field.
func ConditionHasPrefix(v string) predicate.TradeCondition {
	return predicate.TradeCondition(sql.FieldHasPrefix(FieldCondition, v))
}

// ConditionHasSuffix applies the HasSuffix predicate on the "condition" field.
func ConditionHasSuffix(v string) predicate.TradeCondition {
	return predicate.TradeCondition(sql.FieldHasSuffix(FieldCondition, v))
}

// ConditionEqualFold applies the EqualFold predicate on the "condition" field.
func ConditionEqualFold(v string) predicate.TradeCondition {
	return predicate.TradeCondition(sql.FieldEqualFold(FieldCondition, v))
}

// ConditionContainsFold applies the ContainsFold predicate on the "condition" field.
func ConditionContainsFold(v string) predicate.TradeCondition {
	return predicate.TradeCondition(sql.FieldContainsFold(FieldCondition, v))
}

// HasRecord applies the HasEdge predicate on the "record" edge.
func HasRecord() predicate.TradeCondition {
	return predicate.TradeCondition(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, RecordTable, RecordPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasRecordWith applies the HasEdge predicate on the "record" edge with a given conditions (other predicates).
func HasRecordWith(preds ...predicate.TradeRecord) predicate.TradeCondition {
	return predicate.TradeCondition(func(s *sql.Selector) {
		step := newRecordStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.TradeCondition) predicate.TradeCondition {
	return predicate.TradeCondition(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.TradeCondition) predicate.TradeCondition {
	return predicate.TradeCondition(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.TradeCondition) predicate.TradeCondition {
	return predicate.TradeCondition(sql.NotPredicates(p))
}
