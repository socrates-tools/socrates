package skill_test

import (
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/skill"
)

func renderedSkill(name string, rendered, ceiling int) skill.Skill {
	document := strings.Repeat("a", rendered)
	return skill.Skill{Name: name, Host: skill.ClaudeCode, Document: func(skill.Host) string { return document }, Ceiling: ceiling}
}

func TestASkillOverItsCeilingIsReported(t *testing.T) {
	violations := skill.CeilingViolations([]skill.Skill{renderedSkill("big", 10, 5)})
	if len(violations) != 1 {
		t.Fatalf("a skill rendering past its ceiling is reported %d times, want 1: %v", len(violations), violations)
	}
}

func TestASkillWithNoCeilingIsReported(t *testing.T) {
	violations := skill.CeilingViolations([]skill.Skill{renderedSkill("unbounded", 0, 0)})
	if len(violations) != 1 {
		t.Fatalf("a skill stating no ceiling is reported %d times, want 1: %v", len(violations), violations)
	}
}

func TestASkillWithinItsCeilingIsNotReported(t *testing.T) {
	violations := skill.CeilingViolations([]skill.Skill{renderedSkill("fine", 5, 10)})
	if len(violations) != 0 {
		t.Errorf("a skill within its ceiling is reported: %v", violations)
	}
}

func TestAViolationNamesTheHostTheSkillIsRenderedFor(t *testing.T) {
	violations := skill.CeilingViolations([]skill.Skill{renderedSkill("big", 10, 5)})
	if len(violations) != 1 || !strings.Contains(violations[0], skill.ClaudeCode.Name) {
		t.Errorf("a violation does not name the host it was rendered for: %v", violations)
	}
}
