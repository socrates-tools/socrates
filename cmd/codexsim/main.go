package main

import (
	"os"

	"github.com/socrates-tools/socrates-review/internal/codexsim"
)

func main() {
	os.Exit(codexsim.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr, os.Environ()))
}
