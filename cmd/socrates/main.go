package main

import (
	"os"

	"github.com/socrates-tools/socrates/internal/dispatcher"
	"github.com/socrates-tools/socrates/internal/exit"
)

func main() {
	code := exit.NothingAsked

	if err := dispatcher.NewRootCommand(&code).Execute(); err != nil {
		os.Exit(dispatcher.ReportFailure(err))
	}
	os.Exit(code)
}
