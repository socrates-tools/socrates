package main

import (
	"github.com/socrates-tools/socrates/internal/refusal"
	"github.com/socrates-tools/socrates/review/config"
	"github.com/socrates-tools/socrates/review/report"
	"github.com/socrates-tools/socrates/review/terminal"
	"github.com/spf13/cobra"
)

func newConfigCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: configShort,
		Args:  cobra.ArbitraryArgs,
		RunE: func(command *cobra.Command, arguments []string) error {
			if len(arguments) > 0 {
				return refusal.New(configTakesNoArguments, arguments[0])
			}

			path, err := config.Path()
			if err != nil {
				return err
			}
			loaded, err := config.LoadFrom(path)
			if err != nil {
				return err
			}

			shown := report.ConfigurationAt(path, loaded).Lines()
			return terminal.For(command.OutOrStdout(), command.ErrOrStderr()).Write(command.Context(), shown)
		},
	}
}
