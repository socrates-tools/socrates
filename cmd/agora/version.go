package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

const recordedUncommitted = "yes"

var (
	version     string
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
	if version == "" {
		if uncommitted == recordedUncommitted {
			return fmt.Sprintf(builtBetweenReleasesFromDirtyTree, commit)
		}
		return fmt.Sprintf(builtBetweenReleases, commit)
	}
	if uncommitted == recordedUncommitted {
		return fmt.Sprintf(builtFromDirtyTree, version, commit)
	}
	return fmt.Sprintf(builtFrom, version, commit)
}
