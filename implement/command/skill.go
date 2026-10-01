package command

import (
	"fmt"

	"github.com/socrates-tools/socrates/implement/config"
	"github.com/socrates-tools/socrates/implement/skill"
	"github.com/spf13/cobra"
)

func newSkillCommand() *cobra.Command {
	var agentsTarget string

	command := &cobra.Command{
		Use:   "skill",
		Short: skillShort,
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			stated, err := config.Load()
			if err != nil {
				return err
			}

			if agentsTarget == "" {
				chosen, err := skill.AgentsDirectory()
				if err != nil {
					return err
				}
				agentsTarget = chosen
			}

			path, err := skill.Install(agentsTarget, stated)
			if err != nil {
				return err
			}
			fmt.Fprintf(command.OutOrStdout(), installedTheAgent+"\n", skill.AgentName, path)
			return nil
		},
	}

	command.Flags().StringVar(&agentsTarget, "agents-to", "", agentsToUsage)
	return command
}
