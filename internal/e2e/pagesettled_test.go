package e2e_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func whatWasSettled() []sketch {
	return []sketch{
		{"id": "q1", "name": "the surface", "location": "agora/surface", "question": "what is the page called", "chosen": []string{"the surface", "the page"}},
		{"id": "q2", "question": "should a design carry a cap", "words": "nothing here spawns a process for a cap to protect", "attached": []string{"before.png", "trace.txt"}},
		{"id": "q3", "question": "who writes the glossary entry"},
	}
}

func oneCacheDocumentWithWhatWasSettled() sketch {
	document := oneCacheDocument("single-choice")
	document["settled"] = whatWasSettled()
	return document
}

func onlyWhatWasSettled() sketch {
	return sketch{"title": "everything was settled", "settled": whatWasSettled()}
}

const whatWasSettledSays = `Array.from(document.querySelectorAll('.settled-answer')).map((held) => ({
  id: held.dataset.settled,
  text: held.textContent.replace(/\s+/g, ' ').trim(),
  attachments: Array.from(held.querySelectorAll('.settled-attachment'), (file) => file.textContent.trim()),
  lists: held.querySelectorAll('ul').length,
  controls: held.querySelectorAll('input, textarea, button, select, ui-button, ui-field, ui-choice').length,
}))`

type settledRead struct {
	ID          string   `json:"id"`
	Text        string   `json:"text"`
	Attachments []string `json:"attachments"`
	Lists       int      `json:"lists"`
	Controls    int      `json:"controls"`
}

func settledOf(page *reader) []settledRead {
	page.t.Helper()

	var read []settledRead
	page.evaluate(whatWasSettledSays, &read)
	return read
}

func TestWhatWasSettledIsReadBeforeWhatIsStillOpen(t *testing.T) {
	t.Parallel()
	_, page := read(t, written(t, oneCacheDocumentWithWhatWasSettled()))

	var before bool
	page.evaluate(`!!(document.querySelector('.settled').compareDocumentPosition(document.querySelector('.decisions')) & Node.DOCUMENT_POSITION_FOLLOWING)`, &before)
	if !before {
		t.Error("what was settled stands after the decisions, and a person reads what is behind them first")
	}

	answers := settledOf(page)
	if len(answers) != 3 {
		t.Fatalf("three settled answers drew %d", len(answers))
	}
	for position, expected := range [][]string{
		{"the surface", "agora/surface", "what is the page called", "the surface · the page"},
		{"q2", "should a design carry a cap", "“nothing here spawns a process for a cap to protect”", "before.png", "trace.txt"},
		{"q3", "who writes the glossary entry", "Left to whoever asked"},
	} {
		for _, wanted := range expected {
			if !strings.Contains(answers[position].Text, wanted) {
				t.Errorf("settled answer %d does not carry %q: %q", position+1, wanted, answers[position].Text)
			}
		}
	}
	if strings.Contains(answers[0].Text, "Left to whoever asked") || strings.Contains(answers[1].Text, "Left to whoever asked") {
		t.Error("a settled answer that chose something or said something is shown as left to whoever asked")
	}
}

func TestASettledAnswerNamesItsAttachmentsAndOneWithNoneShowsNoList(t *testing.T) {
	t.Parallel()
	_, page := read(t, written(t, oneCacheDocumentWithWhatWasSettled()))

	answers := settledOf(page)
	if diff := cmp.Diff([]string{"before.png", "trace.txt"}, answers[1].Attachments); diff != "" {
		t.Errorf("the settled answer carrying two attachments lists other files (-want +got):\n%s", diff)
	}
	for _, bare := range []int{0, 2} {
		if strings.Contains(answers[bare].Text, ".png") || answers[bare].Lists != 0 {
			t.Errorf("settled answer %d carries no attachment and shows a list: %q", bare+1, answers[bare].Text)
		}
	}
}

func TestASettledAnswerOffersNothingToAnswer(t *testing.T) {
	t.Parallel()
	_, page := read(t, written(t, oneCacheDocumentWithWhatWasSettled()))

	for _, answer := range settledOf(page) {
		if answer.Controls != 0 {
			t.Errorf("settled answer %q offers %d controls, and a settled answer is read rather than answered", answer.ID, answer.Controls)
		}
	}
	for _, card := range cardsOf(page) {
		if strings.HasPrefix(card.ID, "q") {
			t.Errorf("settled answer %q drew as a decision card, and a card is what a person answers", card.ID)
		}
	}
}

func TestTheIndexLeadsToWhatWasSettledWithoutCountingIt(t *testing.T) {
	t.Parallel()
	_, page := read(t, written(t, oneCacheDocumentWithWhatWasSettled()))

	var led struct {
		Said    string `json:"said"`
		Href    string `json:"href"`
		Entries int    `json:"entries"`
		Target  bool   `json:"target"`
	}
	page.evaluate(`(() => {
  const link = document.querySelector('.settled-link');
  return {
    said: link.textContent.replace(/\s+/g, ' ').trim(),
    href: link.getAttribute('href'),
    entries: document.querySelectorAll('.entry').length,
    target: !!document.querySelector('.settled#settled'),
  };
})()`, &led)

	if !strings.Contains(led.Said, "3 settled earlier") || !strings.Contains(led.Said, "the surface, q2, q3") {
		t.Errorf("the index says %q of what was settled, and it says how many and names each", led.Said)
	}
	if led.Href != "#settled" || !led.Target {
		t.Errorf("the index leads to %q for what was settled, and what was settled is drawn at %q: %v", led.Href, "#settled", led.Target)
	}
	if led.Entries != 1 {
		t.Errorf("the index lists %d entries for one decision and three settled answers, and only what is open is counted", led.Entries)
	}
}

func TestAPageHoldingOnlyWhatWasSettledOffersNeitherControl(t *testing.T) {
	t.Parallel()
	held := standing(t)
	page := openingFor(t, held.serve(written(t, onlyWhatWasSettled())), ".settled", ".settled")

	if answers := settledOf(page); len(answers) != 3 {
		t.Fatalf("a page holding only what was settled drew %d settled answers", len(answers))
	}
	var offered []string
	page.evaluate(`Array.from(document.querySelectorAll('input, textarea, button, select, form, output, dialog, ui-button, ui-field, ui-choice, ui-dialog'), (held) => held.localName + '.' + held.className)`, &offered)
	if len(offered) != 0 {
		t.Errorf("a page holding only what was settled offers %v, and there is nothing on it to answer", offered)
	}
}

func TestWhatWasSettledArrivesFoldedWhileAQuestionIsStillOpen(t *testing.T) {
	t.Parallel()
	_, page := read(t, written(t, oneCacheDocumentWithWhatWasSettled()))

	var open bool
	page.evaluate(`document.querySelector('details.settled').open`, &open)
	if open {
		t.Error("what was settled arrives unfolded beside a question still open, and a person reads it only when they ask to")
	}
}

func TestWhatWasSettledArrivesUnfoldedWhenNothingIsStillOpen(t *testing.T) {
	t.Parallel()
	held := standing(t)
	page := openingFor(t, held.serve(written(t, onlyWhatWasSettled())), ".settled", ".settled")

	var open bool
	page.evaluate(`document.querySelector('details.settled').open`, &open)
	if !open {
		t.Error("what was settled arrives folded on a page holding nothing else, and a person would see an empty page")
	}
}
