package skill_test

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/design/skill"
	skillkit "github.com/socrates-tools/socrates/internal/skill"
)

func TestEachDesignSkillIsInstalledUnderItsOwnName(t *testing.T) {
	directory := t.TempDir()

	written, err := skillkit.Install(directory, skill.All(), skillkit.Dispatcher+" design skill")
	if err != nil {
		t.Fatalf("installing the skills: %v", err)
	}

	want := []string{skill.Name, skill.TicketsName}
	var got []string
	for _, installed := range written {
		got = append(got, installed.Name)
	}
	if !slices.Equal(got, want) {
		t.Fatalf("installing wrote %v, and design states %v", got, want)
	}

	for position, installed := range written {
		if named := filepath.Base(filepath.Dir(installed.Path)); named != installed.Name {
			t.Errorf("the %s skill is installed under %q", installed.Name, named)
		}
		if _, err := os.Stat(installed.Path); err != nil {
			t.Errorf("nothing is at %s: %v", installed.Path, err)
		}

		declared := ""
		for line := range strings.SplitSeq(skill.All()[position].Document(), "\n") {
			if stated, isName := strings.CutPrefix(line, "name: "); isName {
				declared = strings.TrimSpace(stated)
				break
			}
		}
		if declared != installed.Name {
			t.Errorf("the skill declares itself %q while it is installed as %q, so one of the two is a second copy of the name", declared, installed.Name)
		}
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
