package usecase

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/madiaslanov/qazTil/internal/domain"
)

type memUsers struct {
	items  map[int64]domain.User
	nextID int64
}

func newMemUsers() *memUsers {
	return &memUsers{items: map[int64]domain.User{}, nextID: 1}
}

func (m *memUsers) Create(_ context.Context, user domain.User) (domain.User, error) {
	for _, existing := range m.items {
		if existing.Email == user.Email {
			return domain.User{}, domain.ErrEmailTaken
		}
	}
	user.ID = m.nextID
	m.nextID++
	m.items[user.ID] = user
	return user, nil
}

func (m *memUsers) GetByEmail(_ context.Context, email string) (domain.User, error) {
	for _, user := range m.items {
		if user.Email == email {
			return user, nil
		}
	}
	return domain.User{}, domain.ErrNotFound
}

func (m *memUsers) GetByID(_ context.Context, id int64) (domain.User, error) {
	user, ok := m.items[id]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	return user, nil
}

// stubHasher stands in for bcrypt so the tests stay fast.
type stubHasher struct{}

func (stubHasher) Hash(plain string) (string, error) { return "hashed:" + plain, nil }

func (stubHasher) Compare(hash, plain string) error {
	if hash != "hashed:"+plain {
		return domain.ErrBadCredentials
	}
	return nil
}

type stubTokens struct{ issued int }

func (s *stubTokens) Issue(userID int64) (string, time.Time, error) {
	s.issued++
	return "token", time.Unix(0, 0).UTC().Add(time.Hour), nil
}

func (s *stubTokens) Parse(string) (int64, error) { return 0, domain.ErrUnauthorized }

func newAuth() (*AuthService, *memUsers) {
	users := newMemUsers()
	return NewAuthService(users, stubHasher{}, &stubTokens{}), users
}

func TestRegisterStoresHashAndIssuesToken(t *testing.T) {
	auth, users := newAuth()

	session, err := auth.Register(context.Background(), "  Aidana@Example.COM ", "qazaqtili2026", "Айдана")
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if session.Token == "" {
		t.Fatal("want a token")
	}
	if session.User.Email != "aidana@example.com" {
		t.Fatalf("email not normalized: %q", session.User.Email)
	}
	if session.User.PasswordHash != "" {
		t.Fatal("password hash must not leave the service")
	}

	stored, err := users.GetByEmail(context.Background(), "aidana@example.com")
	if err != nil {
		t.Fatalf("get stored: %v", err)
	}
	if strings.Contains(stored.PasswordHash, "qazaqtili2026") == false {
		t.Fatalf("unexpected stored hash %q", stored.PasswordHash)
	}
	if stored.PasswordHash == "qazaqtili2026" {
		t.Fatal("password stored in plain text")
	}
}

func TestRegisterFallsBackToEmailName(t *testing.T) {
	auth, _ := newAuth()

	session, err := auth.Register(context.Background(), "erasyl@example.com", "qazaqtili2026", "   ")
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if session.User.DisplayName != "erasyl" {
		t.Fatalf("display name = %q", session.User.DisplayName)
	}
}

func TestRegisterRejectsDuplicateEmail(t *testing.T) {
	auth, _ := newAuth()
	ctx := context.Background()

	if _, err := auth.Register(ctx, "madi@example.com", "qazaqtili2026", ""); err != nil {
		t.Fatalf("first register: %v", err)
	}
	// Same address, different case: normalization must still catch it.
	_, err := auth.Register(ctx, "MADI@example.com", "qazaqtili2026", "")
	if !errors.Is(err, domain.ErrEmailTaken) {
		t.Fatalf("want ErrEmailTaken, got %v", err)
	}
}

func TestRegisterRejectsBadInput(t *testing.T) {
	cases := map[string]struct{ email, password string }{
		"short password": {"shokan@example.com", "qazaq"},
		"no at sign":     {"shokan.example.com", "qazaqtili2026"},
		"no domain dot":  {"shokan@example", "qazaqtili2026"},
		"empty email":    {"", "qazaqtili2026"},
		"long password":  {"shokan@example.com", strings.Repeat("a", 73)},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			auth, _ := newAuth()
			if _, err := auth.Register(context.Background(), tc.email, tc.password, ""); !errors.Is(err, domain.ErrInvalid) {
				t.Fatalf("want ErrInvalid, got %v", err)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	auth, _ := newAuth()
	ctx := context.Background()
	if _, err := auth.Register(ctx, "maria@example.com", "qazaqtili2026", "Мария"); err != nil {
		t.Fatalf("register: %v", err)
	}

	session, err := auth.Login(ctx, "Maria@Example.com", "qazaqtili2026")
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if session.Token == "" || session.User.ID == 0 {
		t.Fatal("want a signed-in session")
	}
}

func TestLoginRejectsWrongPasswordAndUnknownEmailAlike(t *testing.T) {
	auth, _ := newAuth()
	ctx := context.Background()
	if _, err := auth.Register(ctx, "maria@example.com", "qazaqtili2026", ""); err != nil {
		t.Fatalf("register: %v", err)
	}

	if _, err := auth.Login(ctx, "maria@example.com", "wrongpassword"); !errors.Is(err, domain.ErrBadCredentials) {
		t.Fatalf("wrong password: want ErrBadCredentials, got %v", err)
	}
	if _, err := auth.Login(ctx, "nobody@example.com", "qazaqtili2026"); !errors.Is(err, domain.ErrBadCredentials) {
		t.Fatalf("unknown email: want ErrBadCredentials, got %v", err)
	}
	// A malformed address must not leak a different error than a wrong password.
	if _, err := auth.Login(ctx, "broken", "qazaqtili2026"); !errors.Is(err, domain.ErrBadCredentials) {
		t.Fatalf("malformed email: want ErrBadCredentials, got %v", err)
	}
}

func TestMeHidesPasswordHash(t *testing.T) {
	auth, _ := newAuth()
	ctx := context.Background()
	session, err := auth.Register(ctx, "birzhan@example.com", "qazaqtili2026", "")
	if err != nil {
		t.Fatalf("register: %v", err)
	}

	user, err := auth.Me(ctx, session.User.ID)
	if err != nil {
		t.Fatalf("me: %v", err)
	}
	if user.PasswordHash != "" {
		t.Fatal("password hash must not be returned")
	}
	if _, err := auth.Me(ctx, 0); !errors.Is(err, domain.ErrInvalid) {
		t.Fatalf("want ErrInvalid, got %v", err)
	}
	if _, err := auth.Me(ctx, 404); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}
