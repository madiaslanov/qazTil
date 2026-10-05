package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/madiaslanov/qazTil/internal/domain"
)

func TestCheckExercise(t *testing.T) {
	course := &memCourse{exercises: []domain.Exercise{
		{ID: 1, LessonID: 1, Kind: domain.ExerciseChoice, Prompt: "Привет", Answer: "Сәлем", Options: []string{"Сәлем", "Рақмет"}},
		{ID: 2, LessonID: 1, Kind: domain.ExerciseTranslate, Prompt: "Спасибо", Answer: "Рақмет"},
	}}
	svc := NewCourseService(course)

	ok, err := svc.Check(context.Background(), 1, "  сәлем ")
	if err != nil {
		t.Fatal(err)
	}
	if !ok.Correct || ok.Expected != "Сәлем" {
		t.Fatalf("check = %+v", ok)
	}

	bad, err := svc.Check(context.Background(), 2, "сәлем")
	if err != nil {
		t.Fatal(err)
	}
	if bad.Correct || bad.Expected != "Рақмет" {
		t.Fatalf("check = %+v", bad)
	}

	_, err = svc.Check(context.Background(), 2, "   ")
	if !errors.Is(err, domain.ErrInvalid) {
		t.Fatalf("err = %v", err)
	}
	_, err = svc.Check(context.Background(), 9, "сәлем")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("err = %v", err)
	}
}

func TestOpenMissingUnit(t *testing.T) {
	_, _, err := NewCourseService(&memCourse{}).OpenUnit(context.Background(), 4)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("err = %v", err)
	}
	_, _, err = NewSituationService(&memSituations{}).Open(context.Background(), 0)
	if !errors.Is(err, domain.ErrInvalid) {
		t.Fatalf("err = %v", err)
	}
}

type memSituations struct {
	items []domain.Situation
	words []domain.SituationWord
}

func (m *memSituations) List(context.Context) ([]domain.Situation, error) {
	return append([]domain.Situation(nil), m.items...), nil
}

func (m *memSituations) Get(_ context.Context, id int64) (domain.Situation, error) {
	for _, item := range m.items {
		if item.ID == id {
			return item, nil
		}
	}
	return domain.Situation{}, domain.ErrNotFound
}

func (m *memSituations) ListWords(_ context.Context, situationID int64) ([]domain.SituationWord, error) {
	out := []domain.SituationWord{}
	for _, word := range m.words {
		if word.SituationID == situationID {
			out = append(out, word)
		}
	}
	return out, nil
}

type memCourse struct {
	units     []domain.Unit
	lessons   []domain.Lesson
	exercises []domain.Exercise
}

func (m *memCourse) ListUnits(context.Context) ([]domain.Unit, error) {
	return append([]domain.Unit(nil), m.units...), nil
}

func (m *memCourse) GetUnit(_ context.Context, id int64) (domain.Unit, error) {
	for _, item := range m.units {
		if item.ID == id {
			return item, nil
		}
	}
	return domain.Unit{}, domain.ErrNotFound
}

func (m *memCourse) ListLessons(_ context.Context, unitID int64) ([]domain.Lesson, error) {
	out := []domain.Lesson{}
	for _, item := range m.lessons {
		if item.UnitID == unitID {
			out = append(out, item)
		}
	}
	return out, nil
}

func (m *memCourse) GetLesson(_ context.Context, id int64) (domain.Lesson, error) {
	for _, item := range m.lessons {
		if item.ID == id {
			return item, nil
		}
	}
	return domain.Lesson{}, domain.ErrNotFound
}

func (m *memCourse) ListExercises(_ context.Context, lessonID int64) ([]domain.Exercise, error) {
	out := []domain.Exercise{}
	for _, item := range m.exercises {
		if item.LessonID == lessonID {
			out = append(out, item)
		}
	}
	return out, nil
}

func (m *memCourse) GetExercise(_ context.Context, id int64) (domain.Exercise, error) {
	for _, item := range m.exercises {
		if item.ID == id {
			return item, nil
		}
	}
	return domain.Exercise{}, domain.ErrNotFound
}
