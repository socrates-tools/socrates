package presentation_test

import (
	"path/filepath"
	"slices"
	"testing"

	"github.com/socrates-tools/socrates/internal/testing/source"
)

func TestTheRendererDecidesNothingAboutTheRun(t *testing.T) {
	renderer := filepath.Join("internal", "presentation")
	allowed := []string{"fmt", "maps", "slices", "strings", source.ModulePath + "/internal/runlog"}

	root := source.Root(t)
	for _, imported := range source.DirectImports(t, root, renderer) {
		if !slices.Contains(allowed, imported) {
			t.Errorf("the renderer imports %s, and it may import only %v", imported, allowed)
		}
	}

	for _, called := range source.PackagesCalled(t, root, renderer) {
		if called == "runlog" {
			t.Error("the renderer calls a function of the event log package, and it may name only its event type")
		}
	}
}
