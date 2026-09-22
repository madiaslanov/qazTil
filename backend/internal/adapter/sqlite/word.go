package sqlite

import (
	"context"
	"database/sql"

	"github.com/madiaslanov/qazTil/internal/domain"
)

type WordRepo struct {
	db *sql.DB
}

func NewWordRepo(db *sql.DB) *WordRepo {
	return &WordRepo{db: db}
}

func (r *WordRepo) List(ctx context.Context, categoryID *int64) ([]domain.Word, error) {
	query := `
		SELECT id, category_id, kazakh, russian, transcription, example_kk, example_ru
		FROM words`
	var args []any
	if categoryID != nil {
		query += ` WHERE category_id = ?`
		args = append(args, *categoryID)
	}
	query += ` ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []domain.Word{}
	for rows.Next() {
		word, err := scanWord(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, word)
	}
	return out, rows.Err()
}

func (r *WordRepo) Get(ctx context.Context, id int64) (domain.Word, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, category_id, kazakh, russian, transcription, example_kk, example_ru
		FROM words
		WHERE id = ?`, id)
	word, err := scanWord(row)
	if err != nil {
		return domain.Word{}, mapErr(err)
	}
	return word, nil
}

func (r *WordRepo) Create(ctx context.Context, word domain.Word) (domain.Word, error) {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO words (category_id, kazakh, russian, transcription, example_kk, example_ru)
		VALUES (?, ?, ?, ?, ?, ?)`,
		word.CategoryID, word.Kazakh, word.Russian, word.Transcription, word.ExampleKK, word.ExampleRU)
	if err != nil {
		return domain.Word{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.Word{}, err
	}
	return r.Get(ctx, id)
}

func (r *WordRepo) Update(ctx context.Context, word domain.Word) (domain.Word, error) {
	res, err := r.db.ExecContext(ctx, `
		UPDATE words
		SET category_id = ?, kazakh = ?, russian = ?, transcription = ?, example_kk = ?, example_ru = ?
		WHERE id = ?`,
		word.CategoryID, word.Kazakh, word.Russian, word.Transcription, word.ExampleKK, word.ExampleRU, word.ID)
	if err != nil {
		return domain.Word{}, err
	}
	if err := missingRow(res); err != nil {
		return domain.Word{}, err
	}
	return r.Get(ctx, word.ID)
}

func (r *WordRepo) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM words WHERE id = ?`, id)
	if err != nil {
		return err
	}
	return missingRow(res)
}

type scanner interface {
	Scan(dest ...any) error
}

func scanWord(row scanner) (domain.Word, error) {
	var word domain.Word
	err := row.Scan(&word.ID, &word.CategoryID, &word.Kazakh, &word.Russian, &word.Transcription, &word.ExampleKK, &word.ExampleRU)
	return word, err
}
