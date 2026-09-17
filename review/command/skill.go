package command

import (
	"fmt"

	skillkit "github.com/socrates-tools/socrates/internal/skill"
	"github.com/socrates-tools/socrates/review/skill"
	"github.com/spf13/cobra"
)

func newSkillCommand() *cobra.Command {
	var target string

	command := &cobra.Command{
		Use:   "skill",
		Short: skillShort,
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
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

			out := command.OutOrStdout()
			for _, installed := range written {
				fmt.Fprintf(out, installedTheSkill+"\n", installed.Name, installed.Path)
			}
			return nil
		},
	}

	command.Flags().StringVar(&target, "to", "", toUsage)
	return command
}
