package main

import (
	"fmt"
	"os"

	"github.com/ktei/socrates/internal/promotion"
)

func main() {
	dropped, err := report()
	if err != nil {
		fmt.Fprintln(os.Stderr, failurePrefix+err.Error())
		os.Exit(2)
	}
	if dropped {
		os.Exit(1)
	}
}

func report() (bool, error) {
	root, err := os.Getwd()
	if err != nil {
		return false, err
	}

	checked, err := promotion.Check(root)
	if err != nil {
		return false, err
	}

	out := os.Stdout
	spec := ""
	for _, drop := range checked.Drops {
		if drop.Spec != spec {
			spec = drop.Spec
			fmt.Fprintf(out, dropHeading+"\n", spec)
		}
		fmt.Fprintf(out, dropLine+"\n", drop.Requirement, drop.Condition)
		named := scenarioLine
		if drop.Condition == promotion.PromotedAhead {
			named = promotedLine
		}
		for _, scenario := range drop.Scenarios {
			fmt.Fprintf(out, named+"\n", scenario)
		}
	}

	if checked.Changes == 0 {
		fmt.Fprintln(out, nothingUnderWay)
		return false, nil
	}
	if len(checked.Drops) == 0 {
		fmt.Fprintf(out, nothingDropped+"\n", checked.Substituted, checked.Added)
		return false, nil
	}
	fmt.Fprintf(out, somethingDropped+"\n", len(checked.Drops), checked.Substituted, checked.Added)
	return true, nil
}
