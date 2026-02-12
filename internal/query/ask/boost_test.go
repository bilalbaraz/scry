package ask

import "testing"

func TestApplyBoosts(t *testing.T) {
	chunks := []Chunk{
		{ID: "1", FilePath: "pkg/ignore/ignore.go", Score: 1.0},
		{ID: "2", FilePath: "pkg/scan/scan.go", Score: 1.0},
		{ID: "3", FilePath: "README.md", Score: 1.0},
	}
	terms := []string{"ignore"}
	rules := DefaultBoostRules()
	got := ApplyBoosts(chunks, terms, rules)
	if got[0].Score <= 1.0 || got[1].Score <= 1.0 {
		t.Fatalf("expected boosted scores for ignore/scan paths")
	}
	if got[2].Score != 1.0 {
		t.Fatalf("did not expect boost for README")
	}
}

func TestApplyBoostsPenalty(t *testing.T) {
	chunks := []Chunk{
		{ID: "1", FilePath: "internal/query/ask/boost.go", Score: 1.0},
		{ID: "2", FilePath: "cmd/scry/ask.go", Score: 1.0},
		{ID: "3", FilePath: "pkg/scan/scan.go", Score: 1.0},
	}
	terms := []string{"scan"}
	rules := DefaultBoostRules()
	got := ApplyBoosts(chunks, terms, rules)
	if got[0].Score >= 1.0 || got[1].Score >= 1.0 {
		t.Fatalf("expected penalties for internal/query/ask and cmd paths")
	}
	if got[2].Score <= 1.0 {
		t.Fatalf("expected positive boost for pkg/scan path")
	}
}
