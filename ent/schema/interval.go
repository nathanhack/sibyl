package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

//go run -mod=mod entgo.io/ent/cmd/ent init Interval

// Interval holds the schema definition for the Interval entity.
type Interval struct {
	ent.Schema
}

// Fields of the HistoryInterval.
func (Interval) Fields() []ent.Field {
	return []ent.Field{
		field.Bool("active").Default(false),
		field.Enum("interval").Values("trades", "1min", "daily", "monthly", "yearly"),
		field.Int("stock_id"),
		field.Int("data_source_id"),
	}
}

// Edges of the HistoryInterval.
func (Interval) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("data_source", DataSource.Type).
			Ref("intervals").
			Field("data_source_id").
			Unique().
			Required(),
		edge.From("stock", Entity.Type).
			Ref("intervals").
			Field("stock_id").
			Unique().
			Required(),
		edge.To("bars", BarTimeRange.Type),
		edge.To("trades", TradeTimeRange.Type),
	}
}

func (Interval) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("stock_id", "data_source_id", "interval").Unique(),
	}
}
