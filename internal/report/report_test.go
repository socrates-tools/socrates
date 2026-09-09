package report_test

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates-review/internal/dialogue"
	"github.com/socrates-tools/socrates-review/internal/findings"
	"github.com/socrates-tools/socrates-review/internal/report"
	"github.com/socrates-tools/socrates-review/internal/runstate"
	"github.com/socrates-tools/socrates-review/internal/runstore"
)

const escalatingDocument = `{
  "findings": [
    {"id": "f1", "severity": "error", "file": "parse.go", "line": 12, "description": "An empty document parses as a valid one."},
    {"id": "f2", "severity": "warning", "file": "parse.go", "line": 40, "description": "The retry loop never backs off.", "action": "auto-fix"},
    {"id": "f3", "severity": "info", "file": "parse.go", "line": 88, "description": "This package now owns two spellings of the same idea.", "action": "no-op"}
  ],
  "risk_level": "medium",
  "risk_rationale": "one changed path answers wrongly without erroring"
}`

const informationalDocument = `{
  "findings": [
    {"id": "f1", "severity": "info", "file": "parse.go", "line": 3, "description": "The new name reads better than the old one.", "action": "no-op"}
  ],
  "risk_level": "low",
  "risk_rationale": "nothing in this change can answer wrongly"
}`

func TestTheEffectiveActionIsResolvedForTheReader(t *testing.T) {
	store, name := run(t, runstate.Escalated, escalatingDocument)

	reported, err := report.FindingsOf(store, name)
	if err != nil {
		t.Fatalf("reading the findings of %s: %v", name, err)
	}
	if got := reported.Findings[0].Resolved; got != findings.AskUser {
		t.Errorf("a finding that stated no action is reported as %q, want %q", got, findings.AskUser)
	}

	encoded, err := report.AsJSON(reported)
	if err != nil {
		t.Fatalf("encoding the findings: %v", err)
	}

	var read struct {
		Findings []map[string]any `json:"findings"`
	}
	if err := json.Unmarshal([]byte(encoded), &read); err != nil {
		t.Fatalf("reading back what socrates printed: %v", err)
	}

	var fields []string
	for field := range read.Findings[0] {
		fields = append(fields, field)
	}
	slices.Sort(fields)

	want := []string{"action", "description", "file", "id", "line", "severity"}
	if diff := cmp.Diff(want, fields); diff != "" {
		t.Errorf("the reader is given other fields than these (-want +got):\n%s", diff)
	}
	if read.Findings[0]["action"] != string(findings.AskUser) {
		t.Errorf("the printed action is %v, want %q", read.Findings[0]["action"], findings.AskUser)
	}
}

func TestFindingsAreReadable(t *testing.T) {
	store, name := run(t, runstate.Escalated, escalatingDocument)

	reported, err := report.FindingsOf(store, name)
	if err != nil {
		t.Fatalf("reading the findings of %s: %v", name, err)
	}

	printed := strings.Join(reported.Lines(), "\n")
	for _, needed := range []string{
		"f1", "error", "parse.go:12", "An empty document parses as a valid one.", "needs a person",
		"f2", "warning", "parse.go:40", "can be patched",
		"f3", "info", "informational",
		"medium", "one changed path answers wrongly without erroring",
	} {
		if !strings.Contains(printed, needed) {
			t.Errorf("what a person reads does not carry %q:\n%s", needed, printed)
		}
	}
	if reported.Code() != 1 {
		t.Errorf("an escalated run exits %d, want 1", reported.Code())
	}
}

func TestACleanRunHasFindingsToReadToo(t *testing.T) {
	store, name := run(t, runstate.Clean, informationalDocument)

	reported, err := report.FindingsOf(store, name)
	if err != nil {
		t.Fatalf("reading the findings of %s: %v", name, err)
	}
	if len(reported.Findings) != 1 {
		t.Fatalf("a clean run with one informational finding reported %d findings", len(reported.Findings))
	}
	if reported.Code() != 0 {
		t.Errorf("a clean run exits %d, want 0", reported.Code())
	}
	if printed := strings.Join(reported.Lines(), "\n"); !strings.Contains(printed, "The new name reads better than the old one.") {
		t.Errorf("what a person reads does not carry the informational finding:\n%s", printed)
	}
}

func TestAFailedRunHasNoFindingsToGive(t *testing.T) {
	store, name := run(t, runstate.Failed, escalatingDocument)

	reported, err := report.FindingsOf(store, name)
	if err != nil {
		t.Fatalf("reading the findings of %s: %v", name, err)
	}
	if len(reported.Findings) != 0 {
		t.Errorf("a failed run gave %d findings", len(reported.Findings))
	}
	if reported.Reviewed {
		t.Error("a failed run is reported as reviewed")
	}
	if reported.Code() != 2 {
		t.Errorf("a failed run exits %d, want 2", reported.Code())
	}

	printed := strings.Join(reported.Lines(), "\n")
	if !strings.Contains(printed, "no review completed") {
		t.Errorf("what a person reads does not say that no review completed:\n%s", printed)
	}
	if strings.Contains(printed, "An empty document parses as a valid one.") {
		t.Errorf("a failed run printed the findings a reviewer left behind:\n%s", printed)
	}
}

func TestARunningReviewIsAFewLines(t *testing.T) {
	store, name := run(t, runstate.Reviewing, "")

	status, err := report.StatusOf(store, name)
	if err != nil {
		t.Fatalf("reading the status of %s: %v", name, err)
	}

	lines := status.Lines()
	if len(lines) > 5 {
		t.Errorf("a running review prints %d lines:\n%s", len(lines), strings.Join(lines, "\n"))
	}
	if status.Code() != 4 {
		t.Errorf("a running review exits %d, want 4", status.Code())
	}
}

func TestATerminalRunReportsItsVerdict(t *testing.T) {
	store, name := run(t, runstate.Escalated, escalatingDocument)

	status, err := report.StatusOf(store, name)
	if err != nil {
		t.Fatalf("reading the status of %s: %v", name, err)
	}
	if status.Findings == nil {
		t.Fatal("a terminal run reports no counts")
	}

	want := &report.Counts{Total: 3, NeedingAPerson: 1, Patchable: 1, Informational: 1}
	if diff := cmp.Diff(want, status.Findings); diff != "" {
		t.Errorf("the counts differ (-want +got):\n%s", diff)
	}

	printed := strings.Join(status.Lines(), "\n")
	if !strings.Contains(printed, string(runstate.Escalated)) {
		t.Errorf("the status does not name the verdict:\n%s", printed)
	}
	if !strings.Contains(printed, "1 for the person") {
		t.Errorf("the status does not say how many findings need the person:\n%s", printed)
	}
}

func run(t *testing.T, state runstate.State, document string) (runstore.Store, string) {
	t.Helper()

	store := runstore.New(t.TempDir())
	created, err := store.Create(runstore.Start{
		StartedAt:      time.Date(2026, 8, 19, 10, 15, 0, 0, time.UTC),
		RepositoryRoot: t.TempDir(),
		Base:           "0000000000000000000000000000000000000000",
		Snapshot:       "1111111111111111111111111111111111111111",
	})
	if err != nil {
		t.Fatalf("creating a run: %v", err)
	}

	round, err := store.CreateRound(created.Name, 1)
	if err != nil {
		t.Fatalf("creating the round: %v", err)
	}
	if document != "" {
		if err := round.Write(round.FindingsPath(), document); err != nil {
			t.Fatalf("writing the findings: %v", err)
		}
	}

	if state.Terminal() {
		if err := created.End(state, "a reviewer reported what it found"); err != nil {
			t.Fatalf("ending the run as %s: %v", state, err)
		}
	} else if err := created.Advance(state); err != nil {
		t.Fatalf("advancing the run to %s: %v", state, err)
	}
	if err := store.Save(created); err != nil {
		t.Fatalf("saving the run: %v", err)
	}

	return store, created.Name
}

const oneChangedLine = `{
  "touched": {"parse.go": [{"first": 8, "last": 16}]}
}`

func TestTheStatusSaysWhichCommitTheChangeWasReviewedAgainst(t *testing.T) {
	store, name := run(t, runstate.Escalated, escalatingDocument)

	status, err := report.StatusOf(store, name)
	if err != nil {
		t.Fatalf("reading the status of %s: %v", name, err)
	}
	if status.Base != "0000000000000000000000000000000000000000" {
		t.Errorf("the status reports the base %q, want the commit the run recorded", status.Base)
	}
	if printed := strings.Join(status.Lines(), "\n"); !strings.Contains(printed, "base: 0000000000000000000000000000000000000000") {
		t.Errorf("what a person reads does not say which commit the change was reviewed against:\n%s", printed)
	}
}

func TestWhetherAFindingIsInsideTheChangeIsReported(t *testing.T) {
	store, name := run(t, runstate.Escalated, escalatingDocument)
	touching(t, store, name, oneChangedLine)

	reported, err := report.FindingsOf(store, name)
	if err != nil {
		t.Fatalf("reading the findings of %s: %v", name, err)
	}

	inside := map[string]bool{}
	for _, finding := range reported.Findings {
		if finding.Inside == nil {
			t.Fatalf("finding %s does not say whether it is inside the change", finding.ID)
		}
		inside[finding.ID] = *finding.Inside
	}

	want := map[string]bool{"f1": true, "f2": false, "f3": false}
	if diff := cmp.Diff(want, inside); diff != "" {
		t.Errorf("which findings are inside the change differs (-want +got):\n%s", diff)
	}
}

func TestNothingIsClaimedAboutAFindingWhoseChangeWasNeverRead(t *testing.T) {
	store, name := run(t, runstate.Escalated, escalatingDocument)

	reported, err := report.FindingsOf(store, name)
	if err != nil {
		t.Fatalf("reading the findings of %s: %v", name, err)
	}
	for _, finding := range reported.Findings {
		if finding.Inside != nil {
			t.Errorf("finding %s claims to be inside the change: %t, and nothing was read about what the change touched", finding.ID, *finding.Inside)
		}
	}
}

func touching(t *testing.T, store runstore.Store, name, change string) {
	t.Helper()

	round, present, err := store.LatestRound(name)
	if err != nil || !present {
		t.Fatalf("reading the round of %s: %v", name, err)
	}
	if err := round.Write(round.ChangePath(), change); err != nil {
		t.Fatalf("writing what the change touched: %v", err)
	}
}

func TestARecordedResolutionIsReported(t *testing.T) {
	store, name := run(t, runstate.AwaitingResponse, escalatingDocument)
	answering(t, store, name, dialogue.Response{Finding: "f2", Resolution: dialogue.Fixed, Reason: "the loop backs off now"})

	reported, err := report.FindingsOf(store, name)
	if err != nil {
		t.Fatalf("reading the findings of %s: %v", name, err)
	}

	recorded := map[string]dialogue.Resolution{}
	for _, finding := range reported.Findings {
		recorded[finding.ID] = finding.Resolution
	}
	if got := recorded["f2"]; got != dialogue.Fixed {
		t.Errorf("the answer recorded against f2 reads %q, want %q", got, dialogue.Fixed)
	}
	if got := recorded["f1"]; got != "" {
		t.Errorf("f1 has no answer recorded against it and reads %q", got)
	}
}

func TestAChangeThatIsNotMineIsAnsweredByItsOwnRules(t *testing.T) {
	store, name := run(t, runstate.AwaitingResponse, escalatingDocument)
	handingOver(t, store, name)

	reported, err := report.FindingsOf(store, name)
	if err != nil {
		t.Fatalf("reading the findings of %s: %v", name, err)
	}

	if diff := cmp.Diff(report.HandoverRules(), reported.Guidance); diff != "" {
		t.Errorf("the rules a person is given differ (-want +got):\n%s", diff)
	}
	if !reported.AwaitsAPerson {
		t.Error("a round of a change nobody here can patch waits on nobody")
	}

	printed := strings.Join(reported.Lines(), "\n")
	for _, named := range []string{"--forwarded", "--dropped"} {
		if !strings.Contains(printed, named) {
			t.Errorf("what a person reads does not name %s:\n%s", named, printed)
		}
	}
}

func answering(t *testing.T, store runstore.Store, name string, response dialogue.Response) {
	t.Helper()

	round, present, err := store.LatestRound(name)
	if err != nil || !present {
		t.Fatalf("reading the round of %s: %v", name, err)
	}
	if err := round.AppendResponse(response); err != nil {
		t.Fatalf("recording an answer: %v", err)
	}
}

func handingOver(t *testing.T, store runstore.Store, name string) {
	t.Helper()

	recorded, err := store.Load(name)
	if err != nil {
		t.Fatalf("reading the run %s: %v", name, err)
	}
	recorded.NotMyChange = true
	if err := store.Save(recorded); err != nil {
		t.Fatalf("recording that the change is not the person's: %v", err)
	}
}
