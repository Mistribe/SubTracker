package persistence

import (
	"database/sql"
)

func stringToSqlNull(input *string) sql.NullString {
	if input == nil {
		return sql.NullString{
			Valid: false,
		}
	}
	return sql.NullString{
		String: *input,
		Valid:  true,
	}
}

func sqlNullToString(input sql.NullString) *string {
	if !input.Valid {
		return nil
	}
	return &input.String
}
