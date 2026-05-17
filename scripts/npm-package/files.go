package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func resetOutputDir(path string) error {
	clean := filepath.Clean(path)
	if clean == "." || clean == string(filepath.Separator) {
		return fmt.Errorf("refusing to clean unsafe output directory %q", path)
	}
	if err := os.RemoveAll(clean); err != nil {
		return fmt.Errorf("clean output directory: %w", err)
	}
	if err := os.MkdirAll(clean, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	return nil
}

func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal %s: %w", path, err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func copyFile(source string, destination string, mode os.FileMode) error {
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer func() {
		_ = in.Close()
	}()

	out, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(out, in)
	closeErr := out.Close()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}
