package ask

import "testing"

func TestBuildEvidenceSelectsTwo(t *testing.T) {
	chunks := []Chunk{
		{ID: "1", Text: "ignore patterns", FilePath: "pkg/ignore/x.go", Score: 2},
		{ID: "2", Text: "scan rules", FilePath: "pkg/scan/x.go", Score: 1},
		{ID: "3", Text: "other", FilePath: "README.md", Score: 3},
	}
	terms := []string{"ignore", "scan"}
	evidence := BuildEvidence(chunks, terms, Options{MaxEvidence: 2, SnippetChars: 40})
	if len(evidence) != 2 {
		t.Fatalf("expected 2 evidence chunks, got %d", len(evidence))
	}
}

func TestBuildPipelineDelegates(t *testing.T) {
	chunks := []Chunk{
		{ID: "1", Text: "ignore patterns", FilePath: "pkg/ignore/x.go", Score: 2},
	}
	terms := []string{"ignore"}
	dec := BuildPipeline(chunks, terms, AskOptions{MaxEvidence: 1, SnippetChars: 40, MinScore: 0.1})
	if len(dec.Evidence) != 1 {
		t.Fatalf("expected 1 evidence chunk, got %d", len(dec.Evidence))
	}
}
