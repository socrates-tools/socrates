package assets_test

import (
	"io/fs"
	"testing"

	"github.com/ktei/socrates/internal/assets"
	"github.com/ktei/socrates/internal/convention"
	"github.com/ktei/socrates/internal/document"
)

func TestTheLanguagesStatedAreTheLanguagesTheBundleRegisters(t *testing.T) {
	client, err := fs.ReadFile(assets.Files(), convention.ClientBundle)
	if err != nil {
		t.Fatalf("reading the embedded client: %v", err)
	}

	registered := convention.LanguagesTheBundleRegisters(string(client))
	if len(registered) == 0 {
		t.Fatalf("the bundle that ships registers no language at all, so this test proves nothing")
	}

	statedOnly, registeredOnly := convention.LanguagesOnOneSideOnly(document.Languages(), registered)
	for _, named := range statedOnly {
		t.Errorf("a document naming %q is served and renders as plain text, and the bundle registers no such language: run make assets", named)
	}
	for _, named := range registeredOnly {
		t.Errorf("a document naming %q is refused, and the bundle registers it: run make assets", named)
	}
}
