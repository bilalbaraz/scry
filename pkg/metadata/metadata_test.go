package metadata

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func requireSQLite(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skipf("sqlite3 not found: %v", err)
	}
}

func TestOpenRejectsDirectoryPath(t *testing.T) {
	requireSQLite(t)
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "dbdir"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	_, err := Open(filepath.Join(root, "dbdir"))
	if err == nil {
		t.Fatalf("expected error opening directory as db path")
	}
}

func TestOpenMissingSQLite(t *testing.T) {
	orig := os.Getenv("PATH")
	if err := os.Setenv("PATH", ""); err != nil {
		t.Fatalf("setenv: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Setenv("PATH", orig)
	})
	if _, err := Open(filepath.Join(t.TempDir(), "index.db")); err == nil {
		t.Fatalf("expected sqlite3 missing error")
	}
}

func TestDBLifecycle(t *testing.T) {
	requireSQLite(t)
	root := t.TempDir()
	dbPath := filepath.Join(root, "index.db")
	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	file := FileRecord{Path: "src/main.go", Hash: "h1", MTime: 1, Size: 10}
	chunk := ChunkRecord{ID: "c1", FilePath: "src/main.go", StartLine: 1, EndLine: 1, Hash: "ch1", Content: "alpha beta"}
	terms := []TermRecord{
		{Term: "alpha", ChunkID: "c1", TF: 2},
		{Term: "beta", ChunkID: "c1", TF: 1},
	}
	if err := store.ReplaceFileData(file, []ChunkRecord{chunk}, terms); err != nil {
		t.Fatalf("replace file data: %v", err)
	}

	gotFile, ok, err := store.GetFile("src/main.go")
	if err != nil {
		t.Fatalf("get file: %v", err)
	}
	if !ok || gotFile.Path != "src/main.go" {
		t.Fatalf("expected file record")
	}

	files, err := store.ListFiles()
	if err != nil {
		t.Fatalf("list files: %v", err)
	}
	if len(files) != 1 || files[0] != "src/main.go" {
		t.Fatalf("unexpected files: %v", files)
	}

	view, ok, err := store.GetChunk("c1")
	if err != nil {
		t.Fatalf("get chunk: %v", err)
	}
	if !ok || view.Content != "alpha beta" {
		t.Fatalf("unexpected chunk view")
	}

	views, err := store.GetChunksByIDs([]string{"c1"})
	if err != nil {
		t.Fatalf("get chunks: %v", err)
	}
	if len(views) != 1 || views[0].ID != "c1" {
		t.Fatalf("unexpected chunks: %v", views)
	}

	hits, err := store.TermHits("alpha")
	if err != nil {
		t.Fatalf("term hits: %v", err)
	}
	if len(hits) != 1 || hits[0].TF != 2 {
		t.Fatalf("unexpected term hits: %v", hits)
	}

	stats, err := store.Stats()
	if err != nil {
		t.Fatalf("stats: %v", err)
	}
	if stats.Files != 1 || stats.Chunks != 1 || stats.Terms != 2 {
		t.Fatalf("unexpected stats: %+v", stats)
	}

	if err := store.DeleteFile("src/main.go"); err != nil {
		t.Fatalf("delete file: %v", err)
	}
	files, err = store.ListFiles()
	if err != nil {
		t.Fatalf("list files: %v", err)
	}
	if len(files) != 0 {
		t.Fatalf("expected no files after delete, got %v", files)
	}
}

func TestEmptyQueries(t *testing.T) {
	requireSQLite(t)
	root := t.TempDir()
	dbPath := filepath.Join(root, "index.db")
	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	if _, ok, err := store.GetFile("missing.go"); err != nil || ok {
		t.Fatalf("expected missing file, ok=%v err=%v", ok, err)
	}
	if _, ok, err := store.GetChunk("missing"); err != nil || ok {
		t.Fatalf("expected missing chunk, ok=%v err=%v", ok, err)
	}
	if chunks, err := store.GetChunksByIDs(nil); err != nil || chunks != nil {
		t.Fatalf("expected nil chunks for empty ids, got %v err=%v", chunks, err)
	}
	if hits, err := store.TermHits("none"); err != nil || len(hits) != 0 {
		t.Fatalf("expected no term hits, got %v err=%v", hits, err)
	}
	stats, err := store.Stats()
	if err != nil {
		t.Fatalf("stats: %v", err)
	}
	if stats.Files != 0 || stats.Chunks != 0 || stats.Terms != 0 {
		t.Fatalf("unexpected stats: %+v", stats)
	}
}

func TestHelpers(t *testing.T) {
	lines := splitLines("a\n\nb\n")
	if len(lines) != 2 || lines[0] != "a" || lines[1] != "b" {
		t.Fatalf("unexpected lines: %v", lines)
	}
	if got := parseInt64("nope"); got != 0 {
		t.Fatalf("expected 0 for invalid int, got %d", got)
	}
	if got, err := decodeHexString(""); err != nil || got != "" {
		t.Fatalf("expected empty decode, got %q err=%v", got, err)
	}
	if _, err := decodeHexString("zz"); err == nil {
		t.Fatalf("expected decode error")
	}
}

func TestRunQueryError(t *testing.T) {
	requireSQLite(t)
	root := t.TempDir()
	dbPath := filepath.Join(root, "index.db")
	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if _, err := store.runQuery("SELECT nope;"); err == nil {
		t.Fatalf("expected query error")
	}
}

func TestDBSchemaMissingErrors(t *testing.T) {
	requireSQLite(t)
	root := t.TempDir()
	dbPath := filepath.Join(root, "raw.db")
	if err := os.WriteFile(dbPath, []byte(""), 0o644); err != nil {
		t.Fatalf("write db: %v", err)
	}
	store := &DB{Path: dbPath}

	if _, _, err := store.GetFile("x"); err == nil {
		t.Fatalf("expected GetFile error")
	}
	if _, err := store.ListFiles(); err == nil {
		t.Fatalf("expected ListFiles error")
	}
	if _, _, err := store.GetChunk("c1"); err == nil {
		t.Fatalf("expected GetChunk error")
	}
	if _, err := store.GetChunksByIDs([]string{"c1"}); err == nil {
		t.Fatalf("expected GetChunksByIDs error")
	}
	if _, err := store.TermHits("alpha"); err == nil {
		t.Fatalf("expected TermHits error")
	}
	if _, err := store.Stats(); err == nil {
		t.Fatalf("expected Stats error")
	}
}

func TestScalarIntEmptyResult(t *testing.T) {
	requireSQLite(t)
	root := t.TempDir()
	dbPath := filepath.Join(root, "index.db")
	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	got, err := store.scalarInt("SELECT 1 WHERE 0;")
	if err != nil {
		t.Fatalf("scalarInt: %v", err)
	}
	if got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}
