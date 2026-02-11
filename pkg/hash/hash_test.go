package hash

import "testing"

func TestFileHashDeterministic(t *testing.T) {
	data := []byte("hello world")
	got1 := FileHash(data)
	got2 := FileHash(data)
	if got1 != got2 {
		t.Fatalf("expected deterministic hash, got %s and %s", got1, got2)
	}
}

func TestChunkHashChangesWithRange(t *testing.T) {
	fileHash := FileHash([]byte("file"))
	text := "chunk content"
	h1 := ChunkHash(fileHash, 1, 3, text)
	h2 := ChunkHash(fileHash, 1, 4, text)
	if h1 == h2 {
		t.Fatalf("expected different hashes for different ranges")
	}
}
