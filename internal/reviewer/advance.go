package reviewer

import (
	"context"
	"fmt"
	"strings"

	"github.com/socrates-tools/socrates-review/internal/refusal"
	"github.com/socrates-tools/socrates-review/internal/runlog"
	"github.com/socrates-tools/socrates-review/internal/runstate"
	"github.com/socrates-tools/socrates-review/internal/runstore"
	"github.com/socrates-tools/socrates-review/internal/snapshot"
)

type Continuation struct {
	Store  runstore.Store
	Name   string
	Binary string
}

func Continue(ctx context.Context, continuation Continuation) (*runstore.Run, int, error) {
	if outer := OuterRun(); outer != "" {
		return nil, 0, refusal.New(aReviewInsideAReview, outer)
	}

	store := continuation.Store

	run, round, document, err := parked(store, continuation.Name)
	if err != nil {
		return nil, 0, err
	}

	recorded, err := round.Responses()
	if err != nil {
		return nil, 0, err
	}
	if waiting := recorded.Unanswered(document); len(waiting) > 0 {
		if asked := recorded.AwaitingAPerson(document); len(asked) > 0 {
			return nil, 0, refusal.New(stillUnansweredByAPerson, len(waiting), strings.Join(waiting, ", "), len(asked), strings.Join(asked, ", "))
		}
		return nil, 0, refusal.New(stillUnanswered, len(waiting), strings.Join(waiting, ", "))
	}

	next := round.Number + 1

	taken, err := snapshot.TakeAgainst(ctx, run.Repository, run.Base)
	if err != nil {
		return run, next, fail(store, run, err)
	}
	if taken.NothingToReview() {
		return run, next, conclude(store, run, runstate.Clean, snapshot.NothingToReviewReason)
	}

	if _, err := store.CreateRound(run.Name, next); err != nil {
		return run, next, fail(store, run, err)
	}

	run.Snapshot = taken.Commit
	run.Snapshotted(next, taken.Base, taken.Commit)
	run.Reviewer = 0
	run.Agent = 0
	if err := run.Advance(runstate.Starting); err != nil {
		return run, next, err
	}
	if err := store.Save(run); err != nil {
		return run, next, err
	}

	if err := record(store, run, runlog.Event{
		Kind:    EventRoundStarted,
		Subject: run.Name,
		Fields:  map[string]any{"round": next, "base": run.Base, "snapshot": run.Snapshot},
	}); err != nil {
		return run, next, err
	}

	directory := store.Path(run.Name)
	path := snapshot.WorktreePath(directory)
	leftover := snapshot.Worktree{RepositoryRoot: run.Repository, Path: path, Commit: run.Snapshot}
	if err := leftover.Remove(ctx); err != nil {
		return run, next, fail(store, run, fmt.Errorf(clearingTheLastRound, err))
	}

	worktree, err := taken.Checkout(ctx, path)
	if err != nil {
		return run, next, fail(store, run, fmt.Errorf(checkingOutTheSnapshot, err))
	}

	run.Worktree = worktree.Path
	if err := store.Save(run); err != nil {
		return run, next, err
	}

	if err := spawnSupervisor(continuation.Binary, directory, run.Name); err != nil {
		return run, next, fail(store, run, err)
	}
	return run, next, nil
}
