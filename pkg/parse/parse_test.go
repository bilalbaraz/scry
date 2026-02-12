package parse

import "testing"

func TestChunkGoByTopLevel(t *testing.T) {
	src := `package main

type Foo struct {}

func (f Foo) Hello() {}

func main() {}
`
	chunks := ChunkGo("main.go", src)
	if len(chunks) != 3 {
		t.Fatalf("expected 3 chunks, got %d", len(chunks))
	}
	if chunks[0].StartLine != 3 {
		t.Fatalf("expected first chunk start 3, got %d", chunks[0].StartLine)
	}
	if chunks[1].StartLine != 5 {
		t.Fatalf("expected second chunk start 5, got %d", chunks[1].StartLine)
	}
	if chunks[2].StartLine != 7 {
		t.Fatalf("expected third chunk start 7, got %d", chunks[2].StartLine)
	}
}

func TestChunkMarkdownByHeading(t *testing.T) {
	src := "# Title\nIntro\n## Section\nBody\n"
	chunks := ChunkMarkdown("doc.md", src)
	if len(chunks) != 2 {
		t.Fatalf("expected 2 chunks, got %d", len(chunks))
	}
	if chunks[0].StartLine != 1 || chunks[0].EndLine != 2 {
		t.Fatalf("unexpected first chunk range %d-%d", chunks[0].StartLine, chunks[0].EndLine)
	}
	if chunks[1].StartLine != 3 {
		t.Fatalf("unexpected second chunk start %d", chunks[1].StartLine)
	}
}

func TestChunksForFileUnknownExtension(t *testing.T) {
	chunks := ChunksForFile("note.txt", "hello")
	if chunks != nil {
		t.Fatalf("expected nil chunks for unknown extension")
	}
}

func TestChunksForFileGoAndMarkdown(t *testing.T) {
	if chunks := ChunksForFile("main.go", "package main\n\nfunc A() {}\n"); len(chunks) == 0 {
		t.Fatalf("expected go chunks")
	}
	if chunks := ChunksForFile("doc.markdown", "# Title\nBody\n"); len(chunks) == 0 {
		t.Fatalf("expected markdown chunks")
	}
}

func TestChunkGoFallback(t *testing.T) {
	src := "package main\nfunc {"
	chunks := ChunkGo("main.go", src)
	if len(chunks) != 1 {
		t.Fatalf("expected fallback chunk, got %d", len(chunks))
	}
	if chunks[0].StartLine != 1 || chunks[0].EndLine != 2 {
		t.Fatalf("unexpected fallback range %d-%d", chunks[0].StartLine, chunks[0].EndLine)
	}
}

func TestClampLineBounds(t *testing.T) {
	if got := clampLine(-1, 10); got != 1 {
		t.Fatalf("expected clamp to 1, got %d", got)
	}
	if got := clampLine(20, 10); got != 10 {
		t.Fatalf("expected clamp to max, got %d", got)
	}
}

func TestChunkMarkdownEmpty(t *testing.T) {
	chunks := ChunkMarkdown("empty.md", "")
	if len(chunks) != 0 {
		t.Fatalf("expected no chunks, got %d", len(chunks))
	}
}
