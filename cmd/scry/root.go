package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"scry/pkg/config"
)

const (
	exitSuccess          = 0
	exitRuntimeError     = 1
	exitUsageError       = 2
	exitIndexMissing     = 3
	exitOfflineViolation = 4
	exitNoResults        = 5
)

type exitError struct {
	code   int
	err    error
	silent bool
}

func (e exitError) Error() string {
	if e.err == nil {
		return ""
	}
	return e.err.Error()
}

func execute() int {
	root := newRootCmd()
	root.SilenceErrors = true
	root.SilenceUsage = true
	if err := root.Execute(); err != nil {
		var ee exitError
		if errors.As(err, &ee) {
			if !ee.silent && ee.err != nil {
				fmt.Fprintln(os.Stderr, ee.err)
			}
			return ee.code
		}
		fmt.Fprintln(os.Stderr, err)
		return exitUsageError
	}
	return exitSuccess
}

func newRootCmd() *cobra.Command {
	var configPath string
	root := &cobra.Command{
		Use:   "scry",
		Short: "Local-first codebase memory engine",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return exitError{code: exitUsageError, err: err}
			}
			return exitError{code: exitUsageError, silent: true}
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Name() == "scry" {
				return nil
			}
			explicit := cmd.Flags().Changed("config")
			if err := loadConfig(configPath, explicit); err != nil {
				return exitError{code: exitRuntimeError, err: err}
			}
			return nil
		},
	}

	root.PersistentFlags().StringVarP(&configPath, "config", "c", "", "Config file (default: .memengine.yml)")

	root.AddCommand(newIndexCmd())
	root.AddCommand(newSearchCmd())
	root.AddCommand(newAskCmd())
	root.AddCommand(newStatusCmd())
	root.AddCommand(newImpactCmd())

	return root
}

func loadConfig(path string, required bool) error {
	_, err := config.Load(path, required)
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}
	return nil
}

func addCommonFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("json", false, "output JSON")
	cmd.Flags().Bool("quiet", false, "suppress progress output")
}

func runNotImplemented(cmd *cobra.Command, args []string) error {
	fmt.Fprintln(os.Stdout, "not implemented")
	return exitError{code: exitRuntimeError, silent: true}
}
