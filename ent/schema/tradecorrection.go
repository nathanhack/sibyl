package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TradeCorrection holds the schema definition for the TradeCorrection entity.
type TradeCorrection struct {
	ent.Schema
}

// Fields of the StockPointCorrection.
func (TradeCorrection) Fields() []ent.Field {
	return []ent.Field{
		field.String("correction"),
	}
}

// Edges of the StockPointCorrection.
func (TradeCorrection) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("record", TradeRecord.Type).Ref("correction"),
	}
}
