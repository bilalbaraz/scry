package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"scry/pkg/metadata"
	"scry/pkg/workspace"
)

func newStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show repo index health",
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
			jsonOut, _ := cmd.Flags().GetBool("json")
			if !workspace.Exists(paths) {
				if jsonOut {
					_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
						"type":  "status",
						"repo":  root,
						"index": "missing",
					})
					return nil
				}
				fmt.Fprintln(os.Stdout, "Index: missing")
				return exitError{code: exitIndexMissing, silent: true}
			}
			store, err := metadata.Open(paths.IndexDBPath)
			if err != nil {
				return exitError{code: exitRuntimeError, err: err}
			}
			stats, err := store.Stats()
			if err != nil {
				return exitError{code: exitRuntimeError, err: err}
			}
			if jsonOut {
				_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
					"type":   "status",
					"repo":   root,
					"index":  "present",
					"files":  stats.Files,
					"chunks": stats.Chunks,
					"terms":  stats.Terms,
				})
				return nil
			}
			fmt.Fprintf(os.Stdout, "Repo: %s\n", root)
			fmt.Fprintf(os.Stdout, "Index: present\n")
			fmt.Fprintf(os.Stdout, "Files indexed: %d\n", stats.Files)
			fmt.Fprintf(os.Stdout, "Chunks indexed: %d\n", stats.Chunks)
			fmt.Fprintf(os.Stdout, "Terms indexed: %d\n", stats.Terms)
			return nil
		},
	}
	addCommonFlags(cmd)
	return cmd
}
