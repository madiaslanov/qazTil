package security

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/madiaslanov/qazTil/internal/domain"
)

// BcryptHasher hashes passwords with bcrypt. The salt is generated per password
// and stored inside the hash string, so the users table needs no salt column.
type BcryptHasher struct {
	cost int
}

// NewBcryptHasher falls back to the library default when cost is out of range.
func NewBcryptHasher(cost int) *BcryptHasher {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	return &BcryptHasher{cost: cost}
}

func (h *BcryptHasher) Hash(plain string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), h.cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (h *BcryptHasher) Compare(hash, plain string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)); err != nil {
		return domain.ErrBadCredentials
	}
	return nil
}
