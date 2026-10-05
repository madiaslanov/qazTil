package domain

const (
	// ExerciseChoice asks the learner to pick the Kazakh line from options.
	ExerciseChoice = "choice"
	// ExerciseTranslate asks the learner to produce the Kazakh line.
	ExerciseTranslate = "translate"
)

// Unit is a course block built around one situation.
type Unit struct {
	ID          int64
	SituationID int64
	TitleKK     string
	TitleRU     string
	Description string
	Position    int
	LessonCount int
}

// Lesson is an ordered step inside a unit.
type Lesson struct {
	ID          int64
	UnitID      int64
	TitleKK     string
	TitleRU     string
	Description string
	Position    int
}

// Exercise is one prompt inside a lesson. Answer stays on the server until a check.
type Exercise struct {
	ID       int64
	LessonID int64
	Kind     string
	Prompt   string
	Answer   string
	Options  []string
	Position int
}
