package main

import (
	"context"
	"log"
	"os"

	"github.com/madiaslanov/qazTil/internal/adapter/sqlite"
)

// Fills an empty database with starter vocabulary and the situational course.
// Existing rows are left in place. DB_PATH defaults to data/qaztil.db.
func main() {
	path := os.Getenv("DB_PATH")
	if path == "" {
		path = "data/qaztil.db"
	}

	db, err := sqlite.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	if err := sqlite.Migrate(ctx, db); err != nil {
		log.Fatal(err)
	}
	vocab, err := sqlite.Seed(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
	course, err := sqlite.SeedCourse(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("db %s: vocabulary seeded=%v, course seeded=%v", path, vocab, course)
}
