package domain

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrInvalid         = errors.New("invalid")
	ErrNotEnoughWords  = errors.New("not enough words")
	ErrAlreadyAnswered = errors.New("already answered")
)
