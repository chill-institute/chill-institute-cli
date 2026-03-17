package cli

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/chill-institute/chill-institute-cli/internal/config"
)

func TestRunListTopMoviesRejectsBadFieldsAndMissingAuth(t *testing.T) {
	t.Parallel()

	app := &appContext{
		opts:   &appOptions{configPath: filepath.Join(t.TempDir(), "config.json"), output: outputJSON},
		stdin:  strings.NewReader(""),
		stdout: &strings.Builder{},
		stderr: &strings.Builder{},
	}

	if err := runListTopMovies(app, "["); err == nil {
		t.Fatal("runListTopMovies(invalid fields) error = nil, want error")
	}

	if err := app.saveConfig(config.Config{APIBaseURL: "https://api.chill.institute"}); err != nil {
		t.Fatalf("saveConfig() error = %v", err)
	}
	if err := runListTopMovies(app, ""); err == nil {
		t.Fatal("runListTopMovies(missing auth) error = nil, want error")
	}
}
