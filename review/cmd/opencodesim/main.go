package main

import (
	"os"

	"github.com/socrates-tools/socrates/review/opencodesim"
)

func main() {
	os.Exit(opencodesim.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr, os.Environ()))
}
