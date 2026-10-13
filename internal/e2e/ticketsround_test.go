package e2e_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/agora/surface"
	"github.com/socrates-tools/socrates/internal/exit"
)

const breakdownDocument = `{
  "title": "the proposed tickets for socrates-tickets",
  "preamble": "cut from docs/changes/socrates-tickets/spec.md: two tickets, the second blocked by the first",
  "decisions": [
    {
      "id": "1-prefactor",
      "question": "1: the skill machinery takes a second design document",
      "meaning": "delivers: socrates design skill installs two skills",
      "evidence": {"language": "markdown", "code": "# 1: the skill machinery takes a second design document\n\n## What to build\n\nsocrates design skill installs two skills.\n\n## Acceptance criteria\n\n- [ ] socrates design skill writes a second SKILL.md, and before this ticket it writes one\n\n## User stories\n\n- 1\n\n## Blocked by\n\nNone"},
      "mode": "single-choice",
      "options": [{"label": "keep it", "meaning": "cut as proposed", "recommended": true}]
    },
    {
      "id": "2-tickets-skill",
      "question": "2: the tickets skill states the breakdown, the main ticket and the tickets",
      "meaning": "delivers: a session can cut a spec into tickets",
      "evidence": {"language": "markdown", "code": "# 2: the tickets skill states the breakdown, the main ticket and the tickets\n\n## What to build\n\nA session can cut a spec into tickets.\n\n## Acceptance criteria\n\n- [ ] the installed tickets skill names the main ticket, and before this ticket nothing does\n\n## User stories\n\n- 2, 3\n\n## Blocked by\n\n- 1-prefactor.md"},
      "mode": "single-choice",
      "options": [{"label": "keep it", "meaning": "cut as proposed", "recommended": true}]
    }
  ]
}`

const askingForASplit = "split this: the main ticket and the task tickets are two slices"

func TestABreakdownArrivesKeptAndWhatThePersonAskedToChangeReachesTheSession(t *testing.T) {
	t.Parallel()
	held := standing(t)
	address := held.serve(breakdownDocument)
	page := opening(t, address)

	var order []string
	page.evaluate(`Array.from(document.querySelectorAll('.decision')).map((held) => held.dataset.decision)`, &order)
	if diff := cmp.Diff([]string{"1-prefactor", "2-tickets-skill"}, order); diff != "" {
		t.Errorf("the breakdown does not render one card per ticket in the order it proposed them (-want +got):\n%s", diff)
	}

	var shown []string
	page.evaluate(`['1-prefactor', '2-tickets-skill'].map((id) => (document.querySelector('.decision[data-decision="' + id + '"] .evidence.code') || { innerText: '' }).innerText)`, &shown)
	for position, written := range []string{"socrates design skill writes a second SKILL.md", "- 1-prefactor.md"} {
		if position >= len(shown) || !strings.Contains(shown[position], "## Acceptance criteria") || !strings.Contains(shown[position], written) {
			t.Errorf("card %d does not show the whole ticket it would be written as, so its acceptance criteria and blockers would be approved unseen: %q", position+1, shown)
		}
	}

	var scrolled bool
	page.evaluate(`Array.from(document.querySelectorAll('.evidence.code pre')).some((held) => held.scrollWidth > held.clientWidth)`, &scrolled)
	if scrolled {
		t.Error("a ticket on its card runs past the edge of the card, so reading its acceptance criteria takes scrolling sideways")
	}

	var kept []bool
	page.evaluate(`['1-prefactor', '2-tickets-skill'].map((id) => document.querySelector('.decision[data-decision="' + id + '"] .choice[value="0"]').checked)`, &kept)
	if diff := cmp.Diff([]bool{true, true}, kept); diff != "" {
		t.Errorf("keep it does not arrive chosen on every card, so agreeing with the breakdown costs more than one press (-want +got):\n%s", diff)
	}

	page.evaluate(declined("2-tickets-skill"), nil)
	page.write(within("2-tickets-skill", ".words-field"), askingForASplit)
	page.click(".send")
	page.waitFor(`!document.querySelector('.receipt').hidden`)

	printed, failures, code := held.poll()
	if code != exit.NothingAsked {
		t.Fatalf("agora poll exited %d: %s", code, failures)
	}

	var answered surface.Answered
	if err := json.Unmarshal([]byte(printed), &answered); err != nil {
		t.Fatalf("reading back what agora poll printed: %v\n%s", err, printed)
	}
	if len(answered.Decisions) != 2 {
		t.Fatalf("answering the breakdown came back as %d decisions, want 2: %+v", len(answered.Decisions), answered.Decisions)
	}

	keptAsProposed, askedToChange := answered.Decisions[0], answered.Decisions[1]
	if keptAsProposed.ID != "1-prefactor" {
		t.Errorf("the first card to come back is %q, want 1-prefactor", keptAsProposed.ID)
	}
	if diff := cmp.Diff([]string{"keep it"}, keptAsProposed.Chosen); diff != "" {
		t.Errorf("a card left as it arrived came back differently (-want +got):\n%s", diff)
	}
	if len(keptAsProposed.Commands) != 0 {
		t.Errorf("keeping a ticket came back carrying the commands %v, and a breakdown answer runs nothing", keptAsProposed.Commands)
	}

	if askedToChange.ID != "2-tickets-skill" {
		t.Errorf("the second card to come back is %q, want 2-tickets-skill", askedToChange.ID)
	}
	if askedToChange.Words != askingForASplit {
		t.Errorf("what the person asked for came back as %q, want %q verbatim", askedToChange.Words, askingForASplit)
	}
	if len(askedToChange.Chosen) != 0 || len(askedToChange.Commands) != 0 {
		t.Errorf("a card answered in the person's own words came back with %v chosen and the commands %v", askedToChange.Chosen, askedToChange.Commands)
	}

	nothingLeftTheMachine(t, address, page.requests.all())
}
