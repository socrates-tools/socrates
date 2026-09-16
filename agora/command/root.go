package command

import (
	"github.com/socrates-tools/socrates/agora/surface"
	"github.com/socrates-tools/socrates/internal/refusal"
	"github.com/spf13/cobra"
)

func NewRootCommand(code *int, loadSurfaceSettings SurfaceLoader) *cobra.Command {
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
		newServeCommand(loadSurfaceSettings),
		newPollCommand(code),
		newEndCommand(),
		newStopCommand(),
	)
	return root
}
