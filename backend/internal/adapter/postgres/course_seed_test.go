package postgres_test

import (
	"context"
	"testing"

	"github.com/madiaslanov/qazTil/internal/adapter/postgres"
	"github.com/madiaslanov/qazTil/internal/domain"
	"github.com/madiaslanov/qazTil/internal/usecase"
)

func TestSeedCourse(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	if err := postgres.Migrate(ctx, db); err != nil {
		t.Fatal(err)
	}

	seeded, err := postgres.SeedCourse(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	if !seeded {
		t.Fatal("expected course seed")
	}
	seeded, err = postgres.SeedCourse(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	if seeded {
		t.Fatal("course seed ran twice")
	}

	situations, err := postgres.NewSituationRepo(db).List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	// Проверяем состав, а не количество: добавление новой ситуации в сид —
	// обычное дело и ронять этот тест не должно.
	bySlug := map[string]domain.Situation{}
	for _, item := range situations {
		bySlug[item.Slug] = item
	}
	for _, slug := range []string{"cafe", "shop", "road"} {
		if _, ok := bySlug[slug]; !ok {
			t.Fatalf("ситуация %q не засеяна, есть: %+v", slug, situations)
		}
	}
	cafe := bySlug["cafe"]
	if cafe.WordCount != 5 || cafe.TitleRU != "В кафе" {
		t.Fatalf("cafe = %+v", cafe)
	}

	opened, words, err := usecase.NewSituationService(postgres.NewSituationRepo(db)).Open(ctx, cafe.ID)
	if err != nil {
		t.Fatal(err)
	}
	if opened.Slug != "cafe" || len(words) != 5 || words[0].Kazakh == "" {
		t.Fatalf("open = %+v words=%d", opened, len(words))
	}

	course := postgres.NewCourseRepo(db)
	units, err := course.ListUnits(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(units) != len(situations) {
		t.Fatalf("юнитов %d при %d ситуациях", len(units), len(situations))
	}
	var cafeUnit domain.Unit
	for _, unit := range units {
		if unit.SituationID == cafe.ID {
			cafeUnit = unit
			break
		}
	}
	if cafeUnit.ID == 0 || cafeUnit.LessonCount != 2 {
		t.Fatalf("юнит кафе = %+v", cafeUnit)
	}

	_, lessons, err := usecase.NewCourseService(course).OpenUnit(ctx, cafeUnit.ID)
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
