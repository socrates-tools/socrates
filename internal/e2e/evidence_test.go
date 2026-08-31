package e2e_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ktei/socrates/internal/findings"
	"github.com/ktei/socrates/internal/process"
	"github.com/ktei/socrates/internal/runstate"
	"github.com/ktei/socrates/internal/testing/harness"
)

func TestAMissingFindingsDocumentIsAFailure(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario("default:\n  behaviour: no-structured-output\n")
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	run := session.WaitForTerminal(name)

	if run.State != runstate.Failed {
		t.Fatalf("a review that reported nothing ended as %s: %s", run.State, reason(run))
	}
	if !strings.Contains(reason(run), "structured output") {
		t.Errorf("the failure reads %q, without naming the missing findings", reason(run))
	}
	if _, err := os.Stat(session.Round(name, 1).FindingsPath()); err == nil {
		t.Error("a findings document exists for a review that returned none")
	}
}

func TestAnUnparseableFindingsDocumentIsAFailure(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario("default:\n  behaviour: schema-violation\n")
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	run := session.WaitForTerminal(name)

	if run.State != runstate.Failed {
		t.Fatalf("a document violating the schema ended the run as %s: %s", run.State, reason(run))
	}
	if !strings.Contains(reason(run), "schema") {
		t.Errorf("the failure reads %q, without naming the violation", reason(run))
	}

	kept := session.Round(name, 1).FindingsPath()
	content := session.File(kept)
	if _, err := findings.Parse([]byte(content)); err == nil {
		t.Errorf("the kept document %s satisfies the schema after all", content)
	}
	if !strings.Contains(reason(run), kept) {
		t.Errorf("the failure reads %q, without saying where the document is kept", reason(run))
	}
}

func TestAnEmptyFindingsArrayWithEvidenceIsClean(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(`
default:
  structured: '{"findings":[],"risk_level":"low","risk_rationale":"a small and self-contained change"}'
`)
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	run := session.WaitForTerminal(name)

	if run.State != runstate.Clean {
		t.Fatalf("a completed review that found nothing ended as %s: %s", run.State, reason(run))
	}
	if run.Verdict == nil || run.Verdict.Outcome != runstate.Clean {
		t.Errorf("the verdict is %+v", run.Verdict)
	}

	document, err := findings.Parse([]byte(session.File(session.Round(name, 1).FindingsPath())))
	if err != nil {
		t.Fatalf("the findings document does not parse: %v", err)
	}
	if len(document.Findings) != 0 {
		t.Errorf("the clean run reports %d findings", len(document.Findings))
	}
}

func TestADeadReviewerIsAFailure(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario("default:\n  behaviour: silence\n  hold: 5m\n")
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	run := session.WaitFor(name, runstate.Reviewing)

	if err := process.Kill(run.Reviewer); err != nil {
		t.Fatalf("killing the reviewer: %v", err)
	}

	stopped := session.WaitForTerminal(name)
	if stopped.State != runstate.Failed {
		t.Fatalf("a run whose reviewer was killed reads as %s", stopped.State)
	}
	if !strings.Contains(reason(stopped), "went away") {
		t.Errorf("the failure reads %q, without naming that the reviewer went away", reason(stopped))
	}
	if session.Recorded(name).State != runstate.Reviewing {
		t.Error("the judgement was written to disk rather than made when the run was read")
	}
}

func TestAKilledReviewerStopsTheRun(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario("default:\n  behaviour: silence\n  hold: 5m\n")
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	run := session.WaitFor(name, runstate.Reviewing)

	if err := process.Kill(run.Reviewer); err != nil {
		t.Fatalf("killing the reviewer: %v", err)
	}

	stopped := session.WaitForTerminal(name)
	if stopped.State == runstate.Reviewing {
		t.Error("the run still reports that it is reviewing")
	}
	if stopped.State != runstate.Failed {
		t.Errorf("the run reports %s, want %s", stopped.State, runstate.Failed)
	}
}

func TestALiveReviewerIsStillWorking(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario("default:\n  events: 30\n  interval: 100ms\n")
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	run := session.WaitFor(name, runstate.Reviewing)

	for range 5 {
		if reading := session.Run(name); reading.State != runstate.Reviewing {
			t.Fatalf("a run with a live reviewer reads as %s", reading.State)
		}
		time.Sleep(30 * time.Millisecond)
	}

	if !process.Alive(run.Reviewer) {
		t.Error("the reviewer this test needs is not running")
	}
}
