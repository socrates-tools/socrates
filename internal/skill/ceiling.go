package skill

import "fmt"

func CeilingViolations(skills []Skill) []string {
	var violations []string

	for _, one := range skills {
		rendered := len(one.Rendered())

		if rendered > one.Ceiling {
			violations = append(violations, fmt.Sprintf("%s: "+theSkillIsTooLarge, onHost(one), rendered, one.Ceiling))
		}
		if one.Ceiling == 0 {
			violations = append(violations, fmt.Sprintf(theSkillHasNoCeiling, onHost(one)))
		}
	}

	return violations
}

func onHost(one Skill) string {
	return fmt.Sprintf(skillOnAHost, one.Name, one.Host.Name)
}
