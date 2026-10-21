package e2e_test

import (
	"strings"
	"testing"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
	"github.com/google/go-cmp/cmp"
)

const whatTheIndexSays = `(() => ({
  said: document.querySelector('.progress-said').textContent.replace(/\s+/g, ' ').trim(),
  bar: document.querySelector('.progress-bar').value,
  entries: Object.fromEntries(Array.from(document.querySelectorAll('a.entry-link')).map((held) => [held.dataset.decision, held.querySelector('.answered').dataset.answered])),
  cards: Object.fromEntries(Array.from(document.querySelectorAll('.decision')).map((held) => [held.dataset.decision, held.querySelector('.attribution .answered').dataset.answered])),
}))()`

type indexed struct {
	Said    string            `json:"said"`
	Bar     float64           `json:"bar"`
	Entries map[string]string `json:"entries"`
	Cards   map[string]string `json:"cards"`
}

func indexAfter(t *testing.T, page *reader) indexed {
	t.Helper()

	var read indexed
	page.evaluate(whatTheIndexSays, &read)
	if diff := cmp.Diff(read.Cards, read.Entries); diff != "" {
		t.Errorf("the index and the cards disagree about what is answered (-cards +index):\n%s", diff)
	}
	return read
}

func TestTheIndexCountsADecisionAsThePersonAnswersIt(t *testing.T) {
	t.Parallel()
	_, page := read(t, twoDecisionDocument)

	arrived := indexAfter(t, page)
	if arrived.Said != "1 of 2 answered" || arrived.Bar != 1 {
		t.Errorf("a page whose one recommendation arrives chosen says %q with the bar at %v, want 1 of 2", arrived.Said, arrived.Bar)
	}
	if diff := cmp.Diff(map[string]string{"c-1": "no", "c-2": "yes"}, arrived.Cards); diff != "" {
		t.Errorf("what arrives answered differs (-want +got):\n%s", diff)
	}

	page.evaluate(chose("c-1", 0), nil)
	if chosen := indexAfter(t, page); chosen.Said != "2 of 2 answered" || chosen.Cards["c-1"] != "yes" {
		t.Errorf("choosing an option leaves the index saying %q and c-1 marked %q", chosen.Said, chosen.Cards["c-1"])
	}

	page.evaluate(chose("c-2", 2), nil)
	if unticked := indexAfter(t, page); unticked.Said != "1 of 2 answered" || unticked.Cards["c-2"] != "no" {
		t.Errorf("unticking the only box leaves the index saying %q and c-2 marked %q", unticked.Said, unticked.Cards["c-2"])
	}

	page.write(within("c-2", ".words-field"), "only the write")
	if written := indexAfter(t, page); written.Said != "2 of 2 answered" || written.Cards["c-2"] != "yes" {
		t.Errorf("words alone leave the index saying %q and c-2 marked %q, and words alone are an answer", written.Said, written.Cards["c-2"])
	}

	page.evaluate(declined("c-1"), nil)
	if declined := indexAfter(t, page); declined.Said != "1 of 2 answered" || declined.Cards["c-1"] != "no" {
		t.Errorf("taking none of the options leaves the index saying %q and c-1 marked %q, and taking none alone answers nothing", declined.Said, declined.Cards["c-1"])
	}
}

const whatTheTipSays = `(() => {
  const tip = document.querySelector('.tooltip');
  return { shown: !tip.hidden && tip.getBoundingClientRect().height > 0, said: (tip.heading + ' ' + tip.hint).trim() };
})()`

type tipped struct {
	Shown bool   `json:"shown"`
	Said  string `json:"said"`
}

func TestAMarkSaysWhatItMeansWhileItHasTheCursor(t *testing.T) {
	t.Parallel()
	_, page := read(t, glancedDocument)
	page.do(sending(emulation.SetFocusEmulationEnabled, emulation.SetFocusEmulationEnabledParams{Enabled: true}))

	var before tipped
	page.evaluate(whatTheTipSays, &before)
	if before.Shown {
		t.Fatalf("the page opens with a mark's meaning showing: %q", before.Said)
	}

	page.do(chromedp.Focus(chromedp.CSS(within("c-1", ".mark.severity"))))

	var focused tipped
	page.evaluate(whatTheTipSays, &focused)
	if !focused.Shown || !strings.Contains(focused.Said, "Severity:") {
		t.Errorf("a severity mark with the cursor on it shows %+v, and it says what the severity means", focused)
	}

	page.do(chromedp.Focus(chromedp.CSS(within("c-1", ".words-field"))))

	var left tipped
	page.evaluate(whatTheTipSays, &left)
	if left.Shown {
		t.Errorf("a mark's meaning is still showing after the cursor left it: %q", left.Said)
	}
}

const whereTheExchangeStands = `(() => {
  const status = document.querySelector('.status');
  const sends = document.querySelector('.sends');
  return {
    status: status.dataset.status,
    said: status.textContent.trim(),
    presses: getComputedStyle(sends).display !== 'none',
  };
})()`

type standingOf struct {
	Status  string `json:"status"`
	Said    string `json:"said"`
	Presses bool   `json:"presses"`
}

func TestTheHeaderSaysWhereTheExchangeStands(t *testing.T) {
	t.Parallel()
	_, page := read(t, twoDecisionDocument)
	page.evaluate(chose("c-1", 0), nil)

	var open standingOf
	page.evaluate(whereTheExchangeStands, &open)
	if open != (standingOf{Status: "open", Said: "Open", Presses: true}) {
		t.Errorf("an open page stands as %+v", open)
	}

	page.click(".send")
	page.waitFor(`!document.querySelector('.receipt').hidden`)

	var sent standingOf
	page.evaluate(whereTheExchangeStands, &sent)
	if sent != (standingOf{Status: "sent", Said: "Sent, waiting for a reply", Presses: false}) {
		t.Errorf("a page that sent and waits stands as %+v, and the presses give way to the way back", sent)
	}

	page.click(".amend")
	page.waitFor(`document.querySelector('.receipt').hidden`)
	page.click(".finish")
	page.waitFor(`document.querySelector('.send').disabled === true`)

	var over standingOf
	page.evaluate(whereTheExchangeStands, &over)
	if over != (standingOf{Status: "over", Said: "Exchange over", Presses: false}) {
		t.Errorf("a page that finished stands as %+v", over)
	}
}

const settledDocument = `{
  "title": "one question left",
  "settled": [
    {"id": "q1", "question": "what should the loopback server be called", "chosen": ["the surface"]}
  ],
  "decisions": [
    {"id": "q2", "question": "should the glossary gain an entry today", "mode": "single-choice",
     "options": [{"label": "write it now", "meaning": "the term is settled"}, {"label": "wait", "meaning": "batch it"}]}
  ]
}`

func TestWhatWasSettledArrivesFoldedAndOpensFromItsHeadingOrTheIndex(t *testing.T) {
	t.Parallel()
	_, page := read(t, settledDocument)

	const answersShown = `document.querySelector('.settled-answer').checkVisibility()`
	const headingShown = `document.querySelector('.settled-summary').checkVisibility()`

	var shown bool
	page.evaluate(headingShown, &shown)
	if !shown {
		t.Fatal("the heading of what was settled is not shown, and a person could not open it")
	}
	page.evaluate(answersShown, &shown)
	if shown {
		t.Fatal("what was settled arrives unfolded, and it arrives folded while a question is still open")
	}

	page.click(".settled-summary")
	page.evaluate(answersShown, &shown)
	if !shown {
		t.Error("pressing the heading of what was settled left it folded")
	}

	page.click(".settled-summary")
	page.evaluate(answersShown, &shown)
	if shown {
		t.Error("pressing the heading of what was settled again left it showing")
	}

	page.click(".settled-link")
	page.waitFor(answersShown)
}

func TestTheIndexLinkToWhatWasSettledIsNotLeftUnderTheTopBarWhenItIsScrolledToTheTopOfTheWindow(t *testing.T) {
	t.Parallel()
	_, page := read(t, settledDocument)

	page.click(".settled-summary")
	page.evaluate(`document.querySelector('.settled-link').scrollIntoView({ block: 'start' })`, nil)

	var covered bool
	page.evaluate(`(() => {
		const link = document.querySelector('.settled-link')
		const box = link.getBoundingClientRect()
		const hit = document.elementFromPoint(box.x + box.width / 2, box.y + box.height / 2)
		return !link.contains(hit)
	})()`, &covered)
	if covered {
		t.Error("the link to what was settled, scrolled to the top of the window, sits under the top bar and a press on it lands on the bar")
	}
}

const quotedIdentifierDocument = `{
  "title": "one thing to settle",
  "decisions": [
    {"id": "a\"b\\c", "question": "should the cache be read before the lock is taken", "mode": "single-choice",
     "options": [{"label": "take the lock first", "meaning": "correctness first"}, {"label": "leave it", "meaning": "acceptable here"}],
     "evidence": [{"language": "go", "code": "if held {\n\treturn entry\n}\n"}]}
  ]
}`

func TestAnIdentifierHoldingAQuoteLeavesThePageWorking(t *testing.T) {
	t.Parallel()
	_, page := read(t, quotedIdentifierDocument)

	page.waitFor(aCodeBlockIsHighlighted)
	page.evaluate(`document.querySelector('.choice').click()`, nil)

	if chosen := indexAfter(t, page); chosen.Said != "1 of 1 answered" {
		t.Errorf("choosing an option on a decision whose identifier holds a quote leaves the index saying %q", chosen.Said)
	}

	page.click(".send")
	page.waitFor(`!document.querySelector('.receipt').hidden`)

	var said string
	page.evaluate(`document.querySelector('.receipt-decisions').textContent`, &said)
	if !strings.Contains(said, "take the lock first") {
		t.Errorf("the receipt for a decision whose identifier holds a quote says %q, and it names what was chosen", said)
	}
}

const whereTheBarAndTheCardSit = `(() => {
  const bar = document.querySelector('.topbar').getBoundingClientRect();
  const status = document.querySelector('.status').getBoundingClientRect();
  const card = document.querySelector('.decision[data-decision="c-2"]').getBoundingClientRect();
  return { barBottom: bar.bottom, statusTop: status.top, statusBottom: status.bottom, cardTop: card.top, height: window.innerHeight, scrolled: window.scrollY };
})()`

type sitting struct {
	BarBottom    float64 `json:"barBottom"`
	StatusTop    float64 `json:"statusTop"`
	StatusBottom float64 `json:"statusBottom"`
	CardTop      float64 `json:"cardTop"`
	Height       float64 `json:"height"`
	Scrolled     float64 `json:"scrolled"`
}

func TestTheStatusStaysInViewAsThePageScrolls(t *testing.T) {
	t.Parallel()
	_, page := read(t, twoDecisionDocument)

	page.evaluate(`window.scrollTo(0, document.documentElement.scrollHeight)`, nil)

	var read sitting
	page.evaluate(whereTheBarAndTheCardSit, &read)
	if read.Scrolled == 0 {
		t.Fatal("the page did not scroll, and a page that cannot scroll says nothing about what stays in view")
	}
	if read.StatusTop < 0 || read.StatusBottom > read.Height {
		t.Errorf("scrolled to the end, the status sits at %v to %v in a window %v high, and it stays in view", read.StatusTop, read.StatusBottom, read.Height)
	}
}

func TestAnEntryLandsTheCardClearOfTheBar(t *testing.T) {
	t.Parallel()
	_, page := read(t, twoDecisionDocument)
	page.do(sending(emulation.SetEmulatedMedia, emulation.SetEmulatedMediaParams{Features: []*emulation.MediaFeature{{Name: "prefers-reduced-motion", Value: "reduce"}}}))

	page.evaluate(`window.scrollTo(0, document.documentElement.scrollHeight)`, nil)
	page.evaluate(`document.querySelector('a.entry-link[data-decision="c-2"]').click()`, nil)
	page.waitFor(`document.activeElement.closest('.decision')?.dataset.decision === 'c-2'`)

	var read sitting
	page.evaluate(whereTheBarAndTheCardSit, &read)
	if read.CardTop < read.BarBottom {
		t.Errorf("following the entry left the card's top at %vpx, under the bar that ends at %vpx", read.CardTop, read.BarBottom)
	}
	if read.CardTop > read.Height/2 {
		t.Errorf("following the entry left the card's top at %vpx in a window %v high, and it brings the card to the top", read.CardTop, read.Height)
	}
}
