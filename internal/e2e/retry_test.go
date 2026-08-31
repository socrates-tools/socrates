package e2e_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/socrates-devtool/socrates/internal/agent"
	"github.com/socrates-devtool/socrates/internal/claudesim"
	"github.com/socrates-devtool/socrates/internal/process"
	"github.com/socrates-devtool/socrates/internal/runstate"
	"github.com/socrates-devtool/socrates/internal/testing/harness"
)

const talkativeScenario = "default:\n  events: 60\n  interval: 100ms\n"

func TestCancellingStopsTheReviewerThatIsAnswering(t *testing.T) {
	// scenario: Cancelling stops everything
	session := harness.Start(t)
	session.Simulator.Scenario(talkativeScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	session.WaitFor(name, runstate.Reviewing)
	answering := waitForTheReviewerToAnswer(t, session)

	started := session.Recorded(name)
	if started.Agent != answering.Process {
		t.Errorf("the run records the process %d and the reviewer that is answering is %d", started.Agent, answering.Process)
	}

	if _, failures, code := session.Socrates("cancel", name); code != 0 {
		t.Fatalf("socrates cancel exited %d: %s", code, failures)
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

func TestAMachineThatHandsTheReviewerMoreIsRecorded(t *testing.T) {
	// scenario: A machine that hands the reviewer more than it asked for
	session := harness.Start(t)
	session.Simulator.Scenario("default:\n  extra_tools: [Bash, browser_run_code_unsafe]\n")
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	run := session.WaitForTerminal(name)
	if run.State != runstate.Clean {
		t.Fatalf("a session holding more than socrates asked for ended the run as %s: %s", run.State, reason(run))
	}

	recorded := ""
	for _, event := range session.Events(name) {
		if event.Kind == agent.EventUnexpectedTools {
			recorded = fmt.Sprint(event.Fields["tools"])
		}
	}
	if recorded == "" {
		t.Fatal("a run whose reviewer held tools socrates did not ask for records nothing about it")
	}
	for _, named := range []string{"Bash", "browser_run_code_unsafe"} {
		if !strings.Contains(recorded, named) {
			t.Errorf("the recorded event does not name %s: %q", named, recorded)
		}
	}
}

func TestAToolThatWritesCannotReachTheChangeUnderReview(t *testing.T) {
	// scenario: A tool that writes cannot reach the change under review
	session := harness.Start(t)
	session.Simulator.Scenario("default:\n  extra_tools: [Bash]\n  writes:\n    - path: .serena/project.yml\n      content: \"a tool writing its own bookkeeping as it starts\\n\"\n")
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	run := session.WaitForTerminal(name)
	if run.State != runstate.Clean {
		t.Fatalf("a reviewer whose session wrote as it started ended the run as %s: %s", run.State, reason(run))
	}
	if run.Verdict == nil || run.Verdict.Outcome != runstate.Clean {
		t.Errorf("the verdict is %+v, want the one the findings imply", run.Verdict)
	}

	invocation := session.Simulator.OnlyInvocation()
	if len(invocation.WritesRefused) != 1 {
		t.Fatalf("the invocation recorded %v, and this test proves nothing unless the write was attempted and refused", invocation.WritesRefused)
	}

	if _, err := os.Stat(run.Worktree); run.Worktree != "" && !os.IsNotExist(err) {
		t.Errorf("the worktree at %s was kept, so the write reached the change under review", run.Worktree)
	}
	if listed := session.Repository.Git("status", "--porcelain=v1", "--untracked-files=all"); listed != "?? internal/added.go" {
		t.Errorf("the person's own repository reads:\n%s", listed)
	}
}

func waitForTheReviewerToAnswer(t *testing.T, session *harness.Session) claudesim.Invocation {
	t.Helper()

	deadline := time.Now().Add(20 * time.Second)
	for time.Now().Before(deadline) {
		invocations := session.Simulator.Invocations()
		if len(invocations) == 1 && process.Alive(invocations[0].Process) {
			return invocations[0]
		}
		time.Sleep(20 * time.Millisecond)
	}

	t.Fatalf("the reviewer was invoked %d times and never started answering", len(session.Simulator.Invocations()))
	return claudesim.Invocation{}
}

func TestAReviewThatReallyFailedStillReadsAsOne(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario("default:\n  behaviour: no-result\n")
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	if run := session.WaitForTerminal(name); run.State != runstate.Failed {
		t.Fatalf("a reviewer that returned no result ended the run as %s", run.State)
	}

	printed, _, _ := session.Socrates("logs", name)
	if !strings.Contains(printed, "the reviewer failed") {
		t.Errorf("a review that failed does not read as one:\n%s", printed)
	}
}
