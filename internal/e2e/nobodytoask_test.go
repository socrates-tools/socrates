package e2e_test

import (
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/testing/harness"
	"github.com/socrates-tools/socrates/review/runstate"
)

const nobodyToAskIntent = "count the characters of a document"

func beginWithNobodyToAsk(t *testing.T, session *harness.Session, scenario string, more ...string) string {
	t.Helper()

	session.Simulator.Scenario(scenario)
	session.Repository.Write("internal/added.go", answeredWrongly)

	arguments := append([]string{"begin", "--nobody-to-ask", "--intent", nobodyToAskIntent}, more...)
	printed, failures, code := session.Socrates(arguments...)
	if code != 0 && code != 4 {
		t.Fatalf("socrates review %s exited %d: %s\n%s", strings.Join(arguments, " "), code, failures, printed)
	}
	name, _, _ := strings.Cut(strings.TrimSpace(printed), "\n")
	if name == "" {
		t.Fatal("socrates review begin printed no run identifier")
	}
	return name
}

func TestARunWithNobodyToAskIsSettledByTheSessionThatBeganIt(t *testing.T) {
	t.Parallel()
	harness.EachVendor(t, func(t *testing.T, session *harness.Session) {
		name := beginWithNobodyToAsk(t, session, dialogueScenario(questionOnly, nothingFoundAgain))
		session.WaitFor(name, runstate.AwaitingResponse)

		t.Run("The status says nobody is there to ask", func(t *testing.T) {
			printed, failures, code := session.Socrates("status", name)
			if code != 4 {
				t.Fatalf("socrates review status of a parked run exited %d, want 4: %s", code, failures)
			}
			if !strings.Contains(printed, "nobody to ask") {
				t.Errorf("the status of a run begun with --nobody-to-ask does not say nobody is there to ask:\n%s", printed)
			}
			for _, aPerson := range []string{"needing a person rather than a patch", "or for a person the author asks"} {
				if strings.Contains(printed, aPerson) {
					t.Errorf("the status of a run with nobody to ask says it waits on a person (%q):\n%s", aPerson, printed)
				}
			}
		})

		t.Run("findings --json says the run has nobody to ask and does not await a person", func(t *testing.T) {
			read := reportedFor(t, session, name)
			if !read.NobodyToAsk {
				t.Errorf("socrates review findings --json does not say the run has nobody to ask: %+v", read)
			}
			if read.AwaitsAPerson {
				t.Errorf("socrates review findings --json reports a round holding only a question as awaiting a person on a run with nobody to ask: %+v", read)
			}
		})

		t.Run("A decision is refused, naming the settlement it takes", func(t *testing.T) {
			refused := refuse(t, session, "respond", name, "--finding", "f1", "--decided", "--reason", "Vic said the function stays")
			if !strings.Contains(refused, "--settled") {
				t.Errorf("the refusal does not name --settled as what the finding takes: %s", refused)
			}
		})

		t.Run("A settlement with no reason is refused", func(t *testing.T) {
			refused := refuse(t, session, "respond", name, "--finding", "f1", "--settled")
			if !strings.Contains(refused, "settlement") || !strings.Contains(refused, "needs a reason") {
				t.Errorf("the refusal does not say a settlement needs a reason: %s", refused)
			}
		})

		t.Run("Continuing before the question is settled is refused, naming the settlement", func(t *testing.T) {
			refused := refuse(t, session, "continue", name)
			if !strings.Contains(refused, "--settled") {
				t.Errorf("the refusal does not name --settled as how the finding is answered: %s", refused)
			}
		})

		const settled = "the function stays, because the importer reads it"
		respond(t, session, name, "--finding", "f1", "--settled", "--reason", settled)

		t.Run("The settlement is recorded as settled and reported as such", func(t *testing.T) {
			recorded := readResponses(t, session, name)
			answer, present := recorded.For("f1")
			if !present {
				t.Fatalf("nothing was stored against f1: %+v", recorded)
			}
			if string(answer.Resolution) != "settled" || answer.Reason != settled {
				t.Errorf("f1 is recorded as %q with the reason %q, want settled with %q", answer.Resolution, answer.Reason, settled)
			}

			var reported string
			for _, finding := range reportedFor(t, session, name).Findings {
				if finding.ID == "f1" {
					reported = finding.Resolution
				}
			}
			if reported != "settled" {
				t.Errorf("socrates review findings --json reports f1's resolution as %q, want settled", reported)
			}
		})

		t.Run("Continuing starts the second round", func(t *testing.T) {
			printed, failures, code := session.Socrates("continue", name, "--wait")
			if code != 0 {
				t.Fatalf("socrates review continue --wait after the settlement exited %d: %s\n%s", code, failures, printed)
			}
			if invocations := session.Simulator.Invocations(); len(invocations) != 2 {
				t.Fatalf("the reviewer was invoked %d times, want a second round after the settlement", len(invocations))
			}
		})

		t.Run("The second round is told the session decided it and is never told a person did", func(t *testing.T) {
			judging := session.Simulator.Invocations()[1].Prompt
			if !strings.Contains(judging, `"author_response": "The session that started this run decided this finding`) || !strings.Contains(judging, `"author_said": "`+settled+`"`) {
				t.Errorf("the round that judges the settlement was not told the session made it, with its reason:\n%s", judging)
			}
			if strings.Contains(strings.ToLower(judging), "a person answered") {
				t.Errorf("the round that judges the settlement was told a person answered:\n%s", judging)
			}
			if !strings.Contains(judging, "Judge whether the code now matches what was settled") {
				t.Errorf("the round that judges the settlement was not asked to judge the code against what was settled:\n%s", judging)
			}
		})
	})
}

func TestBeginWaitOnARunWithNobodyToAskDoesNotReportItParkedOnAPerson(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	session.Simulator.Scenario(dialogueScenario(questionOnly, nothingFoundAgain))
	session.Repository.Write("internal/added.go", answeredWrongly)

	printed, failures, code := session.Socrates("begin", "--nobody-to-ask", "--wait", "--intent", nobodyToAskIntent)
	if code != 4 {
		t.Fatalf("socrates review begin --nobody-to-ask --wait on a round holding only a question exited %d, want the awaiting-response status 4: %s\n%s", code, failures, printed)
	}
	if !strings.Contains(printed, string(runstate.AwaitingResponse)) {
		t.Errorf("begin --wait does not report the run as awaiting responses:\n%s", printed)
	}
	for _, aPerson := range []string{"needing a person rather than a patch", "or for a person the author asks"} {
		if strings.Contains(printed, aPerson) {
			t.Errorf("begin --wait reports a run with nobody to ask as parked on a person (%q):\n%s", aPerson, printed)
		}
	}
}

func TestASettlementIsRefusedOnARunThatHasSomebodyToAsk(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	name := parkedRun(t, session, dialogueScenario(questionOnly, nothingFoundAgain))

	refused := refuse(t, session, "respond", name, "--finding", "f1", "--settled", "--reason", "the function stays")
	if !strings.Contains(refused, "--decided") {
		t.Errorf("the refusal does not name --decided as what the finding takes: %s", refused)
	}

	respond(t, session, name, "--finding", "f1", "--decided", "--reason", "Vic said the function stays")
	recorded := readResponses(t, session, name)
	answer, present := recorded.For("f1")
	if !present || string(answer.Resolution) != "decided" {
		t.Errorf("a decision on a run with somebody to ask is recorded as %+v, want decided", answer)
	}

	if printed, _, _ := session.Socrates("status", name); strings.Contains(printed, "nobody to ask") {
		t.Errorf("the status of a run begun without --nobody-to-ask says nobody is there to ask:\n%s", printed)
	}
}

func TestNobodyToAskBesideNotMyChangeIsRefusedBeforeAnythingHappens(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	session.Repository.Write("internal/added.go", "package internal\n")

	head := session.Repository.Head()
	status := session.Repository.Status()
	objects := session.Repository.Git("count-objects", "-v")

	_, failures, code := session.Socrates("begin", "--nobody-to-ask", "--not-my-change")
	if code != 3 {
		t.Fatalf("socrates review begin --nobody-to-ask --not-my-change exited %d, want the usage code: %s", code, failures)
	}
	for _, named := range []string{"--nobody-to-ask", "--not-my-change"} {
		if !strings.Contains(failures, named) {
			t.Errorf("the refusal does not name %s: %s", named, failures)
		}
	}

	if names := session.Names(); len(names) != 0 {
		t.Errorf("the refusal created the run directories %v", names)
	}
	if invocations := session.Simulator.Invocations(); len(invocations) != 0 {
		t.Errorf("the refusal spawned %d reviewers", len(invocations))
	}
	if now := session.Repository.Git("count-objects", "-v"); now != objects {
		t.Errorf("the refusal wrote objects into the repository, so a snapshot was taken:\n%s\nwant:\n%s", now, objects)
	}
	if now := session.Repository.Head(); now != head {
		t.Errorf("the refusal moved HEAD from %s to %s", head, now)
	}
	if now := session.Repository.Status(); now != status {
		t.Errorf("git status reads %q after the refusal, want %q", now, status)
	}
}
