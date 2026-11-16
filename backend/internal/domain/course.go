package domain

import "encoding/json"

const (
	// ExerciseChoice asks the learner to pick the Kazakh line from options.
	ExerciseChoice = "choice"
	// ExerciseTranslate asks the learner to produce the Kazakh line.
	ExerciseTranslate = "translate"
	// ExerciseMatch asks the learner to connect Kazakh lines with their meaning.
	ExerciseMatch = "match"
	// ExerciseAssemble asks the learner to build a line out of word tiles.
	ExerciseAssemble = "assemble"
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
	// Answer is the key: a plain line for choice and translate, JSON for the
	// kinds whose answer holds more than one value.
	Answer string
	// Payload is everything the client may see, shaped by Kind.
	Payload json.RawMessage
	// Options is the older flat list, kept for the choice exercises seeded
	// before payload existed.
	Options  []string
	Position int
}
