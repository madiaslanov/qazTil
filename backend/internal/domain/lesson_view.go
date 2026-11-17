package domain

// LessonView is one lesson as the player needs it: the lesson itself, its
// exercises, and where the learner stands inside it.
type LessonView struct {
	Lesson    Lesson
	Exercises []Exercise
	// State is one of PathCompleted, PathCurrent or PathLocked. A locked
	// lesson is never served: the caller gets ErrLocked instead.
	State string
	// ResumePosition is the 0-based index of the first unanswered exercise.
	ResumePosition int
	Correct        int
	Total          int
	Lives          int
	MaxLives       int
}
