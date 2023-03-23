package schema

import "entgo.io/ent/dialect"

var dateTimeSchemaType = map[string]string{
	dialect.MySQL:    "datetime",
	dialect.Postgres: "date",
}
