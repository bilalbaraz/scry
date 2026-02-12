package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaultMissingOptional(t *testing.T) {
	root := t.TempDir()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer func() {
		if err := os.Chdir(cwd); err != nil {
			t.Fatalf("chdir back: %v", err)
		}
	}()
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	cfg, err := Load("", false)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if cfg.Found {
		t.Fatalf("expected not found")
	}
	if cfg.Path != ".scry.yml" {
		t.Fatalf("unexpected path: %s", cfg.Path)
	}
}

func TestLoadExplicitRequiredMissing(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "missing.yml"), true)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestLoadReadsContents(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, ".scry.yml")
	if err := os.WriteFile(path, []byte("foo: bar\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	cfg, err := Load(path, true)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if !cfg.Found || cfg.Raw == "" {
		t.Fatalf("expected found config with raw content")
	}
}
