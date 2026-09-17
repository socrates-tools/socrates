package skill_test

import (
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/design/skill"
)

func TestTheSkillNeverWritesNoneOfTheseAsAnOption(t *testing.T) {
	document := strings.ToLower(skill.Document())

	for _, forbidden := range []string{`"label": "none of these"`, `"label":"none of these"`} {
		if strings.Contains(document, forbidden) {
			t.Errorf("the skill writes %q as an option, and the page already offers that control on a single-choice decision", forbidden)
		}
	}
}

func TestTheSkillStatesNoSeverityInADocumentItWrites(t *testing.T) {
	document := strings.ToLower(skill.Document())

	if strings.Contains(document, `"severity":`) || strings.Contains(document, `"severity" :`) {
		t.Error("the skill writes a severity field into a document it composes, and a design question is earlier or later in the argument rather than more or less severe")
	}
}

func TestNoOptionTheSkillWritesCarriesACommand(t *testing.T) {
	document := strings.ToLower(skill.Document())

	if strings.Contains(document, `"command":`) || strings.Contains(document, `"command" :`) {
		t.Error("the skill writes a command into an option, and a design answer runs nothing")
	}
}
