package e2e_test

import (
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/convention"
	"github.com/socrates-tools/socrates/internal/testing/harness"
)

const theGuaranteeThatIsMissing = "process group"

func TestAnUnsupportedPlatformRefusesTheReviewBeforeAnythingStarts(t *testing.T) {
	session := harness.Start(t).WhereNothingReaps()
	session.Repository.Write("internal/added.go", "package internal\n")

	head := session.Repository.Head()
	status := session.Repository.Status()
	objects := session.Repository.Git("count-objects", "-v")

	out, failures, code := session.Socrates("begin")

	if code != 3 {
		t.Fatalf("socrates-review begin exited %d where nothing reaps, want the usage code: %s%s", code, out, failures)
	}
	for _, said := range []string{theGuaranteeThatIsMissing, convention.TheWayForward} {
		if !strings.Contains(failures, said) {
			t.Errorf("the refusal reads %q, without saying %q", failures, said)
		}
	}

	if names := session.Names(); len(names) != 0 {
		t.Errorf("the refusal created the run directories %v", names)
	}
	if invocations := session.Simulator.Invocations(); len(invocations) != 0 {
		t.Errorf("the refusal spawned %d reviewers", len(invocations))
	}
	if now := session.Repository.Git("count-objects", "-v"); now != objects {
		t.Errorf("the refusal wrote objects into the repository, so a snapshot was taken:\n%s\nwant:\n%s", now, objects)
	}
	if now := session.Repository.Head(); now != head {
		t.Errorf("the refusal moved HEAD from %s to %s", head, now)
	}
	if now := session.Repository.Status(); now != status {
		t.Errorf("git status reads %q after the refusal, want %q", now, status)
	}
	if listed := session.Repository.Worktrees(); strings.Contains(listed, session.Home) {
		t.Errorf("the refusal left a worktree under %s:\n%s", session.Home, listed)
	}
}

func TestReadingARunIsNotRefusedOnAnUnsupportedPlatform(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(sealScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Begin()
	session.WaitForTerminal(name)

	elsewhere := session.WhereNothingReaps()

	for _, arguments := range [][]string{
		{"version"},
		{"status", name},
		{"findings", name},
		{"logs", name},
	} {
		out, failures, code := elsewhere.Socrates(arguments...)
		if code == 3 {
			t.Errorf("socrates %s was refused where nothing reaps, and it neither spawns nor reaps: %s", strings.Join(arguments, " "), failures)
		}
		if strings.TrimSpace(out) == "" {
			t.Errorf("socrates %s printed nothing where nothing reaps: %s", strings.Join(arguments, " "), failures)
		}
	}

	printed, failures, _ := elsewhere.Socrates("findings", name)
	if !strings.Contains(printed, "f1") {
		t.Errorf("what the run recorded reads %q where nothing reaps, without naming the finding: %s", printed, failures)
	}
}
