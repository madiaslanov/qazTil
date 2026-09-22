package usecase

import (
	"context"
	"time"

	"github.com/madiaslanov/qazTil/internal/domain"
)

type memCategories struct {
	items []domain.Category
	next  int64
}

func (m *memCategories) List(context.Context) ([]domain.Category, error) {
	return append([]domain.Category(nil), m.items...), nil
}

func (m *memCategories) Get(_ context.Context, id int64) (domain.Category, error) {
	for _, item := range m.items {
		if item.ID == id {
			return item, nil
		}
	}
	return domain.Category{}, domain.ErrNotFound
}

func (m *memCategories) Create(_ context.Context, category domain.Category) (domain.Category, error) {
	m.next++
	category.ID = m.next
	m.items = append(m.items, category)
	return category, nil
}

func (m *memCategories) Update(_ context.Context, category domain.Category) (domain.Category, error) {
	for i, item := range m.items {
		if item.ID == category.ID {
			m.items[i] = category
			return category, nil
		}
	}
	return domain.Category{}, domain.ErrNotFound
}

func (m *memCategories) Delete(_ context.Context, id int64) error {
	for i, item := range m.items {
		if item.ID == id {
			m.items = append(m.items[:i], m.items[i+1:]...)
			return nil
		}
	}
	return domain.ErrNotFound
}

type memWords struct {
	items []domain.Word
	next  int64
}

func (m *memWords) List(_ context.Context, categoryID *int64) ([]domain.Word, error) {
	out := []domain.Word{}
	for _, item := range m.items {
		if categoryID != nil && item.CategoryID != *categoryID {
			continue
		}
		out = append(out, item)
	}
	return out, nil
}

func (m *memWords) Get(_ context.Context, id int64) (domain.Word, error) {
	for _, item := range m.items {
		if item.ID == id {
			return item, nil
		}
	}
	return domain.Word{}, domain.ErrNotFound
}

func (m *memWords) Create(_ context.Context, word domain.Word) (domain.Word, error) {
	m.next++
	word.ID = m.next
	m.items = append(m.items, word)
	return word, nil
}

func (m *memWords) Update(_ context.Context, word domain.Word) (domain.Word, error) {
	for i, item := range m.items {
		if item.ID == word.ID {
			m.items[i] = word
			return word, nil
		}
	}
	return domain.Word{}, domain.ErrNotFound
}

func (m *memWords) Delete(_ context.Context, id int64) error {
	for i, item := range m.items {
		if item.ID == id {
			m.items = append(m.items[:i], m.items[i+1:]...)
			return nil
		}
	}
	return domain.ErrNotFound
}

type memQuizzes struct {
	items map[int64]domain.Quiz
	next  int64
}

func newMemQuizzes() *memQuizzes {
	return &memQuizzes{items: map[int64]domain.Quiz{}}
}

func (m *memQuizzes) Create(_ context.Context, quiz domain.Quiz) (domain.Quiz, error) {
	m.next++
	quiz.ID = m.next
	for i := range quiz.Questions {
		quiz.Questions[i].ID = m.next*100 + int64(i+1)
	}
	m.items[quiz.ID] = cloneQuiz(quiz)
	return cloneQuiz(quiz), nil
}

func (m *memQuizzes) Get(_ context.Context, id int64) (domain.Quiz, error) {
	quiz, ok := m.items[id]
	if !ok {
		return domain.Quiz{}, domain.ErrNotFound
	}
	return cloneQuiz(quiz), nil
}

func (m *memQuizzes) Save(_ context.Context, quiz domain.Quiz) error {
	if _, ok := m.items[quiz.ID]; !ok {
		return domain.ErrNotFound
	}
	m.items[quiz.ID] = cloneQuiz(quiz)
	return nil
}

func cloneQuiz(quiz domain.Quiz) domain.Quiz {
	copied := quiz
	copied.Questions = append([]domain.QuizQuestion(nil), quiz.Questions...)
	for i := range copied.Questions {
		copied.Questions[i].Options = append([]string(nil), quiz.Questions[i].Options...)
		if quiz.Questions[i].SelectedIndex != nil {
			selected := *quiz.Questions[i].SelectedIndex
			copied.Questions[i].SelectedIndex = &selected
		}
		if quiz.Questions[i].Correct != nil {
			correct := *quiz.Questions[i].Correct
			copied.Questions[i].Correct = &correct
		}
	}
	if quiz.FinishedAt != nil {
		finished := *quiz.FinishedAt
		copied.FinishedAt = &finished
	}
	return copied
}

type memProgress struct {
	items map[int64]domain.Progress
}

func newMemProgress() *memProgress {
	return &memProgress{items: map[int64]domain.Progress{}}
}

func (m *memProgress) List(context.Context) ([]domain.Progress, error) {
	out := make([]domain.Progress, 0, len(m.items))
	for _, item := range m.items {
		out = append(out, item)
	}
	return out, nil
}

func (m *memProgress) Record(_ context.Context, categoryID int64, correct bool, at time.Time) error {
	item := m.items[categoryID]
	item.CategoryID = categoryID
	item.TotalAnswers++
	if correct {
		item.CorrectAnswers++
	}
	studied := at.UTC()
	item.LastStudiedAt = &studied
	m.items[categoryID] = item
	return nil
}
