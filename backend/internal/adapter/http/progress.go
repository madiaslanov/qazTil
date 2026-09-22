package httpapi

import (
	"net/http"

	"github.com/madiaslanov/qazTil/internal/domain"
)

// ListProgress godoc
// @Summary Прогресс по категориям
// @Tags progress
// @Produce json
// @Success 200 {array} ProgressResponse
// @Failure 500 {object} ErrorResponse
// @Router /progress [get]
func (a *API) ListProgress(w http.ResponseWriter, r *http.Request) {
	items, err := a.progress.List(r.Context())
	if err != nil {
		writeErr(w, err)
		return
	}
	out := make([]ProgressResponse, 0, len(items))
	for _, item := range items {
		out = append(out, toProgress(item))
	}
	writeJSON(w, http.StatusOK, out)
}

func toProgress(item domain.Progress) ProgressResponse {
	return ProgressResponse{
		CategoryID:     item.CategoryID,
		CategoryNameKK: item.CategoryNameKK,
		CategoryNameRU: item.CategoryNameRU,
		TotalAnswers:   item.TotalAnswers,
		CorrectAnswers: item.CorrectAnswers,
		LastStudiedAt:  item.LastStudiedAt,
	}
}
