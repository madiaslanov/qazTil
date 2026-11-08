package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/madiaslanov/qazTil/internal/domain"
)

type LearnerRepo struct {
	db *sql.DB
}

func NewLearnerRepo(db *sql.DB) *LearnerRepo {
	return &LearnerRepo{db: db}
}

func (r *LearnerRepo) Current(ctx context.Context, userID int64, now time.Time) (domain.LearnerStats, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.LearnerStats{}, err
	}
	defer rollback(tx)

	stats, err := lockStats(ctx, tx, userID)
	if err != nil {
		return domain.LearnerStats{}, err
	}
	stats = stats.Regenerated(now)
	if err := saveStats(ctx, tx, stats); err != nil {
		return domain.LearnerStats{}, err
	}
	if err := tx.Commit(); err != nil {
		return domain.LearnerStats{}, err
	}
	return stats, nil
}

func (r *LearnerRepo) Apply(ctx context.Context, userID, exerciseID, lessonID int64, correct bool, now time.Time) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer rollback(tx)

	stats, err := lockStats(ctx, tx, userID)
	if err != nil {
		return err
	}
	stats = stats.Regenerated(now)
	if stats.Lives == 0 {
		if err := saveStats(ctx, tx, stats); err != nil {
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
		return domain.ErrNoLives
	}

	var prior int
	if err := tx.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM attempts
		WHERE user_id = $1 AND exercise_id = $2`, userID, exerciseID).Scan(&prior); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO attempts (user_id, exercise_id, correct, created_at)
		VALUES ($1, $2, $3, $4)`, userID, exerciseID, correct, now.UTC()); err != nil {
		return err
	}
	if prior == 0 {
		if correct {
			stats = stats.GainXP()
		} else {
			stats = stats.LoseLife(now)
		}
	}
	finished, err := lessonJustFinished(ctx, tx, userID, lessonID, prior == 0)
	if err != nil {
		return err
	}
	if finished {
		stats = stats.CompleteLesson(now)
	}
	if err := saveStats(ctx, tx, stats); err != nil {
		return err
	}
	return tx.Commit()
}

func lockStats(ctx context.Context, tx *sql.Tx, userID int64) (domain.LearnerStats, error) {
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO learner_stats (user_id, xp, lives, streak)
		VALUES ($1, 0, $2, 0)
		ON CONFLICT (user_id) DO NOTHING`, userID, domain.MaxLives); err != nil {
		return domain.LearnerStats{}, err
	}
	var (
		stats domain.LearnerStats
		next  sql.NullTime
		last  sql.NullTime
	)
	err := tx.QueryRowContext(ctx, `
		SELECT user_id, xp, lives, next_life_at, streak, last_lesson_on
		FROM learner_stats
		WHERE user_id = $1
		FOR UPDATE`, userID).Scan(&stats.UserID, &stats.XP, &stats.Lives, &next, &stats.Streak, &last)
	if err != nil {
		return domain.LearnerStats{}, err
	}
	stats.MaxLives = domain.MaxLives
	stats.NextLifeAt = timePtr(next)
	if last.Valid {
		day := time.Date(last.Time.UTC().Year(), last.Time.UTC().Month(), last.Time.UTC().Day(), 0, 0, 0, 0, time.UTC)
		stats.LastLessonOn = &day
	}
	return stats, nil
}

func saveStats(ctx context.Context, tx *sql.Tx, stats domain.LearnerStats) error {
	_, err := tx.ExecContext(ctx, `
		UPDATE learner_stats
		SET xp = $2, lives = $3, next_life_at = $4, streak = $5, last_lesson_on = $6
		WHERE user_id = $1`,
		stats.UserID, stats.XP, stats.Lives, stats.NextLifeAt, stats.Streak, dateOrNil(stats.LastLessonOn))
	return err
}

func dateOrNil(value *time.Time) any {
	if value == nil {
		return nil
	}
	return value.UTC().Format("2006-01-02")
}

func lessonJustFinished(ctx context.Context, tx *sql.Tx, userID, lessonID int64, firstOnExercise bool) (bool, error) {
	if !firstOnExercise {
		return false, nil
	}
	var total, attempted int
	if err := tx.QueryRowContext(ctx, `
		SELECT COUNT(*),
			COUNT(*) FILTER (
				WHERE EXISTS (
					SELECT 1 FROM attempts a
					WHERE a.exercise_id = exercises.id AND a.user_id = $2
				)
			)
		FROM exercises
		WHERE lesson_id = $1`, lessonID, userID).Scan(&total, &attempted); err != nil {
		return false, err
	}
	return total > 0 && attempted == total, nil
}
