package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type artifact struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	GoOS   string `json:"goos"`
	GoArch string `json:"goarch"`
	Type   string `json:"type"`
}

type metadata struct {
	Version string `json:"version"`
}

func resolveVersion(cfg options) (string, error) {
	version := strings.TrimPrefix(strings.TrimSpace(cfg.version), "v")
	if version != "" {
		return version, nil
	}

	data, err := os.ReadFile(filepath.Join(cfg.distDir, "metadata.json"))
	if err != nil {
		return "", fmt.Errorf("read metadata version: %w", err)
	}
	var meta metadata
	if err := json.Unmarshal(data, &meta); err != nil {
		return "", fmt.Errorf("parse metadata version: %w", err)
	}
	version = strings.TrimPrefix(strings.TrimSpace(meta.Version), "v")
	if version == "" {
		return "", errors.New("metadata version is empty")
	}
	return version, nil
}

func readArtifacts(path string) ([]artifact, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read artifacts: %w", err)
	}
	var artifacts []artifact
	if err := json.Unmarshal(data, &artifacts); err != nil {
		return nil, fmt.Errorf("parse artifacts: %w", err)
	}
	return artifacts, nil
}

func binaryArtifacts(artifacts []artifact) map[string]string {
	binaries := map[string]string{}
	for _, item := range artifacts {
		if item.Type != "Binary" || item.GoOS == "" || item.GoArch == "" || item.Path == "" {
			continue
		}
		binaries[item.GoOS+"/"+item.GoArch] = item.Path
	}
	return binaries
}
