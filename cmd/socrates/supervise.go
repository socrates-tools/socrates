package main

import (
	"os"
	"path/filepath"

	"github.com/socrates-devtool/socrates/internal/agent"
	"github.com/socrates-devtool/socrates/internal/config"
	"github.com/socrates-devtool/socrates/internal/reviewer"
	"github.com/socrates-devtool/socrates/internal/runlog"
	"github.com/spf13/cobra"
)

func newSuperviseCommand() *cobra.Command {
	return &cobra.Command{
		Use:    reviewer.SuperviseCommand + " <run>",
		Short:  superviseShort,
		Args:   cobra.ExactArgs(1),
		Hidden: true,
		RunE: func(command *cobra.Command, arguments []string) error {
			settings, err := config.Load()
			if err != nil {
				return err
			}
			store, err := openStore()
			if err != nil {
				return err
			}

			name := arguments[0]
			events, err := runlog.Open(filepath.Join(store.Path(name), runlog.FileName), settings.LogCeiling)
			if err != nil {
				return err
			}
			defer events.Close()

			return reviewer.Supervise(command.Context(), reviewer.Supervision{
				Store:    store,
				Name:     name,
				Process:  os.Getpid(),
				Ceiling:  settings.DiffCeiling,
				Recorder: events,
				NewAgent: func(recorder runlog.Recorder) (agent.Agent, error) {
					return agent.NewClaude(agent.Options{
						Recorder: recorder,
						Routing:  settings.Roles,
						Deadline: settings.AgentDeadline,
					})
				},
			})
		},
	}
}
