package httpapi

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/madiaslanov/qazTil/internal/domain"
)

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("write json: %v", err)
	}
}

func writeErr(w http.ResponseWriter, err error) {
	status, message := http.StatusInternalServerError, "внутренняя ошибка"
	switch {
	case errors.Is(err, domain.ErrNotFound):
		status, message = http.StatusNotFound, "не найдено"
	case errors.Is(err, domain.ErrInvalid):
		status, message = http.StatusBadRequest, "некорректные данные"
	case errors.Is(err, domain.ErrNotEnoughWords):
		status, message = http.StatusBadRequest, "в категории меньше четырёх разных слов"
	case errors.Is(err, domain.ErrAlreadyAnswered):
		status, message = http.StatusConflict, "на вопрос уже ответили"
	default:
		log.Printf("api error: %v", err)
	}
	writeJSON(w, status, ErrorResponse{Error: message})
}

func decode(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return domain.ErrInvalid
	}
	return nil
}

// plainWriter keeps FileServer on the normal copy path.
type plainWriter struct {
	http.ResponseWriter
}

func pathID(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, domain.ErrInvalid
	}
	return id, nil
}
