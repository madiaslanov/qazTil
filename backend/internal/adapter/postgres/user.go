package postgres

import (
	"context"
	"database/sql"

	"github.com/madiaslanov/qazTil/internal/domain"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user domain.User) (domain.User, error) {
	created, err := r.get(ctx, `
		INSERT INTO users (email, password_hash, display_name, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, password_hash, display_name, created_at`,
		user.Email, user.PasswordHash, user.DisplayName, user.CreatedAt.UTC())
	if err != nil {
		// Two parallel registrations can both pass the "is the email free" check,
		// so the UNIQUE index is the real guard and this maps it to a domain error.
		if isUniqueViolation(err) {
			return domain.User{}, domain.ErrEmailTaken
		}
		return domain.User{}, err
	}
	return created, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	return r.get(ctx, `
		SELECT id, email, password_hash, display_name, created_at
		FROM users
		WHERE email = $1`, email)
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (domain.User, error) {
	return r.get(ctx, `
		SELECT id, email, password_hash, display_name, created_at
		FROM users
		WHERE id = $1`, id)
}

func (r *UserRepo) get(ctx context.Context, query string, args ...any) (domain.User, error) {
	var user domain.User
	err := r.db.QueryRowContext(ctx, query, args...).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.DisplayName, &user.CreatedAt)
	if err != nil {
		return domain.User{}, mapErr(err)
	}
	user.CreatedAt = user.CreatedAt.UTC()
	return user, nil
}
