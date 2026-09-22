package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

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

	res, err := tx.ExecContext(ctx, `
		INSERT INTO quizzes (category_id, created_at, finished_at)
		VALUES (?, ?, ?)`, quiz.CategoryID, formatTime(quiz.CreatedAt), formatTimePtr(quiz.FinishedAt))
	if err != nil {
		return domain.Quiz{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.Quiz{}, err
	}
	quiz.ID = id
	for i := range quiz.Questions {
		raw, err := json.Marshal(quiz.Questions[i].Options)
		if err != nil {
			return domain.Quiz{}, err
		}
		res, err = tx.ExecContext(ctx, `
			INSERT INTO quiz_questions (quiz_id, word_id, prompt, options_json, correct_index, selected_index, correct)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
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
		qid, err := res.LastInsertId()
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
		created  string
		finished sql.NullString
	)
	err := r.db.QueryRowContext(ctx, `
		SELECT id, category_id, created_at, finished_at
		FROM quizzes
		WHERE id = ?`, id).Scan(&quiz.ID, &quiz.CategoryID, &created, &finished)
	if err != nil {
		return domain.Quiz{}, mapErr(err)
	}
	quiz.CreatedAt, err = time.Parse(time.RFC3339Nano, created)
	if err != nil {
		return domain.Quiz{}, err
	}
	quiz.FinishedAt, err = parseTime(finished)
	if err != nil {
		return domain.Quiz{}, err
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, word_id, prompt, options_json, correct_index, selected_index, correct
		FROM quiz_questions
		WHERE quiz_id = ?
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
			correct  sql.NullInt64
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
			value := correct.Int64 != 0
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
		SET finished_at = ?
		WHERE id = ?`, formatTimePtr(quiz.FinishedAt), quiz.ID)
	if err != nil {
		return err
	}
	if err := missingRow(res); err != nil {
		return err
	}
	for _, question := range quiz.Questions {
		if _, err := tx.ExecContext(ctx, `
			UPDATE quiz_questions
			SET selected_index = ?, correct = ?
			WHERE id = ? AND quiz_id = ?`,
			nullInt(question.SelectedIndex), nullBool(question.Correct), question.ID, quiz.ID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func rollback(tx *sql.Tx) {
	_ = tx.Rollback()
}
