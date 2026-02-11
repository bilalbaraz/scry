package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"scry/pkg/indexer"
)

func newIndexCmd() *cobra.Command {
	var (
		clean        bool
		noEmbeddings bool
	)
	cmd := &cobra.Command{
		Use:   "index",
		Short: "Create or update local indexes",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := os.Getwd()
			if err != nil {
				return exitError{code: exitRuntimeError, err: err}
			}
			root, err = filepath.Abs(root)
			if err != nil {
				return exitError{code: exitRuntimeError, err: err}
			}
			jsonOut, _ := cmd.Flags().GetBool("json")
			opts := indexer.Options{
				Root:         root,
				Clean:        clean,
				NoEmbeddings: noEmbeddings,
				JSON:         jsonOut,
			}
			emit := func(p indexer.Progress) {
				if jsonOut {
					enc := json.NewEncoder(os.Stdout)
					_ = enc.Encode(p)
					return
				}
				switch p.Stage {
				case "scan":
					fmt.Fprintf(os.Stdout, "scan: %d files\\n", p.FilesTotal)
				case "index":
					fmt.Fprintf(os.Stdout, "indexed: %s (%d chunks)\\n", p.File, p.Chunks)
				}
			}
			summary, err := indexer.Run(opts, emit)
			if err != nil {
				return exitError{code: exitRuntimeError, err: err}
			}
			if jsonOut {
				_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
					"type":           "summary",
					"files_indexed":  summary.FilesIndexed,
					"chunks_indexed": summary.ChunksIndexed,
				})
			} else {
				fmt.Fprintf(os.Stdout, "done: %d files, %d chunks\\n", summary.FilesIndexed, summary.ChunksIndexed)
			}
			return nil
		},
	}
	addCommonFlags(cmd)
	cmd.Flags().BoolVar(&clean, "clean", false, "rebuild index from scratch")
	cmd.Flags().BoolVar(&noEmbeddings, "no-embeddings", false, "skip embeddings")
	return cmd
}
