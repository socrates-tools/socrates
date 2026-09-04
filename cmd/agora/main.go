package main

import (
	"fmt"
	"os"

	"github.com/socrates-tools/socrates-review/internal/refusal"
	"github.com/socrates-tools/socrates-review/internal/surface"
	"github.com/spf13/cobra"
)

const (
	answered   = 0
	unanswered = 1
	failed     = 2
	usage      = 3
)

func main() {
	code := answered

	if err := newRootCommand(&code).Execute(); err != nil {
		fmt.Fprintln(os.Stderr, failurePrefix+err.Error())
		os.Exit(exitFor(err))
	}
	os.Exit(code)
}

func newRootCommand(code *int) *cobra.Command {
	root := &cobra.Command{
		Use:           "agora",
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

func exitFor(err error) int {
	if refusal.Is(err) {
		return usage
	}
	return failed
}
