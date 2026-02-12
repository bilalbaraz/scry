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
