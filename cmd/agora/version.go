package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

const recordedUncommitted = "yes"

var (
	commit      string
	uncommitted string
)

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: versionShort,
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			fmt.Fprintln(command.OutOrStdout(), described())
			return nil
		},
	}
}

func described() string {
	if commit == "" {
		return recordedNoCommit
	}
	if uncommitted == recordedUncommitted {
		return fmt.Sprintf(builtFromDirtyTree, commit)
	}
	return fmt.Sprintf(builtFrom, commit)
}
