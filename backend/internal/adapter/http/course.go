package httpapi

import (
	"net/http"

	"github.com/madiaslanov/qazTil/internal/domain"
)

// ListUnits godoc
// @Summary Список юнитов
// @Tags course
// @Produce json
// @Success 200 {array} UnitResponse
// @Failure 500 {object} ErrorResponse
// @Router /units [get]
func (a *API) ListUnits(w http.ResponseWriter, r *http.Request) {
	items, err := a.course.ListUnits(r.Context())
	if err != nil {
		writeErr(w, err)
		return
	}
	out := make([]UnitResponse, 0, len(items))
	for _, item := range items {
		out = append(out, toUnit(item))
	}
	writeJSON(w, http.StatusOK, out)
}

// GetUnit godoc
// @Summary Юнит и его уроки
// @Tags course
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} UnitDetailResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /units/{id} [get]
func (a *API) GetUnit(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, err)
		return
	}
	item, lessons, err := a.course.OpenUnit(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, UnitDetailResponse{
		UnitResponse: toUnit(item),
		Lessons:      toLessons(lessons),
	})
}

// ListLessons godoc
// @Summary Уроки юнита
// @Tags course
// @Produce json
// @Param id path int true "ID юнита"
// @Success 200 {array} LessonResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /units/{id}/lessons [get]
func (a *API) ListLessons(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, err)
		return
	}
	lessons, err := a.course.Lessons(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toLessons(lessons))
}

// GetLesson godoc
// @Summary Урок и его упражнения
// @Description Правильный ответ в упражнение не входит. Его возвращает проверка.
// @Tags course
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} LessonDetailResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /lessons/{id} [get]
func (a *API) GetLesson(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, err)
		return
	}
	item, exercises, err := a.course.OpenLesson(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, LessonDetailResponse{
		LessonResponse: toLesson(item),
		Exercises:      toExercises(exercises),
	})
}

// ListExercises godoc
// @Summary Упражнения урока
// @Tags course
// @Produce json
// @Param id path int true "ID урока"
// @Success 200 {array} ExerciseResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /lessons/{id}/exercises [get]
func (a *API) ListExercises(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, err)
		return
	}
	exercises, err := a.course.Exercises(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toExercises(exercises))
}

// AnswerExercise godoc
// @Summary Проверить ответ на упражнение
// @Tags course
// @Accept json
// @Produce json
// @Param id path int true "ID упражнения"
// @Param body body ExerciseAnswerRequest true "Ответ"
// @Success 200 {object} ExerciseAnswerResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /exercises/{id}/answers [post]
func (a *API) AnswerExercise(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, err)
		return
	}
	var req ExerciseAnswerRequest
	if err := decode(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	result, err := a.course.Check(r.Context(), id, req.Answer)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, ExerciseAnswerResponse{
		Correct:  result.Correct,
		Expected: result.Expected,
	})
}

func toUnit(item domain.Unit) UnitResponse {
	return UnitResponse{
		ID:          item.ID,
		SituationID: item.SituationID,
		TitleKK:     item.TitleKK,
		TitleRU:     item.TitleRU,
		Description: item.Description,
		Position:    item.Position,
		LessonCount: item.LessonCount,
	}
}

func toLessons(items []domain.Lesson) []LessonResponse {
	out := make([]LessonResponse, 0, len(items))
	for _, item := range items {
		out = append(out, toLesson(item))
	}
	return out
}

func toLesson(item domain.Lesson) LessonResponse {
	return LessonResponse{
		ID:          item.ID,
		UnitID:      item.UnitID,
		TitleKK:     item.TitleKK,
		TitleRU:     item.TitleRU,
		Description: item.Description,
		Position:    item.Position,
	}
}

func toExercises(items []domain.Exercise) []ExerciseResponse {
	out := make([]ExerciseResponse, 0, len(items))
	for _, item := range items {
		options := item.Options
		if options == nil {
			options = []string{}
		}
		out = append(out, ExerciseResponse{
			ID:       item.ID,
			LessonID: item.LessonID,
			Kind:     item.Kind,
			Prompt:   item.Prompt,
			Options:  options,
			Position: item.Position,
		})
	}
	return out
}
