package skill_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rui-ktei/socrates/internal/convention"
	"github.com/rui-ktei/socrates/internal/skill"
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
		if !strings.Contains(document, skill.Name+" "+invocation.Command) {
			t.Errorf("the skill is read to name %q, which appears in it under no other name than %q", invocation.Command, skill.Name)
		}
	}
}

func TestTheSkillNamesTheFullCommand(t *testing.T) {
	document := skill.Document()

	if strings.Contains(document, convention.ShortName) {
		t.Errorf("the skill names the short name, and every invocation it gives is the full command:\n%s", document)
	}
	if !strings.Contains(document, skill.Name+" review") {
		t.Errorf("the skill names no invocation of %s at all", skill.Name)
	}
}

func TestTheSkillIsInstalledUnderTheProductsName(t *testing.T) {
	directory := t.TempDir()

	path, err := skill.Install(directory)
	if err != nil {
		t.Fatalf("installing the skill: %v", err)
	}

	want := filepath.Join(directory, "socrates", skill.FileName)
	if path != want {
		t.Errorf("the skill is at %s, want %s", path, want)
	}
	if _, err := os.Stat(want); err != nil {
		t.Errorf("nothing is at %s: %v", want, err)
	}
}
