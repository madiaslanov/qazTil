package sqlite_test

import (
	"context"
	"testing"

	"github.com/madiaslanov/qazTil/internal/adapter/sqlite"
	"github.com/madiaslanov/qazTil/internal/domain"
	"github.com/madiaslanov/qazTil/internal/usecase"
)

func TestSeedCourse(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	if err := sqlite.Migrate(ctx, db); err != nil {
		t.Fatal(err)
	}

	seeded, err := sqlite.SeedCourse(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	if !seeded {
		t.Fatal("expected course seed")
	}
	seeded, err = sqlite.SeedCourse(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	if seeded {
		t.Fatal("course seed ran twice")
	}

	situations, err := sqlite.NewSituationRepo(db).List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(situations) != 3 || situations[0].WordCount != 5 || situations[0].TitleRU != "В кафе" {
		t.Fatalf("situations = %+v", situations)
	}

	opened, words, err := usecase.NewSituationService(sqlite.NewSituationRepo(db)).Open(ctx, situations[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	if opened.Slug != "cafe" || len(words) != 5 || words[0].Kazakh == "" {
		t.Fatalf("open = %+v words=%d", opened, len(words))
	}

	course := sqlite.NewCourseRepo(db)
	units, err := course.ListUnits(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(units) != 3 || units[0].SituationID != situations[0].ID || units[0].LessonCount != 2 {
		t.Fatalf("units = %+v", units)
	}

	_, lessons, err := usecase.NewCourseService(course).OpenUnit(ctx, units[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(lessons) != 2 || lessons[0].TitleRU != "Заказ" {
		t.Fatalf("lessons = %+v", lessons)
	}

	_, exercises, err := usecase.NewCourseService(course).OpenLesson(ctx, lessons[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(exercises) != 3 || exercises[0].Kind != domain.ExerciseChoice || exercises[0].Answer == "" {
		t.Fatalf("exercises = %+v", exercises)
	}
	if len(exercises[2].Options) != 0 {
		t.Fatalf("translate options = %v", exercises[2].Options)
	}

	result, err := usecase.NewCourseService(course).Check(ctx, exercises[0].ID, exercises[0].Answer)
	if err != nil || !result.Correct {
		t.Fatalf("check = %+v err=%v", result, err)
	}
}
