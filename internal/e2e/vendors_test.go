package e2e_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/testing/harness"
	"github.com/socrates-tools/socrates/review/agent"
)

func TestAnUnsupportedCLIIsRefusedBeforeAnythingHappens(t *testing.T) {
	unsupported := agent.Unsupported()
	if len(unsupported) == 0 {
		t.Fatal("every stated reviewer CLI is supported, so nothing exercises the refusal")
	}
	turnedAway := unsupported[0]

	session := harness.Start(t)
	nameTheReviewerCLI(t, session, turnedAway.Name)
	session.Repository.Write("internal/added.go", "package internal\n")

	head := session.Repository.Head()
	status := session.Repository.Status()
	objects := session.Repository.Git("count-objects", "-v")

	out, failures, code := session.Socrates("begin")

	if code != 3 {
		t.Fatalf("socrates-review begin exited %d under an unsupported reviewer CLI, want the usage code: %s%s", code, out, failures)
	}
	if !strings.Contains(failures, turnedAway.Name) {
		t.Errorf("the refusal reads %q, without naming the CLI that was configured", failures)
	}
	if lacking := turnedAway.Lacking(); len(lacking) > 0 && !strings.Contains(failures, lacking[0]) {
		t.Errorf("the refusal reads %q, without naming the guarantee %s does not offer", failures, turnedAway.Name)
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

func TestAReviewerCLISocratesDoesNotKnowIsRefusedByName(t *testing.T) {
	session := harness.Start(t)
	nameTheReviewerCLI(t, session, "a-cli-nobody-wrote")
	session.Repository.Write("internal/added.go", "package internal\n")

	_, failures, code := session.Socrates("begin")

	if code != 3 {
		t.Fatalf("socrates-review begin exited %d under a reviewer CLI socrates does not know, want the usage code: %s", code, failures)
	}
	if !strings.Contains(failures, "a-cli-nobody-wrote") {
		t.Errorf("the refusal reads %q, without naming what was configured", failures)
	}
	for _, supported := range agent.SupportedNames() {
		if !strings.Contains(failures, supported) {
			t.Errorf("the refusal reads %q, without naming %s, which socrates does drive", failures, supported)
		}
	}
	if invocations := session.Simulator.Invocations(); len(invocations) != 0 {
		t.Errorf("the refusal spawned %d reviewers", len(invocations))
	}
}

func TestAFinishedRunStaysReadableUnderAnyConfiguration(t *testing.T) {
	unsupported := agent.Unsupported()
	if len(unsupported) == 0 {
		t.Fatal("every stated reviewer CLI is supported, so nothing exercises the refusal")
	}

	session := harness.Start(t)
	session.Simulator.Scenario(sealScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Begin()
	session.WaitForTerminal(name)

	nameTheReviewerCLI(t, session, unsupported[0].Name)

	for _, arguments := range [][]string{
		{"version"},
		{"status", name},
		{"findings", name},
		{"logs", name},
	} {
		out, failures, code := session.Socrates(arguments...)
		if code == 3 {
			t.Errorf("socrates %s was refused under an unsupported reviewer CLI, and it neither spawns nor reaps: %s", strings.Join(arguments, " "), failures)
		}
		if strings.TrimSpace(out) == "" {
			t.Errorf("socrates %s printed nothing under an unsupported reviewer CLI: %s", strings.Join(arguments, " "), failures)
		}
	}
}

func nameTheReviewerCLI(t *testing.T, session *harness.Session, named string) {
	t.Helper()

	path := filepath.Join(session.Home, "config.yaml")
	if err := os.WriteFile(path, []byte("reviewer_cli: "+named+"\n"), 0o644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
}
