package convention_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ktei/socrates/internal/convention"
	"github.com/ktei/socrates/internal/testing/source"
)

func TestTheSurfaceKnowsNothingOfTheOtherTool(t *testing.T) {
	read, err := convention.WordsOfTheOtherTool(source.Root(t), convention.SurfaceDirectories)
	if err != nil {
		t.Fatalf("reading what the surface names: %v", err)
	}

	for _, named := range read.Naming {
		t.Errorf("%s, and the surface takes a document, serves it, waits and prints what came back", named)
	}
}

func TestAWordOfTheOtherToolIsReported(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "held"), 0o755); err != nil {
		t.Fatalf("making a fabricated tree: %v", err)
	}
	write(t, filepath.Join(root, "held", "knowing.go"), "package held\n\nconst said = \"the reviewer read the finding in that run\"\n")

	read, err := convention.WordsOfTheOtherTool(root, []string{"held"})
	if err != nil {
		t.Fatalf("reading what a fabricated tree names: %v", err)
	}

	want := []string{
		`held/knowing.go:3 names "finding"`,
		`held/knowing.go:3 names "reviewer"`,
		`held/knowing.go:3 names "that run"`,
	}
	if diff := cmp.Diff(want, read.Naming); diff != "" {
		t.Errorf("what a file naming the other tool reports differs (-want +got):\n%s", diff)
	}
}
