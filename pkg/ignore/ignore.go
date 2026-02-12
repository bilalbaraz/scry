package ignore

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type Matcher struct {
	root  string
	globs []string
}

func Load(root string) (*Matcher, error) {
	globs := defaultPatterns()
	gitIgnore, err := loadPatterns(filepath.Join(root, ".gitignore"))
	if err != nil {
		return nil, err
	}
	scryIgnore, err := loadPatterns(filepath.Join(root, ".scryignore"))
	if err != nil {
		return nil, err
	}
	globs = append(globs, gitIgnore...)
	globs = append(globs, scryIgnore...)
	return &Matcher{root: root, globs: globs}, nil
}

func (m *Matcher) Ignored(relPath string) bool {
	rel := filepath.ToSlash(relPath)
	for _, g := range m.globs {
		if g == "" {
			continue
		}
		if strings.HasSuffix(g, "/") {
			if strings.HasPrefix(rel, g) {
				return true
			}
			continue
		}
		if ok, _ := filepath.Match(g, rel); ok {
			return true
		}
		if strings.Contains(g, "/") {
			if ok, _ := filepath.Match(g, rel); ok {
				return true
			}
			continue
		}
		base := filepath.Base(rel)
		if ok, _ := filepath.Match(g, base); ok {
			return true
		}
	}
	return false
}

func loadPatterns(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, filepath.ToSlash(line))
	}
	return patterns, scanner.Err()
}

func defaultPatterns() []string {
	return []string{
		"*_test.go",
		"testdata/",
	}
}
