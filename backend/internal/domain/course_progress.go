package domain

// LessonProgress is how far one learner is through a lesson's exercises.
type LessonProgress struct {
	LessonID  int64
	Completed bool
	Correct   int
	Total     int
}

// UnitProgress counts completed lessons inside one unit.
type UnitProgress struct {
	UnitID           int64
	CompletedLessons int
	TotalLessons     int
}

// CourseProgress is the lesson and unit view behind GET /progress.
type CourseProgress struct {
	Lessons []LessonProgress
	Units   []UnitProgress
}

// AttemptMark is what the learner has already done on one exercise.
// A missing mark means there is no attempt yet.
type AttemptMark struct {
	Correct bool
}

// AssembleProgress builds lesson and unit totals from the course outline.
// A lesson is completed when every exercise has at least one attempt.
// Correct counts exercises with at least one correct attempt.
func AssembleProgress(units []Unit, lessons []Lesson, exercises []Exercise, marks map[int64]AttemptMark) CourseProgress {
	order := make([]int64, 0, len(lessons))
	byLesson := make(map[int64]*lessonTally, len(lessons))
	for _, lesson := range lessons {
		order = append(order, lesson.ID)
		byLesson[lesson.ID] = &lessonTally{}
	}
	for _, exercise := range exercises {
		tally := byLesson[exercise.LessonID]
		if tally == nil {
			continue
		}
		tally.total++
		mark, ok := marks[exercise.ID]
		if !ok {
			continue
		}
		tally.attempted++
		if mark.Correct {
			tally.correct++
		}
	}

	out := CourseProgress{
		Lessons: make([]LessonProgress, 0, len(order)),
		Units:   make([]UnitProgress, 0, len(units)),
	}
	completed := make(map[int64]bool, len(order))
	for _, id := range order {
		tally := byLesson[id]
		item := LessonProgress{
			LessonID:  id,
			Completed: tally.total > 0 && tally.attempted == tally.total,
			Correct:   tally.correct,
			Total:     tally.total,
		}
		completed[id] = item.Completed
		out.Lessons = append(out.Lessons, item)
	}
	for _, unit := range units {
		item := UnitProgress{UnitID: unit.ID}
		for _, lesson := range lessons {
			if lesson.UnitID != unit.ID {
				continue
			}
			item.TotalLessons++
			if completed[lesson.ID] {
				item.CompletedLessons++
			}
		}
		out.Units = append(out.Units, item)
	}
	return out
}

type lessonTally struct {
	total     int
	attempted int
	correct   int
}
