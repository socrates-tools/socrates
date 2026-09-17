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
	return skill.Skill{Name: name, Document: func() string { return document }, Ceiling: len(document)}
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
		if string(content) != one.Document() {
			t.Errorf("%s holds %q, want %q", want, content, one.Document())
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

func TestUserDirectoryUnderJoinsTheVendorAndSkillsDirectories(t *testing.T) {
	home := filepath.Join("some", "home")
	want := filepath.Join(home, ".claude", "skills")
	if got := skill.UserDirectoryUnder(home); got != want {
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

	_, err := skill.UserDirectory(fixtureInvocation)
	if err == nil {
		t.Skip("this platform found a home directory with the variable cleared, so this check cannot force the refusal")
	}
	if !strings.Contains(err.Error(), fixtureInvocation) {
		t.Errorf("the refusal does not name the invocation it was given: %v", err)
	}
}
