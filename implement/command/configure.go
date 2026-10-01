package command

import (
	"github.com/socrates-tools/socrates/implement/configure"
	"github.com/socrates-tools/socrates/internal/refusal"
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

			repository, err := repositoryOfTheWorkingDirectory()
			if err != nil {
				return err
			}
			return configure.Run(command.Context(), command.InOrStdin(), command.OutOrStdout(), repository)
		},
	}
}
