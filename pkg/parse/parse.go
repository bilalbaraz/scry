package parse

import (
	"path/filepath"
	"strings"

	"scry/pkg/chunk"
)

func ChunksForFile(path string, content string) []chunk.Chunk {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".go":
		return ChunkGo(path, content)
	case ".md", ".markdown":
		return ChunkMarkdown(path, content)
	default:
		return nil
	}
}
