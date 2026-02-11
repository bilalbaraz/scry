package parse

import (
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
	"strings"

	"scry/pkg/chunk"
)

type goChunk struct {
	start int
	end   int
}

func ChunkGo(path string, content string) []chunk.Chunk {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, content, parser.ParseComments)
	if err != nil {
		return fallbackGoChunks(path, content)
	}

	var spans []goChunk
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			start := fset.Position(d.Pos()).Line
			end := fset.Position(d.End()).Line
			spans = append(spans, goChunk{start: start, end: end})
		case *ast.GenDecl:
			if d.Tok != token.TYPE {
				continue
			}
			start := fset.Position(d.Pos()).Line
			end := fset.Position(d.End()).Line
			spans = append(spans, goChunk{start: start, end: end})
		}
	}

	if len(spans) == 0 {
		return fallbackGoChunks(path, content)
	}

	sort.Slice(spans, func(i, j int) bool { return spans[i].start < spans[j].start })
	lines := strings.Split(content, "\n")
	chunks := make([]chunk.Chunk, 0, len(spans))
	for _, sp := range spans {
		start := clampLine(sp.start, len(lines))
		end := clampLine(sp.end, len(lines))
		text := strings.Join(lines[start-1:end], "\n")
		chunks = append(chunks, chunk.Chunk{
			FilePath:  path,
			StartLine: start,
			EndLine:   end,
			Text:      text,
			Lang:      "go",
		})
	}
	return chunks
}

func fallbackGoChunks(path string, content string) []chunk.Chunk {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return nil
	}
	return []chunk.Chunk{
		{
			FilePath:  path,
			StartLine: 1,
			EndLine:   len(lines),
			Text:      content,
			Lang:      "go",
		},
	}
}

func clampLine(line int, max int) int {
	if line < 1 {
		return 1
	}
	if line > max {
		return max
	}
	return line
}
