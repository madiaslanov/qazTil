package postgres

import (
	"context"
	"database/sql"

	"github.com/madiaslanov/qazTil/internal/domain"
)

type CourseProgressRepo struct {
	db *sql.DB
}

func NewCourseProgressRepo(db *sql.DB) *CourseProgressRepo {
	return &CourseProgressRepo{db: db}
}

func (r *CourseProgressRepo) List(ctx context.Context, userID int64) (domain.CourseProgress, error) {
	units, err := r.units(ctx)
	if err != nil {
		return domain.CourseProgress{}, err
	}
	lessons, err := r.lessons(ctx)
	if err != nil {
		return domain.CourseProgress{}, err
	}
	exercises, marks, err := r.exercises(ctx, userID)
	if err != nil {
		return domain.CourseProgress{}, err
	}
	return domain.AssembleProgress(units, lessons, exercises, marks), nil
}

func (r *CourseProgressRepo) units(ctx context.Context) ([]domain.Unit, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id
		FROM units
		ORDER BY position, id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Unit{}
	for rows.Next() {
		var item domain.Unit
		if err := rows.Scan(&item.ID); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *CourseProgressRepo) lessons(ctx context.Context) ([]domain.Lesson, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, unit_id
		FROM lessons
		ORDER BY unit_id, position, id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Lesson{}
	for rows.Next() {
		var item domain.Lesson
		if err := rows.Scan(&item.ID, &item.UnitID); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *CourseProgressRepo) exercises(ctx context.Context, userID int64) ([]domain.Exercise, map[int64]domain.AttemptMark, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT e.id, e.lesson_id,
			EXISTS (
				SELECT 1 FROM attempts a
				WHERE a.exercise_id = e.id AND a.user_id = $1
			),
			EXISTS (
				SELECT 1 FROM attempts a
				WHERE a.exercise_id = e.id AND a.user_id = $1 AND a.correct
			)
		FROM exercises e
		ORDER BY e.lesson_id, e.position, e.id`, userID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	exercises := []domain.Exercise{}
	marks := map[int64]domain.AttemptMark{}
	for rows.Next() {
		var (
			item      domain.Exercise
			attempted bool
			correct   bool
		)
		if err := rows.Scan(&item.ID, &item.LessonID, &attempted, &correct); err != nil {
			return nil, nil, err
		}
		exercises = append(exercises, item)
		if attempted {
			marks[item.ID] = domain.AttemptMark{Correct: correct}
		}
	}
	return exercises, marks, rows.Err()
}
