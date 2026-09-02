package report_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/ktei/socrates/internal/report"
	"github.com/ktei/socrates/internal/runstate"
)

const beyondAnsweringDocument = `{
  "findings": [
    {"id": "f1", "severity": "error", "file": "parse.go", "line": 12, "description": "This whole approach reads the file twice and there is no version of it that does not.", "action": "no-op"},
    {"id": "f2", "severity": "warning", "file": "parse.go", "line": 40, "description": "The retry loop never backs off.", "action": "auto-fix"}
  ],
  "risk_level": "high",
  "risk_rationale": "the approach may be wrong rather than the code"
}`

func TestARoundThatWaitsOnAPersonSaysSo(t *testing.T) {
	for _, reported := range []struct {
		named    string
		state    runstate.State
		document string
	}{
		{named: "a run that escalated for a finding no response addresses", state: runstate.Escalated, document: beyondAnsweringDocument},
		{named: "a run that escalated at the cap", state: runstate.Escalated, document: patchableDocument},
		{named: "a parked round holding a question", state: runstate.AwaitingResponse, document: questioningDocument},
	} {
		t.Run(reported.named, func(t *testing.T) {
			store, name := run(t, reported.state, reported.document)

			read, err := report.FindingsOf(store, name)
			if err != nil {
				t.Fatalf("reading the findings of %s: %v", name, err)
			}
			if !read.AwaitsAPerson {
				t.Errorf("%s is not reported as a round that waits on a person, so nothing downstream knows to wait for what they answer", reported.named)
			}

			printed := strings.Join(read.Lines(), "\n")
			if !strings.Contains(printed, "awaits a person") {
				t.Errorf("%s does not say the round waits on a person:\n%s", reported.named, printed)
			}

			if field := jsonOf(t, read)["awaits_a_person"]; field != true {
				t.Errorf("%s reports %v in the object socrates prints, so a caller has to read prose", reported.named, field)
			}
		})
	}
}

func TestARoundTheAuthorAnswersAloneSaysNothingAboutWaiting(t *testing.T) {
	for _, reported := range []struct {
		named    string
		state    runstate.State
		document string
	}{
		{named: "a clean run", state: runstate.Clean, document: patchableDocument},
		{named: "a parked round the author answers alone", state: runstate.AwaitingResponse, document: patchableDocument},
		{named: "a run that failed before answering", state: runstate.Failed, document: ""},
		{named: "a cancelled run", state: runstate.Cancelled, document: questioningDocument},
	} {
		t.Run(reported.named, func(t *testing.T) {
			store, name := run(t, reported.state, reported.document)

			read, err := report.FindingsOf(store, name)
			if err != nil {
				t.Fatalf("reading the findings of %s: %v", name, err)
			}
			if read.AwaitsAPerson {
				t.Errorf("%s is reported as a round that waits on a person, so the session stops for an answer nobody was asked for", reported.named)
			}

			printed := strings.Join(read.Lines(), "\n")
			if strings.Contains(printed, "awaits a person") {
				t.Errorf("%s says the round waits on a person:\n%s", reported.named, printed)
			}
			if _, carried := jsonOf(t, read)["awaits_a_person"]; carried {
				t.Errorf("%s carries the field in the object socrates prints", reported.named)
			}

			for field := range jsonOf(t, read) {
				if namesAPage(field) {
					t.Errorf("%s reports %q, and whether a round is handed over as a page is nothing socrates states", reported.named, field)
				}
			}
			for _, line := range read.Lines() {
				if namesAPage(line) {
					t.Errorf("%s prints %q, and nothing socrates says would stop the round being handed over as a page", reported.named, line)
				}
			}
		})
	}
}

func TestScratStatusDoesNotCarryTheWait(t *testing.T) {
	store, name := run(t, runstate.AwaitingResponse, questioningDocument)

	status, err := report.StatusOf(store, name)
	if err != nil {
		t.Fatalf("reading the status of %s: %v", name, err)
	}

	printed := strings.Join(status.Lines(), "\n")
	if strings.Contains(printed, "awaits a person") {
		t.Errorf("status carries what only the command read once needs, and everything it prints lands in a polling session's context:\n%s", printed)
	}
	if _, carried := jsonOf(t, status)["awaits_a_person"]; carried {
		t.Error("the status object carries the field, and status is what a waiting session calls in a loop")
	}
}

func namesAPage(named string) bool {
	for _, page := range []string{"page", "artifact", "html", "publish", "render"} {
		if strings.Contains(strings.ToLower(named), page) {
			return true
		}
	}
	return false
}

func jsonOf(t *testing.T, reported any) map[string]any {
	t.Helper()

	encoded, err := report.AsJSON(reported)
	if err != nil {
		t.Fatalf("encoding what socrates concluded: %v", err)
	}
	var read map[string]any
	if err := json.Unmarshal([]byte(encoded), &read); err != nil {
		t.Fatalf("reading back what socrates printed: %v", err)
	}
	return read
}
