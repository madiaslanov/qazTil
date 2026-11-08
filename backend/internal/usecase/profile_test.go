package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/madiaslanov/qazTil/internal/domain"
)

type memStatsRepo struct {
	stats map[int64]domain.LearnerStats
}

func (m *memStatsRepo) Current(_ context.Context, userID int64, _ time.Time) (domain.LearnerStats, error) {
	if stats, ok := m.stats[userID]; ok {
		stats.MaxLives = domain.MaxLives
		return stats, nil
	}
	return domain.DefaultStats(userID), nil
}

func (m *memStatsRepo) Apply(context.Context, int64, int64, int64, bool, time.Time) error {
	return nil
}

func newProfileService() (*ProfileService, *memUsers, *memProfiles) {
	users := newMemUsers()
	profiles := newMemProfiles()
	stats := &memStatsRepo{stats: map[int64]domain.LearnerStats{}}
	return NewProfileService(users, profiles, stats), users, profiles
}

func TestProfileGetAndUpdate(t *testing.T) {
	svc, users, profiles := newProfileService()
	ctx := context.Background()

	user, err := users.Create(ctx, domain.User{
		Email:       "aidana@example.com",
		DisplayName: "Айдана",
		CreatedAt:   time.Unix(0, 0).UTC(),
	})
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	if err := profiles.Save(ctx, user.ID, domain.DefaultProfilePrefs()); err != nil {
		t.Fatalf("save prefs: %v", err)
	}

	profile, err := svc.Get(ctx, user.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if profile.Email != user.Email || profile.Prefs.DailyGoal != domain.DailyGoalMedium {
		t.Fatalf("profile = %+v", profile)
	}
	if profile.Stats.Lives != domain.MaxLives {
		t.Fatalf("lives = %d", profile.Stats.Lives)
	}

	goal := domain.DailyGoalLong
	level := domain.LevelA2
	name := "Aidana"
	updated, err := svc.Update(ctx, user.ID, ProfileUpdate{
		DisplayName: &name,
		DailyGoal:   &goal,
		Level:       &level,
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.DisplayName != name || updated.Prefs.DailyGoal != goal || updated.Prefs.Level != level {
		t.Fatalf("updated = %+v", updated)
	}
	if updated.Prefs.Motivation != domain.MotivationLife {
		t.Fatalf("motivation changed: %s", updated.Prefs.Motivation)
	}
}

func TestProfileUpdateRejectsBadGoal(t *testing.T) {
	svc, users, profiles := newProfileService()
	ctx := context.Background()
	user, err := users.Create(ctx, domain.User{Email: "a@example.com", DisplayName: "a", CreatedAt: time.Now().UTC()})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	_ = profiles.Save(ctx, user.ID, domain.DefaultProfilePrefs())

	goal := 7
	if _, err := svc.Update(ctx, user.ID, ProfileUpdate{DailyGoal: &goal}); !errors.Is(err, domain.ErrInvalid) {
		t.Fatalf("want ErrInvalid, got %v", err)
	}
}
