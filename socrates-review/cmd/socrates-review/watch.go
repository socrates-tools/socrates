package main

import (
	"fmt"
	"io"

	"github.com/socrates-tools/socrates/socrates-review/follow"
	"github.com/socrates-tools/socrates/socrates-review/presentation"
	"github.com/socrates-tools/socrates/socrates-review/report"
	"github.com/socrates-tools/socrates/socrates-review/runlog"
	"github.com/socrates-tools/socrates/socrates-review/runstate"
	"github.com/spf13/cobra"
)

func newWatchCommand(code *int) *cobra.Command {
	return &cobra.Command{
		Use:   "watch [run]",
		Short: watchShort,
		Args:  cobra.ArbitraryArgs,
		RunE: func(command *cobra.Command, arguments []string) error {
			store, name, err := chosenRun(command.Context(), arguments)
			if err != nil {
				return err
			}

			out := command.OutOrStdout()
			following := follow.Following{Store: store, Name: name}
			sink := follow.Sink{
				Events: func(events []runlog.Event) error {
					return printLines(out, presentation.Lines(events))
				},
				Waiting: func(state runstate.State) error {
					_, err := fmt.Fprintf(out, stillStarting+"\n", name, state)
					return err
				},
			}
			if _, err := following.Until(command.Context(), sink); err != nil {
				return err
			}

			status, err := report.StatusOf(store, name)
			if err != nil {
				return err
			}
			if err := printLines(out, status.Lines()); err != nil {
				return err
			}
			*code = status.Code()
			return nil
		},
	}
}

func printLines(out io.Writer, lines []string) error {
	for _, line := range lines {
		if _, err := fmt.Fprintln(out, line); err != nil {
			return err
		}
	}
	return nil
}
