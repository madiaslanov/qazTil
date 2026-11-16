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

// AttemptResult is what the learner sees after one answer: whether it was
// right, the line they were after, and the state the answer left them in.
type AttemptResult struct {
	Correct      bool
	Expected     string
	Lives        int
	MaxLives     int
	LessonFailed bool
}

// Answer stores one attempt. The first wrong attempt spends a life.
// A correct first attempt awards XP. No lives means the attempt is refused.
//
// The answer arrives in the shape its exercise kind calls for, so a match
// exercise is never compared as a string.
func (s *AttemptService) Answer(ctx context.Context, userID, exerciseID int64, answer domain.ExerciseAnswer) (AttemptResult, error) {
	if userID <= 0 || exerciseID <= 0 {
		return AttemptResult{}, domain.ErrInvalid
	}
	exercise, err := s.course.GetExercise(ctx, exerciseID)
	if err != nil {
		return AttemptResult{}, err
	}
	correct, expected, err := domain.CheckAnswer(exercise, answer)
	if err != nil {
		return AttemptResult{}, err
	}
	if err := s.stats.Apply(ctx, userID, exercise.ID, exercise.LessonID, correct, s.now()); err != nil {
		return AttemptResult{}, err
	}
	stats, err := s.stats.Current(ctx, userID, s.now())
	if err != nil {
		return AttemptResult{}, err
	}
	return AttemptResult{
		Correct:  correct,
		Expected: expected,
		Lives:    stats.Lives,
		MaxLives: stats.MaxLives,
		// Out of lives after a wrong answer means the lesson stops here:
		// the learner has to start it again.
		LessonFailed: !correct && stats.Lives == 0,
	}, nil
}
