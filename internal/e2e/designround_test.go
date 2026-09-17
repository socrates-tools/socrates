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

func TestADesignRoundArrivesChosenAndWhatThePersonAnsweredReachesTheSession(t *testing.T) {
	held := standing(t)
	address := held.serve(designRoundDocument)
	page := opening(t, address)

	var order []string
	page.evaluate(`Array.from(document.querySelectorAll('.decision')).map((held) => held.dataset.decision)`, &order)
	if diff := cmp.Diff([]string{"q1", "q2"}, order); diff != "" {
		t.Errorf("a document stating no severity anywhere does not render in the order it listed the decisions (-want +got):\n%s", diff)
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
