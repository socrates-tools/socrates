package skill_test

import (
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/design/skill"
)

func TestTheSkillNeverWritesNoneOfTheseAsAnOption(t *testing.T) {
	for _, one := range skill.All() {
		document := strings.ToLower(one.Rendered())

		for _, forbidden := range []string{`"label": "none of these"`, `"label":"none of these"`} {
			if strings.Contains(document, forbidden) {
				t.Errorf("%s writes %q as an option, and the page already offers that control on a single-choice decision", one.Name, forbidden)
			}
		}
	}
}

func TestTheSkillStatesNoSeverityInADocumentItWrites(t *testing.T) {
	for _, one := range skill.All() {
		document := strings.ToLower(one.Rendered())

		if strings.Contains(document, `"severity":`) || strings.Contains(document, `"severity" :`) {
			t.Errorf("%s writes a severity field into a document it composes, and a design question is earlier or later in the argument rather than more or less severe", one.Name)
		}
	}
}

func TestNoOptionTheSkillWritesCarriesACommand(t *testing.T) {
	for _, one := range skill.All() {
		document := strings.ToLower(one.Rendered())

		if strings.Contains(document, `"command":`) || strings.Contains(document, `"command" :`) {
			t.Errorf("%s writes a command into an option, and a design answer runs nothing", one.Name)
		}
	}
}
