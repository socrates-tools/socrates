package main

const (
	failurePrefix = "scrat: "

	rootShort      = "An independent reviewer for unreviewed local changes"
	reviewShort    = "Freeze this working tree and have somebody who did not write it review it"
	statusShort    = "Say what a review is doing, in as few lines as possible"
	findingsShort  = "Print what a review found, as socrates classified it"
	respondShort   = "Record that one finding was fixed, or that it is wrong and why"
	continueShort  = "Snapshot the patched working tree and have a fresh reviewer judge the answers"
	cancelShort    = "Stop a review and say that nothing was decided"
	superviseShort = "Conduct a review that scrat review started"

	intentUsage   = "what the person asked for, in their words"
	waitUsage     = "block until the review is terminal and exit with its verdict"
	jsonUsage     = "print one JSON object instead of lines for a person"
	findingUsage  = "the identifier of the one finding this answers"
	fixedUsage    = "the finding was fixed in the working tree"
	rejectedUsage = "the finding is wrong, and --reason says why"
	reasonUsage   = "why the rejected finding is wrong, or - to read it from standard input"
	noteUsage     = "what was done about the fixed finding, or - to read it from standard input"

	notACommand            = "%q is not a scrat command: run scrat --help to see what scrat does"
	unusableFlag           = "%s: run scrat --help to see what scrat accepts"
	reviewTakesNoArguments = "scrat review takes no arguments, and %q is one: pass what the person asked for with --intent"
	oneRunAtATime          = "%s names a second run, and a command reports on one: name one run or none at all"

	nameTheFinding            = "scrat respond answers one finding, and none was named: pass --finding with an identifier scrat findings printed"
	oneAnswerAtATime          = "scrat respond records one answer, and it is either a fix or a rejection: pass --fixed or --rejected, and not both"
	reasonBelongsToARejection = "--reason says why a finding is wrong, and this answer records a fix: pass --note to say what was done about it"
	noteBelongsToAFix         = "--note says what was done about a finding, and this answer rejects one: pass --reason to say why it is wrong"
	readingStandardInput      = "reading the reason from standard input: %w"

	recordedIt           = "finding %s is recorded as %s against round %d"
	replacedIt           = "an earlier response to %s was replaced, and both are kept"
	everythingIsAnswered = "every blocking finding of this round now has a response"
	stillToAnswer        = "%d findings still have none: %s"
	roundStarted         = "round %d is under way"

	followIt      = "follow it with: scrat status %s"
	stoppedIt     = "run %s is cancelled, and nothing was decided about the change"
	alreadyEnded  = "run %s already ended as %s, so there was nothing to stop"
	locatingHere  = "locating the working directory: %w"
	locatingScrat = "locating scrat itself: %w"
)
