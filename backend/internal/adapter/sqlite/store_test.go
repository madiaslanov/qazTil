package sqlite_test

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"github.com/madiaslanov/qazTil/internal/adapter/sqlite"
	"github.com/madiaslanov/qazTil/internal/domain"
)

func TestSeedQuizAndProgress(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	if err := sqlite.Migrate(ctx, db); err != nil {
		t.Fatal(err)
	}
	seeded, err := sqlite.Seed(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	if !seeded {
		t.Fatal("expected seed")
	}
	seeded, err = sqlite.Seed(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	if seeded {
		t.Fatal("seed ran twice")
	}

	categories, err := sqlite.NewCategoryRepo(db).List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(categories) != 3 {
		t.Fatalf("categories = %d", len(categories))
	}
	categoryID := categories[0].ID
	words, err := sqlite.NewWordRepo(db).List(ctx, &categoryID)
	if err != nil {
		t.Fatal(err)
	}
	if len(words) != 6 {
		t.Fatalf("words = %d", len(words))
	}

	now := time.Date(2026, 9, 22, 12, 0, 0, 0, time.UTC)
	created, err := sqlite.NewQuizRepo(db).Create(ctx, domain.Quiz{
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
	if err := sqlite.NewQuizRepo(db).Save(ctx, created); err != nil {
		t.Fatal(err)
	}
	loaded, err := sqlite.NewQuizRepo(db).Get(ctx, created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.FinishedAt == nil || loaded.Questions[0].Correct == nil || !*loaded.Questions[0].Correct {
		t.Fatalf("loaded quiz = %+v", loaded)
	}
	if loaded.Questions[0].Options[0] != words[0].Kazakh {
		t.Fatalf("options = %v", loaded.Questions[0].Options)
	}

	if err := sqlite.NewProgressRepo(db).Record(ctx, categoryID, true, now); err != nil {
		t.Fatal(err)
	}
	if err := sqlite.NewProgressRepo(db).Record(ctx, categoryID, false, now); err != nil {
		t.Fatal(err)
	}
	stats, err := sqlite.NewProgressRepo(db).List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(stats) != 1 || stats[0].TotalAnswers != 2 || stats[0].CorrectAnswers != 1 {
		t.Fatalf("progress = %+v", stats)
	}

	if err := sqlite.NewCategoryRepo(db).Delete(ctx, categoryID); err != nil {
		t.Fatal(err)
	}
	left, err := sqlite.NewWordRepo(db).List(ctx, &categoryID)
	if err != nil {
		t.Fatal(err)
	}
	if len(left) != 0 {
		t.Fatalf("words left after category delete: %d", len(left))
	}
}

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sqlite.Open(filepath.Join(t.TempDir(), "qaztil.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}
