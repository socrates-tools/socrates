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

	awaitingLine   = "awaiting responses: socrates is waiting for the author, and no process of this run is running"
	unansweredLine = "unanswered: %d"
	readyLine      = "ready to continue: run scrat continue to start the next round"

	rulesHeading = "How to answer these findings:"
	ruleLine     = "  - %s"

	didNotReviewIt     = "You did not review this code. An independent reviewer did, and it did not write any of it."
	fixOrReject        = "For each finding, either apply a fix and record it with scrat respond --fixed, or record a rejection with scrat respond --rejected and a concrete reason."
	aRejectionIsAClaim = "A rejection is a claim you are making to the next reviewer, which judges it. It is not a decision, and it does not make the finding go away."
	doNotWeakenTheWork = "Do not resolve a finding by weakening or deleting behaviour you added on purpose. Where the finding is right about a symptom and wrong about the cause, say so and fix the cause."
)
