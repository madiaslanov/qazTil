package usecase

import (
	"context"
	"time"

	"github.com/madiaslanov/qazTil/internal/domain"
)

// AttemptService checks an exercise answer and records the learner's attempt.
type AttemptService struct {
	course   domain.CourseRepository
	stats    domain.LearnerStatsRepository
	progress domain.CourseProgressRepository
	now      func() time.Time
}

func NewAttemptService(
	course domain.CourseRepository,
	stats domain.LearnerStatsRepository,
	progress domain.CourseProgressRepository,
) *AttemptService {
	return &AttemptService{
		course: course, stats: stats, progress: progress,
		now: func() time.Time { return time.Now().UTC() },
	}
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
	// LessonCompleted is true once every exercise of the lesson has an answer.
	LessonCompleted bool
	// XPEarned is what this single answer added. XP comes only from a correct
	// first attempt, so a repeat of the same exercise earns nothing.
	XPEarned int
	Streak   int
	// StreakIncreased is false when the learner already finished a lesson
	// today: the streak moves once a day, not once a lesson.
	StreakIncreased bool
	// NextLessonID is the lesson that opened up, or nil at the end of the path.
	NextLessonID *int64
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
	// Read the stats before the answer so the difference afterwards tells us
	// what this one answer changed, without the repository having to report it.
	before, err := s.stats.Current(ctx, userID, s.now())
	if err != nil {
		return AttemptResult{}, err
	}
	if err := s.stats.Apply(ctx, userID, exercise.ID, exercise.LessonID, correct, s.now()); err != nil {
		return AttemptResult{}, err
	}
	after, err := s.stats.Current(ctx, userID, s.now())
	if err != nil {
		return AttemptResult{}, err
	}
	result := AttemptResult{
		Correct:  correct,
		Expected: expected,
		Lives:    after.Lives,
		MaxLives: after.MaxLives,
		// Out of lives after a wrong answer means the lesson stops here:
		// the learner has to start it again.
		LessonFailed:    !correct && after.Lives == 0,
		XPEarned:        after.XP - before.XP,
		Streak:          after.Streak,
		StreakIncreased: after.Streak > before.Streak,
	}
	if s.progress == nil {
		return result, nil
	}
	lesson, err := s.progress.Lesson(ctx, userID, exercise.LessonID)
	if err != nil {
		return AttemptResult{}, err
	}
	result.LessonCompleted = lesson.Completed
	if lesson.Completed {
		next, err := s.nextLesson(ctx, exercise.LessonID)
		if err != nil {
			return AttemptResult{}, err
		}
		result.NextLessonID = next
	}
	return result, nil
}

// nextLesson walks the path in its own order — units by position, lessons
// inside them — and returns the one after the given lesson.
func (s *AttemptService) nextLesson(ctx context.Context, lessonID int64) (*int64, error) {
	lessons, err := s.course.ListAllLessons(ctx)
	if err != nil {
		return nil, err
	}
	for i, lesson := range lessons {
		if lesson.ID != lessonID {
			continue
		}
		if i+1 < len(lessons) {
			next := lessons[i+1].ID
			return &next, nil
		}
		return nil, nil
	}
	return nil, nil
}
