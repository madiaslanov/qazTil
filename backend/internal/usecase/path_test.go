package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/madiaslanov/qazTil/internal/domain"
)

func TestPathMarksCurrentLesson(t *testing.T) {
	course := &memCourse{
		units: []domain.Unit{
			{ID: 1, TitleRU: "Кафе", Position: 1},
			{ID: 2, TitleRU: "Магазин", Position: 2},
		},
		lessons: []domain.Lesson{
			{ID: 10, UnitID: 1, TitleRU: "Заказ", Position: 1},
			{ID: 11, UnitID: 1, TitleRU: "Счёт", Position: 2},
			{ID: 20, UnitID: 2, TitleRU: "Покупка", Position: 1},
		},
	}
	progress := &memCourseProgress{items: domain.CourseProgress{
		Lessons: []domain.LessonProgress{
			{LessonID: 10, Completed: true, Correct: 2, Total: 2},
			{LessonID: 11, Completed: false, Correct: 0, Total: 2},
			{LessonID: 20, Completed: false, Correct: 0, Total: 1},
		},
	}}
	path, err := NewPathService(course, progress).Get(context.Background(), 7)
	if err != nil {
		t.Fatal(err)
	}
	if path.CurrentLessonID == nil || *path.CurrentLessonID != 11 {
		t.Fatalf("current = %v", path.CurrentLessonID)
	}
	if path.Units[0].Lessons[0].State != domain.PathCompleted {
		t.Fatalf("first = %s", path.Units[0].Lessons[0].State)
	}
	if path.Units[0].Lessons[1].State != domain.PathCurrent {
		t.Fatalf("second = %s", path.Units[0].Lessons[1].State)
	}
	if path.Units[1].Lessons[0].State != domain.PathLocked {
		t.Fatalf("third = %s", path.Units[1].Lessons[0].State)
	}
}

func TestPathRejectsMissingUser(t *testing.T) {
	_, err := NewPathService(&memCourse{}, &memCourseProgress{}).Get(context.Background(), 0)
	if !errors.Is(err, domain.ErrInvalid) {
		t.Fatalf("err = %v", err)
	}
}

type memCourseProgress struct {
	items domain.CourseProgress
}

func (m *memCourseProgress) List(_ context.Context, userID int64) (domain.CourseProgress, error) {
	if userID <= 0 {
		return domain.CourseProgress{}, domain.ErrInvalid
	}
	return m.items, nil
}
