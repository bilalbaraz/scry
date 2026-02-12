package ask

import "testing"

func TestAskRankingIntegration(t *testing.T) {
	chunks := []Chunk{
		{ID: "1", FilePath: "README.md", Text: "scan rules overview", Score: 2.0, StartLine: 1},
		{ID: "2", FilePath: "cmd/scry/ask.go", Text: "ask command scan rules", Score: 2.5, StartLine: 1},
		{ID: "3", FilePath: "internal/query/ask/boost.go", Text: "apply boosts for scan", Score: 2.2, StartLine: 1},
		{ID: "4", FilePath: "pkg/scan/scan.go", Text: "scan rules are defined here", Score: 1.8, StartLine: 10},
		{ID: "5", FilePath: "pkg/ignore/ignore.go", Text: "ignore patterns and scan", Score: 1.7, StartLine: 5},
	}
	terms := TokenizeQuery("ignore nasıl çalışıyor")
	decision := Decide(chunks, terms, AskOptions{MaxEvidence: 2, SnippetChars: 80, MinScore: 0.1})
	if len(decision.Evidence) == 0 {
		t.Fatalf("expected evidence")
	}
	if decision.Evidence[0].Chunk.FilePath != "pkg/ignore/ignore.go" && decision.Evidence[0].Chunk.FilePath != "pkg/scan/scan.go" {
		t.Fatalf("expected top evidence from pkg/ignore or pkg/scan, got %s", decision.Evidence[0].Chunk.FilePath)
	}
}
