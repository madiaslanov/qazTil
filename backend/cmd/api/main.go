package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	httpapi "github.com/madiaslanov/qazTil/internal/adapter/http"
	"github.com/madiaslanov/qazTil/internal/adapter/sqlite"
	"github.com/madiaslanov/qazTil/internal/usecase"
)

// @title qazTil API
// @version 1.0
// @description API для изучения казахского языка: категории, словарь, квиз и прогресс.
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
func main() {
	cfg := loadConfig()
	if _, err := os.Stat(cfg.webDir); err != nil {
		log.Fatalf("web dir: %v", err)
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

	categories := sqlite.NewCategoryRepo(db)
	words := sqlite.NewWordRepo(db)
	quizzes := sqlite.NewQuizRepo(db)
	progress := sqlite.NewProgressRepo(db)

	api := httpapi.NewAPI(
		usecase.NewCategoryService(categories),
		usecase.NewWordService(categories, words),
		usecase.NewQuizService(categories, words, quizzes, progress),
		usecase.NewProgressService(categories, progress),
		cfg.webDir,
	)

	srv := &http.Server{
		Addr:              cfg.addr,
		Handler:           httpapi.NewHandler(api),
		ReadHeaderTimeout: 5 * time.Second,
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

type config struct {
	addr   string
	dbPath string
	webDir string
}

func loadConfig() config {
	return config{
		addr:   listenAddr(),
		dbPath: env("DB_PATH", "data/qaztil.db"),
		webDir: env("WEB_DIR", "../frontend"),
	}
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
