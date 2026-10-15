package skill_test

import (
	skillkit "github.com/socrates-tools/socrates/internal/skill"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/design/skill"
)

const specHeading = "## Writing the spec"

func TestTheSkillNamesTheSpecsPathOnTheRepositorysTemplate(t *testing.T) {
	spec := sectionOf(t, skill.Document(skillkit.ClaudeCode), specHeading)

	if !strings.Contains(spec, "docs/changes/<slug>/spec.md") {
		t.Errorf("the spec section does not name docs/changes/<slug>/spec.md as where the spec goes:\n%s", spec)
	}

	for _, section := range []string{
		"Problem Statement",
		"Solution",
		"User Stories",
		"Implementation Decisions",
		"Testing Decisions",
		"Out of Scope",
		"Further Notes",
	} {
		if !strings.Contains(spec, section) {
			t.Errorf("the spec section does not name %q, and the skill should name the template's sections rather than inventing a shape:\n%s", section, spec)
		}
	}
}

func TestTheSpecSectionsAreNamedInTemplateOrder(t *testing.T) {
	spec := sectionOf(t, skill.Document(skillkit.ClaudeCode), specHeading)

	order := []string{
		"Problem Statement",
		"Solution",
		"User Stories",
		"Implementation Decisions",
		"Testing Decisions",
		"Out of Scope",
		"Further Notes",
	}

	last := -1
	for _, section := range order {
		at := strings.Index(spec, section)
		if at < 0 {
			t.Fatalf("the spec section does not name %q", section)
		}
		if at < last {
			t.Errorf("%q is named before %q, and the template lists them in a fixed order", section, order[0])
		}
		last = at
	}
}

func TestTheSpecSectionSaysFurtherNotesCarriesWhatADesignEndedEarlyLeftOpen(t *testing.T) {
	spec := sectionOf(t, skill.Document(skillkit.ClaudeCode), specHeading)

	if !strings.Contains(spec, "where a design ends early, what it left open") {
		t.Errorf("the spec section does not connect Further Notes to a design that ended early:\n%s", spec)
	}
}

func TestTheSpecSectionComesAfterTheGlossarySection(t *testing.T) {
	document := skill.Document(skillkit.ClaudeCode)

	glossary := strings.Index(document, glossaryHeading)
	spec := strings.Index(document, specHeading)
	if glossary < 0 || spec < 0 {
		t.Fatalf("the skill is missing one of %q or %q", glossaryHeading, specHeading)
	}
	if glossary > spec {
		t.Error("the spec section is stated before the glossary section")
	}
}
