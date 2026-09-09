package skill

import "testing"

func TestASkillOverItsCeilingFailsTheSuite(t *testing.T) {
	for _, one := range All() {
		rendered := len(one.Document())

		if rendered > one.Ceiling {
			t.Errorf("%s: "+theSkillIsTooLarge, one.Name, rendered, one.Ceiling)
		}
		if one.Ceiling == 0 {
			t.Errorf("%s states no ceiling, so nothing bounds what every session that reads it reads", one.Name)
		}
	}
}
