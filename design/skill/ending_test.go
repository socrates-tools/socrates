package skill_test

import (
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/design/skill"
)

const endingHeading = "## Ending early"

func TestTheSkillSaysADesignCanEndAtAnyRound(t *testing.T) {
	ending := sectionOf(t, skill.Document(), endingHeading)

	for _, said := range []string{
		"The person can end a design at any round, before the frontier empties",
		"socrates agora end --session <session>",
	} {
		if !strings.Contains(ending, said) {
			t.Errorf("the ending section does not say %q:\n%s", said, ending)
		}
	}
}

func TestEndingEarlyStillWritesFromWhatIsSettled(t *testing.T) {
	ending := sectionOf(t, skill.Document(), endingHeading)

	if !strings.Contains(ending, "Write the spec anyway") {
		t.Errorf("the ending section does not say to write the spec anyway:\n%s", ending)
	}
	if !strings.Contains(ending, "Everything settled up to that point becomes its answer in the spec exactly as it would on a normal end") {
		t.Errorf("the ending section does not say what was settled still becomes the spec's answer:\n%s", ending)
	}
}

func TestUnsettledQuestionsGoUnderAHeadingSayingTheDesignEndedEarly(t *testing.T) {
	ending := sectionOf(t, skill.Document(), endingHeading)

	for _, said := range []string{
		"is listed under a heading in Further Notes that says the design ended early",
		"### The design ended early",
	} {
		if !strings.Contains(ending, said) {
			t.Errorf("the ending section does not say %q, so an ended-early spec could not be told apart from one that finished:\n%s", said, ending)
		}
	}
}

func TestTheEndingSectionComesAfterTheFrontierEmptiesSection(t *testing.T) {
	document := skill.Document()

	frontier := strings.Index(document, frontierEmptiesHeading)
	ending := strings.Index(document, endingHeading)
	if frontier < 0 || ending < 0 {
		t.Fatalf("the skill is missing one of %q or %q", frontierEmptiesHeading, endingHeading)
	}
	if frontier > ending {
		t.Error("the ending section is stated before the section on the frontier emptying")
	}
}

func TestTheEndingSectionComesBeforeInstallingElsewhere(t *testing.T) {
	document := skill.Document()

	ending := strings.Index(document, endingHeading)
	installing := strings.Index(document, "## Installing elsewhere")
	if ending < 0 || installing < 0 {
		t.Fatalf("the skill is missing one of %q or %q", endingHeading, "## Installing elsewhere")
	}
	if ending > installing {
		t.Error("the ending section is stated after installing elsewhere")
	}
}
