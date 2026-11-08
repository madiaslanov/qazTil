package httpapi

import (
	"net/http"

	"github.com/madiaslanov/qazTil/internal/domain"
	"github.com/madiaslanov/qazTil/internal/usecase"
)

// GetProfile godoc
// @Summary Профиль ученика
// @Tags profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ProfileResponse
// @Failure 401 {object} ErrorResponse
// @Router /me/profile [get]
func (a *API) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserID(r.Context())
	if !ok {
		writeErr(w, domain.ErrUnauthorized)
		return
	}
	profile, err := a.profile.Get(r.Context(), userID)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toProfile(profile))
}

// UpdateProfile godoc
// @Summary Обновить профиль
// @Tags profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body ProfileUpdateRequest true "Поля для изменения"
// @Success 200 {object} ProfileResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /me/profile [put]
func (a *API) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserID(r.Context())
	if !ok {
		writeErr(w, domain.ErrUnauthorized)
		return
	}
	var req ProfileUpdateRequest
	if err := decode(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	profile, err := a.profile.Update(r.Context(), userID, toProfileUpdate(req))
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toProfile(profile))
}

func toProfile(profile domain.Profile) ProfileResponse {
	stats := toStats(profile.Stats)
	return ProfileResponse{
		ID:           profile.UserID,
		Email:        profile.Email,
		DisplayName:  profile.DisplayName,
		CreatedAt:    profile.CreatedAt,
		DailyGoal:    profile.Prefs.DailyGoal,
		Level:        string(profile.Prefs.Level),
		Motivation:   string(profile.Prefs.Motivation),
		XP:           stats.XP,
		Lives:        stats.Lives,
		MaxLives:     stats.MaxLives,
		NextLifeAt:   stats.NextLifeAt,
		Streak:       stats.Streak,
		LastLessonOn: stats.LastLessonOn,
	}
}

func toProfileUpdate(req ProfileUpdateRequest) usecase.ProfileUpdate {
	var patch usecase.ProfileUpdate
	if req.DisplayName != nil {
		patch.DisplayName = req.DisplayName
	}
	if req.DailyGoal != nil {
		patch.DailyGoal = req.DailyGoal
	}
	if req.Level != nil {
		level := domain.Level(*req.Level)
		patch.Level = &level
	}
	if req.Motivation != nil {
		motivation := domain.Motivation(*req.Motivation)
		patch.Motivation = &motivation
	}
	return patch
}

func toPrefs(req RegisterRequest) domain.ProfilePrefs {
	if req.DailyGoal == nil && req.Level == nil && req.Motivation == nil {
		return domain.ProfilePrefs{}
	}
	prefs := domain.DefaultProfilePrefs()
	if req.DailyGoal != nil {
		prefs.DailyGoal = *req.DailyGoal
	}
	if req.Level != nil {
		prefs.Level = domain.Level(*req.Level)
	}
	if req.Motivation != nil {
		prefs.Motivation = domain.Motivation(*req.Motivation)
	}
	return prefs
}
