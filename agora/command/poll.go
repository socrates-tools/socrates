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
	var last bool

	command := &cobra.Command{
		Use:   "poll",
		Short: pollShort,
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if session == "" {
				return refusal.New(nameTheSession, "poll")
			}

			if last {
				return printTheLastCollected(command, session)
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
			return printAnswered(command, read)
		},
	}

	command.Flags().StringVar(&session, sessionFlag, "", sessionUsage)
	command.Flags().BoolVar(&last, lastFlag, false, lastUsage)
	return command
}

func printTheLastCollected(command *cobra.Command, session string) error {
	read, collected, err := surface.ReadAnswer(session)
	if err != nil {
		return err
	}
	if !collected {
		return refusal.New(nothingWasCollected, session, session)
	}
	return printAnswered(command, read)
}

func printAnswered(command *cobra.Command, read surface.Answered) error {
	written, err := json.MarshalIndent(read, "", "  ")
	if err != nil {
		return fmt.Errorf(writingTheAnswer, err)
	}
	fmt.Fprintln(command.OutOrStdout(), string(bytes.TrimSpace(written)))
	return nil
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
