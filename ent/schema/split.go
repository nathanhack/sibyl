package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Split holds the schema definition for the Split entity.
type Split struct {
	ent.Schema
}

// Fields of the Split.
func (Split) Fields() []ent.Field {
	return []ent.Field{
		field.Time("execution_date").
			SchemaType(dateTimeSchemaType), //From Polygon.io
		field.Float("from"), //From Polygon.io
		field.Float("to"),   //From Polygon.io
	}
}

// Edges of the Split.
func (Split) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("stock", Entity.Type).Ref("splits").Unique(),
	}
}
