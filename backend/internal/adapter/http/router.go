package httpapi

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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
	origins    []string
}

func NewAPI(
	categories *usecase.CategoryService,
	words *usecase.WordService,
	quizzes *usecase.QuizService,
	progress *usecase.ProgressService,
	webDir string,
	origins []string,
) *API {
	return &API{
		categories: categories,
		words:      words,
		quizzes:    quizzes,
		progress:   progress,
		webDir:     webDir,
		origins:    origins,
	}
}

// NewHandler wires JSON routes, Swagger UI, the static page and a readiness check.
func NewHandler(api *API, ready func(context.Context) error, webDir string) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// Фронтенд живёт на своём домене (Vercel), поэтому браузеру нужен CORS.
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   api.origins,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/health", healthHandler(ready))

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

	files := http.FileServer(http.Dir(webDir))
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

func healthHandler(ready func(context.Context) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := ready(ctx); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, ErrorResponse{Error: "база недоступна"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}
