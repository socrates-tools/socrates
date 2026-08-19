package runstore_test

import (
	"os"
	"testing"
	"time"

	"github.com/rui-ktei/socrates/internal/runstate"
	"github.com/rui-ktei/socrates/internal/runstore"
)

func TestAnOldRunIsSwept(t *testing.T) {
	store := runstore.New(t.TempDir())
	repository := t.TempDir()
	now := time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC)

	old := start(t, store, repository, now.Add(-20*24*time.Hour))
	if err := old.End(runstate.Clean, "a reviewer read the change and found nothing to report"); err != nil {
		t.Fatalf("ending the old run: %v", err)
	}
	if err := store.Save(old); err != nil {
		t.Fatalf("saving the old run: %v", err)
	}

	recent := start(t, store, repository, now.Add(-time.Hour))
	if err := recent.End(runstate.Escalated, "a reviewer reported something"); err != nil {
		t.Fatalf("ending the recent run: %v", err)
	}
	if err := store.Save(recent); err != nil {
		t.Fatalf("saving the recent run: %v", err)
	}

	if err := store.Sweep(now, 14*24*time.Hour); err != nil {
		t.Fatalf("sweeping: %v", err)
	}

	if _, err := os.Stat(store.Path(old.Name)); !os.IsNotExist(err) {
		t.Errorf("the run %s is older than the retention age and its directory remains", old.Name)
	}
	if _, err := os.Stat(store.Path(recent.Name)); err != nil {
		t.Errorf("the run %s is within the retention age and its directory is gone: %v", recent.Name, err)
	}
}

func TestALiveRunIsNeverSwept(t *testing.T) {
	store := runstore.New(t.TempDir())
	now := time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC)

	working := start(t, store, t.TempDir(), now.Add(-20*24*time.Hour))
	if err := working.Advance(runstate.Reviewing); err != nil {
		t.Fatalf("advancing the old run: %v", err)
	}
	if err := store.Save(working); err != nil {
		t.Fatalf("saving the old run: %v", err)
	}

	if err := store.Sweep(now, 14*24*time.Hour); err != nil {
		t.Fatalf("sweeping: %v", err)
	}

	if _, err := os.Stat(store.Path(working.Name)); err != nil {
		t.Errorf("a run that is still %s was swept: %v", runstate.Reviewing, err)
	}
}

func start(t *testing.T, store runstore.Store, repository string, at time.Time) *runstore.Run {
	t.Helper()

	run, err := store.Create(runstore.Start{
		StartedAt:      at,
		RepositoryRoot: repository,
		Base:           "0000000000000000000000000000000000000000",
		Snapshot:       "1111111111111111111111111111111111111111",
	})
	if err != nil {
		t.Fatalf("creating a run started at %s: %v", at, err)
	}
	return run
}
