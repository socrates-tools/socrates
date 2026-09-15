package main

import (
	"fmt"

	"github.com/socrates-tools/socrates/internal/version"
	"github.com/spf13/cobra"
)

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: versionShort,
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			fmt.Fprintln(command.OutOrStdout(), version.Reported(toolName))
			return nil
		},
	}
}
