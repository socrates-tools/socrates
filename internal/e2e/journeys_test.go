package e2e_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rui-ktei/socrates/internal/corpus"
	"github.com/rui-ktei/socrates/internal/findings"
	"github.com/rui-ktei/socrates/internal/runstate"
	"github.com/rui-ktei/socrates/internal/testing/harness"
)

func TestTheCleanJourney(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(cleanScenario)

	session.Repository.Write("internal/parser/parse.go", "package parser\n\nfunc Parse(document string) error {\n\treturn nil\n}\n")
	before := session.Repository.Status()

	printed, failures, code := session.Scrat("review", "--wait", "--intent", "add a parser that accepts anything")
	if code != 0 {
		t.Fatalf("scrat review --wait exited %d for a clean review: %s\n%s", code, failures, printed)
	}

	name := firstLine(printed)
	run := session.Run(name)
	if run.State != runstate.Clean {
		t.Fatalf("the run ended as %s: %s", run.State, reason(run))
	}

	if _, _, statusCode := session.Scrat("status", name); statusCode != 0 {
		t.Errorf("scrat status of a clean run exited %d, want 0", statusCode)
	}
	if _, _, findingsCode := session.Scrat("findings", name); findingsCode != 0 {
		t.Errorf("scrat findings of a clean run exited %d, want 0", findingsCode)
	}

	if after := session.Repository.Status(); after != before {
		t.Errorf("the repository changed:\nbefore:\n%s\nafter:\n%s", before, after)
	}
	if _, err := os.Stat(run.Worktree); run.Worktree != "" && !os.IsNotExist(err) {
		t.Errorf("the worktree at %s was not cleaned up", run.Worktree)
	}

	invocation := session.Simulator.OnlyInvocation()
	if !strings.Contains(invocation.Prompt, "add a parser that accepts anything") {
		t.Error("the reviewer was not told what the person asked for")
	}
	if invocation.Directory != run.Worktree {
		t.Errorf("the reviewer ran in %s, want the snapshot worktree %s", invocation.Directory, run.Worktree)
	}
}

func TestTheEscalatedJourney(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(escalatingScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	printed, failures, code := session.Scrat("review", "--wait")
	if code != 1 {
		t.Fatalf("scrat review --wait exited %d for a review that escalates: %s\n%s", code, failures, printed)
	}

	name := firstLine(printed)
	run := session.Run(name)
	if run.State != runstate.Escalated {
		t.Fatalf("the run ended as %s: %s", run.State, reason(run))
	}

	document, err := findings.Parse([]byte(session.File(session.Round(name, 1).FindingsPath())))
	if err != nil {
		t.Fatalf("the findings on disk do not parse: %v", err)
	}
	classified := findings.Classify(document)
	if len(classified.NeedingAPerson) != 1 || len(classified.Patchable) != 1 || len(classified.NoOps) != 1 {
		t.Errorf("the findings classified as %d for the person, %d patchable, %d no-ops",
			len(classified.NeedingAPerson), len(classified.Patchable), len(classified.NoOps))
	}

	readable, _, findingsCode := session.Scrat("findings", name)
	if findingsCode != 1 {
		t.Errorf("scrat findings of an escalated run exited %d, want 1", findingsCode)
	}
	for _, needed := range []string{"needs a person", "can be patched", "informational"} {
		if !strings.Contains(readable, needed) {
			t.Errorf("scrat findings does not say who %q is for:\n%s", needed, readable)
		}
	}
}

func TestTheEscalatedJourneyOnARecordedAnswer(t *testing.T) {
	entry, err := corpus.Named("dialogue-layer-first-look")
	if err != nil {
		t.Fatalf("reading the recorded answer: %v", err)
	}

	session := harness.Start(t)
	session.Simulator.Scenario("default:\n  recorded: " + entry.Name + "\n")
	session.Repository.Write("internal/added.go", "package internal\n")

	printed, failures, code := session.Scrat("review", "--wait")
	if code != 1 {
		t.Fatalf("scrat review --wait exited %d for a recorded answer that escalates: %s\n%s", code, failures, printed)
	}

	name := firstLine(printed)
	run := session.Run(name)
	if run.State != runstate.Escalated {
		t.Fatalf("the run ended as %s: %s", run.State, reason(run))
	}

	stored := strings.TrimSuffix(session.File(session.Round(name, 1).FindingsPath()), "\n")
	if stored != string(entry.Document) {
		t.Errorf("what a real reviewer answered did not reach disk intact: %d bytes recorded, %d stored", len(entry.Document), len(stored))
	}

	readable, _, findingsCode := session.Scrat("findings", name)
	if findingsCode != 1 {
		t.Errorf("scrat findings of an escalated run exited %d, want 1", findingsCode)
	}
	if !strings.Contains(readable, "needs a person") {
		t.Errorf("scrat findings does not say who the recorded findings are for:\n%s", readable)
	}
}

func TestNoFailureJourneyReportsClean(t *testing.T) {
	for _, failure := range []struct {
		behaviour string
		names     string
	}{
		{behaviour: "truncated-stream", names: "truncated"},
		{behaviour: "no-result", names: "result event"},
		{behaviour: "no-structured-output", names: "structured output"},
		{behaviour: "schema-violation", names: "schema"},
		{behaviour: "non-zero-exit", names: "exited"},
	} {
		t.Run(failure.behaviour, func(t *testing.T) {
			session := harness.Start(t)
			session.Simulator.Scenario("default:\n  behaviour: " + failure.behaviour + "\n")
			session.Repository.Write("internal/added.go", "package internal\n")

			name := session.Review()
			run := session.WaitForTerminal(name)

			if run.State != runstate.Failed {
				t.Fatalf("%s ended the run as %s: %s", failure.behaviour, run.State, reason(run))
			}
			if !strings.Contains(reason(run), failure.names) {
				t.Errorf("the failure reads %q, without naming %q", reason(run), failure.names)
			}

			printed, _, code := session.Scrat("findings", name)
			if code != 2 {
				t.Errorf("scrat findings of a failed run exited %d, want 2", code)
			}
			if !strings.Contains(printed, "no review completed") {
				t.Errorf("scrat findings does not say that no review completed:\n%s", printed)
			}
		})
	}
}

func TestASilentReviewerIsAFailureJourney(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario("default:\n  behaviour: silence\n  hold: 30s\n")
	session.Repository.Write("internal/added.go", "package internal\n")

	path := filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "socrates", "config.yaml")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("creating the configuration directory: %v", err)
	}
	if err := os.WriteFile(path, []byte("agent_deadline: 2s\n"), 0o644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}

	name := session.Review()
	run := session.WaitForTerminal(name)

	if run.State != runstate.Failed {
		t.Fatalf("a reviewer that answered nothing ended the run as %s: %s", run.State, reason(run))
	}
	if !strings.Contains(reason(run), "2s") {
		t.Errorf("the failure reads %q, without naming the deadline", reason(run))
	}
	waitForDeath(t, name, 0, session.Recorded(name).Agent)
}

func TestTheConcurrentEditingJourney(t *testing.T) {
	const reviewed = "package internal\n\nfunc Answer() int {\n\treturn 1\n}\n"
	const edited = "package internal\n\nfunc Answer() int {\n\treturn 2\n}\n"

	session := harness.Start(t)
	session.Simulator.Scenario(slowScenario)
	session.Repository.Write("internal/added.go", reviewed)

	name := session.Review()
	session.WaitFor(name, runstate.Reviewing)
	worktree := session.Recorded(name).Worktree

	session.Repository.Write("internal/added.go", edited)
	if session.Repository.Read("internal/added.go") != edited {
		t.Fatal("the working tree was not edited")
	}

	if held := session.File(filepath.Join(worktree, "internal", "added.go")); held != reviewed {
		t.Errorf("the reviewer's copy moved under it:\n%s", held)
	}

	session.WaitForTerminal(name)

	invocation := session.Simulator.OnlyInvocation()
	if !strings.Contains(invocation.Prompt, "return 1") {
		t.Error("the reviewer was not given the snapshot's content")
	}
	if strings.Contains(invocation.Prompt, "return 2") {
		t.Error("the reviewer was given content the person wrote after the snapshot")
	}

	asked := session.File(session.Round(name, 1).PromptPath())
	if !strings.Contains(asked, "return 1") || strings.Contains(asked, "return 2") {
		t.Errorf("the prompt on disk is not anchored to the snapshot:\n%s", asked)
	}
}
