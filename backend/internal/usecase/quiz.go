package usecase

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/madiaslanov/qazTil/internal/domain"
)

const (
	defaultQuizSize = 5
	maxQuizSize     = 20
	optionCount     = 4
)

// QuizService builds practice sessions and records answers.
type QuizService struct {
	categories domain.CategoryRepository
	words      domain.WordRepository
	quizzes    domain.QuizRepository
	progress   domain.ProgressRepository
	now        func() time.Time
	mu         sync.Mutex
	rng        *rand.Rand
}

func NewQuizService(
	categories domain.CategoryRepository,
	words domain.WordRepository,
	quizzes domain.QuizRepository,
	progress domain.ProgressRepository,
) *QuizService {
	return &QuizService{
		categories: categories,
		words:      words,
		quizzes:    quizzes,
		progress:   progress,
		now:        func() time.Time { return time.Now().UTC() },
		rng:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *QuizService) Get(ctx context.Context, id int64) (domain.Quiz, error) {
	if id <= 0 {
		return domain.Quiz{}, domain.ErrInvalid
	}
	return s.quizzes.Get(ctx, id)
}

func (s *QuizService) Start(ctx context.Context, categoryID int64, size int) (domain.Quiz, error) {
	if categoryID <= 0 {
		return domain.Quiz{}, domain.ErrInvalid
	}
	switch {
	case size < 0 || size > maxQuizSize:
		return domain.Quiz{}, domain.ErrInvalid
	case size == 0:
		size = defaultQuizSize
	}
	if _, err := s.categories.Get(ctx, categoryID); err != nil {
		return domain.Quiz{}, err
	}
	words, err := s.words.List(ctx, &categoryID)
	if err != nil {
		return domain.Quiz{}, err
	}
	if uniqueKazakh(words) < optionCount {
		return domain.Quiz{}, domain.ErrNotEnoughWords
	}
	if size > len(words) {
		size = len(words)
	}

	s.mu.Lock()
	picked := pickWords(words, size, s.rng)
	questions := make([]domain.QuizQuestion, 0, len(picked))
	for _, word := range picked {
		options, correctIndex, err := buildOptions(word, words, s.rng)
		if err != nil {
			s.mu.Unlock()
			return domain.Quiz{}, err
		}
		questions = append(questions, domain.QuizQuestion{
			WordID:       word.ID,
			Prompt:       word.Russian,
			Options:      options,
			CorrectIndex: correctIndex,
		})
	}
	s.mu.Unlock()

	quiz := domain.Quiz{
		CategoryID: categoryID,
		CreatedAt:  s.now(),
		Questions:  questions,
	}
	return s.quizzes.Create(ctx, quiz)
}

func (s *QuizService) Answer(ctx context.Context, quizID, questionID int64, selectedIndex int) (domain.Quiz, error) {
	if quizID <= 0 || questionID <= 0 {
		return domain.Quiz{}, domain.ErrInvalid
	}
	quiz, err := s.quizzes.Get(ctx, quizID)
	if err != nil {
		return domain.Quiz{}, err
	}
	now := s.now()
	if err := quiz.ApplyAnswer(questionID, selectedIndex, now); err != nil {
		return domain.Quiz{}, err
	}
	if err := s.quizzes.Save(ctx, quiz); err != nil {
		return domain.Quiz{}, err
	}
	correct := false
	for _, question := range quiz.Questions {
		if question.ID == questionID && question.Correct != nil {
			correct = *question.Correct
			break
		}
	}
	if err := s.progress.Record(ctx, quiz.CategoryID, correct, now); err != nil {
		return domain.Quiz{}, err
	}
	return quiz, nil
}

func uniqueKazakh(words []domain.Word) int {
	seen := make(map[string]struct{}, len(words))
	for _, word := range words {
		seen[word.Kazakh] = struct{}{}
	}
	return len(seen)
}

func pickWords(words []domain.Word, n int, rng *rand.Rand) []domain.Word {
	copied := append([]domain.Word(nil), words...)
	rng.Shuffle(len(copied), func(i, j int) {
		copied[i], copied[j] = copied[j], copied[i]
	})
	return copied[:n]
}

func buildOptions(correct domain.Word, all []domain.Word, rng *rand.Rand) ([]string, int, error) {
	seen := make(map[string]struct{}, len(all))
	pool := make([]string, 0, len(all))
	for _, word := range all {
		if word.Kazakh == correct.Kazakh {
			continue
		}
		if _, ok := seen[word.Kazakh]; ok {
			continue
		}
		seen[word.Kazakh] = struct{}{}
		pool = append(pool, word.Kazakh)
	}
	if len(pool) < optionCount-1 {
		return nil, 0, domain.ErrNotEnoughWords
	}
	rng.Shuffle(len(pool), func(i, j int) {
		pool[i], pool[j] = pool[j], pool[i]
	})
	options := []string{correct.Kazakh, pool[0], pool[1], pool[2]}
	rng.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})
	correctIndex := 0
	for i, option := range options {
		if option == correct.Kazakh {
			correctIndex = i
			break
		}
	}
	return options, correctIndex, nil
}
