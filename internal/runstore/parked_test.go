package runstore

import (
	"os/exec"
	"runtime"
	"testing"

	"github.com/socrates-tools/socrates/internal/runstate"
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

func TestARunParkedAsItsReviewerExitedIsNotAFailure(t *testing.T) {
	interleavedOnOneProcessor(t)
	store := New(t.TempDir())

	for attempt := 0; attempt < 100; attempt++ {
		run := started(t, store)
		if err := run.Advance(runstate.Reviewing); err != nil {
			t.Fatalf("advancing to %s: %v", runstate.Reviewing, err)
		}

		reviewer := exec.Command("sh", "-c", "sleep 30")
		if err := reviewer.Start(); err != nil {
			t.Fatalf("starting a reviewer to outlive the record: %v", err)
		}
		run.Reviewer = reviewer.Process.Pid
		if err := store.Save(run); err != nil {
			t.Fatalf("saving the reviewing run: %v", err)
		}

		parked := *run
		if err := parked.Advance(runstate.AwaitingResponse); err != nil {
			t.Fatalf("parking the run: %v", err)
		}

		ready := make(chan struct{})
		saved := make(chan error, 1)
		go func() {
			<-ready
			saved <- store.Save(&parked)
			_ = reviewer.Process.Kill()
			_, _ = reviewer.Process.Wait()
		}()

		close(ready)
		for {
			loaded, err := store.Load(run.Name)
			if err != nil {
				t.Fatalf("reading the run: %v", err)
			}
			if loaded.State == runstate.AwaitingResponse {
				break
			}
			if loaded.State != runstate.Reviewing {
				t.Fatalf("a run parked as its reviewer exited reads as %q, and a reader saw the record from before it parked: %+v", loaded.State, loaded.Verdict)
			}
		}
		if err := <-saved; err != nil {
			t.Fatalf("parking the run on disk: %v", err)
		}
	}
}

func interleavedOnOneProcessor(t *testing.T) {
	t.Helper()

	previous := runtime.GOMAXPROCS(1)
	t.Cleanup(func() { runtime.GOMAXPROCS(previous) })
}
