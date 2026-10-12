package e2e_test

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/agora/surface"
	"github.com/socrates-tools/socrates/internal/exit"
)

const lastOpenDesignQuestionDocument = `{
  "title": "the one question still open",
  "settled": [
    {
      "id": "q1",
      "name": "Question 1",
      "question": "what should the loopback server be called in the glossary",
      "chosen": ["the surface"],
      "words": "matching what agora's own source already calls it"
    }
  ],
  "decisions": [
    {
      "id": "q2",
      "name": "Question 2",
      "question": "should a design round carry a cap the way a review round does",
      "mode": "single-choice",
      "options": [
        {"label": "no cap", "recommended": true},
        {"label": "cap it at two"}
      ]
    }
  ]
}`

func TestAnsweringTheLastOpenQuestionLeavesThePageClosedOnceTheSessionEndsIt(t *testing.T) {
	t.Parallel()
	held := standing(t)
	address := held.serve(lastOpenDesignQuestionDocument)
	page := opening(t, address)

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
	if answered.Ended {
		t.Fatal("sending and waiting for a reply ended the exchange, so the session would have nothing left to close")
	}
	if len(answered.Decisions) != 1 || answered.Decisions[0].ID != "q2" {
		t.Fatalf("answering the last open question delivered %+v, want only q2", answered.Decisions)
	}
	if diff := cmp.Diff([]string{"no cap"}, answered.Decisions[0].Chosen); diff != "" {
		t.Errorf("the recommendation left untouched came back differently (-want +got):\n%s", diff)
	}

	held.end(held.opened)
	page.waitFor(`document.querySelector('.state').textContent.toLowerCase().includes('over')`)

	rendersClosed(t, page, "a design page the session closed once nothing was left to ask")
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
	t.Parallel()
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
	t.Parallel()
	held := standing(t)
	address := held.serve(twoOpenDesignQuestionsDocument)

	held.end(held.opened)

	if delivered, arrived := held.awaiting(shortWindow); arrived {
		t.Fatalf("closing a design session from the command line without an answer delivered one anyway: %+v", delivered)
	}

	rendersClosed(t, opening(t, address), "a design page opened after socrates agora end closed the session from outside")
}
