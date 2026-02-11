package main

import "github.com/spf13/cobra"

func newStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show repo index health",
		RunE:  runNotImplemented,
	}
	addCommonFlags(cmd)
	return cmd
}
