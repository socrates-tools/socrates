package render_test

import (
	"path/filepath"
	"slices"
	"testing"

	"github.com/socrates-devtool/socrates/internal/testing/source"
)

func TestTheRendererOpensNoSocketAndReadsNoClock(t *testing.T) {
	renderer := filepath.Join("internal", "render")
	allowed := []string{
		"bytes",
		"crypto/sha256",
		"embed",
		"encoding/base64",
		"encoding/json",
		"errors",
		"fmt",
		"html/template",
		"io",
		"slices",
		"strconv",
		"strings",
		"github.com/tdewolff/parse/v2",
		"github.com/tdewolff/parse/v2/css",
		source.ModulePath + "/internal/document",
		source.ModulePath + "/internal/sanitise",
		source.ModulePath + "/internal/stylesheet",
		"golang.org/x/net/html",
		"golang.org/x/net/html/atom",
	}

	for _, imported := range source.Imports(t, source.Root(t), renderer) {
		if !slices.Contains(allowed, imported) {
			t.Errorf("the renderer reaches %s, and it renders a document and nothing else", imported)
		}
	}
}
