package main

const (
	failurePrefix = "scrat: "

	rootShort      = "An independent reviewer for unreviewed local changes"
	reviewShort    = "Freeze this working tree and have somebody who did not write it review it"
	statusShort    = "Say what a review is doing, in as few lines as possible"
	findingsShort  = "Print what a review found, as socrates classified it"
	cancelShort    = "Stop a review and say that nothing was decided"
	superviseShort = "Conduct a review that scrat review started"

	intentUsage = "what the person asked for, in their words"
	waitUsage   = "block until the review is terminal and exit with its verdict"
	jsonUsage   = "print one JSON object instead of lines for a person"

	notACommand            = "%q is not a scrat command: run scrat --help to see what scrat does"
	unusableFlag           = "%s: run scrat --help to see what scrat accepts"
	reviewTakesNoArguments = "scrat review takes no arguments, and %q is one: pass what the person asked for with --intent"
	oneRunAtATime          = "%s names a second run, and a command reports on one: name one run or none at all"

	followIt      = "follow it with: scrat status %s"
	stoppedIt     = "run %s is cancelled, and nothing was decided about the change"
	alreadyEnded  = "run %s already ended as %s, so there was nothing to stop"
	locatingHere  = "locating the working directory: %w"
	locatingScrat = "locating scrat itself: %w"
)
