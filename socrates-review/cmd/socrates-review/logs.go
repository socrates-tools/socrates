package main

import (
	"github.com/socrates-tools/socrates/socrates-review/follow"
	"github.com/socrates-tools/socrates/socrates-review/presentation"
	"github.com/socrates-tools/socrates/socrates-review/report"
	"github.com/socrates-tools/socrates/socrates-review/terminal"
	"github.com/spf13/cobra"
)

func newLogsCommand(code *int) *cobra.Command {
	return &cobra.Command{
		Use:   "logs [run]",
		Short: logsShort,
		Args:  cobra.ArbitraryArgs,
		RunE: func(command *cobra.Command, arguments []string) error {
			store, name, err := chosenRun(command.Context(), arguments)
			if err != nil {
				return err
			}

			status, err := report.StatusOf(store, name)
			if err != nil {
				return err
			}
			held, err := follow.Held(store, name)
			if err != nil {
				return err
			}

			shown := append(presentation.Lines(held), status.Lines()...)
			if err := terminal.For(command.OutOrStdout(), command.ErrOrStderr()).Write(command.Context(), shown); err != nil {
				return err
			}
			*code = status.Code()
			return nil
		},
	}
}
