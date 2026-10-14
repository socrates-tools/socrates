package dispatcher

import (
	"fmt"
	"os"

	agoracommand "github.com/socrates-tools/socrates/agora/command"
	designcommand "github.com/socrates-tools/socrates/design/command"
	implementcommand "github.com/socrates-tools/socrates/implement/command"
	reviewcommand "github.com/socrates-tools/socrates/review/command"

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
		reviewcommand.NewRootCommand(code),
		agoracommand.NewRootCommand(code),
		designcommand.NewRootCommand(),
		implementcommand.NewRootCommand(),
		newSkillCommand(),
		newLicencesCommand(),
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
