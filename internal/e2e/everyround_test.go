package e2e_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/agora/surface"
	"github.com/socrates-tools/socrates/socrates-review/runstate"
)

func (j *journey) handedOver(name, title string, decisions []map[string]any) (reportedRound, string) {
	j.t.Helper()

	read := j.reported(name)
	return read, j.serve(documentFor(j.t, title, read, decisions))
}

func (j *journey) whereAPageWouldBe() string {
	j.t.Helper()
	return filepath.Join(j.serving().surfaceHome(), "surface.json")
}

func (j *journey) servedNothing() {
	j.t.Helper()

	if _, err := os.Stat(j.whereAPageWouldBe()); err == nil {
		j.t.Error("a run that ended without a result was handed over as a page, and the person is being answered a question they withdrew")
	}
}

func TestTheCheckForAPageThatWasNotServedReadsWhereAPageLands(t *testing.T) {
	// scenario: A run that ended without a result is handed over as nothing
	driven := driving(t, dialogueScenario(patchableFindings, nothingFoundAgain), answeredWrongly)
	name := driven.review("count the characters of a document")

	if state, _ := driven.waitedFor(name); state != runstate.AwaitingResponse {
		t.Fatalf("the round left the run %s, want a parked run", state)
	}

	where := driven.whereAPageWouldBe()
	if _, err := os.Stat(where); err == nil {
		t.Fatalf("%s exists before anything was served, so the check reads a path something else writes", where)
	}

	driven.handedOver(name, "what an independent reviewer said", nil)

	if _, err := os.Stat(where); err != nil {
		t.Fatalf("a page was served and %s holds nothing: %v\nthe check that a run gets no page reads a path a page never reaches, so it could not fail", where, err)
	}
}

func TestARoundThatCompletedAReviewIsHandedOverAsAPage(t *testing.T) {
	for _, driven := range []struct {
		named     string
		scenario  string
		state     runstate.State
		decisions func(*testing.T, reportedRound) []map[string]any
	}{
		{
			named:    "a clean run",
			scenario: dialogueScenario(nothingFoundAgain, nothingFoundAgain),
			state:    runstate.Clean,
		},
		{
			named:    "a parked round the author answers alone",
			scenario: dialogueScenario(patchableFindings, nothingFoundAgain),
			state:    runstate.AwaitingResponse,
		},
		{
			named:     "a parked round holding a question",
			scenario:  dialogueScenario(questionOfIntent, nothingFoundAgain),
			state:     runstate.AwaitingResponse,
			decisions: oneQuestionAsADecision,
		},
		{
			named:    "an escalated run",
			scenario: dialogueScenario(unanswerableFinding, nothingFoundAgain),
			state:    runstate.Escalated,
			decisions: func(_ *testing.T, read reportedRound) []map[string]any {
				return everyFindingAsADecision(read)
			},
		},
	} {
		t.Run(driven.named, func(t *testing.T) {
			journey := driving(t, driven.scenario, answeredWrongly)
			name := journey.review("count the characters of a document")

			state, _ := journey.waitedFor(name)
			if state != driven.state {
				t.Fatalf("%s left the run %s, want %s", driven.named, state, driven.state)
			}

			var decisions []map[string]any
			if driven.decisions != nil {
				decisions = driven.decisions(t, journey.reported(name))
			}
			read, address := journey.handedOver(name, "what an independent reviewer said", decisions)

			if err := surface.CheckLoopback(originOf(t, address)[len("http://"):]); err != nil {
				t.Fatalf("agora serve handed back %q: %v", address, err)
			}

			page := openingEither(t, address, len(decisions) > 0)

			var headed string
			page.evaluate(`document.querySelector('.heading h1').textContent`, &headed)
			if strings.TrimSpace(headed) != "what an independent reviewer said" {
				t.Errorf("the page the person was handed heads with %q", headed)
			}

			var listed []string
			page.evaluate(`Array.from(document.querySelectorAll('.entry .entry-name')).map((held) => held.textContent.trim())`, &listed)

			if len(listed) != len(read.Findings) {
				t.Errorf("the page lists %v, and socrates reported %d findings for the round", listed, len(read.Findings))
			}
			for _, finding := range read.Findings {
				if !strings.Contains(strings.Join(listed, "\n"), finding.name()) {
					t.Errorf("the page does not list %q, and it is what socrates reported for the round", finding.name())
				}
			}

			nothingLeftTheMachine(t, address, page.requests.all())
			journey.wroteNothingItself()
			journey.namedNothingTheSkillDoesNot()
		})
	}
}

func openingEither(t *testing.T, address string, answerable bool) *reader {
	t.Helper()

	if answerable {
		return opening(t, address)
	}
	return openingNothingToAnswer(t, address)
}

func TestARunThatEndedWithoutAResultIsHandedOverAsNothing(t *testing.T) {
	t.Run("a run that failed", func(t *testing.T) {
		driven := driving(t, "default:\n  behaviour: no-result\n", answeredWrongly)
		name := driven.review("count the characters of a document")

		state, _ := driven.waitedFor(name)
		if state != runstate.Failed {
			t.Fatalf("a reviewer that answered nothing left the run %s, want failed", state)
		}

		reported := driven.findings(name)
		if !strings.Contains(reported, "no review completed") {
			t.Errorf("what the session reads about a failed run does not say that no review completed:\n%s", reported)
		}

		driven.servedNothing()
		driven.wroteNothingItself()
		driven.namedNothingTheSkillDoesNot()
	})

	t.Run("a run stopped after a round that completed a review", func(t *testing.T) {
		driven := driving(t, dialogueScenario(questionOfIntent, nothingFoundAgain), answeredWrongly)
		name := driven.review("count the characters of a document")

		if state, _ := driven.waitedFor(name); state != runstate.AwaitingResponse {
			t.Fatalf("the round left the run %s, want a parked run", state)
		}
		if _, failures, code := driven.session.Socrates("cancel", name); code != 0 {
			t.Fatalf("somebody stopping the review exited %d: %s", code, failures)
		}
		if stopped := driven.session.Run(name); stopped.State != runstate.Cancelled {
			t.Fatalf("the cancelled run reads as %s", stopped.State)
		}

		printed, _, _ := driven.run("status", name)
		if !strings.Contains(printed, string(runstate.Cancelled)) {
			t.Errorf("what the session reads about a cancelled run does not say it was stopped:\n%s", printed)
		}

		driven.servedNothing()
		driven.wroteNothingItself()
		driven.namedNothingTheSkillDoesNot()
	})
}

func TestARoundThatWaitsOnNobodyIsHandedOverAndNotWaitedOn(t *testing.T) {
	for _, driven := range []struct {
		named    string
		scenario string
		state    runstate.State
	}{
		{named: "a clean run", scenario: dialogueScenario(nothingFoundAgain, nothingFoundAgain), state: runstate.Clean},
		{named: "a parked round the author answers alone", scenario: dialogueScenario(patchableFindings, nothingFoundAgain), state: runstate.AwaitingResponse},
	} {
		t.Run(driven.named, func(t *testing.T) {
			journey := driving(t, driven.scenario, answeredWrongly)
			name := journey.review("count the characters of a document")

			state, _ := journey.waitedFor(name)
			if state != driven.state {
				t.Fatalf("%s left the run %s, want %s", driven.named, state, driven.state)
			}
			if reported := journey.findings(name); strings.Contains(reported, awaitsAPerson) {
				t.Fatalf("%s is reported as one that waits on a person:\n%s", driven.named, reported)
			}

			_, address := journey.handedOver(name, "what an independent reviewer said", nil)
			if address == "" {
				t.Fatal("the address handed over is empty, and the person has nothing to open")
			}

			held := journey.serving()
			if delivered, arrived := held.awaiting(shortWindow); arrived {
				t.Errorf("a page nobody was asked to answer delivered %+v", delivered)
			}
		})
	}
}

func TestAPageNobodyWaitsOnDoesNotStopTheLoop(t *testing.T) {
	driven := driving(t, dialogueScenario(patchableFindings, nothingFoundAgain), answeredWrongly)
	name := driven.review("count the characters of a document")

	state, code := driven.waitedFor(name)
	if state != runstate.AwaitingResponse || code != 4 {
		t.Fatalf("the first round left the run %s and socrates-review status exited %d, want a parked run", state, code)
	}

	_, address := driven.handedOver(name, "what an independent reviewer said", nil)
	if address == "" {
		t.Fatal("the address handed over is empty, and the person has nothing to open")
	}

	driven.author("internal/added.go", answeredBetter)
	if _, failures, code := driven.run("respond", name, "--finding", "f1", "--fixed", "--note", "an empty document now answers minus one"); code != 4 {
		t.Fatalf("recording a fix while a page nobody waits on stands open exited %d: %s", code, failures)
	}
	if _, failures, code := driven.run("respond", name, "--finding", "f2", "--rejected", "--reason", "there is no retry loop in this function at all"); code != 4 {
		t.Fatalf("recording a rejection while a page nobody waits on stands open exited %d: %s", code, failures)
	}
	if _, failures, code := driven.run("continue", name); code != 0 {
		t.Fatalf("socrates-review continue while a page nobody waits on stands open exited %d: %s", code, failures)
	}

	state, code = driven.waitedFor(name)
	if state != runstate.Clean || code != 0 {
		t.Fatalf("the second round left the run %s and socrates-review status exited %d, want clean", state, code)
	}

	held := driven.serving()
	if delivered, arrived := held.awaiting(shortWindow); arrived {
		t.Errorf("the page nobody was asked to answer delivered %+v, and nobody ever opened it", delivered)
	}

	driven.namedNothingTheSkillDoesNot()
}
