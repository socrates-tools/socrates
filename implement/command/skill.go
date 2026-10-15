package command

import (
	"fmt"

	"github.com/socrates-tools/socrates/implement/config"
	"github.com/socrates-tools/socrates/implement/skill"
	skillkit "github.com/socrates-tools/socrates/internal/skill"
	"github.com/spf13/cobra"
)

func newSkillCommand() *cobra.Command {
	var host, target, agentsTarget string

	command := &cobra.Command{
		Use:   "skill",
		Short: skillShort,
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			hosts, err := skillkit.Selected(host,
				skillkit.Directory{Flag: skillkit.ToFlag, Value: target},
				skillkit.Directory{Flag: skillkit.AgentsToFlag, Value: agentsTarget})
			if err != nil {
				return err
			}

			stated, err := config.Load()
			if err != nil {
				return err
			}

			written, err := skillkit.InstallFor(hosts, target, skillInvocation, skill.For)
			if err != nil {
				return err
			}

			out := command.OutOrStdout()
			for _, installed := range written {
				fmt.Fprintln(out, installed.Announcement())
			}

			for _, installing := range hosts {
				directory := agentsTarget
				if directory == "" {
					chosen, err := skill.AgentsDirectory(installing)
					if err != nil {
						return err
					}
					directory = chosen
				}
				path, err := skill.Install(installing, directory, stated)
				if err != nil {
					return err
				}
				fmt.Fprintf(out, installedTheAgent+"\n", skill.AgentName, skill.AgentKind(installing), path)
			}
			return nil
		},
	}

	command.Flags().StringVar(&host, skillkit.HostFlag, "", skillkit.HostUsage())
	command.Flags().StringVar(&target, skillkit.ToFlag, "", toUsage)
	command.Flags().StringVar(&agentsTarget, skillkit.AgentsToFlag, "", agentsToUsage)
	return command
}
