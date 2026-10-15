package skill_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/design/skill"
	skillkit "github.com/socrates-tools/socrates/internal/skill"
)

func TestEachHostsRenderingNamesItsOwnDirectoryAndHostWhereItSaysToInstallElsewhere(t *testing.T) {
	for _, one := range skill.All() {
		installing := sectionOf(t, one.Rendered(), "## Installing elsewhere")

		if want := "~/" + filepath.ToSlash(one.Host.SkillsDirectory); !strings.Contains(installing, want) {
			t.Errorf("%s for %s does not name %s under Installing elsewhere:\n%s", one.Name, one.Host.Name, want, installing)
		}
		if want := "socrates design skill --host " + one.Host.Flag + " --to <directory>"; !strings.Contains(installing, want) {
			t.Errorf("%s for %s does not say %q under Installing elsewhere:\n%s", one.Name, one.Host.Name, want, installing)
		}
		for _, other := range skillkit.Hosts() {
			if other != one.Host && strings.Contains(one.Rendered(), filepath.ToSlash(other.SkillsDirectory)) {
				t.Errorf("%s for %s names %s's skills directory", one.Name, one.Host.Name, other.Name)
			}
		}
	}
}
