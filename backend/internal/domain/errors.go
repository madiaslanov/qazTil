package domain

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrInvalid         = errors.New("invalid")
	ErrNotEnoughWords  = errors.New("not enough words")
	ErrAlreadyAnswered = errors.New("already answered")
	ErrEmailTaken      = errors.New("email taken")
	ErrBadCredentials  = errors.New("bad credentials")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrNoLives         = errors.New("no lives")
	// ErrLocked is a lesson the learner has not reached yet.
	ErrLocked = errors.New("locked")
)
