package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/madiaslanov/qazTil/internal/domain"
)

func TestAnswerSpendsALifeOnlyOnTheFirstMiss(t *testing.T) {
	now := time.Date(2026, 10, 10, 12, 0, 0, 0, time.UTC)
	course := &memCourse{exercises: []domain.Exercise{{ID: 4, LessonID: 2, Answer: "Рақмет"}}}
	stats := &memStats{stats: domain.DefaultStats(7)}
	svc := NewAttemptService(course, stats)
	svc.now = func() time.Time { return now }

	result, err := svc.Answer(context.Background(), 7, 4, "жоқ")
	if err != nil || result.Correct {
		t.Fatalf("result = %+v err=%v", result, err)
	}
	if stats.stats.Lives != 4 || stats.stats.XP != 0 || stats.applied != 1 {
		t.Fatalf("stats = %+v applied=%d", stats.stats, stats.applied)
	}

	if _, err := svc.Answer(context.Background(), 7, 4, "жоқ"); err != nil {
		t.Fatal(err)
	}
	if stats.stats.Lives != 4 {
		t.Fatalf("second miss lives = %d", stats.stats.Lives)
	}
}

func TestAnswerAwardsXPOnTheFirstCorrectAttempt(t *testing.T) {
	course := &memCourse{exercises: []domain.Exercise{{ID: 4, LessonID: 2, Answer: "Рақмет"}}}
	stats := &memStats{stats: domain.DefaultStats(7)}
	svc := NewAttemptService(course, stats)
	svc.now = func() time.Time { return time.Date(2026, 10, 10, 12, 0, 0, 0, time.UTC) }

	result, err := svc.Answer(context.Background(), 7, 4, "  рақмет ")
	if err != nil || !result.Correct || stats.stats.XP != domain.XPPerFirstCorrect || stats.stats.Lives != domain.MaxLives {
		t.Fatalf("result = %+v stats = %+v err=%v", result, stats.stats, err)
	}
}

func TestAnswerRejectsWhenNoLivesRemain(t *testing.T) {
	course := &memCourse{exercises: []domain.Exercise{{ID: 4, LessonID: 2, Answer: "Рақмет"}}}
	stats := &memStats{stats: domain.LearnerStats{UserID: 7, Lives: 0, MaxLives: domain.MaxLives}, refuse: true}
	svc := NewAttemptService(course, stats)

	_, err := svc.Answer(context.Background(), 7, 4, "рақмет")
	if !errors.Is(err, domain.ErrNoLives) || stats.applied != 0 {
		t.Fatalf("err = %v applied=%d", err, stats.applied)
	}
}

type memStats struct {
	stats   domain.LearnerStats
	applied int
	refuse  bool
}

func (m *memStats) Current(context.Context, int64, time.Time) (domain.LearnerStats, error) {
	return m.stats, nil
}

func (m *memStats) Apply(_ context.Context, _, _, _ int64, correct bool, now time.Time) error {
	if m.refuse || m.stats.Lives == 0 {
		return domain.ErrNoLives
	}
	m.applied++
	prior := m.applied > 1
	if !prior {
		if correct {
			m.stats = m.stats.GainXP()
		} else {
			m.stats = m.stats.LoseLife(now)
		}
	}
	return nil
}
