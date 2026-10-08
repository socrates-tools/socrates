package e2e_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/testing/harness"
	"github.com/socrates-tools/socrates/review/runstate"
	"github.com/socrates-tools/socrates/review/skill"
)

const awaitsAPerson = "awaits a person"

func TestSocratesSaysWhichRoundsWaitOnAPersonAndOnlyThose(t *testing.T) {
	t.Parallel()
	for _, driven := range []struct {
		named    string
		scenario string
		state    runstate.State
		waits    bool
	}{
		{named: "a round escalating for a finding no answer addresses", scenario: dialogueScenario(unanswerableFinding, nothingFoundAgain), state: runstate.Escalated, waits: true},
		{named: "a round parked on a question about intent", scenario: dialogueScenario(questionOfIntent, nothingFoundAgain), state: runstate.AwaitingResponse, waits: true},
		{named: "a round parked on findings the author patches", scenario: dialogueScenario(patchableFindings, nothingFoundAgain), state: runstate.AwaitingResponse, waits: false},
		{named: "a round that found nothing blocking", scenario: dialogueScenario(nothingFoundAgain, nothingFoundAgain), state: runstate.Clean, waits: false},
	} {
		t.Run(driven.named, func(t *testing.T) {
			t.Parallel()
			journey := driving(t, driven.scenario, answeredWrongly)
			name := journey.begin("count the characters of a document")

			state, _ := journey.waitedFor(name)
			if state != driven.state {
				t.Fatalf("%s left the run %s, want %s", driven.named, state, driven.state)
			}

			reported := journey.findings(name)
			if said := strings.Contains(reported, awaitsAPerson); said != driven.waits {
				t.Errorf("%s says the round waits on a person: %t, want %t:\n%s", driven.named, said, driven.waits, reported)
			}

			polled, _, _ := journey.run("status", name)
			if strings.Contains(polled, awaitsAPerson) {
				t.Errorf("the command a session polls in a loop carries what only the command read once needs:\n%s", polled)
			}

			journey.wroteNothingItself()
			journey.namedNothingTheSkillDoesNot()
		})
	}
}

func TestTheSurfaceDoesNotGrow(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)

	for command, flags := range surfaceOf(t, session) {
		for _, flag := range flags {
			if namesAPage(flag) {
				t.Errorf("socrates %s offers %s, and socrates renders, publishes and stores nothing", command, flag)
			}
		}
		if namesAPage(command) {
			t.Errorf("socrates offers %s, and socrates renders, publishes and stores nothing", command)
		}
	}

	session.Simulator.Scenario(dialogueScenario(unanswerableFinding, nothingFoundAgain))
	session.Repository.Write("internal/added.go", answeredWrongly)
	name := session.Begin("--intent", "count the characters of a document")
	session.WaitForTerminal(name)

	recorded, err := os.ReadFile(filepath.Join(session.Store.Path(name), "run.json"))
	if err != nil {
		t.Fatalf("reading what the run recorded: %v", err)
	}
	var read map[string]any
	if err := json.Unmarshal(recorded, &read); err != nil {
		t.Fatalf("reading back what the run recorded: %v", err)
	}
	for field := range read {
		if namesAPage(field) {
			t.Errorf("the run records %q, and a link socrates neither made nor can refresh is a field that goes stale", field)
		}
	}
}

func TestNoRoleSocratesSpawnsWritesForAReader(t *testing.T) {
	// scenario: Nothing socrates spawns drafts a page
	t.Parallel()
	session := harness.Start(t)
	session.Simulator.Scenario(dialogueScenario(questionOfIntent, nothingFoundAgain))
	session.Repository.Write("internal/added.go", answeredWrongly)

	name := session.Begin("--intent", "count the characters of a document")
	session.WaitFor(name, runstate.AwaitingResponse)

	invocations := session.Simulator.Invocations()
	if len(invocations) == 0 {
		t.Fatal("no reviewer was invoked, so there is no delivered prompt to read")
	}
	for _, invoked := range invocations {
		for _, asked := range []string{"page", "artifact", "html", "diagram", "plain english", "for the person", "summarise for", skill.Surface} {
			if strings.Contains(strings.ToLower(invoked.Prompt), asked) {
				t.Errorf("a prompt socrates delivered names %q, and nothing socrates spawns writes for a person:\n%s", asked, invoked.Prompt)
			}
		}
	}
}

func namesAPage(named string) bool {
	for _, page := range []string{"page", "artifact", "html", "publish", "render"} {
		if strings.Contains(strings.ToLower(named), page) {
			return true
		}
	}
	return false
}
