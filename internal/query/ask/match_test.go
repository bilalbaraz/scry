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

func TestDistinctTermMatchCountEmpty(t *testing.T) {
	if got := DistinctTermMatchCount(nil, "text"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
	if got := DistinctTermMatchCount([]string{"ignore"}, ""); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestDistinctTermMatchCountSkipsEmptyTerm(t *testing.T) {
	got := DistinctTermMatchCount([]string{""}, "alpha")
	if got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}
