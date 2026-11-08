package postgres

import (
	"context"
	"database/sql"

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
	err := r.db.QueryRowContext(ctx, `
		SELECT daily_goal, level, motivation
		FROM learner_profiles
		WHERE user_id = $1`, userID).Scan(&prefs.DailyGoal, &level, &motivation)
	if err != nil {
		return domain.ProfilePrefs{}, mapErr(err)
	}
	prefs.Level = domain.Level(level)
	prefs.Motivation = domain.Motivation(motivation)
	return prefs, nil
}

func (r *ProfileRepo) Save(ctx context.Context, userID int64, prefs domain.ProfilePrefs) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO learner_profiles (user_id, daily_goal, level, motivation)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE
		SET daily_goal = EXCLUDED.daily_goal,
		    level = EXCLUDED.level,
		    motivation = EXCLUDED.motivation`,
		userID, prefs.DailyGoal, string(prefs.Level), string(prefs.Motivation))
	return err
}
