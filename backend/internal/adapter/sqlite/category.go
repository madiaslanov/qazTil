package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/madiaslanov/qazTil/internal/domain"
)

type CategoryRepo struct {
	db *sql.DB
}

func NewCategoryRepo(db *sql.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) List(ctx context.Context) ([]domain.Category, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name_kk, name_ru, description
		FROM categories
		ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []domain.Category{}
	for rows.Next() {
		var category domain.Category
		if err := rows.Scan(&category.ID, &category.NameKK, &category.NameRU, &category.Description); err != nil {
			return nil, err
		}
		out = append(out, category)
	}
	return out, rows.Err()
}

func (r *CategoryRepo) Get(ctx context.Context, id int64) (domain.Category, error) {
	var category domain.Category
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name_kk, name_ru, description
		FROM categories
		WHERE id = ?`, id).Scan(&category.ID, &category.NameKK, &category.NameRU, &category.Description)
	if err != nil {
		return domain.Category{}, mapErr(err)
	}
	return category, nil
}

func (r *CategoryRepo) Create(ctx context.Context, category domain.Category) (domain.Category, error) {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO categories (name_kk, name_ru, description)
		VALUES (?, ?, ?)`, category.NameKK, category.NameRU, category.Description)
	if err != nil {
		return domain.Category{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.Category{}, err
	}
	return r.Get(ctx, id)
}

func (r *CategoryRepo) Update(ctx context.Context, category domain.Category) (domain.Category, error) {
	res, err := r.db.ExecContext(ctx, `
		UPDATE categories
		SET name_kk = ?, name_ru = ?, description = ?
		WHERE id = ?`, category.NameKK, category.NameRU, category.Description, category.ID)
	if err != nil {
		return domain.Category{}, err
	}
	if err := missingRow(res); err != nil {
		return domain.Category{}, err
	}
	return r.Get(ctx, category.ID)
}

func (r *CategoryRepo) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM categories WHERE id = ?`, id)
	if err != nil {
		return err
	}
	return missingRow(res)
}

func missingRow(res sql.Result) error {
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func mapErr(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrNotFound
	}
	return err
}
