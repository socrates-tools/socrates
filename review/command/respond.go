package command

import (
	"fmt"
	"io"
	"strings"

	"github.com/socrates-tools/socrates/internal/refusal"
	"github.com/socrates-tools/socrates/review/dialogue"
	"github.com/socrates-tools/socrates/review/report"
	"github.com/socrates-tools/socrates/review/reviewer"
	"github.com/spf13/cobra"
)

const fromStandardInput = "-"

func newRespondCommand(code *int) *cobra.Command {
	var finding string
	var fixed bool
	var rejected bool
	var decided bool
	var forwarded bool
	var dropped bool
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

			chosen, err := oneResolution(map[dialogue.Resolution]bool{
				dialogue.Fixed:     fixed,
				dialogue.Rejected:  rejected,
				dialogue.Decided:   decided,
				dialogue.Forwarded: forwarded,
				dialogue.Dropped:   dropped,
			})
			if err != nil {
				return err
			}
			if chosen == dialogue.Fixed && reason != "" {
				return refusal.New(reasonBelongsToAClaim)
			}
			if chosen == dialogue.Forwarded && reason != "" {
				return refusal.New(aForwardCarriesTheFinding, dialogue.Forwarded)
			}
			if chosen != dialogue.Fixed && note != "" {
				return refusal.New(noteBelongsToAFix, chosen)
			}

			response := dialogue.Response{Finding: strings.TrimSpace(finding), Resolution: chosen, Reason: reason}
			if chosen == dialogue.Fixed {
				response.Reason = note
			}

			if response.Reason == fromStandardInput {
				read, err := readStandardInput(command)
				if err != nil {
					return err
				}
				response.Reason = read
			}

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
	command.Flags().BoolVar(&decided, "decided", false, decidedUsage)
	command.Flags().BoolVar(&forwarded, "forwarded", false, forwardedUsage)
	command.Flags().BoolVar(&dropped, "dropped", false, droppedUsage)
	command.Flags().StringVar(&reason, "reason", "", reasonUsage)
	command.Flags().StringVar(&note, "note", "", noteUsage)
	return command
}

func oneResolution(chosen map[dialogue.Resolution]bool) (dialogue.Resolution, error) {
	var named []dialogue.Resolution
	for _, resolution := range dialogue.Resolutions() {
		if chosen[resolution] {
			named = append(named, resolution)
		}
	}
	if len(named) != 1 {
		return "", refusal.New(oneAnswerAtATime, len(named), dialogue.Flags(dialogue.Resolutions()))
	}
	return named[0], nil
}

func readStandardInput(command *cobra.Command) (string, error) {
	read, err := io.ReadAll(command.InOrStdin())
	if err != nil {
		return "", fmt.Errorf(readingStandardInput, err)
	}
	return string(read), nil
}
