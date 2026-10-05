package domain

import (
	"context"
	"time"
)

// CategoryRepository stores vocabulary topics.
type CategoryRepository interface {
	List(ctx context.Context) ([]Category, error)
	Get(ctx context.Context, id int64) (Category, error)
	Create(ctx context.Context, category Category) (Category, error)
	Update(ctx context.Context, category Category) (Category, error)
	Delete(ctx context.Context, id int64) error
}

// WordRepository stores Kazakh words.
type WordRepository interface {
	List(ctx context.Context, categoryID *int64) ([]Word, error)
	Get(ctx context.Context, id int64) (Word, error)
	Create(ctx context.Context, word Word) (Word, error)
	Update(ctx context.Context, word Word) (Word, error)
	Delete(ctx context.Context, id int64) error
}

// QuizRepository stores practice sessions.
type QuizRepository interface {
	Create(ctx context.Context, quiz Quiz) (Quiz, error)
	Get(ctx context.Context, id int64) (Quiz, error)
	Save(ctx context.Context, quiz Quiz) error
}

// ProgressRepository stores answer totals for the single local learner.
type ProgressRepository interface {
	List(ctx context.Context) ([]Progress, error)
	Record(ctx context.Context, categoryID int64, correct bool, at time.Time) error
}

// SituationRepository stores real-life contexts and the words that belong to them.
type SituationRepository interface {
	List(ctx context.Context) ([]Situation, error)
	Get(ctx context.Context, id int64) (Situation, error)
	ListWords(ctx context.Context, situationID int64) ([]SituationWord, error)
}

// CourseRepository stores units, lessons and exercises.
type CourseRepository interface {
	ListUnits(ctx context.Context) ([]Unit, error)
	GetUnit(ctx context.Context, id int64) (Unit, error)
	ListLessons(ctx context.Context, unitID int64) ([]Lesson, error)
	GetLesson(ctx context.Context, id int64) (Lesson, error)
	ListExercises(ctx context.Context, lessonID int64) ([]Exercise, error)
	GetExercise(ctx context.Context, id int64) (Exercise, error)
}
