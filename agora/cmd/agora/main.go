package main

import (
	"os"

	"github.com/socrates-tools/socrates/agora/command"
	"github.com/socrates-tools/socrates/internal/exit"
)

func main() {
	code := exit.NothingAsked

	if err := command.NewRootCommand(&code).Execute(); err != nil {
		os.Exit(command.ReportFailure(err))
	}
	os.Exit(code)
}
