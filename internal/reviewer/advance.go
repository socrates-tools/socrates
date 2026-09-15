package reviewer

import (
	"context"
	"fmt"
	"strings"

	"github.com/socrates-tools/socrates/internal/refusal"
	"github.com/socrates-tools/socrates/internal/runlog"
	"github.com/socrates-tools/socrates/internal/runstate"
	"github.com/socrates-tools/socrates/internal/runstore"
	"github.com/socrates-tools/socrates/internal/snapshot"
)

type Continuation struct {
	Store  runstore.Store
	Name   string
	Binary string
}

type Continued struct {
	Run     *runstore.Run
	Round   int
	Started bool
}

func Continue(ctx context.Context, continuation Continuation) (Continued, error) {
	if outer := OuterRun(); outer != "" {
		return Continued{}, refusal.New(aReviewInsideAReview, outer)
	}

	store := continuation.Store

	run, round, document, err := parked(store, continuation.Name)
	if err != nil {
		return Continued{}, err
	}

	recorded, err := round.Responses()
	if err != nil {
		return Continued{}, err
	}
	next := round.Number + 1
	continued := Continued{Run: run, Round: next}

	if waiting := recorded.Unanswered(document, run.Authorship()); len(waiting) > 0 {
		if run.NotMyChange {
			return Continued{}, refusal.New(stillUnansweredForSomebodyElse, len(waiting), strings.Join(waiting, ", "))
		}
		if asked := recorded.AwaitingAPerson(document, run.Authorship()); len(asked) > 0 {
			return Continued{}, refusal.New(stillUnansweredByAPerson, len(waiting), strings.Join(waiting, ", "), len(asked), strings.Join(asked, ", "))
		}
		return Continued{}, refusal.New(stillUnanswered, len(waiting), strings.Join(waiting, ", "))
	}

	if run.NotMyChange {
		return continued, conclude(store, run, runstate.Escalated, theFindingsWereHandedOver)
	}

	taken, err := snapshot.TakeAgainst(ctx, run.Repository, run.Base)
	if err != nil {
		return continued, fail(store, run, err)
	}
	if taken.NothingToReview() {
		return continued, conclude(store, run, runstate.Failed, nothingLeftToReview)
	}

	if _, err := store.CreateRound(run.Name, next); err != nil {
		return continued, fail(store, run, err)
	}

	run.Snapshot = taken.Commit
	run.Snapshotted(next, taken.Base, taken.Commit)
	run.Reviewer = 0
	run.Agent = 0
	if err := run.Advance(runstate.Starting); err != nil {
		return continued, err
	}
	if err := store.Save(run); err != nil {
		return continued, err
	}

	if err := record(store, run, runlog.Event{
		Kind:    EventRoundStarted,
		Subject: run.Name,
		Fields:  map[string]any{"round": next, "base": run.Base, "snapshot": run.Snapshot},
	}); err != nil {
		return continued, err
	}

	directory := store.Path(run.Name)
	path := snapshot.WorktreePath(directory)
	leftover := snapshot.Worktree{RepositoryRoot: run.Repository, Path: path, Commit: run.Snapshot}
	if err := leftover.Remove(ctx); err != nil {
		return continued, fail(store, run, fmt.Errorf(clearingTheLastRound, err))
	}

	worktree, err := taken.Checkout(ctx, path)
	if err != nil {
		return continued, fail(store, run, fmt.Errorf(checkingOutTheSnapshot, err))
	}

	run.Worktree = worktree.Path
	if err := store.Save(run); err != nil {
		return continued, err
	}

	if err := spawnSupervisor(continuation.Binary, directory, run.Name); err != nil {
		return continued, fail(store, run, err)
	}

	continued.Started = true
	return continued, nil
}
