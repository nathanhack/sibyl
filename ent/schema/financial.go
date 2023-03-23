package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
)

// Financial holds the schema definition for the Financial entity.
type Financial struct {
	ent.Schema
}

// Fields of the Financial.
func (Financial) Fields() []ent.Field {
	return []ent.Field{}
}

// Edges of the Financial.
func (Financial) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("stock", Entity.Type).Ref("financials"),
	}
}
