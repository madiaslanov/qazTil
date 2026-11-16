package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/madiaslanov/qazTil/internal/domain"
)

type ProfileRepo struct {
	db *sql.DB
}

func NewProfileRepo(db *sql.DB) *ProfileRepo {
	return &ProfileRepo{db: db}
}

func (r *ProfileRepo) Get(ctx context.Context, userID int64) (domain.ProfilePrefs, error) {
	defaults := domain.DefaultProfilePrefs()
	if _, err := r.db.ExecContext(ctx, `
		INSERT INTO learner_profiles (user_id, daily_goal, level, motivation)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO NOTHING`,
		userID, defaults.DailyGoal, string(defaults.Level), string(defaults.Motivation)); err != nil {
		return domain.ProfilePrefs{}, err
	}

	var prefs domain.ProfilePrefs
	var level, motivation string
	var onboardedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, `
		SELECT daily_goal, level, motivation, timezone, onboarded_at
		FROM learner_profiles
		WHERE user_id = $1`, userID).
		Scan(&prefs.DailyGoal, &level, &motivation, &prefs.Timezone, &onboardedAt)
	if err != nil {
		return domain.ProfilePrefs{}, mapErr(err)
	}
	prefs.Level = domain.Level(level)
	prefs.Motivation = domain.Motivation(motivation)
	prefs.Onboarded = onboardedAt.Valid
	return prefs, nil
}

// Save stores the prefs. It stamps the moment onboarding finished only when
// the caller says the learner actually chose something: registration saves
// defaults for everyone, and that must not count as having been onboarded.
// Once stamped, the moment never moves — changing a daily goal later does not
// make the learner newly onboarded.
func (r *ProfileRepo) Save(ctx context.Context, userID int64, prefs domain.ProfilePrefs) error {
	timezone := prefs.Timezone
	if timezone == "" {
		timezone = domain.DefaultTimezone
	}
	var stamp any
	if prefs.Onboarded {
		stamp = time.Now().UTC()
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO learner_profiles (user_id, daily_goal, level, motivation, timezone, onboarded_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id) DO UPDATE
		SET daily_goal = EXCLUDED.daily_goal,
		    level = EXCLUDED.level,
		    motivation = EXCLUDED.motivation,
		    timezone = EXCLUDED.timezone,
		    onboarded_at = COALESCE(learner_profiles.onboarded_at, EXCLUDED.onboarded_at)`,
		userID, prefs.DailyGoal, string(prefs.Level), string(prefs.Motivation), timezone, stamp)
	return err
}

// Daily counts what the learner earned inside their own day. XP is awarded for
// the first correct attempt on an exercise, so the same attempts that drive the
// stats drive this count, and no extra table is needed.
func (r *ProfileRepo) Daily(ctx context.Context, userID int64, zone string, now time.Time) (domain.DailyProgress, error) {
	if zone == "" {
		zone = domain.DefaultTimezone
	}
	var first int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM (
			SELECT DISTINCT ON (exercise_id) correct, created_at
			FROM attempts
			WHERE user_id = $1
			ORDER BY exercise_id, created_at
		) AS firsts
		WHERE firsts.correct
		  AND (firsts.created_at AT TIME ZONE $2)::date = ($3 AT TIME ZONE $2)::date`,
		userID, zone, now.UTC()).Scan(&first)
	if err != nil {
		return domain.DailyProgress{}, err
	}
	var day string
	if err := r.db.QueryRowContext(ctx,
		`SELECT to_char(($1 AT TIME ZONE $2)::date, 'YYYY-MM-DD')`, now.UTC(), zone).Scan(&day); err != nil {
		return domain.DailyProgress{}, err
	}
	return domain.DailyProgress{Date: day, XPEarned: first * domain.XPPerFirstCorrect}, nil
}
