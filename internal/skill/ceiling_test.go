package skill

import "testing"

func TestASkillOverTheCeilingFailsTheSuite(t *testing.T) {
	rendered := len(Document())

	if rendered > Ceiling {
		t.Errorf(theSkillIsTooLarge, rendered, Ceiling)
	}
}
