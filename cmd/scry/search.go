package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"scry/pkg/metadata"
	"scry/pkg/search"
	"scry/pkg/workspace"
)

func newSearchCmd() *cobra.Command {
	var limit int
	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Hybrid search across indexes",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := os.Getwd()
			if err != nil {
				return exitError{code: exitRuntimeError, err: err}
			}
			root, err = filepath.Abs(root)
			if err != nil {
				return exitError{code: exitRuntimeError, err: err}
			}
			paths := workspace.Resolve(root)
			if !workspace.Exists(paths) {
				return exitError{code: exitIndexMissing, err: fmt.Errorf("index not found; run `scry index`")}
			}
			store, err := metadata.Open(paths.IndexDBPath)
			if err != nil {
				return exitError{code: exitRuntimeError, err: err}
			}
			engine := search.New(store)
			query := strings.Join(args, " ")
			results, err := engine.Search(query, limit)
			if err != nil {
				return exitError{code: exitRuntimeError, err: err}
			}
			jsonOut, _ := cmd.Flags().GetBool("json")
			if len(results) == 0 {
				if jsonOut {
					_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
						"type":    "summary",
						"results": 0,
					})
				} else {
					fmt.Fprintln(os.Stdout, "no results")
				}
				return exitError{code: exitNoResults, silent: true}
			}
			for i, r := range results {
				snippet := formatSnippet(r.Chunk.Content, 200)
				if jsonOut {
					_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
						"type":       "result",
						"rank":       i + 1,
						"score":      r.Score,
						"path":       r.Chunk.FilePath,
						"start_line": r.Chunk.StartLine,
						"end_line":   r.Chunk.EndLine,
						"snippet":    snippet,
					})
				} else {
					fmt.Fprintf(os.Stdout, "%d. %s:%d-%d (score %.2f)\n", i+1, r.Chunk.FilePath, r.Chunk.StartLine, r.Chunk.EndLine, r.Score)
					fmt.Fprintf(os.Stdout, "   %s\n", snippet)
				}
			}
			if jsonOut {
				_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
					"type":    "summary",
					"results": len(results),
				})
			}
			return nil
		},
	}
	addCommonFlags(cmd)
	cmd.Flags().IntVar(&limit, "limit", 20, "max results")
	return cmd
}

func formatSnippet(text string, max int) string {
	s := strings.TrimSpace(text)
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
