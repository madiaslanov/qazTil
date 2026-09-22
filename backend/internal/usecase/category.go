package usecase

import (
	"context"
	"strings"

	"github.com/madiaslanov/qazTil/internal/domain"
)

// CategoryService manages vocabulary topics.
type CategoryService struct {
	categories domain.CategoryRepository
}

func NewCategoryService(categories domain.CategoryRepository) *CategoryService {
	return &CategoryService{categories: categories}
}

func (s *CategoryService) List(ctx context.Context) ([]domain.Category, error) {
	return s.categories.List(ctx)
}

func (s *CategoryService) Get(ctx context.Context, id int64) (domain.Category, error) {
	if id <= 0 {
		return domain.Category{}, domain.ErrInvalid
	}
	return s.categories.Get(ctx, id)
}

func (s *CategoryService) Create(ctx context.Context, category domain.Category) (domain.Category, error) {
	category = normalizeCategory(category)
	if err := validateCategory(category); err != nil {
		return domain.Category{}, err
	}
	return s.categories.Create(ctx, category)
}

func (s *CategoryService) Update(ctx context.Context, category domain.Category) (domain.Category, error) {
	if category.ID <= 0 {
		return domain.Category{}, domain.ErrInvalid
	}
	category = normalizeCategory(category)
	if err := validateCategory(category); err != nil {
		return domain.Category{}, err
	}
	return s.categories.Update(ctx, category)
}

func (s *CategoryService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return domain.ErrInvalid
	}
	return s.categories.Delete(ctx, id)
}

func normalizeCategory(category domain.Category) domain.Category {
	category.NameKK = strings.TrimSpace(category.NameKK)
	category.NameRU = strings.TrimSpace(category.NameRU)
	category.Description = strings.TrimSpace(category.Description)
	return category
}

func validateCategory(category domain.Category) error {
	if category.NameKK == "" || category.NameRU == "" {
		return domain.ErrInvalid
	}
	return nil
}
