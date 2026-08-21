package report_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/rui-ktei/socrates/internal/dialogue"
	"github.com/rui-ktei/socrates/internal/report"
	"github.com/rui-ktei/socrates/internal/runstate"
	"github.com/rui-ktei/socrates/internal/runstore"
)

const patchableDocument = `{
  "findings": [
    {"id": "f1", "severity": "error", "file": "parse.go", "line": 12, "description": "An empty document parses as a valid one.", "action": "auto-fix"},
    {"id": "f2", "severity": "warning", "file": "parse.go", "line": 40, "description": "The retry loop never backs off.", "action": "auto-fix"},
    {"id": "f3", "severity": "info", "file": "parse.go", "line": 88, "description": "This package now owns two spellings of the same idea.", "action": "no-op"}
  ],
  "risk_level": "medium",
  "risk_rationale": "two changed paths answer wrongly without erroring"
}`

const questioningDocument = `{
  "findings": [
    {"id": "f1", "severity": "error", "file": "parse.go", "line": 12, "description": "This endpoint does not look like something the product needs.", "action": "ask-user"},
    {"id": "f2", "severity": "warning", "file": "parse.go", "line": 40, "description": "The retry loop never backs off.", "action": "auto-fix"},
    {"id": "f3", "severity": "info", "file": "parse.go", "line": 88, "description": "This package now owns two spellings of the same idea.", "action": "no-op"}
  ],
  "risk_level": "high",
  "risk_rationale": "the change may not be wanted at all"
}`

func TestAParkedRunReportsWhatItWants(t *testing.T) {
	store, name := run(t, runstate.AwaitingResponse, patchableDocument)

	status, err := report.StatusOf(store, name)
	if err != nil {
		t.Fatalf("reading the status of %s: %v", name, err)
	}
	if status.Unanswered == nil {
		t.Fatal("a parked run reports no count of what it is waiting for")
	}
	if *status.Unanswered != 2 {
		t.Errorf("a round with two blocking findings and no responses reports %d unanswered, want 2", *status.Unanswered)
	}
	if status.Code() != 4 {
		t.Errorf("a parked run exits %d, want 4", status.Code())
	}

	printed := strings.Join(status.Lines(), "\n")
	for _, needed := range []string{string(runstate.AwaitingResponse), "round: 1", "unanswered: 2", "awaiting responses"} {
		if !strings.Contains(printed, needed) {
			t.Errorf("the status does not say %q:\n%s", needed, printed)
		}
	}
	if strings.Contains(printed, "verdict:") {
		t.Errorf("a parked run printed a verdict:\n%s", printed)
	}
}

func TestAParkedRunSaysWhatNeedsAPerson(t *testing.T) {
	store, name := run(t, runstate.AwaitingResponse, questioningDocument)

	status, err := report.StatusOf(store, name)
	if err != nil {
		t.Fatalf("reading the status of %s: %v", name, err)
	}
	if status.Unanswered == nil || *status.Unanswered != 2 {
		t.Fatalf("a round holding a question and a patchable defect reports %v unanswered, want 2", status.Unanswered)
	}
	if status.AwaitingAPerson == nil || *status.AwaitingAPerson != 1 {
		t.Fatalf("a round holding one question reports %v awaiting a person, want 1", status.AwaitingAPerson)
	}

	printed := strings.Join(status.Lines(), "\n")
	for _, needed := range []string{"unanswered: 2", "person"} {
		if !strings.Contains(printed, needed) {
			t.Errorf("the status does not say %q:\n%s", needed, printed)
		}
	}

	encoded, err := report.AsJSON(status)
	if err != nil {
		t.Fatalf("encoding the status: %v", err)
	}
	var read map[string]any
	if err := json.Unmarshal([]byte(encoded), &read); err != nil {
		t.Fatalf("reading back what socrates printed: %v", err)
	}
	if read["unanswered"] != float64(2) {
		t.Errorf("the JSON object reports %v unanswered", read["unanswered"])
	}
	if read["awaiting_a_person"] != float64(1) {
		t.Errorf("the JSON object reports %v awaiting a person, so a caller that must stop and ask has to read prose", read["awaiting_a_person"])
	}
}

func TestTheCountReachesZero(t *testing.T) {
	store, name := run(t, runstate.AwaitingResponse, questioningDocument)
	answer(t, store, name, dialogue.Response{Finding: "f1", Resolution: dialogue.Decided, Reason: "Vic decided the endpoint stays"})
	answer(t, store, name, dialogue.Response{Finding: "f2", Resolution: dialogue.Rejected, Reason: "the loop is bounded by the caller's deadline"})

	parked, err := report.StatusOf(store, name)
	if err != nil {
		t.Fatalf("reading the status of %s: %v", name, err)
	}
	if parked.AwaitingAPerson == nil || *parked.AwaitingAPerson != 0 {
		t.Errorf("the question is answered and the status reports %v awaiting a person", parked.AwaitingAPerson)
	}
}

func TestTheCountOfWhatAPatchableRoundWantsReachesZero(t *testing.T) {
	// scenario: The count reaches zero
	store, name := run(t, runstate.AwaitingResponse, patchableDocument)
	answer(t, store, name, dialogue.Response{Finding: "f1", Resolution: dialogue.Fixed})
	answer(t, store, name, dialogue.Response{Finding: "f2", Resolution: dialogue.Rejected, Reason: "the loop is bounded by the caller's deadline"})

	status, err := report.StatusOf(store, name)
	if err != nil {
		t.Fatalf("reading the status of %s: %v", name, err)
	}
	if status.Unanswered == nil || *status.Unanswered != 0 {
		t.Fatalf("every answerable finding is answered and the status reports %v unanswered", status.Unanswered)
	}

	printed := strings.Join(status.Lines(), "\n")
	if !strings.Contains(printed, "unanswered: 0") {
		t.Errorf("the status does not report zero unanswered:\n%s", printed)
	}
	if !strings.Contains(printed, "scrat continue") {
		t.Errorf("the status does not say the run is ready to continue:\n%s", printed)
	}
}

func TestTheResponseRulesArriveWithTheFindings(t *testing.T) {
	t.Run("The rules arrive with the findings", func(t *testing.T) {
		store, name := run(t, runstate.AwaitingResponse, patchableDocument)

		reported, err := report.FindingsOf(store, name)
		if err != nil {
			t.Fatalf("reading the findings of %s: %v", name, err)
		}
		if len(reported.Findings) != 3 {
			t.Fatalf("a parked run reported %d findings, want the three the round found", len(reported.Findings))
		}

		printed := strings.Join(reported.Lines(), "\n")
		for _, rule := range report.ResponseRules() {
			if !strings.Contains(printed, rule) {
				t.Errorf("what the author reads does not carry the rule %q:\n%s", rule, printed)
			}
		}
		if !strings.Contains(printed, "An empty document parses as a valid one.") {
			t.Errorf("the rules arrived without the findings:\n%s", printed)
		}
	})

	t.Run("The rules do not arrive where they do not apply", func(t *testing.T) {
		for _, state := range []runstate.State{runstate.Clean, runstate.Escalated, runstate.Failed} {
			store, name := run(t, state, patchableDocument)

			reported, err := report.FindingsOf(store, name)
			if err != nil {
				t.Fatalf("reading the findings of a %s run: %v", state, err)
			}
			if reported.Guidance != nil {
				t.Errorf("a %s run is asked for responses: %v", state, reported.Guidance)
			}

			printed := strings.Join(reported.Lines(), "\n")
			for _, rule := range report.ResponseRules() {
				if strings.Contains(printed, rule) {
					t.Errorf("a %s run prints the response rule %q:\n%s", state, rule, printed)
				}
			}
		}
	})
}

func answer(t *testing.T, store runstore.Store, name string, response dialogue.Response) {
	t.Helper()

	if err := store.Round(name, 1).AppendResponse(response); err != nil {
		t.Fatalf("recording a response to %s: %v", name, err)
	}
}
