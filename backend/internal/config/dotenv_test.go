package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDotEnvFromParent(t *testing.T) {
	root := t.TempDir()
	nested := filepath.Join(root, "backend")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	body := "QAZTIL_DOTENV_FROM_FILE=from-file\nQAZTIL_DOTENV_KEEP=from-file\n"
	if err := os.WriteFile(filepath.Join(root, ".env"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	const fromFile = "QAZTIL_DOTENV_FROM_FILE"
	const keep = "QAZTIL_DOTENV_KEEP"
	t.Cleanup(func() {
		_ = os.Unsetenv(fromFile)
		_ = os.Unsetenv(keep)
	})
	t.Setenv(keep, "from-process")

	LoadDotEnvFrom(nested)

	if got := os.Getenv(fromFile); got != "from-file" {
		t.Fatalf("%s = %q", fromFile, got)
	}
	if got := os.Getenv(keep); got != "from-process" {
		t.Fatalf("%s = %q", keep, got)
	}
}
