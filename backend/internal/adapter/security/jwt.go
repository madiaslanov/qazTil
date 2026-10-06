package security

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/madiaslanov/qazTil/internal/domain"
)

// DefaultTokenTTL is how long an access token stays valid when JWT_TTL is unset.
const DefaultTokenTTL = 24 * time.Hour

// JWTIssuer signs HS256 access tokens that carry the learner id in sub.
type JWTIssuer struct {
	secret []byte
	ttl    time.Duration
	now    func() time.Time
}

func NewJWTIssuer(secret string, ttl time.Duration) *JWTIssuer {
	if ttl <= 0 {
		ttl = DefaultTokenTTL
	}
	return &JWTIssuer{secret: []byte(secret), ttl: ttl, now: time.Now}
}

func (i *JWTIssuer) Issue(userID int64) (string, time.Time, error) {
	issuedAt := i.now().UTC()
	expiresAt := issuedAt.Add(i.ttl)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(userID, 10),
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	})
	signed, err := token.SignedString(i.secret)
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, expiresAt, nil
}

// Parse verifies the signature and the expiry. WithValidMethods pins HS256,
// which is what stops a token re-signed with alg=none from being accepted.
func (i *JWTIssuer) Parse(raw string) (int64, error) {
	var claims jwt.RegisteredClaims
	_, err := jwt.ParseWithClaims(raw, &claims, func(*jwt.Token) (any, error) {
		return i.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return 0, domain.ErrUnauthorized
	}
	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil || userID <= 0 {
		return 0, domain.ErrUnauthorized
	}
	return userID, nil
}
