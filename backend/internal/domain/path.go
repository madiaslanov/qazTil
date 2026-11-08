package domain

const (
	// PathCompleted means every exercise in the lesson has an attempt.
	PathCompleted = "completed"
	// PathCurrent is the first lesson that is not completed yet.
	PathCurrent = "current"
	// PathLocked means an earlier lesson is still open.
	PathLocked = "locked"
)

// PathLesson is one node on the learning path.
type PathLesson struct {
	Lesson
	State string
}

// PathUnit is a module on the learning path with its ordered lessons.
type PathUnit struct {
	Unit
	Lessons []PathLesson
}

// LearningPath is the full skill-tree outline for one learner.
type LearningPath struct {
	Units           []PathUnit
	CurrentLessonID *int64
}

// AssemblePath builds units with lesson states from the course outline.
// Lessons stay locked until every earlier lesson is completed.
func AssemblePath(units []Unit, lessons []Lesson, completed map[int64]bool) LearningPath {
	if completed == nil {
		completed = map[int64]bool{}
	}
	ordered := make([]Lesson, 0, len(lessons))
	for _, unit := range units {
		for _, lesson := range lessons {
			if lesson.UnitID == unit.ID {
				ordered = append(ordered, lesson)
			}
		}
	}

	var currentID *int64
	for i := range ordered {
		id := ordered[i].ID
		if completed[id] {
			continue
		}
		currentID = &id
		break
	}

	passedCurrent := false
	states := make(map[int64]string, len(ordered))
	for _, lesson := range ordered {
		switch {
		case completed[lesson.ID]:
			states[lesson.ID] = PathCompleted
		case currentID != nil && lesson.ID == *currentID:
			states[lesson.ID] = PathCurrent
			passedCurrent = true
		case passedCurrent || currentID != nil:
			states[lesson.ID] = PathLocked
		default:
			// Every lesson completed: keep them marked completed above.
			states[lesson.ID] = PathLocked
		}
	}

	out := LearningPath{
		Units:           make([]PathUnit, 0, len(units)),
		CurrentLessonID: currentID,
	}
	for _, unit := range units {
		item := PathUnit{Unit: unit, Lessons: []PathLesson{}}
		for _, lesson := range ordered {
			if lesson.UnitID != unit.ID {
				continue
			}
			item.Lessons = append(item.Lessons, PathLesson{
				Lesson: lesson,
				State:  states[lesson.ID],
			})
		}
		item.LessonCount = len(item.Lessons)
		out.Units = append(out.Units, item)
	}
	return out
}
