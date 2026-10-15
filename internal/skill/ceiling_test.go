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

func TestACeilingIsStatedForEachHostAndAHostItNamesNoneForHasNoCeiling(t *testing.T) {
	ceilings := skill.Ceilings{skill.ClaudeCode: 100, skill.Codex: 150}

	if got := ceilings.For(skill.ClaudeCode); got != 100 {
		t.Errorf("Claude Code's ceiling is %d, want 100", got)
	}
	if got := ceilings.For(skill.Codex); got != 150 {
		t.Errorf("Codex's ceiling is %d, want 150", got)
	}

	unstated := skill.Ceilings{skill.ClaudeCode: 100}
	if got := unstated.For(skill.Codex); got != 0 {
		t.Errorf("a host the ceilings do not name has ceiling %d, want 0 so that the suite reports it", got)
	}
}
