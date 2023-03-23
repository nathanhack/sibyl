package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// MarketHours holds the schema definition for the MarketHours entity.
type MarketHours struct {
	ent.Schema
}

// Fields of the MarketHours.
func (MarketHours) Fields() []ent.Field {
	return []ent.Field{
		field.Time("date").
			SchemaType(dateTimeSchemaType), // From Alpaca
		field.Time("start_time").
			SchemaType(dateTimeSchemaType), // From Alpaca
		field.Time("end_time").
			SchemaType(dateTimeSchemaType), // From Alpaca
	}
}

// Edges of the MarketHours.
func (MarketHours) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("market_info", MarketInfo.Type).
			Ref("hours").Unique(),
	}
}
