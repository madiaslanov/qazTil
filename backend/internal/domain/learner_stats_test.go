package domain

import (
	"testing"
	"time"
)

func TestLivesRegenerateFromTheFirstLoss(t *testing.T) {
	start := time.Date(2026, 10, 10, 12, 0, 0, 0, time.UTC)
	stats := DefaultStats(1).LoseLife(start)
	if stats.Lives != 4 || stats.NextLifeAt == nil || !stats.NextLifeAt.Equal(start.Add(LifeRecovery)) {
		t.Fatalf("first loss = %+v", stats)
	}
	again := stats.LoseLife(start.Add(10 * time.Minute))
	if again.Lives != 3 || !again.NextLifeAt.Equal(*stats.NextLifeAt) {
		t.Fatalf("timer moved = %+v", again)
	}

	due := again.Regenerated(stats.NextLifeAt.Add(LifeRecovery))
	if due.Lives != 5 || due.NextLifeAt != nil {
		t.Fatalf("regenerated = %+v", due)
	}
}

func TestLivesDoNotRegenerateEarly(t *testing.T) {
	start := time.Date(2026, 10, 10, 12, 0, 0, 0, time.UTC)
	stats := DefaultStats(1).LoseLife(start)
	early := stats.Regenerated(start.Add(LifeRecovery - time.Second))
	if early.Lives != 4 || early.NextLifeAt == nil {
		t.Fatalf("early = %+v", early)
	}
}

func TestStreakCountsLessonDays(t *testing.T) {
	monday := time.Date(2026, 10, 5, 18, 0, 0, 0, time.UTC)
	stats := DefaultStats(1).CompleteLesson(monday)
	if stats.Streak != 1 {
		t.Fatalf("first day = %+v", stats)
	}
	same := stats.CompleteLesson(monday.Add(2 * time.Hour))
	if same.Streak != 1 {
		t.Fatalf("same day = %+v", same)
	}
	next := same.CompleteLesson(monday.AddDate(0, 0, 1))
	if next.Streak != 2 {
		t.Fatalf("next day = %+v", next)
	}
	skipped := next.CompleteLesson(monday.AddDate(0, 0, 3))
	if skipped.Streak != 1 {
		t.Fatalf("skipped = %+v", skipped)
	}
}
