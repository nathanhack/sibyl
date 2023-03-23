package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

//go generate ./ent
//go run -mod=mod entgo.io/ent/cmd/ent init StockPointRecord

// TradeRecord holds the schema definition for the TradeRecord entity.
type TradeRecord struct {
	ent.Schema
}

// Fields of the StockPointRecord.
func (TradeRecord) Fields() []ent.Field {
	return []ent.Field{
		field.Float("price"),
		field.Time("timestamp").
			SchemaType(dateTimeSchemaType),
		field.Int32("volume"),
		field.Int("time_range_id"),
	}
}

// Edges of the StockPointRecord.
func (TradeRecord) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("time_range", TradeTimeRange.Type).
			Ref("records").
			Field("time_range_id").
			Unique().Required(),

		edge.To("conditions", TradeCondition.Type),
		edge.To("correction", TradeCorrection.Type),
		edge.To("exchange", Exchange.Type),
	}
}

func (TradeRecord) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("time_range_id", "timestamp").Unique(),
	}
}
