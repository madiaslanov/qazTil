package domain

import "testing"

func TestAssemblePathLocksUntilCurrent(t *testing.T) {
	units := []Unit{
		{ID: 1, TitleRU: "Кафе", Position: 1},
		{ID: 2, TitleRU: "Магазин", Position: 2},
	}
	lessons := []Lesson{
		{ID: 10, UnitID: 1, TitleRU: "Заказ", Position: 1},
		{ID: 11, UnitID: 1, TitleRU: "Счёт", Position: 2},
		{ID: 20, UnitID: 2, TitleRU: "Покупка", Position: 1},
	}
	completed := map[int64]bool{10: true}

	got := AssemblePath(units, lessons, completed)
	if got.CurrentLessonID == nil || *got.CurrentLessonID != 11 {
		t.Fatalf("current = %v", got.CurrentLessonID)
	}
	if len(got.Units) != 2 || len(got.Units[0].Lessons) != 2 || len(got.Units[1].Lessons) != 1 {
		t.Fatalf("shape = %+v", got.Units)
	}
	if got.Units[0].Lessons[0].State != PathCompleted {
		t.Fatalf("lesson 10 = %s", got.Units[0].Lessons[0].State)
	}
	if got.Units[0].Lessons[1].State != PathCurrent {
		t.Fatalf("lesson 11 = %s", got.Units[0].Lessons[1].State)
	}
	if got.Units[1].Lessons[0].State != PathLocked {
		t.Fatalf("lesson 20 = %s", got.Units[1].Lessons[0].State)
	}
}

func TestAssemblePathStartsAtFirstLesson(t *testing.T) {
	units := []Unit{{ID: 1}}
	lessons := []Lesson{{ID: 1, UnitID: 1}, {ID: 2, UnitID: 1}}

	got := AssemblePath(units, lessons, nil)
	if got.CurrentLessonID == nil || *got.CurrentLessonID != 1 {
		t.Fatalf("current = %v", got.CurrentLessonID)
	}
	if got.Units[0].Lessons[0].State != PathCurrent || got.Units[0].Lessons[1].State != PathLocked {
		t.Fatalf("states = %+v", got.Units[0].Lessons)
	}
}

func TestAssemblePathAllCompleted(t *testing.T) {
	units := []Unit{{ID: 1}}
	lessons := []Lesson{{ID: 1, UnitID: 1}, {ID: 2, UnitID: 1}}
	completed := map[int64]bool{1: true, 2: true}

	got := AssemblePath(units, lessons, completed)
	if got.CurrentLessonID != nil {
		t.Fatalf("current = %v", got.CurrentLessonID)
	}
	if got.Units[0].Lessons[0].State != PathCompleted || got.Units[0].Lessons[1].State != PathCompleted {
		t.Fatalf("states = %+v", got.Units[0].Lessons)
	}
}
