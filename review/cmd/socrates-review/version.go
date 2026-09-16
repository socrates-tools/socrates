package main

import (
	"fmt"

	"github.com/socrates-tools/socrates/internal/refusal"
	"github.com/socrates-tools/socrates/internal/version"
	"github.com/spf13/cobra"
)

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: versionShort,
		Args:  cobra.ArbitraryArgs,
		RunE: func(command *cobra.Command, arguments []string) error {
			if len(arguments) > 0 {
				return refusal.New(versionTakesNoArguments, arguments[0])
			}

			fmt.Fprintln(command.OutOrStdout(), version.Reported(toolName))
			return nil
		},
	}
}
