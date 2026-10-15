package command

import (
	"fmt"

	skillkit "github.com/socrates-tools/socrates/internal/skill"
	"github.com/socrates-tools/socrates/review/skill"
	"github.com/spf13/cobra"
)

func newSkillCommand() *cobra.Command {
	var target, host string

	command := &cobra.Command{
		Use:   "skill",
		Short: skillShort,
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			hosts, err := skillkit.Selected(host, skillkit.Directory{Flag: skillkit.ToFlag, Value: target})
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
			return nil
		},
	}

	command.Flags().StringVar(&host, skillkit.HostFlag, "", skillkit.HostUsage())
	command.Flags().StringVar(&target, skillkit.ToFlag, "", toUsage)
	return command
}
