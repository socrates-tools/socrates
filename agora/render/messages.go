package render

import "strings"

const (
	renderingThePage = "rendering the page: %w"

	partialDrawnWithNothing = "the partial %s has no data to be drawn with"
	drawingAPartial         = "drawing the partial %s: %w"

	dictionaryNeedsPairs  = "a dictionary is written in pairs of a word and a value, and this one holds %d arguments"
	dictionaryKeyNotAWord = "a dictionary names each value by a word, and this one names one by %v"

	layoutNeedsPairs   = "a layout is written in pairs of a field and the utilities it takes, and this one holds %d arguments"
	layoutNamesNoClass = "a layout puts utilities on a class its view holds, and %q is not one of %T"

	illustrationTitle = "Illustration for %s"

	recommendedTag = "Recommended"

	answering          = "answering"
	choiceOneLegend    = "Choose one"
	choiceManyLegend   = "Choose as many as apply"
	noneOfThese        = "None of these"
	noneOfTheseMeaning = "Write your own answer below, or leave this to whoever asked."

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
	whatWasSeen = "What was seen"
	reply       = "Reply"
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

	brandAMarkAndAName = "a mark and a name"

	catalogBrandName = "A brand"
	catalogBrandMark = "/assets/socrates-mark.svg"

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

	shellNarrow         = "narrower than 960px, stacked and scrolling with the window"
	shellAtNineSixty    = "960px wide, with the navigation and the content scrolling on their own"
	shellAtTwelveEighty = "1280px wide, with the panel's place at the right"

	catalogShellTitle      = "A shell"
	catalogShellNavigation = "Sections"
	catalogShellPanel      = "Panel"
	catalogShellSections   = 12
	catalogChoiceCount     = 3
	catalogShellHeight     = 480
	catalogNarrowWidth     = 480
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
	choiceWithEverything = "a choice with a meaning and a tag"
	choiceDashed         = "a choice drawn dashed"
	choiceDashedChosen   = "a choice drawn dashed, chosen"
	choiceDisabled       = "a choice that cannot be chosen"

	catalogChoiceLabel   = "A choice"
	catalogChoiceMeaning = "What choosing it means."
	catalogChoiceTag     = "A tag"

	catalogChoiceGroupLabel = "Some choices"

	choiceGroupSingle              = "a group of choices, one at most chosen, none yet"
	choiceGroupSingleChosen        = "a group of choices, one at most chosen, the second chosen"
	choiceGroupWithNoneOfThese     = "a group of choices ending in a choice of none of them"
	choiceGroupWithARecommendation = "a group of choices whose first is recommended and arrives chosen"
	choiceGroupMultiple            = "a group of choices, several chosen at once"
	choiceGroupDisabled            = "a group of choices that cannot be chosen"

	copyResting         = "a copy control holding text"
	copyOfferingNothing = "a copy control holding no text"

	catalogCopyText = "Copy this exactly\n\tas it was written, with a trailing space "

	expandResting  = "an expand control at rest"
	expandExpanded = "an expand control that is expanded"

	expandLabel = "Expand"
	expandHint  = "Lays what it expands over the whole window. Escape closes it."
)

const (
	codeShortLine                = "a short line"
	codeWrappedLine              = "lines wider than its block, unindented and indented"
	codeTabIndentedLine          = "lines indented with a tab"
	codeUnknownLanguage          = "a language the page does not highlight"
	diffPatch                    = "a patch"
	diffPatchInAnUnknownLanguage = "a patch in a language the page does not highlight"

	catalogLanguage        = "go"
	catalogUnknownLanguage = "klingon"
	catalogShortCode       = "return limit\n"
	catalogTabIndentedCode = "func Write() {\n\tmu.Lock()\n\treturn\n}\n"
	catalogUnknownCode     = "qapla'\nbatlh\n"
	catalogPatch           = "--- a/limit.go\n+++ b/limit.go\n@@ -1,3 +1,4 @@ func Limit\n func Limit() int {\n-\treturn 2\n+\t// the limit, raised\n+\treturn 3\n }\n"
)

var (
	catalogLongLine    = `fmt.Println("` + strings.Repeat("a line wider than its block ", 14) + `")`
	catalogWrappedCode = catalogLongLine + "\n\t" + catalogLongLine + "\n      " + catalogLongLine + "\n"
)
