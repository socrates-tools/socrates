package runs_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/rui-ktei/socrates/internal/agent"
	"github.com/rui-ktei/socrates/internal/corpus/runs"
	"github.com/rui-ktei/socrates/internal/reviewer"
	"github.com/rui-ktei/socrates/internal/runlog"
	"github.com/rui-ktei/socrates/internal/runstore"
)

const returned = `{"findings":[{"id":"f1","severity":"warning","file":"a.go","line":3,"description":"d","action":"ask-user"}],"risk_level":"medium","risk_rationale":"r"}`

func TestAnAnswerIsTakenFromACompletedRun(t *testing.T) {
	store, run := completed(t, returned, started("opus", "high"))

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

	for field, pair := range map[string][2]string{
		"model":    {answered.Entry.Model, "opus"},
		"effort":   {answered.Entry.Effort, "high"},
		"base":     {answered.Entry.Base, "base-of-round-1"},
		"snapshot": {answered.Entry.Snapshot, "snapshot-of-round-1"},
		"run":      {answered.Entry.Run, run.Name},
	} {
		if pair[0] != pair[1] {
			t.Errorf("the entry states %s %q, want %q", field, pair[0], pair[1])
		}
	}
	if answered.Entry.Round != 1 {
		t.Errorf("the entry states round %d, want 1", answered.Entry.Round)
	}
}

func TestTheModelRecordedIsTheOneThatAnsweredThatRound(t *testing.T) {
	logged := started("opus", "high")
	logged = append(logged, runlog.Event{
		Kind:    reviewer.EventRoundStarted,
		Subject: "a-run",
		Fields:  map[string]any{"round": 2},
	})
	logged = append(logged, started("sonnet", "medium")...)

	store, run := completed(t, returned, logged)
	rounds(t, store, run, 2, returned)

	for round, want := range map[int][2]string{1: {"opus", "high"}, 2: {"sonnet", "medium"}} {
		answered, err := runs.Answered(store, runs.Request{Run: run.Name, Round: round, Name: "n", Subject: "s", On: "2026-08-21"})
		if err != nil {
			t.Fatalf("taking round %d out of the run: %v", round, err)
		}
		if answered.Entry.Model != want[0] || answered.Entry.Effort != want[1] {
			t.Errorf("round %d was recorded as %s at %s effort, and it was answered by %s at %s",
				round, answered.Entry.Model, answered.Entry.Effort, want[0], want[1])
		}
		if answered.Entry.Base != "base-of-round-"+string(rune('0'+round)) {
			t.Errorf("round %d was recorded against %s", round, answered.Entry.Base)
		}
	}
}

func TestARoundWhoseModelTheLogNoLongerHoldsIsRefused(t *testing.T) {
	logged := []runlog.Event{{
		Kind:    reviewer.EventRoundStarted,
		Subject: "a-run",
		Fields:  map[string]any{"round": 2},
	}}
	logged = append(logged, started("sonnet", "medium")...)

	store, run := completed(t, returned, logged)
	rounds(t, store, run, 2, returned)

	_, err := runs.Answered(store, runs.Request{Run: run.Name, Round: 1, Name: "n", Subject: "s", On: "2026-08-21"})
	if err == nil {
		t.Fatal("a round the log no longer holds a model for is recorded from a later round's")
	}
	if !strings.Contains(err.Error(), "1") {
		t.Errorf("the refusal does not name the round: %v", err)
	}
}

func TestARunWhoseLogNamesNoModelIsRefused(t *testing.T) {
	store, run := completed(t, returned, []runlog.Event{{
		Kind:    agent.EventStarted,
		Subject: "reviewer",
		Fields:  map[string]any{"process": 1},
	}})

	_, err := runs.Answered(store, runs.Request{Run: run.Name, Round: 1, Name: "one-live-round", Subject: "s", On: "2026-08-21"})
	if err == nil {
		t.Fatal("a run whose log names no model is accepted")
	}
	if !strings.Contains(err.Error(), "model") {
		t.Errorf("the refusal does not name what is absent: %v", err)
	}
}

func TestARoundThatAnsweredNothingIsRefused(t *testing.T) {
	store, run := completed(t, "", started("opus", "high"))

	_, err := runs.Answered(store, runs.Request{Run: run.Name, Round: 1, Name: "one-live-round", Subject: "s", On: "2026-08-21"})
	if err == nil {
		t.Fatal("a round that answered nothing is accepted")
	}
	if !strings.Contains(err.Error(), run.Name) {
		t.Errorf("the refusal does not name the run: %v", err)
	}
}

func TestARoundThatAnsweredSomethingUnreadableIsRefused(t *testing.T) {
	store, run := completed(t, `{"findings":"none that matter","risk_level":"low","risk_rationale":"r"}`, started("opus", "high"))

	_, err := runs.Answered(store, runs.Request{Run: run.Name, Round: 1, Name: "one-live-round", Subject: "s", On: "2026-08-21"})
	if err == nil {
		t.Fatal("an answer socrates never routed is recorded as evidence")
	}
	if !strings.Contains(err.Error(), run.Name) {
		t.Errorf("the refusal does not name the run: %v", err)
	}
}

func started(model, effort string) []runlog.Event {
	return []runlog.Event{{
		Kind:    agent.EventStarted,
		Subject: "reviewer",
		Fields:  map[string]any{"model": model, "effort": effort, "process": 1},
	}}
}

func completed(t *testing.T, document string, logged []runlog.Event) (runstore.Store, *runstore.Run) {
	t.Helper()

	store := runstore.New(filepath.Join(t.TempDir(), "runs"))
	run, err := store.Create(runstore.Start{
		StartedAt:      time.Date(2026, 8, 21, 1, 2, 3, 0, time.UTC),
		RepositoryRoot: t.TempDir(),
		Base:           "base-of-round-1",
		Snapshot:       "snapshot-of-round-1",
		Intent:         "one line of intent",
	})
	if err != nil {
		t.Fatalf("creating a run: %v", err)
	}
	rounds(t, store, run, 1, document)

	path := filepath.Join(store.Path(run.Name), runlog.FileName)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatalf("writing the event log: %v", err)
	}
	defer file.Close()

	for _, event := range logged {
		event.At = time.Date(2026, 8, 21, 1, 2, 4, 0, time.UTC)
		line, err := json.Marshal(event)
		if err != nil {
			t.Fatalf("encoding an event: %v", err)
		}
		if _, err := file.Write(append(line, '\n')); err != nil {
			t.Fatalf("writing the event log: %v", err)
		}
	}

	return store, run
}

func rounds(t *testing.T, store runstore.Store, run *runstore.Run, upTo int, document string) {
	t.Helper()

	for round := 1; round <= upTo; round++ {
		numbered := string(rune('0' + round))
		run.Snapshotted(round, "base-of-round-"+numbered, "snapshot-of-round-"+numbered)
		if _, err := store.CreateRound(run.Name, round); err != nil && !os.IsExist(err) {
			t.Fatalf("creating round %d: %v", round, err)
		}
		if document == "" {
			continue
		}
		created := store.Round(run.Name, round)
		if err := created.Write(created.FindingsPath(), document+"\n"); err != nil {
			t.Fatalf("writing the findings of round %d: %v", round, err)
		}
	}
	if err := store.Save(run); err != nil {
		t.Fatalf("saving the run: %v", err)
	}
}
