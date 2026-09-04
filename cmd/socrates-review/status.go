package main

import (
	"context"
	"fmt"
	"os"

	"github.com/socrates-tools/socrates-review/internal/refusal"
	"github.com/socrates-tools/socrates-review/internal/report"
	"github.com/socrates-tools/socrates-review/internal/runstore"
	"github.com/socrates-tools/socrates-review/internal/snapshot"
	"github.com/spf13/cobra"
)

func newStatusCommand(code *int) *cobra.Command {
	var asJSON bool

	command := &cobra.Command{
		Use:   "status [run]",
		Short: statusShort,
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

			if err := print(command, asJSON, status, status.Lines()); err != nil {
				return err
			}
			*code = status.Code()
			return nil
		},
	}

	command.Flags().BoolVar(&asJSON, "json", false, jsonUsage)
	return command
}

func print(command *cobra.Command, asJSON bool, reported any, lines []string) error {
	out := command.OutOrStdout()

	if asJSON {
		encoded, err := report.AsJSON(reported)
		if err != nil {
			return err
		}
		fmt.Fprint(out, encoded)
		return nil
	}

	for _, line := range lines {
		fmt.Fprintln(out, line)
	}
	return nil
}

func chosenRun(ctx context.Context, arguments []string) (runstore.Store, string, error) {
	if len(arguments) > 1 {
		return runstore.Store{}, "", refusal.New(oneRunAtATime, arguments[1])
	}

	store, err := openStore()
	if err != nil {
		return runstore.Store{}, "", err
	}

	identifier := ""
	if len(arguments) == 1 {
		identifier = arguments[0]
	}
	if identifier != "" {
		return store, identifier, nil
	}

	directory, err := os.Getwd()
	if err != nil {
		return runstore.Store{}, "", fmt.Errorf(locatingHere, err)
	}
	root, err := snapshot.Root(ctx, directory)
	if err != nil {
		return runstore.Store{}, "", err
	}

	name, err := store.Select("", root)
	if err != nil {
		return runstore.Store{}, "", err
	}
	return store, name, nil
}
