package runstore_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rui-ktei/socrates/internal/runstate"
	"github.com/rui-ktei/socrates/internal/runstore"
	"github.com/rui-ktei/socrates/internal/snapshot"
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

func TestARunWhoseWorktreeIsStillSealedIsSwept(t *testing.T) {
	store := runstore.New(t.TempDir())
	now := time.Date(2026, 8, 19, 12, 0, 0, 0, time.UTC)

	abandoned := start(t, store, t.TempDir(), now.Add(-20*24*time.Hour))
	if err := abandoned.End(runstate.Failed, "the reviewer died and nothing lifted the seal"); err != nil {
		t.Fatalf("ending the abandoned run: %v", err)
	}
	if err := store.Save(abandoned); err != nil {
		t.Fatalf("saving the abandoned run: %v", err)
	}

	worktree := snapshot.WorktreePath(store.Path(abandoned.Name))
	if err := os.MkdirAll(filepath.Join(worktree, "internal"), 0o755); err != nil {
		t.Fatalf("standing in for a worktree nothing unsealed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(worktree, "internal", "added.go"), []byte("package internal\n"), 0o644); err != nil {
		t.Fatalf("filling the worktree nothing unsealed: %v", err)
	}
	if err := snapshot.Seal(worktree); err != nil {
		t.Fatalf("sealing the worktree nothing unsealed: %v", err)
	}

	if err := store.Sweep(now, 14*24*time.Hour); err != nil {
		t.Fatalf("sweeping past a sealed worktree: %v", err)
	}
	if _, err := os.Stat(store.Path(abandoned.Name)); !os.IsNotExist(err) {
		t.Errorf("the run %s was left on disk by the sweep, so one sealed worktree stops every later review", abandoned.Name)
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
