package command

import (
	"fmt"
	"os"

	"github.com/socrates-tools/socrates/agora/surface"
	"github.com/socrates-tools/socrates/internal/exit"
	"github.com/socrates-tools/socrates/internal/refusal"
	"github.com/spf13/cobra"
)

func NewRootCommand(code *int) *cobra.Command {
	root := &cobra.Command{
		Use:           toolName,
		Short:         rootShort,
		Args:          cobra.ArbitraryArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			return surface.CheckNoHomeOfItsOwn()
		},
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
		newServeCommand(),
		newPollCommand(code),
		newEndCommand(),
		newStopCommand(),
		newVersionCommand(),
	)
	return root
}

func ReportFailure(err error) int {
	fmt.Fprintln(os.Stderr, failurePrefix+err.Error())
	return exitFor(err)
}

func exitFor(err error) int {
	if refusal.Is(err) {
		return exit.Unusable
	}
	return exit.Failed
}
