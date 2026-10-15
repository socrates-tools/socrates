package skill_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/implement/skill"
	skillkit "github.com/socrates-tools/socrates/internal/skill"
)

func TestEachHostsRenderingNamesItsOwnDirectoryAndHostWhereItSaysToInstallElsewhere(t *testing.T) {
	for _, one := range skill.All() {
		_, installing, found := strings.Cut(one.Rendered(), "## Installing elsewhere")
		if !found {
			t.Fatalf("%s for %s has no Installing elsewhere section", one.Name, one.Host.Name)
		}

		if want := "~/" + filepath.ToSlash(one.Host.SkillsDirectory); !strings.Contains(installing, want) {
			t.Errorf("%s for %s does not name %s under Installing elsewhere:\n%s", one.Name, one.Host.Name, want, installing)
		}
		if want := "socrates implement skill --host " + one.Host.Flag + " --to <directory>"; !strings.Contains(installing, want) {
			t.Errorf("%s for %s does not say %q under Installing elsewhere:\n%s", one.Name, one.Host.Name, want, installing)
		}
		for _, other := range skillkit.Hosts() {
			if other != one.Host && strings.Contains(one.Rendered(), filepath.ToSlash(other.SkillsDirectory)) {
				t.Errorf("%s for %s names %s's skills directory", one.Name, one.Host.Name, other.Name)
			}
		}
	}
}

func TestTheCodexOverseerSpawnsEachImplementerSendsFeedbackCollectsReportsAndClosesItOnceItsTicketHasLanded(t *testing.T) {
	rendered := skill.Document(skillkit.Codex)

	for _, want := range []string{
		"`spawn_agent`",
		"`agent_type` `socrates-implementer`",
		"`send_input`",
		"`wait_agent`",
		"answer the person between",
		"`close_agent`",
		"once its ticket has landed",
	} {
		if !strings.Contains(rendered, want) {
			t.Errorf("the Codex overseer is not told %q:\n%s", want, rendered)
		}
	}
}

func TestTheCodexOverseerIsNeverToldToRunAnythingInTheBackground(t *testing.T) {
	if strings.Contains(skill.Document(skillkit.Codex), "background") {
		t.Error("the Codex overseer is told to run something in the background, and Codex has no such thing: it spawns and waits")
	}
}

func TestClaudeCodesOverseerIsToldNothingOfCodexsTools(t *testing.T) {
	rendered := skill.Document(skillkit.ClaudeCode)

	for _, unsaid := range []string{"spawn_agent", "send_input", "wait_agent", "close_agent", "writable_roots", "--add-dir"} {
		if strings.Contains(rendered, unsaid) {
			t.Errorf("the Claude Code overseer is told %q, which is Codex's", unsaid)
		}
	}
}

func TestTheCodexOverseerChecksAtTheStartOfEveryRunThatItCanWriteTheDirectoryBesideTheConfigFileAndStartsNothingWhereItCannot(t *testing.T) {
	rendered := skill.Document(skillkit.Codex)

	for _, want := range []string{
		"at the start of every run",
		"the directory beside the file `socrates implement config` prints on its first line",
		"start nothing",
		"`writable_roots` under `[sandbox_workspace_write]`",
		"`--add-dir <that directory>`",
	} {
		if !strings.Contains(rendered, want) {
			t.Errorf("the Codex overseer is not told %q about the directory it must write:\n%s", want, rendered)
		}
	}
}

func TestTheCodexOverseerTellsThePersonToRunSocratesImplementSkillWhereTheImplementerIsNotARoleItCanSpawn(t *testing.T) {
	want := "Where `socrates-implementer` is not a role you can spawn, tell the person to run `socrates implement skill` and start nothing."
	if rendered := skill.Document(skillkit.Codex); !strings.Contains(rendered, want) {
		t.Errorf("the Codex overseer does not say %q:\n%s", want, rendered)
	}
}

func TestEachHostsOverseerIsHeldToItsOwnCeiling(t *testing.T) {
	for _, host := range skillkit.Hosts() {
		if skill.OverseerCeilings.For(host) == 0 {
			t.Errorf("the overseer states no ceiling for %s", host.Name)
		}
	}
}

func TestOnlyClaudeCodesInstallingElsewhereNamesClaudeConfigDirForTheSkillAndTheDefinition(t *testing.T) {
	for _, one := range skill.All() {
		_, installing, _ := strings.Cut(one.Rendered(), "## Installing elsewhere")

		if one.Host == skillkit.ClaudeCode {
			for _, want := range []string{"$CLAUDE_CONFIG_DIR/skills", "$CLAUDE_CONFIG_DIR/agents"} {
				if !strings.Contains(installing, want) {
					t.Errorf("%s for Claude Code does not say %s under Installing elsewhere:\n%s", one.Name, want, installing)
				}
			}
		} else if strings.Contains(one.Rendered(), "CLAUDE_CONFIG_DIR") {
			t.Errorf("%s for %s names CLAUDE_CONFIG_DIR", one.Name, one.Host.Name)
		}
	}
}
