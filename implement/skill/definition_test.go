package skill_test

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/implement/config"
	"github.com/socrates-tools/socrates/implement/skill"
	"github.com/socrates-tools/socrates/internal/convention"
	skillkit "github.com/socrates-tools/socrates/internal/skill"
	"github.com/socrates-tools/socrates/internal/state"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

func frontMatterOf(t *testing.T, document string) map[string]string {
	t.Helper()

	lines := strings.Split(document, "\n")
	if lines[0] != "---" {
		t.Fatalf("the definition does not open with front matter: %q", lines[0])
	}
	fields := map[string]string{}
	for _, line := range lines[1:] {
		if line == "---" {
			return fields
		}
		key, value, _ := strings.Cut(line, ": ")
		fields[key] = value
	}
	t.Fatal("the definition's front matter is never closed")
	return nil
}

func stated(model, effort string) config.Config {
	return config.Config{Roles: config.Roles{Implementer: config.Role{Model: model, Effort: effort}}}
}

func TestTheDefinitionIsHeldToACeilingOfItsOwn(t *testing.T) {
	for _, violation := range skillkit.CeilingViolations(skill.WithAgentFiles()) {
		t.Error(violation)
	}
}

func TestTheDefinitionRunsTheImplementerInTheBackgroundUnderItsOwnName(t *testing.T) {
	fields := frontMatterOf(t, skill.DefaultDefinition())

	if fields["name"] != skill.AgentName {
		t.Errorf("the definition names itself %q, want %q", fields["name"], skill.AgentName)
	}
	if fields["background"] != "true" {
		t.Errorf("the definition states background %q, and an implementer runs in the background", fields["background"])
	}
	if fields["description"] == "" {
		t.Error("the definition states no description, and Claude Code requires one")
	}
	if _, carries := fields["disable-model-invocation"]; carries {
		t.Error("the definition carries disable-model-invocation, which belongs to a skill")
	}
}

func TestTheDefinitionAtTheDefaultsCarriesTheDefaultModelAndEffort(t *testing.T) {
	fields := frontMatterOf(t, skill.DefinitionFor(config.Config{}))

	if fields["model"] != config.DefaultModel || fields["effort"] != config.DefaultEffort {
		t.Errorf("the definition states model %q and effort %q, want %q and %q", fields["model"], fields["effort"], config.DefaultModel, config.DefaultEffort)
	}
}

func TestTheDefinitionAfterAPersonSetEachCarriesWhatTheyChose(t *testing.T) {
	fields := frontMatterOf(t, skill.DefinitionFor(stated("opus", "max")))

	if fields["model"] != "opus" || fields["effort"] != "max" {
		t.Errorf("the definition states model %q and effort %q, want opus and max", fields["model"], fields["effort"])
	}
}

func TestTheDefinitionStatesTheStandingInstructions(t *testing.T) {
	document := skill.DefaultDefinition()

	for _, promise := range []string{
		"Work in the worktree you were briefed with and nowhere else",
		"Follow the repository's own instructions, including how it tests",
		"Commit on your own branch",
		"Run the repository's own verification before you report",
		"Merge the change's branch into your own branch before you report done",
		"Report what you did and what you did not",
		"When feedback arrives, patch your own work",
		"Never invoke a socrates skill",
		"Never use git stash: the stash is shared by every worktree of the repository",
		"Where the repository states no way of testing, write each test before the code and watch it fail",
		"test behaviour through the public interface",
		"take expected values from the ticket or a worked example rather than computing them the way the code does",
		"mock only at the system's edges",
	} {
		if !strings.Contains(document, promise) {
			t.Errorf("the definition does not say %q", promise)
		}
	}
}

func aSocratesCommand() *regexp.Regexp {
	var namespaces []string
	for _, tool := range convention.Tools {
		namespaces = append(namespaces, regexp.QuoteMeta(tool.Namespace))
	}
	return regexp.MustCompile(`(?i)` + convention.Suite + `\s+(` + strings.Join(namespaces, "|") + `)\b`)
}

func TestTheDefinitionNamesNoSocratesCommand(t *testing.T) {
	if found := aSocratesCommand().FindString(skill.DefaultDefinition()); found != "" {
		t.Errorf("the definition names %q, and an implementer runs no socrates command", found)
	}
}

func TestTheCommittedCopyMatchesWhatIsRenderedAtTheDefaults(t *testing.T) {
	path := filepath.Join(skillkit.ClaudeCode.AgentsDirectoryUnder(source.Root(t)), skill.AgentFileName())

	committed, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("%s is not readable, and it is the copy this repository ships: %v. Run make skill to write it", path, err)
	}
	if string(committed) != skill.DefaultDefinition() {
		t.Errorf("%s is not what the definition's source generates: run make skill to rewrite it", path)
	}
}

func TestInstallingWritesUnderTheNamedDirectoryAndRecordsIt(t *testing.T) {
	t.Setenv(state.HomeVariable, t.TempDir())
	directory := filepath.Join(t.TempDir(), "agents")

	path, err := skill.Install(skillkit.ClaudeCode, directory, stated("opus", "low"))
	if err != nil {
		t.Fatalf("installing: %v", err)
	}

	if want := filepath.Join(directory, skill.AgentFileName()); path != want {
		t.Errorf("installing wrote %s, want %s", path, want)
	}
	written, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading what was written: %v", err)
	}
	if fields := frontMatterOf(t, string(written)); fields["model"] != "opus" || fields["effort"] != "low" {
		t.Errorf("what was written states model %q and effort %q, want opus and low", fields["model"], fields["effort"])
	}

	recorded, installed, err := skill.Recorded(skillkit.ClaudeCode)
	if err != nil || !installed || recorded != directory {
		t.Errorf("the record holds %q (installed %v, error %v), want %q", recorded, installed, err, directory)
	}
}

func TestASecondInstallSomewhereElseReplacesTheRecord(t *testing.T) {
	t.Setenv(state.HomeVariable, t.TempDir())
	first := filepath.Join(t.TempDir(), "first")
	second := filepath.Join(t.TempDir(), "second")

	for _, directory := range []string{first, second} {
		if _, err := skill.Install(skillkit.ClaudeCode, directory, config.Config{}); err != nil {
			t.Fatalf("installing into %s: %v", directory, err)
		}
	}

	if recorded, _, _ := skill.Recorded(skillkit.ClaudeCode); recorded != second {
		t.Errorf("the record holds %q after a second install into %q", recorded, second)
	}
}

func TestNothingIsRecordedWhereNothingWasInstalled(t *testing.T) {
	t.Setenv(state.HomeVariable, t.TempDir())

	if recorded, installed, err := skill.Recorded(skillkit.ClaudeCode); err != nil || installed {
		t.Errorf("the record holds %q (installed %v, error %v) before anything was installed", recorded, installed, err)
	}
	if _, rewrote, err := skill.Reinstall(skillkit.ClaudeCode, config.Config{}); err != nil || rewrote {
		t.Errorf("reinstalling wrote something (%v, error %v) where nothing was installed", rewrote, err)
	}
}

func TestReinstallingWritesTheDefinitionAgainWhereTheRecordSays(t *testing.T) {
	t.Setenv(state.HomeVariable, t.TempDir())
	directory := filepath.Join(t.TempDir(), "agents")
	if _, err := skill.Install(skillkit.ClaudeCode, directory, config.Config{}); err != nil {
		t.Fatalf("installing: %v", err)
	}

	path, rewrote, err := skill.Reinstall(skillkit.ClaudeCode, stated("haiku", "medium"))
	if err != nil || !rewrote {
		t.Fatalf("reinstalling wrote %v with error %v", rewrote, err)
	}
	written, _ := os.ReadFile(path)
	if fields := frontMatterOf(t, string(written)); fields["model"] != "haiku" || fields["effort"] != "medium" {
		t.Errorf("reinstalling left model %q and effort %q, want haiku and medium", fields["model"], fields["effort"])
	}
}

func TestTheDefaultDirectoryIsUnderTheHome(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("CLAUDE_CONFIG_DIR", "")

	directory, err := skill.AgentsDirectory(skillkit.ClaudeCode)
	if err != nil {
		t.Fatalf("resolving the default directory: %v", err)
	}
	if want := filepath.Join(home, ".claude", "agents"); directory != want {
		t.Errorf("the default directory is %s, want %s", directory, want)
	}
}
