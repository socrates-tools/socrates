package skill_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates-review/internal/dialogue"
	"github.com/socrates-tools/socrates-review/internal/runstate"
	"github.com/socrates-tools/socrates-review/internal/skill"
)

const (
	loopHeading = "## Procedure"
	inlineCode  = "`"
)

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

	if !strings.Contains(waiting, "socrates-review status") {
		t.Errorf("the step that waits for a review does not name socrates-review status:\n%s", waiting)
	}
	if !strings.Contains(waiting, "context") {
		t.Errorf("the step that waits for a review does not say why it polls rather than follows:\n%s", waiting)
	}
}

func TestTheWaitingStepSaysHowOftenAndForHowLong(t *testing.T) {
	waiting := strings.ToLower(numberedStep(t, skill.Document(), 3))

	if !strings.Contains(waiting, "five seconds") {
		t.Errorf("the step that waits for a review names no interval, so a session picks its own:\n%s", waiting)
	}
	for _, stated := range []string{"timeout", "deadline", string(runstate.Failed)} {
		if !strings.Contains(waiting, stated) {
			t.Errorf("the step that waits for a review does not say %q, so a session bounds the wait itself:\n%s", stated, waiting)
		}
	}
}

func TestWatchIsForThePerson(t *testing.T) {
	document := skill.Document()

	if named := strings.Count(document, "socrates-review watch"); named != 1 {
		t.Fatalf("the skill names socrates-review watch %d times, and it is a suggestion for the person rather than a step", named)
	}

	suggestion := sentenceNaming(document, "socrates-review watch")
	for _, stated := range []string{"person", "terminal"} {
		if !strings.Contains(suggestion, stated) {
			t.Errorf("the sentence naming socrates-review watch does not say %q, so it does not read as a suggestion for the person: %q", stated, suggestion)
		}
	}

	loop := strings.Index(document, loopHeading)
	if loop < 0 {
		t.Fatal("the skill has no loop, so there is nothing for the suggestion to stay out of")
	}
	if strings.Index(document, "socrates-review watch") > loop {
		t.Errorf("socrates-review watch is named inside the loop the session is told to run: %q", suggestion)
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

	reading := strings.Index(document, "socrates-review findings")
	if reading < 0 {
		t.Fatal("the skill never reads the findings, so there is nothing for the rules to come before")
	}
	if stated > reading {
		t.Error("the skill reads the findings before it states how to answer one, which is the moment the session is defending its own work")
	}
}

func TestAQuestionIsThePersons(t *testing.T) {
	said := sectionAbout(t, skill.Document(), string(runstate.AwaitingResponse))

	for _, stated := range []string{"needing a person", "is the person's", "record what they decided"} {
		if !strings.Contains(said, stated) {
			t.Errorf("the parked branch does not say %q, so a session reads a question as work of its own:\n%s", stated, said)
		}
	}
	if !strings.Contains(said, "You do not answer such a finding on your own behalf.") {
		t.Errorf("the parked branch does not say that the session does not answer a question itself, and nothing else stops it now that a question parks rather than ends the run:\n%s", said)
	}
	if !strings.Contains(said, "--decided") {
		t.Errorf("the parked branch does not say how what a person decided is recorded:\n%s", said)
	}
}

func TestARoundIsSpentWithoutAsking(t *testing.T) {
	said := sectionAbout(t, skill.Document(), string(runstate.AwaitingResponse))

	if !strings.Contains(said, "Do not ask the person whether the round is worth spending.") {
		t.Errorf("the parked branch does not say that the round is spent without asking:\n%s", said)
	}
	for _, stated := range []string{"unreviewed new code", "socrates-review continue"} {
		if !strings.Contains(said, stated) {
			t.Errorf("the parked branch does not say %q, so a session that answered everything has no reason to spend the round:\n%s", stated, said)
		}
	}
}

func TestWhatThePersonIsToldAtTheEnd(t *testing.T) {
	reporting := numberedStep(t, skill.Document(), 5)

	for _, stated := range []string{"how many rounds", "what the last one found"} {
		if !strings.Contains(reporting, stated) {
			t.Errorf("the step that reports a finished run does not say %q:\n%s", stated, reporting)
		}
	}
}

func TestAFurtherLookIsANewRun(t *testing.T) {
	said := sectionAbout(t, skill.Document(), string(runstate.Escalated))

	t.Run("There is no third round", func(t *testing.T) {
		if !strings.Contains(said, "No further round is possible.") {
			t.Errorf("the escalated branch does not say that no further round is possible, so a session looks for the flag that gives it one:\n%s", said)
		}
		if !strings.Contains(said, "fresh "+inlineCode+"socrates-review review") {
			t.Errorf("the escalated branch names no fresh review as the way to look again:\n%s", said)
		}
	})

	t.Run("A fresh review is the person's to ask for", func(t *testing.T) {
		for _, stated := range []string{"Offer that", "only where the person asks"} {
			if !strings.Contains(said, stated) {
				t.Errorf("the escalated branch does not say %q, so a session starts a review nobody asked for:\n%s", stated, said)
			}
		}
	})

	t.Run("A new run knows nothing of the old one", func(t *testing.T) {
		for _, stated := range []string{"told nothing about this run", "nothing about what it found", "nothing about the answers recorded against it"} {
			if !strings.Contains(said, stated) {
				t.Errorf("the escalated branch does not say %q, so the person does not know what a fresh review gives them:\n%s", stated, said)
			}
		}
	})
}

func TestTheLoopCloses(t *testing.T) {
	said := sectionAbout(t, skill.Document(), string(runstate.AwaitingResponse))

	for _, named := range []string{"socrates-review findings", "socrates-review respond", "socrates-review continue", "socrates-review status"} {
		if !strings.Contains(said, named) {
			t.Errorf("the section about a parked run does not name %s, so the loop does not close:\n%s", named, said)
		}
	}
}

func TestTheResponseRulesWriteTheirInvocationsAsCode(t *testing.T) {
	for _, rule := range dialogue.Rules() {
		spans := strings.Split(rule, inlineCode)
		for position := 0; position < len(spans); position += 2 {
			if strings.Contains(spans[position], skill.Name+" respond") {
				t.Errorf("a response rule names an invocation outside code, and the skill writes every invocation it names as code:\n%s", rule)
			}
		}
	}
}
