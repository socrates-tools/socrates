package main

import (
	"os"
	"path/filepath"

	"github.com/rui-ktei/socrates/internal/agent"
	"github.com/rui-ktei/socrates/internal/reviewer"
	"github.com/rui-ktei/socrates/internal/runlog"
	"github.com/spf13/cobra"
)

func newSuperviseCommand() *cobra.Command {
	return &cobra.Command{
		Use:    reviewer.SuperviseCommand + " <run>",
		Short:  "Conduct a review that scrat review started",
		Args:   cobra.ExactArgs(1),
		Hidden: true,
		RunE: func(command *cobra.Command, arguments []string) error {
			store, err := openStore()
			if err != nil {
				return err
			}

			name := arguments[0]
			events, err := runlog.Open(filepath.Join(store.Path(name), runlog.FileName))
			if err != nil {
				return err
			}
			defer events.Close()

			return reviewer.Supervise(command.Context(), reviewer.Supervision{
				Store:    store,
				Name:     name,
				Process:  os.Getpid(),
				Recorder: events,
				NewAgent: func(recorder runlog.Recorder) (agent.Agent, error) {
					return agent.NewClaude(agent.Options{Recorder: recorder})
				},
			})
		},
	}
}
