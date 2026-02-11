package main

import "github.com/spf13/cobra"

func newAskCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ask <question>",
		Short: "RAG-style answers with citations",
		RunE:  runNotImplemented,
	}
	addCommonFlags(cmd)
	return cmd
}
