//go:build integration

package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLiveAuthAndReadSurfaces(t *testing.T) {
	apiURL := strings.TrimSpace(os.Getenv("CHILLY_TEST_API_URL"))
	token := strings.TrimSpace(os.Getenv("CHILLY_TEST_TOKEN"))
	if apiURL == "" || token == "" {
		t.Skip("CHILLY_TEST_API_URL and CHILLY_TEST_TOKEN are required")
	}

	configPath := filepath.Join(t.TempDir(), "config.json")

	loginOutput := runLiveCommand(t, []string{
		"--config", configPath,
		"--api-url", apiURL,
		"auth", "login",
		"--token", token,
		"--output", "json",
	})
	if loginOutput["status"] != "ok" || loginOutput["saved"] != true {
		t.Fatalf("login output = %#v", loginOutput)
	}

	whoamiOutput := runLiveCommand(t, []string{
		"--config", configPath,
		"whoami",
		"--output", "json",
	})
	username, _ := whoamiOutput["username"].(string)
	if strings.TrimSpace(username) == "" {
		t.Fatalf("whoami output = %#v", whoamiOutput)
	}

	doctorOutput := runLiveCommand(t, []string{
		"--config", configPath,
		"doctor",
		"--output", "json",
	})
	if doctorOutput["status"] != "ok" {
		t.Fatalf("doctor output = %#v", doctorOutput)
	}

	searchOutput := runLiveCommand(t, []string{
		"--config", configPath,
		"search",
		"--query", "dune",
		"--fields", "results.title",
		"--output", "json",
	})
	results, ok := searchOutput["results"].([]any)
	if !ok || len(results) == 0 {
		t.Fatalf("search output = %#v", searchOutput)
	}
}

func runLiveCommand(t *testing.T, args []string) map[string]any {
	t.Helper()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exitCode := Run(args, strings.NewReader(""), stdout, stderr)
	if exitCode != int(exitCodeSuccess) {
		t.Fatalf("Run(%v) exitCode = %d, stderr = %q", args, exitCode, stderr.String())
	}

	var output map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
		t.Fatalf("json.Unmarshal(stdout) error = %v; stdout=%q stderr=%q", err, stdout.String(), stderr.String())
	}
	return output
}
