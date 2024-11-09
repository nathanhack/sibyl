package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Dividend holds the schema definition for the Dividend entity.
type Dividend struct {
	ent.Schema
}

// Fields of the Dividend.
func (Dividend) Fields() []ent.Field {
	return []ent.Field{
		field.Float("rate"), 
		field.Time("declaration_date").
			SchemaType(dateTimeSchemaType), 
		field.Time("ex_dividend_date").
			SchemaType(dateTimeSchemaType), 
		field.Time("record_date").
			SchemaType(dateTimeSchemaType), 
		field.Time("pay_date").
			SchemaType(dateTimeSchemaType), 
	}
}

// Edges of the Dividend.
func (Dividend) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("stock", Entity.Type).
			Ref("dividends"),
	}
}
