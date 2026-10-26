package skill_test

import (
	skillkit "github.com/socrates-tools/socrates/internal/skill"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/design/skill"
)

const decisionsHeading = "## Writing the decisions"

func TestTheADRCriterionIsConjunctiveRatherThanDisjunctive(t *testing.T) {
	decisions := sectionOf(t, skill.Document(skillkit.ClaudeCode), decisionsHeading)

	if !strings.Contains(decisions, "all three of these hold together: it is hard to reverse, it is surprising without context, and it is the result of a real trade-off") {
		t.Errorf("the decisions section does not state the three-part criterion joined by and:\n%s", decisions)
	}
	if !strings.Contains(decisions, "Not one or two of the three, all three") {
		t.Errorf("the decisions section does not rule out one or two of the three being enough:\n%s", decisions)
	}

	forbidden := []string{
		"hard to reverse, surprising without context, or the result of a real trade-off",
		"hard to reverse, surprising without context or the result of a real trade-off",
	}
	for _, disjunctive := range forbidden {
		if strings.Contains(decisions, disjunctive) {
			t.Errorf("the decisions section states the criterion disjunctively: %q, and with or almost every decision qualifies", disjunctive)
		}
	}
}

func TestTheADRsHomeAndNumberingAreStated(t *testing.T) {
	decisions := sectionOf(t, skill.Document(skillkit.ClaudeCode), decisionsHeading)

	for _, said := range []string{
		"Number the file after the highest number already in",
		"docs/adr/",
		"docs/adr/NNNN-<sentence>.md",
		"Status: accepted",
		"Part of: <slug>",
	} {
		if !strings.Contains(decisions, said) {
			t.Errorf("the decisions section does not say %q:\n%s", said, decisions)
		}
	}
}

func TestTheRejectedSectionIsBuiltFromDeclinedOptionsOnly(t *testing.T) {
	decisions := sectionOf(t, skill.Document(skillkit.ClaudeCode), decisionsHeading)

	for _, said := range []string{
		"Its Rejected section is built from the options the person declined, and from nothing else",
		"the reason it lost is in the words they wrote beside the choice",
		"An option nobody was offered never appears there",
	} {
		if !strings.Contains(decisions, said) {
			t.Errorf("the decisions section does not say %q, so the Rejected section could be reconstructed rather than recorded:\n%s", said, decisions)
		}
	}
}

func TestTheDecisionsSectionComesAfterTheSpecSection(t *testing.T) {
	document := skill.Document(skillkit.ClaudeCode)

	spec := strings.Index(document, specHeading)
	decisions := strings.Index(document, decisionsHeading)
	if spec < 0 || decisions < 0 {
		t.Fatalf("the skill is missing one of %q or %q", specHeading, decisionsHeading)
	}
	if spec > decisions {
		t.Error("the decisions section is stated before the spec section")
	}
}

func TestADecisionThatFailsTheCriterionStaysInTheSpec(t *testing.T) {
	// scenario: A decision is easy to reverse
	decisions := sectionOf(t, skill.Document(skillkit.ClaudeCode), decisionsHeading)

	if !strings.Contains(decisions, "A decision that fails even one of the three stays inside the spec's Implementation Decisions and gets no file of its own.") {
		t.Errorf("the decisions section does not say where a decision without an ADR goes:\n%s", decisions)
	}
}
