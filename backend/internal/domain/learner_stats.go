package domain

import "time"

const (
	// MaxLives is the full life reserve.
	MaxLives = 5
	// LifeRecovery is how long one life takes to return.
	LifeRecovery = 30 * time.Minute
	// XPPerFirstCorrect is awarded only for a correct first attempt.
	XPPerFirstCorrect = 5
)

// LearnerStats is the server-side reserve of XP, lives and streak.
type LearnerStats struct {
	UserID       int64
	XP           int
	Lives        int
	MaxLives     int
	NextLifeAt   *time.Time
	Streak       int
	LastLessonOn *time.Time
}

// DefaultStats is a new learner with a full reserve and no streak.
func DefaultStats(userID int64) LearnerStats {
	return LearnerStats{UserID: userID, Lives: MaxLives, MaxLives: MaxLives}
}

// Regenerated applies lives that have come due by now.
// The timer starts at the first lost life and then steps forward by LifeRecovery.
func (s LearnerStats) Regenerated(now time.Time) LearnerStats {
	s.MaxLives = MaxLives
	if s.NextLifeAt == nil || now.Before(*s.NextLifeAt) {
		return s
	}
	elapsed := now.Sub(*s.NextLifeAt)
	earned := 1 + int(elapsed/LifeRecovery)
	s.Lives += earned
	if s.Lives >= MaxLives {
		s.Lives = MaxLives
		s.NextLifeAt = nil
		return s
	}
	next := s.NextLifeAt.Add(time.Duration(earned) * LifeRecovery)
	s.NextLifeAt = &next
	return s
}

// LoseLife spends one life. The recovery timer starts once and is not moved
// by later losses.
func (s LearnerStats) LoseLife(now time.Time) LearnerStats {
	s.MaxLives = MaxLives
	if s.Lives > 0 {
		s.Lives--
	}
	if s.NextLifeAt == nil {
		next := now.UTC().Add(LifeRecovery)
		s.NextLifeAt = &next
	}
	return s
}

// GainXP adds the reward for a correct first attempt.
func (s LearnerStats) GainXP() LearnerStats {
	s.XP += XPPerFirstCorrect
	s.MaxLives = MaxLives
	return s
}

// CompleteLesson records that a lesson became complete on this day.
// The same day does not add another streak point. A missed day starts again at 1.
func (s LearnerStats) CompleteLesson(now time.Time) LearnerStats {
	s.MaxLives = MaxLives
	today := dateOnly(now.UTC())
	switch {
	case s.LastLessonOn != nil && sameDate(*s.LastLessonOn, today):
	case s.LastLessonOn != nil && sameDate(*s.LastLessonOn, today.AddDate(0, 0, -1)):
		s.Streak++
	default:
		s.Streak = 1
	}
	s.LastLessonOn = &today
	return s
}

func dateOnly(value time.Time) time.Time {
	year, month, day := value.UTC().Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func sameDate(left, right time.Time) bool {
	return dateOnly(left).Equal(dateOnly(right))
}
