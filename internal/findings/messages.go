package findings

const (
	readingTheDocument      = "reading the findings document: %w"
	statesNothingFor        = "the findings document states no %q, so it does not say what the review found"
	moreThanOneDocument     = "the findings document is followed by more content than one document"
	unknownRiskLevel        = "the findings document states the risk level %q, which is not one of %v"
	findingAt               = "finding %d"
	findingNamed            = "finding %q"
	noIdentifier            = "%s carries no identifier, so a response could not be attached to it"
	repeatedIdentifier      = "two findings carry the identifier %q, so a response attached to it would be ambiguous"
	unknownSeverity         = "%s states the severity %q, which is not one of %v"
	namesTheFile            = "%s names the file %q: %w"
	linesAreCountedFromOne  = "%s names line %d, and lines are counted from one"
	describesNothing        = "%s describes nothing"
	aFindingNamesItsFile    = "a finding names the file it is about"
	relativeToTheRoot       = "a file is named relative to the repository root"
	asTheRepositorySpellsIt = "a file is named as the repository spells it"
	insideTheRepository     = "a file inside the repository is what a review is about"
)
