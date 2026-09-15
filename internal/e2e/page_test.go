package e2e_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/agora/surface"
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

func read(t *testing.T, document string) (*surfaced, *reader) {
	t.Helper()

	held := standing(t)
	address := held.serve(document)
	return held, opening(t, address)
}

func TestChoosingSendsNothing(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.evaluate(chose("c-1", 1), nil)
	page.evaluate(chose("c-2", 2), nil)
	page.evaluate(chose("c-2", 2), nil)
	page.evaluate(chose("c-2", 0), nil)

	var ticked int
	page.evaluate(`document.querySelectorAll('.choice:checked').length`, &ticked)
	if ticked != 2 {
		t.Errorf("selecting, deselecting and reselecting left %d options ticked, want 2", ticked)
	}

	if delivered, arrived := held.awaiting(shortWindow); arrived {
		t.Errorf("choosing delivered an answer: %+v", delivered)
	}

	page.click(".send")
	page.waitFor(`document.querySelector('.state').textContent.length > 0`)

	sent, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("sending delivered nothing at all")
	}
	if len(sent.Decisions) != 0 {
		t.Errorf("the send that waits for a reply carried a selection nobody committed: %+v", sent.Decisions)
	}

	page.click(".finish")
	page.waitFor(`!document.querySelector('.discarding').hidden`)
	page.click(".discard")
	page.waitFor(`document.querySelector('.send').disabled === true`)

	finished, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("finishing delivered nothing at all")
	}
	if len(finished.Decisions) != 0 {
		t.Errorf("the send that ends carried a selection nobody committed: %+v", finished.Decisions)
	}
}

func TestMultipleChoiceCarriesEveryTick(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-2", 2), nil)
	page.evaluate(chose("c-2", 0), nil)
	page.evaluate(chose("c-2", 1), nil)
	page.click(within("c-2", ".commit"))
	page.click(".send")
	page.waitFor(`document.querySelector('.state').textContent.length > 0`)

	delivered, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("a committed decision was sent and nothing was delivered")
	}
	if len(delivered.Decisions) != 1 {
		t.Fatalf("sending carried %d decisions, want the one that was committed: %+v", len(delivered.Decisions), delivered)
	}

	want := []string{"the write", "the metric", "the log line"}
	if diff := cmp.Diff(want, delivered.Decisions[0].Chosen); diff != "" {
		t.Errorf("what three ticks carried differs (-want +got):\n%s", diff)
	}
}

func TestCommittingAgainReplacesTheEarlierAnswer(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.click(within("c-1", ".commit"))
	page.click(within("c-1", ".reopen"))
	page.evaluate(chose("c-1", 1), nil)
	page.click(within("c-1", ".commit"))
	page.click(".send")
	page.waitFor(`document.querySelector('.state').textContent.length > 0`)

	delivered, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("a committed decision was sent and nothing was delivered")
	}
	if len(delivered.Decisions) != 1 {
		t.Fatalf("committing twice delivered %d answers for one decision: %+v", len(delivered.Decisions), delivered)
	}

	want := []string{"agree c-1"}
	if diff := cmp.Diff(want, delivered.Decisions[0].Commands); diff != "" {
		t.Errorf("the answer delivered is not the later one (-want +got):\n%s", diff)
	}
}

func TestACommittedAnswerIsDistinguishable(t *testing.T) {
	_, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.click(within("c-1", ".commit"))
	page.evaluate(chose("c-2", 1), nil)

	var committed, touched string
	page.evaluate(stateOf("c-1"), &committed)
	page.evaluate(stateOf("c-2"), &touched)
	if committed == touched {
		t.Errorf("a committed decision and one merely selected in both render as %q", committed)
	}

	var shown, unshown bool
	page.evaluate(`!document.querySelector('.decision[data-decision="c-1"] .committed').hidden`, &shown)
	page.evaluate(`!document.querySelector('.decision[data-decision="c-2"] .committed').hidden`, &unshown)
	if !shown {
		t.Error("a committed decision does not show what it committed")
	}
	if unshown {
		t.Error("a decision merely selected in shows a committed answer")
	}

	var said string
	page.evaluate(`document.querySelector('.decision[data-decision="c-1"] .committed').textContent`, &said)
	if !strings.Contains(said, "take the lock first") {
		t.Errorf("a committed decision shows %q, and it does not say what was committed", said)
	}
}

func TestACommittedAnswerCanBeTakenBack(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.click(within("c-1", ".commit"))
	page.click(within("c-1", ".reopen"))

	var state string
	page.evaluate(stateOf("c-1"), &state)
	if state != "open" {
		t.Errorf("a reopened decision is in the state %q", state)
	}

	page.click(".send")
	page.waitFor(`document.querySelector('.state').textContent.length > 0`)

	delivered, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("sending delivered nothing at all")
	}
	if len(delivered.Decisions) != 0 {
		t.Errorf("a decision taken back before sending was delivered: %+v", delivered.Decisions)
	}
}

func TestSendingWithoutEndingKeepsThePageOpen(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.click(within("c-1", ".commit"))
	page.write(".message-field", "one more look please")
	page.click(".send")
	page.waitFor(`document.querySelector('.state').textContent.length > 0`)

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

	var closed bool
	page.evaluate(`document.querySelector('.send').disabled`, &closed)
	if closed {
		t.Error("sending without ending closed the page")
	}

	held.serveAgain(oneDecisionReply)
	page.waitFor(`!!document.querySelector('.decision[data-decision="c-3"]')`)

	var title string
	page.evaluate(`document.querySelector('.heading h1').textContent`, &title)
	if title != "one thing left" {
		t.Errorf("the page did not take up the reply: it heads with %q", title)
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

func TestAnAnswerCommittedButNotSentIsNotCarriedIntoNewContent(t *testing.T) {
	// scenario: A re-render reaches an open page
	held, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.click(within("c-1", ".commit"))

	held.serveAgain(oneDecisionReply)
	page.waitFor(`!!document.querySelector('.decision[data-decision="c-3"]')`)

	page.click(".send")
	page.waitFor(`document.querySelector('.state').textContent.length > 0`)

	delivered, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("sending after a re-render delivered nothing")
	}
	for _, decided := range delivered.Decisions {
		if decided.ID == "c-1" {
			t.Errorf("an answer committed before the re-render was carried into content that does not hold it: %+v", decided)
		}
	}
}

func TestEndingClosesTheExchange(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.click(within("c-1", ".commit"))
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
	page.click(within("c-1", ".commit"))

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
		t.Fatal("a committed decision was sent and nothing was delivered")
	}
	if diff := cmp.Diff([]string{"patch c-1"}, delivered.Decisions[0].Commands); diff != "" {
		t.Errorf("what the caller was handed back differs (-want +got):\n%s\nthe command is the caller's, and it still has to arrive", diff)
	}
}

func TestACommittedAnswerIsShownInTheWordsThePersonRead(t *testing.T) {
	_, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-2", 0), nil)
	page.evaluate(chose("c-2", 2), nil)
	page.click(within("c-2", ".commit"))

	var shown string
	page.evaluate(`document.querySelector('.decision[data-decision="c-2"] .committed').textContent`, &shown)

	for _, label := range []string{"the write", "the log line"} {
		if !strings.Contains(shown, label) {
			t.Errorf("what was committed is shown as %q, and it does not name the option %q the person read", shown, label)
		}
	}
	if strings.Contains(shown, "the metric") {
		t.Errorf("what was committed is shown as %q, and it names an option nobody chose", shown)
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

	page.evaluate(chose("c-2", 2), nil)
	page.click(within("c-2", ".commit"))
	page.click(".send")
	page.waitFor(`document.querySelector('.state').textContent.length > 0`)

	printed, failures, code := held.poll()
	if code != 0 {
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

func TestMarkingIsNotChoosing(t *testing.T) {
	held, page := read(t, twoDecisionDocument)

	var ticked int
	page.evaluate(`document.querySelectorAll('.choice:checked').length`, &ticked)
	if ticked != 0 {
		t.Errorf("%d options arrive chosen, and a mark ranks an option rather than answering for the person", ticked)
	}

	page.click(".send")
	page.waitFor(`document.querySelector('.state').textContent.length > 0`)

	delivered, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("sending delivered nothing at all")
	}
	for _, decided := range delivered.Decisions {
		t.Errorf("a decision nobody touched was delivered: %+v", decided)
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
	page.click(within("c-1", ".commit"))
	page.click(".finish")
	page.waitFor(`document.querySelector('.send').disabled === true`)

	printed, failures, code := held.poll()
	if code != 0 {
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

	page.click(`.entry a.entry-name`)
	page.waitFor(`location.hash === '#decision-c-1'`)

	var reached string
	page.evaluate(`document.querySelector(location.hash).dataset.decision`, &reached)
	if reached != "c-1" {
		t.Errorf("following the entry arrived at %q, and the decision covering it is %q", reached, "c-1")
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
