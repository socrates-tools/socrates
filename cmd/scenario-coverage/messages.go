package main

const (
	failurePrefix = "scenario-coverage: "

	gapHeading       = "%s"
	gapLine          = "  %d: %s"
	gapSummary       = "%d scenarios under way, %d of them named by no test"
	coveredSummary   = "%d scenarios under way, every one of them named by a test"
	notStartedLine   = "%d more scenarios belong to changes nothing has started, and they are not counted"
	noSpecifications = "there are no specifications under %s, so there is nothing to check: run this from the repository root"
)
