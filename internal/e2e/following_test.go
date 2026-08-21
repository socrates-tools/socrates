package e2e_test

import (
	"strings"
	"testing"
	"time"

	"github.com/rui-ktei/socrates/internal/runstate"
	"github.com/rui-ktei/socrates/internal/testing/harness"
)

func TestFollowingAReviewToACleanVerdict(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(cleanScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	printed, failures, code := session.Scrat("watch", name)
	if code != 0 {
		t.Fatalf("scrat watch of a clean review exited %d: %s\n%s", code, failures, printed)
	}
	for _, needed := range []string{"the run started", "the reviewer answered", string(runstate.Clean)} {
		if !strings.Contains(printed, needed) {
			t.Errorf("scrat watch does not report %q:\n%s", needed, printed)
		}
	}
}

func TestFollowingAReviewToAnEscalation(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(escalatingScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	printed, failures, code := session.Scrat("watch", name)
	if code != 1 {
		t.Fatalf("scrat watch of an escalating review exited %d, want 1: %s\n%s", code, failures, printed)
	}
	for _, needed := range []string{string(runstate.Escalated), "for the person"} {
		if !strings.Contains(printed, needed) {
			t.Errorf("scrat watch does not report %q:\n%s", needed, printed)
		}
	}
}

func TestEveryVerdictIsFollowedToItsOwnCode(t *testing.T) {
	for _, verdict := range []struct {
		state    runstate.State
		scenario string
	}{
		{state: runstate.Clean, scenario: cleanScenario},
		{state: runstate.Escalated, scenario: escalatingScenario},
		{state: runstate.Failed, scenario: "default:\n  behaviour: no-result\n"},
	} {
		t.Run(string(verdict.state), func(t *testing.T) {
			session := harness.Start(t)
			session.Simulator.Scenario(verdict.scenario)
			session.Repository.Write("internal/added.go", "package internal\n")

			name := session.Review()
			following := session.Starting("watch", name)

			printed, failures, followed := following.Wait()
			run := session.WaitForTerminal(name)
			if run.State != verdict.state {
				t.Fatalf("the run ended as %s, want %s: %s", run.State, verdict.state, reason(run))
			}

			_, _, reported := session.Scrat("status", name)
			if followed != reported {
				t.Errorf("the follower of a %s run exited %d and its own status exited %d: %s\n%s", run.State, followed, reported, failures, printed)
			}
			if !strings.Contains(printed, string(run.State)) {
				t.Errorf("the follower does not report the run as %s:\n%s", run.State, printed)
			}
		})
	}
}

func TestAFollowerDoesNotOutliveTheRun(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(slowScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	following := session.Starting("watch", name)

	run := session.WaitForTerminal(name)
	ended := time.Now()

	printed, failures, code := following.Wait()
	if run.State != runstate.Clean || code != 0 {
		t.Fatalf("a run that ended %s was followed to exit %d: %s\n%s", run.State, code, failures, printed)
	}
	if outlived := time.Since(ended); outlived > 10*time.Second {
		t.Errorf("the follower kept waiting for %s after the run reached %s", outlived, run.State)
	}
}

func TestWatchingImmediatelyAfterStartingARun(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(slowScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	following := session.Starting("watch", name)

	printed, failures, code := following.Wait()
	if run := session.Run(name); run.State != runstate.Clean || code != 0 {
		t.Fatalf("a run that ended %s was followed from its start to exit %d: %s\n%s", run.State, code, failures, printed)
	}
	if !strings.Contains(printed, "the run started") {
		t.Errorf("a follower attached at the start missed the events that followed:\n%s", printed)
	}
}

func TestCancellingWhileFollowing(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(slowScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	session.WaitFor(name, runstate.Reviewing)
	following := session.Starting("watch", name)

	if _, failures, code := session.Scrat("cancel", name); code != 0 {
		t.Fatalf("scrat cancel exited %d: %s", code, failures)
	}

	printed, failures, code := following.Wait()
	if code != 2 {
		t.Fatalf("scrat watch of a cancelled run exited %d, want 2: %s\n%s", code, failures, printed)
	}
	if !strings.Contains(printed, string(runstate.Cancelled)) {
		t.Errorf("scrat watch does not say the run was cancelled:\n%s", printed)
	}
}

func TestAFollowerStopsWhereTheRoundStops(t *testing.T) {
	session := harness.Start(t)
	name := parkedRun(t, session, dialogueScenario(patchableFindings, nothingFoundAgain))

	printed, failures, code := session.Scrat("watch", name)
	if code != 4 {
		t.Fatalf("scrat watch of a parked round exited %d, want 4: %s\n%s", code, failures, printed)
	}
	for _, needed := range []string{string(runstate.AwaitingResponse), "is waiting on you", "unanswered"} {
		if !strings.Contains(printed, needed) {
			t.Errorf("scrat watch does not report %q for a parked round:\n%s", needed, printed)
		}
	}
	if parked := lineSaying(printed, "is waiting on you"); strings.Contains(parked, "person") {
		t.Errorf("a round parked on patches alone reports a person to answer it: %q", parked)
	}
}

func lineSaying(printed, said string) string {
	for line := range strings.SplitSeq(printed, "\n") {
		if strings.Contains(line, said) {
			return line
		}
	}
	return ""
}

func TestAFollowerSaysAQuestionIsWhatTheRoundWaitsOn(t *testing.T) {
	session := harness.Start(t)
	name := parkedRun(t, session, dialogueScenario(questionOnly, nothingFoundAgain))

	printed, failures, code := session.Scrat("watch", name)
	if code != 4 {
		t.Fatalf("scrat watch of a round parked on a question exited %d, want 4: %s\n%s", code, failures, printed)
	}
	parked := lineSaying(printed, "is waiting on you")
	if parked == "" {
		t.Fatalf("scrat watch reports no parked round:\n%s", printed)
	}
	if !strings.Contains(parked, "1 need a person") {
		t.Errorf("the one line a follower prints for a round parked on a question alone does not say what it waits on: %q", parked)
	}
}

func TestReadingAFinishedRun(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(escalatingScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	session.WaitForTerminal(name)

	printed, failures, code := session.Scrat("logs", name)
	if code != 1 {
		t.Fatalf("scrat logs of an escalated run exited %d, want 1: %s", code, failures)
	}

	events := session.Events(name)
	if len(events) < 4 {
		t.Fatalf("the log holds %d events, and this test needs a run that recorded several", len(events))
	}
	if lines := strings.Count(strings.TrimSpace(printed), "\n") + 1; lines < len(events) {
		t.Errorf("scrat logs printed %d lines for %d events:\n%s", lines, len(events), printed)
	}
	if !strings.Contains(printed, string(runstate.Escalated)) {
		t.Errorf("scrat logs does not report the verdict:\n%s", printed)
	}
}

func TestReadingARunStillInProgress(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(slowScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	session.WaitFor(name, runstate.Reviewing)

	started := time.Now()
	printed, failures, code := session.Scrat("logs", name)
	if code != 4 {
		t.Fatalf("scrat logs of a run still reviewing exited %d, want 4: %s\n%s", code, failures, printed)
	}
	if waited := time.Since(started); waited > 5*time.Second {
		t.Errorf("scrat logs waited %s for a run that had not finished", waited)
	}
	if !strings.Contains(printed, string(runstate.Reviewing)) {
		t.Errorf("scrat logs does not say the run is still working:\n%s", printed)
	}
}
