package usecase

import (
	"context"
	"strings"

	"github.com/madiaslanov/qazTil/internal/domain"
)

// CourseService serves units, lessons and exercise checks.
type CourseService struct {
	course domain.CourseRepository
}

func NewCourseService(course domain.CourseRepository) *CourseService {
	return &CourseService{course: course}
}

func (s *CourseService) ListUnits(ctx context.Context) ([]domain.Unit, error) {
	return s.course.ListUnits(ctx)
}

// OpenUnit returns a unit and its lessons in order.
func (s *CourseService) OpenUnit(ctx context.Context, id int64) (domain.Unit, []domain.Lesson, error) {
	if id <= 0 {
		return domain.Unit{}, nil, domain.ErrInvalid
	}
	item, err := s.course.GetUnit(ctx, id)
	if err != nil {
		return domain.Unit{}, nil, err
	}
	lessons, err := s.course.ListLessons(ctx, id)
	if err != nil {
		return domain.Unit{}, nil, err
	}
	if lessons == nil {
		lessons = []domain.Lesson{}
	}
	item.LessonCount = len(lessons)
	return item, lessons, nil
}

func (s *CourseService) Lessons(ctx context.Context, unitID int64) ([]domain.Lesson, error) {
	if unitID <= 0 {
		return nil, domain.ErrInvalid
	}
	if _, err := s.course.GetUnit(ctx, unitID); err != nil {
		return nil, err
	}
	lessons, err := s.course.ListLessons(ctx, unitID)
	if err != nil {
		return nil, err
	}
	if lessons == nil {
		lessons = []domain.Lesson{}
	}
	return lessons, nil
}

// OpenLesson returns a lesson and its exercises. Answers stay on the exercises.
func (s *CourseService) OpenLesson(ctx context.Context, id int64) (domain.Lesson, []domain.Exercise, error) {
	if id <= 0 {
		return domain.Lesson{}, nil, domain.ErrInvalid
	}
	item, err := s.course.GetLesson(ctx, id)
	if err != nil {
		return domain.Lesson{}, nil, err
	}
	exercises, err := s.course.ListExercises(ctx, id)
	if err != nil {
		return domain.Lesson{}, nil, err
	}
	if exercises == nil {
		exercises = []domain.Exercise{}
	}
	return item, exercises, nil
}

func (s *CourseService) Exercises(ctx context.Context, lessonID int64) ([]domain.Exercise, error) {
	if lessonID <= 0 {
		return nil, domain.ErrInvalid
	}
	if _, err := s.course.GetLesson(ctx, lessonID); err != nil {
		return nil, err
	}
	exercises, err := s.course.ListExercises(ctx, lessonID)
	if err != nil {
		return nil, err
	}
	if exercises == nil {
		exercises = []domain.Exercise{}
	}
	return exercises, nil
}

// CheckResult is the outcome of one exercise attempt.
type CheckResult struct {
	Correct  bool
	Expected string
}

// Check compares the learner's answer with the stored line.
func (s *CourseService) Check(ctx context.Context, exerciseID int64, answer string) (CheckResult, error) {
	if exerciseID <= 0 {
		return CheckResult{}, domain.ErrInvalid
	}
	exercise, err := s.course.GetExercise(ctx, exerciseID)
	if err != nil {
		return CheckResult{}, err
	}
	got := normalizeAnswer(answer)
	want := normalizeAnswer(exercise.Answer)
	if got == "" || want == "" {
		return CheckResult{}, domain.ErrInvalid
	}
	return CheckResult{Correct: got == want, Expected: exercise.Answer}, nil
}

func normalizeAnswer(value string) string {
	return strings.ToLower(strings.Join(strings.Fields(value), " "))
}
