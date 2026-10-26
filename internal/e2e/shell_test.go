package e2e_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)

const (
	aWindowWideEnoughForTheShell = 1100
	aWindowWideEnoughForThePanel = 1400
	aShortWindow                 = 700
)

func openingAt(t *testing.T, document string, width, height int) *reader {
	t.Helper()
	held := standing(t)
	return opening(t, held.serve(document), chromedp.WindowSize(width, height))
}

const howTheShellScrolls = `(() => {
  const main = document.querySelector('.main');
  const index = document.querySelector('.index');
  const bar = document.querySelector('.topbar').getBoundingClientRect();
  return {
    pageScrolls: document.documentElement.scrollHeight > window.innerHeight + 1,
    mainScrolls: getComputedStyle(main).overflowY === 'auto' && main.scrollHeight > main.clientHeight,
    indexScrolls: getComputedStyle(index).overflowY === 'auto',
    barLeft: bar.left,
    barRight: bar.right,
    width: document.documentElement.clientWidth,
  };
})()`

type shellScrolling struct {
	PageScrolls  bool    `json:"pageScrolls"`
	MainScrolls  bool    `json:"mainScrolls"`
	IndexScrolls bool    `json:"indexScrolls"`
	BarLeft      float64 `json:"barLeft"`
	BarRight     float64 `json:"barRight"`
	Width        float64 `json:"width"`
}

func TestAWideWindowHoldsTheShellAndOnlyItsPartsScroll(t *testing.T) {
	t.Parallel()
	page := openingAt(t, twoDecisionDocument, aWindowWideEnoughForTheShell, aShortWindow)

	var read shellScrolling
	page.evaluate(howTheShellScrolls, &read)

	if read.PageScrolls {
		t.Error("in a wide window the page itself scrolls, and the shell fills the window with its parts scrolling on their own")
	}
	if !read.MainScrolls {
		t.Error("in a wide window the decisions do not scroll on their own")
	}
	if !read.IndexScrolls {
		t.Error("in a wide window the navigation does not scroll on its own")
	}
	if read.BarLeft != 0 || read.BarRight < read.Width-1 {
		t.Errorf("the top bar runs from %vpx to %vpx in a window %vpx wide, and it runs the full width", read.BarLeft, read.BarRight, read.Width)
	}
}

func TestTheShellOfThePageHasAShadowRoot(t *testing.T) {
	// scenario: The page's shell renders into a shadow root and slots the topbar, the navigation, the content and the panel
	t.Parallel()
	page := openingAt(t, twoDecisionDocument, aWindowWideEnoughForThePanel, aShortWindow)

	var read struct {
		Shadowed bool     `json:"shadowed"`
		Slotted  []string `json:"slotted"`
	}
	page.evaluate(`(() => {
  const shell = document.querySelector('ui-app-shell');
  return {
    shadowed: shell.shadowRoot !== null,
    slotted: Array.from(shell.children).map((child) => child.localName + '[' + child.getAttribute('slot') + ']'),
  };
})()`, &read)

	if !read.Shadowed {
		t.Error("the page's shell has no shadow root")
	}
	if want := []string{"header[topbar]", "nav[navigation]", "main[content]", "aside[panel]"}; !slices.Equal(read.Slotted, want) {
		t.Errorf("the shell holds %v, want %v", read.Slotted, want)
	}
}

func TestTheTopBarCarriesTheTitleAsABreadcrumb(t *testing.T) {
	t.Parallel()
	page := openingAt(t, twoDecisionDocument, aWindowWideEnoughForTheShell, aShortWindow)

	var crumb string
	page.evaluate(`document.querySelector('.topbar .crumb').textContent.trim()`, &crumb)
	if crumb != "two things to settle" {
		t.Errorf("the top bar's breadcrumb says %q, and it carries the exchange's title", crumb)
	}
}

const whereTheReplyPanelSits = `(() => {
  const panel = document.querySelector('.panel');
  const send = document.querySelector('.send').getBoundingClientRect();
  const box = panel.getBoundingClientRect();
  return {
    shown: panel.checkVisibility(),
    holdsTheComposer: !!panel.querySelector('.composer'),
    panelRight: box.right,
    width: document.documentElement.clientWidth,
    sendBottom: send.bottom,
    sendTop: send.top,
    height: window.innerHeight,
  };
})()`

type replyPanel struct {
	Shown            bool    `json:"shown"`
	HoldsTheComposer bool    `json:"holdsTheComposer"`
	PanelRight       float64 `json:"panelRight"`
	Width            float64 `json:"width"`
	SendBottom       float64 `json:"sendBottom"`
	SendTop          float64 `json:"sendTop"`
	Height           float64 `json:"height"`
}

func TestFromTwelveEightyTheReplyPanelKeepsThePressesInReach(t *testing.T) {
	t.Parallel()
	page := openingAt(t, twoDecisionDocument, aWindowWideEnoughForThePanel, aShortWindow)

	var before replyPanel
	page.evaluate(whereTheReplyPanelSits, &before)
	if !before.Shown || !before.HoldsTheComposer {
		t.Fatalf("in a window %dpx wide the reply panel is %+v, and it holds the message and the presses", aWindowWideEnoughForThePanel, before)
	}
	if before.PanelRight < before.Width-1 {
		t.Errorf("the reply panel ends at %vpx in a window %vpx wide, and it sits at the right of the window", before.PanelRight, before.Width)
	}
	if before.SendTop < 0 || before.SendBottom > before.Height {
		t.Errorf("before any scrolling the send press sits at %v to %v in a window %v high, and it is in reach", before.SendTop, before.SendBottom, before.Height)
	}

	page.evaluate(`document.querySelector('.main').scrollTo(0, document.querySelector('.main').scrollHeight)`, nil)

	var after replyPanel
	page.evaluate(whereTheReplyPanelSits, &after)
	if after.SendTop != before.SendTop {
		t.Errorf("scrolling the decisions moved the send press from %vpx to %vpx, and it stays pinned", before.SendTop, after.SendTop)
	}
}

const whereTheComposerSits = `(() => {
  const composer = document.querySelector('.composer');
  const decisions = document.querySelector('.decisions');
  return {
    inTheMain: !!composer.closest('.main'),
    afterTheDecisions: decisions.compareDocumentPosition(composer) === Node.DOCUMENT_POSITION_FOLLOWING,
    panelShown: document.querySelector('.panel').checkVisibility(),
  };
})()`

type composerSitting struct {
	InTheMain         bool `json:"inTheMain"`
	AfterTheDecisions bool `json:"afterTheDecisions"`
	PanelShown        bool `json:"panelShown"`
}

func TestBelowTwelveEightyTheComposerClosesTheDecisions(t *testing.T) {
	t.Parallel()
	page := openingAt(t, twoDecisionDocument, aWindowWideEnoughForTheShell, aShortWindow)

	var read composerSitting
	page.evaluate(whereTheComposerSits, &read)
	if !read.InTheMain || !read.AfterTheDecisions {
		t.Errorf("in a window %dpx wide the composer sits %+v, and it closes the decisions", aWindowWideEnoughForTheShell, read)
	}
	if read.PanelShown {
		t.Errorf("in a window %dpx wide the reply panel is shown, and it joins from 1280px", aWindowWideEnoughForTheShell)
	}
}

func TestTheComposerFollowsTheWindowAcrossTwelveEighty(t *testing.T) {
	// scenario: The composer moves between the reply panel and the end of the decisions as the window changes
	t.Parallel()
	page := openingAt(t, twoDecisionDocument, aWindowWideEnoughForThePanel, aShortWindow)
	page.evaluate(`document.querySelector('.message-field').value = 'kept across the move'`, nil)

	page.do(sending(emulation.SetDeviceMetricsOverride, emulation.SetDeviceMetricsOverrideParams{Width: aWindowWideEnoughForTheShell, Height: aShortWindow, DeviceScaleFactor: 1}))
	page.waitFor(`!!document.querySelector('.main .composer')`)

	var kept string
	page.evaluate(`document.querySelector('.message-field').value`, &kept)
	if kept != "kept across the move" {
		t.Errorf("moving the composer out of the reply panel left the message as %q", kept)
	}

	page.do(sending(emulation.SetDeviceMetricsOverride, emulation.SetDeviceMetricsOverrideParams{Width: aWindowWideEnoughForThePanel, Height: aShortWindow, DeviceScaleFactor: 1}))
	page.waitFor(`!!document.querySelector('.panel .composer')`)
}

const whatTheSummarySays = `Array.from(document.querySelectorAll('.panel .summary-row')).map((row) => ({
  said: row.querySelector('.summary-said').textContent.trim(),
  answered: row.querySelector('.answered').dataset.answered,
}))`

type summaryRow struct {
	Said     string `json:"said"`
	Answered string `json:"answered"`
}

func TestTheReplyPanelSummarisesWhatWillGoAsThePersonAnswers(t *testing.T) {
	t.Parallel()
	page := openingAt(t, twoDecisionDocument, aWindowWideEnoughForThePanel, aShortWindow)

	var before []summaryRow
	page.evaluate(whatTheSummarySays, &before)
	if len(before) != 2 {
		t.Fatalf("the summary holds %d rows for two decisions: %+v", len(before), before)
	}
	if before[0].Answered != "no" || !strings.HasSuffix(before[0].Said, ": nothing") {
		t.Errorf("a decision nobody answered is summarised as %+v", before[0])
	}
	if before[1].Answered != "yes" || !strings.Contains(before[1].Said, "the log line") {
		t.Errorf("a decision whose recommendation arrived chosen is summarised as %+v", before[1])
	}

	page.evaluate(chose("c-1", 0), nil)
	page.evaluate(`(() => { const field = document.querySelector('.decision[data-decision="c-1"] .words-field'); field.value = 'and say why'; field.dispatchEvent(new Event('input', { bubbles: true })); })()`, nil)

	var after []summaryRow
	page.evaluate(whatTheSummarySays, &after)
	if after[0].Answered != "yes" || !strings.Contains(after[0].Said, "take the lock first, and say why") {
		t.Errorf("a decision answered with an option and words is summarised as %+v", after[0])
	}
}

const whichEntryIsCurrent = `Array.from(document.querySelectorAll('a.entry-link[aria-current="location"]')).map((link) => link.dataset.decision)`

func TestTheNavigationMarksTheDecisionInView(t *testing.T) {
	t.Parallel()
	page := openingAt(t, twoDecisionDocument, aWindowWideEnoughForTheShell, aShortWindow)

	page.waitFor(`document.querySelector('a.entry-link[data-decision="c-1"]').getAttribute('aria-current') === 'location'`)

	page.evaluate(`(() => { const main = document.querySelector('.main'); const card = document.querySelector('.decision[data-decision="c-2"]'); main.scrollBy(0, card.getBoundingClientRect().top - main.getBoundingClientRect().top); })()`, nil)
	page.waitFor(`document.querySelector('a.entry-link[data-decision="c-2"]').getAttribute('aria-current') === 'location'`)

	var current []string
	page.evaluate(whichEntryIsCurrent, &current)
	if len(current) != 1 || current[0] != "c-2" {
		t.Errorf("with the second decision in view the navigation marks %v as current", current)
	}
}

const whereTheCardLandsInTheMain = `(() => {
  const main = document.querySelector('.main').getBoundingClientRect();
  const card = document.querySelector('.decision[data-decision="c-2"]').getBoundingClientRect();
  return { mainTop: main.top, cardTop: card.top };
})()`

func TestAnEntryLandsItsCardAtTheTopOfTheDecisions(t *testing.T) {
	t.Parallel()
	page := openingAt(t, twoDecisionDocument, aWindowWideEnoughForTheShell, aShortWindow)
	page.do(sending(emulation.SetEmulatedMedia, emulation.SetEmulatedMediaParams{Features: []*emulation.MediaFeature{{Name: "prefers-reduced-motion", Value: "reduce"}}}))

	page.evaluate(`document.querySelector('a.entry-link[data-decision="c-2"]').click()`, nil)
	page.waitFor(`document.activeElement.closest('.decision')?.dataset.decision === 'c-2'`)

	var landed struct {
		MainTop float64 `json:"mainTop"`
		CardTop float64 `json:"cardTop"`
	}
	page.evaluate(whereTheCardLandsInTheMain, &landed)
	if apart := landed.CardTop - landed.MainTop; apart < 0 || apart > 48 {
		t.Errorf("following the entry left the card %vpx below the top of the decisions, and it lands the card at their top", apart)
	}
}

const howThePanelHoldsThePresses = `(() => {
  const panel = document.querySelector('.panel').getBoundingClientRect();
  const sends = document.querySelector('.panel .sends');
  const send = document.querySelector('.send').getBoundingClientRect();
  const finish = document.querySelector('.finish').getBoundingClientRect();
  return {
    width: panel.width,
    sameRow: send.top === finish.top,
    fits: sends.scrollWidth <= sends.clientWidth && finish.left >= panel.left && send.right <= panel.right,
  };
})()`

func TestTheReplyPanelIsThreeHundredAndFourWideAndHoldsBothPressesOnOneRow(t *testing.T) {
	t.Parallel()
	page := openingAt(t, twoDecisionDocument, aWindowWideEnoughForThePanel, aShortWindow)

	var read struct {
		Width   float64 `json:"width"`
		SameRow bool    `json:"sameRow"`
		Fits    bool    `json:"fits"`
	}
	page.evaluate(howThePanelHoldsThePresses, &read)

	if read.Width != 304 {
		t.Errorf("the reply panel is %vpx wide, and it is 304px", read.Width)
	}
	if !read.SameRow || !read.Fits {
		t.Errorf("the presses in the reply panel sit %+v, and both fit side by side inside it", read)
	}
}

const howThePressSpacesItsWords = `(() => {
  const press = document.querySelector('ui-button.send');
  const walker = document.createTreeWalker(press, NodeFilter.SHOW_TEXT);
  const pieces = [];
  for (let node = walker.nextNode(); node; node = walker.nextNode()) {
    for (let at = 0; at < node.data.length; at++) {
      const range = document.createRange();
      range.setStart(node, at);
      range.setEnd(node, at + 1);
      const box = range.getBoundingClientRect();
      pieces.push({ letter: node.data[at], left: box.left, right: box.right, drawn: box.width > 0 });
    }
  }
  const drawn = pieces.filter((piece) => piece.drawn && piece.letter.trim() !== '');
  const gaps = [];
  for (let at = 1; at < pieces.length; at++) {
    if (pieces[at].letter === ' ') {
      const before = pieces.slice(0, at).filter((piece) => piece.letter.trim() !== '').at(-1);
      const after = pieces.slice(at).find((piece) => piece.letter.trim() !== '');
      if (before && after) {
        gaps.push(after.left - before.right);
      }
    }
  }
  return { said: drawn.map((piece) => piece.letter).join(''), gaps };
})()`

type pressSpacing struct {
	Said string    `json:"said"`
	Gaps []float64 `json:"gaps"`
}

func TestBelowTwelveEightyTheSendPressSpacesForAReplyLikeEveryOtherWord(t *testing.T) {
	t.Parallel()
	page := openingAt(t, twoDecisionDocument, aWindowWideEnoughForTheShell, aShortWindow)

	var read pressSpacing
	page.evaluate(howThePressSpacesItsWords, &read)
	if read.Said != "Sendandwaitforareply" || len(read.Gaps) != 5 {
		t.Fatalf("in a window %dpx wide the send press draws %q with gaps %v, and it says Send and wait for a reply", aWindowWideEnoughForTheShell, read.Said, read.Gaps)
	}
	for at, gap := range read.Gaps[1:] {
		if gap < read.Gaps[0]-0.5 || gap > read.Gaps[0]+0.5 {
			t.Errorf("the send press leaves %vpx after word %d, and every word sits %vpx from the next as Send does from and", gap, at+2, read.Gaps[0])
		}
	}
}
