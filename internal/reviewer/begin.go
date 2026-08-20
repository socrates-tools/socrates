package reviewer

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/rui-ktei/socrates/internal/process"
	"github.com/rui-ktei/socrates/internal/refusal"
	"github.com/rui-ktei/socrates/internal/runlog"
	"github.com/rui-ktei/socrates/internal/runstate"
	"github.com/rui-ktei/socrates/internal/runstore"
	"github.com/rui-ktei/socrates/internal/snapshot"
)

const (
	RunMarker         = "SOCRATES_RUN"
	SuperviseCommand  = "supervise"
	supervisorLogName = "supervisor.log"
)

const (
	EventRunStarted = "run-started"
	EventRunEnded   = "run-ended"
)

type Beginning struct {
	Directory string
	Intent    string
	Store     runstore.Store
	Binary    string
	StartedAt time.Time
	Retention time.Duration
}

func OuterRun() string {
	return os.Getenv(RunMarker)
}

func Begin(ctx context.Context, beginning Beginning) (*runstore.Run, error) {
	if outer := OuterRun(); outer != "" {
		return nil, refusal.New(aReviewInsideAReview, outer)
	}

	taken, err := snapshot.Take(ctx, beginning.Directory)
	if err != nil {
		return nil, err
	}

	if err := beginning.Store.Sweep(beginning.StartedAt, beginning.Retention); err != nil {
		return nil, err
	}

	run, err := beginning.Store.Create(runstore.Start{
		StartedAt:      beginning.StartedAt,
		RepositoryRoot: taken.RepositoryRoot,
		Base:           taken.Base,
		Snapshot:       taken.Commit,
		Intent:         beginning.Intent,
	})
	if err != nil {
		return nil, err
	}

	directory := beginning.Store.Path(run.Name)
	if err := beginning.record(run, runlog.Event{
		Kind:    EventRunStarted,
		Subject: run.Name,
		Fields: map[string]any{
			"repository": run.Repository,
			"base":       run.Base,
			"snapshot":   run.Snapshot,
		},
	}); err != nil {
		return nil, err
	}

	if taken.NothingToReview() {
		if err := run.End(runstate.Clean, snapshot.NothingToReviewReason); err != nil {
			return nil, err
		}
		if err := beginning.Store.Save(run); err != nil {
			return nil, err
		}
		return run, beginning.recordEnd(run)
	}

	if _, err := beginning.Store.CreateRound(run.Name, firstRound); err != nil {
		return run, beginning.fail(run, err)
	}

	worktree, err := taken.Checkout(ctx, snapshot.WorktreePath(directory))
	if err != nil {
		if snapshot.SealNotEnforced(err) {
			return nil, errors.Join(err, beginning.Store.Discard(run.Name))
		}
		return run, beginning.fail(run, fmt.Errorf(checkingOutTheSnapshot, err))
	}
	run.Worktree = worktree.Path
	if err := beginning.Store.Save(run); err != nil {
		return run, err
	}

	if err := spawnSupervisor(beginning.Binary, directory, run.Name); err != nil {
		return run, beginning.fail(run, err)
	}

	return run, nil
}

func (b Beginning) fail(run *runstore.Run, reason error) error {
	if err := run.End(runstate.Failed, reason.Error()); err != nil {
		return err
	}
	if err := b.Store.Save(run); err != nil {
		return err
	}
	if err := liftTheSeal(run); err != nil {
		return errors.Join(reason, err)
	}
	if err := b.recordEnd(run); err != nil {
		return err
	}
	return reason
}

func (b Beginning) recordEnd(run *runstore.Run) error {
	return b.record(run, runlog.Event{
		Kind:    EventRunEnded,
		Subject: run.Name,
		Fields:  map[string]any{"state": string(run.State), "reason": run.Verdict.Reason},
	})
}

func (b Beginning) record(run *runstore.Run, event runlog.Event) error {
	events, err := runlog.Open(filepath.Join(b.Store.Path(run.Name), runlog.FileName), 0)
	if err != nil {
		return err
	}
	defer events.Close()
	return events.Record(event)
}

func spawnSupervisor(binary, directory, name string) error {
	log, err := os.Create(filepath.Join(directory, supervisorLogName))
	if err != nil {
		return fmt.Errorf(openingTheSupervisorLog, err)
	}
	defer log.Close()

	command := exec.Command(binary, SuperviseCommand, name)
	command.Dir = directory
	command.Env = append(os.Environ(), RunMarker+"="+name)
	command.Stdin = nil
	command.Stdout = log
	command.Stderr = log
	process.Detach(command)

	if err := command.Start(); err != nil {
		return fmt.Errorf(startingTheReviewer, err)
	}
	return nil
}
