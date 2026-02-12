package ignore

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMatcherFromFiles(t *testing.T) {
	root := filepath.Join("..", "..", "testdata", "ignore")
	m, err := Load(root)
	if err != nil {
		t.Fatalf("load matcher: %v", err)
	}

	cases := []struct {
		path    string
		ignored bool
	}{
		{"vendor/module/file.go", true},
		{"build/output.o", true},
		{"notes.md", true},
		{"pkg/scan/scan_test.go", true},
		{"testdata/fixtures/input.txt", true},
		{"src/main.go", false},
		{"docs/readme.md", true},
	}

	for _, c := range cases {
		if got := m.Ignored(c.path); got != c.ignored {
			t.Fatalf("path %s ignored=%v want %v", c.path, got, c.ignored)
		}
	}
}

func TestLoadPatternsAndMatcherPaths(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, ".gitignore")
	content := "# comment\nbuild/\n*.log\nnested/file.txt\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	patterns, err := loadPatterns(path)
	if err != nil {
		t.Fatalf("load patterns: %v", err)
	}
	if len(patterns) != 3 {
		t.Fatalf("expected 3 patterns, got %d", len(patterns))
	}
	m := &Matcher{root: root, globs: patterns}
	cases := []struct {
		path    string
		ignored bool
	}{
		{"build/output.o", true},
		{"error.log", true},
		{"nested/file.txt", true},
		{"nested/other.txt", false},
	}
	for _, c := range cases {
		if got := m.Ignored(c.path); got != c.ignored {
			t.Fatalf("path %s ignored=%v want %v", c.path, got, c.ignored)
		}
	}
}

func TestLoadReturnsErrorOnInvalidIgnoreFile(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, ".gitignore"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if _, err := Load(root); err == nil {
		t.Fatalf("expected error")
	}
}

func TestLoadPatternsScannerError(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, ".gitignore")
	longLine := strings.Repeat("a", 70000)
	if err := os.WriteFile(path, []byte(longLine), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if _, err := loadPatterns(path); err == nil {
		t.Fatalf("expected scanner error")
	}
}

func TestLoadPatternsMissingFile(t *testing.T) {
	patterns, err := loadPatterns(filepath.Join(t.TempDir(), "missing"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if patterns != nil {
		t.Fatalf("expected nil patterns, got %v", patterns)
	}
}

func TestMatcherPathWithSlashNoMatch(t *testing.T) {
	m := &Matcher{root: t.TempDir(), globs: []string{"dir/file.txt"}}
	if m.Ignored("dir/other.txt") {
		t.Fatalf("expected no match for non-matching slash pattern")
	}
}

func TestLoadWithScryIgnoreError(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, ".scryignore"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if _, err := Load(root); err == nil {
		t.Fatalf("expected error for invalid .scryignore")
	}
}

func TestMatcherSkipsEmptyGlob(t *testing.T) {
	m := &Matcher{root: t.TempDir(), globs: []string{""}}
	if m.Ignored("file.txt") {
		t.Fatalf("expected empty glob to be skipped")
	}
}
