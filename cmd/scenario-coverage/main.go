package main

import (
	"fmt"
	"os"

	"github.com/socrates-devtool/socrates/internal/coverage"
)

func main() {
	gaps, err := report()
	if err != nil {
		fmt.Fprintln(os.Stderr, failurePrefix+err.Error())
		os.Exit(2)
	}
	if gaps {
		os.Exit(1)
	}
}

func report() (bool, error) {
	root, err := os.Getwd()
	if err != nil {
		return false, err
	}

	scenarios, err := coverage.Scenarios(root)
	if err != nil {
		return false, err
	}
	if len(scenarios) == 0 {
		return false, fmt.Errorf(noSpecifications, root)
	}

	claims, err := coverage.Claims(root)
	if err != nil {
		return false, err
	}

	underway := 0
	for _, scenario := range scenarios {
		if scenario.Underway {
			underway++
		}
	}

	var gaps, waiting []coverage.Scenario
	for _, gap := range coverage.Gaps(scenarios, claims) {
		if gap.Underway {
			gaps = append(gaps, gap)
			continue
		}
		waiting = append(waiting, gap)
	}

	out := os.Stdout
	spec := ""
	for _, gap := range gaps {
		if gap.Spec != spec {
			spec = gap.Spec
			fmt.Fprintf(out, gapHeading+"\n", spec)
		}
		fmt.Fprintf(out, gapLine+"\n", gap.Line, gap.Text)
	}

	if len(gaps) == 0 {
		fmt.Fprintf(out, coveredSummary+"\n", underway)
	} else {
		fmt.Fprintf(out, gapSummary+"\n", underway, len(gaps))
	}
	if len(waiting) > 0 {
		fmt.Fprintf(out, notStartedLine+"\n", len(waiting))
	}

	return len(gaps) > 0, nil
}
