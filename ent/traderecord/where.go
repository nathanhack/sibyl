// Code generated by ent, DO NOT EDIT.

package traderecord

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/nathanhack/sibyl/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldLTE(FieldID, id))
}

// Price applies equality check predicate on the "price" field. It's identical to PriceEQ.
func Price(v float64) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldEQ(FieldPrice, v))
}

// Timestamp applies equality check predicate on the "timestamp" field. It's identical to TimestampEQ.
func Timestamp(v time.Time) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldEQ(FieldTimestamp, v))
}

// Volume applies equality check predicate on the "volume" field. It's identical to VolumeEQ.
func Volume(v int32) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldEQ(FieldVolume, v))
}

// TimeRangeID applies equality check predicate on the "time_range_id" field. It's identical to TimeRangeIDEQ.
func TimeRangeID(v int) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldEQ(FieldTimeRangeID, v))
}

// PriceEQ applies the EQ predicate on the "price" field.
func PriceEQ(v float64) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldEQ(FieldPrice, v))
}

// PriceNEQ applies the NEQ predicate on the "price" field.
func PriceNEQ(v float64) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldNEQ(FieldPrice, v))
}

// PriceIn applies the In predicate on the "price" field.
func PriceIn(vs ...float64) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldIn(FieldPrice, vs...))
}

// PriceNotIn applies the NotIn predicate on the "price" field.
func PriceNotIn(vs ...float64) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldNotIn(FieldPrice, vs...))
}

// PriceGT applies the GT predicate on the "price" field.
func PriceGT(v float64) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldGT(FieldPrice, v))
}

// PriceGTE applies the GTE predicate on the "price" field.
func PriceGTE(v float64) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldGTE(FieldPrice, v))
}

// PriceLT applies the LT predicate on the "price" field.
func PriceLT(v float64) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldLT(FieldPrice, v))
}

// PriceLTE applies the LTE predicate on the "price" field.
func PriceLTE(v float64) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldLTE(FieldPrice, v))
}

// TimestampEQ applies the EQ predicate on the "timestamp" field.
func TimestampEQ(v time.Time) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldEQ(FieldTimestamp, v))
}

// TimestampNEQ applies the NEQ predicate on the "timestamp" field.
func TimestampNEQ(v time.Time) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldNEQ(FieldTimestamp, v))
}

// TimestampIn applies the In predicate on the "timestamp" field.
func TimestampIn(vs ...time.Time) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldIn(FieldTimestamp, vs...))
}

// TimestampNotIn applies the NotIn predicate on the "timestamp" field.
func TimestampNotIn(vs ...time.Time) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldNotIn(FieldTimestamp, vs...))
}

// TimestampGT applies the GT predicate on the "timestamp" field.
func TimestampGT(v time.Time) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldGT(FieldTimestamp, v))
}

// TimestampGTE applies the GTE predicate on the "timestamp" field.
func TimestampGTE(v time.Time) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldGTE(FieldTimestamp, v))
}

// TimestampLT applies the LT predicate on the "timestamp" field.
func TimestampLT(v time.Time) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldLT(FieldTimestamp, v))
}

// TimestampLTE applies the LTE predicate on the "timestamp" field.
func TimestampLTE(v time.Time) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldLTE(FieldTimestamp, v))
}

// VolumeEQ applies the EQ predicate on the "volume" field.
func VolumeEQ(v int32) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldEQ(FieldVolume, v))
}

// VolumeNEQ applies the NEQ predicate on the "volume" field.
func VolumeNEQ(v int32) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldNEQ(FieldVolume, v))
}

// VolumeIn applies the In predicate on the "volume" field.
func VolumeIn(vs ...int32) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldIn(FieldVolume, vs...))
}

// VolumeNotIn applies the NotIn predicate on the "volume" field.
func VolumeNotIn(vs ...int32) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldNotIn(FieldVolume, vs...))
}

// VolumeGT applies the GT predicate on the "volume" field.
func VolumeGT(v int32) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldGT(FieldVolume, v))
}

// VolumeGTE applies the GTE predicate on the "volume" field.
func VolumeGTE(v int32) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldGTE(FieldVolume, v))
}

// VolumeLT applies the LT predicate on the "volume" field.
func VolumeLT(v int32) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldLT(FieldVolume, v))
}

// VolumeLTE applies the LTE predicate on the "volume" field.
func VolumeLTE(v int32) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldLTE(FieldVolume, v))
}

// TimeRangeIDEQ applies the EQ predicate on the "time_range_id" field.
func TimeRangeIDEQ(v int) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldEQ(FieldTimeRangeID, v))
}

// TimeRangeIDNEQ applies the NEQ predicate on the "time_range_id" field.
func TimeRangeIDNEQ(v int) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldNEQ(FieldTimeRangeID, v))
}

// TimeRangeIDIn applies the In predicate on the "time_range_id" field.
func TimeRangeIDIn(vs ...int) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldIn(FieldTimeRangeID, vs...))
}

// TimeRangeIDNotIn applies the NotIn predicate on the "time_range_id" field.
func TimeRangeIDNotIn(vs ...int) predicate.TradeRecord {
	return predicate.TradeRecord(sql.FieldNotIn(FieldTimeRangeID, vs...))
}

// HasTimeRange applies the HasEdge predicate on the "time_range" edge.
func HasTimeRange() predicate.TradeRecord {
	return predicate.TradeRecord(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, TimeRangeTable, TimeRangeColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasTimeRangeWith applies the HasEdge predicate on the "time_range" edge with a given conditions (other predicates).
func HasTimeRangeWith(preds ...predicate.TradeTimeRange) predicate.TradeRecord {
	return predicate.TradeRecord(func(s *sql.Selector) {
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

// HasConditions applies the HasEdge predicate on the "conditions" edge.
func HasConditions() predicate.TradeRecord {
	return predicate.TradeRecord(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, ConditionsTable, ConditionsPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasConditionsWith applies the HasEdge predicate on the "conditions" edge with a given conditions (other predicates).
func HasConditionsWith(preds ...predicate.TradeCondition) predicate.TradeRecord {
	return predicate.TradeRecord(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ConditionsInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, ConditionsTable, ConditionsPrimaryKey...),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasCorrection applies the HasEdge predicate on the "correction" edge.
func HasCorrection() predicate.TradeRecord {
	return predicate.TradeRecord(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, CorrectionTable, CorrectionPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasCorrectionWith applies the HasEdge predicate on the "correction" edge with a given conditions (other predicates).
func HasCorrectionWith(preds ...predicate.TradeCorrection) predicate.TradeRecord {
	return predicate.TradeRecord(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(CorrectionInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, CorrectionTable, CorrectionPrimaryKey...),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasExchange applies the HasEdge predicate on the "exchange" edge.
func HasExchange() predicate.TradeRecord {
	return predicate.TradeRecord(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, ExchangeTable, ExchangeColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasExchangeWith applies the HasEdge predicate on the "exchange" edge with a given conditions (other predicates).
func HasExchangeWith(preds ...predicate.Exchange) predicate.TradeRecord {
	return predicate.TradeRecord(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ExchangeInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, ExchangeTable, ExchangeColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.TradeRecord) predicate.TradeRecord {
	return predicate.TradeRecord(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.TradeRecord) predicate.TradeRecord {
	return predicate.TradeRecord(func(s *sql.Selector) {
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
func Not(p predicate.TradeRecord) predicate.TradeRecord {
	return predicate.TradeRecord(func(s *sql.Selector) {
		p(s.Not())
	})
}
