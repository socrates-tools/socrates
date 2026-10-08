package e2e_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/testing/harness"
	"github.com/socrates-tools/socrates/review/corpus"
	"github.com/socrates-tools/socrates/review/corpus/runs"
	"github.com/socrates-tools/socrates/review/runstate"
)

func TestAnAnswerIsRecordedFromTheRunThatProducedIt(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	session.Simulator.Scenario(escalatingScenario)
	session.Repository.Write("go.mod", "module "+corpus.ModulePath+"\n\ngo 1.24\n")
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Begin()
	run := session.WaitForTerminal(name)
	if run.State != runstate.Escalated {
		t.Fatalf("the review ended as %s: %s", run.State, reason(run))
	}

	answered, err := runs.Answered(session.Store, runs.Request{
		Run:     name,
		Round:   1,
		Name:    "one-reviewed-change",
		Subject: "one added file in a repository that exists for this test",
		On:      "2026-08-21",
	})
	if err != nil {
		t.Fatalf("taking the answer out of the run: %v", err)
	}

	into := t.TempDir()
	if _, err := corpus.Record(into, answered.Repository, answered.Entry); err != nil {
		t.Fatalf("recording the answer: %v", err)
	}

	kept, err := corpus.Read(filepath.Join(into, corpus.Directory))
	if err != nil {
		t.Fatalf("reading the recorded answer back: %v", err)
	}
	if len(kept) != 1 {
		t.Fatalf("%d answers were kept, want one", len(kept))
	}

	entry := kept[0]
	stored := strings.TrimSuffix(session.File(session.Round(name, 1).FindingsPath()), "\n")
	if string(entry.Document) != stored {
		t.Errorf("the kept document is not what the reviewer answered:\nanswered %q\nkept     %q", stored, entry.Document)
	}

	invocation := session.Simulator.OnlyInvocation()
	if entry.Model != invocation.Model || entry.Effort != invocation.Effort {
		t.Errorf("the answer was kept as %s at %s effort, and the reviewer was asked as %s at %s", entry.Model, entry.Effort, invocation.Model, invocation.Effort)
	}
	if entry.Round != 1 || entry.Run != name {
		t.Errorf("the answer was kept as round %d of %s, want round 1 of %s", entry.Round, entry.Run, name)
	}
	if entry.Base != run.Base || entry.Snapshot != run.Snapshot {
		t.Errorf("the answer was kept against %s and %s, and the round reviewed %s and %s", entry.Base, entry.Snapshot, run.Base, run.Snapshot)
	}
}
