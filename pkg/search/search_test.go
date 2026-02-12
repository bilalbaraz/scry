package search

import (
	"os/exec"
	"path/filepath"
	"testing"

	"scry/pkg/metadata"
)

func requireSQLite(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("sqlite3"); err != nil {
		t.Skipf("sqlite3 not found: %v", err)
	}
}

func TestSearchRanksAndLimits(t *testing.T) {
	requireSQLite(t)
	root := t.TempDir()
	dbPath := filepath.Join(root, "index.db")
	store, err := metadata.Open(dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	file1 := metadata.FileRecord{Path: "b/file.go", Hash: "h1", MTime: 1, Size: 10}
	chunk1 := metadata.ChunkRecord{ID: "c1", FilePath: "b/file.go", StartLine: 1, EndLine: 1, Hash: "ch1", Content: "alpha beta"}
	terms1 := []metadata.TermRecord{
		{Term: "alpha", ChunkID: "c1", TF: 2},
	}
	if err := store.ReplaceFileData(file1, []metadata.ChunkRecord{chunk1}, terms1); err != nil {
		t.Fatalf("replace file1: %v", err)
	}

	file2 := metadata.FileRecord{Path: "a/file.go", Hash: "h2", MTime: 1, Size: 10}
	chunk2 := metadata.ChunkRecord{ID: "c2", FilePath: "a/file.go", StartLine: 1, EndLine: 1, Hash: "ch2", Content: "alpha beta"}
	terms2 := []metadata.TermRecord{
		{Term: "alpha", ChunkID: "c2", TF: 1},
		{Term: "beta", ChunkID: "c2", TF: 1},
	}
	if err := store.ReplaceFileData(file2, []metadata.ChunkRecord{chunk2}, terms2); err != nil {
		t.Fatalf("replace file2: %v", err)
	}

	engine := New(store)
	results, err := engine.Search("alpha beta", 1)
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Chunk.FilePath != "a/file.go" {
		t.Fatalf("expected path tiebreaker, got %s", results[0].Chunk.FilePath)
	}
}

func TestSearchEmptyOrNoHits(t *testing.T) {
	requireSQLite(t)
	root := t.TempDir()
	dbPath := filepath.Join(root, "index.db")
	store, err := metadata.Open(dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	engine := New(store)

	results, err := engine.Search("a", 10)
	if err != nil {
		t.Fatalf("search short query: %v", err)
	}
	if results != nil {
		t.Fatalf("expected nil results for empty terms")
	}

	results, err = engine.Search("gamma", 10)
	if err != nil {
		t.Fatalf("search no hits: %v", err)
	}
	if results != nil {
		t.Fatalf("expected nil results for no hits")
	}
}

func TestSearchNoLimit(t *testing.T) {
	requireSQLite(t)
	root := t.TempDir()
	dbPath := filepath.Join(root, "index.db")
	store, err := metadata.Open(dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	file := metadata.FileRecord{Path: "file.go", Hash: "h1", MTime: 1, Size: 10}
	chunk := metadata.ChunkRecord{ID: "c1", FilePath: "file.go", StartLine: 1, EndLine: 1, Hash: "ch1", Content: "alpha"}
	terms := []metadata.TermRecord{{Term: "alpha", ChunkID: "c1", TF: 1}}
	if err := store.ReplaceFileData(file, []metadata.ChunkRecord{chunk}, terms); err != nil {
		t.Fatalf("replace file: %v", err)
	}

	engine := New(store)
	results, err := engine.Search("alpha", 0)
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

type fakeStore struct {
	termHits map[string][]metadata.TermHit
	termErr  error
	chunks   []metadata.ChunkView
	chunkErr error
}

func (f *fakeStore) TermHits(term string) ([]metadata.TermHit, error) {
	if f.termErr != nil {
		return nil, f.termErr
	}
	return f.termHits[term], nil
}

func (f *fakeStore) GetChunksByIDs(ids []string) ([]metadata.ChunkView, error) {
	if f.chunkErr != nil {
		return nil, f.chunkErr
	}
	return f.chunks, nil
}

func TestSearchTermHitsError(t *testing.T) {
	engine := New(&fakeStore{termErr: errSentinel{}})
	if _, err := engine.Search("alpha", 10); err == nil {
		t.Fatalf("expected error")
	}
}

func TestSearchChunkFetchError(t *testing.T) {
	engine := New(&fakeStore{
		termHits: map[string][]metadata.TermHit{
			"alpha": {{ChunkID: "c1", TF: 1}},
		},
		chunkErr: errSentinel{},
	})
	if _, err := engine.Search("alpha", 10); err == nil {
		t.Fatalf("expected error")
	}
}

type errSentinel struct{}

func (errSentinel) Error() string { return "boom" }
