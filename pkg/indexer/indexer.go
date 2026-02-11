package indexer

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"scry/pkg/hash"
	"scry/pkg/index/lexical"
	"scry/pkg/metadata"
	"scry/pkg/parse"
	"scry/pkg/scan"
	"scry/pkg/workspace"
)

type Options struct {
	Root         string
	Clean        bool
	NoEmbeddings bool
	JSON         bool
}

type Progress struct {
	Type       string `json:"type"`
	Stage      string `json:"stage,omitempty"`
	FilesTotal int    `json:"files_total,omitempty"`
	File       string `json:"file,omitempty"`
	Chunks     int    `json:"chunks,omitempty"`
	Message    string `json:"message,omitempty"`
}

type Summary struct {
	FilesIndexed  int
	ChunksIndexed int
}

func Run(opts Options, emit func(Progress)) (Summary, error) {
	paths := workspace.Resolve(opts.Root)
	if err := workspace.Ensure(paths); err != nil {
		return Summary{}, err
	}
	if opts.Clean {
		_ = os.Remove(paths.IndexDBPath)
	}
	store, err := metadata.Open(paths.IndexDBPath)
	if err != nil {
		return Summary{}, err
	}

	scanner, err := scan.New(opts.Root)
	if err != nil {
		return Summary{}, err
	}
	files, err := scanner.ListFiles()
	if err != nil {
		return Summary{}, err
	}
	files = filterSupported(files)
	emit(Progress{Type: "progress", Stage: "scan", FilesTotal: len(files)})

	// Remove deleted files
	indexed, err := store.ListFiles()
	if err != nil {
		return Summary{}, err
	}
	indexedSet := map[string]struct{}{}
	for _, p := range indexed {
		indexedSet[p] = struct{}{}
	}
	currentSet := map[string]struct{}{}
	for _, f := range files {
		rel, _ := filepath.Rel(opts.Root, f.Path)
		rel = filepath.ToSlash(rel)
		currentSet[rel] = struct{}{}
	}
	for p := range indexedSet {
		if _, ok := currentSet[p]; !ok {
			if err := store.DeleteFile(p); err != nil {
				return Summary{}, err
			}
		}
	}

	lex := lexical.New()
	summary := Summary{}
	for _, f := range files {
		rel, _ := filepath.Rel(opts.Root, f.Path)
		rel = filepath.ToSlash(rel)
		data, err := os.ReadFile(f.Path)
		if err != nil {
			return Summary{}, err
		}
		fileHash := hash.FileHash(data)
		rec, ok, err := store.GetFile(rel)
		if err != nil {
			return Summary{}, err
		}
		if ok && rec.Hash == fileHash {
			continue
		}

		chunks := parse.ChunksForFile(rel, string(data))
		if len(chunks) == 0 {
			continue
		}
		var chunkRecords []metadata.ChunkRecord
		var termRecords []metadata.TermRecord
		for _, c := range chunks {
			ch := c
			ch.FilePath = rel
			chunkHash := hash.ChunkHash(fileHash, ch.StartLine, ch.EndLine, ch.Text)
			chunkID := chunkHash
			chunkRecords = append(chunkRecords, metadata.ChunkRecord{
				ID:        chunkID,
				FilePath:  rel,
				StartLine: ch.StartLine,
				EndLine:   ch.EndLine,
				Hash:      chunkHash,
				Content:   ch.Text,
			})
			postings := lex.Add(chunkID, ch.Text)
			for _, p := range postings {
				termRecords = append(termRecords, metadata.TermRecord{Term: p.Term, ChunkID: p.ChunkID, TF: p.TF})
			}
		}

		fr := metadata.FileRecord{
			Path:  rel,
			Hash:  fileHash,
			MTime: f.Info.ModTime().Unix(),
			Size:  f.Info.Size(),
		}
		if err := store.ReplaceFileData(fr, chunkRecords, termRecords); err != nil {
			return Summary{}, err
		}
		summary.FilesIndexed++
		summary.ChunksIndexed += len(chunkRecords)
		emit(Progress{Type: "progress", Stage: "index", File: rel, Chunks: len(chunkRecords)})
	}

	return summary, nil
}

func filterSupported(files []scan.File) []scan.File {
	var out []scan.File
	for _, f := range files {
		ext := strings.ToLower(filepath.Ext(f.Path))
		switch ext {
		case ".go", ".md", ".markdown":
			out = append(out, f)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out
}

// Placeholder to avoid unused import for time in future progress timestamps.
