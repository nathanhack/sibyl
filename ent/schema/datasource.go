package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// DataSource holds the schema definition for the DataSource entity.
type DataSource struct {
	ent.Schema
}

// Fields of the DataSource.
func (DataSource) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().Unique(),
		field.String("address").Default(""),
	}
}

// Edges of the DataSource.
func (DataSource) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("intervals", Interval.Type), // which intervals it will download
	}
}
