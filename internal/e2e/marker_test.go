package e2e_test

import (
	"strings"
	"testing"

	"github.com/rui-ktei/socrates/internal/reviewer"
	"github.com/rui-ktei/socrates/internal/testing/harness"
)

func TestTheReviewersEnvironmentIsMarked(t *testing.T) {
	session := harness.Start(t)
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	session.WaitForTerminal(name)

	marker := session.Simulator.OnlyInvocation().Environment[reviewer.RunMarker]
	if marker != name {
		t.Errorf("the reviewer's environment carries %s=%q, want the run %q", reviewer.RunMarker, marker, name)
	}
}

func TestAReviewInsideAReviewIsRefused(t *testing.T) {
	session := harness.Start(t)
	session.Repository.Write("internal/added.go", "package internal\n")

	const outer = "20260819T041530Z-ee66fe34"
	out, failures, code := session.ScratWith(
		session.Repository.Root,
		[]string{reviewer.RunMarker + "=" + outer},
		"review",
	)

	if code != 3 {
		t.Fatalf("scrat review exited %d, want the usage code: %s%s", code, out, failures)
	}
	if !strings.Contains(failures, outer) {
		t.Errorf("the refusal reads %q, without naming the outer run", failures)
	}
	if names := session.Names(); len(names) != 0 {
		t.Errorf("the refusal created the run directories %v", names)
	}
}
