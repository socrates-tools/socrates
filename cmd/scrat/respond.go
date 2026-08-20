package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rui-ktei/socrates/internal/dialogue"
	"github.com/rui-ktei/socrates/internal/refusal"
	"github.com/rui-ktei/socrates/internal/report"
	"github.com/rui-ktei/socrates/internal/reviewer"
	"github.com/spf13/cobra"
)

const fromStandardInput = "-"

func newRespondCommand(code *int) *cobra.Command {
	var finding string
	var fixed bool
	var rejected bool
	var reason string
	var note string

	command := &cobra.Command{
		Use:   "respond [run]",
		Short: respondShort,
		Args:  cobra.ArbitraryArgs,
		RunE: func(command *cobra.Command, arguments []string) error {
			if strings.TrimSpace(finding) == "" {
				return refusal.New(nameTheFinding)
			}
			if fixed == rejected {
				return refusal.New(oneAnswerAtATime)
			}

			response := dialogue.Response{Finding: strings.TrimSpace(finding), Resolution: dialogue.Fixed, Reason: note}
			if rejected {
				response.Resolution = dialogue.Rejected
				response.Reason = reason
			}

			read, err := reasonFor(command, response)
			if err != nil {
				return err
			}
			response.Reason = read

			store, name, err := chosenRun(command.Context(), arguments)
			if err != nil {
				return err
			}

			recording, err := reviewer.Respond(store, name, response)
			if err != nil {
				return err
			}

			out := command.OutOrStdout()
			fmt.Fprintf(out, recordedIt+"\n", recording.Response.Finding, recording.Response.Resolution, recording.Round)
			if recording.Replaced {
				fmt.Fprintf(out, replacedIt+"\n", recording.Response.Finding)
			}
			if len(recording.Unanswered) == 0 {
				fmt.Fprintln(out, everythingIsAnswered)
			} else {
				fmt.Fprintf(out, stillToAnswer+"\n", len(recording.Unanswered), strings.Join(recording.Unanswered, ", "))
			}

			status, err := report.StatusOf(store, name)
			if err != nil {
				return err
			}
			*code = status.Code()
			return nil
		},
	}

	command.Flags().StringVar(&finding, "finding", "", findingUsage)
	command.Flags().BoolVar(&fixed, "fixed", false, fixedUsage)
	command.Flags().BoolVar(&rejected, "rejected", false, rejectedUsage)
	command.Flags().StringVar(&reason, "reason", "", reasonUsage)
	command.Flags().StringVar(&note, "note", "", noteUsage)
	return command
}

func reasonFor(command *cobra.Command, response dialogue.Response) (string, error) {
	if response.Reason == fromStandardInput {
		return readStandardInput(command)
	}
	if response.Reason != "" {
		return response.Reason, nil
	}
	if response.Resolution == dialogue.Rejected && standardInputIsPiped() {
		return readStandardInput(command)
	}
	return response.Reason, nil
}

func readStandardInput(command *cobra.Command) (string, error) {
	read, err := io.ReadAll(command.InOrStdin())
	if err != nil {
		return "", fmt.Errorf(readingStandardInput, err)
	}
	return string(read), nil
}

func standardInputIsPiped() bool {
	stated, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return stated.Mode()&os.ModeCharDevice == 0
}
