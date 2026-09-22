package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/madiaslanov/qazTil/internal/domain"
)

func TestStartQuizBuildsOptions(t *testing.T) {
	svc, _, _ := newQuizFixture(sampleWords())
	quiz, err := svc.Start(context.Background(), 1, 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(quiz.Questions) != 3 {
		t.Fatalf("questions = %d, want 3", len(quiz.Questions))
	}
	byRussian := map[string]domain.Word{}
	for _, word := range sampleWords() {
		byRussian[word.Russian] = word
	}
	for _, question := range quiz.Questions {
		word, ok := byRussian[question.Prompt]
		if !ok {
			t.Fatalf("unexpected prompt %q", question.Prompt)
		}
		if len(question.Options) != optionCount {
			t.Fatalf("options = %d", len(question.Options))
		}
		if question.Options[question.CorrectIndex] != word.Kazakh {
			t.Fatalf("correct option = %q, want %q", question.Options[question.CorrectIndex], word.Kazakh)
		}
		seen := map[string]int{}
		for _, option := range question.Options {
			seen[option]++
		}
		if len(seen) != optionCount {
			t.Fatalf("options are not unique: %v", question.Options)
		}
	}
}

func TestStartQuizDefaultsAndValidation(t *testing.T) {
	svc, _, _ := newQuizFixture(sampleWords())
	ctx := context.Background()

	quiz, err := svc.Start(ctx, 1, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(quiz.Questions) != defaultQuizSize {
		t.Fatalf("questions = %d, want %d", len(quiz.Questions), defaultQuizSize)
	}

	_, err = svc.Start(ctx, 1, -1)
	if !errors.Is(err, domain.ErrInvalid) {
		t.Fatalf("size -1: %v", err)
	}
	_, err = svc.Start(ctx, 1, maxQuizSize+1)
	if !errors.Is(err, domain.ErrInvalid) {
		t.Fatalf("size too large: %v", err)
	}
	_, err = svc.Start(ctx, 99, 1)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("missing category: %v", err)
	}

	short := sampleWords()[:3]
	shortSvc, _, _ := newQuizFixture(short)
	_, err = shortSvc.Start(ctx, 1, 1)
	if !errors.Is(err, domain.ErrNotEnoughWords) {
		t.Fatalf("short category: %v", err)
	}
}

func TestAnswerUpdatesProgressAndRejectsRepeat(t *testing.T) {
	fixed := time.Date(2026, 9, 22, 12, 0, 0, 0, time.UTC)
	svc, progress, _ := newQuizFixture(sampleWords())
	svc.now = func() time.Time { return fixed }
	ctx := context.Background()

	quiz, err := svc.Start(ctx, 1, 2)
	if err != nil {
		t.Fatal(err)
	}

	first := quiz.Questions[0]
	wrong := 0
	if wrong == first.CorrectIndex {
		wrong = 1
	}
	quiz, err = svc.Answer(ctx, quiz.ID, first.ID, wrong)
	if err != nil {
		t.Fatal(err)
	}
	if progress.items[1].TotalAnswers != 1 || progress.items[1].CorrectAnswers != 0 {
		t.Fatalf("progress after wrong answer = %+v", progress.items[1])
	}
	if progress.items[1].LastStudiedAt == nil || !progress.items[1].LastStudiedAt.Equal(fixed) {
		t.Fatalf("last studied = %v", progress.items[1].LastStudiedAt)
	}

	_, err = svc.Answer(ctx, quiz.ID, first.ID, first.CorrectIndex)
	if !errors.Is(err, domain.ErrAlreadyAnswered) {
		t.Fatalf("repeat answer: %v", err)
	}
	if progress.items[1].TotalAnswers != 1 {
		t.Fatalf("progress changed after repeat: %+v", progress.items[1])
	}

	second := quiz.Questions[1]
	quiz, err = svc.Answer(ctx, quiz.ID, second.ID, second.CorrectIndex)
	if err != nil {
		t.Fatal(err)
	}
	if quiz.FinishedAt == nil {
		t.Fatal("quiz should be finished")
	}
	correct, total := quiz.Score()
	if correct != 1 || total != 2 {
		t.Fatalf("score = %d/%d", correct, total)
	}
	if progress.items[1].TotalAnswers != 2 || progress.items[1].CorrectAnswers != 1 {
		t.Fatalf("progress = %+v", progress.items[1])
	}
}

func newQuizFixture(words []domain.Word) (*QuizService, *memProgress, *memQuizzes) {
	categories := &memCategories{items: []domain.Category{{
		ID:     1,
		NameKK: "Сәлемдесу",
		NameRU: "Приветствия",
	}}}
	progress := newMemProgress()
	quizzes := newMemQuizzes()
	svc := NewQuizService(categories, &memWords{items: words}, quizzes, progress)
	return svc, progress, quizzes
}

func sampleWords() []domain.Word {
	pairs := [][2]string{
		{"Сәлем", "Привет"},
		{"Рақмет", "Спасибо"},
		{"Иә", "Да"},
		{"Жоқ", "Нет"},
		{"Су", "Вода"},
		{"Нан", "Хлеб"},
	}
	words := make([]domain.Word, len(pairs))
	for i, pair := range pairs {
		words[i] = domain.Word{
			ID:            int64(i + 1),
			CategoryID:    1,
			Kazakh:        pair[0],
			Russian:       pair[1],
			Transcription: "t",
		}
	}
	return words
}
