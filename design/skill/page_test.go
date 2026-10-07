package skill_test

import (
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/design/skill"
)

const (
	roundsHeading = "## Rounds"
	pageHeading   = "## Putting a round on a page"
	factsHeading  = "## Facts"
)

func TestTheSkillStatesTheFrontierAndTheRound(t *testing.T) {
	document := skill.Document()

	for _, said := range []string{
		"The frontier is every question whose prerequisites are settled.",
		"A round is the whole frontier",
		"A question whose answer depends on one still open belongs to a later round",
	} {
		if !strings.Contains(document, said) {
			t.Errorf("the skill does not say %q, so a session has no rule for which questions belong in front of a person at once", said)
		}
	}
}

func TestTheAnswerArrivesChosenAndAtMostOneIsMarked(t *testing.T) {
	page := sectionOf(t, skill.Document(), pageHeading)

	for _, said := range []string{
		`"recommended": true`,
		"at most one is ever marked",
		"A document marking two is refused before anything is served",
		"it arrives chosen when the page renders",
	} {
		if !strings.Contains(page, said) {
			t.Errorf("the page section does not say %q, so a session cannot tell whether agreeing with the round costs one press:\n%s", said, page)
		}
	}
}

func TestTheSessionIsMintedOnceAndKeptForTheWholeDesign(t *testing.T) {
	page := sectionOf(t, skill.Document(), pageHeading)

	for _, said := range []string{
		"Mint the session once, with the first",
		"socrates agora serve",
		"so the person is given one address",
		"Keep that session for the whole design",
		"socrates agora serve --session",
	} {
		if !strings.Contains(page, said) {
			t.Errorf("the page section does not say %q, so a person could be handed a second address for a later round:\n%s", said, page)
		}
	}
}

func TestTheSkillSaysWhenAnIllustrationEarnsItsPlace(t *testing.T) {
	page := sectionOf(t, skill.Document(), pageHeading)

	for _, said := range []string{
		"A decision may also carry one illustration, painted by the document's own stylesheet",
		"Draw one where the choice is a shape",
		"Skip it where the choice is a name, a threshold or a policy",
		"An illustration that only repeats what the words already said costs the person attention and costs the skill its ceiling for nothing",
	} {
		if !strings.Contains(page, said) {
			t.Errorf("the page section does not say %q, so a session has no rule for when an illustration earns its place:\n%s", said, page)
		}
	}
}

func TestFactsAreFoundRatherThanAsked(t *testing.T) {
	facts := sectionOf(t, skill.Document(), factsHeading)

	for _, said := range []string{
		"Facts are found rather than asked for",
		"No question reaches the page that reading the repository would have answered",
	} {
		if !strings.Contains(facts, said) {
			t.Errorf("the facts section does not say %q, so a session could put a question on a page that reading the repository would have answered:\n%s", said, facts)
		}
	}
}

func sectionOf(t *testing.T, document, heading string) string {
	t.Helper()

	_, after, found := strings.Cut(document, heading)
	if !found {
		t.Fatalf("the skill has no %s section", heading)
	}

	depth := strings.Count(strings.Fields(heading)[0], "#")
	for line := range strings.SplitSeq(after, "\n") {
		if !strings.HasPrefix(line, "#") {
			continue
		}
		if strings.Count(strings.Fields(line)[0], "#") > depth {
			continue
		}
		if next := strings.Index(after, "\n"+line); next >= 0 {
			return after[:next]
		}
	}
	return after
}

func TestTheRoundsSectionComesBeforeThePageSection(t *testing.T) {
	document := skill.Document()

	rounds := strings.Index(document, roundsHeading)
	page := strings.Index(document, pageHeading)
	if rounds < 0 || page < 0 {
		t.Fatalf("the skill is missing one of %q or %q", roundsHeading, pageHeading)
	}
	if rounds > page {
		t.Error("the page section is stated before the rounds section explains what a round is")
	}
}

func TestTheSkillNumbersQuestionsFromOneAndNamesEachQuestionN(t *testing.T) {
	page := sectionOf(t, skill.Document(), pageHeading)

	for _, said := range []string{
		`{"id": "q2", "name": "Question 2",`,
		"Number questions from 1, never 0, as they are first asked",
		"identify each |q<N>| and name it |Question <N>|",
	} {
		said = strings.ReplaceAll(said, "|", "`")
		if !strings.Contains(page, said) {
			t.Errorf("the page section does not say %q, so a session titles its cards q0, q1 rather than Question 1, Question 2:\n%s", said, page)
		}
	}
}
