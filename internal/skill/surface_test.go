package skill_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/skill"
)

var testNamespaces = map[string]string{"one": "binary-one", "two": "binary-two"}

func TestNamedReadsAnIndentedInvocation(t *testing.T) {
	document := "some prose\n\n    socrates one begin --intent \"x\"\n"
	named := skill.Named([]skill.Skill{documentSkill("s", document)}, testNamespaces)

	want := []skill.Invocation{
		{Binary: "binary-one", Command: "begin"},
		{Binary: "binary-one", Command: "begin", Flag: "--intent"},
	}
	if !slices.Equal(named, want) {
		t.Errorf("Named() = %+v, want %+v", named, want)
	}
}

func TestNamedReadsAnInlineInvocation(t *testing.T) {
	document := "run `socrates two poll --session` to read what they answered.\n"
	named := skill.Named([]skill.Skill{documentSkill("s", document)}, testNamespaces)

	want := []skill.Invocation{
		{Binary: "binary-two", Command: "poll"},
		{Binary: "binary-two", Command: "poll", Flag: "--session"},
	}
	if !slices.Equal(named, want) {
		t.Errorf("Named() = %+v, want %+v", named, want)
	}
}

func TestNamedIgnoresANamespaceItWasNotToldAbout(t *testing.T) {
	document := "    socrates three begin\n"
	named := skill.Named([]skill.Skill{documentSkill("s", document)}, testNamespaces)

	if len(named) != 0 {
		t.Errorf("Named() read %+v from a namespace nothing maps, and it should read nothing", named)
	}
}

func TestNamedDeduplicatesRepeatedInvocations(t *testing.T) {
	document := "    socrates one begin\n    socrates one begin\n"
	named := skill.Named([]skill.Skill{documentSkill("s", document)}, testNamespaces)

	if len(named) != 1 {
		t.Errorf("Named() read the same invocation %d times, want 1: %+v", len(named), named)
	}
}

func TestUnavailableReportsABinaryThisRepositoryDoesNotBuild(t *testing.T) {
	document := "    socrates one begin\n"
	skills := []skill.Skill{documentSkill("s", document)}

	absent := skill.Unavailable(skills, testNamespaces, map[string]map[string][]string{})
	if len(absent) != 1 {
		t.Fatalf("Unavailable() reported %d entries for a binary that is not built, want 1: %v", len(absent), absent)
	}
	if !strings.Contains(absent[0], "binary-one") {
		t.Errorf("the report does not name the missing binary: %q", absent[0])
	}
}

func TestUnavailableReportsACommandTheBinaryDoesNotOffer(t *testing.T) {
	document := "    socrates one begin\n"
	skills := []skill.Skill{documentSkill("s", document)}

	absent := skill.Unavailable(skills, testNamespaces, map[string]map[string][]string{"binary-one": {}})
	if len(absent) != 1 {
		t.Fatalf("Unavailable() reported %d entries for a command that is not offered, want 1: %v", len(absent), absent)
	}
	if !strings.Contains(absent[0], "begin") {
		t.Errorf("the report does not name the missing command: %q", absent[0])
	}
}

func TestUnavailableReportsAFlagTheCommandDoesNotOffer(t *testing.T) {
	document := "    socrates one begin --intent\n"
	skills := []skill.Skill{documentSkill("s", document)}

	absent := skill.Unavailable(skills, testNamespaces, map[string]map[string][]string{"binary-one": {"begin": {}}})
	if len(absent) != 1 {
		t.Fatalf("Unavailable() reported %d entries for a flag that is not offered, want 1: %v", len(absent), absent)
	}
	if !strings.Contains(absent[0], "--intent") {
		t.Errorf("the report does not name the missing flag: %q", absent[0])
	}
}

func TestUnavailableReportsNothingForACommandAndFlagBothOffered(t *testing.T) {
	document := "    socrates one begin --intent\n"
	skills := []skill.Skill{documentSkill("s", document)}

	absent := skill.Unavailable(skills, testNamespaces, map[string]map[string][]string{"binary-one": {"begin": {"--intent"}}})
	if len(absent) != 0 {
		t.Errorf("Unavailable() reported %v for a command and flag both offered", absent)
	}
}
