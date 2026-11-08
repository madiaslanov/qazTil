package usecase

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/madiaslanov/qazTil/internal/domain"
)

// ProfileUpdate is a partial edit of identity and prefs. Nil fields stay as-is.
type ProfileUpdate struct {
	DisplayName *string
	DailyGoal   *int
	Level       *domain.Level
	Motivation  *domain.Motivation
}

// ProfileService reads and updates the learner profile shown on /profile.
type ProfileService struct {
	users    domain.UserRepository
	profiles domain.ProfileRepository
	stats    domain.LearnerStatsRepository
	now      func() time.Time
}

func NewProfileService(
	users domain.UserRepository,
	profiles domain.ProfileRepository,
	stats domain.LearnerStatsRepository,
) *ProfileService {
	return &ProfileService{users: users, profiles: profiles, stats: stats, now: time.Now}
}

// Get assembles identity, prefs and live stats for one learner.
func (s *ProfileService) Get(ctx context.Context, userID int64) (domain.Profile, error) {
	if userID <= 0 {
		return domain.Profile{}, domain.ErrInvalid
	}
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return domain.Profile{}, err
	}
	prefs, err := s.profiles.Get(ctx, userID)
	if err != nil {
		return domain.Profile{}, err
	}
	stats, err := s.stats.Current(ctx, userID, s.now())
	if err != nil {
		return domain.Profile{}, err
	}
	return domain.Profile{
		UserID:      user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		CreatedAt:   user.CreatedAt,
		Prefs:       prefs,
		Stats:       stats,
	}, nil
}

// Update applies a partial profile edit and returns the fresh profile.
func (s *ProfileService) Update(ctx context.Context, userID int64, patch ProfileUpdate) (domain.Profile, error) {
	if userID <= 0 {
		return domain.Profile{}, domain.ErrInvalid
	}
	prefs, err := s.profiles.Get(ctx, userID)
	if err != nil {
		return domain.Profile{}, err
	}

	if patch.DailyGoal != nil {
		prefs.DailyGoal = *patch.DailyGoal
	}
	if patch.Level != nil {
		prefs.Level = *patch.Level
	}
	if patch.Motivation != nil {
		prefs.Motivation = *patch.Motivation
	}
	if !prefs.Valid() {
		return domain.Profile{}, domain.ErrInvalid
	}
	if err := s.profiles.Save(ctx, userID, prefs); err != nil {
		return domain.Profile{}, err
	}

	if patch.DisplayName != nil {
		name := strings.TrimSpace(*patch.DisplayName)
		if name == "" || utf8.RuneCountInString(name) > maxDisplayNameRune {
			return domain.Profile{}, domain.ErrInvalid
		}
		if _, err := s.users.UpdateDisplayName(ctx, userID, name); err != nil {
			return domain.Profile{}, err
		}
	}
	return s.Get(ctx, userID)
}
