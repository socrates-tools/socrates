package corpus

const (
	readingTheCorpus      = "reading the recorded answers under %s: %w"
	readingAnAnswer       = "reading the recorded answer %s: %w"
	writingAnAnswer       = "writing the recorded answer %s: %w"
	noSuchAnswer          = "no recorded answer is named %s, and the corpus holds: %s"
	answerAlreadyRecorded = "a recorded answer is already named %s, at %s, and a recorded answer is never rewritten"
	entryStatesNothingFor = "the recorded answer %s states nothing for: %s, and an answer without them is not evidence of anything"
	nameIsAPath           = "the recorded answer is named %q, which is a path rather than a name"
	notSocrates           = "reading the module the repository %s declares: %w"
	anotherRepository     = "the repository %s is %s rather than socrates, and a findings document quotes the code it was written about"
)
