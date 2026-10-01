package command

import (
	"fmt"

	"github.com/socrates-tools/socrates/implement/config"
	"github.com/socrates-tools/socrates/internal/refusal"
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

			for _, line := range config.ReportAt(path, loaded).Lines() {
				fmt.Fprintln(command.OutOrStdout(), line)
			}
			return nil
		},
	}
}
