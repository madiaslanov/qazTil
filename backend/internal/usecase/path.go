package usecase

import (
	"context"
	"time"

	"github.com/madiaslanov/qazTil/internal/domain"
)

// PathService builds the learning-path skill tree for one learner.
type PathService struct {
	course   domain.CourseRepository
	progress domain.CourseProgressRepository
	stats    domain.LearnerStatsRepository
	now      func() time.Time
}

func NewPathService(
	course domain.CourseRepository,
	progress domain.CourseProgressRepository,
	stats domain.LearnerStatsRepository,
) *PathService {
	return &PathService{
		course: course, progress: progress, stats: stats,
		now: func() time.Time { return time.Now().UTC() },
	}
}

// Get returns units with nested lessons and lock/current/completed states.
func (s *PathService) Get(ctx context.Context, userID int64) (domain.LearningPath, error) {
	if userID <= 0 {
		return domain.LearningPath{}, domain.ErrInvalid
	}
	units, err := s.course.ListUnits(ctx)
	if err != nil {
		return domain.LearningPath{}, err
	}
	if units == nil {
		units = []domain.Unit{}
	}
	lessons, err := s.course.ListAllLessons(ctx)
	if err != nil {
		return domain.LearningPath{}, err
	}
	if lessons == nil {
		lessons = []domain.Lesson{}
	}
	progress, err := s.progress.List(ctx, userID)
	if err != nil {
		return domain.LearningPath{}, err
	}
	completed := make(map[int64]bool, len(progress.Lessons))
	for _, item := range progress.Lessons {
		if item.Completed {
			completed[item.LessonID] = true
		}
	}
	path := domain.AssemblePath(units, lessons, completed)
	if path.Units == nil {
		path.Units = []domain.PathUnit{}
	}
	return path, nil
}

// Lesson opens one lesson for the player: its exercises, where to continue and
// how many lives are left. A lesson the learner has not reached is refused
// rather than served, so the path cannot be skipped by guessing an id.
func (s *PathService) Lesson(ctx context.Context, userID, lessonID int64) (domain.LessonView, error) {
	if userID <= 0 || lessonID <= 0 {
		return domain.LessonView{}, domain.ErrInvalid
	}
	path, err := s.Get(ctx, userID)
	if err != nil {
		return domain.LessonView{}, err
	}
	state := ""
	for _, unit := range path.Units {
		for _, lesson := range unit.Lessons {
			if lesson.ID == lessonID {
				state = lesson.State
			}
		}
	}
	if state == "" {
		return domain.LessonView{}, domain.ErrNotFound
	}
	if state == domain.PathLocked {
		return domain.LessonView{}, domain.ErrLocked
	}

	lesson, err := s.course.GetLesson(ctx, lessonID)
	if err != nil {
		return domain.LessonView{}, err
	}
	exercises, err := s.course.ListExercises(ctx, lessonID)
	if err != nil {
		return domain.LessonView{}, err
	}
	if exercises == nil {
		exercises = []domain.Exercise{}
	}
	progress, err := s.progress.Lesson(ctx, userID, lessonID)
	if err != nil {
		return domain.LessonView{}, err
	}
	resume, err := s.progress.ResumePosition(ctx, userID, lessonID)
	if err != nil {
		return domain.LessonView{}, err
	}
	// Начатый заново урок продолжать неоткуда: ученик уже прошёл его целиком.
	if state == domain.PathCompleted {
		resume = 0
	}
	stats, err := s.stats.Current(ctx, userID, s.now())
	if err != nil {
		return domain.LessonView{}, err
	}
	return domain.LessonView{
		Lesson:         lesson,
		Exercises:      exercises,
		State:          state,
		ResumePosition: resume,
		Correct:        progress.Correct,
		Total:          progress.Total,
		Lives:          stats.Lives,
		MaxLives:       stats.MaxLives,
	}, nil
}
