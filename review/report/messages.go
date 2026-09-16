package report

const (
	runLine               = "run: %s"
	stateLine             = "state: %s"
	roundLine             = "round: %d"
	baseLine              = "base: %s"
	reviewerLine          = "reviewer: %s (%s)"
	verdictLine           = "verdict: %s"
	countsLine            = "findings: %d in total, %d for the person, %d to patch, %d informational"
	findingLine           = "%s  %s  %s:%d  %s"
	descriptionLine       = "    %s"
	riskLine              = "risk: %s"
	rationaleLine         = "why: %s"
	nothingToReportLine   = "no findings: a reviewer read the change and reported nothing"
	noReviewLine          = "no review completed: %s"
	stillWorkingReason    = "the run is %s, so wait for it to end and read the findings then"
	personAudience        = "needs a person"
	patchableAudience     = "can be patched"
	informationalAudience = "informational"
	encodingFailed        = "encoding what socrates review concluded: %w"

	awaitingLine        = "awaiting responses: socrates review is waiting for the author, or for a person the author asks, and no process of this run is running"
	unansweredLine      = "unanswered: %d"
	awaitingAPersonLine = "of those, needing a person rather than a patch: %d"
	readyLine           = "ready to continue: run socrates review continue to start the next round"

	awaitsAPersonLine = "awaits a person: this round holds what only a person can act on, so wait for what they answer"

	rulesHeading = "How to answer these findings:"
	ruleLine     = "  - %s"

	settingLine = "%-*s  %-*s  %s"
)
