package httpapi

import (
	"net/http"
	"strconv"

	"github.com/madiaslanov/qazTil/internal/domain"
)

// ListWords godoc
// @Summary Список слов
// @Tags words
// @Produce json
// @Param category_id query int false "Фильтр по категории"
// @Success 200 {array} WordResponse
// @Failure 400 {object} ErrorResponse
// @Router /words [get]
func (a *API) ListWords(w http.ResponseWriter, r *http.Request) {
	var categoryID *int64
	if raw := r.URL.Query().Get("category_id"); raw != "" {
		id, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || id <= 0 {
			writeErr(w, domain.ErrInvalid)
			return
		}
		categoryID = &id
	}
	items, err := a.words.List(r.Context(), categoryID)
	if err != nil {
		writeErr(w, err)
		return
	}
	out := make([]WordResponse, 0, len(items))
	for _, item := range items {
		out = append(out, toWord(item))
	}
	writeJSON(w, http.StatusOK, out)
}

// GetWord godoc
// @Summary Слово по id
// @Tags words
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} WordResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /words/{id} [get]
func (a *API) GetWord(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, err)
		return
	}
	item, err := a.words.Get(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toWord(item))
}

// CreateWord godoc
// @Summary Создать слово
// @Tags words
// @Accept json
// @Produce json
// @Param body body WordRequest true "Слово"
// @Success 201 {object} WordResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /words [post]
func (a *API) CreateWord(w http.ResponseWriter, r *http.Request) {
	var req WordRequest
	if err := decode(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	item, err := a.words.Create(r.Context(), wordFromRequest(0, req))
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toWord(item))
}

// UpdateWord godoc
// @Summary Обновить слово
// @Tags words
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param body body WordRequest true "Слово"
// @Success 200 {object} WordResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /words/{id} [put]
func (a *API) UpdateWord(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, err)
		return
	}
	var req WordRequest
	if err := decode(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	item, err := a.words.Update(r.Context(), wordFromRequest(id, req))
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toWord(item))
}

// DeleteWord godoc
// @Summary Удалить слово
// @Tags words
// @Param id path int true "ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /words/{id} [delete]
func (a *API) DeleteWord(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, err)
		return
	}
	if err := a.words.Delete(r.Context(), id); err != nil {
		writeErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func wordFromRequest(id int64, req WordRequest) domain.Word {
	return domain.Word{
		ID:            id,
		CategoryID:    req.CategoryID,
		Kazakh:        req.Kazakh,
		Russian:       req.Russian,
		Transcription: req.Transcription,
		ExampleKK:     req.ExampleKK,
		ExampleRU:     req.ExampleRU,
	}
}

func toWord(item domain.Word) WordResponse {
	return WordResponse{
		ID:            item.ID,
		CategoryID:    item.CategoryID,
		Kazakh:        item.Kazakh,
		Russian:       item.Russian,
		Transcription: item.Transcription,
		ExampleKK:     item.ExampleKK,
		ExampleRU:     item.ExampleRU,
	}
}
