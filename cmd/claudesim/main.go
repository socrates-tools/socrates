package main

import (
	"os"

	"github.com/socrates-tools/socrates-review/internal/claudesim"
)

func main() {
	os.Exit(claudesim.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr, os.Environ()))
}
