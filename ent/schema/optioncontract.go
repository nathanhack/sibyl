package schema

import "entgo.io/ent"

//go generate ./...
//go run -mod=mod entgo.io/ent/cmd/ent init Entity

// Entity holds the schema definition for the Entity(Stock) entity.
type OptionContract struct {
	ent.Schema
}


/ Fields of the StockEntity.
func (Entity) Fields() []ent.Field {
	return []ent.Field{
		field.Bool("active"),
		field.String("ticker").NotEmpty(),
		field.String("name").NotEmpty(),
		field.String("description").
			Default("").
			MaxLen(1000), //from Polygon.io ticker details
		field.Time("list_date"). //from Polygon.io ticker details
						SchemaType(dateTimeSchemaType).
						Default(time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local)),
		field.Bool("options"),
		field.Bool("tradable"),
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