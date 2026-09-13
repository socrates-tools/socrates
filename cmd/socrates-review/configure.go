package main

import (
	"github.com/socrates-tools/socrates-review/internal/configure"
	"github.com/socrates-tools/socrates-review/internal/refusal"
	"github.com/spf13/cobra"
)

func newConfigureCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "configure",
		Short: configureShort,
		Args:  cobra.ArbitraryArgs,
		RunE: func(command *cobra.Command, arguments []string) error {
			if len(arguments) > 0 {
				return refusal.New(configureTakesNoArguments, arguments[0])
			}
			return configure.Run(command.Context(), command.InOrStdin(), command.OutOrStdout())
		},
	}
}
