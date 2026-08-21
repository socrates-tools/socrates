package runs_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/rui-ktei/socrates/internal/corpus/runs"
	"github.com/rui-ktei/socrates/internal/runlog"
	"github.com/rui-ktei/socrates/internal/runstore"
)

const returned = `{"findings":[{"id":"f1","severity":"warning","file":"a.go","line":3,"description":"d","action":"ask-user"}],"risk_level":"medium","risk_rationale":"r"}`

func TestAnAnswerIsTakenFromACompletedRun(t *testing.T) {
	store, run := completed(t, returned, map[string]any{"model": "opus", "effort": "high"})

	answered, err := runs.Answered(store, runs.Request{
		Run:     run.Name,
		Round:   1,
		Name:    "one-live-round",
		Subject: "a change worth one line",
		On:      "2026-08-21",
	})
	if err != nil {
		t.Fatalf("taking the answer from the run: %v", err)
	}

	if answered.Repository != run.Repository {
		t.Errorf("the repository is %s, want %s", answered.Repository, run.Repository)
	}
	if err := answered.Entry.Validate(); err != nil {
		t.Errorf("the entry taken from a completed run is incomplete: %v", err)
	}
	if !bytes.Equal(answered.Entry.Document, []byte(returned)) {
		t.Errorf("the document is not what the reviewer returned:\nreturned %q\ntaken    %q", returned, answered.Entry.Document)
	}

	for field, want := range map[string]string{
		"model":    "opus",
		"effort":   "high",
		"base":     "base-of-round-one",
		"snapshot": "snapshot-of-round-one",
		"run":      run.Name,
	} {
		got := map[string]string{
			"model":    answered.Entry.Model,
			"effort":   answered.Entry.Effort,
			"base":     answered.Entry.Base,
			"snapshot": answered.Entry.Snapshot,
			"run":      answered.Entry.Run,
		}[field]
		if got != want {
			t.Errorf("the entry states %s %q, want %q", field, got, want)
		}
	}
	if answered.Entry.Round != 1 {
		t.Errorf("the entry states round %d, want 1", answered.Entry.Round)
	}
}

func TestARunWhoseLogNamesNoModelIsRefused(t *testing.T) {
	store, run := completed(t, returned, map[string]any{"process": 1})

	_, err := runs.Answered(store, runs.Request{Run: run.Name, Round: 1, Name: "one-live-round", Subject: "s", On: "2026-08-21"})
	if err == nil {
		t.Fatal("a run whose log names no model is accepted")
	}
	if !strings.Contains(err.Error(), "model") {
		t.Errorf("the refusal does not name what is absent: %v", err)
	}
}

func TestARoundThatAnsweredNothingIsRefused(t *testing.T) {
	store, run := completed(t, "", map[string]any{"model": "opus", "effort": "high"})

	_, err := runs.Answered(store, runs.Request{Run: run.Name, Round: 1, Name: "one-live-round", Subject: "s", On: "2026-08-21"})
	if err == nil {
		t.Fatal("a round that answered nothing is accepted")
	}
	if !strings.Contains(err.Error(), run.Name) {
		t.Errorf("the refusal does not name the run: %v", err)
	}
}

func completed(t *testing.T, document string, fields map[string]any) (runstore.Store, *runstore.Run) {
	t.Helper()

	store := runstore.New(filepath.Join(t.TempDir(), "runs"))
	run, err := store.Create(runstore.Start{
		StartedAt:      time.Date(2026, 8, 21, 1, 2, 3, 0, time.UTC),
		RepositoryRoot: t.TempDir(),
		Base:           "base-of-round-one",
		Snapshot:       "snapshot-of-round-one",
		Intent:         "one line of intent",
	})
	if err != nil {
		t.Fatalf("creating a run: %v", err)
	}

	run.Snapshotted(1, "base-of-round-one", "snapshot-of-round-one")
	if err := store.Save(run); err != nil {
		t.Fatalf("saving the run: %v", err)
	}

	round, err := store.CreateRound(run.Name, 1)
	if err != nil {
		t.Fatalf("creating a round: %v", err)
	}
	if document != "" {
		if err := round.Write(round.FindingsPath(), document+"\n"); err != nil {
			t.Fatalf("writing the findings: %v", err)
		}
	}

	logged, err := json.Marshal(runlog.Event{
		At:      time.Date(2026, 8, 21, 1, 2, 4, 0, time.UTC),
		Kind:    "agent-started",
		Subject: "reviewer",
		Fields:  fields,
	})
	if err != nil {
		t.Fatalf("encoding an event: %v", err)
	}
	if err := os.WriteFile(filepath.Join(store.Path(run.Name), runlog.FileName), append(logged, '\n'), 0o644); err != nil {
		t.Fatalf("writing the event log: %v", err)
	}

	return store, run
}
