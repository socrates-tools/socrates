package main

import (
	"fmt"

	"github.com/socrates-tools/socrates-review/internal/reviewer"
	"github.com/spf13/cobra"
)

func newCancelCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "cancel [run]",
		Short: cancelShort,
		Args:  cobra.ArbitraryArgs,
		RunE: func(command *cobra.Command, arguments []string) error {
			store, name, err := chosenRun(command.Context(), arguments)
			if err != nil {
				return err
			}

			before, err := store.Load(name)
			if err != nil {
				return err
			}
			finishedAlready := before.State.Terminal()

			run, err := reviewer.Cancel(command.Context(), store, name)
			if err != nil {
				return err
			}

			out := command.OutOrStdout()
			if finishedAlready {
				fmt.Fprintf(out, alreadyEnded+"\n", run.Name, before.State)
				return nil
			}
			fmt.Fprintf(out, stoppedIt+"\n", run.Name)
			return nil
		},
	}
}
