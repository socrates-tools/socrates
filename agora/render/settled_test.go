package render_test

import (
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/agora/document"
	"github.com/socrates-tools/socrates/agora/render"
	"golang.org/x/net/html"
)

func withWhatWasSettled() document.Document {
	read := oneDecision(document.One)
	read.Settled = []document.Settled{
		{ID: "q1", Name: "the surface", Location: "agora/surface", Question: "what is the page called", Chosen: []string{"the surface", "the page"}},
		{ID: "q2", Question: "should a design carry a cap", Words: "nothing here spawns a process for a cap to protect"},
		{ID: "q3", Question: "who writes the glossary entry"},
	}
	return read
}

func settledIn(page *html.Node) []*html.Node {
	var found []*html.Node
	for _, held := range elements(page, "li") {
		if hasClass(held, "settled-answer") {
			found = append(found, held)
		}
	}
	return found
}

func TestWhatWasSettledIsReadBeforeWhatIsStillOpen(t *testing.T) {
	read := withWhatWasSettled()

	written, err := render.Page(read, render.Frame{Base: "/sessions/a1b2c3d4", Token: "a-token", Revision: 1})
	if err != nil {
		t.Fatalf("rendering the page: %v", err)
	}
	settled := strings.Index(written, `class="settled"`)
	decided := strings.Index(written, `class="decisions"`)
	if settled < 0 || decided < 0 || settled > decided {
		t.Errorf("what was settled is at %d and the decisions at %d, and a person reads what is behind them first", settled, decided)
	}

	page := rendered(t, read)
	answers := settledIn(page)
	if len(answers) != 3 {
		t.Fatalf("three settled answers rendered %d", len(answers))
	}
	for position, expected := range [][]string{
		{"the surface", "agora/surface", "what is the page called", "the surface · the page"},
		{"q2", "should a design carry a cap", "“nothing here spawns a process for a cap to protect”"},
		{"q3", "who writes the glossary entry", "Left to whoever asked"},
	} {
		for _, wanted := range expected {
			if !strings.Contains(text(answers[position]), wanted) {
				t.Errorf("settled answer %d does not carry %q:\n%s", position+1, wanted, text(answers[position]))
			}
		}
	}
	if strings.Contains(text(answers[0]), "Left to whoever asked") || strings.Contains(text(answers[1]), "Left to whoever asked") {
		t.Error("a settled answer that chose something or said something is shown as left to whoever asked")
	}
}

func TestASettledAnswerOffersNothingToAnswer(t *testing.T) {
	page := rendered(t, withWhatWasSettled())

	for _, answer := range settledIn(page) {
		for _, control := range elements(answer, "input", "textarea", "button", "select") {
			t.Errorf("settled answer %q offers a <%s class=%q>, and a settled answer is read rather than answered",
				attribute(answer, "data-settled"), control.Data, attribute(control, "class"))
		}
	}
	for _, card := range cardsIn(page) {
		if named := attribute(card, "data-decision"); strings.HasPrefix(named, "q") {
			t.Errorf("settled answer %q rendered as a decision card, and a card is what a person answers", named)
		}
	}
}

func TestTheIndexLeadsToWhatWasSettledWithoutCountingIt(t *testing.T) {
	page := rendered(t, withWhatWasSettled())

	var led string
	for _, link := range elements(page, "a") {
		if hasClass(link, "settled-link") {
			led = attribute(link, "href")
			if said := strings.Join(strings.Fields(text(link)), " "); !strings.Contains(said, "3 settled earlier") || !strings.Contains(said, "the surface, q2, q3") {
				t.Errorf("the index says %q of what was settled, and it says how many and names each", said)
			}
		}
	}
	if led != "#settled" {
		t.Errorf("the index leads to %q for what was settled, and what was settled is rendered at %q", led, "#settled")
	}

	if entries := entriesIn(page); len(entries) != 1 {
		t.Errorf("the index lists %d entries for one decision and three settled answers, and only what is open is counted", len(entries))
	}
}

func TestAPageHoldingOnlyWhatWasSettledOffersNeitherControl(t *testing.T) {
	read := withWhatWasSettled()
	read.Decisions = nil

	page := rendered(t, read)
	if len(settledIn(page)) != 3 {
		t.Fatalf("a page holding only what was settled rendered %d settled answers", len(settledIn(page)))
	}
	for _, control := range elements(page, "input", "textarea", "button", "select", "form", "output", "dialog") {
		t.Errorf("a page holding only what was settled offers a <%s class=%q>, and there is nothing on it to answer", control.Data, attribute(control, "class"))
	}
}
