package main

import "github.com/spf13/cobra"

func newSearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Hybrid search across indexes",
		RunE:  runNotImplemented,
	}
	addCommonFlags(cmd)
	return cmd
}
