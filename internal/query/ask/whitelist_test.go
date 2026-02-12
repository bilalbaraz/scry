package ask

import "testing"

func TestApplyWhitelistPromotion(t *testing.T) {
	chunks := []Chunk{
		{ID: "1", FilePath: "pkg/scan/scan.go", Score: 1.0},
		{ID: "2", FilePath: "pkg/indexer/indexer.go", Score: 2.0},
	}
	terms := []string{"scan"}
	got := ApplyWhitelistPromotion(chunks, terms, DefaultWhitelistRules())
	if got[0].Score <= 1.0 {
		t.Fatalf("expected whitelist boost for pkg/scan")
	}
}

func TestIsScanRelated(t *testing.T) {
	if !isScanRelated([]string{"ignore"}) {
		t.Fatalf("expected scan-related true")
	}
	if isScanRelated([]string{"search"}) {
		t.Fatalf("expected scan-related false")
	}
}
