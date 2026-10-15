package e2e_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	skillkit "github.com/socrates-tools/socrates/internal/skill"
	"github.com/socrates-tools/socrates/internal/testing/harness"
	"github.com/socrates-tools/socrates/review/agent"
	"github.com/socrates-tools/socrates/review/runstate"
)

func startWithNoReviewerCLISet(t *testing.T, vendor agent.Vendor) *harness.Session {
	t.Helper()

	session := harness.StartFor(t, vendor)
	if err := os.Remove(filepath.Join(session.Home, "config.yaml")); err != nil && !os.IsNotExist(err) {
		t.Fatalf("leaving reviewer_cli unset: %v", err)
	}
	session.Repository.Write("internal/added.go", "package internal\n")
	return session
}

func TestAHostWithNoReviewerCLISetDecidesWhichReviewerRuns(t *testing.T) {
	t.Parallel()

	codex, _ := agent.Named("codex")
	for _, begun := range []struct {
		named     string
		simulates agent.Vendor
		arguments []string
		reviewer  string
	}{
		{"begin --host codex spawns codex", codex, []string{"--host", skillkit.Codex.Flag}, "codex"},
		{"begin --host claude spawns claude", agent.TheDefault(), []string{"--host", skillkit.ClaudeCode.Flag}, "claude"},
		{"begin with no host spawns claude", agent.TheDefault(), nil, "claude"},
	} {
		t.Run(begun.named, func(t *testing.T) {
			t.Parallel()
			session := startWithNoReviewerCLISet(t, begun.simulates)

			name := session.Begin(begun.arguments...)
			run := session.WaitForTerminal(name)

			if run.ReviewerCLI != begun.reviewer {
				t.Errorf("the run records the reviewer CLI %q, want %q (%s)", run.ReviewerCLI, begun.reviewer, run.State)
			}
			if invocations := session.Simulator.Invocations(); len(invocations) != 1 {
				t.Errorf("the %s simulator was invoked %d times, want once", begun.reviewer, len(invocations))
			}
		})
	}
}

func TestAReviewerCLIAPersonSetWinsWhateverTheHost(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	nameTheReviewerCLI(t, session, "claude")
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Begin("--host", skillkit.Codex.Flag)
	run := session.WaitForTerminal(name)

	if run.ReviewerCLI != "claude" {
		t.Errorf("the run records the reviewer CLI %q under reviewer_cli: claude, want claude (%s)", run.ReviewerCLI, run.State)
	}
	if invocations := session.Simulator.Invocations(); len(invocations) != 1 {
		t.Errorf("the claude simulator was invoked %d times, want once", len(invocations))
	}
}

func TestEveryRoundOfARunIsReviewedByTheCLIItsFirstRoundWas(t *testing.T) {
	t.Parallel()

	codex, _ := agent.Named("codex")
	for _, given := range [][]string{
		{"--host", skillkit.Codex.Flag},
		{"--host", skillkit.ClaudeCode.Flag},
		nil,
	} {
		t.Run("continue "+strings.Join(given, " "), func(t *testing.T) {
			t.Parallel()
			session := startWithNoReviewerCLISet(t, codex)
			session.Simulator.Scenario(dialogueScenario(patchableFindings, nothingFoundAgain))
			session.Repository.Write("internal/added.go", answeredWrongly)

			name := session.Begin("--host", skillkit.Codex.Flag, "--intent", "count the characters of a document")
			session.WaitFor(name, runstate.AwaitingResponse)
			respond(t, session, name, "--finding", "f1", "--fixed")
			respond(t, session, name, "--finding", "f2", "--rejected", "--reason", "the retry loop is bounded by the caller's deadline")

			arguments := append([]string{"continue", name, "--wait"}, given...)
			if out, failures, code := session.Socrates(arguments...); code != 0 {
				t.Fatalf("socrates review %s exited %d: %s\n%s", strings.Join(arguments, " "), code, failures, out)
			}

			run := session.Run(name)
			if run.ReviewerCLI != "codex" {
				t.Errorf("the run records the reviewer CLI %q after a continue, want codex", run.ReviewerCLI)
			}
			if invocations := session.Simulator.Invocations(); len(invocations) != 2 {
				t.Errorf("the codex simulator was invoked %d times across two rounds, want twice", len(invocations))
			}
		})
	}
}

func TestAHostNamingNeitherHostIsRefusedBeforeASnapshotIsTaken(t *testing.T) {
	t.Parallel()
	session := startWithNoReviewerCLISet(t, agent.TheDefault())

	head := session.Repository.Head()
	status := session.Repository.Status()
	objects := session.Repository.Git("count-objects", "-v")

	out, failures, code := session.Socrates("begin", "--host", "a-host-nobody-has")

	if code != 3 {
		t.Fatalf("socrates review begin --host naming neither host exited %d, want the usage code: %s%s", code, out, failures)
	}
	for _, named := range append([]string{"a-host-nobody-has"}, skillkit.ClaudeCode.Flag, skillkit.Codex.Flag) {
		if !strings.Contains(failures, named) {
			t.Errorf("the refusal reads %q, without naming %s", failures, named)
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
}
