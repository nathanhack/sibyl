package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

//TODO move these notes to where they belong
// polygon.io  Exchanges:operating_mic  == Tickers:primary_exchange

// Exchange holds the schema definition for the Exchange entity.
type Exchange struct {
	ent.Schema
}

// Fields of the StockExchange.
func (Exchange) Fields() []ent.Field {
	return []ent.Field{
		field.String("code"), //alpaca:exchange code  == polygon.io:participant_id
		field.String("name"), //from alpaca's "Name of Exchange" https://alpaca.markets/docs/market-data/
	}
}

// Edges of the StockExchange.
func (Exchange) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("stocks", Entity.Type).Ref("exchanges"),
	}
}
