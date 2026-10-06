package sqlite

import (
	"context"
	"database/sql"
	"strings"

	"github.com/madiaslanov/qazTil/internal/domain"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user domain.User) (domain.User, error) {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO users (email, password_hash, display_name, created_at)
		VALUES (?, ?, ?, ?)`,
		user.Email, user.PasswordHash, user.DisplayName, formatTime(user.CreatedAt))
	if err != nil {
		// Two parallel registrations can both pass the "is the email free" check,
		// so the UNIQUE index is the real guard and this maps it to a domain error.
		if isUniqueViolation(err) {
			return domain.User{}, domain.ErrEmailTaken
		}
		return domain.User{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.User{}, err
	}
	return r.GetByID(ctx, id)
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	return r.get(ctx, "email = ?", email)
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (domain.User, error) {
	return r.get(ctx, "id = ?", id)
}

func (r *UserRepo) get(ctx context.Context, where string, arg any) (domain.User, error) {
	var (
		user    domain.User
		created sql.NullString
	)
	err := r.db.QueryRowContext(ctx, `
		SELECT id, email, password_hash, display_name, created_at
		FROM users
		WHERE `+where, arg).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.DisplayName, &created)
	if err != nil {
		return domain.User{}, mapErr(err)
	}
	parsed, err := parseTime(created)
	if err != nil {
		return domain.User{}, err
	}
	if parsed != nil {
		user.CreatedAt = *parsed
	}
	return user, nil
}

// isUniqueViolation keeps the driver's wording in the adapter layer.
func isUniqueViolation(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "unique constraint failed")
}
