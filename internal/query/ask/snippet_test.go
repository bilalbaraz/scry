package ask

import (
	"strings"
	"testing"
)

func TestSelectTopEvidenceMaxTwo(t *testing.T) {
	chunks := []Chunk{{ID: "1"}, {ID: "2"}, {ID: "3"}}
	got := SelectTopEvidence(chunks, 2)
	if len(got) != 2 {
		t.Fatalf("expected 2 chunks, got %d", len(got))
	}
}

func TestSelectTopEvidenceZeroLimit(t *testing.T) {
	chunks := []Chunk{{ID: "1"}}
	got := SelectTopEvidence(chunks, 0)
	if got != nil {
		t.Fatalf("expected nil for zero limit")
	}
}

func TestSnippetAroundTerm(t *testing.T) {
	text := "alpha beta gamma delta epsilon"
	terms := []string{"gamma"}
	snippet := SnippetAroundTerm(text, terms, 10)
	if snippet == "" {
		t.Fatalf("expected snippet")
	}
	if snippet == text {
		t.Fatalf("expected trimmed snippet")
	}
}

func TestSnippetAroundTermNoMatchUsesTrim(t *testing.T) {
	text := "alpha beta gamma delta epsilon"
	terms := []string{"omega"}
	got := SnippetAroundTerm(text, terms, 8)
	if got != "alpha be..." {
		t.Fatalf("unexpected snippet: %q", got)
	}
}

func TestSnippetAroundTermZeroMax(t *testing.T) {
	got := SnippetAroundTerm("alpha", []string{"alpha"}, 0)
	if got != "" {
		t.Fatalf("expected empty snippet, got %q", got)
	}
}

func TestTrimSnippetShortText(t *testing.T) {
	got := trimSnippet("short", 10)
	if got != "short" {
		t.Fatalf("unexpected trim: %q", got)
	}
}

func TestSnippetAroundTermAddsEllipses(t *testing.T) {
	text := "alpha beta gamma delta epsilon zeta"
	terms := []string{"gamma"}
	got := SnippetAroundTerm(text, terms, 10)
	if !strings.HasPrefix(got, "...") || !strings.HasSuffix(got, "...") {
		t.Fatalf("expected ellipses around snippet, got %q", got)
	}
}

func TestSnippetAroundTermEmptyTerms(t *testing.T) {
	got := SnippetAroundTerm("alpha beta", nil, 5)
	if got != "alpha..." {
		t.Fatalf("unexpected snippet: %q", got)
	}
}

func TestSnippetAroundTermAtStart(t *testing.T) {
	text := "alpha beta"
	got := SnippetAroundTerm(text, []string{"alpha"}, 20)
	if strings.HasPrefix(got, "...") || strings.HasSuffix(got, "...") {
		t.Fatalf("did not expect ellipses, got %q", got)
	}
}

func TestSnippetAroundTermAtEnd(t *testing.T) {
	text := "alpha beta gamma"
	got := SnippetAroundTerm(text, []string{"gamma"}, 10)
	if !strings.HasPrefix(got, "...") {
		t.Fatalf("expected prefix ellipsis, got %q", got)
	}
	if strings.HasSuffix(got, "...") {
		t.Fatalf("did not expect suffix ellipsis, got %q", got)
	}
}

func TestSnippetAroundTermSkipsMissingTerm(t *testing.T) {
	text := "alpha beta gamma"
	got := SnippetAroundTerm(text, []string{"missing", "beta"}, 10)
	if !strings.Contains(got, "beta") {
		t.Fatalf("expected snippet containing beta, got %q", got)
	}
}
