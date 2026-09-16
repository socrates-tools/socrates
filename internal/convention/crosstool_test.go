package convention_test

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/internal/convention"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

func TestNoToolImportsAnotherTool(t *testing.T) {
	root := source.Root(t)

	module, err := convention.ModulePath(root)
	if err != nil {
		t.Fatalf("reading the module path: %v", err)
	}

	crossing, err := convention.CrossToolImports(root, module, convention.Tools)
	if err != nil {
		t.Fatalf("reading what each tool imports: %v", err)
	}
	for _, crossed := range crossing {
		t.Errorf("%s is %s's own and imports %s, which is %s's: no tool imports another tool", crossed.File, crossed.Tool, crossed.Imported, crossed.Other)
	}
}

func TestAToolImportingAnotherIsReported(t *testing.T) {
	root := t.TempDir()
	tools := []convention.Tool{{Namespace: "review", Directory: "review"}, {Namespace: "agora", Directory: "agora"}}

	mkdir(t, filepath.Join(root, "review", "command"))
	mkdir(t, filepath.Join(root, "agora", "surface"))

	write(t, filepath.Join(root, "review", "command", "root.go"), "package command\n\nimport \"example.com/two/agora/surface\"\n\nvar _ = surface.Serve\n")
	write(t, filepath.Join(root, "agora", "surface", "surface.go"), "package surface\n\nfunc Serve() {}\n")

	crossing, err := convention.CrossToolImports(root, "example.com/two", tools)
	if err != nil {
		t.Fatalf("reading a fabricated tree: %v", err)
	}

	want := []convention.CrossToolImport{{
		File:     filepath.Join("review", "command", "root.go"),
		Tool:     "review",
		Imported: "example.com/two/agora/surface",
		Other:    "agora",
	}}
	if diff := cmp.Diff(want, crossing); diff != "" {
		t.Errorf("what a tool importing another reports differs (-want +got):\n%s", diff)
	}
}

func TestAToolImportingItsOwnPackageIsNotReported(t *testing.T) {
	root := t.TempDir()
	tools := []convention.Tool{{Namespace: "review", Directory: "review"}}

	mkdir(t, filepath.Join(root, "review", "command"))
	mkdir(t, filepath.Join(root, "review", "config"))

	write(t, filepath.Join(root, "review", "command", "root.go"), "package command\n\nimport \"example.com/two/review/config\"\n\nvar _ = config.Load\n")
	write(t, filepath.Join(root, "review", "config", "config.go"), "package config\n\nfunc Load() {}\n")

	crossing, err := convention.CrossToolImports(root, "example.com/two", tools)
	if err != nil {
		t.Fatalf("reading a fabricated tree: %v", err)
	}
	if len(crossing) != 0 {
		t.Errorf("a tool importing its own package is reported as crossing: %v", crossing)
	}
}

func TestABuildTaggedFileImportingAnotherToolIsReported(t *testing.T) {
	root := t.TempDir()
	tools := []convention.Tool{{Namespace: "review", Directory: "review"}, {Namespace: "agora", Directory: "agora"}}

	mkdir(t, filepath.Join(root, "review", "fidelity"))
	mkdir(t, filepath.Join(root, "agora", "surface"))

	write(t, filepath.Join(root, "review", "fidelity", "live.go"),
		"//go:build live\n\npackage fidelity\n\nimport \"example.com/two/agora/surface\"\n\nvar _ = surface.Serve\n")
	write(t, filepath.Join(root, "agora", "surface", "surface.go"), "package surface\n\nfunc Serve() {}\n")

	crossing, err := convention.CrossToolImports(root, "example.com/two", tools)
	if err != nil {
		t.Fatalf("reading a fabricated tree: %v", err)
	}

	want := []convention.CrossToolImport{{
		File:     filepath.Join("review", "fidelity", "live.go"),
		Tool:     "review",
		Imported: "example.com/two/agora/surface",
		Other:    "agora",
	}}
	if diff := cmp.Diff(want, crossing); diff != "" {
		t.Errorf("what a build-tagged non-test file importing another tool reports differs (-want +got):\n%s", diff)
	}
}

func TestATestFileImportingAnotherToolIsNotReported(t *testing.T) {
	root := t.TempDir()
	tools := []convention.Tool{{Namespace: "review", Directory: "review"}, {Namespace: "agora", Directory: "agora"}}

	mkdir(t, filepath.Join(root, "review", "fidelity"))
	mkdir(t, filepath.Join(root, "agora", "surface"))

	write(t, filepath.Join(root, "review", "fidelity", "livesmoke_test.go"),
		"//go:build livesmoke\n\npackage fidelity_test\n\nimport \"example.com/two/agora/surface\"\n\nvar _ = surface.Serve\n")
	write(t, filepath.Join(root, "agora", "surface", "surface.go"), "package surface\n\nfunc Serve() {}\n")

	crossing, err := convention.CrossToolImports(root, "example.com/two", tools)
	if err != nil {
		t.Fatalf("reading a fabricated tree: %v", err)
	}
	if len(crossing) != 0 {
		t.Errorf("a test file importing another tool is reported as crossing: %v", crossing)
	}
}

func TestAPackageImportingNeitherToolIsNotReported(t *testing.T) {
	root := t.TempDir()
	tools := []convention.Tool{{Namespace: "review", Directory: "review"}}

	mkdir(t, filepath.Join(root, "internal", "exit"))
	mkdir(t, filepath.Join(root, "review", "command"))

	write(t, filepath.Join(root, "internal", "exit", "code.go"), "package exit\n\nimport \"example.com/two/review/command\"\n\nvar _ = command.NewRootCommand\n")
	write(t, filepath.Join(root, "review", "command", "root.go"), "package command\n\nfunc NewRootCommand() {}\n")

	crossing, err := convention.CrossToolImports(root, "example.com/two", tools)
	if err != nil {
		t.Fatalf("reading a fabricated tree: %v", err)
	}
	if len(crossing) != 0 {
		t.Errorf("internal, which is not a tool, is reported as one: %v", crossing)
	}
}
