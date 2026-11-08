package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/madiaslanov/qazTil/internal/adapter/postgres"
	"github.com/madiaslanov/qazTil/internal/domain"
)

func TestProfilePrefsRoundTrip(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	if err := postgres.Migrate(ctx, db); err != nil {
		t.Fatal(err)
	}

	user, err := postgres.NewUserRepo(db).Create(ctx, domain.User{
		Email:        "profile@example.com",
		PasswordHash: "hash",
		DisplayName:  "Профиль",
		CreatedAt:    time.Now().UTC(),
	})
	if err != nil {
		t.Fatal(err)
	}

	repo := postgres.NewProfileRepo(db)
	defaults, err := repo.Get(ctx, user.ID)
	if err != nil {
		t.Fatal(err)
	}
	if defaults != domain.DefaultProfilePrefs() {
		t.Fatalf("defaults = %+v", defaults)
	}

	prefs := domain.ProfilePrefs{
		DailyGoal:  domain.DailyGoalLong,
		Level:      domain.LevelA1,
		Motivation: domain.MotivationStudy,
	}
	if err := repo.Save(ctx, user.ID, prefs); err != nil {
		t.Fatal(err)
	}
	loaded, err := repo.Get(ctx, user.ID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded != prefs {
		t.Fatalf("loaded = %+v", loaded)
	}
}
