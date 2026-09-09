package main

import (
	"fmt"

	"github.com/socrates-tools/socrates-review/internal/skill"
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
				chosen, err := skill.UserDirectory()
				if err != nil {
					return err
				}
				target = chosen
			}

			written, err := skill.Install(target)
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
