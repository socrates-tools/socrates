package skill_test

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/socrates-tools/socrates/implement/config"
	"github.com/socrates-tools/socrates/implement/skill"
	skillkit "github.com/socrates-tools/socrates/internal/skill"
	"github.com/socrates-tools/socrates/internal/state"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

var keysCodexAccepts = []string{"name", "description", "developer_instructions", "model_reasoning_effort", "model"}

func roleRead(t *testing.T, document string) map[string]any {
	t.Helper()

	var read map[string]any
	if _, err := toml.Decode(document, &read); err != nil {
		t.Fatalf("the role is not TOML: %v\n%s", err, document)
	}
	return read
}

func TestTheCodexRoleHoldsOnlyTheKeysCodexAccepts(t *testing.T) {
	for name, document := range map[string]string{
		"at the defaults":  skill.DefaultCodexRole(),
		"with a model set": skill.CodexRole("gpt-5.5", "low"),
	} {
		t.Run(name, func(t *testing.T) {
			for key := range roleRead(t, document) {
				if !slices.Contains(keysCodexAccepts, key) {
					t.Errorf("the role holds %q, and Codex refuses a role file holding a key it does not know", key)
				}
			}
		})
	}
}

func TestTheCodexRoleAtTheDefaultsNamesItselfAndCarriesTheDefaultEffortAndNoModel(t *testing.T) {
	read := roleRead(t, skill.DefaultCodexRole())

	if read["name"] != skill.AgentName {
		t.Errorf("the role names itself %v, want %q", read["name"], skill.AgentName)
	}
	if read["description"] == "" || read["description"] == nil {
		t.Error("the role states no description, and Codex requires one")
	}
	if read["model_reasoning_effort"] != "high" {
		t.Errorf("the role states effort %v, want high", read["model_reasoning_effort"])
	}
	if model, stated := read["model"]; stated {
		t.Errorf("the role states model %v where none was set, and the implementer should run on the session's model", model)
	}
}

func TestTheCodexRoleCarriesAModelWhereOneIsSet(t *testing.T) {
	read := roleRead(t, skill.CodexRole("gpt-5.5", "xhigh"))

	if read["model"] != "gpt-5.5" || read["model_reasoning_effort"] != "xhigh" {
		t.Errorf("the role states model %v and effort %v, want gpt-5.5 and xhigh", read["model"], read["model_reasoning_effort"])
	}
}

func TestTheCodexRolesInstructionsAreTheOnesTheClaudeCodeDefinitionCarries(t *testing.T) {
	instructions, _ := roleRead(t, skill.DefaultCodexRole())["developer_instructions"].(string)
	_, claudes, found := strings.Cut(skill.DefaultDefinition(), "---\n\n")
	if !found {
		t.Fatal("the Claude Code definition has no body after its front matter")
	}

	if strings.TrimSpace(instructions) != strings.TrimSpace(claudes) {
		t.Errorf("the role's instructions are not the definition's:\n%s\n---\n%s", instructions, claudes)
	}
}

func TestTheCodexRoleNamesNoSocratesCommand(t *testing.T) {
	if found := aSocratesCommand().FindString(skill.DefaultCodexRole()); found != "" {
		t.Errorf("the role names %q, and an implementer runs no socrates command", found)
	}
}

func TestTheCodexRoleIsHeldToACeilingOfItsOwn(t *testing.T) {
	held := false
	for _, one := range skill.WithAgentFiles() {
		if one.Name == skill.AgentName && one.Host == skillkit.Codex {
			held = true
		}
	}
	if !held {
		t.Error("no ceiling holds the Codex role")
	}
}

func TestTheCommittedCodexRoleMatchesWhatIsRenderedAtTheDefaults(t *testing.T) {
	path := filepath.Join(skillkit.Codex.AgentsDirectoryUnder(source.Root(t)), skill.CodexRoleFileName())

	committed, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("%s is not readable, and it is the copy this repository ships: %v. Run make skill to write it", path, err)
	}
	if string(committed) != skill.DefaultCodexRole() {
		t.Errorf("%s is not what the role's source generates: run make skill to rewrite it", path)
	}
}

func TestTheCodexRoleFileIsNamedForTheImplementerWithATomlSuffix(t *testing.T) {
	if got, want := skill.CodexRoleFileName(), "socrates-implementer.toml"; got != want {
		t.Errorf("the role file is %s, want %s", got, want)
	}
}

func TestInstallingForCodexWritesTheRoleUnderTheNamedDirectoryCreatingItAndRecordsItBesideClaudeCodes(t *testing.T) {
	t.Setenv(state.HomeVariable, t.TempDir())
	claudes := filepath.Join(t.TempDir(), "claude")
	codexes := filepath.Join(t.TempDir(), "not", "yet", "there")

	if _, err := skill.Install(skillkit.ClaudeCode, claudes, config.Config{}); err != nil {
		t.Fatalf("installing for Claude Code: %v", err)
	}
	path, err := skill.Install(skillkit.Codex, codexes, config.Config{})
	if err != nil {
		t.Fatalf("installing for Codex: %v", err)
	}

	if want := filepath.Join(codexes, skill.CodexRoleFileName()); path != want {
		t.Errorf("installing wrote %s, want %s", path, want)
	}
	written, err := os.ReadFile(path)
	if err != nil || string(written) != skill.DefaultCodexRole() {
		t.Errorf("what was written is not the default role (error %v):\n%s", err, written)
	}
	for host, want := range map[skillkit.Host]string{skillkit.ClaudeCode: claudes, skillkit.Codex: codexes} {
		if recorded, installed, err := skill.Recorded(host); err != nil || !installed || recorded != want {
			t.Errorf("the record holds %q for %s (installed %v, error %v), want %q", recorded, host.Name, installed, err, want)
		}
	}
}

func TestAnInstallForOneHostLeavesTheOtherHostsRecordAlone(t *testing.T) {
	t.Setenv(state.HomeVariable, t.TempDir())
	first, second := filepath.Join(t.TempDir(), "first"), filepath.Join(t.TempDir(), "second")

	for _, directory := range []string{first, second} {
		if _, err := skill.Install(skillkit.Codex, directory, config.Config{}); err != nil {
			t.Fatalf("installing into %s: %v", directory, err)
		}
	}
	if _, installed, _ := skill.Recorded(skillkit.ClaudeCode); installed {
		t.Error("the record holds a Claude Code directory though only Codex was installed")
	}
	if recorded, _, _ := skill.Recorded(skillkit.Codex); recorded != second {
		t.Errorf("the record holds %q for Codex, want the last install %q", recorded, second)
	}
}

func TestARecordInTodaysShapeIsReadAsClaudeCodes(t *testing.T) {
	root := t.TempDir()
	t.Setenv(state.HomeVariable, root)
	path, err := skill.RecordPath()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("agents_directory: /somewhere/.claude/agents\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if recorded, installed, err := skill.Recorded(skillkit.ClaudeCode); err != nil || !installed || recorded != "/somewhere/.claude/agents" {
		t.Errorf("Claude Code's record holds %q (installed %v, error %v)", recorded, installed, err)
	}
	if _, installed, err := skill.Recorded(skillkit.Codex); err != nil || installed {
		t.Errorf("Codex's record is held (installed %v, error %v) where only today's shape exists", installed, err)
	}

	codexes := filepath.Join(t.TempDir(), "codex")
	if _, err := skill.Install(skillkit.Codex, codexes, config.Config{}); err != nil {
		t.Fatalf("installing for Codex: %v", err)
	}
	if recorded, installed, _ := skill.Recorded(skillkit.ClaudeCode); !installed || recorded != "/somewhere/.claude/agents" {
		t.Errorf("recording Codex's directory lost Claude Code's, which now holds %q", recorded)
	}
}

func TestReinstallingForCodexWritesTheRoleWhereTheRecordSaysAndTouchesNoDefinition(t *testing.T) {
	t.Setenv(state.HomeVariable, t.TempDir())
	claudes, codexes := filepath.Join(t.TempDir(), "claude"), filepath.Join(t.TempDir(), "codex")
	for host, directory := range map[skillkit.Host]string{skillkit.ClaudeCode: claudes, skillkit.Codex: codexes} {
		if _, err := skill.Install(host, directory, config.Config{}); err != nil {
			t.Fatalf("installing for %s: %v", host.Name, err)
		}
	}
	if err := os.Remove(filepath.Join(codexes, skill.CodexRoleFileName())); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(claudes, skill.AgentFileName())); err != nil {
		t.Fatal(err)
	}

	path, rewrote, err := skill.Reinstall(skillkit.Codex, config.Config{})
	if err != nil || !rewrote || path != filepath.Join(codexes, skill.CodexRoleFileName()) {
		t.Fatalf("reinstalling wrote %s (%v, error %v)", path, rewrote, err)
	}
	if _, err := os.Stat(filepath.Join(claudes, skill.AgentFileName())); !os.IsNotExist(err) {
		t.Error("reinstalling for Codex rewrote Claude Code's definition")
	}
}

func TestCodexsAgentsDirectoryFollowsCodexHomeAndFallsBackToTheHomeDirectory(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	t.Setenv("CODEX_HOME", "")
	if got, err := skill.AgentsDirectory(skillkit.Codex); err != nil || got != filepath.Join(home, ".codex", "agents") {
		t.Errorf("with CODEX_HOME unset the directory is %s (error %v), want it under ~/.codex", got, err)
	}

	stated := t.TempDir()
	t.Setenv("CODEX_HOME", stated)
	if got, err := skill.AgentsDirectory(skillkit.Codex); err != nil || got != filepath.Join(stated, "agents") {
		t.Errorf("with CODEX_HOME set the directory is %s (error %v), want %s", got, err, filepath.Join(stated, "agents"))
	}
}

func statedForCodex(model, effort string) config.Config {
	return config.Config{Roles: config.Roles{Implementer: config.Role{Codex: config.CodexRole{Model: model, Effort: effort}}}}
}

func TestTheCodexRoleFollowsTheCodexSettingsAndIgnoresClaudeCodes(t *testing.T) {
	cases := map[string]struct {
		stated        config.Config
		model, effort string
	}{
		"nothing set":                     {config.Config{}, "", "high"},
		"only Claude Code's set":          {stated("opus", "max"), "", "high"},
		"a Codex model and effort set":    {statedForCodex("gpt-5.5", "xhigh"), "gpt-5.5", "xhigh"},
		"a Codex effort set and no model": {statedForCodex("", "ultra"), "", "ultra"},
		"both hosts set apart":            {config.Config{Roles: config.Roles{Implementer: config.Role{Model: "haiku", Effort: "low", Codex: config.CodexRole{Model: "gpt-5.5", Effort: "none"}}}}, "gpt-5.5", "none"},
	}
	for name, one := range cases {
		t.Run(name, func(t *testing.T) {
			read := roleRead(t, skill.CodexRoleFor(one.stated))

			model, stated := read["model"]
			if one.model == "" && stated {
				t.Errorf("the role carries model %v where none is set", model)
			}
			if one.model != "" && model != one.model {
				t.Errorf("the role carries model %v, want %s", model, one.model)
			}
			if read["model_reasoning_effort"] != one.effort {
				t.Errorf("the role carries effort %v, want %s", read["model_reasoning_effort"], one.effort)
			}
		})
	}
}

func TestInstallingForCodexWritesTheCodexSettingsInForce(t *testing.T) {
	t.Setenv(state.HomeVariable, t.TempDir())
	directory := filepath.Join(t.TempDir(), "agents")

	path, err := skill.Install(skillkit.Codex, directory, statedForCodex("gpt-5.5", "low"))
	if err != nil {
		t.Fatalf("installing: %v", err)
	}

	written, _ := os.ReadFile(path)
	if read := roleRead(t, string(written)); read["model"] != "gpt-5.5" || read["model_reasoning_effort"] != "low" {
		t.Errorf("what was written carries %v and %v, want gpt-5.5 and low", read["model"], read["model_reasoning_effort"])
	}
}

func TestReinstallingForCodexWritesTheNewSettingsAndLeavesClaudeCodesDefinitionAsItWas(t *testing.T) {
	t.Setenv(state.HomeVariable, t.TempDir())
	claudes, codexes := filepath.Join(t.TempDir(), "claude"), filepath.Join(t.TempDir(), "codex")
	for host, directory := range map[skillkit.Host]string{skillkit.ClaudeCode: claudes, skillkit.Codex: codexes} {
		if _, err := skill.Install(host, directory, config.Config{}); err != nil {
			t.Fatalf("installing for %s: %v", host.Name, err)
		}
	}

	path, rewrote, err := skill.Reinstall(skillkit.Codex, statedForCodex("gpt-5.5", "max"))
	if err != nil || !rewrote {
		t.Fatalf("reinstalling wrote %v with error %v", rewrote, err)
	}
	written, _ := os.ReadFile(path)
	if read := roleRead(t, string(written)); read["model"] != "gpt-5.5" || read["model_reasoning_effort"] != "max" {
		t.Errorf("the role carries %v and %v, want gpt-5.5 and max", read["model"], read["model_reasoning_effort"])
	}
	untouched, _ := os.ReadFile(filepath.Join(claudes, skill.AgentFileName()))
	if string(untouched) != skill.DefaultDefinition() {
		t.Errorf("reinstalling for Codex changed Claude Code's definition:\n%s", untouched)
	}
}

func TestClaudeCodesAgentsDirectoryFollowsClaudeConfigDirAndFallsBackToTheHomeDirectory(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	t.Setenv("CLAUDE_CONFIG_DIR", "")
	if got, err := skill.AgentsDirectory(skillkit.ClaudeCode); err != nil || got != filepath.Join(home, ".claude", "agents") {
		t.Errorf("with CLAUDE_CONFIG_DIR unset the directory is %s (error %v), want it under ~/.claude", got, err)
	}

	stated := t.TempDir()
	t.Setenv("CLAUDE_CONFIG_DIR", stated)
	if got, err := skill.AgentsDirectory(skillkit.ClaudeCode); err != nil || got != filepath.Join(stated, "agents") {
		t.Errorf("with CLAUDE_CONFIG_DIR set the directory is %s (error %v), want %s", got, err, filepath.Join(stated, "agents"))
	}
}
