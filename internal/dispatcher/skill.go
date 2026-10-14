package dispatcher

import (
	"github.com/socrates-tools/socrates/internal/refusal"
	"github.com/spf13/cobra"
)

const skillOperation = "skill"

func newSkillCommand() *cobra.Command {
	return &cobra.Command{
		Use:   skillOperation,
		Short: skillShort,
		Args:  cobra.ArbitraryArgs,
		RunE: func(command *cobra.Command, arguments []string) error {
			if len(arguments) > 0 {
				return refusal.New(skillTakesNoArguments, arguments[0])
			}

			for _, installing := range skillOperationsUnder(command.Root(), command) {
				if err := installing.RunE(installing, nil); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func skillOperationsUnder(root, topLevel *cobra.Command) []*cobra.Command {
	var operations []*cobra.Command
	for _, namespace := range root.Commands() {
		if namespace == topLevel {
			continue
		}
		for _, operation := range namespace.Commands() {
			if operation.Name() == skillOperation && operation.RunE != nil {
				operations = append(operations, operation)
			}
		}
	}
	return operations
}
