package skill_test

import (
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/skill"
)

const aPage = "the design questions page"

func TestClaudeCodeNeedsNoGuidanceBeyondTheSkill(t *testing.T) {
	if got := skill.ClaudeCode.Guidance(aPage, true); got != "" {
		t.Errorf("Claude Code is given guidance %q, and its sandbox and its subagents need none", got)
	}
}

func TestCodexIsToldToEscalateServePollAndEndWithAJustificationNamingThePage(t *testing.T) {
	got := skill.Codex.Guidance(aPage, false)

	for _, want := range []string{"socrates agora serve", "socrates agora poll", "socrates agora end", "sandbox_permissions", "require_escalated", "justification", aPage} {
		if !strings.Contains(got, want) {
			t.Errorf("Codex's guidance does not say %q:\n%s", want, got)
		}
	}
}

func TestCodexIsToldToInvokeNoSkill(t *testing.T) {
	if got := skill.Codex.Guidance(aPage, false); !strings.Contains(got, "Invoke no skill") {
		t.Errorf("Codex's guidance does not forbid invoking a skill:\n%s", got)
	}
}

func TestCodexIsToldToSpawnAndWaitOnlyWhereTheSkillDelegates(t *testing.T) {
	delegating := skill.Codex.Guidance(aPage, true)
	for _, want := range []string{"spawn_agent", "wait_agent"} {
		if !strings.Contains(delegating, want) {
			t.Errorf("Codex's guidance for a skill that delegates does not say %q:\n%s", want, delegating)
		}
	}

	staying := skill.Codex.Guidance(aPage, false)
	for _, unsaid := range []string{"spawn_agent", "wait_agent"} {
		if strings.Contains(staying, unsaid) {
			t.Errorf("Codex's guidance for a skill that delegates nothing says %q:\n%s", unsaid, staying)
		}
	}
}

func TestCodexIsToldToWaitBrieflyEnoughToAnswerThePerson(t *testing.T) {
	if got := skill.Codex.Guidance(aPage, true); !strings.Contains(got, "answer the person between") {
		t.Errorf("Codex's guidance does not keep its waits short enough to answer the person:\n%s", got)
	}
}
