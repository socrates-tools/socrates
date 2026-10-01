package configure

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/settingsmenu"
	"github.com/socrates-tools/socrates/internal/state"
	"github.com/socrates-tools/socrates/review/config"
)

func TestTheSurfaceBlockThatIsWrongAsAWholeIsResetAsTheBlockUnderReviewsOwnLoader(t *testing.T) {
	t.Setenv(state.HomeVariable, t.TempDir())
	path := filepath.Join(t.TempDir(), "config.yaml")
	held := "diff_ceiling: 4096\nsurface:\n  bind: 127.0.0.1:18019\n  reached_at: http://localhost:18019/\n"

	written := repairedUnder(t, path, held, "surface", 2)
	if want := "diff_ceiling: 4096\n"; string(written) != want {
		t.Errorf("the repaired file holds:\n%s\nwant:\n%s", written, want)
	}
}

func TestTheNestedSettingThatIsWrongIsResetUnderReviewsOwnLoader(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	held := "diff_ceiling: 4096\nroles:\n  reviewer:\n    model: sonnet\n    temperature: 0.5\n"

	written := repairedUnder(t, path, held, "roles.reviewer.temperature", 5)
	if err := config.Parse(path, written); err != nil {
		t.Errorf("the repaired file is one the loader refuses: %v\n%s", err, written)
	}
	for _, kept := range []string{"diff_ceiling: 4096", "model: sonnet"} {
		if !strings.Contains(string(written), kept) {
			t.Errorf("%q was dropped along with the setting that was reset:\n%s", kept, written)
		}
	}
	if strings.Contains(string(written), "temperature") {
		t.Errorf("the setting that was reset is still named:\n%s", written)
	}
}

func repairedUnder(t *testing.T, path, held, key string, line int) []byte {
	t.Helper()

	repair, found, err := settingsmenu.RepairFirst(definitionAt(path), []byte(held))
	if err != nil {
		t.Fatalf("repairing a file review's loader refuses: %v", err)
	}
	if !found {
		t.Fatalf("review's loader refuses the file and nothing was named:\n%s", held)
	}
	if repair.Key != key || repair.Line != line {
		t.Errorf("the setting named is %q at line %d, and what is wrong is %q at line %d", repair.Key, repair.Line, key, line)
	}
	return repair.Written
}
