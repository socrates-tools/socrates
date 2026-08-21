package skill_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rui-ktei/socrates/internal/skill"
	"github.com/rui-ktei/socrates/internal/testing/source"
)

func committedPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(source.Root(t), ".claude", "skills", skill.Name, skill.FileName)
}

func TestDriftFailsTheSuite(t *testing.T) {
	path := committedPath(t)

	committed, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("%s is not readable, and it is the copy this repository ships: %v. Run make skill to write it", path, err)
	}
	if string(committed) != skill.Document() {
		t.Errorf("%s is not what the skill's source generates: run make skill to rewrite it", path)
	}
}
