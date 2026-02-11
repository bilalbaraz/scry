package scan

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"scry/pkg/ignore"
)

type File struct {
	Path string
	Info os.FileInfo
}

type Scanner struct {
	Root    string
	Matcher *ignore.Matcher
}

func New(root string) (*Scanner, error) {
	matcher, err := ignore.Load(root)
	if err != nil {
		return nil, err
	}
	return &Scanner{Root: root, Matcher: matcher}, nil
}

func (s *Scanner) ListFiles() ([]File, error) {
	var files []File
	err := filepath.WalkDir(s.Root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == s.Root {
			return nil
		}
		rel, err := filepath.Rel(s.Root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if strings.HasPrefix(rel, ".git/") || rel == ".git" {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasPrefix(rel, ".scry/") || rel == ".scry" {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if s.Matcher != nil && s.Matcher.Ignored(rel) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		files = append(files, File{Path: path, Info: info})
		return nil
	})
	return files, err
}
