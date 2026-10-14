package dispatcher

import (
	"fmt"

	"github.com/socrates-tools/socrates"
	"github.com/socrates-tools/socrates/internal/refusal"
	"github.com/spf13/cobra"
)

func newLicencesCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "licences",
		Aliases: []string{"licenses"},
		Short:   licencesShort,
		Args:    cobra.ArbitraryArgs,
		RunE: func(command *cobra.Command, arguments []string) error {
			if len(arguments) > 0 {
				return refusal.New(licencesTakesNoArguments, arguments[0])
			}

			fmt.Fprint(command.OutOrStdout(), socrates.Licence+socrates.Notice)
			return nil
		},
	}
}
