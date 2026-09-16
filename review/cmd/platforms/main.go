package main

import (
	"fmt"
	"os"

	"github.com/socrates-tools/socrates/internal/platform"
)

func main() {
	if err := state(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}

func state() error {
	rendered, err := platform.Rendered()
	if err != nil {
		return err
	}
	if err := os.WriteFile(platform.ReleaseTargetsFile, rendered, 0o644); err != nil {
		return fmt.Errorf(writingTheTargets, platform.ReleaseTargetsFile, err)
	}
	return nil
}
