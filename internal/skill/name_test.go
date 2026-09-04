package skill_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates-review/internal/convention"
	"github.com/socrates-tools/socrates-review/internal/skill"
)

func TestOneConstantNamesTheSkillAndReadsItsInvocations(t *testing.T) {
	directory := t.TempDir()

	path, err := skill.Install(directory)
	if err != nil {
		t.Fatalf("installing the skill: %v", err)
	}
	if named := filepath.Base(filepath.Dir(path)); named != skill.Name {
		t.Errorf("the skill is installed under %q rather than %q", named, skill.Name)
	}

	document := skill.Document()
	declared := ""
	for line := range strings.SplitSeq(document, "\n") {
		if stated, isName := strings.CutPrefix(line, "name: "); isName {
			declared = strings.TrimSpace(stated)
			break
		}
	}
	if declared != skill.Name {
		t.Errorf("the skill declares itself %q while it is installed as %q, so one of the two is a second copy of the name", declared, skill.Name)
	}

	named := skill.Named()
	if len(named) == 0 {
		t.Fatal("the skill is read to name no invocation at all, so this check reads nothing")
	}
	for _, invocation := range named {
		if !strings.Contains(document, invocation.Binary+" "+invocation.Command) {
			t.Errorf("the skill is read to name %q, which appears in it under no other name than %q", invocation.Command, invocation.Binary)
		}
	}
}

func TestTheSkillNamesTheOneNameTheToolHas(t *testing.T) {
	document := skill.Document()

	if skill.Name != convention.Repository {
		t.Errorf("the skill is named %q and the tool is named %q, so the skill teaches a command that is not the one the tool installs", skill.Name, convention.Repository)
	}
	if !strings.Contains(document, skill.Name+" review") {
		t.Errorf("the skill names no invocation of %s at all", skill.Name)
	}
	for _, invocation := range skill.Named() {
		if invocation.Binary == skill.Name || invocation.Binary == skill.Surface {
			continue
		}
		t.Errorf("the skill writes an invocation under %q, and the tool has no name but %q for one to use instead", invocation.Binary, skill.Name)
	}
}

func TestTheSkillIsInstalledUnderTheProductsName(t *testing.T) {
	directory := t.TempDir()

	path, err := skill.Install(directory)
	if err != nil {
		t.Fatalf("installing the skill: %v", err)
	}

	want := filepath.Join(directory, convention.Repository, skill.FileName)
	if path != want {
		t.Errorf("the skill is at %s, want %s", path, want)
	}
	if _, err := os.Stat(want); err != nil {
		t.Errorf("nothing is at %s: %v", want, err)
	}
}
