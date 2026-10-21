package render

import (
	"bytes"
	_ "embed"
	"html/template"
	"slices"
	"strings"
)

//go:embed catalog.gohtml
var catalogSource string

type dialogView struct {
	ID          string
	Class       string
	Interaction string
	Heading     string
	Body        string
	Actions     string
}

type buttonView struct {
	Class       string
	Variant     string
	Type        string
	Interaction string
	Ends        string
	Label       string
	Content     string
	Raised      bool
	Disabled    bool
	Hidden      bool
	Autofocus   bool
}

type fieldView struct {
	Class        string
	CaptionClass string
	ControlClass string
	Caption      string
	Name         string
	Rows         int
	Placeholder  string
	Interaction  string
	Decision     string
	Disabled     bool
}

type choiceView struct {
	Class        string
	ControlClass string
	Family       string
	Variant      string
	Type         string
	Name         string
	Value        string
	Interaction  string
	Decision     string
	Label        string
	Meaning      string
	Tag          string
	Recommended  bool
	Checked      bool
	Disabled     bool
}

type copyView struct {
	Label       string
	Hint        string
	Copied      string
	Failed      string
	Text        string
	Interaction string
}

type expandView struct {
	Label       string
	Hint        string
	Interaction string
	Expanded    bool
}

type anchorView struct {
	ID      string
	Tooltip string
	Label   string
	Title   string
	Hint    string
}

type state struct {
	Key    string
	Name   string
	Opens  string
	Anchor *anchorView
	View   any
}

func badgeState(name string, view badgeView) state {
	return state{Name: name, View: view}
}

func cardState(name string, view cardView) state {
	return state{Name: name, View: view}
}

func overlineState(name string, view overlineView) state {
	return state{Name: name, View: view}
}

func progressBarState(name string, view progressBarView) state {
	return state{Name: name, View: view}
}

func statusDotState(name string, view statusDotView) state {
	return state{Name: name, View: view}
}

func tooltipState(name string, id string, anchor anchorView) state {
	anchor.ID = "anchor-" + id
	anchor.Tooltip = id
	return state{Name: name, Anchor: &anchor, View: tooltipView{ID: id, Class: "tooltip"}}
}

func markedBadge(tone string) badgeView {
	return badgeView{Tone: tone, Icon: icon("info"), Label: catalogBadgeLabel}
}

func dialogState(key, name string, view dialogView) state {
	return state{Key: key, Name: name, Opens: view.ID, View: view}
}

func stateDrawn(key, name string, view any) state {
	return state{Key: key, Name: name, View: view}
}

var catalogued = map[string][]state{
	"ui-badge": {
		badgeState(badgeNeutral, markedBadge("neutral")),
		badgeState(badgeAccent, markedBadge("accent")),
		badgeState(badgeWarning, markedBadge("warning")),
		badgeState(badgeError, markedBadge("error")),
		badgeState(badgeIconAlone, badgeView{Tone: "neutral", Icon: icon("info")}),
		badgeState(badgeExplained, badgeView{
			Tone:      "warning",
			Icon:      icon("triangle-alert"),
			Label:     catalogBadgeLabel,
			Title:     catalogBadgeTitle,
			Hint:      catalogBadgeHint,
			Spoken:    catalogBadgeTitle + ". " + catalogBadgeHint,
			Focusable: true,
		}),
	},
	"ui-card": {
		cardState(cardWithSections, cardView{Content: "catalog-card-sections"}),
		cardState(cardWithABodyAlone, cardView{Content: "catalog-card-body"}),
		cardState(cardWithNoSections, cardView{Content: "catalog-card-plain"}),
	},
	"ui-dialog": {
		dialogState("with-everything", dialogWithEverything, dialogView{
			ID:      "catalog-dialog",
			Class:   "catalog-dialog",
			Heading: catalogDialogHeading,
			Body:    "catalog-dialog-body",
			Actions: "catalog-dialog-actions",
		}),
		dialogState("heading-alone", dialogWithAHeadingAlone, dialogView{
			ID:      "catalog-dialog-bare",
			Class:   "catalog-dialog-bare",
			Heading: catalogDialogHeading,
		}),
	},
	"ui-button": {
		stateDrawn("primary", buttonPrimary, buttonView{Variant: "primary", Type: "button", Label: catalogButtonLabel}),
		stateDrawn("raised", buttonRaised, buttonView{Variant: "primary", Raised: true, Type: "button", Label: catalogButtonLabel}),
		stateDrawn("secondary", buttonSecondary, buttonView{Variant: "secondary", Type: "button", Label: catalogButtonLabel}),
		stateDrawn("quiet", buttonQuiet, buttonView{Variant: "quiet", Type: "button", Label: catalogButtonLabel}),
		stateDrawn("link", buttonLink, buttonView{Variant: "link", Type: "button", Label: catalogButtonLabel}),
		stateDrawn("disabled", buttonDisabled, buttonView{Variant: "primary", Type: "button", Label: catalogButtonLabel, Disabled: true}),
	},
	"ui-field": {
		stateDrawn("empty", fieldEmpty, fieldView{Class: "catalog-field", Caption: catalogFieldCaption, Name: "catalog-empty", Rows: 2}),
		stateDrawn("placeholder", fieldWithAPlaceholder, fieldView{Class: "catalog-field", Caption: catalogFieldCaption, Name: "catalog-placeholder", Rows: 2, Placeholder: catalogFieldPlaceholder}),
		stateDrawn("disabled", fieldDisabled, fieldView{Class: "catalog-field", Caption: catalogFieldCaption, Name: "catalog-disabled", Rows: 2, Disabled: true}),
	},
	"ui-choice": {
		stateDrawn("radio", choiceRadio, choiceView{Family: "catalog-choice", Variant: "solid", Type: "radio", Name: "catalog-radio", Value: "0", Label: catalogChoiceLabel}),
		stateDrawn("radio-chosen", choiceRadioChosen, choiceView{Family: "catalog-choice", Variant: "solid", Type: "radio", Name: "catalog-radio", Value: "1", Label: catalogChoiceLabel, Checked: true}),
		stateDrawn("checkbox", choiceCheckbox, choiceView{Family: "catalog-choice", Variant: "solid", Type: "checkbox", Name: "catalog-checkbox", Value: "0", Label: catalogChoiceLabel}),
		stateDrawn("with-everything", choiceWithEverything, choiceView{Family: "catalog-choice", Variant: "solid", Type: "radio", Name: "catalog-everything", Value: "0", Label: catalogChoiceLabel, Meaning: catalogChoiceMeaning, Tag: catalogChoiceTag, Recommended: true}),
		stateDrawn("dashed", choiceDashed, choiceView{Family: "catalog-choice", Variant: "dashed", Type: "radio", Name: "catalog-dashed", Label: catalogChoiceLabel, Meaning: catalogChoiceMeaning}),
		stateDrawn("disabled", choiceDisabled, choiceView{Family: "catalog-choice", Variant: "solid", Type: "checkbox", Name: "catalog-disabled-choice", Value: "0", Label: catalogChoiceLabel, Disabled: true}),
	},
	"ui-copy": {
		stateDrawn("resting", copyResting, copyView{Label: copyLabel, Hint: copyHint, Copied: copiedWord, Failed: notCopied, Text: catalogCopyText}),
		stateDrawn("offering-nothing", copyOfferingNothing, copyView{Label: copyLabel, Hint: copyHint, Copied: copiedWord, Failed: notCopied}),
	},
	"ui-expand": {
		stateDrawn("resting", expandResting, expandView{Label: expandLabel, Hint: expandHint}),
		stateDrawn("expanded", expandExpanded, expandView{Label: expandLabel, Hint: expandHint, Expanded: true}),
	},
	"ui-overline": {
		overlineState(overlineWordsAlone, overlineOf("", catalogOverlineWords)),
		overlineState(overlineWithACount, overlineCounting("", catalogOverlineWords, 3)),
	},
	"ui-progress-bar": {
		progressBarState(progressNothingDone, progressBarOf("", 4, 0)),
		progressBarState(progressPartlyDone, progressBarOf("", 4, 2)),
		progressBarState(progressAllDone, progressBarOf("", 4, 4)),
		progressBarState(progressNamed, progressBarView{Max: 4, Value: 2, Label: catalogProgressLabel}),
	},
	"ui-status-dot": {
		statusDotState(dotInTheColourAround, statusDotOf("", "")),
		statusDotState(dotNeutral, statusDotOf("", "neutral")),
		statusDotState(dotAccent, statusDotOf("", "accent")),
		statusDotState(dotResolved, statusDotOf("", "resolved")),
		statusDotState(dotWarning, statusDotOf("", "warning")),
		statusDotState(dotError, statusDotOf("", "error")),
	},
	"ui-tooltip": {
		tooltipState(tooltipWithATitleAndAHint, "catalog-tooltip", anchorView{Label: catalogAnchorLabel, Title: catalogBadgeTitle, Hint: catalogBadgeHint}),
		tooltipState(tooltipWithATitleAlone, "catalog-tooltip-bare", anchorView{Label: catalogAnchorLabel, Title: catalogBadgeTitle}),
	},
}

type catalogState struct {
	Key    string
	Name   string
	Opens  string
	Anchor *anchorView
	Drawn  template.HTML
}

type catalogPrimitive struct {
	Name   string
	States []catalogState
}

type catalogView struct {
	Title      string
	Policy     string
	Primitives []catalogPrimitive
}

func Catalog() (string, error) {
	names := make([]string, 0, len(catalogued))
	for name := range catalogued {
		names = append(names, name)
	}
	slices.Sort(names)

	view := catalogView{Title: catalogTitle, Policy: policy}
	for _, name := range names {
		primitive := catalogPrimitive{Name: name}
		for _, held := range catalogued[name] {
			markup, err := partialDrawn(name, held.View)
			if err != nil {
				return "", err
			}
			primitive.States = append(primitive.States, catalogState{Key: keyOf(held), Name: held.Name, Opens: held.Opens, Anchor: held.Anchor, Drawn: markup})
		}
		view.Primitives = append(view.Primitives, primitive)
	}

	var written bytes.Buffer
	if err := page.ExecuteTemplate(&written, "catalog", view); err != nil {
		return "", err
	}
	return written.String(), nil
}

var keySpelling = strings.NewReplacer(" ", "-", ",", "", ".", "", "'", "")

func keyOf(held state) string {
	if held.Key != "" {
		return held.Key
	}
	return keySpelling.Replace(strings.ToLower(held.Name))
}
