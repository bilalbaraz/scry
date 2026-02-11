package parse

import (
	"strings"

	"scry/pkg/chunk"
)

func ChunkMarkdown(path string, content string) []chunk.Chunk {
	lines := strings.Split(content, "\n")
	var chunks []chunk.Chunk
	start := 1
	current := ""
	for i, line := range lines {
		ln := i + 1
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			if current != "" {
				chunks = append(chunks, chunk.Chunk{
					FilePath:  path,
					StartLine: start,
					EndLine:   ln - 1,
					Text:      current,
					Lang:      "md",
				})
			}
			start = ln
			current = line
			continue
		}
		if current == "" {
			current = line
			start = ln
		} else {
			current += "\n" + line
		}
	}
	if strings.TrimSpace(current) != "" {
		chunks = append(chunks, chunk.Chunk{
			FilePath:  path,
			StartLine: start,
			EndLine:   len(lines),
			Text:      current,
			Lang:      "md",
		})
	}
	return chunks
}
