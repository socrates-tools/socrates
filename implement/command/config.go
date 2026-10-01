package command

import (
	"errors"
	"fmt"
	"os"

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

			repository, err := repositoryOfTheWorkingDirectory()
			if err != nil {
				return err
			}

			for _, line := range config.ReportIn(path, loaded, repository).Lines() {
				fmt.Fprintln(command.OutOrStdout(), line)
			}
			return nil
		},
	}
}

func repositoryOfTheWorkingDirectory() (string, error) {
	here, err := os.Getwd()
	if err != nil {
		return "", err
	}
	repository, err := config.RepositoryAt(here)
	if errors.Is(err, config.ErrNotARepository) || errors.Is(err, config.ErrBareRepository) {
		return "", nil
	}
	return repository, err
}
