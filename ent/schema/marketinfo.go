package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// MarketInfo holds the schema definition for the MarketInfo entity.
type MarketInfo struct {
	ent.Schema
}

// Fields of the MarketHours.
func (MarketInfo) Fields() []ent.Field {
	return []ent.Field{
		field.Time("hours_start").
			SchemaType(dateTimeSchemaType), // Start of MarketHours interval
		field.Time("hours_end").
			SchemaType(dateTimeSchemaType), //  End of MarketHours interval (last day) inclusive
	}
}

// Edges of the MarketInfo.
func (MarketInfo) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("hours", MarketHours.Type),
	}
}
