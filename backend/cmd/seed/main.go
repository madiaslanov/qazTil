package main

import (
	"context"
	"log"
	"os"

	"github.com/madiaslanov/qazTil/internal/adapter/postgres"
	"github.com/madiaslanov/qazTil/internal/config"
)

// Fills an empty database with starter vocabulary and the situational course.
// Existing rows are left in place. DATABASE_URL is required.
func main() {
	config.LoadDotEnv()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := postgres.Open(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	if err := postgres.Migrate(ctx, db); err != nil {
		log.Fatal(err)
	}
	vocab, err := postgres.Seed(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
	course, err := postgres.SeedCourse(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("vocabulary seeded=%v, course seeded=%v", vocab, course)
}
