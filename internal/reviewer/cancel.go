package reviewer

import (
	"context"

	"github.com/socrates-tools/socrates-review/internal/process"
	"github.com/socrates-tools/socrates-review/internal/runstate"
	"github.com/socrates-tools/socrates-review/internal/runstore"
	"github.com/socrates-tools/socrates-review/internal/snapshot"
)

func Cancel(ctx context.Context, store runstore.Store, name string) (*runstore.Run, error) {
	run, err := store.Load(name)
	if err != nil {
		return nil, err
	}
	if run.State.Terminal() {
		if err := store.Save(run); err != nil {
			return nil, err
		}
		return run, liftTheSeal(run)
	}

	if err := run.End(runstate.Cancelled, cancelledReason); err != nil {
		return nil, err
	}
	if err := store.Save(run); err != nil {
		return nil, err
	}

	if err := process.KillGroup(run.Agent); err != nil {
		return nil, err
	}
	if err := process.Kill(run.Reviewer); err != nil {
		return nil, err
	}

	if run.Worktree != "" {
		worktree := snapshot.Worktree{RepositoryRoot: run.Repository, Path: run.Worktree, Commit: run.Snapshot}
		if err := worktree.Remove(ctx); err != nil {
			return run, err
		}
	}

	return run, nil
}
