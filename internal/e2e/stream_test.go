package e2e_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/socrates-tools/socrates/internal/runstate"
	"github.com/socrates-tools/socrates/internal/testing/harness"
)

func TestOutputIsVisibleWhileTheReviewIsRunning(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario("default:\n  events: 40\n  interval: 60ms\n")
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	session.WaitFor(name, runstate.Reviewing)

	path := session.Round(name, 1).StreamPath()
	deadline := time.Now().Add(15 * time.Second)

	for time.Now().Before(deadline) {
		content, err := os.ReadFile(path)
		if err == nil && strings.Count(string(content), "\n") >= 2 {
			if session.Run(name).State.Terminal() {
				t.Fatal("the run finished before anything could be read from a running one")
			}
			if strings.Contains(string(content), `"type":"result"`) {
				t.Errorf("the stream already holds the result event, so it was not read part-way:\n%s", content)
			}
			return
		}
		time.Sleep(20 * time.Millisecond)
	}

	t.Errorf("nothing reached %s while the review was running", path)
}

func TestAFailedRoundKeepsWhatItGot(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario("default:\n  behaviour: truncated-stream\n  events: 3\n")
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	run := session.WaitForTerminal(name)

	if run.State != runstate.Failed {
		t.Fatalf("a truncated stream ended the run as %s: %s", run.State, reason(run))
	}
	if !strings.Contains(reason(run), "truncated") {
		t.Errorf("the failure reads %q, without naming the truncated stream", reason(run))
	}

	kept := session.File(session.Round(name, 1).StreamPath())
	if strings.Count(kept, `"type":"assistant"`) != 3 {
		t.Errorf("the round kept %d assistant events, want the 3 that arrived:\n%s",
			strings.Count(kept, `"type":"assistant"`), kept)
	}
	if strings.HasSuffix(kept, "\n") {
		t.Errorf("the kept stream ends with a newline, so the truncated line was not kept:\n%s", kept)
	}
}
