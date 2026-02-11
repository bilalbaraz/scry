package main

import "github.com/spf13/cobra"

func newImpactCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "impact <path|range|commit>",
		Short: "Show change impact (placeholder)",
		RunE:  runNotImplemented,
	}
	addCommonFlags(cmd)
	return cmd
}
