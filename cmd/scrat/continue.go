package main

import (
	"fmt"

	"github.com/rui-ktei/socrates/internal/reviewer"
	"github.com/spf13/cobra"
)

func newContinueCommand(code *int) *cobra.Command {
	var wait bool

	command := &cobra.Command{
		Use:   "continue [run]",
		Short: continueShort,
		Args:  cobra.ArbitraryArgs,
		RunE: func(command *cobra.Command, arguments []string) error {
			store, name, err := chosenRun(command.Context(), arguments)
			if err != nil {
				return err
			}
			binary, err := ownBinary()
			if err != nil {
				return err
			}

			run, round, err := reviewer.Continue(command.Context(), reviewer.Continuation{
				Store:  store,
				Name:   name,
				Binary: binary,
			})
			if err != nil {
				return err
			}

			out := command.OutOrStdout()
			fmt.Fprintln(out, run.Name)
			fmt.Fprintf(out, roundStarted+"\n", round)

			if !wait {
				fmt.Fprintf(out, followIt+"\n", run.Name)
				return nil
			}

			status, err := awaitVerdict(command.Context(), store, run.Name)
			if err != nil {
				return err
			}
			for _, line := range status.Lines() {
				fmt.Fprintln(out, line)
			}
			*code = status.Code()
			return nil
		},
	}

	command.Flags().BoolVar(&wait, "wait", false, waitUsage)
	return command
}
