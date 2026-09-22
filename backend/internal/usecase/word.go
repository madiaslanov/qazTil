package usecase

import (
	"context"
	"strings"

	"github.com/madiaslanov/qazTil/internal/domain"
)

// WordService manages Kazakh vocabulary.
type WordService struct {
	categories domain.CategoryRepository
	words      domain.WordRepository
}

func NewWordService(categories domain.CategoryRepository, words domain.WordRepository) *WordService {
	return &WordService{categories: categories, words: words}
}

func (s *WordService) List(ctx context.Context, categoryID *int64) ([]domain.Word, error) {
	if categoryID != nil && *categoryID <= 0 {
		return nil, domain.ErrInvalid
	}
	return s.words.List(ctx, categoryID)
}

func (s *WordService) Get(ctx context.Context, id int64) (domain.Word, error) {
	if id <= 0 {
		return domain.Word{}, domain.ErrInvalid
	}
	return s.words.Get(ctx, id)
}

func (s *WordService) Create(ctx context.Context, word domain.Word) (domain.Word, error) {
	word = normalizeWord(word)
	if err := validateWord(word); err != nil {
		return domain.Word{}, err
	}
	if _, err := s.categories.Get(ctx, word.CategoryID); err != nil {
		return domain.Word{}, err
	}
	return s.words.Create(ctx, word)
}

func (s *WordService) Update(ctx context.Context, word domain.Word) (domain.Word, error) {
	if word.ID <= 0 {
		return domain.Word{}, domain.ErrInvalid
	}
	word = normalizeWord(word)
	if err := validateWord(word); err != nil {
		return domain.Word{}, err
	}
	if _, err := s.categories.Get(ctx, word.CategoryID); err != nil {
		return domain.Word{}, err
	}
	return s.words.Update(ctx, word)
}

func (s *WordService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return domain.ErrInvalid
	}
	return s.words.Delete(ctx, id)
}

func normalizeWord(word domain.Word) domain.Word {
	word.Kazakh = strings.TrimSpace(word.Kazakh)
	word.Russian = strings.TrimSpace(word.Russian)
	word.Transcription = strings.TrimSpace(word.Transcription)
	word.ExampleKK = strings.TrimSpace(word.ExampleKK)
	word.ExampleRU = strings.TrimSpace(word.ExampleRU)
	return word
}

func validateWord(word domain.Word) error {
	if word.CategoryID <= 0 || word.Kazakh == "" || word.Russian == "" || word.Transcription == "" {
		return domain.ErrInvalid
	}
	return nil
}
