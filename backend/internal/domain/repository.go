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
