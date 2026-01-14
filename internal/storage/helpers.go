package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
)

// ErrNotFound is returned when a record is not found
var ErrNotFound = errors.New("record not found")

// toNullString converts a string to sql.NullString
func toNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}

// marshalJSON marshals a value to JSON string
func marshalJSON(v interface{}) string {
	if v == nil {
		return ""
	}
	data, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(data)
}
