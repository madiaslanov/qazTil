package domain

import "time"

// Quiz is one practice session: a Russian prompt and Kazakh options.
type Quiz struct {
	ID         int64
	CategoryID int64
	CreatedAt  time.Time
	FinishedAt *time.Time
	Questions  []QuizQuestion
}

// QuizQuestion stores the options and, after an answer, the learner's choice.
type QuizQuestion struct {
	ID            int64
	WordID        int64
	Prompt        string
	Options       []string
	CorrectIndex  int
	SelectedIndex *int
	Correct       *bool
}

// ApplyAnswer records a choice. The correct option stays on the question
// so the HTTP layer can hide it until the learner answers.
func (q *Quiz) ApplyAnswer(questionID int64, selectedIndex int, now time.Time) error {
	if q.FinishedAt != nil {
		return ErrAlreadyAnswered
	}
	for i := range q.Questions {
		question := &q.Questions[i]
		if question.ID != questionID {
			continue
		}
		if selectedIndex < 0 || selectedIndex >= len(question.Options) {
			return ErrInvalid
		}
		if question.SelectedIndex != nil {
			return ErrAlreadyAnswered
		}
		selected := selectedIndex
		correct := selectedIndex == question.CorrectIndex
		question.SelectedIndex = &selected
		question.Correct = &correct
		q.markFinished(now)
		return nil
	}
	return ErrNotFound
}

func (q *Quiz) markFinished(now time.Time) {
	for _, question := range q.Questions {
		if question.SelectedIndex == nil {
			return
		}
	}
	finished := now.UTC()
	q.FinishedAt = &finished
}

// Score returns how many answered questions were correct and how many questions exist.
func (q Quiz) Score() (correct int, total int) {
	total = len(q.Questions)
	for _, question := range q.Questions {
		if question.Correct != nil && *question.Correct {
			correct++
		}
	}
	return correct, total
}
