package runstore

import (
	"os/exec"
	"testing"

	"github.com/rui-ktei/socrates/internal/runstate"
)

func TestAParkedRunWithNoProcessesIsStillParked(t *testing.T) {
	store := New(t.TempDir())
	run := started(t, store)

	if err := run.Advance(runstate.Reviewing); err != nil {
		t.Fatalf("advancing to %s: %v", runstate.Reviewing, err)
	}
	run.Reviewer = reapedProcess(t)
	if err := store.Save(run); err != nil {
		t.Fatalf("saving the reviewing run: %v", err)
	}

	reviewing, err := store.Load(run.Name)
	if err != nil {
		t.Fatalf("reading the reviewing run: %v", err)
	}
	if reviewing.State != runstate.Failed {
		t.Fatalf("a run recorded as %s with no live reviewer reads as %q, so the liveness rule is not what this exempts from", runstate.Reviewing, reviewing.State)
	}

	if err := run.Advance(runstate.AwaitingResponse); err != nil {
		t.Fatalf("parking the run: %v", err)
	}
	if err := store.Save(run); err != nil {
		t.Fatalf("saving the parked run: %v", err)
	}

	parked, err := store.Load(run.Name)
	if err != nil {
		t.Fatalf("reading the parked run: %v", err)
	}
	if parked.State != runstate.AwaitingResponse {
		t.Errorf("a parked run with no process of its own reads as %q, want %q", parked.State, runstate.AwaitingResponse)
	}
	if parked.Verdict != nil {
		t.Errorf("a parked run carries the verdict %+v", *parked.Verdict)
	}
}

func reapedProcess(t *testing.T) int {
	t.Helper()

	command := exec.Command("sh", "-c", "exit 0")
	if err := command.Start(); err != nil {
		t.Fatalf("starting a process to reap: %v", err)
	}

	pid := command.Process.Pid
	if err := command.Wait(); err != nil {
		t.Fatalf("waiting for the process to be reaped: %v", err)
	}
	return pid
}
