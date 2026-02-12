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

func newAskCmd() *cobra.Command {
	var limit int
	cmd := &cobra.Command{
		Use:   "ask <question>",
		Short: "RAG-style answers with citations",
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
			question := strings.Join(args, " ")
			results, err := engine.Search(question, limit)
			if err != nil {
				return exitError{code: exitRuntimeError, err: err}
			}
			jsonOut, _ := cmd.Flags().GetBool("json")
			if len(results) == 0 {
				if jsonOut {
					_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
						"type":   "answer",
						"text":   "I don't know.",
						"reason": "no_evidence",
					})
				} else {
					fmt.Fprintln(os.Stdout, "I don't know.")
				}
				return nil
			}

			answer := buildExtractiveAnswer(results)
			if jsonOut {
				_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
					"type": "answer",
					"text": answer,
				})
				for i, r := range results {
					_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
						"type":       "citation",
						"id":         i + 1,
						"path":       r.Chunk.FilePath,
						"start_line": r.Chunk.StartLine,
						"end_line":   r.Chunk.EndLine,
					})
				}
				_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
					"type": "summary",
					"k":    len(results),
				})
				return nil
			}

			fmt.Fprintln(os.Stdout, "Answer:")
			fmt.Fprintln(os.Stdout, answer)
			fmt.Fprintln(os.Stdout, "")
			fmt.Fprintln(os.Stdout, "Citations:")
			for i, r := range results {
				fmt.Fprintf(os.Stdout, "[%d] %s:%d-%d\n", i+1, r.Chunk.FilePath, r.Chunk.StartLine, r.Chunk.EndLine)
			}
			return nil
		},
	}
	addCommonFlags(cmd)
	cmd.Flags().IntVar(&limit, "k", 6, "number of context chunks")
	return cmd
}

func buildExtractiveAnswer(results []search.Result) string {
	if len(results) == 0 {
		return "I don't know."
	}
	var parts []string
	for i, r := range results {
		if i >= 3 {
			break
		}
		parts = append(parts, formatSnippet(r.Chunk.Content, 240))
	}
	return strings.Join(parts, "\n\n")
}
