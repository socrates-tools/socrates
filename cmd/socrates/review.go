package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ktei/socrates/internal/config"
	"github.com/ktei/socrates/internal/refusal"
	"github.com/ktei/socrates/internal/report"
	"github.com/ktei/socrates/internal/reviewer"
	"github.com/ktei/socrates/internal/runstate"
	"github.com/ktei/socrates/internal/runstore"
	"github.com/ktei/socrates/internal/state"
	"github.com/spf13/cobra"
)

const verdictPoll = 50 * time.Millisecond

func newReviewCommand(code *int) *cobra.Command {
	var intent string
	var wait bool

	command := &cobra.Command{
		Use:   "review",
		Short: reviewShort,
		Args:  cobra.ArbitraryArgs,
		RunE: func(command *cobra.Command, arguments []string) error {
			if len(arguments) > 0 {
				return refusal.New(reviewTakesNoArguments, arguments[0])
			}

			settings, err := config.Load()
			if err != nil {
				return err
			}
			store, err := openStore()
			if err != nil {
				return err
			}

			directory, err := os.Getwd()
			if err != nil {
				return fmt.Errorf(locatingHere, err)
			}
			binary, err := ownBinary()
			if err != nil {
				return err
			}

			run, err := reviewer.Begin(command.Context(), reviewer.Beginning{
				Directory:       directory,
				Intent:          intent,
				Store:           store,
				Binary:          binary,
				StartedAt:       time.Now(),
				Retention:       settings.RunRetention,
				GuidanceCeiling: int64(settings.GuidanceCeiling),
			})
			if err != nil {
				return err
			}

			out := command.OutOrStdout()
			fmt.Fprintln(out, run.Name)

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

	command.Flags().StringVar(&intent, "intent", "", intentUsage)
	command.Flags().BoolVar(&wait, "wait", false, waitUsage)
	return command
}

func awaitVerdict(ctx context.Context, store runstore.Store, name string) (report.Status, error) {
	for {
		status, err := report.StatusOf(store, name)
		if err != nil {
			return report.Status{}, err
		}
		if status.State.Terminal() || status.State == runstate.AwaitingResponse {
			return status, nil
		}

		select {
		case <-ctx.Done():
			return report.Status{}, ctx.Err()
		case <-time.After(verdictPoll):
		}
	}
}

func ownBinary() (string, error) {
	binary, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf(locatingSocrates, err)
	}
	return binary, nil
}

func openStore() (runstore.Store, error) {
	runs, err := state.Runs()
	if err != nil {
		return runstore.Store{}, err
	}
	return runstore.New(runs), nil
}
