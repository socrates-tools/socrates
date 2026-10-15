package skill_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/skill"
)

const fixtureInvocation = "socrates design skill"

func documentSkill(name, document string) skill.Skill {
	return skill.Skill{Name: name, Host: skill.ClaudeCode, Document: func(skill.Host) string { return document }, Ceiling: len(document)}
}

func TestInstallWritesEachSkillUnderItsOwnName(t *testing.T) {
	directory := t.TempDir()
	skills := []skill.Skill{documentSkill("one", "# One\n"), documentSkill("two", "# Two\n")}

	written, err := skill.Install(directory, skills, fixtureInvocation)
	if err != nil {
		t.Fatalf("installing: %v", err)
	}
	if len(written) != len(skills) {
		t.Fatalf("installing wrote %d skills, and %d were named", len(written), len(skills))
	}

	for position, one := range skills {
		want := filepath.Join(directory, one.Name, skill.FileName)
		if written[position].Path != want {
			t.Errorf("skill %s was written to %s, want %s", one.Name, written[position].Path, want)
		}
		content, err := os.ReadFile(want)
		if err != nil {
			t.Fatalf("%s is not readable: %v", want, err)
		}
		if string(content) != one.Rendered() {
			t.Errorf("%s holds %q, want %q", want, content, one.Rendered())
		}
	}
}

func TestInstallRefusesADirectoryThatCannotBeWritten(t *testing.T) {
	sealed := filepath.Join(t.TempDir(), "sealed")
	if err := os.Mkdir(sealed, 0o500); err != nil {
		t.Fatalf("creating a directory nothing can write in: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(sealed, 0o700) })

	_, err := skill.Install(sealed, []skill.Skill{documentSkill("one", "# One\n")}, fixtureInvocation)
	if err == nil {
		t.Fatal("installing into a sealed directory did not fail")
	}
}

func TestInstallsRefusalNamesTheInvocationItWasGiven(t *testing.T) {
	sealed := filepath.Join(t.TempDir(), "sealed")
	if err := os.Mkdir(sealed, 0o500); err != nil {
		t.Fatalf("creating a directory nothing can write in: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(sealed, 0o700) })

	_, err := skill.Install(sealed, []skill.Skill{documentSkill("one", "# One\n")}, fixtureInvocation)
	if err == nil {
		t.Fatal("installing into a sealed directory did not fail")
	}
	if !strings.Contains(err.Error(), fixtureInvocation) {
		t.Errorf("the refusal does not name the invocation it was given: %v", err)
	}
}

func TestUserDirectoryUnderJoinsTheHomeAndTheHostsSkillsDirectory(t *testing.T) {
	home := filepath.Join("some", "home")
	want := filepath.Join(home, ".claude", "skills")
	if got := skill.UserDirectoryUnder(skill.ClaudeCode, home); got != want {
		t.Errorf("UserDirectoryUnder(%q) = %q, want %q", home, got, want)
	}
}

func TestUserDirectorysRefusalNamesTheInvocationItWasGiven(t *testing.T) {
	homeVariable := "HOME"
	if runtime.GOOS == "windows" {
		homeVariable = "USERPROFILE"
	}
	if runtime.GOOS == "plan9" {
		homeVariable = "home"
	}
	t.Setenv(homeVariable, "")

	_, err := skill.UserDirectory(skill.ClaudeCode, fixtureInvocation)
	if err == nil {
		t.Skip("this platform found a home directory with the variable cleared, so this check cannot force the refusal")
	}
	if !strings.Contains(err.Error(), fixtureInvocation) {
		t.Errorf("the refusal does not name the invocation it was given: %v", err)
	}
}

func TestFrontMatterCarriesThePersonOnlyFieldAndTheNameAndDescription(t *testing.T) {
	got := skill.FrontMatter(skill.ClaudeCode, "socrates-x", "Does x.")
	want := "---\nname: socrates-x\ndescription: Does x.\ndisable-model-invocation: true\n---"
	if got != want {
		t.Errorf("front matter is %q, want %q", got, want)
	}
}

func TestClaudeCodeAndCodexAreTheHostsStated(t *testing.T) {
	hosts := skill.Hosts()
	if len(hosts) != 2 || hosts[0] != skill.ClaudeCode || hosts[1] != skill.Codex {
		t.Errorf("Hosts() = %+v, want Claude Code and then Codex", hosts)
	}
}

func TestCodexStatesItsDirectoriesAndHowAPersonStartsASkill(t *testing.T) {
	home := filepath.Join("some", "home")
	if got, want := skill.UserDirectoryUnder(skill.Codex, home), filepath.Join(home, ".agents", "skills"); got != want {
		t.Errorf("Codex's skills directory is %q, want %q", got, want)
	}
	if got, want := skill.Codex.AgentsDirectoryUnder(home), filepath.Join(home, ".codex", "agents"); got != want {
		t.Errorf("Codex's agents directory is %q, want %q", got, want)
	}
	if got, want := skill.Codex.Starts("socrates-x"), "$socrates-x"; got != want {
		t.Errorf("Codex starts a skill with %q, want %q", got, want)
	}
}

func TestAHostIsFoundByTheValueOfTheHostFlag(t *testing.T) {
	for _, host := range skill.Hosts() {
		found, ok := skill.HostFlagged(host.Flag)
		if !ok || found != host {
			t.Errorf("HostFlagged(%q) = %+v, %t, want %+v", host.Flag, found, ok, host)
		}
	}
	if _, ok := skill.HostFlagged("emacs"); ok {
		t.Error("HostFlagged found a host for a value no host states")
	}
}

func TestCodexFrontMatterCarriesTheNameAndTheDescriptionAndNothingElse(t *testing.T) {
	got := skill.FrontMatter(skill.Codex, "socrates-x", "Does x.")
	want := "---\nname: socrates-x\ndescription: Does x.\n---"
	if got != want {
		t.Errorf("front matter is %q, want %q", got, want)
	}
}

func TestInstallWritesThePolicyBesideACodexSkillAndBesideNoOther(t *testing.T) {
	codexDirectory := t.TempDir()
	claudeDirectory := t.TempDir()
	codexSkill := skill.Skill{Name: "one", Host: skill.Codex, Document: func(host skill.Host) string { return skill.FrontMatter(host, "one", "x") + "\n" }, Ceiling: 100}
	claudeSkill := skill.Skill{Name: "one", Host: skill.ClaudeCode, Document: func(host skill.Host) string { return skill.FrontMatter(host, "one", "x") + "\n" }, Ceiling: 100}

	if _, err := skill.Install(codexDirectory, []skill.Skill{codexSkill}, fixtureInvocation); err != nil {
		t.Fatalf("installing for Codex: %v", err)
	}
	if _, err := skill.Install(claudeDirectory, []skill.Skill{claudeSkill}, fixtureInvocation); err != nil {
		t.Fatalf("installing for Claude Code: %v", err)
	}

	policy, err := os.ReadFile(filepath.Join(codexDirectory, "one", "agents", "openai.yaml"))
	if err != nil {
		t.Fatalf("Codex's skill has no agents/openai.yaml beside it: %v", err)
	}
	if want := "policy:\n  allow_implicit_invocation: false\n"; string(policy) != want {
		t.Errorf("agents/openai.yaml holds %q, want %q", policy, want)
	}
	if _, err := os.Stat(filepath.Join(claudeDirectory, "one", "agents")); !os.IsNotExist(err) {
		t.Errorf("Claude Code's skill has an agents directory beside it: %v", err)
	}
}

func TestCodexReadsItsPersonOnlySwitchFromThePolicyBesideTheSkill(t *testing.T) {
	installed := func(t *testing.T, policy *string) string {
		t.Helper()
		directory := t.TempDir()
		if err := os.WriteFile(filepath.Join(directory, skill.FileName), []byte(skill.FrontMatter(skill.Codex, "x", "x")+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if policy != nil {
			if err := os.MkdirAll(filepath.Join(directory, "agents"), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(directory, "agents", "openai.yaml"), []byte(*policy), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		return directory
	}
	held := "policy:\n  allow_implicit_invocation: false\n"
	allowed := "policy:\n  allow_implicit_invocation: true\n"
	silent := "policy: {}\n"
	broken := "policy: [\n"

	if !skill.Codex.PersonOnly(installed(t, &held)) {
		t.Error("Codex does not read a skill whose policy forbids implicit invocation as person-only")
	}
	for name, policy := range map[string]*string{"no policy file": nil, "an allowing policy": &allowed, "a policy that says nothing": &silent, "a policy that cannot be read": &broken} {
		if skill.Codex.PersonOnly(installed(t, policy)) {
			t.Errorf("Codex reads a skill with %s as person-only", name)
		}
	}
}

func TestCodexDoesNotReadClaudeCodesFieldAsItsSwitchNorClaudeCodeTheCodexPolicy(t *testing.T) {
	directory := t.TempDir()
	if err := os.WriteFile(filepath.Join(directory, skill.FileName), []byte(skill.FrontMatter(skill.ClaudeCode, "x", "x")+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if skill.Codex.PersonOnly(directory) {
		t.Error("Codex reads Claude Code's front matter field as its own switch")
	}

	other := t.TempDir()
	if err := os.WriteFile(filepath.Join(other, skill.FileName), []byte(skill.FrontMatter(skill.Codex, "x", "x")+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(other, "agents"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(other, "agents", "openai.yaml"), []byte("policy:\n  allow_implicit_invocation: false\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if skill.ClaudeCode.PersonOnly(other) {
		t.Error("Claude Code reads Codex's policy as its own switch")
	}
}

func TestAHostIsChosenFromTheHostFlagAndTheDirectoriesNamed(t *testing.T) {
	everyHost, err := skill.Selected("")
	if err != nil || len(everyHost) != len(skill.Hosts()) {
		t.Errorf("no host named and no directory chosen selected %v, %v, want every host", everyHost, err)
	}

	only, err := skill.Selected(skill.Codex.Flag, skill.Directory{Flag: "to", Value: "somewhere"})
	if err != nil || len(only) != 1 || only[0] != skill.Codex {
		t.Errorf("--host codex with --to selected %v, %v, want only Codex", only, err)
	}

	for _, directory := range []skill.Directory{{Flag: "to", Value: "somewhere"}, {Flag: "agents-to", Value: "somewhere"}} {
		_, err = skill.Selected("", directory)
		if err == nil || !strings.Contains(err.Error(), "--host") || !strings.Contains(err.Error(), "--"+directory.Flag) {
			t.Errorf("--%s without --host was not refused naming both flags: %v", directory.Flag, err)
		}
	}

	_, err = skill.Selected("emacs")
	if err == nil || !strings.Contains(err.Error(), "emacs") || !strings.Contains(err.Error(), "--host codex") {
		t.Errorf("--host emacs was not refused naming the hosts there are: %v", err)
	}
}

func TestAHostStatesWhereItsAgentsAreAndHowAPersonStartsASkill(t *testing.T) {
	home := filepath.Join("some", "home")
	if got, want := skill.ClaudeCode.AgentsDirectoryUnder(home), filepath.Join(home, ".claude", "agents"); got != want {
		t.Errorf("Claude Code's agents directory is %q, want %q", got, want)
	}
	if got, want := skill.ClaudeCode.Starts("socrates-x"), "/socrates-x"; got != want {
		t.Errorf("Claude Code starts a skill with %q, want %q", got, want)
	}
}

func TestAnInstalledSkillIsAnnouncedWithHowItsHostStartsIt(t *testing.T) {
	written, err := skill.Install(t.TempDir(), []skill.Skill{documentSkill("one", "# One\n")}, fixtureInvocation)
	if err != nil {
		t.Fatalf("installing: %v", err)
	}
	if !strings.HasSuffix(written[0].Announcement(), "by typing /one") {
		t.Errorf("the announcement is %q, and it should end with how Claude Code starts the skill", written[0].Announcement())
	}
}

func TestClaudeCodeReadsItsOwnPersonOnlySwitchFromAnInstalledSkill(t *testing.T) {
	held := t.TempDir()
	open := t.TempDir()
	if err := os.WriteFile(filepath.Join(held, skill.FileName), []byte(skill.FrontMatter(skill.ClaudeCode, "held", "x")+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(open, skill.FileName), []byte("---\nname: open\ndescription: x\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if !skill.ClaudeCode.PersonOnly(held) {
		t.Error("Claude Code does not read a skill carrying its switch as person-only")
	}
	if skill.ClaudeCode.PersonOnly(open) {
		t.Error("Claude Code reads a skill without its switch as person-only")
	}
	if skill.ClaudeCode.PersonOnly(filepath.Join(held, "absent")) {
		t.Error("Claude Code reads a directory holding no skill as person-only")
	}
}

func TestClaudeCodesSkillsDirectoryFollowsClaudeConfigDirAndFallsBackToTheHomeDirectory(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	t.Setenv("CLAUDE_CONFIG_DIR", "")
	if got, err := skill.UserDirectory(skill.ClaudeCode, fixtureInvocation); err != nil || got != filepath.Join(home, ".claude", "skills") {
		t.Errorf("with CLAUDE_CONFIG_DIR unset the directory is %s (error %v), want it under ~/.claude", got, err)
	}

	stated := t.TempDir()
	t.Setenv("CLAUDE_CONFIG_DIR", stated)
	if got, err := skill.UserDirectory(skill.ClaudeCode, fixtureInvocation); err != nil || got != filepath.Join(stated, "skills") {
		t.Errorf("with CLAUDE_CONFIG_DIR set the directory is %s (error %v), want %s", got, err, filepath.Join(stated, "skills"))
	}
}

func TestCodexsSkillsDirectoryIgnoresBothConfigurationVariables(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("CLAUDE_CONFIG_DIR", t.TempDir())
	t.Setenv("CODEX_HOME", t.TempDir())

	if got, err := skill.UserDirectory(skill.Codex, fixtureInvocation); err != nil || got != filepath.Join(home, ".agents", "skills") {
		t.Errorf("Codex's skills directory is %s (error %v), want it under ~/.agents", got, err)
	}
}
