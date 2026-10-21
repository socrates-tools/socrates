package render

const (
	renderingThePage = "rendering the page: %w"

	partialDrawnWithNothing = "the partial %s has no data to be drawn with"
	drawingAPartial         = "drawing the partial %s: %w"

	illustrationTitle = "Illustration for %s"

	copyLabel  = "Copy code"
	copyHint   = "Puts this code on your clipboard, exactly as it was written. Nothing is sent."
	copiedWord = "Copied"
	notCopied  = "Could not copy"

	severityErrorName   = "Error"
	severityError       = "Severity: Error"
	severityErrorHint   = "Likely wrong; the change should not land as it is."
	severityWarningName = "Warning"
	severityWarning     = "Severity: Warning"
	severityWarningHint = "Worth a look; may be fine."
	severityInfoName    = "Info"
	severityInfo        = "Severity: Info"
	severityInfoHint    = "No action expected."
	severityNamed       = "Severity: %s"
	audienceNamed       = "Decides: %s"
	spokenWithItsHint   = "%s. %s"
)

const (
	catalogTitle = "Catalog"

	dialogWithEverything    = "a heading, a body and actions"
	dialogWithAHeadingAlone = "a heading alone"
	catalogDialogHeading    = "A dialog"

	badgeNeutral   = "a neutral tone with a label"
	badgeAccent    = "an accent tone with a label"
	badgeWarning   = "a warning tone with a label"
	badgeError     = "an error tone with a label"
	badgeIconAlone = "an icon alone"
	badgeExplained = "a label that says what it means on hover, on focus and to a screen reader"

	catalogBadgeLabel = "A badge"
	catalogBadgeTitle = "What this badge is"
	catalogBadgeHint  = "A sentence saying what it means."

	cardWithSections   = "a body and a footer"
	cardWithABodyAlone = "a body alone"
	cardWithNoSections = "no sections"

	overlineWordsAlone = "words alone"
	overlineWithACount = "words and a count"

	catalogOverlineWords = "An overline"

	progressNothingDone = "nothing done"
	progressPartlyDone  = "partly done"
	progressAllDone     = "all done"
	progressNamed       = "named to a screen reader"

	catalogProgressLabel = "Progress"

	dotInTheColourAround = "the colour around it"
	dotNeutral           = "a neutral tone"
	dotAccent            = "an accent tone"
	dotResolved          = "a resolved tone"
	dotWarning           = "a warning tone"
	dotError             = "an error tone"

	tooltipWithATitleAndAHint = "a title and a hint"
	tooltipWithATitleAlone    = "a title alone"

	catalogAnchorLabel = "Hover or focus me"
)
