package e2e_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const whatTheCardsAre = `Array.from(document.querySelectorAll('.decision')).map((held) => ({
  id: held.dataset.decision,
  anchor: held.id,
  severity: (` + severityOfAnElement + `)(held),
  text: held.textContent.replace(/\s+/g, ' ').trim(),
}))`

const whatTheEntriesAre = `Array.from(document.querySelectorAll('.entry')).map((held) => ({
  severity: (` + severityOfAnElement + `)(held),
  name: ((held.querySelector('.entry-name') || {}).textContent || '').trim(),
  text: held.textContent.replace(/\s+/g, ' ').trim(),
  href: held.querySelector('a') ? held.querySelector('a').getAttribute('href') : '',
  linked: held.querySelectorAll('a').length,
}))`

const theListComesBeforeTheDecisions = `!!(document.querySelector('.entries').compareDocumentPosition(document.querySelector('.decisions')) & Node.DOCUMENT_POSITION_FOLLOWING)`

type cardRead struct {
	ID       string `json:"id"`
	Anchor   string `json:"anchor"`
	Severity string `json:"severity"`
	Text     string `json:"text"`
}

type entryRead struct {
	Severity string `json:"severity"`
	Name     string `json:"name"`
	Text     string `json:"text"`
	Href     string `json:"href"`
	Linked   int    `json:"linked"`
}

func cardsOf(page *reader) []cardRead {
	page.t.Helper()

	var read []cardRead
	page.evaluate(whatTheCardsAre, &read)
	return read
}

func entriesOf(page *reader) []entryRead {
	page.t.Helper()

	var read []entryRead
	page.evaluate(whatTheEntriesAre, &read)
	return read
}

func idsOf(cards []cardRead) []string {
	var named []string
	for _, card := range cards {
		named = append(named, card.ID)
	}
	return named
}

func severitiesOfCards(cards []cardRead) []string {
	var named []string
	for _, card := range cards {
		named = append(named, card.Severity)
	}
	return named
}

func severitiesOfEntries(entries []entryRead) []string {
	var named []string
	for _, entry := range entries {
		named = append(named, entry.Severity)
	}
	return named
}

func TestTheMostSevereCardIsTheFirstAPersonMeets(t *testing.T) {
	t.Parallel()
	_, page := read(t, written(t, leastSevereFirst()))

	cards := cardsOf(page)
	if diff := cmp.Diff([]string{"c-1", "c-2", "c-3"}, idsOf(cards)); diff != "" {
		t.Errorf("the cards arrive in an order a person cannot triage by severity (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]string{"error", "warning", "info"}, severitiesOfCards(cards)); diff != "" {
		t.Errorf("the cards arrive carrying other severities than their own (-want +got):\n%s", diff)
	}
}

func TestTheSourceListIsOrderedTheSameWay(t *testing.T) {
	t.Parallel()
	_, page := read(t, written(t, leastSevereFirst()))

	if diff := cmp.Diff([]string{"error", "warning", "info"}, severitiesOfEntries(entriesOf(page))); diff != "" {
		t.Errorf("the entries are listed in another order than the cards (-want +got):\n%s", diff)
	}

	var before bool
	page.evaluate(theListComesBeforeTheDecisions, &before)
	if !before {
		t.Error("the list stands after the decisions, and ordering them changed which comes first")
	}
}

func TestEqualSeveritiesKeepTheOrderTheDocumentListedThem(t *testing.T) {
	t.Parallel()
	source := leastSevereFirstSource()
	for _, entry := range source {
		entry["severity"] = "warning"
	}
	_, page := read(t, written(t, aDocument("three things of one severity", source, aDecision("c-3", nil), aDecision("c-2", nil), aDecision("c-1", nil))))

	if diff := cmp.Diff([]string{"c-3", "c-2", "c-1"}, idsOf(cardsOf(page))); diff != "" {
		t.Errorf("three decisions of one severity do not arrive in the order the document listed them (-want +got):\n%s", diff)
	}
	if entries := entriesOf(page); len(entries) != 3 || entries[0].Name != "the log line says nothing" {
		t.Errorf("three entries of one severity are not listed in the order the document listed them: %+v", entries)
	}
}

func TestTheCallerCannotOrderThePage(t *testing.T) {
	t.Parallel()
	_, page := read(t, written(t, aDocument("three things to settle", leastSevereFirstSource(),
		aDecision("c-3", sketch{"question": "the informational one the caller listed first"}), aDecision("c-2", nil), aDecision("c-1", nil))))

	cards := cardsOf(page)
	if len(cards) != 3 || cards[0].ID != "c-1" {
		t.Fatalf("the first card a person meets is %+v, and the error is c-1", cards)
	}
	if strings.Contains(cards[0].Text, "the informational one the caller listed first") {
		t.Errorf("the card the caller listed first is the card a person meets first: %q", cards[0].Text)
	}
}

func TestAnEntryWhoseSourceStatesNoSeveritySortsLast(t *testing.T) {
	t.Parallel()
	source := append([]sketch{anEntry("c-0", "nobody classified this one", "", "the author patches", "")}, leastSevereFirstSource()...)
	_, page := read(t, written(t, aDocument("three things to settle", source, aDecision("c-3", nil), aDecision("c-2", nil), aDecision("c-1", nil))))

	entries := entriesOf(page)
	if diff := cmp.Diff([]string{"error", "warning", "info", ""}, severitiesOfEntries(entries)); diff != "" {
		t.Errorf("an unclassified entry is not listed last, and it is not evidence of unimportance (-want +got):\n%s", diff)
	}
	if last := entries[len(entries)-1]; !strings.Contains(last.Text, "nobody classified this one") {
		t.Errorf("the entry listed last is %q", last.Text)
	}
}

func TestTheSourceIsReadBeforeTheDecisions(t *testing.T) {
	t.Parallel()
	_, page := read(t, written(t, citingThreeEntries()))

	entries := entriesOf(page)
	if len(entries) != 3 {
		t.Fatalf("a document citing three source entries listed %d of them", len(entries))
	}
	for position, expected := range [][]string{
		{"the read is unlocked", "internal/cache/read.go:41"},
		{"the write counts twice", "internal/cache/write.go:12"},
		{"c-3", "internal/cache/read.go:80"},
	} {
		for _, wanted := range expected {
			if !strings.Contains(entries[position].Text, wanted) {
				t.Errorf("entry %d does not carry %q: %q", position+1, wanted, entries[position].Text)
			}
		}
	}
	if diff := cmp.Diff([]string{"error", "warning", "info"}, severitiesOfEntries(entries)); diff != "" {
		t.Errorf("the entries do not each carry their own source's severity (-want +got):\n%s", diff)
	}

	var before bool
	page.evaluate(theListComesBeforeTheDecisions, &before)
	if !before {
		t.Error("the list stands after the decisions, and a person reads the list first")
	}
}

func TestAnEntryNoDecisionCitesIsStillShown(t *testing.T) {
	t.Parallel()
	_, page := read(t, written(t, citingThreeEntries()))

	if entries := entriesOf(page); len(entries) != 3 {
		t.Errorf("three entries and one decision listed %d entries: %+v", len(entries), entries)
	}
	if cards := cardsOf(page); len(cards) != 1 {
		t.Errorf("three entries and one decision drew %d decisions, and an entry a decision covers is not a second decision", len(cards))
	}
}

func TestAnEntryLeadsToTheDecisionThatCoversIt(t *testing.T) {
	t.Parallel()
	_, page := read(t, written(t, citingThreeEntries()))

	leading := map[string]string{}
	for _, entry := range entriesOf(page) {
		if entry.Linked > 0 {
			leading[entry.Name] = entry.Href
		}
	}
	if diff := cmp.Diff(map[string]string{"the read is unlocked": "#decision-c-1"}, leading); diff != "" {
		t.Errorf("only the entry a decision covers has anywhere to lead (-want +got):\n%s", diff)
	}

	cards := cardsOf(page)
	if len(cards) != 1 || cards[0].Anchor != "decision-c-1" {
		t.Errorf("the entry leads to an identifier the page does not carry: %+v", cards)
	}
}

func TestTheListIsNotTheCallersProse(t *testing.T) {
	t.Parallel()
	document := citingThreeEntries()
	document["preamble"] = "there were three things and none of them matters"
	_, page := read(t, written(t, document))

	entries := entriesOf(page)
	for _, entry := range entries {
		if strings.Contains(entry.Text, "none of them matters") {
			t.Errorf("an entry carries the preamble the caller wrote: %q", entry.Text)
		}
	}
	for _, wanted := range []string{"the read is unlocked", "the write counts twice", "c-3"} {
		if !slices.ContainsFunc(entries, func(entry entryRead) bool { return strings.Contains(entry.Text, wanted) }) {
			t.Errorf("the list does not carry %q, and what it lists comes from the entries", wanted)
		}
	}
}

func TestADocumentCitingNoSourceIndexesItsDecisions(t *testing.T) {
	t.Parallel()
	_, page := read(t, written(t, aDocument("three things to settle", nil,
		aDecision("c-3", sketch{"severity": "info", "name": "the log line says nothing"}),
		aDecision("c-2", sketch{"severity": "warning", "name": "the write counts twice"}),
		aDecision("c-1", sketch{"severity": "error", "name": "the read is unlocked"}))))

	entries := entriesOf(page)
	if diff := cmp.Diff([]string{"error", "warning", "info"}, severitiesOfEntries(entries)); diff != "" {
		t.Fatalf("a document citing no source does not index one entry per decision, most severe first (-want +got):\n%s", diff)
	}
	for position, card := range cardsOf(page) {
		if entries[position].Href != "#"+card.Anchor {
			t.Errorf("entry %d leads to %q, and not to the card at the same place, %q", position+1, entries[position].Href, card.Anchor)
		}
	}
}
