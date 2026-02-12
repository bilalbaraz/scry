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

func TestTermsMatch(t *testing.T) {
	if !termsMatch(nil, map[string]struct{}{}) {
		t.Fatalf("expected empty rule terms to match")
	}
	if termsMatch([]string{"alpha"}, map[string]struct{}{"beta": {}}) {
		t.Fatalf("expected no match")
	}
	if !termsMatch([]string{"alpha"}, map[string]struct{}{"alpha": {}}) {
		t.Fatalf("expected match")
	}
}

func TestApplyBoostsEarlyReturn(t *testing.T) {
	chunks := []Chunk{{ID: "1", FilePath: "pkg/scan/scan.go", Score: 1.0}}
	if got := ApplyBoosts(chunks, nil, DefaultBoostRules()); got == nil {
		t.Fatalf("expected chunks back for empty terms")
	}
	if got := ApplyBoosts(nil, []string{"scan"}, DefaultBoostRules()); got != nil {
		t.Fatalf("expected nil for empty chunks")
	}
	if got := ApplyBoosts(chunks, []string{"scan"}, nil); got == nil {
		t.Fatalf("expected chunks back for empty rules")
	}
}

func TestApplyBoostsNoPathMatch(t *testing.T) {
	chunks := []Chunk{{ID: "1", FilePath: "docs/readme.md", Score: 1.0}}
	terms := []string{"scan"}
	got := ApplyBoosts(chunks, terms, DefaultBoostRules())
	if got[0].Score != 1.0 {
		t.Fatalf("expected no boost for non-matching path")
	}
}
