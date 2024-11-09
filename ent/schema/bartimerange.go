package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// BarTimeRange holds the schema definition for the BarTimeRange entity. The over arching time range. Each will have a set of groups
// which consist of more accurate time ranges. Ex the BarTimeRange(1970-Now) would have a Group(for mins for range 0900-1600)
type BarTimeRange struct {
	ent.Schema
}

func (BarTimeRange) Fields() []ent.Field {
	return []ent.Field{
		field.Time("start").
			SchemaType(dateTimeSchemaType),
		field.Time("end").
			SchemaType(dateTimeSchemaType),
		field.Int("count").Default(0).Comment("The number of BarGroups"),
		field.Int("interval_id"),
		field.Enum("status").
			Values("pending", "created", "clean", "consolidated").
			Default("pending"),
		field.Time("update_time").
			Default(time.Now).
			UpdateDefault(time.Now).
			Annotations(entsql.Annotation{Default: "CURRENT_TIMESTAMP"}).
			SchemaType(dateTimeSchemaType),
	}
}

func (BarTimeRange) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("interval", Interval.Type).
			Ref("bars").
			Field("interval_id").
			Unique().Required(),

		edge.To("groups", BarGroup.Type),
	}
}

func (BarTimeRange) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("interval_id", "start", "end").Unique(),
		index.Fields("status"),
	}
}
