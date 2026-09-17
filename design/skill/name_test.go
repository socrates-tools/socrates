package skill_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/design/skill"
	skillkit "github.com/socrates-tools/socrates/internal/skill"
)

func TestTheSkillIsInstalledUnderItsOwnName(t *testing.T) {
	directory := t.TempDir()

	written, err := skillkit.Install(directory, skill.All(), skillkit.Dispatcher+" design skill")
	if err != nil {
		t.Fatalf("installing the skill: %v", err)
	}
	if len(written) != 1 {
		t.Fatalf("installing wrote %d skills, and design states one", len(written))
	}

	installed := written[0]
	if installed.Name != skill.Name {
		t.Errorf("the skill was installed as %q, want %q", installed.Name, skill.Name)
	}
	if named := filepath.Base(filepath.Dir(installed.Path)); named != skill.Name {
		t.Errorf("the skill is installed under %q", named)
	}
	if _, err := os.Stat(installed.Path); err != nil {
		t.Errorf("nothing is at %s: %v", installed.Path, err)
	}

	declared := ""
	for line := range strings.SplitSeq(skill.Document(), "\n") {
		if stated, isName := strings.CutPrefix(line, "name: "); isName {
			declared = strings.TrimSpace(stated)
			break
		}
	}
	if declared != skill.Name {
		t.Errorf("the skill declares itself %q while it is installed as %q, so one of the two is a second copy of the name", declared, skill.Name)
	}
}

func TestTheSkillNamesInvocationsInBothNamespacesItCovers(t *testing.T) {
	named := skillkit.Named(skill.All(), skill.Namespaces)
	if len(named) == 0 {
		t.Fatal("the skill names no invocation at all, so this check reads nothing")
	}

	seen := map[string]bool{}
	for _, invocation := range named {
		seen[invocation.Binary] = true
	}
	for _, binary := range skill.Binaries() {
		if !seen[binary] {
			t.Errorf("the skill names nothing of %s, and a binary nothing is attributed to is a binary this check never reads", binary)
		}
	}
}
