package e2e_test

import (
	"encoding/json"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/agora/surface"
	"github.com/socrates-tools/socrates/internal/exit"
)

const shortWindow = 500 * time.Millisecond

const twoDecisionDocument = `{
  "title": "two things to settle",
  "preamble": "what follows is what was seen",
  "decisions": [
    {
      "id": "c-1",
      "question": "should the cache be read before the lock is taken",
      "meaning": "two callers arriving together read the same stale entry",
      "mode": "single-choice",
      "options": [
        {"label": "take the lock first", "meaning": "correctness first", "command": "patch c-1"},
        {"label": "leave it", "meaning": "acceptable here", "command": "agree c-1"}
      ],
      "evidence": {
        "language": "go",
        "code": "entry, held := cache.Read(key)\nif held {\n\treturn entry\n}\n",
        "patch": "--- a/internal/cache/write.go\n+++ b/internal/cache/write.go\n@@ -10,3 +10,4 @@ func Write(key string, entry Entry) {\n func Write(key string, entry Entry) {\n-\tstore[key] = entry\n+\tmu.Lock()\n+\tstore[key] = entry\n }\n"
      }
    },
    {
      "id": "c-2",
      "question": "which of these should the write path keep",
      "meaning": "the write path does three things and one of them is what the caller asked for",
      "mode": "multiple-choice",
      "options": [
        {"label": "the write", "meaning": "what the caller asked for", "command": "keep c-2 write"},
        {"label": "the metric", "meaning": "counted nowhere else", "command": "keep c-2 metric"},
        {"label": "the log line", "meaning": "one line per write", "command": "keep c-2 log", "recommended": true}
      ],
      "fragment": {"markup": "<figure><svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 120 40\" width=\"120\" height=\"40\" role=\"img\" aria-label=\"the write path\"><rect x=\"1\" y=\"1\" width=\"40\" height=\"20\" fill=\"none\" stroke=\"#8a6d3b\"/><text x=\"6\" y=\"34\" font-size=\"11\">write</text></svg><figcaption>the write path</figcaption></figure>"}
    }
  ]
}`

const oneDecisionReply = `{
  "title": "one thing left",
  "preamble": "the earlier answers were taken",
  "decisions": [
    {
      "id": "c-3",
      "question": "is the lock the right shape",
      "meaning": "one mutex for the whole map",
      "mode": "single-choice",
      "options": [
        {"label": "yes", "meaning": "good enough", "command": "agree c-3"},
        {"label": "shard it", "meaning": "contention matters", "command": "patch c-3"}
      ]
    }
  ]
}`

const reusingIdentifierReply = `{
  "title": "an unrelated question reusing the identifier",
  "decisions": [
    {
      "id": "c-1",
      "question": "is the timeout long enough for a cold start",
      "mode": "single-choice",
      "options": [
        {"label": "yes", "meaning": "measured under load", "command": "agree c-1"},
        {"label": "raise it", "meaning": "cold starts blow past it", "command": "patch c-1"}
      ]
    }
  ]
}`

func read(t *testing.T, document string) (*surfaced, *reader) {
	t.Helper()

	held := standing(t)
	address := held.serve(document)
	return held, opening(t, address)
}

func TestAnAnswerLeftInAControlIsSentWithoutAnyFurtherPress(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.evaluate(chose("c-1", 1), nil)
	page.evaluate(chose("c-2", 2), nil)
	page.evaluate(chose("c-2", 2), nil)
	page.evaluate(chose("c-2", 0), nil)

	var ticked int
	page.evaluate(`document.querySelectorAll('.choice:checked').length`, &ticked)
	if ticked != 3 {
		t.Errorf("selecting, deselecting and reselecting left %d options ticked, want 3", ticked)
	}

	if delivered, arrived := held.awaiting(shortWindow); arrived {
		t.Errorf("choosing delivered an answer before any press at all: %+v", delivered)
	}

	page.click(".send")
	page.waitFor(`document.querySelector('.state').textContent.length > 0`)

	sent, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("sending delivered nothing at all")
	}
	if len(sent.Decisions) != 2 {
		t.Fatalf("what was left in the controls was sent as %d decisions, want the two that were chosen: %+v", len(sent.Decisions), sent.Decisions)
	}
}

func TestMultipleChoiceCarriesEveryTick(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-2", 0), nil)
	page.evaluate(chose("c-2", 1), nil)
	pressPastTheConfirmation(page, ".send")
	page.waitFor(`document.querySelector('.state').textContent.length > 0`)

	delivered, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("a decision holding a selection was sent and nothing was delivered")
	}
	if len(delivered.Decisions) != 1 {
		t.Fatalf("sending carried %d decisions, want the one holding a selection: %+v", len(delivered.Decisions), delivered)
	}

	want := []string{"the write", "the metric", "the log line"}
	if diff := cmp.Diff(want, delivered.Decisions[0].Chosen); diff != "" {
		t.Errorf("what three ticks carried differs (-want +got):\n%s", diff)
	}
}

func TestAChangedSelectionReplacesTheEarlierOneWithNoSeparateStep(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.evaluate(chose("c-1", 1), nil)
	page.evaluate(chose("c-2", 2), nil)
	pressPastTheConfirmation(page, ".send")
	page.waitFor(`document.querySelector('.state').textContent.length > 0`)

	delivered, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("a changed selection was sent and nothing was delivered")
	}
	if len(delivered.Decisions) != 1 {
		t.Fatalf("changing a selection before sending delivered %d answers for one decision: %+v", len(delivered.Decisions), delivered)
	}

	want := []string{"agree c-1"}
	if diff := cmp.Diff(want, delivered.Decisions[0].Commands); diff != "" {
		t.Errorf("the answer delivered is not the later selection (-want +got):\n%s", diff)
	}
}

func TestADecisionWithNothingSelectedAndNothingTypedIsOmittedFromWhatTheCallerCollects(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-2", 2), nil)
	page.evaluate(chose("c-2", 0), nil)
	page.evaluate(chose("c-2", 0), nil)

	var ticked int
	page.evaluate(`document.querySelectorAll('.choice:checked').length`, &ticked)
	if ticked != 0 {
		t.Fatalf("ticking and unticking a box left %d options ticked, want none", ticked)
	}

	pressPastTheConfirmation(page, ".send")
	page.waitFor(`document.querySelector('.state').textContent.length > 0`)

	delivered, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("sending delivered nothing at all")
	}
	if len(delivered.Decisions) != 0 {
		t.Errorf("a decision with nothing selected and nothing typed was sent anyway: %+v", delivered.Decisions)
	}
}

func TestSendingWithoutEndingKeepsThePageOpen(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.write(".message-field", "one more look please")
	page.click(".send")
	page.waitFor(`!document.querySelector('.receipt').hidden`)

	var said string
	page.evaluate(`document.querySelector('.state').textContent`, &said)
	if !strings.Contains(strings.ToLower(said), "waiting") {
		t.Errorf("sending without ending does not say it is waiting for a reply: %q", said)
	}

	var amendable bool
	page.evaluate(`!document.querySelector('.amend').hidden && !document.querySelector('.amend').disabled`, &amendable)
	if !amendable {
		t.Error("sending without ending offers no way back, and a misclick has nowhere to go")
	}

	delivered, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("sending without ending delivered nothing")
	}
	if delivered.Ended {
		t.Error("sending without ending reports the exchange as over")
	}
	if delivered.Message != "one more look please" {
		t.Errorf("the message a person wrote came back as %q", delivered.Message)
	}

	held.serveAgain(oneDecisionReply)
	page.waitFor(`!!document.querySelector('.decision[data-decision="c-3"]')`)

	var title string
	page.evaluate(`document.querySelector('.heading h1').textContent`, &title)
	if title != "one thing left" {
		t.Errorf("the page did not take up the reply: it heads with %q", title)
	}
}

func TestSendingShowsAReceiptAndPutsEveryControlReadOnly(t *testing.T) {
	_, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.click(".send")
	page.waitFor(`!document.querySelector('.receipt').hidden`)

	var said string
	page.evaluate(`document.querySelector('.receipt').textContent`, &said)
	for _, wanted := range []string{"take the lock first", "the log line"} {
		if !strings.Contains(said, wanted) {
			t.Errorf("the receipt does not say %q, and it is meant to list what went:\n%s", wanted, said)
		}
	}

	var readOnly bool
	page.evaluate(`Array.from(document.querySelectorAll('.choice, .decline, .words-field, .message-field, .send, .finish')).every((held) => held.disabled)`, &readOnly)
	if !readOnly {
		t.Error("sending left a control open to a further change")
	}

	var amendable bool
	page.evaluate(`!document.querySelector('.amend').hidden && !document.querySelector('.amend').disabled`, &amendable)
	if !amendable {
		t.Error("sending without ending offers no way back")
	}
}

func TestAmendingReturnsTheControlsAndASecondSendReplacesTheFirst(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.click(".send")
	page.waitFor(`!document.querySelector('.receipt').hidden`)

	page.click(".amend")
	page.waitFor(`document.querySelector('.receipt').hidden`)

	var readOnly bool
	page.evaluate(`Array.from(document.querySelectorAll('.choice, .words-field, .message-field, .send, .finish')).some((held) => held.disabled)`, &readOnly)
	if readOnly {
		t.Error("amending left a control read only")
	}

	page.evaluate(chose("c-1", 1), nil)
	page.click(".send")
	page.waitFor(`!document.querySelector('.receipt').hidden`)

	delivered, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("sending again after amending delivered nothing")
	}

	var commands []string
	for _, decided := range delivered.Decisions {
		if decided.ID == "c-1" {
			commands = decided.Commands
		}
	}
	if diff := cmp.Diff([]string{"agree c-1"}, commands); diff != "" {
		t.Errorf("what was collected for c-1 is not the amended answer, and a caller collects the last word (-want +got):\n%s", diff)
	}
}

func TestTheWayBackIsWithdrawnOnceTheCallerCollectsAndThePageSaysSo(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.click(".send")
	page.waitFor(`!document.querySelector('.receipt').hidden`)

	var waiting string
	page.evaluate(`document.querySelector('.state').textContent`, &waiting)
	if !strings.Contains(strings.ToLower(waiting), "waiting") {
		t.Errorf("the page does not say it is waiting for a reply before the caller collects: %q", waiting)
	}

	var amendable bool
	page.evaluate(`!document.querySelector('.amend').hidden`, &amendable)
	if !amendable {
		t.Fatal("the way back is not offered before the caller collects")
	}

	if _, failures, code := held.poll(); code != exit.NothingAsked {
		t.Fatalf("agora poll exited %d: %s", code, failures)
	}

	page.waitFor(`document.querySelector('.amend').hidden`)

	var said string
	page.evaluate(`document.querySelector('.state').textContent`, &said)
	if !strings.Contains(strings.ToLower(said), "gone") {
		t.Errorf("the page does not say the answer has gone once the caller collects: %q", said)
	}
	if strings.Contains(strings.ToLower(said), "waiting") {
		t.Errorf("the page still says it is waiting after the caller collected: %q", said)
	}
}

func TestAConnectTimeUncollectedIsNotMistakenForACollection(t *testing.T) {
	_, page := read(t, twoDecisionDocument)
	time.Sleep(150 * time.Millisecond)

	var said string
	page.evaluate(`document.querySelector('.state').textContent`, &said)
	if strings.Contains(strings.ToLower(said), "gone") {
		t.Errorf("a page that never sent anything says %q, and the stream's own connect-time report is not a collection", said)
	}

	var hidden bool
	page.evaluate(`document.querySelector('.receipt').hidden`, &hidden)
	if !hidden {
		t.Error("a page that never sent anything shows a receipt")
	}
}

func TestANewDocumentReplacesTheReceiptToo(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.click(".send")
	page.waitFor(`!document.querySelector('.receipt').hidden`)

	held.serveAgain(oneDecisionReply)
	page.waitFor(`!!document.querySelector('.decision[data-decision="c-3"]')`)

	var hidden bool
	page.evaluate(`document.querySelector('.receipt').hidden`, &hidden)
	if !hidden {
		t.Error("a new document arrived and the receipt from the round before it is still showing")
	}

	var readOnly bool
	page.evaluate(`Array.from(document.querySelectorAll('.choice, .words-field, .message-field, .send, .finish')).some((held) => held.disabled)`, &readOnly)
	if readOnly {
		t.Error("a new document arrived carrying a control the round before it left read only")
	}
}

func TestFinishingShowsAReceiptThatOffersNoWayBack(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.click(".finish")
	page.waitFor(`document.querySelector('.send').disabled === true`)

	if _, arrived := held.awaiting(5 * time.Second); !arrived {
		t.Fatal("finishing delivered nothing")
	}

	var hidden bool
	page.evaluate(`document.querySelector('.receipt').hidden`, &hidden)
	if hidden {
		t.Error("finishing shows no receipt of what went")
	}

	var amendOffered bool
	page.evaluate(`!document.querySelector('.amend').hidden`, &amendOffered)
	if amendOffered {
		t.Error("finishing offers a way back, and nothing further goes")
	}
}

func TestARerenderReachesAnOpenPage(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	var loads int
	page.evaluate(`(window.__loads = (window.__loads || 0) + 1)`, &loads)

	held.serveAgain(oneDecisionReply)
	page.waitFor(`!!document.querySelector('.decision[data-decision="c-3"]')`)

	var again int
	page.evaluate(`window.__loads`, &again)
	if again != loads {
		t.Errorf("the page reloaded to show the reply: %d, then %d", loads, again)
	}

	var gone bool
	page.evaluate(`!document.querySelector('.decision[data-decision="c-1"]')`, &gone)
	if !gone {
		t.Error("the page still holds a decision the new document does not")
	}
}

func TestANewDocumentReplacesEverySelectionEveryTypedWordAndTheMessageFieldAndNothingIsRestoredByIdentifier(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.write(within("c-2", ".words-field"), "a note before the reply arrived")
	page.write(".message-field", "one more look please")

	held.serveAgain(reusingIdentifierReply)
	page.waitFor(`document.querySelector('.heading h1').textContent === 'an unrelated question reusing the identifier'`)

	var ticked int
	page.evaluate(`document.querySelectorAll('.choice:checked').length`, &ticked)
	if ticked != 0 {
		t.Errorf("a new document that reuses the identifier %q arrived with %d options selected, and the earlier answer was about a different question", "c-1", ticked)
	}

	var words, message string
	page.evaluate(`(document.querySelector('.words-field') || {value: ""}).value`, &words)
	page.evaluate(`document.querySelector('.message-field').value`, &message)
	if words != "" {
		t.Errorf("a new document arrived carrying the words %q a person wrote on the one before it", words)
	}
	if message != "" {
		t.Errorf("a new document arrived carrying the message %q a person wrote on the one before it", message)
	}

	pressPastTheConfirmation(page, ".send")
	page.waitFor(`document.querySelector('.state').textContent.length > 0`)

	delivered, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("sending after a new document arrived delivered nothing")
	}
	if len(delivered.Decisions) != 0 {
		t.Errorf("an answer to the question a new document replaced was carried onto the card that reuses its identifier: %+v", delivered.Decisions)
	}
}

func TestEndingClosesTheExchange(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.click(".finish")
	page.waitFor(`document.querySelector('.send').disabled === true`)

	delivered, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("finishing delivered nothing")
	}
	if !delivered.Ended {
		t.Error("finishing does not report the exchange as over")
	}

	var said string
	page.evaluate(`document.querySelector('.state').textContent`, &said)
	if !strings.Contains(strings.ToLower(said), "over") {
		t.Errorf("the page does not say the exchange is over: %q", said)
	}

	var accepting bool
	page.evaluate(`Array.from(document.querySelectorAll('button, .choice, textarea')).some((held) => !held.disabled)`, &accepting)
	if accepting {
		t.Error("a closed exchange still accepts an answer")
	}
}

var everyCommandTheDocumentCarries = []string{
	"patch c-1",
	"agree c-1",
	"keep c-2 write",
	"keep c-2 metric",
	"keep c-2 log",
}

func TestNoCommandReachesThePage(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)

	var written string
	page.evaluate(`document.documentElement.outerHTML`, &written)
	for _, command := range everyCommandTheDocumentCarries {
		if strings.Contains(written, command) {
			t.Errorf("the page carries the command %q, and the person reading it is not at a prompt", command)
		}
	}

	page.click(".send")
	page.waitFor(`document.querySelector('.state').textContent.length > 0`)

	delivered, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("a decision holding a selection was sent and nothing was delivered")
	}
	if diff := cmp.Diff([]string{"patch c-1"}, delivered.Decisions[0].Commands); diff != "" {
		t.Errorf("what the caller was handed back differs (-want +got):\n%s\nthe command is the caller's, and it still has to arrive", diff)
	}
}

func TestWordsAloneOnADecisionAreACompleteAnswerAndComeBackWithTheWordsAndNoCommand(t *testing.T) {
	const inTheirOwnWords = "leave it to whoever asked"

	held, page := read(t, twoDecisionDocument)

	page.write(within("c-1", ".words-field"), inTheirOwnWords)
	page.evaluate(chose("c-2", 2), nil)
	pressPastTheConfirmation(page, ".send")
	page.waitFor(`document.querySelector('.state').textContent.length > 0`)

	delivered, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("words alone on a decision delivered nothing at all")
	}
	if len(delivered.Decisions) != 1 {
		t.Fatalf("words alone on one decision came back as %d answers: %+v", len(delivered.Decisions), delivered.Decisions)
	}

	said := delivered.Decisions[0]
	if said.ID != "c-1" {
		t.Errorf("what came back is attached to %q, and the words were written on %q", said.ID, "c-1")
	}
	if said.Words != inTheirOwnWords {
		t.Errorf("the words came back as %q, want %q verbatim", said.Words, inTheirOwnWords)
	}
	if len(said.Chosen) != 0 || len(said.Commands) != 0 {
		t.Errorf("words alone came back carrying %v chosen and the commands %v, and a complete answer of words alone carries no command", said.Chosen, said.Commands)
	}
}

func TestEnterSubmitsTheForm(t *testing.T) {
	held, page := read(t, oneDecisionDocument)

	page.do(chromedp.Click(within("c-1", ".choice"), chromedp.ByQuery, chromedp.NodeVisible))
	page.do(chromedp.KeyEvent(kb.Enter))
	page.waitFor(`document.querySelector('.state').textContent.length > 0`)

	delivered, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("pressing enter on a control that does not swallow it delivered nothing")
	}
	if delivered.Ended {
		t.Error("pressing enter reached the secondary press rather than the primary that is first in the document")
	}

	if violations := page.violations(); len(violations) != 0 {
		t.Errorf("a real form navigation was attempted and blocked by the policy, and the submit handler is meant to prevent it: %v", violations)
	}
}

func TestTheBarSitsAtTheEndOfThePageRatherThanFixedToTheViewport(t *testing.T) {
	_, page := read(t, twoDecisionDocument)
	page.narrow(360, 480)

	var positioned struct {
		Position  string `json:"position"`
		BelowFold bool   `json:"belowFold"`
		Last      bool   `json:"last"`
	}
	page.evaluate(`(() => {
  const bar = document.querySelector('.bar')
  const style = getComputedStyle(bar)
  return {
    position: style.position,
    belowFold: bar.getBoundingClientRect().top > window.innerHeight,
    last: bar.parentElement === document.querySelector('form.exchange .composer') && bar.nextElementSibling === null && bar.parentElement.nextElementSibling === null,
  }
})()`, &positioned)

	if positioned.Position == "fixed" || positioned.Position == "sticky" {
		t.Errorf("the bar is positioned %q, and it sits in the normal flow of the page rather than fixed to the viewport", positioned.Position)
	}
	if !positioned.BelowFold {
		t.Error("the bar is visible without scrolling, and a page long enough to need scrolling was not given one")
	}
	if !positioned.Last {
		t.Error("the bar is not the last thing on the page")
	}
}

func TestAPageWhoseEveryDecisionAsksNothingShowsTheSameBarWithNoCounter(t *testing.T) {
	normal := standing(t)
	normalPage := opening(t, normal.serve(twoDecisionDocument))

	var normalBar string
	normalPage.evaluate(`document.querySelector('.bar').textContent.replace(/\s+/g, ' ').trim()`, &normalBar)

	held := standing(t)
	page := opening(t, held.serve(askingNothingDocument))

	var barText string
	page.evaluate(`document.querySelector('.bar').textContent.replace(/\s+/g, ' ').trim()`, &barText)

	if barText != normalBar {
		t.Errorf("a page whose every decision asks nothing shows the bar as %q, and a page with something to choose shows %q", barText, normalBar)
	}
	if regexp.MustCompile(`\d`).MatchString(barText) {
		t.Errorf("the bar carries a digit, and there is no counter and no meter: %q", barText)
	}
}

const optionsRead = `Array.from(document.querySelectorAll('.decision[data-decision="c-2"] .option .option-label')).map((held) => held.textContent.trim())`

func TestTheOptionTheCallerWouldTakeIsReadFirst(t *testing.T) {
	_, page := read(t, twoDecisionDocument)

	var offered []string
	page.evaluate(optionsRead, &offered)

	want := []string{"the log line", "the write", "the metric"}
	if diff := cmp.Diff(want, offered); diff != "" {
		t.Errorf("the options a person meets differ (-want +got):\n%s", diff)
	}

	var marked []string
	page.evaluate(`Array.from(document.querySelectorAll('.option[data-recommended="yes"] .option-label')).map((held) => held.textContent.trim())`, &marked)
	if diff := cmp.Diff([]string{"the log line"}, marked); diff != "" {
		t.Errorf("what the page marks as the option the caller would take differs (-want +got):\n%s", diff)
	}
}

func TestRearrangingTheOptionsDoesNotRearrangeTheAnswer(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	pressPastTheConfirmation(page, ".send")
	page.waitFor(`document.querySelector('.state').textContent.length > 0`)

	printed, failures, code := held.poll()
	if code != exit.NothingAsked {
		t.Fatalf("agora poll exited %d: %s", code, failures)
	}
	if !strings.Contains(printed, "keep c-2 log") {
		t.Errorf("choosing the marked option handed the caller back:\n%s\nand the option carries the command %q", printed, "keep c-2 log")
	}
	for _, others := range []string{"keep c-2 write", "keep c-2 metric"} {
		if strings.Contains(printed, others) {
			t.Errorf("choosing the marked option handed back %q, which nobody chose:\n%s", others, printed)
		}
	}
	if !strings.Contains(printed, "the log line") {
		t.Errorf("what came back does not name the option in the words the person read:\n%s", printed)
	}
}

func TestTheRecommendationArrivesChosenAndADecisionMarkedNothingDoesNot(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	var ticked []string
	page.evaluate(`Array.from(document.querySelectorAll('.choice:checked')).map((held) => held.closest('.decision').dataset.decision + ':' + held.value)`, &ticked)
	if diff := cmp.Diff([]string{"c-2:2"}, ticked); diff != "" {
		t.Errorf("what arrives chosen before anybody touches the page differs (-want +got):\n%s", diff)
	}

	pressPastTheConfirmation(page, ".send")
	page.waitFor(`document.querySelector('.state').textContent.length > 0`)

	delivered, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("sending delivered nothing at all")
	}
	if len(delivered.Decisions) != 1 || delivered.Decisions[0].ID != "c-2" {
		t.Fatalf("sending without touching anything carried %+v, want only the recommendation the document marked", delivered.Decisions)
	}
	if diff := cmp.Diff([]string{"the log line"}, delivered.Decisions[0].Chosen); diff != "" {
		t.Errorf("what the recommendation carried differs (-want +got):\n%s", diff)
	}
}

const variedSeverityDocument = `{
  "title": "three things to settle",
  "decisions": [
    {
      "id": "d-1",
      "question": "an error needing a person",
      "severity": "error",
      "audience": "a person",
      "mode": "single-choice",
      "options": [
        {"label": "fix it", "command": "patch d-1"},
        {"label": "leave it", "command": "agree d-1", "recommended": true}
      ]
    },
    {
      "id": "d-2",
      "question": "a warning the author patches",
      "severity": "warning",
      "audience": "the author patches",
      "mode": "single-choice",
      "options": [
        {"label": "fix it", "command": "patch d-2", "recommended": true},
        {"label": "leave it", "command": "agree d-2"}
      ]
    },
    {
      "id": "d-3",
      "question": "an informational note",
      "severity": "info",
      "audience": "the author patches",
      "mode": "single-choice",
      "options": [
        {"label": "fix it", "command": "patch d-3"},
        {"label": "leave it", "command": "agree d-3", "recommended": true}
      ]
    }
  ]
}`

func TestARecommendationArrivesChosenWhateverItsSeverityOrAudience(t *testing.T) {
	_, page := read(t, variedSeverityDocument)

	var ticked []string
	page.evaluate(`Array.from(document.querySelectorAll('.choice:checked')).map((held) => held.closest('.decision').dataset.decision)`, &ticked)

	want := []string{"d-1", "d-2", "d-3"}
	if diff := cmp.Diff(want, ticked); diff != "" {
		t.Errorf("which decisions arrive with their recommendation chosen differs (-want +got):\n%s", diff)
	}
}

const noSeverityAnywhereDocument = `{
  "title": "an argument with no severity",
  "source": [
    {"id": "d-3", "name": "the third thing settled"},
    {"id": "d-1", "name": "the first thing settled"},
    {"id": "d-2", "name": "the second thing settled"}
  ],
  "decisions": [
    {"id": "d-3", "question": "the third question", "mode": "single-choice", "options": [{"label": "yes", "command": "agree d-3"}]},
    {"id": "d-1", "question": "the first question", "mode": "single-choice", "options": [{"label": "yes", "command": "agree d-1"}]},
    {"id": "d-2", "question": "the second question", "mode": "single-choice", "options": [{"label": "yes", "command": "agree d-2"}]}
  ]
}`

func TestADocumentStatingNoSeverityAnywhereRendersCardsInTheOrderItListedThem(t *testing.T) {
	_, page := read(t, noSeverityAnywhereDocument)

	var cards, listed []string
	page.evaluate(`Array.from(document.querySelectorAll('.decision')).map((held) => held.dataset.decision)`, &cards)
	page.evaluate(`Array.from(document.querySelectorAll('.entry .entry-name')).map((held) => held.textContent.trim())`, &listed)

	if diff := cmp.Diff([]string{"d-3", "d-1", "d-2"}, cards); diff != "" {
		t.Errorf("a document stating no severity anywhere did not render its decisions in the order it listed them (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]string{"the third thing settled", "the first thing settled", "the second thing settled"}, listed); diff != "" {
		t.Errorf("a document stating no severity anywhere did not render its source entries in the order it listed them (-want +got):\n%s", diff)
	}
}

func TestChoosingNoneOfTheseClearsTheSelectionAndMovesToTheField(t *testing.T) {
	_, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.evaluate(declined("c-1"), nil)

	var ticked int
	page.evaluate(`document.querySelectorAll('.decision[data-decision="c-1"] .choice:checked').length`, &ticked)
	if ticked != 0 {
		t.Errorf("choosing none of the options on c-1 left %d options ticked, want none", ticked)
	}

	var focused bool
	page.evaluate(`document.activeElement === document.querySelector('.decision[data-decision="c-1"] .words-field')`, &focused)
	if !focused {
		t.Error("choosing none of the options does not move the cursor to the words field below it")
	}
}

func TestChoosingNoneOfTheseAndTypingSendsTheWordsAndNoCommand(t *testing.T) {
	const inTheirOwnWords = "I would rather write my own answer here"

	held, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.evaluate(declined("c-1"), nil)
	page.write(within("c-1", ".words-field"), inTheirOwnWords)
	page.evaluate(chose("c-2", 2), nil)

	pressPastTheConfirmation(page, ".send")
	page.waitFor(`document.querySelector('.state').textContent.length > 0`)

	delivered, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("sending after declining every option and writing an answer delivered nothing")
	}
	if len(delivered.Decisions) != 1 || delivered.Decisions[0].ID != "c-1" {
		t.Fatalf("what was sent is %+v, want only the decision that was declined and answered in words", delivered.Decisions)
	}

	said := delivered.Decisions[0]
	if said.Words != inTheirOwnWords {
		t.Errorf("the words came back as %q, want %q verbatim", said.Words, inTheirOwnWords)
	}
	if len(said.Chosen) != 0 || len(said.Commands) != 0 {
		t.Errorf("declining every option and writing an answer came back carrying %v chosen and the commands %v", said.Chosen, said.Commands)
	}
}

func TestChoosingNoneOfTheseAndTypingNothingLeavesTheDecisionBlank(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.evaluate(declined("c-1"), nil)

	page.click(".send")
	page.waitFor(confirmationIsOpen)

	if delivered, arrived := held.awaiting(shortWindow); arrived {
		t.Fatalf("sending over a decision left blank by declining every option delivered an answer: %+v", delivered)
	}

	var said string
	page.evaluate(whatTheConfirmationSays, &said)
	if !strings.Contains(said, "c-1") {
		t.Errorf("the confirmation does not name the decision declining every option left blank: %q", said)
	}
}

func TestAMultipleChoiceDecisionOffersNoWayToDeclineAndUntickingEveryBoxLeavesItBlank(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	var none int
	page.evaluate(`document.querySelectorAll('.decision[data-decision="c-2"] .none').length`, &none)
	if none != 0 {
		t.Error("a multiple-choice decision renders a control for taking none of the options")
	}

	page.evaluate(chose("c-1", 0), nil)
	page.evaluate(chose("c-2", 2), nil)

	page.click(".send")
	page.waitFor(confirmationIsOpen)

	if delivered, arrived := held.awaiting(shortWindow); arrived {
		t.Fatalf("sending over a multiple-choice decision left blank by unticking every box delivered an answer: %+v", delivered)
	}

	var said string
	page.evaluate(whatTheConfirmationSays, &said)
	if !strings.Contains(said, "c-2") {
		t.Errorf("the confirmation does not name the multiple-choice decision unticking every box left blank: %q", said)
	}
}

const askingNothingDocument = `{
  "title": "what was put right",
  "preamble": "what follows is what was seen and what was done about it",
  "source": [
    {"id": "c-1", "name": "the lock deleted itself", "severity": "error", "audience": "the author patches", "location": "internal/cache/read.go:41"},
    {"id": "c-2", "name": "the write counts twice", "severity": "warning", "audience": "the author patches", "location": "internal/cache/write.go:12"}
  ],
  "decisions": [
    {
      "id": "c-1",
      "question": "the lock deleted itself while another request held it",
      "meaning": "the deletion is gone, and the entry is dropped by whichever reader took it last",
      "evidence": {
        "language": "go",
        "patch": "--- a/internal/cache/read.go\n+++ b/internal/cache/read.go\n@@ -39,3 +39,2 @@ func Read(key string) Entry {\n func Read(key string) Entry {\n-\tdelete(locks, key)\n \treturn store[key]\n }\n"
      }
    },
    {
      "id": "c-2",
      "question": "the write counted itself twice",
      "meaning": "the metric is now incremented once, at the end of the write"
    }
  ]
}`

func TestAnAnnotationOnADecisionThatAskedNothingComesBack(t *testing.T) {
	const inTheirOwnWords = "the second one reads better than the finding did"

	held := standing(t)
	page := opening(t, held.serve(askingNothingDocument))

	page.write(within("c-1", ".words-field"), inTheirOwnWords)
	pressPastTheConfirmation(page, ".finish")
	page.waitFor(`document.querySelector('.send').disabled === true`)

	printed, failures, code := held.poll()
	if code != exit.NothingAsked {
		t.Fatalf("agora poll exited %d: %s", code, failures)
	}

	var answered surface.Answered
	if err := json.Unmarshal([]byte(printed), &answered); err != nil {
		t.Fatalf("reading back what agora poll printed: %v\n%s", err, printed)
	}
	if len(answered.Decisions) != 1 {
		t.Fatalf("writing on one decision that asked nothing came back as %d answers:\n%s", len(answered.Decisions), printed)
	}

	said := answered.Decisions[0]
	if said.ID != "c-1" {
		t.Errorf("what came back is attached to %q, and it was written on %q", said.ID, "c-1")
	}
	if said.Words != inTheirOwnWords {
		t.Errorf("what the person wrote came back as %q, want %q verbatim", said.Words, inTheirOwnWords)
	}
	if len(said.Chosen) != 0 {
		t.Errorf("a decision that offered nothing came back with %v chosen", said.Chosen)
	}
	if len(said.Commands) != 0 {
		t.Errorf("a decision that offered nothing came back carrying the commands %v", said.Commands)
	}
}

func TestAPageOfCardsThatAskNothingOffersNoControls(t *testing.T) {
	held := standing(t)
	page := opening(t, held.serve(askingNothingDocument))

	var cards []string
	page.evaluate(`Array.from(document.querySelectorAll('.decision')).map((held) => held.dataset.decision)`, &cards)
	if diff := cmp.Diff([]string{"c-1", "c-2"}, cards); diff != "" {
		t.Errorf("what a document whose every decision asks nothing renders differs (-want +got):\n%s", diff)
	}

	var offered []string
	page.evaluate(`Array.from(document.querySelectorAll('.decision .choice, .decision .options, .decision legend')).map((held) => held.tagName + '.' + String(held.className))`, &offered)
	if len(offered) != 0 {
		t.Errorf("a page whose every card asks nothing offers %v, and there is nothing on any of them to choose among", offered)
	}

	var carried bool
	page.evaluate(`Array.from(document.querySelectorAll('.decision')).every((held) => !!held.querySelector('.words-field'))`, &carried)
	if !carried {
		t.Error("a card that asks nothing carries no field for a person to write in")
	}
}

const optionsWithNoCommandDocument = `{
  "title": "one design question",
  "decisions": [
    {
      "id": "d-1",
      "question": "what should the surface be called",
      "meaning": "the vocabulary settles what the term means later",
      "mode": "single-choice",
      "options": [
        {"label": "the surface", "meaning": "what the loopback server is"},
        {"label": "the page", "meaning": "what a person reads"}
      ]
    }
  ]
}`

func TestADecisionWhoseOptionsAllCarryNoCommandIsServedAndRendered(t *testing.T) {
	held := standing(t)
	page := opening(t, held.serve(optionsWithNoCommandDocument))

	var offered []string
	page.evaluate(`Array.from(document.querySelectorAll('.decision[data-decision="d-1"] .choice')).map((held) => held.value)`, &offered)
	if diff := cmp.Diff([]string{"0", "1"}, offered); diff != "" {
		t.Errorf("options that carry no command do not render as controls to choose between (-want +got):\n%s", diff)
	}

	var none int
	page.evaluate(`document.querySelectorAll('.decision[data-decision="d-1"] .none').length`, &none)
	if none != 1 {
		t.Errorf("a decision offering options that carry no command offers %d controls to take none of them, want 1", none)
	}

	page.evaluate(chose("d-1", 1), nil)
	page.click(".send")
	page.waitFor(`document.querySelector('.send').disabled === true`)

	printed, failures, code := held.poll()
	if code != exit.NothingAsked {
		t.Fatalf("agora poll exited %d: %s", code, failures)
	}

	var answered surface.Answered
	if err := json.Unmarshal([]byte(printed), &answered); err != nil {
		t.Fatalf("reading back what agora poll printed: %v\n%s", err, printed)
	}
	if len(answered.Decisions) != 1 {
		t.Fatalf("choosing the one decision's option came back as %d answers:\n%s", len(answered.Decisions), printed)
	}

	decided := answered.Decisions[0]
	if diff := cmp.Diff([]string{"the page"}, decided.Chosen); diff != "" {
		t.Errorf("what was chosen came back differently (-want +got):\n%s", diff)
	}
	if len(decided.Commands) != 0 {
		t.Errorf("an option that carries no command came back carrying the commands %v", decided.Commands)
	}
}

const namedAndUnnamed = `{
  "title": "two things to settle",
  "source": [
    {"id": "c-1", "name": "Finding 1", "severity": "error", "audience": "needs a person", "location": "internal/cache/read.go:41"},
    {"id": "c-2", "severity": "warning", "audience": "can be patched", "location": "internal/cache/write.go:12"}
  ],
  "decisions": [
    {
      "id": "c-1",
      "name": "whatever the session felt like",
      "question": "should the cache be read before the lock is taken",
      "mode": "single-choice",
      "options": [{"label": "take the lock first", "meaning": "correctness first", "command": "patch c-1"}]
    },
    {
      "id": "c-2",
      "question": "which of these should the write path keep",
      "mode": "multiple-choice",
      "options": [{"label": "the write", "meaning": "what the caller asked for", "command": "keep c-2 write"}]
    }
  ]
}`

const glancedDocument = `{
  "title": "one thing to settle",
  "preamble": "what follows is what was seen",
  "source": [
    {"id": "c-1", "name": "the read is unlocked", "severity": "error", "audience": "needs a person", "location": "internal/cache/read.go:41"},
    {"id": "c-2", "name": "the write counts twice", "severity": "warning", "audience": "the author patches", "location": "internal/cache/write.go:12"},
    {"id": "c-3", "name": "the log line says nothing", "severity": "info", "audience": "the author patches", "location": "internal/cache/read.go:80"}
  ],
  "decisions": [
    {
      "id": "c-1",
      "question": "should the cache be read before the lock is taken",
      "meaning": "two callers arriving together read the same stale entry",
      "mode": "single-choice",
      "options": [
        {"label": "take the lock first", "meaning": "correctness first", "command": "patch c-1"},
        {"label": "leave it", "meaning": "acceptable here", "command": "agree c-1"}
      ]
    }
  ]
}`

const nothingToAnswerDocument = `{
  "title": "nothing to settle",
  "preamble": "what follows is what was seen",
  "source": [
    {"id": "c-1", "name": "the read is unlocked", "severity": "warning", "audience": "the author patches", "location": "internal/cache/read.go:41"},
    {"id": "c-2", "name": "the write counts twice", "severity": "info", "audience": "the author patches", "location": "internal/cache/write.go:12"}
  ]
}`

const framing = `(() => {
  const surface = getComputedStyle(document.querySelector('main.surface'));
  const heading = getComputedStyle(document.querySelector('.heading h1'));
  const entry = getComputedStyle(document.querySelector('.entry'));
  return {
    paper: getComputedStyle(document.body).backgroundColor,
    ink: getComputedStyle(document.body).color,
    family: getComputedStyle(document.body).fontFamily,
    width: surface.maxWidth,
    padding: surface.padding,
    title: heading.fontSize,
    entryFace: entry.backgroundColor,
    entryRule: entry.borderTopColor,
    entryText: entry.color,
  };
})()`

type framed struct {
	Paper     string `json:"paper"`
	Ink       string `json:"ink"`
	Family    string `json:"family"`
	Width     string `json:"width"`
	Padding   string `json:"padding"`
	Title     string `json:"title"`
	EntryFace string `json:"entryFace"`
	EntryRule string `json:"entryRule"`
	EntryText string `json:"entryText"`
}

func TestThePageSaysNobodyIsWaiting(t *testing.T) {
	held := standing(t)
	page := openingNothingToAnswer(t, held.serve(nothingToAnswerDocument))

	var said string
	page.evaluate(`document.querySelector('.nothing-said').textContent`, &said)
	for _, wanted := range []string{"nothing", "waiting"} {
		if !strings.Contains(strings.ToLower(said), wanted) {
			t.Errorf("a page with nothing to answer says %q, and it says there is nothing to answer and that nobody is waiting", said)
		}
	}

	var offered []string
	page.evaluate(`Array.from(document.querySelectorAll('input, textarea, button, select, form')).map((held) => held.tagName + '.' + String(held.className))`, &offered)
	if len(offered) != 0 {
		t.Errorf("a page with nothing to answer offers %v, and there is nothing on it to answer and nothing to send", offered)
	}
}

func TestTheFrameIsTheSameFrame(t *testing.T) {
	held := standing(t)
	asking := opening(t, held.serve(glancedDocument))

	var whenAsked framed
	asking.evaluate(framing, &whenAsked)

	quiet := standing(t)
	reading := openingNothingToAnswer(t, quiet.serve(nothingToAnswerDocument))

	var whenRead framed
	reading.evaluate(framing, &whenRead)

	if diff := cmp.Diff(whenAsked, whenRead); diff != "" {
		t.Errorf("a page that is read and a page that is answered are framed differently (-answered +read):\n%s", diff)
	}

	for _, preference := range []string{"light", "dark"} {
		reading.emulate(preference)

		var weighed weighed
		reading.evaluate(measuringTheContrast, &weighed)
		if weighed.Measured < 8 {
			t.Fatalf("under a %s preference only %d texts were measured on a page with nothing to answer", preference, weighed.Measured)
		}
		for _, failing := range weighed.Failing {
			t.Errorf("under a %s preference %s, and the floor the page is held to is %v", preference, failing, contrastFloor)
		}
	}
}

func TestAnEntryTakesThePersonToTheDecision(t *testing.T) {
	// scenario: An entry leads to the decision that covers it
	held := standing(t)
	page := opening(t, held.serve(glancedDocument))

	var listed []string
	page.evaluate(`Array.from(document.querySelectorAll('.entry .entry-name')).map((held) => held.textContent.trim())`, &listed)

	want := []string{"the read is unlocked", "the write counts twice", "the log line says nothing"}
	if diff := cmp.Diff(want, listed); diff != "" {
		t.Errorf("what a person sees at a glance differs (-want +got):\n%s", diff)
	}

	page.click(`.entry a.entry-link`)
	page.waitFor(`!!document.activeElement.closest('.decision')`)

	var reached landed
	page.evaluate(whereTheCursorLanded, &reached)
	if reached.Decision != "c-1" {
		t.Errorf("following the entry arrived at %q, and the decision covering it is %q", reached.Decision, "c-1")
	}
	if reached.Control != "choice" {
		t.Errorf("following the entry left the cursor on a %q, and it lands on the card's first control", reached.Control)
	}
}

func TestADecisionIsShownUnderTheNameItsSourceGivesIt(t *testing.T) {
	_, page := read(t, namedAndUnnamed)

	var named, unnamed string
	page.evaluate(shownIn("c-1", ".name"), &named)
	page.evaluate(shownIn("c-2", ".name"), &unnamed)

	if strings.TrimSpace(named) != "Finding 1" {
		t.Errorf("a decision whose source names it is shown as %q, and the source named it %q", named, "Finding 1")
	}
	if strings.TrimSpace(unnamed) != "c-2" {
		t.Errorf("a decision whose source names it nothing is shown as %q, and it has only its identifier to be shown under", unnamed)
	}
}
