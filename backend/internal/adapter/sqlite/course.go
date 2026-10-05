package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/madiaslanov/qazTil/internal/domain"
)

type CourseRepo struct {
	db *sql.DB
}

func NewCourseRepo(db *sql.DB) *CourseRepo {
	return &CourseRepo{db: db}
}

func (r *CourseRepo) ListUnits(ctx context.Context) ([]domain.Unit, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT u.id, u.situation_id, u.title_kk, u.title_ru, u.description, u.position,
			(SELECT COUNT(*) FROM lessons l WHERE l.unit_id = u.id)
		FROM units u
		ORDER BY u.position, u.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []domain.Unit{}
	for rows.Next() {
		item, err := scanUnit(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *CourseRepo) GetUnit(ctx context.Context, id int64) (domain.Unit, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT u.id, u.situation_id, u.title_kk, u.title_ru, u.description, u.position,
			(SELECT COUNT(*) FROM lessons l WHERE l.unit_id = u.id)
		FROM units u
		WHERE u.id = ?`, id)
	item, err := scanUnit(row)
	if err != nil {
		return domain.Unit{}, mapErr(err)
	}
	return item, nil
}

func (r *CourseRepo) ListLessons(ctx context.Context, unitID int64) ([]domain.Lesson, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, unit_id, title_kk, title_ru, description, position
		FROM lessons
		WHERE unit_id = ?
		ORDER BY position, id`, unitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []domain.Lesson{}
	for rows.Next() {
		item, err := scanLesson(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *CourseRepo) GetLesson(ctx context.Context, id int64) (domain.Lesson, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, unit_id, title_kk, title_ru, description, position
		FROM lessons
		WHERE id = ?`, id)
	item, err := scanLesson(row)
	if err != nil {
		return domain.Lesson{}, mapErr(err)
	}
	return item, nil
}

func (r *CourseRepo) ListExercises(ctx context.Context, lessonID int64) ([]domain.Exercise, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, lesson_id, kind, prompt, answer, options_json, position
		FROM exercises
		WHERE lesson_id = ?
		ORDER BY position, id`, lessonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []domain.Exercise{}
	for rows.Next() {
		item, err := scanExercise(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *CourseRepo) GetExercise(ctx context.Context, id int64) (domain.Exercise, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, lesson_id, kind, prompt, answer, options_json, position
		FROM exercises
		WHERE id = ?`, id)
	item, err := scanExercise(row)
	if err != nil {
		return domain.Exercise{}, mapErr(err)
	}
	return item, nil
}

func scanUnit(row situationScanner) (domain.Unit, error) {
	var item domain.Unit
	err := row.Scan(&item.ID, &item.SituationID, &item.TitleKK, &item.TitleRU, &item.Description, &item.Position, &item.LessonCount)
	return item, err
}

func scanLesson(row situationScanner) (domain.Lesson, error) {
	var item domain.Lesson
	err := row.Scan(&item.ID, &item.UnitID, &item.TitleKK, &item.TitleRU, &item.Description, &item.Position)
	return item, err
}

func scanExercise(row situationScanner) (domain.Exercise, error) {
	var (
		item domain.Exercise
		raw  string
	)
	if err := row.Scan(&item.ID, &item.LessonID, &item.Kind, &item.Prompt, &item.Answer, &raw, &item.Position); err != nil {
		return domain.Exercise{}, err
	}
	options, err := decodeOptions(raw)
	if err != nil {
		return domain.Exercise{}, err
	}
	item.Options = options
	return item, nil
}

func decodeOptions(raw string) ([]string, error) {
	if raw == "" {
		return []string{}, nil
	}
	var options []string
	if err := json.Unmarshal([]byte(raw), &options); err != nil {
		return nil, err
	}
	if options == nil {
		return []string{}, nil
	}
	return options, nil
}
