package ask

import "testing"

func TestTokenizeQuery(t *testing.T) {
	q := "Ignore, scan! gitignore; Exclude? ignore"
	got := TokenizeQuery(q)
	want := []string{"ignore", "scan", "gitignore", "exclude"}
	if len(got) != len(want) {
		t.Fatalf("expected %d terms, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("term %d expected %q got %q", i, want[i], got[i])
		}
	}
}

func TestTokenizeQueryDropsShort(t *testing.T) {
	q := "go a an the" // only "the" should remain
	got := TokenizeQuery(q)
	if len(got) != 1 || got[0] != "the" {
		t.Fatalf("expected [the], got %v", got)
	}
}
