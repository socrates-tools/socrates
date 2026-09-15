package e2e_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/runstate"
	"github.com/socrates-tools/socrates/internal/testing/harness"
)

const loudScenario = "default:\n  events: 60\n  interval: 1ms\n"

func TestStatusStaysSmallWhenTheReviewerIsLoud(t *testing.T) {
	quiet := lineCountOfStatus(t, cleanScenario)
	loud := lineCountOfStatus(t, loudScenario)

	if loud != quiet {
		t.Errorf("socrates-review status printed %d lines for a loud reviewer and %d for a quiet one", loud, quiet)
	}
}

func lineCountOfStatus(t *testing.T, scenario string) int {
	t.Helper()

	session := harness.Start(t)
	session.Simulator.Scenario(scenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	run := session.WaitForTerminal(name)
	if run.State != runstate.Clean {
		t.Fatalf("the run ended as %s, and this comparison needs two runs that ended alike: %s", run.State, reason(run))
	}

	printed, failures, code := session.Socrates("status", name)
	if code != 0 {
		t.Fatalf("socrates-review status exited %d: %s", code, failures)
	}
	return strings.Count(strings.TrimSpace(printed), "\n") + 1
}

func TestNeitherCommandGrowsTheOthersFlag(t *testing.T) {
	session := harness.Start(t)

	for command, want := range map[string][]string{
		"status": {"--help", "--json"},
		"watch":  {"--help"},
		"logs":   {"--help"},
	} {
		printed, failures, code := session.Socrates(command, "--help")
		if code != 0 {
			t.Fatalf("socrates %s --help exited %d: %s", command, code, failures)
		}
		if offered := flagsIn(printed); !slices.Equal(offered, want) {
			t.Errorf("socrates %s offers %v, want %v", command, offered, want)
		}
	}
}

func flagsIn(help string) []string {
	var offered []string
	inside := false

	for line := range strings.SplitSeq(help, "\n") {
		if strings.HasPrefix(line, "Flags:") {
			inside = true
			continue
		}
		if !inside {
			continue
		}
		if strings.TrimSpace(line) == "" {
			break
		}
		for _, word := range strings.Fields(line) {
			if strings.HasPrefix(word, "--") {
				offered = append(offered, strings.TrimSuffix(word, ","))
				break
			}
		}
	}

	slices.Sort(offered)
	return offered
}

func TestAllThreeAgreeAboutTheRun(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(escalatingScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	run := session.WaitForTerminal(name)
	if run.State != runstate.Escalated {
		t.Fatalf("the run ended as %s, want escalated: %s", run.State, reason(run))
	}

	codes := map[string]int{}
	for _, command := range []string{"status", "watch", "logs"} {
		printed, failures, code := session.Socrates(command, name)
		codes[command] = code
		if !strings.Contains(printed, string(run.State)) {
			t.Errorf("socrates %s does not report the run as %s: %s\n%s", command, run.State, failures, printed)
		}
	}

	if codes["status"] != codes["watch"] || codes["status"] != codes["logs"] {
		t.Errorf("the three commands exited %v for one run", codes)
	}
}
