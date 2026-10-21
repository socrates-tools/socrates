package render_test

import (
	"testing"

	"github.com/socrates-tools/socrates/agora/document"
	"golang.org/x/net/html"
)

func firstWithClass(t *testing.T, page *html.Node, tag, class string) *html.Node {
	t.Helper()

	for _, held := range elements(page, tag) {
		if hasClass(held, class) {
			return held
		}
	}
	t.Fatalf("the page draws no %s carrying %q", tag, class)
	return nil
}

func TestAPieceKeepsTheClassThatNamesItBesideTheUtilitiesThatLayItOut(t *testing.T) {
	// scenario: A decision, the composer and the shell's panel carry their name class and their layout utilities
	page := rendered(t, oneDecision(document.One))

	for tag, class := range map[string]string{"ui-card": "decision", "ui-choice-group": "options-list"} {
		held := firstWithClass(t, page, tag, class)
		if !hasClass(held, "flex") && !hasClass(held, "overflow-hidden") && !hasClass(held, "panel:@pair:grid") {
			t.Errorf("the %s named %q carries no utility laying it out: %q", tag, class, attribute(held, "class"))
		}
	}
	if composer := firstWithClass(t, page, "ui-card", "composer"); !hasClass(composer, "flex-col") {
		t.Errorf("the composer carries %q, and its layout is written beside its name", attribute(composer, "class"))
	}
	if panel := firstWithClass(t, page, "aside", "panel"); !hasClass(panel, "panel:holding-unsent-composer:hidden!") {
		t.Errorf("the panel carries %q, and the shell's part takes its layout from the page", attribute(panel, "class"))
	}
}

func TestAnIndexOfEntriesNoneOfWhichIsDecidedLaysOutOneColumnAndHidesTheMarks(t *testing.T) {
	// scenario: An index whose entries have no decision behind them draws one column with its answered marks hidden
	page := rendered(t, document.Document{
		Title:  "what was seen",
		Source: []document.Entry{{ID: "s-1", Name: "Cache", Severity: "info", Location: "a.go:1"}},
	})

	link := firstWithClass(t, page, "span", "entry-link")
	if !hasClass(link, "grid-cols-1") || hasClass(link, "grid-cols-check") {
		t.Errorf("an undecided entry carries %q, and its index has nothing to track", attribute(link, "class"))
	}
	if mark := firstWithClass(t, page, "span", "answered"); !hasClass(mark, "hidden!") {
		t.Errorf("an undecided entry's mark carries %q, and it is hidden where nothing is tracked", attribute(mark, "class"))
	}
}

func TestAnIndexWithADecidedEntryLaysOutTheMarkColumnAndShowsTheMarks(t *testing.T) {
	// scenario: An index with a decided entry draws the mark column on every entry
	read := oneDecision(document.One)
	read.Source = []document.Entry{{ID: "c-1", Name: "Cache", Severity: "error"}}
	page := rendered(t, read)

	link := firstWithClass(t, page, "a", "entry-link")
	if !hasClass(link, "grid-cols-check") {
		t.Errorf("a decided entry carries %q, and its index tracks answers", attribute(link, "class"))
	}
	if mark := firstWithClass(t, page, "span", "answered"); hasClass(mark, "hidden!") {
		t.Errorf("a tracked entry's mark carries %q, and it is shown", attribute(mark, "class"))
	}
}
