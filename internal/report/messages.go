package report

const (
	runLine               = "run: %s"
	stateLine             = "state: %s"
	roundLine             = "round: %d"
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
	encodingFailed        = "encoding what socrates concluded: %w"
)
