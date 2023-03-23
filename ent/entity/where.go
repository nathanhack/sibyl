// Code generated by ent, DO NOT EDIT.

package entity

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/nathanhack/sibyl/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Entity {
	return predicate.Entity(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Entity {
	return predicate.Entity(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Entity {
	return predicate.Entity(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Entity {
	return predicate.Entity(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Entity {
	return predicate.Entity(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Entity {
	return predicate.Entity(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Entity {
	return predicate.Entity(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Entity {
	return predicate.Entity(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Entity {
	return predicate.Entity(sql.FieldLTE(FieldID, id))
}

// Active applies equality check predicate on the "active" field. It's identical to ActiveEQ.
func Active(v bool) predicate.Entity {
	return predicate.Entity(sql.FieldEQ(FieldActive, v))
}

// Ticker applies equality check predicate on the "ticker" field. It's identical to TickerEQ.
func Ticker(v string) predicate.Entity {
	return predicate.Entity(sql.FieldEQ(FieldTicker, v))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Entity {
	return predicate.Entity(sql.FieldEQ(FieldName, v))
}

// Description applies equality check predicate on the "description" field. It's identical to DescriptionEQ.
func Description(v string) predicate.Entity {
	return predicate.Entity(sql.FieldEQ(FieldDescription, v))
}

// ListDate applies equality check predicate on the "list_date" field. It's identical to ListDateEQ.
func ListDate(v time.Time) predicate.Entity {
	return predicate.Entity(sql.FieldEQ(FieldListDate, v))
}

// Delisted applies equality check predicate on the "delisted" field. It's identical to DelistedEQ.
func Delisted(v time.Time) predicate.Entity {
	return predicate.Entity(sql.FieldEQ(FieldDelisted, v))
}

// ActiveEQ applies the EQ predicate on the "active" field.
func ActiveEQ(v bool) predicate.Entity {
	return predicate.Entity(sql.FieldEQ(FieldActive, v))
}

// ActiveNEQ applies the NEQ predicate on the "active" field.
func ActiveNEQ(v bool) predicate.Entity {
	return predicate.Entity(sql.FieldNEQ(FieldActive, v))
}

// TickerEQ applies the EQ predicate on the "ticker" field.
func TickerEQ(v string) predicate.Entity {
	return predicate.Entity(sql.FieldEQ(FieldTicker, v))
}

// TickerNEQ applies the NEQ predicate on the "ticker" field.
func TickerNEQ(v string) predicate.Entity {
	return predicate.Entity(sql.FieldNEQ(FieldTicker, v))
}

// TickerIn applies the In predicate on the "ticker" field.
func TickerIn(vs ...string) predicate.Entity {
	return predicate.Entity(sql.FieldIn(FieldTicker, vs...))
}

// TickerNotIn applies the NotIn predicate on the "ticker" field.
func TickerNotIn(vs ...string) predicate.Entity {
	return predicate.Entity(sql.FieldNotIn(FieldTicker, vs...))
}

// TickerGT applies the GT predicate on the "ticker" field.
func TickerGT(v string) predicate.Entity {
	return predicate.Entity(sql.FieldGT(FieldTicker, v))
}

// TickerGTE applies the GTE predicate on the "ticker" field.
func TickerGTE(v string) predicate.Entity {
	return predicate.Entity(sql.FieldGTE(FieldTicker, v))
}

// TickerLT applies the LT predicate on the "ticker" field.
func TickerLT(v string) predicate.Entity {
	return predicate.Entity(sql.FieldLT(FieldTicker, v))
}

// TickerLTE applies the LTE predicate on the "ticker" field.
func TickerLTE(v string) predicate.Entity {
	return predicate.Entity(sql.FieldLTE(FieldTicker, v))
}

// TickerContains applies the Contains predicate on the "ticker" field.
func TickerContains(v string) predicate.Entity {
	return predicate.Entity(sql.FieldContains(FieldTicker, v))
}

// TickerHasPrefix applies the HasPrefix predicate on the "ticker" field.
func TickerHasPrefix(v string) predicate.Entity {
	return predicate.Entity(sql.FieldHasPrefix(FieldTicker, v))
}

// TickerHasSuffix applies the HasSuffix predicate on the "ticker" field.
func TickerHasSuffix(v string) predicate.Entity {
	return predicate.Entity(sql.FieldHasSuffix(FieldTicker, v))
}

// TickerEqualFold applies the EqualFold predicate on the "ticker" field.
func TickerEqualFold(v string) predicate.Entity {
	return predicate.Entity(sql.FieldEqualFold(FieldTicker, v))
}

// TickerContainsFold applies the ContainsFold predicate on the "ticker" field.
func TickerContainsFold(v string) predicate.Entity {
	return predicate.Entity(sql.FieldContainsFold(FieldTicker, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Entity {
	return predicate.Entity(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Entity {
	return predicate.Entity(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Entity {
	return predicate.Entity(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Entity {
	return predicate.Entity(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Entity {
	return predicate.Entity(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Entity {
	return predicate.Entity(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Entity {
	return predicate.Entity(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Entity {
	return predicate.Entity(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Entity {
	return predicate.Entity(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Entity {
	return predicate.Entity(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Entity {
	return predicate.Entity(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Entity {
	return predicate.Entity(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Entity {
	return predicate.Entity(sql.FieldContainsFold(FieldName, v))
}

// DescriptionEQ applies the EQ predicate on the "description" field.
func DescriptionEQ(v string) predicate.Entity {
	return predicate.Entity(sql.FieldEQ(FieldDescription, v))
}

// DescriptionNEQ applies the NEQ predicate on the "description" field.
func DescriptionNEQ(v string) predicate.Entity {
	return predicate.Entity(sql.FieldNEQ(FieldDescription, v))
}

// DescriptionIn applies the In predicate on the "description" field.
func DescriptionIn(vs ...string) predicate.Entity {
	return predicate.Entity(sql.FieldIn(FieldDescription, vs...))
}

// DescriptionNotIn applies the NotIn predicate on the "description" field.
func DescriptionNotIn(vs ...string) predicate.Entity {
	return predicate.Entity(sql.FieldNotIn(FieldDescription, vs...))
}

// DescriptionGT applies the GT predicate on the "description" field.
func DescriptionGT(v string) predicate.Entity {
	return predicate.Entity(sql.FieldGT(FieldDescription, v))
}

// DescriptionGTE applies the GTE predicate on the "description" field.
func DescriptionGTE(v string) predicate.Entity {
	return predicate.Entity(sql.FieldGTE(FieldDescription, v))
}

// DescriptionLT applies the LT predicate on the "description" field.
func DescriptionLT(v string) predicate.Entity {
	return predicate.Entity(sql.FieldLT(FieldDescription, v))
}

// DescriptionLTE applies the LTE predicate on the "description" field.
func DescriptionLTE(v string) predicate.Entity {
	return predicate.Entity(sql.FieldLTE(FieldDescription, v))
}

// DescriptionContains applies the Contains predicate on the "description" field.
func DescriptionContains(v string) predicate.Entity {
	return predicate.Entity(sql.FieldContains(FieldDescription, v))
}

// DescriptionHasPrefix applies the HasPrefix predicate on the "description" field.
func DescriptionHasPrefix(v string) predicate.Entity {
	return predicate.Entity(sql.FieldHasPrefix(FieldDescription, v))
}

// DescriptionHasSuffix applies the HasSuffix predicate on the "description" field.
func DescriptionHasSuffix(v string) predicate.Entity {
	return predicate.Entity(sql.FieldHasSuffix(FieldDescription, v))
}

// DescriptionEqualFold applies the EqualFold predicate on the "description" field.
func DescriptionEqualFold(v string) predicate.Entity {
	return predicate.Entity(sql.FieldEqualFold(FieldDescription, v))
}

// DescriptionContainsFold applies the ContainsFold predicate on the "description" field.
func DescriptionContainsFold(v string) predicate.Entity {
	return predicate.Entity(sql.FieldContainsFold(FieldDescription, v))
}

// ListDateEQ applies the EQ predicate on the "list_date" field.
func ListDateEQ(v time.Time) predicate.Entity {
	return predicate.Entity(sql.FieldEQ(FieldListDate, v))
}

// ListDateNEQ applies the NEQ predicate on the "list_date" field.
func ListDateNEQ(v time.Time) predicate.Entity {
	return predicate.Entity(sql.FieldNEQ(FieldListDate, v))
}

// ListDateIn applies the In predicate on the "list_date" field.
func ListDateIn(vs ...time.Time) predicate.Entity {
	return predicate.Entity(sql.FieldIn(FieldListDate, vs...))
}

// ListDateNotIn applies the NotIn predicate on the "list_date" field.
func ListDateNotIn(vs ...time.Time) predicate.Entity {
	return predicate.Entity(sql.FieldNotIn(FieldListDate, vs...))
}

// ListDateGT applies the GT predicate on the "list_date" field.
func ListDateGT(v time.Time) predicate.Entity {
	return predicate.Entity(sql.FieldGT(FieldListDate, v))
}

// ListDateGTE applies the GTE predicate on the "list_date" field.
func ListDateGTE(v time.Time) predicate.Entity {
	return predicate.Entity(sql.FieldGTE(FieldListDate, v))
}

// ListDateLT applies the LT predicate on the "list_date" field.
func ListDateLT(v time.Time) predicate.Entity {
	return predicate.Entity(sql.FieldLT(FieldListDate, v))
}

// ListDateLTE applies the LTE predicate on the "list_date" field.
func ListDateLTE(v time.Time) predicate.Entity {
	return predicate.Entity(sql.FieldLTE(FieldListDate, v))
}

// DelistedEQ applies the EQ predicate on the "delisted" field.
func DelistedEQ(v time.Time) predicate.Entity {
	return predicate.Entity(sql.FieldEQ(FieldDelisted, v))
}

// DelistedNEQ applies the NEQ predicate on the "delisted" field.
func DelistedNEQ(v time.Time) predicate.Entity {
	return predicate.Entity(sql.FieldNEQ(FieldDelisted, v))
}

// DelistedIn applies the In predicate on the "delisted" field.
func DelistedIn(vs ...time.Time) predicate.Entity {
	return predicate.Entity(sql.FieldIn(FieldDelisted, vs...))
}

// DelistedNotIn applies the NotIn predicate on the "delisted" field.
func DelistedNotIn(vs ...time.Time) predicate.Entity {
	return predicate.Entity(sql.FieldNotIn(FieldDelisted, vs...))
}

// DelistedGT applies the GT predicate on the "delisted" field.
func DelistedGT(v time.Time) predicate.Entity {
	return predicate.Entity(sql.FieldGT(FieldDelisted, v))
}

// DelistedGTE applies the GTE predicate on the "delisted" field.
func DelistedGTE(v time.Time) predicate.Entity {
	return predicate.Entity(sql.FieldGTE(FieldDelisted, v))
}

// DelistedLT applies the LT predicate on the "delisted" field.
func DelistedLT(v time.Time) predicate.Entity {
	return predicate.Entity(sql.FieldLT(FieldDelisted, v))
}

// DelistedLTE applies the LTE predicate on the "delisted" field.
func DelistedLTE(v time.Time) predicate.Entity {
	return predicate.Entity(sql.FieldLTE(FieldDelisted, v))
}

// DelistedIsNil applies the IsNil predicate on the "delisted" field.
func DelistedIsNil() predicate.Entity {
	return predicate.Entity(sql.FieldIsNull(FieldDelisted))
}

// DelistedNotNil applies the NotNil predicate on the "delisted" field.
func DelistedNotNil() predicate.Entity {
	return predicate.Entity(sql.FieldNotNull(FieldDelisted))
}

// HasExchanges applies the HasEdge predicate on the "exchanges" edge.
func HasExchanges() predicate.Entity {
	return predicate.Entity(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, ExchangesTable, ExchangesPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasExchangesWith applies the HasEdge predicate on the "exchanges" edge with a given conditions (other predicates).
func HasExchangesWith(preds ...predicate.Exchange) predicate.Entity {
	return predicate.Entity(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ExchangesInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, ExchangesTable, ExchangesPrimaryKey...),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasIntervals applies the HasEdge predicate on the "intervals" edge.
func HasIntervals() predicate.Entity {
	return predicate.Entity(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, IntervalsTable, IntervalsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasIntervalsWith applies the HasEdge predicate on the "intervals" edge with a given conditions (other predicates).
func HasIntervalsWith(preds ...predicate.Interval) predicate.Entity {
	return predicate.Entity(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(IntervalsInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, IntervalsTable, IntervalsColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasDividends applies the HasEdge predicate on the "dividends" edge.
func HasDividends() predicate.Entity {
	return predicate.Entity(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, DividendsTable, DividendsPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasDividendsWith applies the HasEdge predicate on the "dividends" edge with a given conditions (other predicates).
func HasDividendsWith(preds ...predicate.Dividend) predicate.Entity {
	return predicate.Entity(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(DividendsInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, DividendsTable, DividendsPrimaryKey...),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasSplits applies the HasEdge predicate on the "splits" edge.
func HasSplits() predicate.Entity {
	return predicate.Entity(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, SplitsTable, SplitsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasSplitsWith applies the HasEdge predicate on the "splits" edge with a given conditions (other predicates).
func HasSplitsWith(preds ...predicate.Split) predicate.Entity {
	return predicate.Entity(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(SplitsInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, SplitsTable, SplitsColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasFinancials applies the HasEdge predicate on the "financials" edge.
func HasFinancials() predicate.Entity {
	return predicate.Entity(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, FinancialsTable, FinancialsPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasFinancialsWith applies the HasEdge predicate on the "financials" edge with a given conditions (other predicates).
func HasFinancialsWith(preds ...predicate.Financial) predicate.Entity {
	return predicate.Entity(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(FinancialsInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, FinancialsTable, FinancialsPrimaryKey...),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Entity) predicate.Entity {
	return predicate.Entity(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Entity) predicate.Entity {
	return predicate.Entity(func(s *sql.Selector) {
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
func Not(p predicate.Entity) predicate.Entity {
	return predicate.Entity(func(s *sql.Selector) {
		p(s.Not())
	})
}