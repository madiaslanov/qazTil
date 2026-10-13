package postgres

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/madiaslanov/qazTil/internal/domain"
)

type QuizRepo struct {
	db *sql.DB
}

func NewQuizRepo(db *sql.DB) *QuizRepo {
	return &QuizRepo{db: db}
}

func (r *QuizRepo) Create(ctx context.Context, quiz domain.Quiz) (domain.Quiz, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Quiz{}, err
	}
	defer rollback(tx)

	id, err := insertID(ctx, tx, `
		INSERT INTO quizzes (category_id, created_at, finished_at)
		VALUES ($1, $2, $3)`, quiz.CategoryID, quiz.CreatedAt.UTC(), quiz.FinishedAt)
	if err != nil {
		return domain.Quiz{}, err
	}
	quiz.ID = id
	for i := range quiz.Questions {
		raw, err := json.Marshal(quiz.Questions[i].Options)
		if err != nil {
			return domain.Quiz{}, err
		}
		qid, err := insertID(ctx, tx, `
			INSERT INTO quiz_questions (quiz_id, word_id, prompt, options_json, correct_index, selected_index, correct)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			quiz.ID,
			quiz.Questions[i].WordID,
			quiz.Questions[i].Prompt,
			string(raw),
			quiz.Questions[i].CorrectIndex,
			nullInt(quiz.Questions[i].SelectedIndex),
			nullBool(quiz.Questions[i].Correct),
		)
		if err != nil {
			return domain.Quiz{}, err
		}
		quiz.Questions[i].ID = qid
	}
	if err := tx.Commit(); err != nil {
		return domain.Quiz{}, err
	}
	return quiz, nil
}

func (r *QuizRepo) Get(ctx context.Context, id int64) (domain.Quiz, error) {
	var (
		quiz     domain.Quiz
		finished sql.NullTime
	)
	err := r.db.QueryRowContext(ctx, `
		SELECT id, category_id, created_at, finished_at
		FROM quizzes
		WHERE id = $1`, id).Scan(&quiz.ID, &quiz.CategoryID, &quiz.CreatedAt, &finished)
	if err != nil {
		return domain.Quiz{}, mapErr(err)
	}
	quiz.CreatedAt = quiz.CreatedAt.UTC()
	quiz.FinishedAt = timePtr(finished)

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, word_id, prompt, options_json, correct_index, selected_index, correct
		FROM quiz_questions
		WHERE quiz_id = $1
		ORDER BY id`, id)
	if err != nil {
		return domain.Quiz{}, err
	}
	defer rows.Close()

	quiz.Questions = []domain.QuizQuestion{}
	for rows.Next() {
		var (
			question domain.QuizQuestion
			raw      string
			selected sql.NullInt64
			correct  sql.NullBool
		)
		if err := rows.Scan(&question.ID, &question.WordID, &question.Prompt, &raw, &question.CorrectIndex, &selected, &correct); err != nil {
			return domain.Quiz{}, err
		}
		if err := json.Unmarshal([]byte(raw), &question.Options); err != nil {
			return domain.Quiz{}, err
		}
		if selected.Valid {
			value := int(selected.Int64)
			question.SelectedIndex = &value
		}
		if correct.Valid {
			value := correct.Bool
			question.Correct = &value
		}
		quiz.Questions = append(quiz.Questions, question)
	}
	return quiz, rows.Err()
}

func (r *QuizRepo) Save(ctx context.Context, quiz domain.Quiz) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer rollback(tx)

	res, err := tx.ExecContext(ctx, `
		UPDATE quizzes
		SET finished_at = $1
		WHERE id = $2`, quiz.FinishedAt, quiz.ID)
	if err != nil {
		return err
	}
	if err := missingRow(res); err != nil {
		return err
	}
	for _, question := range quiz.Questions {
		if _, err := tx.ExecContext(ctx, `
			UPDATE quiz_questions
			SET selected_index = $1, correct = $2
			WHERE id = $3 AND quiz_id = $4`,
			nullInt(question.SelectedIndex), nullBool(question.Correct), question.ID, quiz.ID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func rollback(tx *sql.Tx) {
	_ = tx.Rollback()
}
