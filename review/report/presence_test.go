package report_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/review/dialogue"
	"github.com/socrates-tools/socrates/review/report"
	"github.com/socrates-tools/socrates/review/runstate"
	"github.com/socrates-tools/socrates/review/runstore"
)

func withNobodyToAsk(t *testing.T, store runstore.Store, name string) {
	t.Helper()

	recorded, err := store.Load(name)
	if err != nil {
		t.Fatalf("reading the run %s: %v", name, err)
	}
	recorded.NobodyToAsk = true
	if err := store.Save(recorded); err != nil {
		t.Fatalf("recording that nobody is there to ask: %v", err)
	}
}

func TestAParkedRunWithNobodyToAskSaysSoAndWaitsOnNoPerson(t *testing.T) {
	store, name := run(t, runstate.AwaitingResponse, questioningDocument)
	withNobodyToAsk(t, store, name)

	status, err := report.StatusOf(store, name)
	if err != nil {
		t.Fatalf("reading the status of %s: %v", name, err)
	}
	if !status.NobodyToAsk {
		t.Error("the status of a run begun with nobody to ask does not say so")
	}
	if status.Unanswered == nil || *status.Unanswered != 2 {
		t.Errorf("a round holding a question and a patchable defect reports %v unanswered, want 2: the question is still the session's to settle", status.Unanswered)
	}
	if status.AwaitingAPerson != nil {
		t.Errorf("a run with nobody to ask reports %d awaiting a person", *status.AwaitingAPerson)
	}

	printed := strings.Join(status.Lines(), "\n")
	if !strings.Contains(printed, "nobody to ask") {
		t.Errorf("the status does not say nobody is there to ask:\n%s", printed)
	}
	for _, absent := range []string{"needing a person rather than a patch", "or for a person the author asks"} {
		if strings.Contains(printed, absent) {
			t.Errorf("the status of a run with nobody to ask says %q:\n%s", absent, printed)
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
	if read["nobody_to_ask"] != true {
		t.Errorf("the JSON object reports nobody_to_ask as %v", read["nobody_to_ask"])
	}
}

func TestARunWithSomebodyToAskDoesNotSayNobodyIs(t *testing.T) {
	store, name := run(t, runstate.AwaitingResponse, questioningDocument)

	status, err := report.StatusOf(store, name)
	if err != nil {
		t.Fatalf("reading the status of %s: %v", name, err)
	}
	if status.NobodyToAsk || strings.Contains(strings.Join(status.Lines(), "\n"), "nobody to ask") {
		t.Errorf("a run begun without --nobody-to-ask says nobody is there to ask:\n%s", strings.Join(status.Lines(), "\n"))
	}
	encoded, _ := report.AsJSON(status)
	if strings.Contains(encoded, "nobody_to_ask") {
		t.Errorf("the JSON object of a run with somebody to ask carries nobody_to_ask: %s", encoded)
	}
}

func TestTheFindingsOfARunWithNobodyToAskSayItAndReportASettledResolution(t *testing.T) {
	store, name := run(t, runstate.AwaitingResponse, questioningDocument)
	withNobodyToAsk(t, store, name)
	answer(t, store, name, dialogue.Response{Finding: "f1", Resolution: dialogue.Settled, Reason: "the endpoint stays, because the migration reads it"})

	reported, err := report.FindingsOf(store, name)
	if err != nil {
		t.Fatalf("reading the findings of %s: %v", name, err)
	}
	if !reported.NobodyToAsk {
		t.Error("the findings of a run begun with nobody to ask do not say so")
	}
	if reported.AwaitsAPerson {
		t.Error("a parked round whose question the session settled awaits a person on a run with nobody to ask")
	}

	resolutions := map[string]dialogue.Resolution{}
	for _, finding := range reported.Findings {
		resolutions[finding.ID] = finding.Resolution
	}
	if resolutions["f1"] != dialogue.Settled {
		t.Errorf("f1 is reported as %q, want settled", resolutions["f1"])
	}

	encoded, err := report.AsJSON(reported)
	if err != nil {
		t.Fatalf("encoding the findings: %v", err)
	}
	var read map[string]any
	if err := json.Unmarshal([]byte(encoded), &read); err != nil {
		t.Fatalf("reading back what socrates printed: %v", err)
	}
	if read["nobody_to_ask"] != true {
		t.Errorf("the JSON object reports nobody_to_ask as %v", read["nobody_to_ask"])
	}
}
