package assets_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/agora/assets"
)

var customPropertyDeclared = regexp.MustCompile(`(--[a-z0-9-]+)\s*:`)

func TestTheTokensStylesheetHoldsEveryCustomPropertyAnIllustrationWasWrittenAgainst(t *testing.T) {
	// scenario: An illustration loading /assets/tokens.css finds every custom property it found before the tokens moved into the library
	held, err := os.ReadFile(filepath.Join("testdata", "custom-properties.txt"))
	if err != nil {
		t.Fatalf("reading the custom properties illustrations rely on: %v", err)
	}
	served, err := fs.ReadFile(assets.Files(), "tokens.css")
	if err != nil {
		t.Fatalf("reading the embedded tokens stylesheet: %v", err)
	}

	declared := map[string]bool{}
	for _, found := range customPropertyDeclared.FindAllStringSubmatch(string(served), -1) {
		declared[found[1]] = true
	}

	for _, name := range strings.Fields(string(held)) {
		if !declared[name] {
			t.Errorf("tokens.css declares no %s, and an illustration written against it would lose it", name)
		}
	}
}

func TestTheTokensStylesheetStatesTheSpacingStepSoThatStepFourIsSixteenPixels(t *testing.T) {
	// scenario: The spacing scale is one step of four pixels, so p-4 is 16px
	served, err := fs.ReadFile(assets.Files(), "tokens.css")
	if err != nil {
		t.Fatalf("reading the embedded tokens stylesheet: %v", err)
	}

	found := regexp.MustCompile(`--spacing\s*:\s*([0-9.]+)px`).FindStringSubmatch(string(served))
	if found == nil {
		t.Fatal("tokens.css states no --spacing in pixels, and a spacing utility multiplies it")
	}
	step, err := strconv.ParseFloat(found[1], 64)
	if err != nil {
		t.Fatalf("reading the spacing step %q: %v", found[1], err)
	}
	if got := step * 4; got != 16 {
		t.Errorf("step 4 of the spacing scale is %vpx, want 16px", got)
	}
}
