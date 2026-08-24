package main

const (
	failurePrefix = "socrates: "

	rootShort      = "An independent reviewer for unreviewed local changes"
	reviewShort    = "Freeze this working tree and have somebody who did not write it review it"
	statusShort    = "Say what a review is doing, in as few lines as possible"
	findingsShort  = "Print what a review found, as socrates classified it"
	respondShort   = "Record that one finding was fixed, that it is wrong and why, or what a person decided about it"
	continueShort  = "Snapshot the patched working tree and have a fresh reviewer judge the answers"
	watchShort     = "Follow a review as it happens and exit with its verdict"
	logsShort      = "Print what a review recorded, for a person to read afterwards"
	cancelShort    = "Stop a review and say that nothing was decided"
	skillShort     = "Install the skill that teaches a coding session to drive a review"
	versionShort   = "Say which build this is, by the commit it was built from"
	superviseShort = "Conduct a review that socrates review started"

	intentUsage   = "what the person asked for, in their words"
	waitUsage     = "block until the review concludes or parks on you, and exit with its status"
	jsonUsage     = "print one JSON object instead of lines for a person"
	findingUsage  = "the identifier of the one finding this answers"
	fixedUsage    = "the finding was fixed in the working tree"
	rejectedUsage = "the finding is wrong, and --reason says why"
	decidedUsage  = "a person answered the question the finding asks, and --reason says what they decided"
	reasonUsage   = "why the rejected finding is wrong or what was decided about it, or - to read it from standard input"
	noteUsage     = "what was done about the fixed finding, or - to read it from standard input"
	toUsage       = "the skills directory to write to, instead of the one under your home directory"

	notACommand             = "%q is not a socrates command: run socrates --help to see what socrates does"
	unusableFlag            = "%s: run socrates --help to see what socrates accepts"
	reviewTakesNoArguments  = "socrates review takes no arguments, and %q is one: pass what the person asked for with --intent"
	versionTakesNoArguments = "socrates version takes no arguments, and %q is one: run socrates version on its own to see which build this is"
	oneRunAtATime           = "%s names a second run, and a command reports on one: name one run or none at all"

	nameTheFinding        = "socrates respond answers one finding, and none was named: pass --finding with an identifier socrates findings printed"
	oneAnswerAtATime      = "socrates respond records one answer, and it is a fix, a rejection or a decision: pass exactly one of --%s, --%s or --%s"
	reasonBelongsToAClaim = "--reason says why a finding is wrong or what a person decided about it, and this answer records a fix: pass --note to say what was done about it"
	noteBelongsToAFix     = "--note says what was done about a fixed finding, and this answer records that it is %s: pass --reason instead, because that is the claim the next reviewer judges"
	readingStandardInput  = "reading the reason from standard input: %w"

	recordedIt           = "finding %s is recorded as %s against round %d"
	replacedIt           = "an earlier response to %s was replaced, and both are kept"
	everythingIsAnswered = "every finding of this round that a response can address now has one"
	stillToAnswer        = "%d findings still have none: %s"
	roundStarted         = "round %d is under way"

	installedTheSkill = "the socrates skill is at %s, and a session that reads it can drive a review itself"

	stillStarting = "run %s is %s, and its reviewer has written nothing yet"

	followIt         = "follow it with: socrates status %s"
	stoppedIt        = "run %s is cancelled, and nothing was decided about the change"
	alreadyEnded     = "run %s already ended as %s, so there was nothing to stop"
	locatingHere     = "locating the working directory: %w"
	locatingSocrates = "locating socrates itself: %w"
)
