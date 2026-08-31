package runstore

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-devtool/socrates/internal/runstate"
)

func TestATerminalRunDoesNotMove(t *testing.T) {
	for _, outcome := range []runstate.State{runstate.Clean, runstate.Escalated, runstate.Failed, runstate.Cancelled} {
		t.Run(string(outcome), func(t *testing.T) {
			store := New(t.TempDir())
			run := started(t, store)

			if err := run.Advance(runstate.Reviewing); err != nil {
				t.Fatalf("advancing to %s: %v", runstate.Reviewing, err)
			}
			if err := run.End(outcome, "the reviewer said so"); err != nil {
				t.Fatalf("ending as %s: %v", outcome, err)
			}
			if err := store.Save(run); err != nil {
				t.Fatalf("saving the ended run: %v", err)
			}

			if err := run.Advance(runstate.Reviewing); err == nil {
				t.Errorf("a run that ended as %s went back to %s", outcome, runstate.Reviewing)
			}
			if err := run.End(runstate.Clean, "a second opinion"); err == nil {
				t.Errorf("a run that ended as %s was ended again as %s", outcome, runstate.Clean)
			}

			moved := *run
			moved.State = runstate.Reviewing
			moved.Verdict = nil
			if err := store.Save(&moved); err == nil {
				t.Errorf("a run that ended as %s was written back as %s", outcome, runstate.Reviewing)
			}

			loaded, err := store.Load(run.Name)
			if err != nil {
				t.Fatalf("reading the run: %v", err)
			}
			if loaded.State != outcome {
				t.Errorf("the run reads as %q, want %q", loaded.State, outcome)
			}
			if loaded.Verdict == nil || loaded.Verdict.Outcome != outcome {
				t.Errorf("the run's verdict is %+v, want one naming %q", loaded.Verdict, outcome)
			}

			if err := store.Save(run); err != nil {
				t.Errorf("writing the same terminal state again was refused: %v", err)
			}
		})
	}
}

func TestAFailedRunReportsNoFindingsAsItsVerdict(t *testing.T) {
	store := New(t.TempDir())
	run := started(t, store)

	if err := run.End(runstate.Failed, ""); err == nil {
		t.Error("a run failed without saying why, which reads as a review that found nothing")
	}

	const reason = "the reviewer went away before it reported anything"
	if err := run.End(runstate.Failed, reason); err != nil {
		t.Fatalf("failing the run: %v", err)
	}
	if err := store.Save(run); err != nil {
		t.Fatalf("saving the failed run: %v", err)
	}

	loaded, err := store.Load(run.Name)
	if err != nil {
		t.Fatalf("reading the run: %v", err)
	}
	if loaded.Verdict == nil {
		t.Fatal("a failed run has no verdict, so nothing states that no review happened")
	}
	if loaded.Verdict.Outcome != runstate.Failed {
		t.Errorf("the verdict names %q, want %q", loaded.Verdict.Outcome, runstate.Failed)
	}
	if loaded.Verdict.Reason != reason {
		t.Errorf("the verdict reads %q, want %q", loaded.Verdict.Reason, reason)
	}

	content, err := os.ReadFile(filepath.Join(store.Path(run.Name), recordName))
	if err != nil {
		t.Fatalf("reading %s: %v", recordName, err)
	}
	var record map[string]any
	if err := json.Unmarshal(content, &record); err != nil {
		t.Fatalf("parsing %s: %v", recordName, err)
	}
	for key := range record {
		if strings.Contains(strings.ToLower(key), "finding") {
			t.Errorf("a failed run records %q, which presents an absence of findings as a result", key)
		}
	}
}

func TestAdvanceAndEndRefuseEachOthersStates(t *testing.T) {
	store := New(t.TempDir())
	run := started(t, store)

	if err := run.Advance(runstate.Clean); err == nil {
		t.Errorf("a terminal state was advanced to without a verdict")
	}
	if err := run.End(runstate.Reviewing, "because"); err == nil {
		t.Errorf("a state that is not terminal was recorded as a verdict")
	}
	if err := run.Advance(runstate.State("reviewed")); err == nil {
		t.Errorf("a state socrates does not have was accepted")
	}
	if run.State != runstate.Starting {
		t.Errorf("a refused transition moved the run to %q", run.State)
	}
}
