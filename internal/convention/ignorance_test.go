package convention_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/internal/convention"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

func TestTheSurfaceKnowsNothingOfTheOtherTool(t *testing.T) {
	read, err := convention.WordsOfTheOtherTool(source.Root(t), convention.SurfaceDirectories, convention.TheOtherToolNamedIn)
	if err != nil {
		t.Fatalf("reading what the surface names: %v", err)
	}

	for _, named := range read.Naming {
		t.Errorf("%s, and the surface takes a document, serves it, waits and prints what came back", named)
	}
	for _, stale := range read.Stale {
		t.Errorf("%s is named in convention.TheOtherToolNamedIn and names nothing of the other tool, so the exemption outlived its reason", stale)
	}
}

func TestAFileNamedAsNamingTheOtherToolIsExemptAndAStaleOneIsReported(t *testing.T) {
	root := t.TempDir()
	writeNested(t, filepath.Join(root, "held", "stated.go"), "package held\n\nconst hint = \"the reviewer will not answer it\"\n")
	writeNested(t, filepath.Join(root, "held", "other.go"), "package held\n\nconst said = \"the reviewer read it\"\n")
	writeNested(t, filepath.Join(root, "held", "quiet.go"), "package held\n\nconst said = \"nothing to decide\"\n")

	read, err := convention.WordsOfTheOtherTool(root, []string{"held"}, map[string]string{
		filepath.Join("held", "stated.go"): "what a caller's audience is announced as",
		filepath.Join("held", "quiet.go"):  "an exemption nothing needs",
	})
	if err != nil {
		t.Fatalf("reading what a fabricated tree names: %v", err)
	}

	if diff := cmp.Diff([]string{`held/other.go:3 names "reviewer"`}, read.Naming); diff != "" {
		t.Errorf("what a tree with a named file reports differs (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]string{filepath.Join("held", "quiet.go")}, read.Stale); diff != "" {
		t.Errorf("what a stale exemption reports differs (-want +got):\n%s", diff)
	}
}

func TestAWordOfTheOtherToolIsReported(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "held"), 0o755); err != nil {
		t.Fatalf("making a fabricated tree: %v", err)
	}
	write(t, filepath.Join(root, "held", "knowing.go"), "package held\n\nconst said = \"the reviewer read the finding in that run\"\n")

	read, err := convention.WordsOfTheOtherTool(root, []string{"held"}, nil)
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

func TestALineDrawnRoundIsNotAWordOfTheOtherTool(t *testing.T) {
	root := t.TempDir()
	writeNested(t, filepath.Join(root, "held", "drawn.go"), "package held\n\nconst drawn = `<svg stroke-linecap=\"round\" stroke-linejoin=\"round\"></svg>`\n")
	writeNested(t, filepath.Join(root, "held", "drawn.css"), ".icon {\n  stroke-linecap: round;\n  stroke-linejoin:round;\n}\n\n.said::after {\n  content: 'each round';\n}\n")

	read, err := convention.WordsOfTheOtherTool(root, []string{"held"}, nil)
	if err != nil {
		t.Fatalf("reading what a fabricated tree names: %v", err)
	}

	want := []string{`held/drawn.css:7 names "round"`}
	if diff := cmp.Diff(want, read.Naming); diff != "" {
		t.Errorf("what a file drawing round line ends reports differs (-want +got):\n%s", diff)
	}
}

func TestAnImportOfAnotherToolsPackageIsStillReported(t *testing.T) {
	root := t.TempDir()
	writeNested(t, filepath.Join(root, "go.mod"), "module example.com/suite\n")
	writeNested(t, filepath.Join(root, "held", "knowing.go"),
		"package held\n\nimport (\n\t_ \"example.com/suite/other-tool/document\"\n\t_ \"example.com/suite/other-tool/reviewer\"\n)\n")

	read, err := convention.WordsOfTheOtherTool(root, []string{"held"}, nil)
	if err != nil {
		t.Fatalf("reading what a fabricated tree names: %v", err)
	}

	want := []string{`held/knowing.go:5 names "reviewer"`}
	if diff := cmp.Diff(want, read.Naming); diff != "" {
		t.Errorf("what a file importing the other tool reports differs (-want +got):\n%s", diff)
	}
}

func TestEveryPackageOfTheSurfaceIsReadForTheOtherToolsWords(t *testing.T) {
	root := source.Root(t)

	unread, err := convention.PackagesOutside(root, convention.ToolDirectory("agora"), convention.SurfaceDirectories)
	if err != nil {
		t.Fatalf("reading which of the surface's packages are read: %v", err)
	}
	for _, unread := range unread {
		t.Errorf("%s holds Go source and is not in convention.SurfaceDirectories, so what it names goes unread", unread)
	}
}

func TestAPackageOfTheSurfaceNobodyListedIsReported(t *testing.T) {
	root := t.TempDir()
	writeNested(t, filepath.Join(root, "surface", "listed", "listed.go"), "package listed\n")
	writeNested(t, filepath.Join(root, "surface", "added", "added.go"), "package added\n")
	writeNested(t, filepath.Join(root, "surface", "web", "node_modules", "left", "left.go"), "package left\n")

	unread, err := convention.PackagesOutside(root, "surface", []string{filepath.Join("surface", "listed")})
	if err != nil {
		t.Fatalf("reading which packages of a fabricated tree are read: %v", err)
	}

	if diff := cmp.Diff([]string{filepath.Join("surface", "added")}, unread); diff != "" {
		t.Errorf("what a package nobody listed reports differs (-want +got):\n%s", diff)
	}
}
