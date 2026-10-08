package command

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/socrates-tools/socrates/agora/config"
	"github.com/socrates-tools/socrates/agora/surface"
	"github.com/socrates-tools/socrates/internal/exit"
	"github.com/socrates-tools/socrates/internal/refusal"
	"github.com/spf13/cobra"
)

func newPollCommand(code *int) *cobra.Command {
	var session string

	command := &cobra.Command{
		Use:   "poll",
		Short: pollShort,
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if session == "" {
				return refusal.New(nameTheSession, "poll")
			}

			open, err := whereTheSurfaceStands(command)
			if err != nil {
				return err
			}

			read, arrived, err := surface.Reach(open).Await(command.Context(), session, surface.PollDeadline)
			if err != nil {
				return err
			}
			if !arrived {
				fmt.Fprintln(command.OutOrStdout(), nobodyAnswered)
				*code = exit.SomethingAsked
				return nil
			}

			written, err := json.MarshalIndent(read, "", "  ")
			if err != nil {
				return fmt.Errorf(writingTheAnswer, err)
			}
			fmt.Fprintln(command.OutOrStdout(), string(bytes.TrimSpace(written)))
			return nil
		},
	}

	command.Flags().StringVar(&session, sessionFlag, "", sessionUsage)
	return command
}

func newEndCommand() *cobra.Command {
	var session string

	command := &cobra.Command{
		Use:   "end",
		Short: endShort,
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if session == "" {
				return refusal.New(nameTheSession, "end")
			}

			open, err := whereTheSurfaceStands(command)
			if err != nil {
				return err
			}
			return surface.Reach(open).End(command.Context(), session)
		},
	}

	command.Flags().StringVar(&session, sessionFlag, "", sessionUsage)
	return command
}

func newStopCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: stopShort,
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			open, err := whereTheSurfaceStands(command)
			if err != nil {
				return err
			}
			if err := surface.Reach(open).Stop(command.Context()); err != nil {
				return err
			}
			return surface.ForgetRecord()
		},
	}
}

func whereTheSurfaceStands(command *cobra.Command) (surface.Record, error) {
	file, err := config.Path()
	if err != nil {
		return surface.Record{}, err
	}
	loaded, err := config.LoadFrom(file)
	if err != nil {
		return surface.Record{}, err
	}

	found, err := surface.Locate(command.Context(), loaded.Surface.Bind)
	if err != nil {
		return surface.Record{}, err
	}
	return found.Reached()
}
