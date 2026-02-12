package workspace

import (
	"os"
	"path/filepath"
)

type Paths struct {
	Root        string
	Workspace   string
	IndexDBPath string
}

func Resolve(root string) Paths {
	ws := filepath.Join(root, ".scry")
	return Paths{
		Root:        root,
		Workspace:   ws,
		IndexDBPath: filepath.Join(ws, "index.db"),
	}
}

func Ensure(paths Paths) error {
	return os.MkdirAll(paths.Workspace, 0o755)
}

func Exists(paths Paths) bool {
	_, err := os.Stat(paths.IndexDBPath)
	return err == nil
}
