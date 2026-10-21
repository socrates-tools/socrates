package e2e_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/agora/surface"
	"github.com/socrates-tools/socrates/internal/exit"
)

const designRoundDocument = `{
  "title": "two questions before the spec is written",
  "preamble": "reading the repository already ruled out one option",
  "decisions": [
    {
      "id": "q1",
      "name": "Question 1",
      "question": "what should the loopback server be called in the glossary",
      "meaning": "the term settles what CONTEXT.md writes down",
      "mode": "single-choice",
      "options": [
        {"label": "the surface", "meaning": "matches what agora's own source already calls it", "recommended": true},
        {"label": "the host", "meaning": "a new word for something already named"}
      ]
    },
    {
      "id": "q2",
      "name": "Question 2",
      "question": "should a design round carry a cap the way a review round does",
      "meaning": "review caps a round because it spawns a reviewer against a frozen snapshot, and a design spawns nothing",
      "mode": "single-choice",
      "options": [
        {"label": "no cap", "meaning": "the frontier decides when the interview stops"},
        {"label": "cap it at two", "meaning": "matches review"}
      ]
    }
  ]
}`

const inTheirOwnWordsAboutTheCap = "no cap, because nothing here spawns a process for a cap to protect"

const designRoundTwoDocument = `{
  "title": "one question left before the spec is written",
  "preamble": "Question 1: the surface, matching what agora's own source already calls it. Question 2: no round cap, because nothing here spawns a process for one to protect.",
  "settled": [
    {"id": "q1", "name": "Question 1", "question": "what should the loopback server be called in the glossary", "chosen": ["the surface"]},
    {"id": "q2", "name": "Question 2", "question": "should a design round carry a cap the way a review round does", "words": "no cap, because nothing here spawns a process for a cap to protect"}
  ],
  "decisions": [
    {
      "id": "q3",
      "name": "Question 3",
      "question": "should CONTEXT.md gain an entry for the surface today or wait until the spec is written",
      "meaning": "the glossary is maintained as terms settle",
      "mode": "single-choice",
      "options": [
        {"label": "write it now", "meaning": "the term is already settled", "recommended": true},
        {"label": "wait for the spec", "meaning": "batch every glossary change with the write"}
      ]
    }
  ]
}`

func TestADesignRoundArrivesChosenAndWhatThePersonAnsweredReachesTheSession(t *testing.T) {
	t.Parallel()
	held := standing(t)
	address := held.serve(designRoundDocument)
	page := opening(t, address)

	var order []string
	page.evaluate(`Array.from(document.querySelectorAll('.decision')).map((held) => held.dataset.decision)`, &order)
	if diff := cmp.Diff([]string{"q1", "q2"}, order); diff != "" {
		t.Errorf("a document stating no severity anywhere does not render in the order it listed the decisions (-want +got):\n%s", diff)
	}

	var titles []string
	page.evaluate(`Array.from(document.querySelectorAll('.decision .name')).map((held) => held.textContent.trim())`, &titles)
	if diff := cmp.Diff([]string{"Question 1", "Question 2"}, titles); diff != "" {
		t.Errorf("a design round's cards are not titled by the question number counted from one (-want +got):\n%s", diff)
	}

	var noSeverityShown bool
	page.evaluate(`!document.querySelector('.severity')`, &noSeverityShown)
	if !noSeverityShown {
		t.Error("the page shows a severity, and a design document states none")
	}

	var checked bool
	page.evaluate(`document.querySelector('.decision[data-decision="q1"] .choice[value="0"]').checked`, &checked)
	if !checked {
		t.Error("the recommended option does not arrive chosen when the page renders")
	}

	page.write(within("q2", ".words-field"), inTheirOwnWordsAboutTheCap)
	page.click(".finish")
	page.waitFor(`document.querySelector('.send').disabled === true`)

	printed, failures, code := held.poll()
	if code != exit.NothingAsked {
		t.Fatalf("agora poll exited %d: %s", code, failures)
	}

	var answered surface.Answered
	if err := json.Unmarshal([]byte(printed), &answered); err != nil {
		t.Fatalf("reading back what agora poll printed: %v\n%s", err, printed)
	}
	if len(answered.Decisions) != 2 {
		t.Fatalf("answering the round came back as %d decisions, want 2: %+v", len(answered.Decisions), answered.Decisions)
	}

	chosen, written := answered.Decisions[0], answered.Decisions[1]
	if chosen.ID != "q1" {
		t.Errorf("the first decision to come back is %q, want q1", chosen.ID)
	}
	if diff := cmp.Diff([]string{"the surface"}, chosen.Chosen); diff != "" {
		t.Errorf("the recommendation left untouched came back differently (-want +got):\n%s", diff)
	}
	if len(chosen.Commands) != 0 {
		t.Errorf("a design answer came back carrying the commands %v, and a design answer runs nothing", chosen.Commands)
	}

	if written.ID != "q2" {
		t.Errorf("the second decision to come back is %q, want q2", written.ID)
	}
	if written.Words != inTheirOwnWordsAboutTheCap {
		t.Errorf("what the person wrote came back as %q, want %q verbatim", written.Words, inTheirOwnWordsAboutTheCap)
	}
	if len(written.Chosen) != 0 || len(written.Commands) != 0 {
		t.Errorf("a decision answered in the person's own words came back with %v chosen and the commands %v", written.Chosen, written.Commands)
	}
	if !strings.Contains(printed, "q2") {
		t.Fatalf("what agora poll printed does not name q2 at all:\n%s", printed)
	}

	nothingLeftTheMachine(t, address, page.requests.all())
}

func TestASecondDesignRoundOpensWithWhatWasSettledAheadOfWhatIsStillOpen(t *testing.T) {
	t.Parallel()
	held := standing(t)
	address := held.serve(designRoundDocument)
	page := opening(t, address)

	page.write(within("q2", ".words-field"), inTheirOwnWordsAboutTheCap)
	page.click(".send")
	page.waitFor(`!document.querySelector('.receipt').hidden`)

	if _, failures, code := held.poll(); code != exit.NothingAsked {
		t.Fatalf("agora poll exited %d: %s", code, failures)
	}

	held.serveAgain(designRoundTwoDocument)
	page.waitFor(reachedThePage("q3"))

	var read struct {
		Settled   []string `json:"settled"`
		Decisions []string `json:"decisions"`
		First     bool     `json:"first"`
		Controls  int      `json:"controls"`
		Said      string   `json:"said"`
		Folded    bool     `json:"folded"`
	}
	page.evaluate(`(() => {
  const settled = document.querySelector('.settled');
  const decisions = document.querySelector('.decisions');
  return {
    settled: Array.from(document.querySelectorAll('.settled-answer')).map((held) => held.dataset.settled),
    decisions: Array.from(document.querySelectorAll('.decision')).map((held) => held.dataset.decision),
    first: !!settled && !!decisions && !!(settled.compareDocumentPosition(decisions) & Node.DOCUMENT_POSITION_FOLLOWING),
    controls: settled ? settled.querySelectorAll('input, textarea, button, select, ui-button, ui-field').length : -1,
    said: (document.querySelector('.settled-answer[data-settled="q2"]') || { textContent: '' }).textContent,
    folded: !!settled && !settled.open && settled.querySelector('.settled-summary').getBoundingClientRect().height > 0,
  };
})()`, &read)

	if diff := cmp.Diff([]string{"q1", "q2"}, read.Settled); diff != "" {
		t.Errorf("what was settled does not render as the document stated it (-want +got):\n%s", diff)
	}
	var settledNames string
	page.evaluate(`document.querySelector('.settled-link-names').textContent`, &settledNames)
	if settledNames != "Question 1, Question 2" {
		t.Errorf("what was settled is summarised as %q, want Question 1, Question 2", settledNames)
	}
	if diff := cmp.Diff([]string{"q3"}, read.Decisions); diff != "" {
		t.Errorf("round two renders something other than the one open question as a card (-want +got):\n%s", diff)
	}
	if !read.First {
		t.Errorf("what was settled is not shown ahead of what is still open: %+v", read)
	}
	if !read.Folded {
		t.Errorf("what was settled does not arrive folded under a heading the person can open: %+v", read)
	}
	if read.Controls != 0 {
		t.Errorf("what was settled offers %d controls, and a settled answer offers nothing to answer", read.Controls)
	}
	if !strings.Contains(read.Said, "no cap") || !strings.Contains(read.Said, "cap to protect") {
		t.Errorf("the settled answer for q2 does not carry the person's own words: %q", read.Said)
	}

	var openOptions int
	page.evaluate(`document.querySelectorAll('.decision[data-decision="q3"] .choice').length`, &openOptions)
	if openOptions == 0 {
		t.Error("the open question in round two offers no options at all")
	}

	var preamble string
	page.evaluate(`document.querySelector('.preamble').textContent`, &preamble)
	for _, said := range []string{"Question 1:", "the surface", "Question 2:", "no round cap"} {
		if !strings.Contains(preamble, said) {
			t.Errorf("the preamble does not carry %q, and the argument should read without scrolling past the settled block: %q", said, preamble)
		}
	}

	nothingLeftTheMachine(t, address, page.requests.all())
}

func TestASettledAnswerTheSessionParaphrasedIsRefusedAndThePageStaysWhereItWas(t *testing.T) {
	t.Parallel()
	held := standing(t)
	address := held.serve(designRoundDocument)
	page := opening(t, address)

	page.write(within("q2", ".words-field"), inTheirOwnWordsAboutTheCap)
	page.click(".send")
	page.waitFor(`!document.querySelector('.receipt').hidden`)

	if _, failures, code := held.poll(); code != exit.NothingAsked {
		t.Fatalf("agora poll exited %d: %s", code, failures)
	}

	paraphrased := strings.Replace(designRoundTwoDocument, inTheirOwnWordsAboutTheCap, "no cap, since a design spawns nothing a cap would protect", 1)
	_, failures, code := held.reading(paraphrased, held.serving("--session", held.opened)...)
	if code == exit.NothingAsked {
		t.Fatal("a second round settling q2 in words the person never wrote was served")
	}
	for _, named := range []string{"q2", "words", inTheirOwnWordsAboutTheCap, "agora poll"} {
		if !strings.Contains(failures, named) {
			t.Errorf("the refusal does not name %q:\n%s", named, failures)
		}
	}

	var stillOpen bool
	page.evaluate(reachedThePage("q2")+` && !document.querySelector('.decision[data-decision="q3"]')`, &stillOpen)
	if !stillOpen {
		t.Error("a refused round moved the page the person holds")
	}

	held.serveAgain(designRoundTwoDocument)
	page.waitFor(reachedThePage("q3"))
}
