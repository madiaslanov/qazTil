package httpapi

import (
	"net/http"

	"github.com/madiaslanov/qazTil/internal/domain"
)

// GetPath godoc
// @Summary Путь обучения
// @Description Юниты с уроками и статусами completed/current/locked для текущего ученика.
// @Tags course
// @Produce json
// @Security BearerAuth
// @Success 200 {object} PathResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /path [get]
func (a *API) GetPath(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserID(r.Context())
	if !ok {
		writeErr(w, domain.ErrUnauthorized)
		return
	}
	path, err := a.path.Get(r.Context(), userID)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toPath(path))
}

func toPath(path domain.LearningPath) PathResponse {
	out := PathResponse{
		Units:           make([]PathUnitResponse, 0, len(path.Units)),
		CurrentLessonID: path.CurrentLessonID,
	}
	for _, unit := range path.Units {
		item := PathUnitResponse{
			ID:          unit.ID,
			SituationID: unit.SituationID,
			TitleKK:     unit.TitleKK,
			TitleRU:     unit.TitleRU,
			Description: unit.Description,
			Position:    unit.Position,
			LessonCount: unit.LessonCount,
			Lessons:     make([]PathLessonResponse, 0, len(unit.Lessons)),
		}
		for _, lesson := range unit.Lessons {
			item.Lessons = append(item.Lessons, PathLessonResponse{
				ID:          lesson.ID,
				UnitID:      lesson.UnitID,
				TitleKK:     lesson.TitleKK,
				TitleRU:     lesson.TitleRU,
				Description: lesson.Description,
				Position:    lesson.Position,
				State:       lesson.State,
			})
		}
		out.Units = append(out.Units, item)
	}
	return out
}
