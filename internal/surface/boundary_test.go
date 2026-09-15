package surface_test

import (
	"path/filepath"
	"slices"
	"testing"

	"github.com/socrates-tools/socrates/internal/testing/source"
)

func TestTheSurfaceStartsNoProcess(t *testing.T) {
	// scenario: The surface runs nothing
	forbidden := []string{"os/exec", "os/signal", "syscall", "plugin"}

	for _, imported := range source.Imports(t, source.Root(t), filepath.Join("internal", "surface")) {
		if slices.Contains(forbidden, imported) {
			t.Errorf("the surface imports %s, and it composes a command and never runs one", imported)
		}
	}
}
