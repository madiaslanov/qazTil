package httpapi

import (
	"net/http"

	"github.com/madiaslanov/qazTil/internal/domain"
)

// ListSituations godoc
// @Summary Список ситуаций
// @Tags situations
// @Produce json
// @Success 200 {array} SituationResponse
// @Failure 500 {object} ErrorResponse
// @Router /situations [get]
func (a *API) ListSituations(w http.ResponseWriter, r *http.Request) {
	items, err := a.situations.List(r.Context())
	if err != nil {
		writeErr(w, err)
		return
	}
	out := make([]SituationResponse, 0, len(items))
	for _, item := range items {
		out = append(out, toSituation(item))
	}
	writeJSON(w, http.StatusOK, out)
}

// GetSituation godoc
// @Summary Ситуация и её слова
// @Tags situations
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} SituationDetailResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /situations/{id} [get]
func (a *API) GetSituation(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, err)
		return
	}
	item, words, err := a.situations.Open(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, SituationDetailResponse{
		SituationResponse: toSituation(item),
		Words:             toSituationWords(words),
	})
}

// ListSituationWords godoc
// @Summary Слова выбранной ситуации
// @Tags situations
// @Produce json
// @Param id path int true "ID ситуации"
// @Success 200 {array} SituationWordResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /situations/{id}/words [get]
func (a *API) ListSituationWords(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, err)
		return
	}
	_, words, err := a.situations.Open(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toSituationWords(words))
}

func toSituation(item domain.Situation) SituationResponse {
	return SituationResponse{
		ID:          item.ID,
		Slug:        item.Slug,
		TitleKK:     item.TitleKK,
		TitleRU:     item.TitleRU,
		Description: item.Description,
		Position:    item.Position,
		WordCount:   item.WordCount,
	}
}

func toSituationWords(items []domain.SituationWord) []SituationWordResponse {
	out := make([]SituationWordResponse, 0, len(items))
	for _, item := range items {
		out = append(out, SituationWordResponse{
			ID:            item.ID,
			SituationID:   item.SituationID,
			Kazakh:        item.Kazakh,
			Russian:       item.Russian,
			Transcription: item.Transcription,
			Position:      item.Position,
		})
	}
	return out
}
