package usecase

import (
	"context"
	"time"

	"github.com/madiaslanov/qazTil/internal/domain"
)

// AttemptService checks an exercise answer and records the learner's attempt.
type AttemptService struct {
	course domain.CourseRepository
	stats  domain.LearnerStatsRepository
	now    func() time.Time
}

func NewAttemptService(course domain.CourseRepository, stats domain.LearnerStatsRepository) *AttemptService {
	return &AttemptService{course: course, stats: stats, now: func() time.Time { return time.Now().UTC() }}
}

func (s *AttemptService) Stats(ctx context.Context, userID int64) (domain.LearnerStats, error) {
	if userID <= 0 {
		return domain.LearnerStats{}, domain.ErrInvalid
	}
	return s.stats.Current(ctx, userID, s.now())
}

// Answer stores one attempt. The first wrong attempt spends a life.
// A correct first attempt awards XP. No lives means the attempt is refused.
func (s *AttemptService) Answer(ctx context.Context, userID, exerciseID int64, answer string) (CheckResult, error) {
	if userID <= 0 || exerciseID <= 0 {
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
	correct := got == want
	if err := s.stats.Apply(ctx, userID, exercise.ID, exercise.LessonID, correct, s.now()); err != nil {
		return CheckResult{}, err
	}
	return CheckResult{Correct: correct, Expected: exercise.Answer}, nil
}
