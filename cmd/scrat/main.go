package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const exitUsage = 3

func main() {
	if err := newRootCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "scrat: "+err.Error())
		os.Exit(exitUsage)
	}
}

func newRootCommand() *cobra.Command {
	return &cobra.Command{
		Use:           "scrat",
		Short:         "An independent reviewer for unreviewed local changes",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
}
