package source_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/socrates-tools/socrates/internal/convention"
	"github.com/socrates-tools/socrates/internal/platform"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

const noticeThatNamesNoPlatformSpecificModule = "NOTICE\n"

func TestAModuleLinkedForOnePlatformAloneIsListedAndTheNoticeLackingItIsReported(t *testing.T) {
	root := aModuleLinkingADependencyOnlyOn(t, aBuildOtherThanTheOneRunning(t))

	linked := source.LinkedModulesIn(t, root, ".")

	var paths []string
	for _, module := range linked {
		paths = append(paths, module.Path)
	}
	missing := convention.ModulesWithoutANotice(paths, noticeThatNamesNoPlatformSpecificModule)

	if len(missing) != 1 || missing[0] != "example.test/only" {
		t.Fatalf("the notice lacking a module linked for one supported platform alone reports %v, want [example.test/only]", missing)
	}
	for _, module := range linked {
		if module.Path == "example.test/only" && module.Directory != filepath.Join(root, "only") {
			t.Errorf("the module resolves to %s, want %s", module.Directory, filepath.Join(root, "only"))
		}
	}
}

func aBuildOtherThanTheOneRunning(t *testing.T) platform.Build {
	t.Helper()

	running := platform.Build{OperatingSystem: runtime.GOOS, Architecture: runtime.GOARCH}
	for _, build := range platform.Builds() {
		if build != running {
			return build
		}
	}
	t.Fatalf("the supported set holds no build other than the one running, so this test proves nothing")
	return platform.Build{}
}

func aModuleLinkingADependencyOnlyOn(t *testing.T, build platform.Build) string {
	t.Helper()

	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	write := func(name, content string) {
		t.Helper()
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	write("go.mod", "module example.test/app\n\ngo 1.21\n\nrequire example.test/only v0.0.0\n\nreplace example.test/only => ./only\n")
	write("main.go", "package main\n\nfunc main() {}\n")
	write("only_"+build.OperatingSystem+"_"+build.Architecture+".go", "package main\n\nimport _ \"example.test/only\"\n")
	write("only/go.mod", "module example.test/only\n\ngo 1.21\n")
	write("only/only.go", "package only\n")
	return root
}
