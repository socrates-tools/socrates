package skill_test

import (
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/design/skill"
)

const settledHeading = "## What was settled"

func TestTheSkillStatesRoundTwoOpensWithWhatWasSettled(t *testing.T) {
	document := skill.Document()

	for _, said := range []string{
		"Round two opens with what round one settled, above what is still open",
		"Nothing carries across a revision, so a settled card is not carried over from the round before it",
	} {
		if !strings.Contains(document, said) {
			t.Errorf("the skill does not say %q, so a session has no rule for opening a later round with what is already settled", said)
		}
	}
}

func TestTheSkillStatesASettledCardCarriesNoOptionsAndNoMode(t *testing.T) {
	settled := sectionOf(t, skill.Document(), settledHeading)

	for _, said := range []string{
		"A settled card states the question, the answer and the person's own words, and offers nothing",
		"Such a card carries no options and states no mode either",
	} {
		if !strings.Contains(settled, said) {
			t.Errorf("the settled section does not say %q:\n%s", said, settled)
		}
	}

	lowered := strings.ToLower(settled)
	if strings.Contains(lowered, `"options"`) {
		t.Errorf("the settled section's own example carries an options field, and a settled card offers nothing:\n%s", settled)
	}
	if strings.Contains(lowered, `"mode"`) {
		t.Errorf("the settled section's own example states a mode, and a settled card states none:\n%s", settled)
	}
}

func TestTheSkillStatesThePreambleCarriesOneLinePerSettledQuestion(t *testing.T) {
	settled := sectionOf(t, skill.Document(), settledHeading)

	if !strings.Contains(settled, "Carry one sentence per settled question into the preamble too") {
		t.Errorf("the settled section does not say the preamble carries one sentence per settled question:\n%s", settled)
	}
	if !strings.Contains(settled, `"preamble":`) {
		t.Errorf("the settled section shows no example preamble carrying a settled question's line:\n%s", settled)
	}
}

func TestTheSkillStatesIdentifiersAreStableForTheLifeOfTheDesign(t *testing.T) {
	settled := sectionOf(t, skill.Document(), settledHeading)

	for _, said := range []string{
		"Identifiers are stable for the life of the design",
		"q2 is q2 from the round it is first asked to the round the spec is written",
		"A person reopening something they settled two rounds ago names it",
	} {
		if !strings.Contains(settled, said) {
			t.Errorf("the settled section does not say %q, so a session has no reason to keep an identifier from round to round:\n%s", said, settled)
		}
	}
}

func TestTheSkillStatesTheSettledBlockStaysFirstAndContiguousBecauseOfListOrder(t *testing.T) {
	settled := sectionOf(t, skill.Document(), settledHeading)

	if !strings.Contains(settled, "List every settled card ahead of the round's open questions") {
		t.Error("the settled section does not say to list the settled cards ahead of the open questions")
	}
	if !strings.Contains(settled, "the position you list it in is the position it renders in") {
		t.Error("the settled section does not explain that list order is what keeps the block first and contiguous")
	}
}

func TestTheSettledSectionComesAfterThePageSectionAndBeforeFacts(t *testing.T) {
	document := skill.Document()

	page := strings.Index(document, pageHeading)
	settled := strings.Index(document, settledHeading)
	facts := strings.Index(document, factsHeading)
	if page < 0 || settled < 0 || facts < 0 {
		t.Fatalf("the skill is missing one of %q, %q or %q", pageHeading, settledHeading, factsHeading)
	}
	if page >= settled || settled >= facts {
		t.Error("the settled section is not stated after the page section and before the facts section")
	}
}
