package command

import (
	"fmt"

	"github.com/socrates-tools/socrates/implement/config"
	"github.com/socrates-tools/socrates/implement/skill"
	skillkit "github.com/socrates-tools/socrates/internal/skill"
	"github.com/spf13/cobra"
)

func newSkillCommand() *cobra.Command {
	var target, agentsTarget string

	command := &cobra.Command{
		Use:   "skill",
		Short: skillShort,
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			stated, err := config.Load()
			if err != nil {
				return err
			}

			if target == "" {
				chosen, err := skillkit.UserDirectory(skillInvocation)
				if err != nil {
					return err
				}
				target = chosen
			}

			written, err := skillkit.Install(target, skill.All(), skillInvocation)
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
			out := command.OutOrStdout()
			for _, installed := range written {
				fmt.Fprintf(out, installedTheSkill+"\n", installed.Name, installed.Path)
			}
			fmt.Fprintf(out, installedTheAgent+"\n", skill.AgentName, path)
			return nil
		},
	}

	command.Flags().StringVar(&target, "to", "", toUsage)
	command.Flags().StringVar(&agentsTarget, "agents-to", "", agentsToUsage)
	return command
}
