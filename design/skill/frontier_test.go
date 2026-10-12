package skill_test

import (
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/design/skill"
)

const frontierEmptiesHeading = "## When the frontier empties"

func TestAnEmptyFrontierWritesTheSpecWithoutAFurtherPage(t *testing.T) {
	frontier := sectionOf(t, skill.Document(), frontierEmptiesHeading)

	for _, said := range []string{
		"there is nothing left to ask, so serve no further page",
		"socrates agora end --session <session>",
		"write the spec and its decisions straight away",
		"docs/changes/<slug>/spec.md",
		"say in the terminal every path you wrote",
	} {
		if !strings.Contains(frontier, said) {
			t.Errorf("the frontier section does not say %q, so a person could be asked to approve a spec nothing is left to decide about:\n%s", said, frontier)
		}
	}
	if strings.Contains(frontier, `"label"`) {
		t.Errorf("the frontier section still composes a card for the person to approve:\n%s", frontier)
	}
}

func TestTheSpecIsPublishedInNoRepository(t *testing.T) {
	frontier := sectionOf(t, skill.Document(), frontierEmptiesHeading)

	for _, unsaid := range []string{"issue", "tracker", "publish"} {
		if strings.Contains(strings.ToLower(frontier), unsaid) {
			t.Errorf("the frontier section names %q, so a design could publish a spec somewhere other than the change's own folder:\n%s", unsaid, frontier)
		}
	}
	if !strings.Contains(frontier, "in the change's own folder and nowhere else") {
		t.Errorf("the frontier section does not keep the spec in the change's own folder:\n%s", frontier)
	}
}

func TestADesignHasNoRoundCap(t *testing.T) {
	frontier := sectionOf(t, skill.Document(), frontierEmptiesHeading)

	for _, said := range []string{
		"There is no round cap",
		"The frontier decides when a design stops, and it stops once the frontier is empty",
	} {
		if !strings.Contains(frontier, said) {
			t.Errorf("the frontier section does not say %q:\n%s", said, frontier)
		}
	}
}

func TestTheFrontierEmptiesSectionComesAfterTheDecisionsSection(t *testing.T) {
	document := skill.Document()

	decisions := strings.Index(document, decisionsHeading)
	frontier := strings.Index(document, frontierEmptiesHeading)
	if decisions < 0 || frontier < 0 {
		t.Fatalf("the skill is missing one of %q or %q", decisionsHeading, frontierEmptiesHeading)
	}
	if decisions > frontier {
		t.Error("the frontier section is stated before the decisions section")
	}
}

func TestAWrittenSpecStartsNoOtherSkill(t *testing.T) {
	frontier := sectionOf(t, skill.Document(), frontierEmptiesHeading)

	for _, unsaid := range []string{"cut its tickets", skill.TicketsName, "cut the tickets", "tickets next"} {
		if strings.Contains(frontier, unsaid) {
			t.Errorf("the frontier section names %q, so a design could start another skill:\n%s", unsaid, frontier)
		}
	}
	if !strings.Contains(frontier, "Start no other skill") {
		t.Errorf("the frontier section does not forbid starting another skill:\n%s", frontier)
	}
}
