package domain

import "time"

// User is a registered learner.
type User struct {
	ID           int64
	Email        string
	PasswordHash string
	DisplayName  string
	CreatedAt    time.Time
}
