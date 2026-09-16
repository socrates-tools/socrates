package e2e_test

import (
	"encoding/json"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/testing/harness"
	"github.com/socrates-tools/socrates/review/dialogue"
	"github.com/socrates-tools/socrates/review/prompt"
	"github.com/socrates-tools/socrates/review/report"
	"github.com/socrates-tools/socrates/review/runstate"
)

const answeredWrongly = "package internal\n\nfunc Answer(document string) int {\n\treturn len(document)\n}\n"

const answeredBetter = "package internal\n\nfunc Answer(document string) int {\n\tif document == \"\" {\n\t\treturn -1\n\t}\n\treturn len(document)\n}\n"

const patchableFindings = `{"findings":[{"id":"f1","severity":"error","file":"internal/added.go","line":3,"description":"Answer returns zero for an empty document, which is the same answer as for a document of length zero.","action":"auto-fix"},{"id":"f2","severity":"warning","file":"internal/added.go","line":3,"description":"The retry loop never backs off.","action":"auto-fix"},{"id":"f3","severity":"info","file":"internal/added.go","line":1,"description":"This package now owns two spellings of one idea.","action":"no-op"}],"risk_level":"medium","risk_rationale":"one changed path answers wrongly without erroring"}`

const nothingFoundAgain = `{"findings":[],"risk_level":"low","risk_rationale":"the patches answer what the first round found"}`

const foundAgain = `{"findings":[{"id":"g1","severity":"error","file":"internal/added.go","line":3,"description":"The rejected finding is still true of the code as it stands.","action":"auto-fix"}],"risk_level":"medium","risk_rationale":"what the author rejected is still there"}`

const unanswerableFinding = `{"findings":[{"id":"f1","severity":"error","file":"internal/added.go","line":3,"description":"The change takes away a behaviour the importer depends on, and nothing here can put it back.","action":"no-op"},{"id":"f2","severity":"warning","file":"internal/added.go","line":3,"description":"The retry loop never backs off.","action":"auto-fix"},{"id":"f3","severity":"info","file":"internal/added.go","line":1,"description":"This package now owns two spellings of one idea.","action":"no-op"}],"risk_level":"high","risk_rationale":"what holds the change is not something a patch closes"}`

const questionOnly = `{"findings":[{"id":"f1","severity":"error","file":"internal/added.go","line":3,"description":"This function does not look like something the product needs.","action":"ask-user"},{"id":"f2","severity":"info","file":"internal/added.go","line":1,"description":"This package now owns two spellings of one idea.","action":"no-op"}],"risk_level":"high","risk_rationale":"the change may not be wanted at all"}`

const questionAtTheCap = `{"findings":[{"id":"g1","severity":"error","file":"internal/added.go","line":3,"description":"The patch answers minus one for an empty document, and it is not clear the product wants a sentinel here at all.","action":"ask-user"}],"risk_level":"high","risk_rationale":"what the patch chose is a product decision nobody has made"}`

const questionOfIntent = `{"findings":[{"id":"f1","severity":"error","file":"internal/added.go","line":3,"description":"This function does not look like something the product needs.","action":"ask-user"},{"id":"f2","severity":"error","file":"internal/added.go","line":3,"description":"The guard is missing.","action":"auto-fix"},{"id":"f3","severity":"warning","file":"internal/added.go","line":3,"description":"The retry loop never backs off.","action":"auto-fix"}],"risk_level":"high","risk_rationale":"the change may not be wanted at all"}`

func dialogueScenario(first, second string) string {
	return "responses:\n  - when: \"" + prompt.ProvenanceHeading + "\"\n    structured: '" + second +
		"'\ndefault:\n  structured: '" + first + "'\n"
}

func parkedRun(t *testing.T, session *harness.Session, scenario string) string {
	t.Helper()

	session.Simulator.Scenario(scenario)
	session.Repository.Write("internal/added.go", answeredWrongly)

	name := session.Begin("--intent", "count the characters of a document")
	session.WaitFor(name, runstate.AwaitingResponse)
	return name
}

func respond(t *testing.T, session *harness.Session, arguments ...string) string {
	t.Helper()

	printed, failures, code := session.Socrates(append([]string{"respond"}, arguments...)...)
	if code != 4 {
		t.Fatalf("socrates-review respond %s exited %d: %s\n%s", strings.Join(arguments, " "), code, failures, printed)
	}
	return printed
}

func TestTheDialogueJourney(t *testing.T) {
	session := harness.Start(t)
	name := parkedRun(t, session, dialogueScenario(patchableFindings, nothingFoundAgain))

	t.Run("A parked run is not finished", func(t *testing.T) {
		for _, command := range []string{"status", "findings"} {
			printed, failures, code := session.Socrates(command, name)
			if code != 4 {
				t.Errorf("socrates %s of a parked run exited %d, want 4: %s", command, code, failures)
			}
			for _, verdict := range []string{string(runstate.Clean), string(runstate.Escalated), string(runstate.Failed)} {
				if strings.Contains(printed, "verdict: "+verdict) {
					t.Errorf("socrates %s of a parked run printed a verdict:\n%s", command, printed)
				}
			}
		}
	})

	t.Run("A parked run has nothing left to tear down", func(t *testing.T) {
		if _, err := os.Stat(snapshotWorktree(session, name)); !os.IsNotExist(err) {
			t.Errorf("the round's worktree is still at %s when the run reads parked: the next round checks out that exact path, and socrates-review continue races the teardown", snapshotWorktree(session, name))
		}
	})

	t.Run("A parked run holds no process identifiers", func(t *testing.T) {
		recorded := session.Recorded(name)
		if recorded.Reviewer != 0 || recorded.Agent != 0 {
			t.Errorf("a parked run records reviewer %d and agent %d, and both processes have exited: cancelling it later would signal whatever holds those numbers by then", recorded.Reviewer, recorded.Agent)
		}
	})

	t.Run("The rules arrive with the findings", func(t *testing.T) {
		printed, _, _ := session.Socrates("findings", name)
		for _, rule := range report.ResponseRules() {
			if !strings.Contains(printed, rule) {
				t.Errorf("socrates-review findings does not carry the rule %q:\n%s", rule, printed)
			}
		}
	})

	session.Repository.Write("internal/added.go", answeredBetter)
	respond(t, session, name, "--finding", "f1", "--fixed", "--note", "an empty document now answers minus one")
	printed := respond(t, session, name, "--finding", "f2", "--rejected", "--reason", "there is no retry loop in this function at all")

	if !strings.Contains(printed, "a response can address") {
		t.Errorf("socrates-review respond does not say the round is answered:\n%s", printed)
	}

	t.Run("Continuing when everything is answered", func(t *testing.T) {
		out, failures, code := session.Socrates("continue", name, "--wait")
		if code != 0 {
			t.Fatalf("socrates-review continue --wait exited %d: %s\n%s", code, failures, out)
		}
		if !strings.Contains(out, string(runstate.Clean)) {
			t.Errorf("the second round did not end the run clean:\n%s", out)
		}
	})

	run := session.Run(name)
	if run.State != runstate.Clean {
		t.Fatalf("the run ended as %s: %s", run.State, reason(run))
	}

	invocations := session.Simulator.Invocations()
	if len(invocations) != 2 {
		t.Fatalf("the reviewer was invoked %d times across two rounds, want twice", len(invocations))
	}

	t.Run("The claims travel explicitly", func(t *testing.T) {
		delivered := invocations[1].Prompt
		block, present := prompt.BodyOf(prompt.Sections(delivered), prompt.PriorHeading)
		if !present {
			t.Fatalf("the second round was given no prior-round block:\n%s", delivered)
		}
		for _, carried := range []string{
			"f1", "Answer returns zero for an empty document",
			"f2", "The retry loop never backs off.",
			"an empty document now answers minus one",
			"there is no retry loop in this function at all",
		} {
			if !strings.Contains(block, carried) {
				t.Errorf("the prior-round block does not carry %q:\n%s", carried, block)
			}
		}
		if strings.Count(delivered, "there is no retry loop in this function at all") != 1 {
			t.Error("the author's claim reached the reviewer outside the labelled block")
		}
	})

	t.Run("A second round is permitted", func(t *testing.T) {
		if _, err := os.Stat(session.Round(name, 2).PromptPath()); err != nil {
			t.Errorf("the second round kept no prompt: %v", err)
		}
	})
}

func TestWaitingReturnsWhenTheRunParks(t *testing.T) {
	// scenario: A parked run is not finished
	session := harness.Start(t)
	session.Simulator.Scenario(dialogueScenario(patchableFindings, nothingFoundAgain))
	session.Repository.Write("internal/added.go", answeredWrongly)

	printed, failures, code := session.Socrates("begin", "--wait")
	if code != 4 {
		t.Fatalf("socrates-review begin --wait exited %d for a run that parks: %s\n%s", code, failures, printed)
	}
	if !strings.Contains(printed, string(runstate.AwaitingResponse)) {
		t.Errorf("socrates-review begin --wait does not say the run is awaiting responses:\n%s", printed)
	}
	for _, verdict := range []string{string(runstate.Clean), string(runstate.Escalated), string(runstate.Failed)} {
		if strings.Contains(printed, "verdict: "+verdict) {
			t.Errorf("a parked run printed a verdict:\n%s", printed)
		}
	}
}

func TestTheRejectionAcceptedJourney(t *testing.T) {
	session := harness.Start(t)
	name := parkedRun(t, session, dialogueScenario(patchableFindings, nothingFoundAgain))

	respond(t, session, name, "--finding", "f1", "--rejected", "--reason", "an empty document is length zero, and zero is the right answer for it")
	respond(t, session, name, "--finding", "f2", "--rejected", "--reason", "there is no retry loop in this function at all")

	out, failures, code := session.Socrates("continue", name, "--wait")
	if code != 0 {
		t.Fatalf("socrates-review continue --wait exited %d after two rejections: %s\n%s", code, failures, out)
	}

	run := session.Run(name)
	if run.State != runstate.Clean {
		t.Fatalf("a second reviewer that accepted both rejections left the run %s: %s", run.State, reason(run))
	}

	second, present, err := findingsOf(session, name, 2)
	if err != nil {
		t.Fatalf("reading the second round's findings: %v", err)
	}
	if !present || len(second.Findings) != 0 {
		t.Errorf("the second round reported %d findings, want none", len(second.Findings))
	}

	first := session.File(session.Round(name, 1).FindingsPath())
	if !strings.Contains(first, "The retry loop never backs off.") {
		t.Error("the first round's findings were overwritten by the second round")
	}
}

func TestTheRejectionRefusedJourney(t *testing.T) {
	session := harness.Start(t)
	name := parkedRun(t, session, dialogueScenario(patchableFindings, foundAgain))

	respond(t, session, name, "--finding", "f1", "--fixed")
	respond(t, session, name, "--finding", "f2", "--rejected", "--reason", "the retry loop is bounded by the caller's deadline")

	out, failures, code := session.Socrates("continue", name, "--wait")
	if code != 1 {
		t.Fatalf("socrates-review continue --wait exited %d where the second round found something blocking: %s\n%s", code, failures, out)
	}

	run := session.Run(name)

	t.Run("A third round is not", func(t *testing.T) {
		if run.State != runstate.Escalated {
			t.Fatalf("the second round left the run %s, want escalated: %s", run.State, reason(run))
		}
		if !strings.Contains(reason(run), "cap") {
			t.Errorf("the verdict does not say the round cap was reached: %q", reason(run))
		}
		if _, err := os.Stat(session.Round(name, 3).Directory); !os.IsNotExist(err) {
			t.Error("a third round was started")
		}
	})

	t.Run("A disputed finding may be found again", func(t *testing.T) {
		second, present, err := findingsOf(session, name, 2)
		if err != nil || !present {
			t.Fatalf("reading the second round's findings: %v", err)
		}
		if len(second.Findings) != 1 || second.Findings[0].ID != "g1" {
			t.Fatalf("the second round reported %+v", second.Findings)
		}
		if !strings.Contains(second.Findings[0].Description, "still true of the code") {
			t.Errorf("the second round did not report the disputed finding again: %q", second.Findings[0].Description)
		}
	})

	t.Run("An escalated run keeps everything", func(t *testing.T) {
		for _, path := range []string{
			session.Round(name, 1).PromptPath(),
			session.Round(name, 1).StreamPath(),
			session.Round(name, 1).FindingsPath(),
			session.Round(name, 1).ResponsesPath(),
			session.Round(name, 2).PromptPath(),
			session.Round(name, 2).StreamPath(),
			session.Round(name, 2).FindingsPath(),
		} {
			written, err := os.Stat(path)
			if err != nil {
				t.Errorf("%s is not readable after the run escalated: %v", path, err)
				continue
			}
			if written.Size() == 0 {
				t.Errorf("%s is empty", path)
			}
		}
	})
}

func TestTheIntentQuestionJourney(t *testing.T) {
	// scenario: A question about intent stops the loop at once
	session := harness.Start(t)
	session.Simulator.Scenario(dialogueScenario(questionOfIntent, nothingFoundAgain))
	session.Repository.Write("internal/added.go", answeredWrongly)

	printed, failures, code := session.Socrates("begin", "--wait")
	if code != 4 {
		t.Fatalf("a round holding a question of intent exited %d: %s\n%s", code, failures, printed)
	}

	name := firstLine(printed)
	run := session.Run(name)
	if run.State != runstate.AwaitingResponse {
		t.Fatalf("a round holding a question of intent left the run %s: %s", run.State, reason(run))
	}

	readable, _, _ := session.Socrates("findings", name)
	for _, handed := range []string{"f1", "f2", "f3"} {
		if !strings.Contains(readable, handed) {
			t.Errorf("the person was not handed %s:\n%s", handed, readable)
		}
	}
	for _, rule := range report.ResponseRules() {
		if !strings.Contains(readable, rule) {
			t.Errorf("a run parked on a question does not carry the rule %q:\n%s", rule, readable)
		}
	}

	_, refused, respondCode := session.Socrates("respond", name, "--finding", "f1", "--fixed")
	if respondCode != 3 {
		t.Errorf("answering the question with a claimed patch exited %d, want 3", respondCode)
	}
	if !strings.Contains(refused, string(dialogue.Decided)) {
		t.Errorf("the refusal does not name the resolution that applies: %s", refused)
	}
	if _, err := os.Stat(session.Round(name, 1).ResponsesPath()); !os.IsNotExist(err) {
		t.Error("a patch claimed against a question stored something")
	}
}

func TestAnEscalatedRunSaysWhichOfTheTwoReasonsItWas(t *testing.T) {
	t.Run("a finding no response addresses", func(t *testing.T) {
		session := harness.Start(t)
		session.Simulator.Scenario(escalatingScenario)
		session.Repository.Write("internal/added.go", answeredWrongly)

		name := session.Begin()
		run := session.WaitForTerminal(name)
		if run.State != runstate.Escalated {
			t.Fatalf("the run ended as %s: %s", run.State, reason(run))
		}

		if !strings.Contains(reason(run), "asking for nothing") {
			t.Errorf("the verdict does not say that what escalated the run is a finding no response addresses: %q", reason(run))
		}
		if strings.Contains(reason(run), "cap") {
			t.Errorf("a run that escalated at its first round names the round cap: %q", reason(run))
		}

		printed, _, _ := session.Socrates("status", name)
		if !strings.Contains(printed, "asking for nothing") {
			t.Errorf("socrates-review status does not print which of the two reasons it was:\n%s", printed)
		}
	})

	t.Run("the round cap", func(t *testing.T) {
		session := harness.Start(t)
		name := parkedRun(t, session, dialogueScenario(patchableFindings, foundAgain))

		respond(t, session, name, "--finding", "f1", "--fixed")
		respond(t, session, name, "--finding", "f2", "--fixed")
		if _, failures, code := session.Socrates("continue", name, "--wait"); code != 1 {
			t.Fatalf("a second round that still found something blocking exited %d: %s", code, failures)
		}

		run := session.Run(name)
		if !strings.Contains(reason(run), "cap") {
			t.Errorf("the verdict does not say the round cap was reached: %q", reason(run))
		}
		if strings.Contains(reason(run), "asking for nothing") {
			t.Errorf("a run that escalated at the cap names a finding no response addresses: %q", reason(run))
		}
	})
}

func TestAnEscalatedRunTakesNoResponse(t *testing.T) {
	// scenario: Responding to a run that is not awaiting responses
	session := harness.Start(t)
	session.Simulator.Scenario(escalatingScenario)
	session.Repository.Write("internal/added.go", answeredWrongly)

	name := session.Begin()
	run := session.WaitForTerminal(name)
	if run.State != runstate.Escalated {
		t.Fatalf("a round holding a finding no response addresses left the run %s: %s", run.State, reason(run))
	}

	readable, _, findingsCode := session.Socrates("findings", name)
	if findingsCode != 1 {
		t.Errorf("socrates-review findings of an escalated run exited %d, want 1", findingsCode)
	}
	for _, rule := range report.ResponseRules() {
		if strings.Contains(readable, rule) {
			t.Errorf("an escalated run offered its findings for a response:\n%s", readable)
		}
	}

	_, refused, respondCode := session.Socrates("respond", name, "--finding", "f2", "--fixed")
	if respondCode != 3 {
		t.Errorf("socrates-review respond against an escalated run exited %d, want 3", respondCode)
	}
	if !strings.Contains(refused, string(runstate.Escalated)) {
		t.Errorf("the refusal does not name the run's state: %s", refused)
	}
	if _, err := os.Stat(session.Round(name, 1).ResponsesPath()); !os.IsNotExist(err) {
		t.Error("a response was stored against a run that escalated")
	}
}

func TestTheAuthorEditsBetweenRounds(t *testing.T) {
	session := harness.Start(t)
	name := parkedRun(t, session, dialogueScenario(patchableFindings, nothingFoundAgain))

	base := session.Recorded(name).Base
	firstSnapshot := session.Recorded(name).Snapshot

	session.Repository.Write("internal/added.go", answeredBetter)
	session.Repository.Write("internal/added_test.go", "package internal\n\nfunc TestAnEmptyDocument(t *testing.T) {}\n")
	respond(t, session, name, "--finding", "f1", "--fixed", "--note", "an empty document now answers minus one")
	respond(t, session, name, "--finding", "f2", "--fixed", "--note", "there was no loop, so nothing to back off")

	if _, failures, code := session.Socrates("continue", name, "--wait"); code != 0 {
		t.Fatalf("socrates-review continue --wait exited %d: %s", code, failures)
	}

	run := session.Recorded(name)
	first, hasFirst := run.Reviewed(1)
	second, hasSecond := run.Reviewed(2)
	if !hasFirst || !hasSecond {
		t.Fatalf("the run records %+v, want one entry per round", run.Reviews)
	}

	t.Run("Each round records what it reviewed", func(t *testing.T) {
		if first.Snapshot != firstSnapshot {
			t.Errorf("the first round records the snapshot %s, want %s", first.Snapshot, firstSnapshot)
		}
		if second.Snapshot == first.Snapshot {
			t.Error("the working tree changed between the rounds and both rounds record the same snapshot")
		}
	})

	t.Run("The base does not move between rounds", func(t *testing.T) {
		if first.Base != base || second.Base != base {
			t.Errorf("the rounds record the bases %s and %s, want %s", first.Base, second.Base, base)
		}
	})

	t.Run("The second round sees the patches", func(t *testing.T) {
		if held := session.Repository.Contents(second.Snapshot, "internal/added.go"); held != strings.TrimSpace(answeredBetter) {
			t.Errorf("the second round's snapshot does not hold the patched file:\n%s", held)
		}
		delivered := session.Simulator.Invocations()[1].Prompt
		if !strings.Contains(delivered, "+\t\treturn -1") {
			t.Errorf("the second round was not shown the patch as a diff:\n%s", delivered)
		}
	})

	t.Run("The second round sees the whole change", func(t *testing.T) {
		delivered := session.Simulator.Invocations()[1].Prompt
		for _, carried := range []string{base, second.Snapshot, "internal/added.go", "internal/added_test.go", "+func Answer(document string) int {"} {
			if !strings.Contains(delivered, carried) {
				t.Errorf("the second round's prompt does not carry %q:\n%s", carried, delivered)
			}
		}
	})

	t.Run("Two rounds, two directories", func(t *testing.T) {
		for _, path := range []string{
			session.Round(name, 1).PromptPath(),
			session.Round(name, 1).FindingsPath(),
			session.Round(name, 1).StreamPath(),
			session.Round(name, 1).ResponsesPath(),
			session.Round(name, 2).PromptPath(),
			session.Round(name, 2).FindingsPath(),
			session.Round(name, 2).StreamPath(),
		} {
			if _, err := os.Stat(path); err != nil {
				t.Errorf("%s is missing: %v", path, err)
			}
		}
		if _, err := os.Stat(session.Round(name, 2).ResponsesPath()); !os.IsNotExist(err) {
			t.Error("the first round's responses were written under the second round")
		}
		if first := session.File(session.Round(name, 1).PromptPath()); strings.Contains(first, prompt.PriorHeading) {
			t.Error("the first round's prompt was overwritten by the second round's")
		}
	})
}

func TestARoundWithNothingPatchedStillSnapshots(t *testing.T) {
	session := harness.Start(t)
	name := parkedRun(t, session, dialogueScenario(patchableFindings, nothingFoundAgain))

	respond(t, session, name, "--finding", "f1", "--rejected", "--reason", "zero is the right answer for an empty document")
	respond(t, session, name, "--finding", "f2", "--rejected", "--reason", "there is no retry loop in this function at all")

	before := session.Repository.Status()
	if _, failures, code := session.Socrates("continue", name, "--wait"); code != 0 {
		t.Fatalf("socrates-review continue --wait exited %d with nothing patched: %s", code, failures)
	}
	if after := session.Repository.Status(); after != before {
		t.Errorf("the repository changed:\nbefore:\n%s\nafter:\n%s", before, after)
	}

	run := session.Recorded(name)
	second, present := run.Reviewed(2)
	if !present {
		t.Fatalf("the second round recorded no snapshot: %+v", run.Reviews)
	}
	if second.Snapshot == "" {
		t.Error("the second round took no snapshot because nothing was patched")
	}

	delivered := session.Simulator.Invocations()[1].Prompt
	if !strings.Contains(delivered, second.Snapshot) {
		t.Errorf("the second round's reviewer was not anchored to the second round's snapshot:\n%s", delivered)
	}
	if session.Simulator.Invocations()[1].Directory == "" {
		t.Error("the second round's reviewer was given no working directory")
	}
}

func TestALeftoverWorktreeDoesNotStopTheNextRound(t *testing.T) {
	session := harness.Start(t)
	name := parkedRun(t, session, dialogueScenario(patchableFindings, nothingFoundAgain))

	session.Repository.Write("internal/added.go", answeredBetter)
	respond(t, session, name, "--finding", "f1", "--fixed", "--note", "an empty document now answers minus one")
	respond(t, session, name, "--finding", "f2", "--rejected", "--reason", "there is no retry loop in this function at all")

	leftover := snapshotWorktree(session, name)
	if err := os.MkdirAll(filepath.Join(leftover, "what", "the", "teardown", "left"), 0o755); err != nil {
		t.Fatalf("leaving a worktree behind: %v", err)
	}

	printed, failures, code := session.Socrates("continue", name, "--wait")
	if code != 0 {
		t.Fatalf("a round whose predecessor could not be torn down exited %d: %s\n%s", code, failures, printed)
	}
	if run := session.Run(name); run.State != runstate.Clean {
		t.Errorf("the run ended %s: %s", run.State, reason(run))
	}
}

func TestACleanSecondRoundDoesNotNameTheCap(t *testing.T) {
	const informationalOnly = `{"findings":[{"id":"g1","severity":"info","file":"internal/added.go","line":1,"description":"This package now owns two spellings of one idea.","action":"no-op"}],"risk_level":"low","risk_rationale":"nothing the patches introduced holds the change"}`

	session := harness.Start(t)
	name := parkedRun(t, session, dialogueScenario(patchableFindings, informationalOnly))

	session.Repository.Write("internal/added.go", answeredBetter)
	respond(t, session, name, "--finding", "f1", "--fixed", "--note", "an empty document now answers minus one")
	respond(t, session, name, "--finding", "f2", "--rejected", "--reason", "there is no retry loop in this function at all")

	if _, failures, code := session.Socrates("continue", name, "--wait"); code != 0 {
		t.Fatalf("a second round that found nothing blocking exited %d: %s", code, failures)
	}

	run := session.Run(name)
	if run.State != runstate.Clean {
		t.Fatalf("the run ended %s: %s", run.State, reason(run))
	}
	if strings.Contains(reason(run), "cap") {
		t.Errorf("a clean run's verdict names the round cap: %q", reason(run))
	}
}

func TestNothingLeftToReviewEndsTheRunWithoutPassingIt(t *testing.T) {
	session := harness.Start(t)
	name := parkedRun(t, session, dialogueScenario(patchableFindings, nothingFoundAgain))

	respond(t, session, name, "--finding", "f1", "--fixed", "--note", "I took the function out")
	respond(t, session, name, "--finding", "f2", "--fixed", "--note", "it went with the function")
	session.Repository.Delete("internal/added.go")

	printed, failures, code := session.Socrates("continue", name, "--wait")
	if code != 2 {
		t.Fatalf("continuing with the change reverted exited %d, want 2: %s\n%s", code, failures, printed)
	}

	run := session.Run(name)
	if run.State != runstate.Failed {
		t.Errorf("continuing with the change reverted ended the run %s, and no reviewer read the round: %s", run.State, reason(run))
	}
	if strings.Contains(reason(run), "HEAD") {
		t.Errorf("the recorded reason names HEAD, and what was compared is the run's base: %q", reason(run))
	}
	if !strings.Contains(reason(run), "no reviewer read this round") {
		t.Errorf("the recorded reason does not say that nobody read the round: %q", reason(run))
	}
	if invocations := session.Simulator.Invocations(); len(invocations) != 1 {
		t.Errorf("a second reviewer was invoked for a tree holding nothing to review: %d invocations", len(invocations))
	}
	if printed, _, _ := session.Socrates("findings", name); !strings.Contains(printed, "no review completed") {
		t.Errorf("what a person reads does not say that no review completed:\n%s", printed)
	}
}

func TestTheBaseStillExistsWhenThePersonHasMovedOn(t *testing.T) {
	session := harness.Start(t)
	name := parkedRun(t, session, dialogueScenario(patchableFindings, nothingFoundAgain))

	base := session.Recorded(name).Base
	respond(t, session, name, "--finding", "f1", "--fixed")
	respond(t, session, name, "--finding", "f2", "--fixed")

	session.Repository.Write("internal/added.go", answeredBetter)
	session.Repository.Commit("the author committed while the run waited")
	if session.Repository.Head() == base {
		t.Fatal("the person's commit did not move HEAD")
	}

	if _, failures, code := session.Socrates("continue", name, "--wait"); code != 0 {
		t.Fatalf("socrates-review continue --wait exited %d after the person committed: %s", code, failures)
	}

	second, present := session.Recorded(name).Reviewed(2)
	if !present || second.Base != base {
		t.Fatalf("the second round records the base %q, want the run's original base %s", second.Base, base)
	}

	delivered := session.Simulator.Invocations()[1].Prompt
	if !strings.Contains(delivered, "+func Answer(document string) int {") {
		t.Errorf("the second round was diffed against something later than the original base:\n%s", delivered)
	}
}

func TestABaseThatHasGoneFailsTheRound(t *testing.T) {
	// scenario: The base still exists when the person has moved on
	session := harness.Start(t)
	name := parkedRun(t, session, dialogueScenario(patchableFindings, nothingFoundAgain))

	base := session.Recorded(name).Base
	respond(t, session, name, "--finding", "f1", "--fixed")
	respond(t, session, name, "--finding", "f2", "--fixed")

	repository := session.Repository
	repository.Git("checkout", "--quiet", "--orphan", "rewritten")
	repository.Git("add", "--all")
	repository.Git("commit", "--quiet", "--message", "a history the original base is not in")
	repository.Git("branch", "--delete", "--force", "main")
	repository.Git("reflog", "expire", "--expire=now", "--expire-unreachable=now", "--all")
	repository.Git("gc", "--prune=now", "--quiet")

	if reachable(t, repository.Root, base) {
		t.Skipf("the base %s is still in the repository, so there is nothing to report as gone", base)
	}

	_, failures, code := session.Socrates("continue", name)
	if code != 2 {
		t.Errorf("socrates-review continue with a base that has gone exited %d, want 2: %s", code, failures)
	}
	if !strings.Contains(failures, base) {
		t.Errorf("the failure does not name the base that has gone: %s", failures)
	}

	run := session.Recorded(name)
	if run.State != runstate.Failed {
		t.Errorf("a round whose base has gone left the run %s", run.State)
	}
	if !strings.Contains(reason(run), base) {
		t.Errorf("the verdict does not name the base that has gone: %q", reason(run))
	}
}

func reachable(t *testing.T, root, commit string) bool {
	t.Helper()

	command := exec.Command("git", "-C", root, "rev-parse", "--verify", commit+"^{commit}")
	return command.Run() == nil
}

func TestNoSessionIsCarried(t *testing.T) {
	session := harness.Start(t)
	name := parkedRun(t, session, dialogueScenario(patchableFindings, nothingFoundAgain))

	respond(t, session, name, "--finding", "f1", "--fixed")
	respond(t, session, name, "--finding", "f2", "--fixed")
	if _, failures, code := session.Socrates("continue", name, "--wait"); code != 0 {
		t.Fatalf("socrates-review continue --wait exited %d: %s", code, failures)
	}

	invocations := session.Simulator.Invocations()
	if len(invocations) != 2 {
		t.Fatalf("the reviewer was invoked %d times, want twice", len(invocations))
	}

	first, second := invocations[0], invocations[1]
	if second.Process == first.Process {
		t.Error("the second round ran in the first round's process")
	}

	for _, resuming := range []string{"--resume", "--continue", "--session-id", "--fork-session", "--replay-user-messages"} {
		if second.Has(resuming) {
			t.Errorf("the second round was invoked with %s", resuming)
		}
	}
	identifier := sessionIn(t, session.File(session.Round(name, 1).StreamPath()))
	if identifier == "" {
		t.Fatal("the first round's transcript names no session, so this test would prove nothing")
	}
	for _, argument := range append(second.Arguments, second.Prompt) {
		if strings.Contains(argument, identifier) {
			t.Errorf("the second round was invoked naming the first round's session %s", identifier)
		}
	}

	said := strings.ToLower(second.Prompt)
	for _, referring := range []string{"resume", "conversation", "your earlier session", "the previous session", "continue where"} {
		if strings.Contains(said, referring) {
			t.Errorf("the second round's prompt refers to %q", referring)
		}
	}
	if strings.Contains(second.Prompt, first.Prompt) {
		t.Error("the second round's prompt carries the first round's prompt")
	}
}

func sessionIn(t *testing.T, transcript string) string {
	t.Helper()

	for line := range strings.SplitSeq(strings.TrimSpace(transcript), "\n") {
		var event struct {
			SessionID string `json:"session_id"`
		}
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			continue
		}
		if event.SessionID != "" {
			return event.SessionID
		}
	}
	return ""
}

func TestRespondingIsOneFindingAtATime(t *testing.T) {
	session := harness.Start(t)
	name := parkedRun(t, session, dialogueScenario(patchableFindings, nothingFoundAgain))

	t.Run("One finding per invocation", func(t *testing.T) {
		documented, _, code := session.Socrates("respond", "--help")
		if code != 0 {
			t.Fatalf("socrates-review respond --help exited %d", code)
		}
		for line := range strings.SplitSeq(documented, "\n") {
			fields := strings.Fields(line)
			for position, field := range fields {
				if !strings.HasPrefix(field, "--") || position+1 >= len(fields) {
					continue
				}
				if fields[position+1] == "stringArray" || fields[position+1] == "stringSlice" {
					t.Errorf("socrates-review respond takes several values for %s: %s", field, line)
				}
			}
		}

		respond(t, session, name, "--finding", "f1", "--finding", "f2", "--fixed")
		if recorded := readResponses(t, session, name); len(recorded.Recorded) != 1 {
			t.Errorf("naming two findings in one invocation recorded %d responses: %+v", len(recorded.Recorded), recorded)
		}

		_, failures, refused := session.Socrates("respond", "f1", "f2", "--fixed")
		if refused != 3 {
			t.Errorf("socrates-review respond took findings as arguments and exited %d, want 3", refused)
		}
		if !strings.Contains(failures, "--finding") {
			t.Errorf("the refusal does not say that the one finding is named with a flag: %s", failures)
		}
	})

	t.Run("A reason from standard input", func(t *testing.T) {
		const reason = "The first line of why this finding is wrong.\nThe second line, which a shell would have eaten.\n"

		printed, failures, code := session.SocratesReading(reason, "respond", name, "--finding", "f1", "--rejected", "--reason", "-")
		if code != 4 {
			t.Fatalf("socrates-review respond reading its reason from standard input exited %d: %s\n%s", code, failures, printed)
		}

		recorded := readResponses(t, session, name)
		answer, present := recorded.For("f1")
		if !present {
			t.Fatalf("nothing was stored against f1: %+v", recorded)
		}
		if answer.Reason != reason {
			t.Errorf("the reason reads %q, want %q verbatim", answer.Reason, reason)
		}
	})
}

func TestADecisionsReasonFromStandardInput(t *testing.T) {
	const decided = "Vic decided the function stays.\nThe migration reads it, and nothing else does.\n"

	session := harness.Start(t)
	name := parkedRun(t, session, dialogueScenario(questionOfIntent, nothingFoundAgain))

	printed, failures, code := session.SocratesReading(decided, "respond", name, "--finding", "f1", "--decided", "--reason", "-")
	if code != 4 {
		t.Fatalf("socrates-review respond --decided reading its reason from standard input exited %d: %s\n%s", code, failures, printed)
	}

	recorded := readResponses(t, session, name)
	answer, present := recorded.For("f1")
	if !present {
		t.Fatalf("nothing was stored against f1: %+v", recorded)
	}
	if answer.Resolution != dialogue.Decided {
		t.Errorf("f1 is recorded as %q, want %q", answer.Resolution, dialogue.Decided)
	}
	if answer.Reason != decided {
		t.Errorf("what the person decided reads %q, want %q verbatim", answer.Reason, decided)
	}

	t.Run("A reason is read only where it is asked for", func(t *testing.T) {
		_, failures, code := session.SocratesReading("bytes nobody meant as a reason", "respond", name, "--finding", "f1", "--decided")
		if code != 3 {
			t.Fatalf("a decision with no reason and a pipe on standard input exited %d, want 3: %s", code, failures)
		}
		if !strings.Contains(failures, "claim") {
			t.Errorf("the refusal does not say a decision is a claim: %s", failures)
		}
		if again := readResponses(t, session, name); len(again.Recorded) != 1 {
			t.Errorf("the run holds %d responses, so what was on standard input was stored as the person's decision", len(again.Recorded))
		}
	})
}

func TestARunParkedOnAQuestionIsNotAVerdictEither(t *testing.T) {
	session := harness.Start(t)
	name := parkedRun(t, session, dialogueScenario(questionOfIntent, nothingFoundAgain))

	for _, command := range []string{"status", "findings"} {
		printed, failures, code := session.Socrates(command, name)
		if code != 4 {
			t.Errorf("socrates %s of a run parked on a question exited %d, want 4: %s", command, code, failures)
		}
		for _, verdict := range []string{string(runstate.Clean), string(runstate.Escalated), string(runstate.Failed)} {
			if strings.Contains(printed, "verdict: "+verdict) {
				t.Errorf("socrates %s of a run parked on a question printed a verdict:\n%s", command, printed)
			}
		}
	}

	if recorded := session.Recorded(name); recorded.Verdict != nil {
		t.Errorf("a run parked on a question carries the verdict %+v", *recorded.Verdict)
	}

	printed, _, _ := session.Socrates("status", name)
	if !strings.Contains(printed, "needing a person rather than a patch: 1") {
		t.Errorf("the status of a run parked on a question does not say how many of what it waits on need a person:\n%s", printed)
	}
}

func TestChangingADecisionKeepsWhatWasSaidFirst(t *testing.T) {
	// scenario: Changing a response
	session := harness.Start(t)
	name := parkedRun(t, session, dialogueScenario(questionOfIntent, nothingFoundAgain))

	respond(t, session, name, "--finding", "f1", "--decided", "--reason", "Vic said to take the function out")
	printed := respond(t, session, name, "--finding", "f1", "--decided", "--reason", "Vic changed their mind, and the function stays")

	if !strings.Contains(printed, "replaced") {
		t.Errorf("socrates-review respond does not say that an earlier decision was replaced:\n%s", printed)
	}

	recorded := readResponses(t, session, name)
	if len(recorded.Recorded) != 2 {
		t.Errorf("the run holds %d responses, so the change of mind is not visible", len(recorded.Recorded))
	}
	answer, _ := recorded.For("f1")
	if answer.Reason != "Vic changed their mind, and the function stays" {
		t.Errorf("the decision the next reviewer receives reads %q", answer.Reason)
	}
	if !strings.Contains(session.File(session.Round(name, 1).ResponsesPath()), "take the function out") {
		t.Error("the decision that was replaced is gone from the record")
	}
}

func TestRespondRefusesWhatItCannotRecord(t *testing.T) {
	session := harness.Start(t)
	name := parkedRun(t, session, dialogueScenario(patchableFindings, nothingFoundAgain))

	t.Run("A rejection with no reason", func(t *testing.T) {
		for _, blank := range []string{"", "   "} {
			_, failures, code := session.Socrates("respond", name, "--finding", "f1", "--rejected", "--reason", blank)
			if code != 3 {
				t.Errorf("a rejection reasoned %q exited %d, want 3", blank, code)
			}
			if !strings.Contains(failures, "claim") {
				t.Errorf("the refusal does not say a rejection is a claim: %s", failures)
			}
			if _, err := os.Stat(session.Round(name, 1).ResponsesPath()); !os.IsNotExist(err) {
				t.Error("a rejection with no reason stored something")
			}
		}
	})

	t.Run("A flag that does not apply to the answer", func(t *testing.T) {
		for _, attempt := range []struct {
			named     string
			arguments []string
			names     string
		}{
			{named: "a fix reasoned", arguments: []string{"--fixed", "--reason", "I moved the guard above the loop"}, names: "--note"},
			{named: "a rejection noted", arguments: []string{"--rejected", "--note", "the loop is bounded by the caller's deadline"}, names: "--reason"},
			{named: "a decision noted", arguments: []string{"--decided", "--note", "Vic said to leave the endpoint alone"}, names: "--reason"},
			{named: "a fix and a decision at once", arguments: []string{"--fixed", "--decided", "--reason", "Vic said to leave it"}, names: "--decided"},
			{named: "no answer at all", arguments: []string{"--reason", "Vic said to leave it"}, names: "--decided"},
		} {
			_, failures, code := session.Socrates(append([]string{"respond", name, "--finding", "f1"}, attempt.arguments...)...)
			if code != 3 {
				t.Errorf("%s exited %d, want 3: what the author said would have been dropped without a word", attempt.named, code)
			}
			if !strings.Contains(failures, attempt.names) {
				t.Errorf("the refusal of %s does not name %s: %s", attempt.named, attempt.names, failures)
			}
			if _, err := os.Stat(session.Round(name, 1).ResponsesPath()); !os.IsNotExist(err) {
				t.Errorf("%s stored something", attempt.named)
			}
		}
	})

	t.Run("A reason is read only where it is asked for", func(t *testing.T) {
		_, failures, code := session.SocratesReading("bytes nobody meant as a reason", "respond", name, "--finding", "f1", "--rejected")
		if code != 3 {
			t.Fatalf("a rejection with no reason and a pipe on standard input exited %d, want 3: %s", code, failures)
		}
		if !strings.Contains(failures, "claim") {
			t.Errorf("the refusal does not say a rejection is a claim: %s", failures)
		}
		if _, err := os.Stat(session.Round(name, 1).ResponsesPath()); !os.IsNotExist(err) {
			t.Error("a rejection with no reason stored what was on standard input as the author's claim")
		}
	})

	t.Run("A response to a finding that does not exist", func(t *testing.T) {
		_, failures, code := session.Socrates("respond", name, "--finding", "f9", "--fixed")
		if code != 3 {
			t.Errorf("responding to a finding the round never produced exited %d, want 3", code)
		}
		if !strings.Contains(failures, "f9") {
			t.Errorf("the refusal does not name the identifier: %s", failures)
		}
		if _, err := os.Stat(session.Round(name, 1).ResponsesPath()); !os.IsNotExist(err) {
			t.Error("a response to a finding that does not exist stored something")
		}
	})

	t.Run("Changing a response", func(t *testing.T) {
		respond(t, session, name, "--finding", "f1", "--rejected", "--reason", "zero is the right answer for an empty document")
		printed := respond(t, session, name, "--finding", "f1", "--fixed", "--note", "on reflection the reviewer was right")

		if !strings.Contains(printed, "replaced") {
			t.Errorf("socrates-review respond does not say that an earlier response was replaced:\n%s", printed)
		}

		recorded := readResponses(t, session, name)
		answer, _ := recorded.For("f1")
		if answer.Resolution != dialogue.Fixed {
			t.Errorf("the later response is %q, want %q", answer.Resolution, dialogue.Fixed)
		}
		if len(recorded.Recorded) != 2 {
			t.Errorf("the run holds %d responses, so the change of mind is not visible", len(recorded.Recorded))
		}
		if !strings.Contains(session.File(session.Round(name, 1).ResponsesPath()), "zero is the right answer") {
			t.Error("the response that was replaced is gone from the record")
		}
	})
}

func TestRespondingToARunThatIsNotAwaitingResponses(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(cleanScenario)
	session.Repository.Write("internal/added.go", answeredWrongly)

	name := session.Begin()
	session.WaitForTerminal(name)

	_, failures, code := session.Socrates("respond", name, "--finding", "f1", "--fixed")
	if code != 3 {
		t.Errorf("socrates-review respond against a clean run exited %d, want 3", code)
	}
	if !strings.Contains(failures, string(runstate.Clean)) {
		t.Errorf("the refusal does not name the run's state: %s", failures)
	}
	if _, err := os.Stat(session.Round(name, 1).ResponsesPath()); !os.IsNotExist(err) {
		t.Error("a response was stored against a clean run")
	}
}

func TestContinuingTooEarly(t *testing.T) {
	session := harness.Start(t)
	name := parkedRun(t, session, dialogueScenario(patchableFindings, nothingFoundAgain))

	_, failures, code := session.Socrates("continue", name)
	if code != 3 {
		t.Errorf("socrates-review continue with nothing answered exited %d, want 3: %s", code, failures)
	}
	for _, named := range []string{"f1", "f2"} {
		if !strings.Contains(failures, named) {
			t.Errorf("the refusal does not name %s: %s", named, failures)
		}
	}
	if _, err := os.Stat(session.Round(name, 2).Directory); !os.IsNotExist(err) {
		t.Error("a second round was started with nothing answered")
	}
	if session.Recorded(name).State != runstate.AwaitingResponse {
		t.Errorf("refusing to continue moved the run to %s", session.Recorded(name).State)
	}

	t.Run("Non-blocking findings need no response", func(t *testing.T) {
		respond(t, session, name, "--finding", "f1", "--fixed")
		respond(t, session, name, "--finding", "f2", "--fixed")

		if _, failures, code := session.Socrates("continue", name, "--wait"); code != 0 {
			t.Fatalf("socrates-review continue with only the blocking findings answered exited %d: %s", code, failures)
		}
		if recorded := readResponses(t, session, name); len(recorded.Recorded) != 2 {
			t.Errorf("the informational finding was made to wait for a response: %+v", recorded)
		}
	})
}

func TestContinuingWithAQuestionUnanswered(t *testing.T) {
	session := harness.Start(t)
	name := parkedRun(t, session, dialogueScenario(questionOfIntent, nothingFoundAgain))

	session.Repository.Write("internal/added.go", answeredBetter)
	respond(t, session, name, "--finding", "f2", "--fixed", "--note", "the guard is above the loop now")
	respond(t, session, name, "--finding", "f3", "--rejected", "--reason", "there is no retry loop in this function at all")

	_, failures, code := session.Socrates("continue", name)
	if code != 3 {
		t.Fatalf("continuing with the question unanswered exited %d, want 3: %s", code, failures)
	}
	if !strings.Contains(failures, "f1") {
		t.Errorf("the refusal does not name the finding that asks a person: %s", failures)
	}
	if !strings.Contains(failures, "person") {
		t.Errorf("the refusal does not say that a person is what is waited on: %s", failures)
	}
	if _, err := os.Stat(session.Round(name, 2).Directory); !os.IsNotExist(err) {
		t.Error("a second round was started with the question unanswered")
	}
	if session.Recorded(name).State != runstate.AwaitingResponse {
		t.Errorf("refusing to continue moved the run to %s", session.Recorded(name).State)
	}

	respond(t, session, name, "--finding", "f1", "--decided", "--reason", "Vic decided the function stays, because the migration reads it")
	if _, failures, code := session.Socrates("continue", name, "--wait"); code != 0 {
		t.Fatalf("continuing with the question answered exited %d: %s", code, failures)
	}
	if run := session.Run(name); run.State != runstate.Clean {
		t.Errorf("the second round left the run %s: %s", run.State, reason(run))
	}
}

func TestAParkedRunCanStillBeCancelled(t *testing.T) {
	session := harness.Start(t)
	name := parkedRun(t, session, dialogueScenario(patchableFindings, nothingFoundAgain))
	respond(t, session, name, "--finding", "f1", "--fixed")

	printed, failures, code := session.Socrates("cancel", name)
	if code != 0 {
		t.Fatalf("socrates-review cancel of a parked run exited %d: %s", code, failures)
	}
	if !strings.Contains(printed, name) {
		t.Errorf("socrates-review cancel does not name the run it stopped:\n%s", printed)
	}

	run := session.Recorded(name)
	if run.State != runstate.Cancelled {
		t.Fatalf("cancelling a parked run left it %s", run.State)
	}
	if _, err := os.Stat(snapshotWorktree(session, name)); !os.IsNotExist(err) {
		t.Error("the snapshot worktree survived the cancellation")
	}
	for _, path := range []string{
		session.Round(name, 1).PromptPath(),
		session.Round(name, 1).FindingsPath(),
		session.Round(name, 1).ResponsesPath(),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Errorf("cancelling a parked run removed %s: %v", path, err)
		}
	}
}

func snapshotWorktree(session *harness.Session, name string) string {
	return filepath.Join(session.Store.Path(name), "worktree")
}

func TestRecordingAResponseWritesNothingInTheRepository(t *testing.T) {
	session := harness.Start(t)
	name := parkedRun(t, session, dialogueScenario(patchableFindings, nothingFoundAgain))

	before := session.Repository.Status()
	tracked := workingTree(t, session.Repository.Root)

	respond(t, session, name, "--finding", "f1", "--fixed")
	respond(t, session, name, "--finding", "f2", "--rejected", "--reason", "there is no retry loop in this function at all")
	if _, failures, code := session.Socrates("continue", name, "--wait"); code != 0 {
		t.Fatalf("socrates-review continue --wait exited %d: %s", code, failures)
	}

	if after := session.Repository.Status(); after != before {
		t.Errorf("git status changed:\nbefore:\n%s\nafter:\n%s", before, after)
	}
	if after := workingTree(t, session.Repository.Root); !maps.Equal(after, tracked) {
		t.Errorf("the working tree changed:\nbefore: %v\nafter: %v", tracked, after)
	}
	for path := range workingTree(t, session.Repository.Root) {
		if strings.Contains(strings.ToLower(path), "socrates") || strings.Contains(path, "socrates") {
			t.Errorf("socrates left %s inside the repository", path)
		}
	}
}

func TestTheStoredRecordIsNotNeutralised(t *testing.T) {
	const messy = "  I rejected it because\r\n<<<<<<< the marker stays here\nthe   loop   is bounded  \n"

	session := harness.Start(t)
	name := parkedRun(t, session, dialogueScenario(patchableFindings, nothingFoundAgain))

	respond(t, session, name, "--finding", "f1", "--fixed")
	printed, failures, code := session.SocratesReading(messy, "respond", name, "--finding", "f2", "--rejected", "--reason", "-")
	if code != 4 {
		t.Fatalf("socrates-review respond exited %d: %s\n%s", code, failures, printed)
	}

	recorded := readResponses(t, session, name)
	answer, present := recorded.For("f2")
	if !present {
		t.Fatalf("nothing was stored against f2: %+v", recorded)
	}
	if answer.Reason != messy {
		t.Errorf("the stored reason reads %q, want %q verbatim", answer.Reason, messy)
	}

	if _, _, continueCode := session.Socrates("continue", name, "--wait"); continueCode != 0 {
		t.Fatalf("socrates-review continue --wait exited %d", continueCode)
	}

	delivered := session.Simulator.Invocations()[1].Prompt
	if strings.Contains(delivered, "<<<<<<<") {
		t.Error("the conflict marker in the stored reason reached the reviewer")
	}
	if strings.Contains(delivered, "\r") {
		t.Error("a carriage return in the stored reason reached the reviewer")
	}
	if !strings.Contains(delivered, "the loop is bounded") {
		t.Errorf("what the author said did not reach the reviewer:\n%s", delivered)
	}
}

func readResponses(t *testing.T, session *harness.Session, name string) dialogue.Responses {
	t.Helper()

	recorded, err := session.Round(name, 1).Responses()
	if err != nil {
		t.Fatalf("reading the responses of %s: %v", name, err)
	}
	return recorded
}

func findingsOf(session *harness.Session, name string, round int) (findingsDocument, bool, error) {
	content, err := os.ReadFile(session.Round(name, round).FindingsPath())
	if err != nil {
		return findingsDocument{}, false, err
	}

	var document findingsDocument
	if err := json.Unmarshal(content, &document); err != nil {
		return findingsDocument{}, false, err
	}
	return document, true, nil
}

type findingsDocument struct {
	Findings []struct {
		ID          string `json:"id"`
		Description string `json:"description"`
	} `json:"findings"`
}
