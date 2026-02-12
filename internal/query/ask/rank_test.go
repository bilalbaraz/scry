package ask

import "testing"

func TestApplyMatchPreference(t *testing.T) {
	chunks := []Chunk{
		{ID: "1", Text: "ignore scan", FilePath: "pkg/scan/x.go"},
		{ID: "2", Text: "ignore", FilePath: "pkg/ignore/x.go"},
	}
	terms := []string{"ignore", "scan"}
	got := ApplyMatchPreference(chunks, terms)
	if len(got) != 1 {
		t.Fatalf("expected 1 strong chunk, got %d", len(got))
	}
	if got[0].ID != "1" {
		t.Fatalf("expected strong chunk id 1, got %s", got[0].ID)
	}
}

func TestApplyMatchPreferenceFallback(t *testing.T) {
	chunks := []Chunk{
		{ID: "1", Text: "ignore", FilePath: "pkg/ignore/x.go"},
		{ID: "2", Text: "scan", FilePath: "pkg/scan/x.go"},
	}
	terms := []string{"ignore", "scan"}
	got := ApplyMatchPreference(chunks, terms)
	if len(got) != 2 {
		t.Fatalf("expected fallback to weak matches, got %d", len(got))
	}
}

func TestApplyMatchPreferenceEmpty(t *testing.T) {
	if got := ApplyMatchPreference(nil, []string{"ignore"}); got != nil {
		t.Fatalf("expected nil for empty chunks")
	}
	if got := ApplyMatchPreference([]Chunk{{ID: "1", Text: "ignore"}}, nil); got != nil {
		t.Fatalf("expected nil for empty terms")
	}
}

func TestSortCandidatesTiebreakers(t *testing.T) {
	chunks := []Chunk{
		{ID: "1", FilePath: "b/file.go", Score: 2.0, StartLine: 10},
		{ID: "2", FilePath: "a/file.go", Score: 2.0, StartLine: 20},
		{ID: "3", FilePath: "a/file.go", Score: 2.0, StartLine: 5},
		{ID: "4", FilePath: "c/file.go", Score: 3.0, StartLine: 1},
	}
	got := SortCandidates(chunks)
	if got[0].ID != "4" {
		t.Fatalf("expected highest score first, got %s", got[0].ID)
	}
	if got[1].ID != "3" || got[2].ID != "2" {
		t.Fatalf("expected path/line tiebreaker order, got %s then %s", got[1].ID, got[2].ID)
	}
}
