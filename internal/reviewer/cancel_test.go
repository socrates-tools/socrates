package reviewer_test

import (
	"context"
	"os/exec"
	"runtime"
	"testing"
	"time"

	"github.com/socrates-devtool/socrates/internal/reviewer"
	"github.com/socrates-devtool/socrates/internal/runstate"
	"github.com/socrates-devtool/socrates/internal/runstore"
)

func TestACancelledRunIsNeverReadAsFailed(t *testing.T) {
	// scenario: Cancelling while following
	interleavedOnOneProcessor(t)
	store := runstore.New(t.TempDir())

	for attempt := 0; attempt < 100; attempt++ {
		run := underReview(t, store)

		cancelled := make(chan error, 1)
		ready := make(chan struct{})
		go func() {
			<-ready
			_, err := reviewer.Cancel(context.Background(), store, run.Name)
			cancelled <- err
		}()

		close(ready)
		for {
			loaded, err := store.Load(run.Name)
			if err != nil {
				t.Fatalf("reading the run: %v", err)
			}
			if loaded.State == runstate.Cancelled {
				break
			}
			if loaded.State != runstate.Reviewing {
				t.Fatalf("a run being cancelled reads as %q, so a reader saw it between its reviewer dying and the record saying why: %+v", loaded.State, loaded.Verdict)
			}
		}
		if err := <-cancelled; err != nil {
			t.Fatalf("cancelling the run: %v", err)
		}
	}
}

func underReview(t *testing.T, store runstore.Store) *runstore.Run {
	t.Helper()

	run, err := store.Create(runstore.Start{
		StartedAt:      time.Now(),
		RepositoryRoot: t.TempDir(),
		Base:           "base-commit",
		Snapshot:       "snapshot-commit",
		Intent:         "what the person asked for",
	})
	if err != nil {
		t.Fatalf("creating the run: %v", err)
	}
	if err := run.Advance(runstate.Reviewing); err != nil {
		t.Fatalf("advancing to %s: %v", runstate.Reviewing, err)
	}

	reviewing := exec.Command("sh", "-c", "sleep 30")
	if err := reviewing.Start(); err != nil {
		t.Fatalf("starting a reviewer to cancel: %v", err)
	}
	go func() { _ = reviewing.Wait() }()
	t.Cleanup(func() { _ = reviewing.Process.Kill() })

	run.Reviewer = reviewing.Process.Pid
	if err := store.Save(run); err != nil {
		t.Fatalf("saving the reviewing run: %v", err)
	}
	return run
}

func interleavedOnOneProcessor(t *testing.T) {
	t.Helper()

	previous := runtime.GOMAXPROCS(1)
	t.Cleanup(func() { runtime.GOMAXPROCS(previous) })
}
