package e2e_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates-review/internal/runstate"
	"github.com/socrates-tools/socrates-review/internal/testing/harness"
)

func TestAReviewerThatWroteFailsTheRound(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario("default:\n  writes:\n    - path: README.md\n      content: \"the reviewer helped itself\\n\"\n      defeats_the_seal: true\n")
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

	t.Run("The seal is lifted when the round ends", func(t *testing.T) {
		worktree := session.Recorded(name).Worktree

		inspected := filepath.Join(worktree, "internal", "added.go")
		if err := os.WriteFile(inspected, []byte("package internal\n\nfunc Answer() int { return 1 }\n"), 0o644); err != nil {
			t.Errorf("a person cannot edit %s in the worktree the run kept: %v", inspected, err)
		}
		if err := os.RemoveAll(worktree); err != nil {
			t.Errorf("a person cannot delete the worktree the run kept at %s: %v", worktree, err)
		}
	})
}

func TestAReviewerThatCreatedAFileFailsTheRound(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario("default:\n  writes:\n    - path: internal/reviewer-notes.md\n      content: \"some thoughts\\n\"\n      defeats_the_seal: true\n")
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

const sealScenario = `
default:
  structured: '{"findings":[{"id":"f1","severity":"warning","file":"internal/added.go","line":1,"description":"this file has no caller, and whether it should exist at all is not something a patch settles","action":"no-op"}],"risk_level":"medium","risk_rationale":"one small file"}'
`

func eventNaming(t *testing.T, session *harness.Session, name, said string) string {
	t.Helper()

	printed, failures, _ := session.Socrates("logs", name)
	if printed == "" {
		t.Fatalf("socrates-review logs printed nothing for %s: %s", name, failures)
	}
	return lineSaying(printed, said)
}

func TestACleanReviewerPassesTheCheck(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(sealScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	run := session.WaitForTerminal(name)

	if named := eventNaming(t, session, name, "wrote to the snapshot worktree"); named != "" {
		t.Errorf("a reviewer that wrote nothing was recorded as having written: %s", named)
	}
	if strings.Contains(reason(run), "wrote to the snapshot worktree") {
		t.Errorf("a reviewer that wrote nothing was reported as having written: %q", reason(run))
	}
	if run.State != runstate.Escalated {
		t.Fatalf("the run ended as %s, want the verdict the findings imply: %s", run.State, reason(run))
	}
}

func TestAWriteTheSealRefusedIsNotAFindingAgainstTheReviewer(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(sealScenario + `  writes:
    - path: .serena/project.yml
      content: "a tool writing its own bookkeeping\n"
    - path: internal/added.go
      content: "package internal // and a fix nobody asked for\n"
`)
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	run := session.WaitForTerminal(name)

	if named := eventNaming(t, session, name, ".serena"); named != "" {
		t.Errorf("a write the seal refused was recorded against the reviewer: %s", named)
	}
	if strings.Contains(reason(run), ".serena") || strings.Contains(reason(run), "wrote to the snapshot worktree") {
		t.Errorf("a write the seal refused was reported against the reviewer: %q", reason(run))
	}
	if run.State != runstate.Escalated {
		t.Fatalf("the run ended as %s, want the verdict the findings imply: %s", run.State, reason(run))
	}

	refused := session.Simulator.OnlyInvocation().WritesRefused
	if len(refused) != 2 {
		t.Fatalf("the reviewer's session recorded %v, and this test proves nothing unless both writes were attempted and refused", refused)
	}
	for _, attempt := range refused {
		if !strings.Contains(attempt, "permission denied") {
			t.Errorf("a write into the sealed worktree failed for another reason: %s", attempt)
		}
	}

	if _, err := os.Stat(run.Worktree); run.Worktree != "" && !os.IsNotExist(err) {
		t.Errorf("the worktree at %s was kept, so something the seal refused was still counted as a difference", run.Worktree)
	}
	if content := session.Repository.Read("internal/added.go"); content != "package internal\n" {
		t.Errorf("the person's own file reads %q after a reviewer tried to fix it", content)
	}
}
