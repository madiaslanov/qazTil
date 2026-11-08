package httpapi

import (
	"net/http"
	"time"

	"github.com/madiaslanov/qazTil/internal/domain"
)

// ListProgress godoc
// @Summary Прогресс по урокам и модулям
// @Tags progress
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ProgressResponse
// @Failure 401 {object} ErrorResponse
// @Router /progress [get]
func (a *API) ListProgress(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserID(r.Context())
	if !ok {
		writeErr(w, domain.ErrUnauthorized)
		return
	}
	progress, err := a.progress.List(r.Context(), userID)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toCourseProgress(progress))
}

// Stats godoc
// @Summary Опыт, жизни и стрик
// @Tags progress
// @Produce json
// @Security BearerAuth
// @Success 200 {object} StatsResponse
// @Failure 401 {object} ErrorResponse
// @Router /me/stats [get]
func (a *API) Stats(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserID(r.Context())
	if !ok {
		writeErr(w, domain.ErrUnauthorized)
		return
	}
	stats, err := a.attempts.Stats(r.Context(), userID)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toStats(stats))
}

func toCourseProgress(progress domain.CourseProgress) ProgressResponse {
	out := ProgressResponse{
		Lessons: make([]LessonProgressResponse, 0, len(progress.Lessons)),
		Units:   make([]UnitProgressResponse, 0, len(progress.Units)),
	}
	for _, item := range progress.Lessons {
		out.Lessons = append(out.Lessons, LessonProgressResponse{
			LessonID:  item.LessonID,
			Completed: item.Completed,
			Correct:   item.Correct,
			Total:     item.Total,
		})
	}
	for _, item := range progress.Units {
		out.Units = append(out.Units, UnitProgressResponse{
			UnitID:           item.UnitID,
			CompletedLessons: item.CompletedLessons,
			TotalLessons:     item.TotalLessons,
		})
	}
	return out
}

func toStats(stats domain.LearnerStats) StatsResponse {
	var last *string
	if stats.LastLessonOn != nil {
		value := stats.LastLessonOn.UTC().Format(time.DateOnly)
		last = &value
	}
	return StatsResponse{
		XP:           stats.XP,
		Lives:        stats.Lives,
		MaxLives:     domain.MaxLives,
		NextLifeAt:   stats.NextLifeAt,
		Streak:       stats.Streak,
		LastLessonOn: last,
	}
}
