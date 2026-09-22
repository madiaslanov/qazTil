package httpapi

import (
	"net/http"

	"github.com/madiaslanov/qazTil/internal/domain"
)

// StartQuiz godoc
// @Summary Начать квиз
// @Description Русское слово и четыре казахских варианта. Верный ответ скрыт, пока ученик не ответит.
// @Tags quizzes
// @Accept json
// @Produce json
// @Param body body StartQuizRequest true "Параметры"
// @Success 201 {object} QuizResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /quizzes [post]
func (a *API) StartQuiz(w http.ResponseWriter, r *http.Request) {
	var req StartQuizRequest
	if err := decode(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	quiz, err := a.quizzes.Start(r.Context(), req.CategoryID, req.Size)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toQuiz(quiz))
}

// GetQuiz godoc
// @Summary Состояние квиза и счёт
// @Tags quizzes
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} QuizResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /quizzes/{id} [get]
func (a *API) GetQuiz(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, err)
		return
	}
	quiz, err := a.quizzes.Get(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toQuiz(quiz))
}

// AnswerQuiz godoc
// @Summary Ответить на вопрос
// @Description Записывает выбор и обновляет прогресс категории.
// @Tags quizzes
// @Accept json
// @Produce json
// @Param id path int true "ID квиза"
// @Param body body AnswerRequest true "Ответ"
// @Success 200 {object} AnswerResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /quizzes/{id}/answers [post]
func (a *API) AnswerQuiz(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, err)
		return
	}
	var req AnswerRequest
	if err := decode(w, r, &req); err != nil {
		writeErr(w, err)
		return
	}
	quiz, err := a.quizzes.Answer(r.Context(), id, req.QuestionID, req.SelectedIndex)
	if err != nil {
		writeErr(w, err)
		return
	}
	var (
		correct      bool
		correctIndex int
	)
	for _, question := range quiz.Questions {
		if question.ID != req.QuestionID || question.Correct == nil {
			continue
		}
		correct = *question.Correct
		correctIndex = question.CorrectIndex
		break
	}
	writeJSON(w, http.StatusOK, AnswerResponse{
		Correct:      correct,
		CorrectIndex: correctIndex,
		Quiz:         toQuiz(quiz),
	})
}

func toQuiz(quiz domain.Quiz) QuizResponse {
	correct, total := quiz.Score()
	questions := make([]QuestionResponse, 0, len(quiz.Questions))
	for _, question := range quiz.Questions {
		questions = append(questions, toQuestion(question))
	}
	return QuizResponse{
		ID:         quiz.ID,
		CategoryID: quiz.CategoryID,
		CreatedAt:  quiz.CreatedAt,
		FinishedAt: quiz.FinishedAt,
		Score:      ScoreResponse{Correct: correct, Total: total},
		Questions:  questions,
	}
}

func toQuestion(question domain.QuizQuestion) QuestionResponse {
	resp := QuestionResponse{
		ID:            question.ID,
		Prompt:        question.Prompt,
		Options:       question.Options,
		SelectedIndex: question.SelectedIndex,
		Correct:       question.Correct,
	}
	if question.SelectedIndex != nil {
		idx := question.CorrectIndex
		resp.CorrectIndex = &idx
	}
	return resp
}
