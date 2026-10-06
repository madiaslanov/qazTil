package usecase

import (
	"context"
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/madiaslanov/qazTil/internal/domain"
)

const (
	minPasswordLen = 8
	// bcrypt silently ignores everything past 72 bytes, so refuse longer passwords
	// instead of accepting a password that is only partly checked on login.
	maxPasswordLen     = 72
	maxEmailLen        = 254
	maxDisplayNameRune = 60
)

// Session is a signed token together with the learner it belongs to.
type Session struct {
	User      domain.User
	Token     string
	ExpiresAt time.Time
}

// AuthService registers learners and issues access tokens.
type AuthService struct {
	users  domain.UserRepository
	hasher domain.PasswordHasher
	tokens domain.TokenIssuer
	now    func() time.Time
}

func NewAuthService(users domain.UserRepository, hasher domain.PasswordHasher, tokens domain.TokenIssuer) *AuthService {
	return &AuthService{users: users, hasher: hasher, tokens: tokens, now: time.Now}
}

// Register creates a learner and signs them in straight away.
func (s *AuthService) Register(ctx context.Context, email, password, displayName string) (Session, error) {
	email = normalizeEmail(email)
	if err := validateCredentials(email, password); err != nil {
		return Session{}, err
	}
	displayName = strings.TrimSpace(displayName)
	if displayName == "" {
		displayName = nameFromEmail(email)
	}
	if utf8.RuneCountInString(displayName) > maxDisplayNameRune {
		return Session{}, domain.ErrInvalid
	}

	switch _, err := s.users.GetByEmail(ctx, email); {
	case err == nil:
		return Session{}, domain.ErrEmailTaken
	case !errors.Is(err, domain.ErrNotFound):
		return Session{}, err
	}

	hash, err := s.hasher.Hash(password)
	if err != nil {
		return Session{}, err
	}
	user, err := s.users.Create(ctx, domain.User{
		Email:        email,
		PasswordHash: hash,
		DisplayName:  displayName,
		CreatedAt:    s.now().UTC(),
	})
	if err != nil {
		return Session{}, err
	}
	return s.session(user)
}

// Login checks the password and issues a fresh token.
// Every failure returns ErrBadCredentials so the response never reveals
// whether the email is registered.
func (s *AuthService) Login(ctx context.Context, email, password string) (Session, error) {
	email = normalizeEmail(email)
	if err := validateCredentials(email, password); err != nil {
		return Session{}, domain.ErrBadCredentials
	}
	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return Session{}, domain.ErrBadCredentials
		}
		return Session{}, err
	}
	if err := s.hasher.Compare(user.PasswordHash, password); err != nil {
		return Session{}, domain.ErrBadCredentials
	}
	return s.session(user)
}

// Me returns the learner behind a verified token.
func (s *AuthService) Me(ctx context.Context, userID int64) (domain.User, error) {
	if userID <= 0 {
		return domain.User{}, domain.ErrInvalid
	}
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return domain.User{}, err
	}
	user.PasswordHash = ""
	return user, nil
}

func (s *AuthService) session(user domain.User) (Session, error) {
	token, expiresAt, err := s.tokens.Issue(user.ID)
	if err != nil {
		return Session{}, err
	}
	// The hash never leaves the service.
	user.PasswordHash = ""
	return Session{User: user, Token: token, ExpiresAt: expiresAt}, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func nameFromEmail(email string) string {
	if at := strings.IndexByte(email, '@'); at > 0 {
		return email[:at]
	}
	return email
}

func validateCredentials(email, password string) error {
	if !validEmail(email) {
		return domain.ErrInvalid
	}
	if len(password) < minPasswordLen || len(password) > maxPasswordLen {
		return domain.ErrInvalid
	}
	return nil
}

// validEmail is deliberately small: one @, something on both sides, a dot in the domain.
// Anything stricter belongs to a confirmation email, not to a regexp.
func validEmail(email string) bool {
	if email == "" || len(email) > maxEmailLen {
		return false
	}
	if strings.ContainsAny(email, " \t\r\n") {
		return false
	}
	at := strings.IndexByte(email, '@')
	if at <= 0 || at != strings.LastIndexByte(email, '@') {
		return false
	}
	host := email[at+1:]
	if len(host) < 3 || !strings.Contains(host, ".") {
		return false
	}
	return !strings.HasPrefix(host, ".") && !strings.HasSuffix(host, ".")
}
