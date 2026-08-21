package outcomes

const (
	whatSocratesWouldDo  = "what socrates would do with each answer a real reviewer gave"
	reachedBy            = "%s: %d recorded"
	reachedByAndAsserted = "%s: %d recorded, %d asserted by the suite"
	answerLine           = "  %s"
	onlyReachedByHand    = "  no answer a real reviewer gave reaches %s, so every scenario that asserts it was written by hand"
	recordedAnswers      = "recorded answers: %d"

	answerDoesNotParse = "the recorded answer %s does not parse as a findings document: %w"
	answerNamesNoRound = "the recorded answer %s names no round, and what a round leads to depends on which round it was"

	readingATest = "reading %s: %w"
)
