package ignore

import (
	"path/filepath"
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
