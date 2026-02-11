package main

import "github.com/spf13/cobra"

func newIndexCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "index",
		Short: "Create or update local indexes",
		RunE:  runNotImplemented,
	}
	addCommonFlags(cmd)
	return cmd
}
