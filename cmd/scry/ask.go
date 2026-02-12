package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	askquery "scry/internal/query/ask"
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
			terms := askquery.TokenizeQuery(question)
			var candidates []askquery.Chunk
			for _, r := range results {
				candidates = append(candidates, askquery.Chunk{
					ID:        r.Chunk.ID,
					FilePath:  r.Chunk.FilePath,
					StartLine: r.Chunk.StartLine,
					EndLine:   r.Chunk.EndLine,
					Text:      r.Chunk.Content,
					Score:     r.Score,
				})
			}

			decision := askquery.BuildPipeline(candidates, terms, askquery.AskOptions{
				MaxEvidence:  2,
				SnippetChars: 240,
				MinScore:     1.0,
			})

			if len(decision.Evidence) == 0 {
				hint := "No relevant evidence found in indexed chunks."
				if jsonOut {
					_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
						"type":   "answer",
						"text":   "I don't know.",
						"reason": decision.Reason,
						"hint":   hint,
					})
				} else {
					fmt.Fprintln(os.Stdout, "I don't know.")
					fmt.Fprintln(os.Stdout, hint)
				}
				return nil
			}

			answerHeader := askquery.AnswerHeader(decision.Evidence)
			if jsonOut {
				_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
					"type": "answer",
					"text": answerHeader,
				})
				for i, ev := range decision.Evidence {
					_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
						"type":       "evidence",
						"id":         i + 1,
						"snippet":    ev.Snippet,
						"path":       ev.Chunk.FilePath,
						"start_line": ev.Chunk.StartLine,
						"end_line":   ev.Chunk.EndLine,
					})
				}
				_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
					"type": "summary",
					"k":    len(decision.Evidence),
				})
				return nil
			}

			fmt.Fprintln(os.Stdout, "Answer:")
			fmt.Fprintln(os.Stdout, answerHeader)
			fmt.Fprintln(os.Stdout, "")
			fmt.Fprintln(os.Stdout, "Evidence:")
			for i, ev := range decision.Evidence {
				fmt.Fprintf(os.Stdout, "[%d] %s:%d-%d\n", i+1, ev.Chunk.FilePath, ev.Chunk.StartLine, ev.Chunk.EndLine)
				fmt.Fprintf(os.Stdout, "    %s\n", ev.Snippet)
			}
			return nil
		},
	}
	addCommonFlags(cmd)
	cmd.Flags().IntVar(&limit, "k", 6, "number of context chunks")
	return cmd
}
