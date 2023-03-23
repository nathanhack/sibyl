package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// BarGroup holds the schema definition for the BarGroup entity.  These are more accurate time ranges. For "Day" interval
// these ranges would not contain weekends nor holidays. For "Min" interval it would only be the time during the day that data was available.

type BarGroup struct {
	ent.Schema
}

// Fields of the BarGroup.
func (BarGroup) Fields() []ent.Field {
	return []ent.Field{
		field.Time("first").
			SchemaType(dateTimeSchemaType),
		field.Time("last").
			SchemaType(dateTimeSchemaType),
		field.Int("count"),
		field.Int("time_range_id"),
	}
}

// Edges of the BarGroup.
func (BarGroup) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("time_range", BarTimeRange.Type).
			Ref("groups").
			Field("time_range_id").
			Unique().Required(),

		edge.To("records", BarRecord.Type),
	}
}

func (BarGroup) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("time_range_id", "first", "last").Unique(),
	}
}
