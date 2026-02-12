package lexical

import "testing"

func TestTokenize(t *testing.T) {
	got := Tokenize("Go, go! 1 a 22.")
	want := []string{"go", "go", "22"}
	if len(got) != len(want) {
		t.Fatalf("expected %d terms, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("term %d expected %q got %q", i, want[i], got[i])
		}
	}
}
