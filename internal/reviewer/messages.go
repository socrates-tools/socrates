package reviewer

const (
	aReviewInsideAReview    = "this process is already the reviewer for run %s, and a review inside a review spawns without end: run scrat review from your own session instead"
	checkingOutTheSnapshot  = "checking out the snapshot: %w"
	openingTheSupervisorLog = "opening the supervisor's own log: %w"
	startingTheReviewer     = "starting the reviewer: %w"

	alreadyEnded         = "run %s already ended as %s"
	openingTheStream     = "opening %s: %w"
	closingTheStream     = "closing %s: %w"
	noChangeList         = "socrates could not list what the change touches, so no reviewer was told what to review: %w"
	noDiff               = "socrates could not read the diff of the change, so no reviewer was told what to review: %w"
	couldNotCompare      = "socrates could not check whether the reviewer wrote to the worktree: %s"
	theReviewerWrote     = "the reviewer wrote to the snapshot worktree, which it may not do: %s. The worktree is kept at %s"
	andTheReviewFailed   = ". The review also failed: %s"
	whatItReturnedIsKept = ". What the reviewer returned is kept at %s"
	schemaNotSatisfied   = "the findings document does not satisfy the schema: %s. It is kept at %s"
	nothingToReport      = "a reviewer read the change and found nothing to report"
	whatAReviewerFound   = "a reviewer reported %d findings, %d of them for the person, and rated the change %s risk"
	cancelledReason      = "the person stopped it"
)
