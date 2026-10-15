package skill_test

import (
	skillkit "github.com/socrates-tools/socrates/internal/skill"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/design/skill"
)

const glossaryHeading = "## Keeping the glossary"

func TestTheSkillStatesTheGlossaryIsMaintainedDuringTheInterview(t *testing.T) {
	glossary := sectionOf(t, skill.Document(skillkit.ClaudeCode), glossaryHeading)

	for _, said := range []string{
		"is maintained during the interview, not after it",
		"A term worth writing down is written down as soon as it settles",
	} {
		if !strings.Contains(glossary, said) {
			t.Errorf("the glossary section does not say %q, so a session has no rule for when CONTEXT.md is written:\n%s", said, glossary)
		}
	}
}

func TestTheSkillPutsAConflictingOrRenamingTermOnACardOfItsOwn(t *testing.T) {
	glossary := sectionOf(t, skill.Document(skillkit.ClaudeCode), glossaryHeading)

	for _, said := range []string{
		"before writing an entry",
		"conflicts with an entry already there, or that renames something the glossary already calls something else",
		"is never applied silently: put it on a card of its own and let the person decide",
	} {
		if !strings.Contains(glossary, said) {
			t.Errorf("the glossary section does not say %q, so a session could rename a term without the person agreeing to it:\n%s", said, glossary)
		}
	}
}

func TestTheGlossarySectionComesAfterFacts(t *testing.T) {
	document := skill.Document(skillkit.ClaudeCode)

	facts := strings.Index(document, factsHeading)
	glossary := strings.Index(document, glossaryHeading)
	if facts < 0 || glossary < 0 {
		t.Fatalf("the skill is missing one of %q or %q", factsHeading, glossaryHeading)
	}
	if facts > glossary {
		t.Error("the glossary section is stated before the facts section")
	}
}
