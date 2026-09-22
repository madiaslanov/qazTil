package httpapi

import (
	"net/http"

	"github.com/madiaslanov/qazTil/internal/domain"
)

// ListCategories godoc
// @Summary Список категорий
// @Tags categories
// @Produce json
// @Success 200 {array} CategoryResponse
// @Failure 500 {object} ErrorResponse
// @Router /categories [get]
func (a *API) ListCategories(w http.ResponseWriter, r *http.Request) {
	items, err := a.categories.List(r.Context())
	if err != nil {
		writeErr(w, err)
		return
	}
	out := make([]CategoryResponse, 0, len(items))
	for _, item := range items {
		out = append(out, toCategory(item))
	}
	writeJSON(w, http.StatusOK, out)
}

// GetCategory godoc
// @Summary Категория по id
// @Tags categories
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} CategoryResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /categories/{id} [get]
func (a *API) GetCategory(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, err)
		return
	}
	item, err := a.categories.Get(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toCategory(item))
}

// CreateCategory godoc
// @Summary Создать категорию
// @Tags categories
// @Accept json
// @Produce json
// @Param body body CategoryRequest true "Категория"
// @Success 201 {object} CategoryResponse
// @Failure 400 {object} ErrorResponse
// @Router /categories [post]
func (a *API) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req CategoryRequest
	if err := decode(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	item, err := a.categories.Create(r.Context(), domain.Category{
		NameKK:      req.NameKK,
		NameRU:      req.NameRU,
		Description: req.Description,
	})
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toCategory(item))
}

// UpdateCategory godoc
// @Summary Обновить категорию
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param body body CategoryRequest true "Категория"
// @Success 200 {object} CategoryResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /categories/{id} [put]
func (a *API) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, err)
		return
	}
	var req CategoryRequest
	if err := decode(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	item, err := a.categories.Update(r.Context(), domain.Category{
		ID:          id,
		NameKK:      req.NameKK,
		NameRU:      req.NameRU,
		Description: req.Description,
	})
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toCategory(item))
}

// DeleteCategory godoc
// @Summary Удалить категорию
// @Tags categories
// @Param id path int true "ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /categories/{id} [delete]
func (a *API) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, err)
		return
	}
	if err := a.categories.Delete(r.Context(), id); err != nil {
		writeErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func toCategory(item domain.Category) CategoryResponse {
	return CategoryResponse{
		ID:          item.ID,
		NameKK:      item.NameKK,
		NameRU:      item.NameRU,
		Description: item.Description,
	}
}
