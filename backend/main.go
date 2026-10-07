package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	httpapi "github.com/madiaslanov/qazTil/internal/adapter/http"
	"github.com/madiaslanov/qazTil/internal/adapter/security"
	"github.com/madiaslanov/qazTil/internal/adapter/sqlite"
	"github.com/madiaslanov/qazTil/internal/usecase"
)

// @title qazTil API
// @version 1.0
// @description API для изучения казахского языка: категории, словарь, ситуации, юниты, уроки, упражнения, квиз и прогресс.
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Токен из /auth/login в виде "Bearer <token>"
func main() {
	cfg := loadConfig()
	if len(cfg.jwtSecret) < minSecretLen {
		log.Fatalf("JWT_SECRET must be at least %d characters", minSecretLen)
	}
	if cfg.webDir != "" {
		if _, err := os.Stat(cfg.webDir); err != nil {
			log.Fatalf("web dir: %v", err)
		}
	}

	db, err := sqlite.Open(cfg.dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	if err := sqlite.Migrate(ctx, db); err != nil {
		log.Fatal(err)
	}
	seeded, err := sqlite.Seed(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
	if seeded {
		log.Print("seeded starter vocabulary")
	}
	courseSeeded, err := sqlite.SeedCourse(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
	if courseSeeded {
		log.Print("seeded situations, units, lessons and exercises")
	}

	categories := sqlite.NewCategoryRepo(db)
	words := sqlite.NewWordRepo(db)
	quizzes := sqlite.NewQuizRepo(db)
	progress := sqlite.NewProgressRepo(db)
	situations := sqlite.NewSituationRepo(db)
	course := sqlite.NewCourseRepo(db)
	users := sqlite.NewUserRepo(db)

	tokens := security.NewJWTIssuer(cfg.jwtSecret, cfg.jwtTTL)
	hasher := security.NewBcryptHasher(cfg.bcryptCost)

	api := httpapi.NewAPI(
		usecase.NewCategoryService(categories),
		usecase.NewWordService(categories, words),
		usecase.NewQuizService(categories, words, quizzes, progress),
		usecase.NewProgressService(categories, progress),
		usecase.NewSituationService(situations),
		usecase.NewCourseService(course),
		usecase.NewAuthService(users, hasher, tokens),
		tokens,
		cfg.webDir,
		cfg.origins,
	)

	srv := &http.Server{
		Addr:              cfg.addr,
		Handler:           httpapi.NewHandler(api, db.PingContext, cfg.webDir),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("qazTil listening on %s", cfg.addr)
		log.Printf("swagger %s/swagger/index.html", cfg.addr)
		errCh <- srv.ListenAndServe()
	}()

	stop, cancelSignals := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancelSignals()

	select {
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	case <-stop.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Fatal(err)
		}
	}
}

// minSecretLen keeps a throwaway JWT_SECRET out of a deployed service.
const minSecretLen = 32

type config struct {
	addr       string
	dbPath     string
	webDir     string
	origins    []string
	jwtSecret  string
	jwtTTL     time.Duration
	bcryptCost int
}

func loadConfig() config {
	return config{
		addr:       listenAddr(),
		dbPath:     env("DB_PATH", "data/qaztil.db"),
		webDir:     env("WEB_DIR", ""),
		origins:    allowedOrigins(),
		jwtSecret:  os.Getenv("JWT_SECRET"),
		jwtTTL:     durationEnv("JWT_TTL", security.DefaultTokenTTL),
		bcryptCost: intEnv("BCRYPT_COST", 0),
	}
}

func intEnv(key string, fallback int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		log.Printf("ignoring %s=%q: not a number", key, raw)
		return fallback
	}
	return value
}

func durationEnv(key string, fallback time.Duration) time.Duration {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	value, err := time.ParseDuration(raw)
	if err != nil || value <= 0 {
		log.Printf("ignoring %s=%q: not a positive duration", key, raw)
		return fallback
	}
	return value
}

// allowedOrigins задаётся через ALLOWED_ORIGINS списком через запятую.
// По умолчанию открыты локальный фронтенд и страница на Vercel.
func allowedOrigins() []string {
	raw := os.Getenv("ALLOWED_ORIGINS")
	if raw == "" {
		return []string{
			"http://localhost:3000",
			"https://qaz-til.vercel.app",
		}
	}
	origins := make([]string, 0, 4)
	for _, origin := range strings.Split(raw, ",") {
		if trimmed := strings.TrimSpace(origin); trimmed != "" {
			origins = append(origins, trimmed)
		}
	}
	return origins
}

func listenAddr() string {
	if addr := os.Getenv("ADDR"); addr != "" {
		return addr
	}
	if port := os.Getenv("PORT"); port != "" {
		if strings.HasPrefix(port, ":") {
			return port
		}
		return ":" + port
	}
	return ":8080"
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
