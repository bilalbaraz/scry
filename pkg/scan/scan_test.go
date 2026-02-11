package scan

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScannerRespectsIgnore(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte("vendor/\n*.md\n"), 0o644); err != nil {
		t.Fatalf("write gitignore: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, ".scryignore"), []byte("notes.md\n"), 0o644); err != nil {
		t.Fatalf("write scryignore: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, "vendor", "mod"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, "src"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "vendor", "mod", "x.go"), []byte("package mod"), 0o644); err != nil {
		t.Fatalf("write vendor file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "notes.md"), []byte("notes"), 0o644); err != nil {
		t.Fatalf("write notes: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("readme"), 0o644); err != nil {
		t.Fatalf("write readme: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "src", "main.go"), []byte("package main"), 0o644); err != nil {
		t.Fatalf("write main: %v", err)
	}

	scanner, err := New(root)
	if err != nil {
		t.Fatalf("new scanner: %v", err)
	}
	files, err := scanner.ListFiles()
	if err != nil {
		t.Fatalf("list files: %v", err)
	}
	for _, f := range files {
		rel, _ := filepath.Rel(root, f.Path)
		if filepath.ToSlash(rel) == "src/main.go" {
			return
		}
		if filepath.ToSlash(rel) == "README.md" || filepath.ToSlash(rel) == "notes.md" || filepath.ToSlash(rel) == "vendor/mod/x.go" {
			t.Fatalf("unexpected file listed: %s", rel)
		}
	}
	t.Fatalf("expected src/main.go to be included")
}
