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

func TestClaudeCodeIsTheOnlyHostStated(t *testing.T) {
	hosts := skill.Hosts()
	if len(hosts) != 1 || hosts[0] != skill.ClaudeCode {
		t.Errorf("Hosts() = %+v, want only Claude Code", hosts)
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

func TestAHostReadsItsOwnPersonOnlySwitchFromAnInstalledSkill(t *testing.T) {
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
