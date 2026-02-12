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

func TestApplyWhitelistPromotionEarlyReturn(t *testing.T) {
	chunks := []Chunk{{ID: "1", FilePath: "pkg/scan/scan.go", Score: 1.0}}
	if got := ApplyWhitelistPromotion(chunks, nil, DefaultWhitelistRules()); got == nil {
		t.Fatalf("expected chunks back for empty terms")
	}
	if got := ApplyWhitelistPromotion(nil, []string{"scan"}, DefaultWhitelistRules()); got != nil {
		t.Fatalf("expected nil for empty chunks")
	}
	if got := ApplyWhitelistPromotion(chunks, []string{"scan"}, nil); got == nil {
		t.Fatalf("expected chunks back for empty rules")
	}
}

func TestApplyWhitelistPromotionNoMatch(t *testing.T) {
	chunks := []Chunk{{ID: "1", FilePath: "pkg/indexer/indexer.go", Score: 1.0}}
	got := ApplyWhitelistPromotion(chunks, []string{"search"}, DefaultWhitelistRules())
	if got[0].Score != 1.0 {
		t.Fatalf("expected no boost for non-matching terms/path")
	}
}
