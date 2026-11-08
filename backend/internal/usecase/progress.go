package usecase

import (
	"context"

	"github.com/madiaslanov/qazTil/internal/domain"
)

// ProgressService reports lesson and unit progress from exercise attempts.
type ProgressService struct {
	progress domain.CourseProgressRepository
}

func NewProgressService(progress domain.CourseProgressRepository) *ProgressService {
	return &ProgressService{progress: progress}
}

func (s *ProgressService) List(ctx context.Context, userID int64) (domain.CourseProgress, error) {
	if userID <= 0 {
		return domain.CourseProgress{}, domain.ErrInvalid
	}
	progress, err := s.progress.List(ctx, userID)
	if err != nil {
		return domain.CourseProgress{}, err
	}
	if progress.Lessons == nil {
		progress.Lessons = []domain.LessonProgress{}
	}
	if progress.Units == nil {
		progress.Units = []domain.UnitProgress{}
	}
	return progress, nil
}
