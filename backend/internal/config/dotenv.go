package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// LoadDotEnv reads the nearest .env, walking up from the working directory.
// Variables already present in the process are left unchanged, so a platform
// secret wins over a file that happened to be deployed next to the binary.
func LoadDotEnv() {
	dir, err := os.Getwd()
	if err != nil {
		log.Printf("env file: %v", err)
		return
	}
	LoadDotEnvFrom(dir)
}

// LoadDotEnvFrom walks up from dir and loads the first .env it finds.
func LoadDotEnvFrom(dir string) {
	for range 4 {
		path := filepath.Join(dir, ".env")
		info, err := os.Stat(path)
		if err == nil && !info.IsDir() {
			if err := godotenv.Load(path); err != nil {
				log.Printf("env file %s: %v", path, err)
				return
			}
			log.Printf("loaded env %s", path)
			return
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return
		}
		dir = parent
	}
}
