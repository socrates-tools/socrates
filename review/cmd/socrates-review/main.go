package main

import (
	"fmt"
	"os"

	"github.com/socrates-tools/socrates/internal/refusal"
	"github.com/socrates-tools/socrates/review/exit"
	"github.com/spf13/cobra"
)

func main() {
	code := exit.Clean

	if err := newRootCommand(&code).Execute(); err != nil {
		fmt.Fprintln(os.Stderr, failurePrefix+err.Error())
		os.Exit(exitFor(err))
	}
	os.Exit(code)
}

func newRootCommand(code *int) *cobra.Command {
	root := &cobra.Command{
		Use:           toolName,
		Short:         rootShort,
		Args:          cobra.ArbitraryArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(command *cobra.Command, arguments []string) error {
			if len(arguments) > 0 {
				return refusal.New(notACommand, arguments[0])
			}
			return command.Help()
		},
	}
	root.CompletionOptions.DisableDefaultCmd = true
	root.SetFlagErrorFunc(func(_ *cobra.Command, err error) error {
		return refusal.New(unusableFlag, err)
	})
	root.AddCommand(
		newBeginCommand(code),
		newStatusCommand(code),
		newFindingsCommand(code),
		newRespondCommand(code),
		newContinueCommand(code),
		newWatchCommand(code),
		newLogsCommand(code),
		newCancelCommand(),
		newSkillCommand(),
		newConfigCommand(),
		newConfigureCommand(),
		newVersionCommand(),
		newSuperviseCommand(),
	)
	return root
}

func exitFor(err error) int {
	if refusal.Is(err) {
		return exit.Usage
	}
	return exit.Failed
}
