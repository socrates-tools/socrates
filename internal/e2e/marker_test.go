package e2e_test

import (
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/testing/harness"
	"github.com/socrates-tools/socrates/review/reviewer"
	"github.com/socrates-tools/socrates/review/runstate"
)

func TestTheReviewersEnvironmentIsMarked(t *testing.T) {
	session := harness.Start(t)
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Begin()
	session.WaitForTerminal(name)

	marker := session.Simulator.OnlyInvocation().Environment[reviewer.RunMarker]
	if marker != name {
		t.Errorf("the reviewer's environment carries %s=%q, want the run %q", reviewer.RunMarker, marker, name)
	}
}

func TestARoundInsideAReviewIsRefused(t *testing.T) {
	session := harness.Start(t)
	name := parkedRun(t, session, dialogueScenario(patchableFindings, nothingFoundAgain))

	respond(t, session, name, "--finding", "f1", "--fixed", "--note", "an empty document now answers minus one")
	respond(t, session, name, "--finding", "f2", "--rejected", "--reason", "there is no retry loop in this function at all")

	const outer = "20260819T041530Z-ee66fe34"
	out, failures, code := session.SocratesWith(
		session.Repository.Root,
		[]string{reviewer.RunMarker + "=" + outer},
		"continue", name,
	)

	if code != 3 {
		t.Fatalf("socrates review continue exited %d, want the usage code: %s%s", code, out, failures)
	}
	if !strings.Contains(failures, outer) {
		t.Errorf("the refusal reads %q, without naming the outer run", failures)
	}
	if run := session.Run(name); run.State != runstate.AwaitingResponse {
		t.Errorf("the refusal moved the run to %s", run.State)
	}
	if invocations := session.Simulator.Invocations(); len(invocations) != 1 {
		t.Errorf("the refusal started a round: %d invocations", len(invocations))
	}
}

func TestAReviewInsideAReviewIsRefused(t *testing.T) {
	session := harness.Start(t)
	session.Repository.Write("internal/added.go", "package internal\n")

	const outer = "20260819T041530Z-ee66fe34"
	out, failures, code := session.SocratesWith(
		session.Repository.Root,
		[]string{reviewer.RunMarker + "=" + outer},
		"begin",
	)

	if code != 3 {
		t.Fatalf("socrates review begin exited %d, want the usage code: %s%s", code, out, failures)
	}
	if !strings.Contains(failures, outer) {
		t.Errorf("the refusal reads %q, without naming the outer run", failures)
	}
	if names := session.Names(); len(names) != 0 {
		t.Errorf("the refusal created the run directories %v", names)
	}
}
