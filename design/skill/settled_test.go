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
		"Nothing carries across a revision, so what was settled is not carried over from the round before it",
	} {
		if !strings.Contains(document, said) {
			t.Errorf("the skill does not say %q, so a session has no rule for opening a later round with what is already settled", said)
		}
	}
}

func TestTheSkillStatesASettledAnswerGoesInTheSettledListRatherThanAmongTheDecisions(t *testing.T) {
	settled := sectionOf(t, skill.Document(), settledHeading)

	for _, said := range []string{
		"A settled answer goes in the document's `settled` list rather than among its decisions",
		"states the question, the labels the person chose and their own words",
		"An identifier is either settled or still open",
	} {
		if !strings.Contains(settled, said) {
			t.Errorf("the settled section does not say %q:\n%s", said, settled)
		}
	}

	for _, field := range []string{`"settled":`, `"chosen":`, `"words":`} {
		if !strings.Contains(settled, field) {
			t.Errorf("the settled section's own example carries no %s, and it is the shape a session writes:\n%s", field, settled)
		}
	}
	lowered := strings.ToLower(settled)
	for _, field := range []string{`"options"`, `"mode"`, `"decisions"`} {
		if strings.Contains(lowered, field) {
			t.Errorf("the settled section's own example carries %s, and a settled answer is not a decision:\n%s", field, settled)
		}
	}
}

func TestTheSkillStatesASettledAnswerIsReopenedThroughTheMessage(t *testing.T) {
	settled := sectionOf(t, skill.Document(), settledHeading)

	for _, said := range []string{
		"A settled answer offers no field of its own",
		"A person reopening one names it in the message the page carries for the whole round",
		"the next round asks it again as a question",
	} {
		if !strings.Contains(settled, said) {
			t.Errorf("the settled section does not say %q, so a design with no field on a settled answer would be a ratchet:\n%s", said, settled)
		}
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

func TestTheSkillStatesThePageReadsWhatWasSettledFirstWhereverItIsListed(t *testing.T) {
	settled := sectionOf(t, skill.Document(), settledHeading)

	if !strings.Contains(settled, "The page reads what was settled ahead of every open question, as a block that offers nothing to answer") {
		t.Errorf("the settled section does not say the page reads what was settled first:\n%s", settled)
	}
	if strings.Contains(settled, "the position you list it in is the position it renders in") {
		t.Errorf("the settled section still says list order places a settled card, and what was settled is no longer a card:\n%s", settled)
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
