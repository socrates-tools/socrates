package skill_test

import (
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/design/skill"
)

const lastRoundHeading = "## The last round"

func TestTheLastRoundIsOneCardNamingTheShapeWhatWillBeWrittenAndThePath(t *testing.T) {
	last := sectionOf(t, skill.Document(), lastRoundHeading)

	for _, said := range []string{
		"the last round is one card, and it is the only card on that page",
		"the shape of what was agreed, what will be written and the path it will be written to",
		"A person approves a location rather than discovers one",
		"docs/changes/<slug>/spec.md",
	} {
		if !strings.Contains(last, said) {
			t.Errorf("the last round section does not say %q, so a person could not approve a location before it is written:\n%s", said, last)
		}
	}
}

func TestTheLastCardOffersNoPublishingInAnyRepository(t *testing.T) {
	last := sectionOf(t, skill.Document(), lastRoundHeading)

	for _, unsaid := range []string{"issue", "tracker", "publish"} {
		if strings.Contains(strings.ToLower(last), unsaid) {
			t.Errorf("the last round section names %q, so a design could publish a spec somewhere other than the change's own folder:\n%s", unsaid, last)
		}
	}
	if !strings.Contains(last, "The card offers one option in every repository: write it") {
		t.Errorf("the last round section does not say the card offers the same two options in every repository:\n%s", last)
	}
}

func TestTheLastRoundHasNoRoundCap(t *testing.T) {
	last := sectionOf(t, skill.Document(), lastRoundHeading)

	for _, said := range []string{
		"There is no round cap",
		"The frontier decides when a design stops, and it stops once the frontier is empty",
	} {
		if !strings.Contains(last, said) {
			t.Errorf("the last round section does not say %q:\n%s", said, last)
		}
	}
}

func TestTheLastRoundSectionComesAfterTheDecisionsSection(t *testing.T) {
	document := skill.Document()

	decisions := strings.Index(document, decisionsHeading)
	last := strings.Index(document, lastRoundHeading)
	if decisions < 0 || last < 0 {
		t.Fatalf("the skill is missing one of %q or %q", decisionsHeading, lastRoundHeading)
	}
	if decisions > last {
		t.Error("the last round section is stated before the decisions section")
	}
}

func TestTheLastCardOffersWritingItAndNothingThatStartsAnotherSkill(t *testing.T) {
	last := sectionOf(t, skill.Document(), lastRoundHeading)

	if !strings.Contains(last, `"label": "write it", "recommended": true`) {
		t.Errorf("the last round section does not offer write it, recommended:\n%s", last)
	}
	for _, unsaid := range []string{"cut its tickets", skill.TicketsName, "cut the tickets", "tickets next"} {
		if strings.Contains(last, unsaid) {
			t.Errorf("the last round section names %q, so the last card could offer to start another skill:\n%s", unsaid, last)
		}
	}
	if count := strings.Count(last, `"label"`); count != 1 {
		t.Errorf("the last card offers %d options, want write it alone:\n%s", count, last)
	}
	if !strings.Contains(last, "Offer no option that starts another skill") {
		t.Errorf("the last round section does not forbid an option that starts another skill:\n%s", last)
	}
}
