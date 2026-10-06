package domain

import "time"

// PasswordHasher keeps the hashing algorithm out of the use cases.
type PasswordHasher interface {
	Hash(plain string) (string, error)
	Compare(hash, plain string) error
}

// TokenIssuer signs access tokens and reads the learner id back out of them.
type TokenIssuer interface {
	Issue(userID int64) (token string, expiresAt time.Time, err error)
	Parse(token string) (userID int64, err error)
}
