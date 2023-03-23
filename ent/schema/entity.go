package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

//go generate ./...
//go run -mod=mod entgo.io/ent/cmd/ent init Entity

// Entity holds the schema definition for the Entity(Stock) entity.
type Entity struct {
	ent.Schema
}

// Fields of the StockEntity.
func (Entity) Fields() []ent.Field {
	return []ent.Field{
		field.Bool("active"),                     //from Polygon.io ticker details
		field.String("ticker").NotEmpty(),        //from Polygon.io ticker details
		field.String("name").NotEmpty(),          //from Polygon.io ticker details
		field.String("description").MaxLen(1000), //from Polygon.io ticker details
		field.Time("list_date").
			SchemaType(dateTimeSchemaType), //from Polygon.io ticker details
		field.Time("delisted").
			Nillable().Optional().
			SchemaType(dateTimeSchemaType), //from Polygon.io ticker details
	}
}

// Edges of the StockEntity.
func (Entity) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("exchanges", Exchange.Type), // stocks can be listed on multiple exchanges
		edge.To("intervals", Interval.Type), // stocks will have multiple intervals (trades, 1min, 1day, etc)
		edge.To("dividends", Dividend.Type),
		edge.To("splits", Split.Type),
		edge.To("financials", Financial.Type),
	}
}

func (Entity) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("ticker", "list_date").Unique(),
	}
}
