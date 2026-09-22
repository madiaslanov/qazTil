package domain

import "time"

// Progress is how many quiz answers a learner has given in one category.
type Progress struct {
	CategoryID     int64
	CategoryNameKK string
	CategoryNameRU string
	TotalAnswers   int
	CorrectAnswers int
	LastStudiedAt  *time.Time
}
