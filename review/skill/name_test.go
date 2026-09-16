package skill_test

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/convention"
	"github.com/socrates-tools/socrates/review/skill"
)

func documentOf(t *testing.T, named string) string {
	t.Helper()

	for _, one := range skill.All() {
		if one.Name == named {
			return one.Document()
		}
	}
	t.Fatalf("this repository states no skill called %s", named)
	return ""
}

func TestOneConstantNamesTheSkillAndReadsItsInvocations(t *testing.T) {
	directory := t.TempDir()

	written, err := skill.Install(directory)
	if err != nil {
		t.Fatalf("installing the skills: %v", err)
	}
	if len(written) != len(skill.All()) {
		t.Fatalf("installing wrote %d skills, and this repository states %d", len(written), len(skill.All()))
	}

	everything := ""
	for _, installed := range written {
		if named := filepath.Base(filepath.Dir(installed.Path)); named != installed.Name {
			t.Errorf("the %s skill is installed under %q", installed.Name, named)
		}

		document := documentOf(t, installed.Name)
		everything += document

		declared := ""
		for line := range strings.SplitSeq(document, "\n") {
			if stated, isName := strings.CutPrefix(line, "name: "); isName {
				declared = strings.TrimSpace(stated)
				break
			}
		}
		if declared != installed.Name {
			t.Errorf("the skill declares itself %q while it is installed as %q, so one of the two is a second copy of the name", declared, installed.Name)
		}
	}

	named := skill.Named()
	if len(named) == 0 {
		t.Fatal("the skills are read to name no invocation at all, so this check reads nothing")
	}
	for _, invocation := range named {
		if !strings.Contains(everything, typedInvocation(invocation.Binary, invocation.Command)) {
			t.Errorf("a skill is read to name %q, which appears under no other name than %q", invocation.Command, invocation.Binary)
		}
	}
}

func typedInvocation(binary, command string) string {
	switch binary {
	case skill.Name:
		return skill.Dispatcher + " review " + command
	case skill.Surface:
		return skill.Dispatcher + " agora " + command
	}
	return binary + " " + command
}

func TestTheSkillNamesTheOneNameTheToolHas(t *testing.T) {
	document := skill.Document()

	if skill.Name != convention.Repository {
		t.Errorf("the skill is named %q and the repository is named %q, and the skill keeps that name because .claude/skills is a flat namespace with no dispatcher of its own", skill.Name, convention.Repository)
	}
	if !strings.Contains(document, skill.Dispatcher+" review begin") {
		t.Error("the skill names no invocation of socrates review at all")
	}
	for _, invocation := range skill.Named() {
		if invocation.Binary == skill.Name || invocation.Binary == skill.Surface {
			continue
		}
		t.Errorf("the skill writes an invocation under %q, and the tool has no namespace but review or agora for one to use instead", invocation.Binary)
	}
}

func TestTheSkillIsInstalledUnderTheProductsName(t *testing.T) {
	directory := t.TempDir()

	written, err := skill.Install(directory)
	if err != nil {
		t.Fatalf("installing the skills: %v", err)
	}

	for _, want := range []string{
		filepath.Join(directory, convention.Repository, skill.FileName),
		filepath.Join(directory, convention.Repository+"-handover", skill.FileName),
	} {
		if !slices.ContainsFunc(written, func(installed skill.Installed) bool { return installed.Path == want }) {
			t.Errorf("nothing was installed at %s, and what was installed is %+v", want, written)
		}
		if _, err := os.Stat(want); err != nil {
			t.Errorf("nothing is at %s: %v", want, err)
		}
	}
}
