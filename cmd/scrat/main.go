package main

import (
	"fmt"
	"os"

	"github.com/rui-ktei/socrates/internal/refusal"
	"github.com/spf13/cobra"
)

const (
	exitClean   = 0
	exitFailed  = 2
	exitUsage   = 3
	exitWorking = 4
)

func main() {
	if err := newRootCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "scrat: "+err.Error())
		os.Exit(exitFor(err))
	}
}

func newRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:           "scrat",
		Short:         "An independent reviewer for unreviewed local changes",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(newReviewCommand(), newSuperviseCommand())
	return root
}

func exitFor(err error) int {
	if refusal.Is(err) {
		return exitUsage
	}
	return exitFailed
}
