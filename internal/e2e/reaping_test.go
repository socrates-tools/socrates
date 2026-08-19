package e2e_test

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/rui-ktei/socrates/internal/process"
	"github.com/rui-ktei/socrates/internal/reviewer"
	"github.com/rui-ktei/socrates/internal/runstate"
	"github.com/rui-ktei/socrates/internal/testing/harness"
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

	waitUntilReaped(t, session.Simulator.OnlyInvocation().Child)
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

	waitUntilReaped(t, session.Simulator.OnlyInvocation().Child)
}

func TestReapingHappensOnCancellation(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario("default:\n  behaviour: silence\n  hold: 5m\n  child_command: [sleep, \"120\"]\n")
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	session.WaitFor(name, runstate.Reviewing)

	child := waitForChild(t, session)

	stopped, err := reviewer.Cancel(context.Background(), session.Store, name)
	if err != nil {
		t.Fatalf("cancelling the run: %v", err)
	}
	if stopped.State != runstate.Cancelled {
		t.Errorf("the cancelled run reads as %s", stopped.State)
	}

	waitUntilReaped(t, child)
	if process.Alive(stopped.Reviewer) {
		t.Errorf("the reviewer %d is still running after cancellation", stopped.Reviewer)
	}
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

func waitUntilReaped(t *testing.T, pid int) {
	t.Helper()

	if pid == 0 {
		t.Fatal("the simulator recorded no child, so this test reaches no process group")
	}

	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		if !process.Alive(pid) {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Errorf("the child %d the reviewer spawned is still running", pid)
}
