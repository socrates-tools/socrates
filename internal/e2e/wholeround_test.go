package e2e_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/socrates-review/runstate"
)

const fourFindings = `{"findings":[{"id":"f1","severity":"warning","file":"internal/added.go","line":3,"description":"The retry loop never backs off.","action":"auto-fix"},{"id":"f2","severity":"info","file":"internal/added.go","line":1,"description":"This package now owns two spellings of one idea.","action":"no-op"},{"id":"f3","severity":"error","file":"internal/added.go","line":3,"description":"This function does not look like something the product needs.","action":"ask-user"},{"id":"f4","severity":"error","file":"internal/added.go","line":5,"description":"The lock deletes itself while another request still holds it.","action":"auto-fix"}],"risk_level":"high","risk_rationale":"the change may not be wanted at all, and one path answers wrongly without erroring"}`

func theQuestionAsADecision(read reportedRound) []map[string]any {
	var decisions []map[string]any

	for _, finding := range read.Findings {
		decision := map[string]any{
			"id":       finding.ID,
			"question": finding.Description,
			"meaning":  "this is what it costs the change",
		}
		switch finding.ID {
		case "f3":
			decision["mode"] = "single-choice"
			decision["options"] = []map[string]any{
				{"label": "the function stays", "meaning": "the migration reads it", "command": `socrates-review respond --finding f3 --decided --reason "the endpoint stays, because the migration reads it"`, "recommended": true},
				{"label": "take it out", "meaning": "nothing reads it", "command": `socrates-review respond --finding f3 --decided --reason "nothing reads it"`},
			}
		case "f4":
			decision["meaning"] = "the deletion is gone, and the entry is dropped by whichever reader took it last"
			decision["evidence"] = map[string]any{
				"language": "go",
				"patch":    "--- a/internal/added.go\n+++ b/internal/added.go\n@@ -3,3 +3,2 @@ func Answer(document string) int {\n func Answer(document string) int {\n-\tdelete(locks, document)\n \treturn len(document)\n }\n",
			}
		}
		decisions = append(decisions, decision)
	}
	return decisions
}

func TestTheWholeRoundIsHandedOverAsCardsInSeverityOrder(t *testing.T) {
	const inTheirOwnWords = "the deletion was the bug I was chasing all morning"

	driven := driving(t, dialogueScenario(fourFindings, nothingFoundAgain), answeredWrongly)
	name := driven.review("count the characters of a document")

	state, _ := driven.waitedFor(name)
	if state != runstate.AwaitingResponse {
		t.Fatalf("a round holding a question of intent left the run %s, want a parked run", state)
	}
	if reported := driven.findings(name); !strings.Contains(reported, awaitsAPerson) {
		t.Fatalf("socrates does not report the round as one that waits on a person:\n%s", reported)
	}

	read := driven.reported(name)
	if len(read.Findings) != 4 {
		t.Fatalf("socrates reported %d findings for the round, want the four the reviewer returned", len(read.Findings))
	}

	address := driven.serve(documentFor(t, "what an independent reviewer said", read, theQuestionAsADecision(read)))
	page := opening(t, address)

	var cards, listed []string
	page.evaluate(`Array.from(document.querySelectorAll('.decision')).map((held) => held.dataset.decision)`, &cards)
	page.evaluate(`Array.from(document.querySelectorAll('.entry .entry-name')).map((held) => held.textContent.trim())`, &listed)

	if diff := cmp.Diff([]string{"f3", "f4", "f1", "f2"}, cards); diff != "" {
		t.Errorf("the cards a person meets differ (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]string{"Finding 3", "Finding 4", "Finding 1", "Finding 2"}, listed); diff != "" {
		t.Errorf("what a person sees at a glance differs (-want +got):\n%s", diff)
	}

	var asking []string
	page.evaluate(`Array.from(new Set(Array.from(document.querySelectorAll('.decision .choice')).map((held) => held.dataset.decision)))`, &asking)
	if diff := cmp.Diff([]string{"f3"}, asking); diff != "" {
		t.Errorf("the cards offering a choice differ (-want +got):\n%s\na control on a settled finding reads as a question still owed an answer", diff)
	}

	var patched bool
	page.evaluate(`!!document.querySelector('.decision[data-decision="f4"] .hunk .line[data-kind="removed"]')`, &patched)
	if !patched {
		t.Error("the card for the finding the session resolved carries no diff, and the diff is what makes the resolution clear")
	}

	page.evaluate(chose("f3", 0), nil)
	page.write(within("f4", ".words-field"), inTheirOwnWords)
	pressPastTheNotice(page, ".finish")
	page.waitFor(`document.querySelector('.send').disabled === true`)

	answered := driven.polled()
	if len(answered.Decisions) != 2 {
		t.Fatalf("answering one question and writing on one settled card came back as %d answers: %+v", len(answered.Decisions), answered.Decisions)
	}

	decided, annotated := answered.Decisions[0], answered.Decisions[1]
	if diff := cmp.Diff([]string{"the function stays"}, decided.Chosen); diff != "" {
		t.Errorf("what the person chose came back differently (-want +got):\n%s", diff)
	}
	if len(decided.Commands) != 1 || !strings.Contains(decided.Commands[0], "--finding f3 --decided") {
		t.Errorf("the option chosen came back carrying %v", decided.Commands)
	}
	if annotated.ID != "f4" || annotated.Words != inTheirOwnWords {
		t.Errorf("what the person wrote on the settled card came back as %+v", annotated)
	}
	if len(annotated.Chosen) != 0 || len(annotated.Commands) != 0 {
		t.Errorf("the settled card came back with %v chosen and the commands %v", annotated.Chosen, annotated.Commands)
	}

	nothingLeftTheMachine(t, address, page.requests.all())
	driven.wroteNothingItself()
	driven.namedNothingTheSkillDoesNot()
}
