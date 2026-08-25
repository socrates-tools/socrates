package document

const (
	readingTheDocument             = "reading the document: %w"
	moreThanOneDocument            = "the document is followed by more content than one document"
	theDocumentStatesNothingFor    = "the document states no %q, so the page would have nothing to head it with"
	noDecisions                    = "the document holds no decision, so there would be nothing for a person to answer"
	decisionAt                     = "decision %d"
	decisionNamed                  = "decision %q"
	noIdentifier                   = "%s carries no identifier, so an answer could not be attached to it"
	repeatedIdentifier             = "two decisions carry the identifier %q, so an answer attached to it would be ambiguous"
	statesNothingFor               = "%s states no %q"
	unknownMode                    = "%s states the mode %q, which is not one of %v"
	noOptions                      = "%s offers no option, and the field a person writes in is the surface's rather than the document's"
	optionStatesNothingFor         = "option %d of %s states no %q"
	optionCarriesNoCommand         = "option %q of %s carries no command, so what the page shows and what the terminal runs could not be the same string"
	namesNothingInTheSource        = "%s names nothing in the source it cites"
	sourceEntryCarriesNoIdentifier = "source entry %d carries no identifier, so nothing could cite it"
	repeatedSourceIdentifier       = "two source entries carry the identifier %q, so a decision citing it would be ambiguous"
)
