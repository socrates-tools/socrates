package skill_test

import (
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/design/skill"
	skillkit "github.com/socrates-tools/socrates/internal/skill"
)

func refactorRenderings(t *testing.T) []skillkit.Skill {
	t.Helper()
	var renderings []skillkit.Skill
	for _, one := range skill.All() {
		if one.Name == skill.RefactorName {
			renderings = append(renderings, one)
		}
	}
	if len(renderings) != len(skillkit.Hosts()) {
		t.Fatalf("design renders %s for %d hosts, want one per host", skill.RefactorName, len(renderings))
	}
	return renderings
}

func TestTheRefactorSkillNamesEverySkillItComposesAndWhereEachComesFrom(t *testing.T) {
	for _, one := range refactorRenderings(t) {
		for _, want := range []string{"improve-codebase-architecture", "codebase-design", "grilling", "domain-modeling", "lavish", "mattpocock/skills", "kunchenguid/lavish-axi"} {
			if !strings.Contains(one.Rendered(), want) {
				t.Errorf("%s for %s does not name %s", one.Name, one.Host.Name, want)
			}
		}
	}
}

func TestTheRefactorSkillInstallsAMissingSkillForItsOwnHostOnly(t *testing.T) {
	for _, one := range refactorRenderings(t) {
		installing := sectionOf(t, one.Rendered(), "## Installing what is missing")
		for _, source := range []string{"npx -y skills add mattpocock/skills", "npx -y skills add kunchenguid/lavish-axi"} {
			if !strings.Contains(installing, source) {
				t.Errorf("%s for %s does not say to run %q when a skill is missing:\n%s", one.Name, one.Host.Name, source, installing)
			}
		}
		for _, host := range skillkit.Hosts() {
			agent := host.SkillsCLIAgent
			named := strings.Count(installing, "-a "+agent+" ")
			if host == one.Host && named != 2 {
				t.Errorf("%s for %s names -a %s %d times, want once per install:\n%s", one.Name, one.Host.Name, agent, named, installing)
			}
			if host != one.Host && named != 0 {
				t.Errorf("%s for %s installs for %s", one.Name, one.Host.Name, host.Name)
			}
		}
	}
}

func TestTheRefactorSkillPutsTheReportAndThePickOnALavishPageRatherThanTheTerminal(t *testing.T) {
	for _, one := range refactorRenderings(t) {
		running := sectionOf(t, one.Rendered(), "## Running it")
		for _, want := range []string{"npx -y lavish-axi", "HTML-REPORT.md", "Which of these would you like to explore?", "never in the terminal"} {
			if !strings.Contains(running, want) {
				t.Errorf("%s for %s does not say %q under Running it:\n%s", one.Name, one.Host.Name, want, running)
			}
		}
	}
}

func TestTheRefactorSkillReachesForNoOtherSocratesSkill(t *testing.T) {
	for _, one := range refactorRenderings(t) {
		for _, other := range skill.All() {
			if other.Name == skill.RefactorName {
				continue
			}
			if strings.Contains(one.Rendered(), other.Name) {
				t.Errorf("%s for %s names %s, and every socrates skill is a person's to start", one.Name, one.Host.Name, other.Name)
			}
		}
	}
}

func TestOnlyCodexsRefactorRenderingEscalatesWhatLeavesTheSandbox(t *testing.T) {
	for _, one := range refactorRenderings(t) {
		escalates := strings.Contains(one.Rendered(), "require_escalated")
		if one.Host == skillkit.Codex && !escalates {
			t.Errorf("%s for Codex does not escalate the installs and the lavish page out of the sandbox", one.Name)
		}
		if one.Host != skillkit.Codex && escalates {
			t.Errorf("%s for %s names Codex's escalation", one.Name, one.Host.Name)
		}
	}
}
