package command

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/socrates-tools/socrates/review/agent"
	"github.com/socrates-tools/socrates/review/config"
	"github.com/socrates-tools/socrates/review/reviewer"
	"github.com/socrates-tools/socrates/review/runlog"
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
			run, err := store.Load(name)
			if err != nil {
				return errors.Join(err, events.Close())
			}
			vendor, known := agent.Named(run.Vendor())
			if !known {
				return errors.Join(reviewer.ReviewerCLIIsSupported(run.Vendor()), events.Close())
			}

			supervised := reviewer.Supervise(command.Context(), reviewer.Supervision{
				Store:    store,
				Name:     name,
				Process:  os.Getpid(),
				Ceiling:  settings.DiffCeiling,
				Recorder: events,
				Vendor:   vendor,
				NewAgent: func(recorder runlog.Recorder) (agent.Agent, error) {
					return agent.New(vendor.Name, agent.Options{
						Recorder: recorder,
						Routing:  settings.Roles,
						Deadline: settings.AgentDeadline,
					})
				},
			})
			return errors.Join(supervised, events.Close())
		},
	}
}
