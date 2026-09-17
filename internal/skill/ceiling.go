package skill

import "fmt"

func CeilingViolations(skills []Skill) []string {
	var violations []string

	for _, one := range skills {
		rendered := len(one.Document())

		if rendered > one.Ceiling {
			violations = append(violations, fmt.Sprintf("%s: "+theSkillIsTooLarge, one.Name, rendered, one.Ceiling))
		}
		if one.Ceiling == 0 {
			violations = append(violations, fmt.Sprintf(theSkillHasNoCeiling, one.Name))
		}
	}

	return violations
}
