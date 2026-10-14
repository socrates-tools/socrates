package e2e_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/testing/harness"
	"github.com/socrates-tools/socrates/review/agent"
	"github.com/socrates-tools/socrates/review/reviewer"
	"github.com/socrates-tools/socrates/review/runlog"
	"github.com/socrates-tools/socrates/review/runstate"
)

func TestAStartedRunIsOnDiskImmediately(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	session.Simulator.Scenario(slowScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Begin()
	directory := session.Store.Path(name)

	run := session.Recorded(name)
	if run.Repository != session.Repository.Root {
		t.Errorf("run.json names the repository %s, want %s", run.Repository, session.Repository.Root)
	}
	if run.Snapshot == "" {
		t.Error("run.json names no snapshot commit")
	}
	if run.Base == "" {
		t.Error("run.json names no base commit")
	}

	for _, expected := range []string{"run.json", runlog.FileName, "r1"} {
		if _, err := os.Stat(filepath.Join(directory, expected)); err != nil {
			t.Errorf("%s is not in the run directory when socrates review begin returned: %v", expected, err)
		}
	}
}

func TestRoundOutputIsUnderTheRound(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	session.Simulator.Scenario(escalatingScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Begin("--intent", "add a thing")
	session.WaitForTerminal(name)

	round := session.Round(name, 1)
	for _, path := range []string{round.PromptPath(), round.StreamPath(), round.FindingsPath()} {
		written, err := os.Stat(path)
		if err != nil {
			t.Errorf("%s is not under the round: %v", path, err)
			continue
		}
		if written.Size() == 0 {
			t.Errorf("%s is empty", path)
		}
		if filepath.Dir(path) != round.Directory {
			t.Errorf("%s is not under %s", path, round.Directory)
		}
	}
}

func TestEventsAreOneJSONObjectPerLineThroughARun(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	session.Simulator.Scenario(escalatingScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Begin()
	run := session.WaitForTerminal(name)
	if run.State != runstate.Escalated {
		t.Fatalf("the run ended as %s: %s", run.State, reason(run))
	}
	session.WaitForTheLogToRecordTheEnd(name)

	content := session.File(filepath.Join(session.Store.Path(name), runlog.FileName))
	for line := range strings.SplitSeq(strings.TrimSpace(content), "\n") {
		var read map[string]any
		if err := json.Unmarshal([]byte(line), &read); err != nil {
			t.Fatalf("a line of the event log is not one JSON object: %v:\n%s", err, line)
		}
		for _, field := range []string{"at", "kind", "subject"} {
			if read[field] == nil || read[field] == "" {
				t.Errorf("the event %s carries no %s", line, field)
			}
		}
	}

	var kinds []string
	for _, event := range session.Events(name) {
		kinds = append(kinds, event.Kind)
	}
	for _, needed := range []string{reviewer.EventRunStarted, agent.EventStarted, agent.EventCompleted, reviewer.EventRunEnded} {
		if !slices.Contains(kinds, needed) {
			t.Errorf("the event log holds no %s event: %v", needed, kinds)
		}
	}
}

func TestWhatTheAgentCostIsRecordedInTheLog(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	session.Simulator.Scenario("default:\n  input_tokens: 1234\n  output_tokens: 567\n" + strings.TrimPrefix(cleanScenario, "\ndefault:\n"))
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Begin()
	session.WaitForTerminal(name)

	for _, event := range session.Events(name) {
		if event.Kind != agent.EventCompleted {
			continue
		}
		if event.Fields["model"] == "" || event.Fields["model"] == nil {
			t.Errorf("the completion event names no model: %+v", event)
		}
		if event.Fields["input_tokens"] != float64(1234) || event.Fields["output_tokens"] != float64(567) {
			t.Errorf("the completion event reports %+v, want what the agent said it cost", event.Fields)
		}
		return
	}
	t.Error("the event log holds no completion event")
}
