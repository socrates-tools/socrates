package outcomes

const (
	whatSocratesWouldDo = "what socrates-review would do with each answer a real reviewer gave"
	reachedBy           = "%s: %d recorded"
	reachedByAndNamed   = "%s: %d recorded, named by %d tests"
	answerLine          = "  %s"
	onlyReachedByHand   = "  no answer a real reviewer gave reaches %s, so every test that names it was written against input somebody wrote by hand"
	recordedAnswers     = "recorded answers: %d"

	answerDoesNotParse = "the recorded answer %s does not parse as a findings document: %w"
	answerNamesNoRound = "the recorded answer %s names no round, and what a round leads to depends on which round it was"

	readingATest = "reading %s: %w"
)
