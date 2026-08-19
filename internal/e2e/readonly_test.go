package e2e_test

import (
	"os"
	"strings"
	"testing"

	"github.com/rui-ktei/socrates/internal/runstate"
	"github.com/rui-ktei/socrates/internal/testing/harness"
)

func TestAReviewerThatWroteFailsTheRound(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario("default:\n  writes:\n    - path: README.md\n      content: \"the reviewer helped itself\\n\"\n")
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	run := session.WaitForTerminal(name)

	if run.State != runstate.Failed {
		t.Fatalf("a reviewer that wrote ended the run as %s: %s", run.State, reason(run))
	}
	if !strings.Contains(reason(run), "README.md") {
		t.Errorf("the failure reads %q, without naming the path that differs", reason(run))
	}
	if run.Verdict == nil || run.Verdict.Outcome != runstate.Failed {
		t.Errorf("the verdict is %+v, and an absence of findings is not the result here", run.Verdict)
	}

	t.Run("A dirtied worktree is kept", func(t *testing.T) {
		worktree := session.Recorded(name).Worktree
		if _, err := os.Stat(worktree); err != nil {
			t.Errorf("the worktree at %s was removed, so the difference cannot be inspected: %v", worktree, err)
		}
		if !strings.Contains(reason(run), worktree) {
			t.Errorf("the failure reads %q, without recording where the worktree is", reason(run))
		}
		if listed := session.Repository.Worktrees(); !strings.Contains(listed, worktree) {
			t.Errorf("git worktree list no longer names %s:\n%s", worktree, listed)
		}
	})
}

func TestAReviewerThatCreatedAFileFailsTheRound(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario("default:\n  writes:\n    - path: internal/reviewer-notes.md\n      content: \"some thoughts\\n\"\n")
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	run := session.WaitForTerminal(name)

	if run.State != runstate.Failed {
		t.Fatalf("a reviewer that created a file ended the run as %s: %s", run.State, reason(run))
	}
	if !strings.Contains(reason(run), "internal/reviewer-notes.md") {
		t.Errorf("the failure reads %q, without naming the created path", reason(run))
	}
}

func TestACleanReviewerPassesTheCheck(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(`
default:
  structured: '{"findings":[{"id":"f1","severity":"warning","file":"internal/added.go","line":1,"description":"this package has no name anybody would guess","action":"ask-user"}],"risk_level":"medium","risk_rationale":"one small file"}'
`)
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	run := session.WaitForTerminal(name)

	if run.State != runstate.Escalated {
		t.Fatalf("the run ended as %s, want the verdict the findings imply: %s", run.State, reason(run))
	}
	if strings.Contains(reason(run), "wrote to the snapshot worktree") {
		t.Errorf("a reviewer that wrote nothing was reported as having written: %q", reason(run))
	}
}
