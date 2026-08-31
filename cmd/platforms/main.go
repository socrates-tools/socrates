package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/socrates-devtool/socrates/internal/platform"
)

func main() {
	if err := state(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}

func state() error {
	stated, err := json.MarshalIndent(platform.OperatingSystems(), "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(platform.ReleaseTargetsFile, append(stated, '\n'), 0o644); err != nil {
		return fmt.Errorf(writingTheTargets, platform.ReleaseTargetsFile, err)
	}
	return nil
}
