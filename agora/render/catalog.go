package render

import (
	"bytes"
	_ "embed"
	"html/template"
	"slices"
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

type anchorView struct {
	ID      string
	Tooltip string
	Label   string
	Title   string
	Hint    string
}

type state struct {
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

func dialogState(name string, view dialogView) state {
	return state{Name: name, Opens: view.ID, View: view}
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
		dialogState(dialogWithEverything, dialogView{
			ID:      "catalog-dialog",
			Class:   "catalog-dialog",
			Heading: catalogDialogHeading,
			Body:    "catalog-dialog-body",
			Actions: "catalog-dialog-actions",
		}),
		dialogState(dialogWithAHeadingAlone, dialogView{
			ID:      "catalog-dialog-bare",
			Class:   "catalog-dialog-bare",
			Heading: catalogDialogHeading,
		}),
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
			primitive.States = append(primitive.States, catalogState{Name: held.Name, Opens: held.Opens, Anchor: held.Anchor, Drawn: markup})
		}
		view.Primitives = append(view.Primitives, primitive)
	}

	var written bytes.Buffer
	if err := page.ExecuteTemplate(&written, "catalog", view); err != nil {
		return "", err
	}
	return written.String(), nil
}
