package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolvePaths(t *testing.T) {
	root := filepath.Join("root", "repo")
	paths := Resolve(root)
	if paths.Root != root {
		t.Fatalf("unexpected root: %s", paths.Root)
	}
	if paths.Workspace != filepath.Join(root, ".scry") {
		t.Fatalf("unexpected workspace: %s", paths.Workspace)
	}
	if paths.IndexDBPath != filepath.Join(root, ".scry", "index.db") {
		t.Fatalf("unexpected index db: %s", paths.IndexDBPath)
	}
}

func TestEnsureAndExists(t *testing.T) {
	root := t.TempDir()
	paths := Resolve(root)
	if Exists(paths) {
		t.Fatalf("expected index to not exist")
	}
	if err := Ensure(paths); err != nil {
		t.Fatalf("ensure: %v", err)
	}
	if Exists(paths) {
		t.Fatalf("expected index db to not exist before creation")
	}
	if err := os.WriteFile(paths.IndexDBPath, []byte(""), 0o644); err != nil {
		t.Fatalf("write index: %v", err)
	}
	if !Exists(paths) {
		t.Fatalf("expected index db to exist")
	}
}
