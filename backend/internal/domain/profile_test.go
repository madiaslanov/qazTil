package domain

import "testing"

func TestProfilePrefsValid(t *testing.T) {
	ok := DefaultProfilePrefs()
	if !ok.Valid() {
		t.Fatal("defaults must be valid")
	}

	cases := []ProfilePrefs{
		{DailyGoal: 7, Level: LevelA0, Motivation: MotivationLife},
		{DailyGoal: DailyGoalMedium, Level: "B1", Motivation: MotivationLife},
		{DailyGoal: DailyGoalMedium, Level: LevelA0, Motivation: "travel"},
		{},
	}
	for _, prefs := range cases {
		if prefs.Valid() {
			t.Fatalf("expected invalid prefs %+v", prefs)
		}
	}
}
