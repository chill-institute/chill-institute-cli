package cli

import (
	"strings"
	"testing"

	"github.com/chill-institute/cli/internal/buildinfo"
)

func TestReadLineTrimsInput(t *testing.T) {
	t.Parallel()

	stdout := &strings.Builder{}
	app := &appContext{
		opts:   &appOptions{output: outputJSON},
		stdin:  strings.NewReader(" token-123 \n"),
		stdout: stdout,
		stderr: &strings.Builder{},
	}

	line, err := app.readLine("token: ")
	if err != nil {
		t.Fatalf("readLine() error = %v", err)
	}
	if line != "token-123" {
		t.Fatalf("line = %q", line)
	}
	if stdout.String() != "token: " {
		t.Fatalf("prompt = %q", stdout.String())
	}
}

func TestOpenBrowserRejectsEmptyURL(t *testing.T) {
	t.Parallel()

	if err := openBrowser(" "); err == nil {
		t.Fatal("expected error")
	}
}

func TestActiveProfileUsesDevDefaultForDevBuilds(t *testing.T) {
	t.Parallel()

	original := currentBuildInfo
	currentBuildInfo = func() buildinfo.Info {
		return buildinfo.Info{Version: "dev", Commit: "test", BuildDate: "test"}
	}
	t.Cleanup(func() { currentBuildInfo = original })

	app := &appContext{opts: &appOptions{}}
	if profile := app.activeProfile(); profile != "dev" {
		t.Fatalf("profile = %q, want %q", profile, "dev")
	}
}

func TestActiveProfileUsesExplicitProfile(t *testing.T) {
	t.Parallel()

	app := &appContext{opts: &appOptions{profile: "staging"}}
	if profile := app.activeProfile(); profile != "staging" {
		t.Fatalf("profile = %q, want %q", profile, "staging")
	}
}
