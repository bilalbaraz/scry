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
