package follow

const (
	watchingTheLog   = "reading the event log %s: %w"
	unreadableLine   = "a line of the event log %s is not an event: %w"
	nothingEverWrote = "the run %s wrote no event log within %s and no reviewer of it is alive, so nothing has reviewed the change: start a review with socrates review begin"
)
