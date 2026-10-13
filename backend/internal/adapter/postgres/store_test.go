package postgres_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/madiaslanov/qazTil/internal/adapter/postgres"
	"github.com/madiaslanov/qazTil/internal/domain"
)

func TestSeedQuizAndProgress(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	if err := postgres.Migrate(ctx, db); err != nil {
		t.Fatal(err)
	}
	seeded, err := postgres.Seed(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	if !seeded {
		t.Fatal("expected seed")
	}
	seeded, err = postgres.Seed(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	if seeded {
		t.Fatal("seed ran twice")
	}

	categories, err := postgres.NewCategoryRepo(db).List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(categories) != 3 {
		t.Fatalf("categories = %d", len(categories))
	}
	categoryID := categories[0].ID
	words, err := postgres.NewWordRepo(db).List(ctx, &categoryID)
	if err != nil {
		t.Fatal(err)
	}
	if len(words) != 6 {
		t.Fatalf("words = %d", len(words))
	}

	now := time.Date(2026, 9, 22, 12, 0, 0, 0, time.UTC)
	created, err := postgres.NewQuizRepo(db).Create(ctx, domain.Quiz{
		CategoryID: categoryID,
		CreatedAt:  now,
		Questions: []domain.QuizQuestion{{
			WordID:       words[0].ID,
			Prompt:       words[0].Russian,
			Options:      []string{words[0].Kazakh, words[1].Kazakh, words[2].Kazakh, words[3].Kazakh},
			CorrectIndex: 0,
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := created.ApplyAnswer(created.Questions[0].ID, 0, now); err != nil {
		t.Fatal(err)
	}
	if err := postgres.NewQuizRepo(db).Save(ctx, created); err != nil {
		t.Fatal(err)
	}
	loaded, err := postgres.NewQuizRepo(db).Get(ctx, created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.FinishedAt == nil || loaded.Questions[0].Correct == nil || !*loaded.Questions[0].Correct {
		t.Fatalf("loaded quiz = %+v", loaded)
	}
	if loaded.Questions[0].Options[0] != words[0].Kazakh {
		t.Fatalf("options = %v", loaded.Questions[0].Options)
	}

	if err := postgres.NewProgressRepo(db).Record(ctx, categoryID, true, now); err != nil {
		t.Fatal(err)
	}
	if err := postgres.NewProgressRepo(db).Record(ctx, categoryID, false, now); err != nil {
		t.Fatal(err)
	}
	stats, err := postgres.NewProgressRepo(db).List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(stats) != 1 || stats[0].TotalAnswers != 2 || stats[0].CorrectAnswers != 1 {
		t.Fatalf("progress = %+v", stats)
	}

	if err := postgres.NewCategoryRepo(db).Delete(ctx, categoryID); err != nil {
		t.Fatal(err)
	}
	left, err := postgres.NewWordRepo(db).List(ctx, &categoryID)
	if err != nil {
		t.Fatal(err)
	}
	if len(left) != 0 {
		t.Fatalf("words left after category delete: %d", len(left))
	}
}

func TestUserEmailIsUnique(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	if err := postgres.Migrate(ctx, db); err != nil {
		t.Fatal(err)
	}
	repo := postgres.NewUserRepo(db)
	user := domain.User{
		Email:        "aida@example.com",
		PasswordHash: "hash",
		DisplayName:  "Аида",
		CreatedAt:    time.Date(2026, 10, 13, 9, 0, 0, 0, time.UTC),
	}
	created, err := repo.Create(ctx, user)
	if err != nil {
		t.Fatal(err)
	}
	if created.ID == 0 || created.Email != user.Email {
		t.Fatalf("created = %+v", created)
	}
	_, err = repo.Create(ctx, user)
	if !errors.Is(err, domain.ErrEmailTaken) {
		t.Fatalf("duplicate err = %v", err)
	}
	loaded, err := repo.GetByEmail(ctx, user.Email)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.ID != created.ID || !loaded.CreatedAt.Equal(user.CreatedAt) {
		t.Fatalf("loaded = %+v", loaded)
	}
}

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	adminDSN := os.Getenv("DATABASE_URL")
	if adminDSN == "" {
		adminDSN = "postgres://qaztil:qaztil@localhost:5432/qaztil?sslmode=disable"
	}
	admin, err := postgres.Open(adminDSN)
	if err != nil {
		t.Fatalf("postgres is required for tests: %v", err)
	}

	name := fmt.Sprintf("qaztil_t_%d", time.Now().UnixNano())
	if _, err := admin.ExecContext(context.Background(), "CREATE DATABASE "+name); err != nil {
		_ = admin.Close()
		t.Fatal(err)
	}

	testDSN, err := withDatabase(adminDSN, name)
	if err != nil {
		_, _ = admin.ExecContext(context.Background(), "DROP DATABASE IF EXISTS "+name+" WITH (FORCE)")
		_ = admin.Close()
		t.Fatal(err)
	}
	db, err := postgres.Open(testDSN)
	if err != nil {
		_, _ = admin.ExecContext(context.Background(), "DROP DATABASE IF EXISTS "+name+" WITH (FORCE)")
		_ = admin.Close()
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_ = db.Close()
		_, _ = admin.ExecContext(context.Background(), "DROP DATABASE IF EXISTS "+name+" WITH (FORCE)")
		_ = admin.Close()
	})
	return db
}

func withDatabase(dsn, name string) (string, error) {
	parsed, err := url.Parse(dsn)
	if err != nil {
		return "", err
	}
	parsed.Path = "/" + name
	return parsed.String(), nil
}
