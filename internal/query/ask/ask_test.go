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
