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

func TestScannerSkipsGitAndScryDirs(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".git", "objects"), 0o755); err != nil {
		t.Fatalf("mkdir git: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, ".scry", "cache"), 0o755); err != nil {
		t.Fatalf("mkdir scry: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, ".git", "objects", "obj"), []byte("x"), 0o644); err != nil {
		t.Fatalf("write git file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, ".scry", "cache", "x"), []byte("x"), 0o644); err != nil {
		t.Fatalf("write scry file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "main.go"), []byte("package main"), 0o644); err != nil {
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
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}
	rel, _ := filepath.Rel(root, files[0].Path)
	if filepath.ToSlash(rel) != "main.go" {
		t.Fatalf("unexpected file: %s", rel)
	}
}

func TestScannerNewError(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, ".gitignore"), 0o755); err != nil {
		t.Fatalf("mkdir gitignore dir: %v", err)
	}
	if _, err := New(root); err == nil {
		t.Fatalf("expected error for invalid .gitignore")
	}
}

func TestScannerListFilesWithoutMatcher(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "main.go"), []byte("package main"), 0o644); err != nil {
		t.Fatalf("write main: %v", err)
	}
	scanner := &Scanner{Root: root, Matcher: nil}
	files, err := scanner.ListFiles()
	if err != nil {
		t.Fatalf("list files: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}
}

func TestScannerListFilesError(t *testing.T) {
	root := t.TempDir()
	blocked := filepath.Join(root, "blocked")
	if err := os.Mkdir(blocked, 0o000); err != nil {
		t.Fatalf("mkdir blocked: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chmod(blocked, 0o755)
	})

	scanner := &Scanner{Root: root, Matcher: nil}
	if _, err := scanner.ListFiles(); err == nil {
		t.Fatalf("expected error walking blocked dir")
	}
}
