package reviewer

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/rui-ktei/socrates/internal/process"
	"github.com/rui-ktei/socrates/internal/refusal"
	"github.com/rui-ktei/socrates/internal/runstate"
	"github.com/rui-ktei/socrates/internal/runstore"
	"github.com/rui-ktei/socrates/internal/snapshot"
)

const (
	RunMarker         = "SOCRATES_RUN"
	SuperviseCommand  = "supervise"
	supervisorLogName = "supervisor.log"
)

type Beginning struct {
	Directory string
	Intent    string
	Store     runstore.Store
	Binary    string
	StartedAt time.Time
}

func OuterRun() string {
	return os.Getenv(RunMarker)
}

func Begin(ctx context.Context, beginning Beginning) (*runstore.Run, error) {
	if outer := OuterRun(); outer != "" {
		return nil, refusal.New("this process is already the reviewer for run %s, and a review inside a review spawns without end", outer)
	}

	taken, err := snapshot.Take(ctx, beginning.Directory)
	if err != nil {
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

	if taken.NothingToReview() {
		if err := run.End(runstate.Clean, snapshot.NothingToReviewReason); err != nil {
			return nil, err
		}
		if err := beginning.Store.Save(run); err != nil {
			return nil, err
		}
		return run, nil
	}

	directory := beginning.Store.Path(run.Name)

	worktree, err := taken.Checkout(ctx, snapshot.WorktreePath(directory))
	if err != nil {
		return run, beginning.fail(run, fmt.Errorf("checking out the snapshot: %w", err))
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
	return reason
}

func spawnSupervisor(binary, directory, name string) error {
	log, err := os.Create(filepath.Join(directory, supervisorLogName))
	if err != nil {
		return fmt.Errorf("opening the supervisor's own log: %w", err)
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
		return fmt.Errorf("starting the reviewer: %w", err)
	}
	return nil
}
