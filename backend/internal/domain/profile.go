package domain

import "time"

// Daily goals in minutes — same values the onboarding picker offers.
const (
	DailyGoalShort  = 5
	DailyGoalMedium = 10
	DailyGoalLong   = 15
)

// Level is the starting CEFR-ish band the learner picks on onboarding.
type Level string

const (
	LevelA0 Level = "A0"
	LevelA1 Level = "A1"
	LevelA2 Level = "A2"
)

// Motivation is why the learner wants Kazakh; used to bias topics later.
type Motivation string

const (
	MotivationWork  Motivation = "work"
	MotivationStudy Motivation = "study"
	MotivationLife  Motivation = "life"
)

// ProfilePrefs are the settings the learner chooses; not gamification.
type ProfilePrefs struct {
	DailyGoal  int
	Level      Level
	Motivation Motivation
}

// DefaultProfilePrefs is what a new account gets before onboarding is saved.
func DefaultProfilePrefs() ProfilePrefs {
	return ProfilePrefs{
		DailyGoal:  DailyGoalMedium,
		Level:      LevelA0,
		Motivation: MotivationLife,
	}
}

// Valid reports whether every preference is one of the allowed values.
func (p ProfilePrefs) Valid() bool {
	switch p.DailyGoal {
	case DailyGoalShort, DailyGoalMedium, DailyGoalLong:
	default:
		return false
	}
	switch p.Level {
	case LevelA0, LevelA1, LevelA2:
	default:
		return false
	}
	switch p.Motivation {
	case MotivationWork, MotivationStudy, MotivationLife:
	default:
		return false
	}
	return true
}

// Profile is the learner as shown on /profile: identity, prefs and stats.
type Profile struct {
	UserID      int64
	Email       string
	DisplayName string
	CreatedAt   time.Time
	Prefs       ProfilePrefs
	Stats       LearnerStats
}
