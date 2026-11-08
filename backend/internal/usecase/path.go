package usecase

import (
	"context"

	"github.com/madiaslanov/qazTil/internal/domain"
)

// PathService builds the learning-path skill tree for one learner.
type PathService struct {
	course   domain.CourseRepository
	progress domain.CourseProgressRepository
}

func NewPathService(course domain.CourseRepository, progress domain.CourseProgressRepository) *PathService {
	return &PathService{course: course, progress: progress}
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
