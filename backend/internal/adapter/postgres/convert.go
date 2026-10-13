package postgres

import (
	"database/sql"
	"time"
)

func timePtr(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}
	parsed := value.Time.UTC()
	return &parsed
}

func nullInt(value *int) any {
	if value == nil {
		return nil
	}
	return *value
}

func nullBool(value *bool) any {
	if value == nil {
		return nil
	}
	return *value
}
