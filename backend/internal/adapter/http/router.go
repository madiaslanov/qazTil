package httpapi

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "github.com/madiaslanov/qazTil/internal/adapter/http/docs"
	"github.com/madiaslanov/qazTil/internal/usecase"
)

// API is the HTTP adapter over the use cases.
type API struct {
	categories *usecase.CategoryService
	words      *usecase.WordService
	quizzes    *usecase.QuizService
	progress   *usecase.ProgressService
	webDir     string
}

func NewAPI(
	categories *usecase.CategoryService,
	words *usecase.WordService,
	quizzes *usecase.QuizService,
	progress *usecase.ProgressService,
	webDir string,
) *API {
	return &API{
		categories: categories,
		words:      words,
		quizzes:    quizzes,
		progress:   progress,
		webDir:     webDir,
	}
}

// NewHandler wires JSON routes, Swagger UI and the static frontend.
func NewHandler(api *API) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/categories", api.ListCategories)
		r.Post("/categories", api.CreateCategory)
		r.Get("/categories/{id}", api.GetCategory)
		r.Put("/categories/{id}", api.UpdateCategory)
		r.Delete("/categories/{id}", api.DeleteCategory)

		r.Get("/words", api.ListWords)
		r.Post("/words", api.CreateWord)
		r.Get("/words/{id}", api.GetWord)
		r.Put("/words/{id}", api.UpdateWord)
		r.Delete("/words/{id}", api.DeleteWord)

		r.Post("/quizzes", api.StartQuiz)
		r.Get("/quizzes/{id}", api.GetQuiz)
		r.Post("/quizzes/{id}/answers", api.AnswerQuiz)

		r.Get("/progress", api.ListProgress)
	})

	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))

	files := http.FileServer(http.Dir(api.webDir))
	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		if strings.HasPrefix(req.URL.Path, "/api/") {
			writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "не найдено"})
			return
		}
		// Hide io.ReaderFrom. chi's writer makes FileServer stop after the first 512 bytes.
		files.ServeHTTP(plainWriter{w}, req)
	})
	return r
}
