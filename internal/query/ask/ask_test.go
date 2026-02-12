package ask

import "testing"

func TestDecideIDKNoEvidence(t *testing.T) {
	dec := Decide(nil, []string{"ignore"}, AskOptions{MaxEvidence: 2, SnippetChars: 80, MinScore: 1})
	if dec.Reason != "no_evidence" {
		t.Fatalf("expected no_evidence, got %s", dec.Reason)
	}
}

func TestDecideIDKLowScore(t *testing.T) {
	chunks := []Chunk{{ID: "1", Text: "ignore patterns", FilePath: "pkg/ignore/x.go", Score: 0.2}}
	dec := Decide(chunks, []string{"ignore"}, AskOptions{MaxEvidence: 2, SnippetChars: 80, MinScore: 5})
	if dec.Reason != "low_score" {
		t.Fatalf("expected low_score, got %s", dec.Reason)
	}
}

func TestDecideReturnsEvidence(t *testing.T) {
	chunks := []Chunk{
		{ID: "1", Text: "ignore patterns", FilePath: "pkg/ignore/x.go", Score: 2.0},
		{ID: "2", Text: "scan rules", FilePath: "pkg/scan/x.go", Score: 1.5},
	}
	dec := Decide(chunks, []string{"ignore"}, AskOptions{MaxEvidence: 2, SnippetChars: 80, MinScore: 1.0})
	if dec.Reason != "" {
		t.Fatalf("expected no reason, got %s", dec.Reason)
	}
	if len(dec.Evidence) == 0 {
		t.Fatalf("expected evidence")
	}
}

func TestAnswerHeader(t *testing.T) {
	evidence := []Evidence{{Chunk: Chunk{ID: "1"}}}
	got := AnswerHeader(evidence)
	if got != "Found 1 relevant evidence chunk(s)." {
		t.Fatalf("unexpected header: %s", got)
	}
}
