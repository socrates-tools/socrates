package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rui-ktei/socrates/internal/reviewer"
	"github.com/rui-ktei/socrates/internal/runstore"
	"github.com/rui-ktei/socrates/internal/state"
	"github.com/spf13/cobra"
)

func newReviewCommand() *cobra.Command {
	var intent string

	command := &cobra.Command{
		Use:   "review",
		Short: "Freeze this working tree and have somebody who did not write it review it",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			store, err := openStore()
			if err != nil {
				return err
			}

			directory, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("locating the working directory: %w", err)
			}

			binary, err := os.Executable()
			if err != nil {
				return fmt.Errorf("locating scrat itself: %w", err)
			}

			run, err := reviewer.Begin(command.Context(), reviewer.Beginning{
				Directory: directory,
				Intent:    intent,
				Store:     store,
				Binary:    binary,
				StartedAt: time.Now(),
			})
			if err != nil {
				return err
			}

			fmt.Fprintln(command.OutOrStdout(), run.Name)
			return nil
		},
	}

	command.Flags().StringVar(&intent, "intent", "", "what the person asked for, in their words")
	return command
}

func openStore() (runstore.Store, error) {
	runs, err := state.Runs()
	if err != nil {
		return runstore.Store{}, err
	}
	return runstore.New(runs), nil
}
