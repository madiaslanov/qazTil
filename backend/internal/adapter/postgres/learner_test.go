package postgres_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/madiaslanov/qazTil/internal/adapter/postgres"
	"github.com/madiaslanov/qazTil/internal/domain"
)

func TestAttemptsDriveProgressAndLives(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	if err := postgres.Migrate(ctx, db); err != nil {
		t.Fatal(err)
	}
	if _, err := postgres.SeedCourse(ctx, db); err != nil {
		t.Fatal(err)
	}
	var userID int64
	if err := db.QueryRowContext(ctx, `
		INSERT INTO users (email, password_hash, display_name, created_at)
		VALUES ('aida@example.com', 'hash', 'Аида', $1)
		RETURNING id`, time.Now().UTC()).Scan(&userID); err != nil {
		t.Fatal(err)
	}

	course := postgres.NewCourseRepo(db)
	units, err := course.ListUnits(ctx)
	if err != nil || len(units) == 0 {
		t.Fatal(err)
	}
	lessons, err := course.ListLessons(ctx, units[0].ID)
	if err != nil || len(lessons) == 0 {
		t.Fatal(err)
	}
	exercises, err := course.ListExercises(ctx, lessons[0].ID)
	if err != nil || len(exercises) < 3 {
		t.Fatalf("exercises = %d err=%v", len(exercises), err)
	}

	learner := postgres.NewLearnerRepo(db)
	now := time.Date(2026, 10, 10, 12, 0, 0, 0, time.UTC)
	if err := learner.Apply(ctx, userID, exercises[0].ID, exercises[0].LessonID, false, now); err != nil {
		t.Fatal(err)
	}
	stats, err := learner.Current(ctx, userID, now)
	if err != nil {
		t.Fatal(err)
	}
	if stats.Lives != 4 || stats.XP != 0 || stats.NextLifeAt == nil || !stats.NextLifeAt.Equal(now.Add(domain.LifeRecovery)) {
		t.Fatalf("after first miss = %+v", stats)
	}
	if err := learner.Apply(ctx, userID, exercises[0].ID, exercises[0].LessonID, false, now.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	stats, err = learner.Current(ctx, userID, now.Add(time.Minute))
	if err != nil || stats.Lives != 4 || stats.XP != 0 {
		t.Fatalf("second miss = %+v err=%v", stats, err)
	}

	if err := learner.Apply(ctx, userID, exercises[1].ID, exercises[1].LessonID, true, now); err != nil {
		t.Fatal(err)
	}
	if err := learner.Apply(ctx, userID, exercises[2].ID, exercises[2].LessonID, true, now); err != nil {
		t.Fatal(err)
	}
	// The first exercise already has an attempt, so a later correct answer adds no XP,
	// but it does count as a correct exercise and finishes the lesson.
	if err := learner.Apply(ctx, userID, exercises[0].ID, exercises[0].LessonID, true, now); err != nil {
		t.Fatal(err)
	}
	stats, err = learner.Current(ctx, userID, now)
	if err != nil {
		t.Fatal(err)
	}
	if stats.XP != domain.XPPerFirstCorrect*2 || stats.Streak != 1 || stats.Lives != 4 {
		t.Fatalf("finished lesson = %+v", stats)
	}

	progress, err := postgres.NewCourseProgressRepo(db).List(ctx, userID)
	if err != nil {
		t.Fatal(err)
	}
	var lesson domain.LessonProgress
	for _, item := range progress.Lessons {
		if item.LessonID == lessons[0].ID {
			lesson = item
		}
	}
	if !lesson.Completed || lesson.Correct != 3 || lesson.Total != 3 {
		t.Fatalf("lesson progress = %+v", lesson)
	}
	if progress.Units[0].UnitID != units[0].ID || progress.Units[0].CompletedLessons != 1 {
		t.Fatalf("unit progress = %+v", progress.Units[0])
	}

	restored, err := learner.Current(ctx, userID, now.Add(domain.LifeRecovery))
	if err != nil {
		t.Fatal(err)
	}
	if restored.Lives != domain.MaxLives || restored.NextLifeAt != nil {
		t.Fatalf("restored = %+v", restored)
	}
}

func TestNoLivesRejectsTheAttempt(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	if err := postgres.Migrate(ctx, db); err != nil {
		t.Fatal(err)
	}
	if _, err := postgres.SeedCourse(ctx, db); err != nil {
		t.Fatal(err)
	}
	var userID int64
	if err := db.QueryRowContext(ctx, `
		INSERT INTO users (email, password_hash, display_name, created_at)
		VALUES ('erasyl@example.com', 'hash', 'Ерасыл', now())
		RETURNING id`).Scan(&userID); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, `
		INSERT INTO learner_stats (user_id, xp, lives, streak)
		VALUES ($1, 0, 0, 0)`, userID); err != nil {
		t.Fatal(err)
	}
	course := postgres.NewCourseRepo(db)
	units, err := course.ListUnits(ctx)
	if err != nil {
		t.Fatal(err)
	}
	lessons, err := course.ListLessons(ctx, units[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	exercises, err := course.ListExercises(ctx, lessons[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	err = postgres.NewLearnerRepo(db).Apply(ctx, userID, exercises[0].ID, exercises[0].LessonID, true, time.Now().UTC())
	if !errors.Is(err, domain.ErrNoLives) {
		t.Fatalf("err = %v", err)
	}
	var attempts int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM attempts WHERE user_id = $1`, userID).Scan(&attempts); err != nil {
		t.Fatal(err)
	}
	if attempts != 0 {
		t.Fatalf("attempts = %d", attempts)
	}
}
