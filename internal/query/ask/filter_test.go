package ask

import "testing"

func TestFilterByQueryTerms(t *testing.T) {
	chunks := []Chunk{
		{ID: "1", Text: "This talks about ignore patterns."},
		{ID: "2", Text: "Unrelated content."},
		{ID: "3", Text: "Scan rules are defined here."},
	}
	terms := []string{"ignore", "scan"}
	got := FilterByQueryTerms(chunks, terms)
	if len(got) != 2 {
		t.Fatalf("expected 2 chunks, got %d", len(got))
	}
	if got[0].ID != "1" || got[1].ID != "3" {
		t.Fatalf("unexpected chunk IDs: %v, %v", got[0].ID, got[1].ID)
	}
}

func TestFilterByQueryTermsEmptyInputs(t *testing.T) {
	if got := FilterByQueryTerms(nil, []string{"ignore"}); got != nil {
		t.Fatalf("expected nil for empty chunks")
	}
	if got := FilterByQueryTerms([]Chunk{{ID: "1", Text: "ignore"}}, nil); got != nil {
		t.Fatalf("expected nil for empty terms")
	}
}

func TestHasAnyTermSkipsEmpty(t *testing.T) {
	if hasAnyTerm("alpha", []string{""}) {
		t.Fatalf("expected no match for empty term")
	}
	if !hasAnyTerm("alpha beta", []string{"", "beta"}) {
		t.Fatalf("expected match for beta")
	}
}
