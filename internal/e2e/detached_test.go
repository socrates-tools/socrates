package e2e_test

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/socrates-tools/socrates/internal/process"
	"github.com/socrates-tools/socrates/internal/testing/harness"
	"github.com/socrates-tools/socrates/review/runstate"
)

func TestTheCommandReturnsWhileTheReviewRuns(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario("default:\n  events: 40\n  interval: 60ms\n")
	session.Repository.Write("internal/added.go", "package internal\n")

	started := time.Now()
	out, failures, code := session.Socrates("begin", "--intent", "add a thing")
	if code != 0 {
		t.Fatalf("socrates review begin exited %d: %s", code, failures)
	}
	returned := time.Since(started)

	name, _, _ := strings.Cut(strings.TrimSpace(out), "\n")
	if name == "" {
		t.Fatal("socrates review begin printed no run identifier")
	}

	run := session.WaitFor(name, runstate.Reviewing)
	if run.Reviewer == 0 {
		t.Fatal("the run records no reviewer process")
	}
	if !process.Alive(run.Reviewer) {
		t.Errorf("the reviewer %d is not running, so the command did not return early", run.Reviewer)
	}
	if returned > 20*time.Second {
		t.Errorf("socrates review begin took %s to return", returned)
	}

	t.Run("the reviewer holds none of the invoking terminal's streams", func(t *testing.T) {
		for _, descriptor := range []string{"0", "1", "2"} {
			held, err := os.Readlink(filepath.Join("/proc", strconv.Itoa(run.Reviewer), "fd", descriptor))
			if err != nil {
				t.Skipf("this machine does not report what the reviewer holds open: %v", err)
			}
			if strings.Contains(held, "/dev/pts/") || strings.Contains(held, "/dev/tty") {
				t.Errorf("the reviewer holds the terminal on descriptor %s: %s", descriptor, held)
			}
		}
	})

	session.WaitForTerminal(name)
}

func TestTheReviewerSurvivesItsCaller(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario("default:\n  events: 6\n  interval: 80ms\n")
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Begin()

	if process.Alive(0) {
		t.Fatal("liveness answers yes for no process at all")
	}

	run := session.WaitForTerminal(name)
	if run.State != runstate.Clean {
		t.Errorf("the run ended as %s: %s", run.State, reason(run))
	}
	if session.Recorded(name).Reviewer == 0 {
		t.Error("the run does not record which process reviewed it")
	}
}

func TestTheRunDirectoryHoldsTheRound(t *testing.T) {
	session := harness.Start(t)
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Begin("--intent", "add a thing")
	session.WaitForTerminal(name)

	round := session.Round(name, 1)
	if asked := session.File(round.PromptPath()); !strings.Contains(asked, "add a thing") {
		t.Errorf("the prompt on disk does not carry the intent:\n%s", asked)
	}
	if stream := session.File(round.StreamPath()); !strings.Contains(stream, `"type":"result"`) {
		t.Errorf("the stream on disk holds no result event:\n%s", stream)
	}
	if reported := session.File(round.FindingsPath()); !strings.Contains(reported, "risk_level") {
		t.Errorf("the findings on disk are not a findings document:\n%s", reported)
	}

	events := session.Events(name)
	if len(events) == 0 {
		t.Error("the run recorded no events")
	}
}
