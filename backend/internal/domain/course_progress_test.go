package domain

import "testing"

func TestAssembleProgressCountsAttempts(t *testing.T) {
	units := []Unit{{ID: 1}, {ID: 2}}
	lessons := []Lesson{{ID: 5, UnitID: 1}, {ID: 6, UnitID: 1}, {ID: 7, UnitID: 2}}
	exercises := []Exercise{
		{ID: 1, LessonID: 5},
		{ID: 2, LessonID: 5},
		{ID: 3, LessonID: 6},
		{ID: 4, LessonID: 7},
	}
	marks := map[int64]AttemptMark{
		1: {Correct: true},
		2: {Correct: false},
		4: {Correct: true},
	}

	got := AssembleProgress(units, lessons, exercises, marks)
	if len(got.Lessons) != 3 || !got.Lessons[0].Completed || got.Lessons[0].Correct != 1 || got.Lessons[0].Total != 2 {
		t.Fatalf("lesson 5 = %+v", got.Lessons[0])
	}
	if got.Lessons[1].Completed || got.Lessons[1].Correct != 0 || got.Lessons[1].Total != 1 {
		t.Fatalf("lesson 6 = %+v", got.Lessons[1])
	}
	if !got.Lessons[2].Completed || got.Lessons[2].Correct != 1 {
		t.Fatalf("lesson 7 = %+v", got.Lessons[2])
	}
	if got.Units[0].CompletedLessons != 1 || got.Units[0].TotalLessons != 2 {
		t.Fatalf("unit 1 = %+v", got.Units[0])
	}
	if got.Units[1].CompletedLessons != 1 || got.Units[1].TotalLessons != 1 {
		t.Fatalf("unit 2 = %+v", got.Units[1])
	}
}
