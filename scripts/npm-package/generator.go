package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func run(cfg options) error {
	if cfg.distDir == "" {
		return errors.New("dist directory is required")
	}
	if cfg.outDir == "" {
		return errors.New("output directory is required")
	}

	version, err := resolveVersion(cfg)
	if err != nil {
		return err
	}

	artifacts, err := readArtifacts(filepath.Join(cfg.distDir, "artifacts.json"))
	if err != nil {
		return err
	}

	binaries := binaryArtifacts(artifacts)
	if err := resetOutputDir(cfg.outDir); err != nil {
		return err
	}

	if err := writeRootPackage(cfg.outDir, version); err != nil {
		return err
	}
	for _, t := range targets {
		source, ok := binaries[t.goOS+"/"+t.goArch]
		if !ok {
			return fmt.Errorf("missing GoReleaser binary artifact for %s/%s", t.goOS, t.goArch)
		}
		if err := writePlatformPackage(cfg.outDir, version, t, source); err != nil {
			return err
		}
	}

	return nil
}

func writeRootPackage(outDir string, version string) error {
	packageDir := filepath.Join(outDir, "cli")
	if err := os.MkdirAll(filepath.Join(packageDir, "bin"), 0o755); err != nil {
		return fmt.Errorf("create root package: %w", err)
	}

	if err := writeJSON(filepath.Join(packageDir, "package.json"), rootPackage(version)); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(packageDir, "README.md"), []byte(rootReadme()), 0o644); err != nil {
		return fmt.Errorf("write root README: %w", err)
	}
	if err := os.WriteFile(filepath.Join(packageDir, "bin", "chilly.js"), []byte(wrapperScript()), 0o755); err != nil {
		return fmt.Errorf("write launcher: %w", err)
	}
	return nil
}

func writePlatformPackage(outDir string, version string, t target, source string) error {
	packageDir := filepath.Join(outDir, platformDirName(t))
	binDir := filepath.Join(packageDir, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		return fmt.Errorf("create platform package %s: %w", t.suffix, err)
	}

	if err := writeJSON(filepath.Join(packageDir, "package.json"), platformPackage(version, t)); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(packageDir, "README.md"), []byte(platformReadme(t)), 0o644); err != nil {
		return fmt.Errorf("write platform README: %w", err)
	}
	if err := copyFile(source, filepath.Join(binDir, t.binaryFile), 0o755); err != nil {
		return fmt.Errorf("copy platform binary %s: %w", t.suffix, err)
	}
	return nil
}
