package indexer

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"scry/pkg/metadata"
	"scry/pkg/workspace"
)

func requireSQLite(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skipf("sqlite3 not found: %v", err)
	}
}

func TestRunErrorsOnMissingRoot(t *testing.T) {
	requireSQLite(t)
	root := t.TempDir()
	rootFile := filepath.Join(root, "rootfile")
	if err := os.WriteFile(rootFile, []byte("not a dir"), 0o644); err != nil {
		t.Fatalf("write root file: %v", err)
	}
	_, err := Run(Options{Root: rootFile}, func(Progress) {})
	if err == nil {
		t.Fatalf("expected error for missing root")
	}
}

func TestRunEmptyRepo(t *testing.T) {
	requireSQLite(t)
	root := t.TempDir()
	var stages []string
	summary, err := Run(Options{Root: root}, func(p Progress) {
		if p.Type == "progress" {
			stages = append(stages, p.Stage)
		}
	})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if summary.FilesIndexed != 0 || summary.ChunksIndexed != 0 {
		t.Fatalf("unexpected summary: %+v", summary)
	}
	paths := workspace.Resolve(root)
	if _, err := os.Stat(paths.IndexDBPath); err != nil {
		t.Fatalf("expected index db created: %v", err)
	}
	if len(stages) == 0 || stages[0] != "scan" {
		t.Fatalf("expected scan stage progress, got %v", stages)
	}
}

func TestRunIncrementalAndDelete(t *testing.T) {
	requireSQLite(t)
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "a.go"), []byte("package main\n\nfunc A() {}\n"), 0o644); err != nil {
		t.Fatalf("write a.go: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "b.md"), []byte("# Title\nBody\n"), 0o644); err != nil {
		t.Fatalf("write b.md: %v", err)
	}

	summary, err := Run(Options{Root: root}, func(Progress) {})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if summary.FilesIndexed != 2 || summary.ChunksIndexed == 0 {
		t.Fatalf("unexpected summary: %+v", summary)
	}

	summary, err = Run(Options{Root: root}, func(Progress) {})
	if err != nil {
		t.Fatalf("run again: %v", err)
	}
	if summary.FilesIndexed != 0 || summary.ChunksIndexed != 0 {
		t.Fatalf("expected no reindex, got %+v", summary)
	}

	if err := os.Remove(filepath.Join(root, "a.go")); err != nil {
		t.Fatalf("remove: %v", err)
	}
	summary, err = Run(Options{Root: root}, func(Progress) {})
	if err != nil {
		t.Fatalf("run after delete: %v", err)
	}
	if summary.FilesIndexed != 0 {
		t.Fatalf("expected no new files indexed after delete, got %+v", summary)
	}

	store, err := metadata.Open(filepath.Join(root, ".scry", "index.db"))
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	files, err := store.ListFiles()
	if err != nil {
		t.Fatalf("list files: %v", err)
	}
	if len(files) != 1 || files[0] != "b.md" {
		t.Fatalf("expected only b.md in index, got %v", files)
	}
}

func TestRunCleanAndSkipEmptyChunks(t *testing.T) {
	requireSQLite(t)
	root := t.TempDir()
	paths := workspace.Resolve(root)
	if err := os.MkdirAll(paths.Workspace, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(paths.IndexDBPath, []byte("junk"), 0o644); err != nil {
		t.Fatalf("write junk db: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "empty.md"), []byte(""), 0o644); err != nil {
		t.Fatalf("write empty md: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "ok.go"), []byte("package main\n\nfunc A() {}\n"), 0o644); err != nil {
		t.Fatalf("write ok.go: %v", err)
	}

	summary, err := Run(Options{Root: root, Clean: true}, func(Progress) {})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if summary.FilesIndexed != 1 {
		t.Fatalf("expected 1 file indexed, got %+v", summary)
	}
}

func TestRunReadFileError(t *testing.T) {
	requireSQLite(t)
	root := t.TempDir()
	path := filepath.Join(root, "bad.go")
	if err := os.WriteFile(path, []byte("package main\n"), 0o000); err != nil {
		t.Fatalf("write bad.go: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chmod(path, 0o644)
	})

	_, err := Run(Options{Root: root}, func(Progress) {})
	if err == nil {
		t.Fatalf("expected read error")
	}
}

func TestRunMissingSQLite(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "main.go"), []byte("package main\n"), 0o644); err != nil {
		t.Fatalf("write main.go: %v", err)
	}
	orig := os.Getenv("PATH")
	if err := os.Setenv("PATH", ""); err != nil {
		t.Fatalf("setenv: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Setenv("PATH", orig)
	})

	_, err := Run(Options{Root: root}, func(Progress) {})
	if err == nil {
		t.Fatalf("expected sqlite3 missing error")
	}
}
