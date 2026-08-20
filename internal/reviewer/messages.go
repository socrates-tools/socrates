package reviewer

const (
	aReviewInsideAReview    = "this process is already the reviewer for run %s, and a review inside a review spawns without end: run scrat review from your own session instead"
	checkingOutTheSnapshot  = "checking out the snapshot: %w"
	clearingTheLastRound    = "clearing what the last round left where this round's worktree goes: %w"
	storedAndUnreadable     = "the answer was appended to %s and the document could not be read back afterwards (%w), so the answer is recorded and what is still unanswered is not: read the last lines of that file and remove any partial one"
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
	theCapWasReached     = ". The round cap of %d was reached, so what is left is for the person rather than for another round"
	cancelledReason      = "the person stopped it"
	noRoundToReview      = "the run has no round to review, so nothing was given to a reviewer"

	notAwaitingResponses  = "run %s is %s rather than awaiting responses, so there is nothing to answer: read what it found with scrat findings"
	noRoundToRespondTo    = "run %s has no round to answer: start a review with scrat review"
	noFindingsToRespondTo = "round %d of run %s reported no findings to answer: read what it concluded with scrat status"
	noSuchFinding         = "round %d of run %s reported no finding named %s: name one of the identifiers scrat findings prints"
	stillUnanswered       = "%d findings of this round have no response yet, and silence is neither a fix nor a rejection: run scrat respond for %s first"
)
