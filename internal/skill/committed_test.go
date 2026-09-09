package skill_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/socrates-tools/socrates-review/internal/skill"
	"github.com/socrates-tools/socrates-review/internal/testing/source"
)

func committedPath(t *testing.T, named string) string {
	t.Helper()
	return filepath.Join(source.Root(t), ".claude", "skills", named, skill.FileName)
}

func TestDriftFailsTheSuite(t *testing.T) {
	for _, one := range skill.All() {
		path := committedPath(t, one.Name)

		committed, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("%s is not readable, and it is the copy this repository ships: %v. Run make skill to write it", path, err)
		}
		if string(committed) != one.Document() {
			t.Errorf("%s is not what the skill's source generates: run make skill to rewrite it", path)
		}
	}
}
