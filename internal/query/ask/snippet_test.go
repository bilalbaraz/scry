package ask

import "testing"

func TestSelectTopEvidenceMaxTwo(t *testing.T) {
	chunks := []Chunk{{ID: "1"}, {ID: "2"}, {ID: "3"}}
	got := SelectTopEvidence(chunks, 2)
	if len(got) != 2 {
		t.Fatalf("expected 2 chunks, got %d", len(got))
	}
}

func TestSnippetAroundTerm(t *testing.T) {
	text := "alpha beta gamma delta epsilon"
	terms := []string{"gamma"}
	snippet := SnippetAroundTerm(text, terms, 10)
	if snippet == "" {
		t.Fatalf("expected snippet")
	}
	if snippet == text {
		t.Fatalf("expected trimmed snippet")
	}
}
