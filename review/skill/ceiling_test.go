package skill_test

import (
	"testing"

	skillkit "github.com/socrates-tools/socrates/internal/skill"
	"github.com/socrates-tools/socrates/review/skill"
)

func TestASkillOverItsCeilingFailsTheSuite(t *testing.T) {
	for _, violation := range skillkit.CeilingViolations(skill.All()) {
		t.Error(violation)
	}
}
