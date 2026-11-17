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
	if cafeUnit.ID == 0 || cafeUnit.LessonCount != 4 {
		t.Fatalf("юнит кафе = %+v", cafeUnit)
	}

	_, lessons, err := usecase.NewCourseService(course).OpenUnit(ctx, cafeUnit.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(lessons) != 4 || lessons[0].TitleRU != "Заказ" {
		t.Fatalf("lessons = %+v", lessons)
	}

	_, exercises, err := usecase.NewCourseService(course).OpenLesson(ctx, lessons[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(exercises) != 6 || exercises[0].Kind != domain.ExerciseChoice || exercises[0].Answer == "" {
		t.Fatalf("exercises = %+v", exercises)
	}
	// Печатать перевод больше не нужно: третье упражнение собирается из слов.
	if exercises[2].Kind != domain.ExerciseAssemble {
		t.Fatalf("третье упражнение = %q", exercises[2].Kind)
	}
	last := lessons[len(lessons)-1]
	if last.TitleRU != "Диалог: в кафе" {
		t.Fatalf("последний урок модуля = %q, ждали диалог", last.TitleRU)
	}

	result, err := usecase.NewCourseService(course).Check(ctx, exercises[0].ID, exercises[0].Answer)
	if err != nil || !result.Correct {
		t.Fatalf("check = %+v err=%v", result, err)
	}
}

func TestSeedCourseReplacesOnlyWhenCourseChanged(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	if err := postgres.Migrate(ctx, db); err != nil {
		t.Fatal(err)
	}
	if seeded, err := postgres.SeedCourse(ctx, db); err != nil || !seeded {
		t.Fatalf("первый сид: seeded=%v err=%v", seeded, err)
	}

	// Ученик с XP и ответом на первое упражнение.
	var userID, exerciseID int64
	if err := db.QueryRowContext(ctx, `INSERT INTO users (email, password_hash) VALUES ('a@b.c', 'x') RETURNING id`).Scan(&userID); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, `INSERT INTO learner_stats (user_id, xp, lives, streak) VALUES ($1, 40, 3, 2)`, userID); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRowContext(ctx, `SELECT id FROM exercises ORDER BY id LIMIT 1`).Scan(&exerciseID); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, `INSERT INTO attempts (user_id, exercise_id, correct, created_at) VALUES ($1, $2, true, now())`, userID, exerciseID); err != nil {
		t.Fatal(err)
	}

	count := func(query string) int {
		t.Helper()
		var n int
		if err := db.QueryRowContext(ctx, query).Scan(&n); err != nil {
			t.Fatal(err)
		}
		return n
	}
	exercisesBefore := count(`SELECT COUNT(*) FROM exercises`)

	// Тот же курс: ничего не пишем и ничего не стираем.
	if seeded, err := postgres.SeedCourse(ctx, db); err != nil || seeded {
		t.Fatalf("повторный сид: seeded=%v err=%v", seeded, err)
	}
	if got := count(`SELECT COUNT(*) FROM attempts`); got != 1 {
		t.Fatalf("ответ пропал без причины: %d", got)
	}

	// Курс в коде изменился (в базе лежит другой отпечаток): засеваем заново.
	if _, err := db.ExecContext(ctx, `UPDATE course_meta SET value = 'old' WHERE key = 'course_hash'`); err != nil {
		t.Fatal(err)
	}
	if seeded, err := postgres.SeedCourse(ctx, db); err != nil || !seeded {
		t.Fatalf("сид после изменения: seeded=%v err=%v", seeded, err)
	}
	if got := count(`SELECT COUNT(*) FROM exercises`); got != exercisesBefore {
		t.Fatalf("упражнений %d, было %d: курс задвоился или потерялся", got, exercisesBefore)
	}
	if got := count(`SELECT COUNT(*) FROM attempts`); got != 0 {
		t.Fatalf("ответы на старые упражнения должны уйти, осталось %d", got)
	}
	// Ученик, XP и стрик остаются.
	if got := count(`SELECT COUNT(*) FROM users`); got != 1 {
		t.Fatalf("пользователей %d", got)
	}
	if got := count(`SELECT xp FROM learner_stats`); got != 40 {
		t.Fatalf("xp = %d, ждали 40", got)
	}
	if seeded, err := postgres.SeedCourse(ctx, db); err != nil || seeded {
		t.Fatalf("после пересева отпечаток должен совпасть: seeded=%v err=%v", seeded, err)
	}
}

func TestSeedCourseUpgradesDatabaseWithoutFingerprint(t *testing.T) {
	// База, засеянная до появления отпечатка: старые уроки есть, отпечатка нет.
	ctx := context.Background()
	db := openTestDB(t)
	if err := postgres.Migrate(ctx, db); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, `INSERT INTO situations (slug, title_kk, title_ru, description, position) VALUES ('cafe', 'Дәмхана', 'В кафе', 'старый сид', 1)`); err != nil {
		t.Fatal(err)
	}
	if seeded, err := postgres.SeedCourse(ctx, db); err != nil || !seeded {
		t.Fatalf("seeded=%v err=%v", seeded, err)
	}
	var description string
	if err := db.QueryRowContext(ctx, `SELECT description FROM situations WHERE slug = 'cafe'`).Scan(&description); err != nil {
		t.Fatal(err)
	}
	if description == "старый сид" {
		t.Fatal("старая ситуация осталась, курс не обновился")
	}
	var n int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM situations WHERE slug = 'cafe'`).Scan(&n); err != nil || n != 1 {
		t.Fatalf("ситуация кафе должна быть одна: n=%d err=%v", n, err)
	}
}
