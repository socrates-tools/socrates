package e2e_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates-review/internal/corpus"
	"github.com/socrates-tools/socrates-review/internal/findings"
	"github.com/socrates-tools/socrates-review/internal/runstate"
	"github.com/socrates-tools/socrates-review/internal/testing/harness"
)

func TestTheCleanJourney(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(cleanScenario)

	session.Repository.Write("internal/parser/parse.go", "package parser\n\nfunc Parse(document string) error {\n\treturn nil\n}\n")
	before := session.Repository.Status()

	printed, failures, code := session.Socrates("review", "--wait", "--intent", "add a parser that accepts anything")
	if code != 0 {
		t.Fatalf("socrates-review review --wait exited %d for a clean review: %s\n%s", code, failures, printed)
	}

	name := firstLine(printed)
	run := session.Run(name)
	if run.State != runstate.Clean {
		t.Fatalf("the run ended as %s: %s", run.State, reason(run))
	}

	if _, _, statusCode := session.Socrates("status", name); statusCode != 0 {
		t.Errorf("socrates-review status of a clean run exited %d, want 0", statusCode)
	}
	if _, _, findingsCode := session.Socrates("findings", name); findingsCode != 0 {
		t.Errorf("socrates-review findings of a clean run exited %d, want 0", findingsCode)
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

	printed, failures, code := session.Socrates("review", "--wait")
	if code != 1 {
		t.Fatalf("socrates-review review --wait exited %d for a review that escalates: %s\n%s", code, failures, printed)
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

	readable, _, findingsCode := session.Socrates("findings", name)
	if findingsCode != 1 {
		t.Errorf("socrates-review findings of an escalated run exited %d, want 1", findingsCode)
	}
	for _, needed := range []string{"needs a person", "can be patched", "informational"} {
		if !strings.Contains(readable, needed) {
			t.Errorf("socrates-review findings does not say who %q is for:\n%s", needed, readable)
		}
	}
}

func TestTheJourneyOnARecordedAnswer(t *testing.T) {
	entry, err := corpus.Named("dialogue-layer-first-look")
	if err != nil {
		t.Fatalf("reading the recorded answer: %v", err)
	}

	session := harness.Start(t)
	session.Simulator.Scenario("default:\n  recorded: " + entry.Name + "\n")
	session.Repository.Write("internal/added.go", "package internal\n")

	printed, failures, code := session.Socrates("review", "--wait")
	if code != 4 {
		t.Fatalf("socrates-review review --wait exited %d for a recorded answer that parks: %s\n%s", code, failures, printed)
	}

	name := firstLine(printed)
	run := session.Run(name)
	if run.State != runstate.AwaitingResponse {
		t.Fatalf("the run ended as %s: %s", run.State, reason(run))
	}

	stored := strings.TrimSuffix(session.File(session.Round(name, 1).FindingsPath()), "\n")
	if stored != string(entry.Document) {
		t.Errorf("what a real reviewer answered did not reach disk intact: %d bytes recorded, %d stored", len(entry.Document), len(stored))
	}

	readable, _, findingsCode := session.Socrates("findings", name)
	if findingsCode != 4 {
		t.Errorf("socrates-review findings of a run parked on a recorded answer exited %d, want 4", findingsCode)
	}
	if !strings.Contains(readable, "needs a person") {
		t.Errorf("socrates-review findings does not say who the recorded findings are for:\n%s", readable)
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
		{behaviour: "error-result", names: "529"},
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

			printed, _, code := session.Socrates("findings", name)
			if code != 2 {
				t.Errorf("socrates-review findings of a failed run exited %d, want 2", code)
			}
			if !strings.Contains(printed, "no review completed") {
				t.Errorf("socrates-review findings does not say that no review completed:\n%s", printed)
			}
		})
	}
}

func TestNothingBranchesOnWhatTheAgentSaidWentWrong(t *testing.T) {
	accounts := []string{
		"API Error: 529 Overloaded",
		"Prompt is too long for the context window",
	}

	reported := make([]string, 0, len(accounts))
	for _, account := range accounts {
		session := harness.Start(t)
		session.Simulator.Scenario("default:\n  behaviour: error-result\n  text: \"" + account + "\"\n")
		session.Repository.Write("internal/added.go", "package internal\n")

		name := session.Review()
		run := session.WaitForTerminal(name)

		if run.State != runstate.Failed {
			t.Fatalf("the account %q ended the run as %s, and a failure is a failure whatever the agent said", account, run.State)
		}
		if !strings.Contains(reason(run), account) {
			t.Errorf("the run's reason is %q, without the account the agent gave", reason(run))
		}

		_, _, findingsCode := session.Socrates("findings", name)
		_, _, statusCode := session.Socrates("status", name)
		if findingsCode != 2 || statusCode != 2 {
			t.Errorf("the account %q exited findings %d and status %d, want 2 for both", account, findingsCode, statusCode)
		}
		reported = append(reported, strings.Replace(reason(run), account, "the account", 1))
	}

	if reported[0] != reported[1] {
		t.Errorf("two failures differing only in what the agent said read as %q and %q, so something branched on the account", reported[0], reported[1])
	}
}

func TestASilentReviewerIsAFailureJourney(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario("default:\n  behaviour: silence\n  hold: 30s\n")
	session.Repository.Write("internal/added.go", "package internal\n")

	path := configurationPath(t, session)
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
