package usecase

import (
	"context"

	"github.com/madiaslanov/qazTil/internal/domain"
)

// ProgressService reports answer totals for every category.
type ProgressService struct {
	categories domain.CategoryRepository
	progress   domain.ProgressRepository
}

func NewProgressService(categories domain.CategoryRepository, progress domain.ProgressRepository) *ProgressService {
	return &ProgressService{categories: categories, progress: progress}
}

func (s *ProgressService) List(ctx context.Context) ([]domain.Progress, error) {
	categories, err := s.categories.List(ctx)
	if err != nil {
		return nil, err
	}
	stats, err := s.progress.List(ctx)
	if err != nil {
		return nil, err
	}
	byCategory := make(map[int64]domain.Progress, len(stats))
	for _, item := range stats {
		byCategory[item.CategoryID] = item
	}
	out := make([]domain.Progress, 0, len(categories))
	for _, category := range categories {
		item := byCategory[category.ID]
		item.CategoryID = category.ID
		item.CategoryNameKK = category.NameKK
		item.CategoryNameRU = category.NameRU
		out = append(out, item)
	}
	return out, nil
}
