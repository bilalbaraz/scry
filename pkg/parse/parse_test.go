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
