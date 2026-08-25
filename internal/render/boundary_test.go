package render_test

import (
	"path/filepath"
	"slices"
	"testing"

	"github.com/rui-ktei/socrates/internal/testing/source"
)

func TestTheRendererOpensNoSocketAndReadsNoClock(t *testing.T) {
	renderer := filepath.Join("internal", "render")
	allowed := []string{
		"bytes",
		"embed",
		"encoding/json",
		"errors",
		"fmt",
		"html/template",
		"slices",
		"strconv",
		"strings",
		source.ModulePath + "/internal/document",
		source.ModulePath + "/internal/sanitise",
		"golang.org/x/net/html",
		"golang.org/x/net/html/atom",
	}

	for _, imported := range source.Imports(t, source.Root(t), renderer) {
		if !slices.Contains(allowed, imported) {
			t.Errorf("the renderer reaches %s, and it renders a document and nothing else", imported)
		}
	}
}
