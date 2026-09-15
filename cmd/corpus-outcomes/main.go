package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/socrates-tools/socrates/internal/corpus"
	"github.com/socrates-tools/socrates/internal/corpus/outcomes"
)

func main() {
	if err := report(); err != nil {
		fmt.Fprintln(os.Stderr, failurePrefix+err.Error())
		os.Exit(2)
	}
}

func report() error {
	rewrite := flag.Bool("write", false, writeFlag)
	flag.Parse()

	entries, err := corpus.Entries()
	if err != nil {
		return err
	}

	reached, err := outcomes.Of(entries)
	if err != nil {
		return err
	}

	root, err := os.Getwd()
	if err != nil {
		return err
	}

	if *rewrite {
		path := filepath.Join(root, outcomes.DistributionPath)
		if err := os.WriteFile(path, []byte(reached.Distribution()), 0o644); err != nil {
			return fmt.Errorf(writingTheDistribution, path, err)
		}
		fmt.Printf(wroteTheDistribution+"\n", path)
		return nil
	}

	naming, err := outcomes.Naming(root)
	if err != nil {
		return err
	}

	fmt.Print(reached.Render(naming))
	return nil
}
