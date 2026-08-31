package e2e_test

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/socrates-devtool/socrates/internal/reviewer"
	"github.com/socrates-devtool/socrates/internal/runstate"
	"github.com/socrates-devtool/socrates/internal/testing/harness"
)

const (
	childOfTheReviewer = "child the reviewer spawned"
	theReviewer        = "reviewer"
)

func TestAChildThatOutlivesTheAgentIsKilled(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario("default:\n  child_command: [sleep, \"120\"]\n")
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	run := session.WaitForTerminal(name)
	if run.State != runstate.Clean {
		t.Fatalf("the run ended as %s: %s", run.State, reason(run))
	}

	waitUntilReaped(t, session.Simulator.OnlyInvocation().Child, childOfTheReviewer)
}

func TestReapingHappensOnErrorToo(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario("default:\n  behaviour: no-result\n  child_command: [sleep, \"120\"]\n")
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	run := session.WaitForTerminal(name)
	if run.State != runstate.Failed {
		t.Fatalf("the run ended as %s, want %s", run.State, runstate.Failed)
	}

	waitUntilReaped(t, session.Simulator.OnlyInvocation().Child, childOfTheReviewer)
}

func TestReapingHappensOnCancellation(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario("default:\n  behaviour: silence\n  hold: 5m\n  child_command: [sleep, \"120\"]\n")
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	session.WaitFor(name, runstate.Reviewing)

	child := waitForChild(t, session)
	waitUntilTheRunNamesItsAgent(t, session, name)

	stopped, err := reviewer.Cancel(context.Background(), session.Store, name)
	if err != nil {
		t.Fatalf("cancelling the run: %v", err)
	}
	if stopped.State != runstate.Cancelled {
		t.Errorf("the cancelled run reads as %s", stopped.State)
	}

	waitUntilReaped(t, child, childOfTheReviewer)
	waitUntilReaped(t, stopped.Reviewer, theReviewer)
}

func TestTheWorktreeCanBeRemovedAfterTheRound(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario("default:\n  child_command: [sleep, \"120\"]\n")
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	run := session.WaitForTerminal(name)

	worktree := session.Recorded(name).Worktree
	if worktree == "" {
		t.Fatal("the run records no worktree")
	}

	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(worktree); os.IsNotExist(err) {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}

	if _, err := os.Stat(worktree); err == nil {
		t.Errorf("%s is still on disk after the round ended as %s", worktree, run.State)
	}
	if listed := session.Repository.Worktrees(); strings.Contains(listed, worktree) {
		t.Errorf("git worktree list still names %s:\n%s", worktree, listed)
	}
}

func waitForChild(t *testing.T, session *harness.Session) int {
	t.Helper()

	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		invocations := session.Simulator.Invocations()
		if len(invocations) == 1 && invocations[0].Child != 0 {
			return invocations[0].Child
		}
		time.Sleep(20 * time.Millisecond)
	}

	t.Fatal("the reviewer never spawned the child this test needs")
	return 0
}

func waitUntilTheRunNamesItsAgent(t *testing.T, session *harness.Session, name string) {
	t.Helper()

	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		if session.Recorded(name).Agent != 0 {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}

	t.Fatalf("the run %s never named its agent, so cancelling it reaches no process", name)
}

func waitUntilReaped(t *testing.T, pid int, named string) {
	t.Helper()

	if pid == 0 {
		t.Fatalf("no %s was recorded, so this test reaches no process", named)
	}
	if !waitUntilGone(pid) {
		t.Errorf("the %s, %d, is still running", named, pid)
	}
}
