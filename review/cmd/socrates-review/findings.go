package main

import (
	"github.com/socrates-tools/socrates/review/report"
	"github.com/spf13/cobra"
)

func newFindingsCommand(code *int) *cobra.Command {
	var asJSON bool

	command := &cobra.Command{
		Use:   "findings [run]",
		Short: findingsShort,
		Args:  cobra.ArbitraryArgs,
		RunE: func(command *cobra.Command, arguments []string) error {
			store, name, err := chosenRun(command.Context(), arguments)
			if err != nil {
				return err
			}

			reported, err := report.FindingsOf(store, name)
			if err != nil {
				return err
			}

			if err := print(command, asJSON, reported, reported.Lines()); err != nil {
				return err
			}
			*code = reported.Code()
			return nil
		},
	}

	command.Flags().BoolVar(&asJSON, "json", false, jsonUsage)
	return command
}
