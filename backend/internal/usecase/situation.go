package usecase

import (
	"context"

	"github.com/madiaslanov/qazTil/internal/domain"
)

// SituationService lists real-life contexts and the words inside them.
type SituationService struct {
	situations domain.SituationRepository
}

func NewSituationService(situations domain.SituationRepository) *SituationService {
	return &SituationService{situations: situations}
}

func (s *SituationService) List(ctx context.Context) ([]domain.Situation, error) {
	return s.situations.List(ctx)
}

// Open returns a situation together with its words.
func (s *SituationService) Open(ctx context.Context, id int64) (domain.Situation, []domain.SituationWord, error) {
	if id <= 0 {
		return domain.Situation{}, nil, domain.ErrInvalid
	}
	item, err := s.situations.Get(ctx, id)
	if err != nil {
		return domain.Situation{}, nil, err
	}
	words, err := s.situations.ListWords(ctx, id)
	if err != nil {
		return domain.Situation{}, nil, err
	}
	if words == nil {
		words = []domain.SituationWord{}
	}
	return item, words, nil
}
