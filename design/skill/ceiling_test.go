package skill_test

import (
	"testing"

	"github.com/socrates-tools/socrates/design/skill"
	skillkit "github.com/socrates-tools/socrates/internal/skill"
)

func TestASkillOverItsCeilingFailsTheSuite(t *testing.T) {
	for _, violation := range skillkit.CeilingViolations(skill.All()) {
		t.Error(violation)
	}
}
