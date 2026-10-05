package sqlite

import (
	"context"
	"database/sql"

	"github.com/madiaslanov/qazTil/internal/domain"
)

type SituationRepo struct {
	db *sql.DB
}

func NewSituationRepo(db *sql.DB) *SituationRepo {
	return &SituationRepo{db: db}
}

func (r *SituationRepo) List(ctx context.Context) ([]domain.Situation, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT s.id, s.slug, s.title_kk, s.title_ru, s.description, s.position,
			(SELECT COUNT(*) FROM situation_words w WHERE w.situation_id = s.id)
		FROM situations s
		ORDER BY s.position, s.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []domain.Situation{}
	for rows.Next() {
		item, err := scanSituation(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *SituationRepo) Get(ctx context.Context, id int64) (domain.Situation, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT s.id, s.slug, s.title_kk, s.title_ru, s.description, s.position,
			(SELECT COUNT(*) FROM situation_words w WHERE w.situation_id = s.id)
		FROM situations s
		WHERE s.id = ?`, id)
	item, err := scanSituation(row)
	if err != nil {
		return domain.Situation{}, mapErr(err)
	}
	return item, nil
}

func (r *SituationRepo) ListWords(ctx context.Context, situationID int64) ([]domain.SituationWord, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, situation_id, kazakh, russian, transcription, position
		FROM situation_words
		WHERE situation_id = ?
		ORDER BY position, id`, situationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []domain.SituationWord{}
	for rows.Next() {
		var word domain.SituationWord
		if err := rows.Scan(&word.ID, &word.SituationID, &word.Kazakh, &word.Russian, &word.Transcription, &word.Position); err != nil {
			return nil, err
		}
		out = append(out, word)
	}
	return out, rows.Err()
}

type situationScanner interface {
	Scan(dest ...any) error
}

func scanSituation(row situationScanner) (domain.Situation, error) {
	var item domain.Situation
	err := row.Scan(&item.ID, &item.Slug, &item.TitleKK, &item.TitleRU, &item.Description, &item.Position, &item.WordCount)
	return item, err
}
