package corpus

const (
	readingTheCorpus            = "reading the recorded answers under %s: %w"
	readingAnAnswer             = "reading the recorded answer %s: %w"
	writingAnAnswer             = "writing the recorded answer %s: %w"
	noSuchAnswer                = "no recorded answer is named %s, and the corpus holds: %s"
	answerAlreadyRecorded       = "a recorded answer is already named %s, at %s, and a recorded answer is never rewritten"
	entryStatesNothingFor       = "the recorded answer %s states nothing for: %s, and an answer without them is not evidence of anything"
	nameIsAPath                 = "the recorded answer is named %q, which is a path rather than a name"
	nameTheCorpusWouldNotSee    = "the recorded answer is named %q, and the corpus never reads a name opening with a dot or an underscore, so give it a name that starts with a letter"
	documentSocratesCannotRoute = "the recorded answer %s holds a document socrates never routed, so it is evidence of nothing: %w"
	noModuleDeclared            = "the repository %s declares no Go module, so it is not socrates, and a findings document quotes the code it was written about: run corpus-record for a run that reviewed socrates itself"
	unreadableModule            = "reading %s to see which repository an answer is about: %w"
	anotherRepository           = "the repository %s is %s rather than socrates, and a findings document quotes the code it was written about"
)
