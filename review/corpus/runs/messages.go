package runs

const (
	noSuchRound                      = "the run reviewed nothing in round %d, and %s holds the rounds it did review"
	roundAnsweredNothing             = "round %d of %s answered no findings document: %w"
	roundAnsweredSomethingUnreadable = "round %d of %s answered something socrates-review never routed, so it is not evidence of what a reviewer reports: %w"
	readingTheLog                    = "reading the event log %s: %w"
	logStatesNoRound                 = "the event log %s states a round it does not number, so which model answered which round cannot be told"
	logStatesNothingFor              = "the event log holds nothing for round %d in %s: %s, so the answer cannot be recorded with the model that produced it"
)
