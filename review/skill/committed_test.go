package skill_test

import (
	"os"
	"path/filepath"
	"testing"

	skillkit "github.com/socrates-tools/socrates/internal/skill"
	"github.com/socrates-tools/socrates/internal/testing/source"
	"github.com/socrates-tools/socrates/review/skill"
)

func committedPath(t *testing.T, host skillkit.Host, named string) string {
	t.Helper()
	return filepath.Join(host.SkillsDirectoryUnder(source.Root(t)), named, skillkit.FileName)
}

func TestDriftFailsTheSuite(t *testing.T) {
	for _, one := range skill.All() {
		path := committedPath(t, one.Host, one.Name)

		committed, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("%s is not readable, and it is the copy this repository ships: %v. Run make skill to write it", path, err)
		}
		if string(committed) != one.Rendered() {
			t.Errorf("%s is not what the skill's source generates: run make skill to rewrite it", path)
		}
	}
}
