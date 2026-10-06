package httpapi

import (
	"net/http"

	"github.com/madiaslanov/qazTil/internal/domain"
	"github.com/madiaslanov/qazTil/internal/usecase"
)

// Register godoc
// @Summary Регистрация
// @Tags auth
// @Accept json
// @Produce json
// @Param body body RegisterRequest true "Почта и пароль"
// @Success 201 {object} SessionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /auth/register [post]
func (a *API) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := decode(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	session, err := a.auth.Register(r.Context(), req.Email, req.Password, req.DisplayName)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toSession(session))
}

// Login godoc
// @Summary Вход
// @Tags auth
// @Accept json
// @Produce json
// @Param body body LoginRequest true "Почта и пароль"
// @Success 200 {object} SessionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/login [post]
func (a *API) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := decode(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	session, err := a.auth.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toSession(session))
}

// Me godoc
// @Summary Текущий пользователь
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} UserResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/me [get]
func (a *API) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserID(r.Context())
	if !ok {
		writeErr(w, domain.ErrUnauthorized)
		return
	}
	user, err := a.auth.Me(r.Context(), userID)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toUser(user))
}

func toUser(user domain.User) UserResponse {
	return UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		CreatedAt:   user.CreatedAt,
	}
}

func toSession(session usecase.Session) SessionResponse {
	return SessionResponse{
		Token:     session.Token,
		ExpiresAt: session.ExpiresAt,
		User:      toUser(session.User),
	}
}
