package e2e_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates-review/internal/convention"
	"github.com/socrates-tools/socrates-review/internal/testing/source"
)

func TestEveryImportResolvesUnderTheOrganisation(t *testing.T) {
	root := source.Root(t)
	module := "github.com/" + convention.Organisation + "/" + convention.Repository

	own := 0
	for path, directory := range packagesResolved(t, root) {
		if !strings.HasPrefix(directory, root+string(filepath.Separator)) && directory != root {
			continue
		}
		own++
		if path != module && !strings.HasPrefix(path, module+"/") {
			t.Errorf("the package in %s resolves as %s, and the project answers to %s", directory, path, module)
		}
	}

	if own == 0 {
		t.Fatalf("no package of this project resolved under %s", root)
	}
}

func TestTheBuildReadsTheModulePathRatherThanRepeatingIt(t *testing.T) {
	root := requireMake(t)

	declared := moduleDeclared(t, root)
	stated, occurrences := modulePathTheMakefileStates(t, root, declared)

	if stated != declared {
		t.Errorf("the Makefile builds against %q, and go.mod declares %q", stated, declared)
	}
	if occurrences != 1 {
		t.Errorf("the Makefile writes the module path %d times, and a second copy is a copy that can disagree", occurrences)
	}

	printed := versionOf(t, made(t, root))
	if commit := gitSaid(t, root, "rev-parse", "--short=7", "HEAD"); !strings.Contains(printed, commit) {
		t.Errorf("the built binary reports %q, which names no commit, so the stamp reached no symbol under %s", printed, declared)
	}
}

func packagesResolved(t *testing.T, root string) map[string]string {
	t.Helper()

	command := exec.Command("go", "list", "-deps", "-f", "{{.ImportPath}}\t{{.Dir}}", "./...")
	command.Dir = root
	command.Env = append(os.Environ(), source.Toolchain()...)

	var out, failures bytes.Buffer
	command.Stdout = &out
	command.Stderr = &failures
	if err := command.Run(); err != nil {
		t.Fatalf("listing what the module resolves: %v: %s", err, failures.String())
	}

	resolved := map[string]string{}
	for _, line := range strings.Split(strings.TrimSpace(out.String()), "\n") {
		path, directory, split := strings.Cut(line, "\t")
		if !split {
			continue
		}
		resolved[path] = directory
	}
	return resolved
}

func moduleDeclared(t *testing.T, root string) string {
	t.Helper()

	held, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		t.Fatalf("reading go.mod: %v", err)
	}

	for _, line := range strings.Split(string(held), "\n") {
		if declared, stated := strings.CutPrefix(strings.TrimSpace(line), "module "); stated {
			return strings.TrimSpace(declared)
		}
	}
	t.Fatalf("go.mod declares no module")
	return ""
}

func modulePathTheMakefileStates(t *testing.T, root, declared string) (string, int) {
	t.Helper()

	held, err := os.ReadFile(filepath.Join(root, "Makefile"))
	if err != nil {
		t.Fatalf("reading the Makefile: %v", err)
	}

	stated := ""
	for _, line := range strings.Split(string(held), "\n") {
		if assigned, is := strings.CutPrefix(strings.TrimSpace(line), "MODULE :="); is {
			stated = strings.TrimSpace(assigned)
		}
	}
	return stated, strings.Count(string(held), declared)
}
