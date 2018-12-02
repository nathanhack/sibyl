package core

import (
	"database/sql"
	"fmt"
)

func nullFloat64ToString(v sql.NullFloat64, nullString string) string {
	if v.Valid {
		return fmt.Sprintf("%v", v.Float64)
	}
	return nullString
}

func nullInt64ToString(v sql.NullInt64, nullString string) string {
	if v.Valid {
		return fmt.Sprintf("%v", v.Int64)
	}
	return nullString
}

func nullStringToString(v sql.NullString, nullString string, stringEscapes bool) string {
	if v.Valid {
		if stringEscapes {
			return fmt.Sprintf("'%v'", v.String)
		} else {
			return fmt.Sprintf("%v", v.String)
		}
	}
	return nullString
}
