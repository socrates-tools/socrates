package document_test

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates-review/internal/testing/source"
)

func TestTheDocumentContractImportsNothingOutward(t *testing.T) {
	allowed := []string{"bytes", "embed", "encoding/json", "errors", "fmt", "slices", "strings"}

	for _, imported := range source.Imports(t, source.Root(t), filepath.Join("internal", "document")) {
		if strings.HasPrefix(imported, source.ModulePath) {
			t.Errorf("the document contract imports %s, and what a caller writes depends on nothing else in this module", imported)
			continue
		}
		if !slices.Contains(allowed, imported) {
			t.Errorf("the document contract imports %s, and it may import only %v", imported, allowed)
		}
	}
}
