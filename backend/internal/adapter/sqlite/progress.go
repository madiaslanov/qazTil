package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/madiaslanov/qazTil/internal/domain"
)

type ProgressRepo struct {
	db *sql.DB
}

func NewProgressRepo(db *sql.DB) *ProgressRepo {
	return &ProgressRepo{db: db}
}

func (r *ProgressRepo) List(ctx context.Context) ([]domain.Progress, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT category_id, total_answers, correct_answers, last_studied_at
		FROM progress
		ORDER BY category_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []domain.Progress{}
	for rows.Next() {
		var (
			item    domain.Progress
			studied sql.NullString
		)
		if err := rows.Scan(&item.CategoryID, &item.TotalAnswers, &item.CorrectAnswers, &studied); err != nil {
			return nil, err
		}
		item.LastStudiedAt, err = parseTime(studied)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *ProgressRepo) Record(ctx context.Context, categoryID int64, correct bool, at time.Time) error {
	increment := 0
	if correct {
		increment = 1
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO progress (category_id, total_answers, correct_answers, last_studied_at)
		VALUES (?, 1, ?, ?)
		ON CONFLICT(category_id) DO UPDATE SET
			total_answers = total_answers + 1,
			correct_answers = correct_answers + excluded.correct_answers,
			last_studied_at = excluded.last_studied_at`,
		categoryID, increment, formatTime(at))
	return err
}
