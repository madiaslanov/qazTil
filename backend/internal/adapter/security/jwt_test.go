package security

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/madiaslanov/qazTil/internal/domain"
)

const testSecret = "test-secret-0123456789abcdefghijklmno"

func TestIssueThenParse(t *testing.T) {
	issuer := NewJWTIssuer(testSecret, time.Hour)

	token, expiresAt, err := issuer.Issue(42)
	if err != nil {
		t.Fatalf("issue: %v", err)
	}
	if time.Until(expiresAt) <= 0 {
		t.Fatal("token already expired")
	}

	userID, err := issuer.Parse(token)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if userID != 42 {
		t.Fatalf("user id = %d, want 42", userID)
	}
}

func TestParseRejectsOtherSecret(t *testing.T) {
	token, _, err := NewJWTIssuer(testSecret, time.Hour).Issue(42)
	if err != nil {
		t.Fatalf("issue: %v", err)
	}

	if _, err := NewJWTIssuer("another-secret-0123456789abcdefghij", time.Hour).Parse(token); !errors.Is(err, domain.ErrUnauthorized) {
		t.Fatalf("want ErrUnauthorized, got %v", err)
	}
}

func TestParseRejectsExpired(t *testing.T) {
	issuer := NewJWTIssuer(testSecret, time.Hour)
	issuer.now = func() time.Time { return time.Now().Add(-2 * time.Hour) }

	token, _, err := issuer.Issue(42)
	if err != nil {
		t.Fatalf("issue: %v", err)
	}

	fresh := NewJWTIssuer(testSecret, time.Hour)
	if _, err := fresh.Parse(token); !errors.Is(err, domain.ErrUnauthorized) {
		t.Fatalf("want ErrUnauthorized, got %v", err)
	}
}

// A token re-signed with alg=none is the classic JWT bypass; WithValidMethods stops it.
func TestParseRejectsNoneAlgorithm(t *testing.T) {
	unsigned := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.RegisteredClaims{
		Subject:   "42",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	})
	token, err := unsigned.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("sign none: %v", err)
	}

	if _, err := NewJWTIssuer(testSecret, time.Hour).Parse(token); !errors.Is(err, domain.ErrUnauthorized) {
		t.Fatalf("want ErrUnauthorized, got %v", err)
	}
}

func TestParseRejectsGarbageAndEmptySubject(t *testing.T) {
	issuer := NewJWTIssuer(testSecret, time.Hour)
	for _, raw := range []string{"", "not.a.token", "a.b.c"} {
		if _, err := issuer.Parse(raw); !errors.Is(err, domain.ErrUnauthorized) {
			t.Fatalf("%q: want ErrUnauthorized, got %v", raw, err)
		}
	}

	zero := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   "0",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	})
	token, err := zero.SignedString([]byte(testSecret))
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	if _, err := issuer.Parse(token); !errors.Is(err, domain.ErrUnauthorized) {
		t.Fatalf("subject 0: want ErrUnauthorized, got %v", err)
	}
}

func TestBcryptHasher(t *testing.T) {
	hasher := NewBcryptHasher(4) // low cost keeps the test fast

	hash, err := hasher.Hash("qazaqtili2026")
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	if hash == "qazaqtili2026" || hash == "" {
		t.Fatal("password was not hashed")
	}
	if err := hasher.Compare(hash, "qazaqtili2026"); err != nil {
		t.Fatalf("compare: %v", err)
	}
	if err := hasher.Compare(hash, "wrongpassword"); !errors.Is(err, domain.ErrBadCredentials) {
		t.Fatalf("want ErrBadCredentials, got %v", err)
	}

	// The same password must produce a different hash: bcrypt salts every call.
	other, err := hasher.Hash("qazaqtili2026")
	if err != nil {
		t.Fatalf("hash again: %v", err)
	}
	if other == hash {
		t.Fatal("identical hashes: salt is missing")
	}
}
