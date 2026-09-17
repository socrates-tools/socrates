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

func TestPublishingIsOfferedOnTheLastCardWhereTheRepositoryStatesATracker(t *testing.T) {
	last := sectionOf(t, skill.Document(), lastRoundHeading)

	for _, said := range []string{
		"Offer publishing on this same card where the repository states an issue tracker",
		"docs/agents/issue-tracker.md",
		"State nothing about publishing where it does not",
	} {
		if !strings.Contains(last, said) {
			t.Errorf("the last round section does not say %q, so publishing could be offered where nothing takes it:\n%s", said, last)
		}
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
