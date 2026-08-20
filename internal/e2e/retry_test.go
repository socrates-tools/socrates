package e2e_test

import (
	"os"
	"testing"
	"time"

	"github.com/rui-ktei/socrates/internal/agent"
	"github.com/rui-ktei/socrates/internal/claudesim"
	"github.com/rui-ktei/socrates/internal/process"
	"github.com/rui-ktei/socrates/internal/runstate"
	"github.com/rui-ktei/socrates/internal/testing/harness"
)

const refusedStrictnessScenario = "default:\n  behaviour: strict-tool-servers-refused\n  events: 60\n  interval: 100ms\n"

func TestCancellingStopsEverythingAfterASecondInvocation(t *testing.T) {
	// scenario: Cancelling stops everything
	session := harness.Start(t)
	session.Simulator.Scenario(refusedStrictnessScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	session.WaitFor(name, runstate.Reviewing)
	answering := waitForTheSecondInvocation(t, session)

	started := session.Recorded(name)
	if started.Agent != answering.Process {
		t.Errorf("the run records the process %d and the reviewer that is answering is %d", started.Agent, answering.Process)
	}

	if _, failures, code := session.Scrat("cancel", name); code != 0 {
		t.Fatalf("scrat cancel exited %d: %s", code, failures)
	}

	run := session.Recorded(name)
	if run.State != runstate.Cancelled {
		t.Fatalf("the run is %s after being cancelled", run.State)
	}
	if _, err := os.Stat(started.Worktree); !os.IsNotExist(err) {
		t.Errorf("the worktree at %s survived the cancellation", started.Worktree)
	}

	for _, invocation := range session.Simulator.Invocations() {
		waitForDeath(t, name, invocation.Process, 0)
	}
	waitForDeath(t, name, started.Reviewer, 0)
}

func TestTheReviewIsRunAgainWithoutWhatTheCLIRefused(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario("default:\n  behaviour: strict-tool-servers-refused\n")
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	run := session.WaitForTerminal(name)
	if run.State != runstate.Clean {
		t.Fatalf("a CLI that refused to ignore the machine's tool servers ended the run as %s: %s", run.State, reason(run))
	}

	invocations := session.Simulator.Invocations()
	if len(invocations) != 2 {
		t.Fatalf("the reviewer was invoked %d times, want the refused one and the one that ran", len(invocations))
	}
	if !invocations[0].StrictMCP || invocations[1].StrictMCP {
		t.Errorf("the invocations asked for strictness %v and %v", invocations[0].StrictMCP, invocations[1].StrictMCP)
	}

	loaded := 0
	for _, event := range session.Events(name) {
		if event.Kind == agent.EventToolServersLoaded {
			loaded++
		}
	}
	if loaded != 1 {
		t.Errorf("%d events in the run's log record that the machine's tool servers were loaded, want one", loaded)
	}
}

func TestAToolServerThatWritesCannotReachTheChangeUnderReview(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario("default:\n  behaviour: strict-tool-servers-refused\n  writes:\n    - path: .serena/project.yml\n      content: \"a tool server writing its own bookkeeping as it starts\\n\"\n")
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	run := session.WaitForTerminal(name)
	if run.State != runstate.Clean {
		t.Fatalf("a reviewer whose session wrote as it started ended the run as %s: %s", run.State, reason(run))
	}
	if run.Verdict == nil || run.Verdict.Outcome != runstate.Clean {
		t.Errorf("the verdict is %+v, want the one the findings imply", run.Verdict)
	}

	invocations := session.Simulator.Invocations()
	if len(invocations) != 2 {
		t.Fatalf("the reviewer was invoked %d times, want the refused one and the one that ran with the machine's servers", len(invocations))
	}
	if invocations[1].StrictMCP {
		t.Error("the review that completed still asked the CLI to ignore the machine's tool servers")
	}
	for position, invocation := range invocations {
		if len(invocation.WritesRefused) != 1 {
			t.Fatalf("invocation %d recorded %v, and this test proves nothing unless the write was attempted and refused", position, invocation.WritesRefused)
		}
	}

	if _, err := os.Stat(run.Worktree); run.Worktree != "" && !os.IsNotExist(err) {
		t.Errorf("the worktree at %s was kept, so a tool server's write reached the change under review", run.Worktree)
	}
	if listed := session.Repository.Git("status", "--porcelain=v1", "--untracked-files=all"); listed != "?? internal/added.go" {
		t.Errorf("the person's own repository reads:\n%s", listed)
	}
}

func waitForTheSecondInvocation(t *testing.T, session *harness.Session) claudesim.Invocation {
	t.Helper()

	deadline := time.Now().Add(20 * time.Second)
	for time.Now().Before(deadline) {
		invocations := session.Simulator.Invocations()
		if len(invocations) == 2 && process.Alive(invocations[1].Process) {
			return invocations[1]
		}
		time.Sleep(20 * time.Millisecond)
	}

	t.Fatalf("the reviewer was invoked %d times, and the second one never started answering", len(session.Simulator.Invocations()))
	return claudesim.Invocation{}
}
