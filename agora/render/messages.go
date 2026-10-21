package render

const (
	renderingThePage = "rendering the page: %w"

	partialDrawnWithNothing = "the partial %s has no data to be drawn with"
	drawingAPartial         = "drawing the partial %s: %w"

	dictionaryNeedsPairs  = "a dictionary is written in pairs of a word and a value, and this one holds %d arguments"
	dictionaryKeyNotAWord = "a dictionary names each value by a word, and this one names one by %v"

	illustrationTitle = "Illustration for %s"

	recommendedTag = "Recommended"

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
)

const (
	buttonPrimary   = "a primary button"
	buttonRaised    = "a primary button that is raised"
	buttonSecondary = "a secondary button"
	buttonQuiet     = "a quiet button"
	buttonLink      = "a button drawn as a link"
	buttonDisabled  = "a button that cannot be pressed"

	catalogButtonLabel = "A button"

	fieldEmpty            = "an empty field"
	fieldWithAPlaceholder = "a field with a placeholder"
	fieldDisabled         = "a field that cannot be typed into"

	catalogFieldCaption     = "A field"
	catalogFieldPlaceholder = "A placeholder"

	choiceRadio          = "a radio, one of a group, not chosen"
	choiceRadioChosen    = "a radio, one of a group, chosen"
	choiceCheckbox       = "a checkbox"
	choiceWithEverything = "a choice with a meaning and a tag"
	choiceDashed         = "a choice drawn dashed"
	choiceDisabled       = "a choice that cannot be chosen"

	catalogChoiceLabel   = "A choice"
	catalogChoiceMeaning = "What choosing it means."
	catalogChoiceTag     = "A tag"

	copyResting         = "a copy control holding text"
	copyOfferingNothing = "a copy control holding no text"

	catalogCopyText = "Copy this exactly\n\tas it was written, with a trailing space "

	expandResting  = "an expand control at rest"
	expandExpanded = "an expand control that is expanded"

	expandLabel = "Expand"
	expandHint  = "Lays what it expands over the whole window. Escape closes it."
)
