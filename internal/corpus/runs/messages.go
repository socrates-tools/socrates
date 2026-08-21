package runs

const (
	noSuchRound          = "the run reviewed nothing in round %d, and %s holds the rounds it did review"
	roundAnsweredNothing = "round %d of %s answered no findings document: %w"
	readingTheLog        = "reading the event log %s: %w"
	logStatesNothingFor  = "the event log %s states nothing for: %s, so the answer cannot be recorded with the model that produced it"
)
