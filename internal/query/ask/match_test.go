package ask

import "testing"

func TestDistinctTermMatchCount(t *testing.T) {
	terms := []string{"ignore", "scan", "gitignore"}
	text := "We scan ignore patterns in gitignore files."
	got := DistinctTermMatchCount(terms, text)
	if got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}

func TestDistinctTermMatchCountDedup(t *testing.T) {
	terms := []string{"ignore", "ignore"}
	text := "ignore this"
	got := DistinctTermMatchCount(terms, text)
	if got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}
