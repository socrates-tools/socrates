package surface_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/agora/surface"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

func TestTheCommittedClientTypesAreWhatGoWouldWrite(t *testing.T) {
	// scenario: The client's types are the ones its Go structs generate, and no Node is needed to know
	generated, err := surface.ClientTypes()
	if err != nil {
		t.Fatalf("generating the client's types: %v", err)
	}

	committed, err := os.ReadFile(filepath.Join(source.Root(t), filepath.FromSlash(surface.ClientTypesFile)))
	if err != nil {
		t.Fatalf("reading %s: %v", surface.ClientTypesFile, err)
	}

	if difference := cmp.Diff(string(generated), string(committed)); difference != "" {
		t.Errorf("%s differs from what the Go structs generate (-generated +committed). Run make assets and commit what it writes:\n%s", surface.ClientTypesFile, difference)
	}
}

func TestTheClientTypesNameEveryThingThatCrosses(t *testing.T) {
	// scenario: The generated types hold the shape of the document and of the server part
	generated, err := surface.ClientTypes()
	if err != nil {
		t.Fatalf("generating the client's types: %v", err)
	}

	for _, expected := range []string{
		"export interface Served {\n  document: ServedDocument\n  standing: Standing\n}\n",
		"export interface Standing {\n  revision: number\n  over: boolean\n  illustrations: Record<string, string>\n  answered: string[]\n  files_per_tray: number\n}\n",
		"export type Mode = 'single-choice' | 'multiple-choice'\n",
		"  attached?: string[]\n",
		"  recommended: boolean\n",
	} {
		if !strings.Contains(string(generated), expected) {
			t.Errorf("the generated types do not hold %q:\n%s", expected, generated)
		}
	}
	if strings.Contains(string(generated), "command") {
		t.Errorf("the generated types name a command, and none reaches the page:\n%s", generated)
	}
}
