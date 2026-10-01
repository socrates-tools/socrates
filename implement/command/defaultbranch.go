package command

import (
	"errors"
	"fmt"
	"os"

	"github.com/socrates-tools/socrates/implement/config"
	"github.com/socrates-tools/socrates/internal/refusal"
	"github.com/spf13/cobra"
)

func newDefaultBranchCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "default-branch <branch>",
		Short: defaultBranchShort,
		Args:  cobra.ArbitraryArgs,
		RunE: func(command *cobra.Command, arguments []string) error {
			if len(arguments) != 1 {
				return refusal.New(defaultBranchTakesOneBranch, len(arguments))
			}
			branch := arguments[0]

			here, err := os.Getwd()
			if err != nil {
				return err
			}
			repository, err := config.RepositoryAt(here)
			if errors.Is(err, config.ErrNotARepository) {
				return refusal.New(defaultBranchNeedsARepository, here)
			}
			if errors.Is(err, config.ErrBareRepository) {
				return refusal.New(defaultBranchNeedsAWorktree, here)
			}
			if err != nil {
				return err
			}
			if !config.HoldsBranch(repository, branch) {
				return refusal.New(repositoryHoldsNoSuchBranch, repository, branch)
			}

			path, err := config.Path()
			if err != nil {
				return err
			}
			if _, err := config.LoadFrom(path); err != nil {
				return err
			}
			if err := config.RememberDefaultBranch(path, repository, branch); err != nil {
				return err
			}
			fmt.Fprintf(command.OutOrStdout(), defaultBranchRemembered+"\n", repository, branch, path)
			return nil
		},
	}
}
