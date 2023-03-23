package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

//go generate ./ent
//go run -mod=mod entgo.io/ent/cmd/ent init BarRecord

// BarRecord holds the schema definition for the BarRecord entity.
type BarRecord struct {
	ent.Schema
}

// Fields of the StockRangeRecord.
func (BarRecord) Fields() []ent.Field {
	return []ent.Field{
		field.Float("close"),
		field.Float("high"),
		field.Float("low"),
		field.Float("open"),
		field.Time("timestamp").
			SchemaType(dateTimeSchemaType),
		field.Float("volume"),
		field.Int32("transactions").
			Comment("the number of trades during this bar"),
	}
}

// Edges of the StockRangeRecord.
func (BarRecord) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("group", BarGroup.Type).
			Ref("records").Unique(),
	}
}
