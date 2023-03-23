package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TradeCondition holds the schema definition for the TradeCondition entity.
type TradeCondition struct {
	ent.Schema
}

// Fields of the StockPointCondition.
func (TradeCondition) Fields() []ent.Field {
	return []ent.Field{
		field.String("condition"),
	}
}

// Edges of the StockPointCondition.
func (TradeCondition) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("record", TradeRecord.Type).Ref("conditions"),
	}
}
