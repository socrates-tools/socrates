package main

import (
	"fmt"
	"os"

	"github.com/socrates-tools/socrates/agora/surface"
)

func main() {
	if err := state(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}

func state() error {
	generated, err := surface.ClientTypes()
	if err != nil {
		return err
	}
	if err := os.WriteFile(surface.ClientTypesFile, generated, 0o644); err != nil {
		return fmt.Errorf(writingTheTypes, surface.ClientTypesFile, err)
	}
	return nil
}
