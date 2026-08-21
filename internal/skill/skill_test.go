package skill_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/rui-ktei/socrates/internal/dialogue"
	"github.com/rui-ktei/socrates/internal/runstate"
	"github.com/rui-ktei/socrates/internal/skill"
)

const loopHeading = "## The loop"

func sectionHeading(subject string) string {
	return "### " + subject
}

func sectionAbout(t *testing.T, document, subject string) string {
	t.Helper()

	_, after, found := strings.Cut(document, sectionHeading(subject)+"\n")
	if !found {
		t.Fatalf("the skill has no section about %s:\n%s", subject, document)
	}

	body := after
	for _, boundary := range []string{"\n### ", "\n## "} {
		if within, _, found := strings.Cut(body, boundary); found {
			body = within
		}
	}
	return body
}

func numberedStep(t *testing.T, document string, number int) string {
	t.Helper()

	opening := "\n" + strconv.Itoa(number) + ". "
	_, after, found := strings.Cut(document, opening)
	if !found {
		t.Fatalf("the skill has no step %d:\n%s", number, document)
	}

	step, _, _ := strings.Cut(after, "\n"+strconv.Itoa(number+1)+". ")
	return step
}

func TestTheWaitingCommandIsStatus(t *testing.T) {
	waiting := numberedStep(t, skill.Document(), 3)

	if !strings.Contains(waiting, "scrat status") {
		t.Errorf("the step that waits for a review does not name scrat status:\n%s", waiting)
	}
	if !strings.Contains(waiting, "context") {
		t.Errorf("the step that waits for a review does not say why it polls rather than follows:\n%s", waiting)
	}
}

func TestWatchIsForThePerson(t *testing.T) {
	document := skill.Document()

	if named := strings.Count(document, "scrat watch"); named != 1 {
		t.Fatalf("the skill names scrat watch %d times, and it is a suggestion for the person rather than a step", named)
	}

	suggestion := sentenceNaming(document, "scrat watch")
	for _, stated := range []string{"person", "terminal"} {
		if !strings.Contains(suggestion, stated) {
			t.Errorf("the sentence naming scrat watch does not say %q, so it does not read as a suggestion for the person: %q", stated, suggestion)
		}
	}

	loop := strings.Index(document, loopHeading)
	if loop < 0 {
		t.Fatal("the skill has no loop, so there is nothing for the suggestion to stay out of")
	}
	if strings.Index(document, "scrat watch") > loop {
		t.Errorf("scrat watch is named inside the loop the session is told to run: %q", suggestion)
	}
}

func sentenceNaming(document, named string) string {
	before, after, _ := strings.Cut(document, named)

	opening := strings.LastIndexAny(before, ".\n")
	closing := strings.IndexAny(after, ".\n")
	if closing < 0 {
		closing = len(after) - 1
	}
	return strings.TrimSpace(before[opening+1:]) + named + after[:closing+1]
}

func TestEveryTerminalStateIsCovered(t *testing.T) {
	document := skill.Document()

	for _, state := range []runstate.State{runstate.Clean, runstate.Failed, runstate.Escalated, runstate.AwaitingResponse} {
		t.Run(string(state), func(t *testing.T) {
			said := sectionAbout(t, document, string(state))
			if len(strings.Fields(said)) < 15 {
				t.Errorf("the skill names %s and says next to nothing about what to do: %q", state, said)
			}
		})
	}
}

func TestEveryStateARunRestsInIsNamed(t *testing.T) {
	document := skill.Document()

	for _, state := range runstate.States() {
		if !state.Terminal() && state != runstate.AwaitingResponse {
			continue
		}
		if !strings.Contains(document, sectionHeading(string(state))) {
			t.Errorf("a run rests in %s and the skill says nothing about it, so a loop waiting for the states it names would poll through it", state)
		}
	}
}

func TestACancelledRunIsNotAVerdict(t *testing.T) {
	said := strings.ToLower(sectionAbout(t, skill.Document(), string(runstate.Cancelled)))

	for _, stated := range []string{"nothing was decided", "not report the change as reviewed", "on purpose"} {
		if !strings.Contains(said, stated) {
			t.Errorf("the section about a cancelled run does not say %q:\n%s", stated, said)
		}
	}
}

func TestAFailedRunIsStatedNotToBeAPass(t *testing.T) {
	said := sectionAbout(t, skill.Document(), string(runstate.Failed))

	for _, stated := range []string{"no review completed", "nothing has been checked", "not a clean"} {
		if !strings.Contains(strings.ToLower(said), stated) {
			t.Errorf("the section about a failed run does not say %q:\n%s", stated, said)
		}
	}
}

func TestFailedComesBeforeEscalated(t *testing.T) {
	// scenario: A failed run is stated not to be a pass
	document := skill.Document()

	failed := strings.Index(document, sectionHeading(string(runstate.Failed)))
	escalated := strings.Index(document, sectionHeading(string(runstate.Escalated)))
	if failed < 0 || escalated < 0 {
		t.Fatalf("the skill describes failed at %d and escalated at %d", failed, escalated)
	}
	if failed > escalated {
		t.Error("the skill describes escalated before failed, and a session scanning for what to do about findings reads past the state that has none")
	}
}

func TestTheRulesAreStatedBeforeTheFindingsAreRead(t *testing.T) {
	document := skill.Document()

	stated := 0
	for _, rule := range dialogue.Rules() {
		at := strings.Index(document, rule)
		if at < 0 {
			t.Fatalf("the skill does not state the rule %q", rule)
		}
		stated = max(stated, at)
	}

	reading := strings.Index(document, "scrat findings")
	if reading < 0 {
		t.Fatal("the skill never reads the findings, so there is nothing for the rules to come before")
	}
	if stated > reading {
		t.Error("the skill reads the findings before it states how to answer one, which is the moment the session is defending its own work")
	}
}

func TestTheLoopCloses(t *testing.T) {
	said := sectionAbout(t, skill.Document(), string(runstate.AwaitingResponse))

	for _, named := range []string{"scrat findings", "scrat respond", "scrat continue", "scrat status"} {
		if !strings.Contains(said, named) {
			t.Errorf("the section about a parked run does not name %s, so the loop does not close:\n%s", named, said)
		}
	}
}
