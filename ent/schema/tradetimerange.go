package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// TradeTimeRange holds the schema definition for the TradeTimeRange entity.
type TradeTimeRange struct {
	ent.Schema
}

// Fields of the TradeTimeRange.
func (TradeTimeRange) Fields() []ent.Field {
	return []ent.Field{
		field.Time("start").
			SchemaType(dateTimeSchemaType),
		field.Time("end").
			SchemaType(dateTimeSchemaType),
		field.Int("interval_id"),
	}
}

// Edges of the TradeTimeRange.
func (TradeTimeRange) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("interval", Interval.Type).
			Ref("trades").
			Field("interval_id").
			Unique().Required(),

		edge.To("records", TradeRecord.Type),
	}
}

func (TradeTimeRange) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("interval_id", "start", "end").Unique(),
	}
}
