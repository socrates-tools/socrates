package e2e_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/agora/surface"
	"github.com/socrates-tools/socrates/internal/exit"
)

const lastDesignCardDocument = `{
  "title": "one card before the spec is written",
  "preamble": "q1: the surface, matching what agora's own source already calls it. q2: no round cap, because nothing here spawns a process for one to protect.",
  "decisions": [
    {
      "id": "q1",
      "question": "what should the loopback server be called in the glossary",
      "meaning": "settled: the surface, matching what agora's own source already calls it"
    },
    {
      "id": "q2",
      "question": "should a design round carry a cap the way a review round does",
      "meaning": "settled: no cap, because nothing here spawns a process for a cap to protect"
    },
    {
      "id": "write",
      "question": "write the spec now",
      "meaning": "shape: Problem Statement, Solution, User Stories, Implementation Decisions, Testing Decisions, Out of Scope, Further Notes. path: docs/changes/socrates-design/spec.md",
      "mode": "single-choice",
      "options": [
        {"label": "write it", "recommended": true}
      ]
    }
  ]
}`

func TestTheLastDesignRoundIsOneCardApprovingTheShapeAndThePath(t *testing.T) {
	held := standing(t)
	address := held.serve(lastDesignCardDocument)
	page := opening(t, address)

	var offering []string
	page.evaluate(`Array.from(document.querySelectorAll('.decision')).filter((held) => held.querySelectorAll('.choice').length > 0).map((held) => held.dataset.decision)`, &offering)
	if diff := cmp.Diff([]string{"write"}, offering); diff != "" {
		t.Errorf("the last round offers options on more than one card, and it should be one card only (-want +got):\n%s", diff)
	}

	var meaning string
	page.evaluate(`document.querySelector('.decision[data-decision="write"] .meaning').textContent`, &meaning)
	for _, said := range []string{"Problem Statement", "Further Notes", "docs/changes/socrates-design/spec.md"} {
		if !strings.Contains(meaning, said) {
			t.Errorf("the last card does not say %q, so a person could not approve a shape and a path they cannot see: %q", said, meaning)
		}
	}

	var checked bool
	page.evaluate(`document.querySelector('.decision[data-decision="write"] .choice[value="0"]').checked`, &checked)
	if !checked {
		t.Error("the recommended way to write the spec does not arrive chosen on the last card")
	}

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
	if !answered.Ended {
		t.Error("approving the last card does not end the design")
	}

	var approved *surface.Decided
	for position := range answered.Decisions {
		if answered.Decisions[position].ID == "write" {
			approved = &answered.Decisions[position]
		}
	}
	if approved == nil {
		t.Fatalf("what came back does not carry the write decision at all: %+v", answered.Decisions)
	}
	if diff := cmp.Diff([]string{"write it"}, approved.Chosen); diff != "" {
		t.Errorf("what the person approved on the last card differs (-want +got):\n%s", diff)
	}
	if len(approved.Commands) != 0 {
		t.Errorf("the last card came back carrying the commands %v, and a design answer runs nothing", approved.Commands)
	}
}

const twoOpenDesignQuestionsDocument = `{
  "title": "two questions before either is settled",
  "decisions": [
    {
      "id": "q1",
      "question": "what should the loopback server be called in the glossary",
      "mode": "single-choice",
      "options": [
        {"label": "the surface", "recommended": true},
        {"label": "the host"}
      ]
    },
    {
      "id": "q2",
      "question": "should a design round carry a cap the way a review round does",
      "mode": "single-choice",
      "options": [
        {"label": "no cap"},
        {"label": "cap it at two"}
      ]
    }
  ]
}`

func TestEndingADesignRoundEarlyDeliversOnlyWhatWasAnsweredAndClosesTheExchange(t *testing.T) {
	held := standing(t)
	address := held.serve(twoOpenDesignQuestionsDocument)
	page := opening(t, address)

	page.evaluate(`document.querySelector('.decision[data-decision="q1"] .choice[value="0"]').click()`, nil)
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
	if !answered.Ended {
		t.Fatal("ending the design early does not end the exchange")
	}
	if len(answered.Decisions) != 1 || answered.Decisions[0].ID != "q1" {
		t.Fatalf("ending early delivered %+v, want only q1, so a session cannot tell what stayed unsettled", answered.Decisions)
	}

	rendersClosed(t, page, "a design page whose round was ended before the frontier was empty")
}

func TestEndingADesignSessionFromTheCommandLineLeavesNothingForTheSessionToPoll(t *testing.T) {
	held := standing(t)
	address := held.serve(twoOpenDesignQuestionsDocument)

	held.end(held.opened)

	if delivered, arrived := held.awaiting(shortWindow); arrived {
		t.Fatalf("closing a design session from the command line without an answer delivered one anyway: %+v", delivered)
	}

	rendersClosed(t, opening(t, address), "a design page opened after socrates agora end closed the session from outside")
}
