package e2e_test

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates-review/internal/convention"
	"github.com/socrates-tools/socrates-review/internal/testing/harness"
	"github.com/socrates-tools/socrates-review/internal/testing/source"
)

func TestOneNameIsInstalledAndNothingIsPutBesideIt(t *testing.T) {
	root := requireMake(t)
	directory := installed(t, root)

	tool := filepath.Join(directory, convention.Repository)
	reported := versionOf(t, tool)
	if commit := gitSaid(t, root, "rev-parse", "--short=7", "HEAD"); !strings.Contains(reported, commit) {
		t.Errorf("the installed binary reports %q, which does not name the commit %s", reported, commit)
	}

	installedTool, err := os.Stat(tool)
	if err != nil {
		t.Fatalf("the tool was not installed under %s: %v", convention.Repository, err)
	}

	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatalf("reading what installing put on the path: %v", err)
	}
	for _, entry := range entries {
		if entry.Name() == convention.Repository {
			continue
		}
		beside, err := os.Stat(filepath.Join(directory, entry.Name()))
		if err != nil {
			continue
		}
		if os.SameFile(installedTool, beside) {
			t.Errorf("installing put %s beside %s and the two are one binary, and the tool answers to one name", entry.Name(), convention.Repository)
		}
	}
}

func TestWhatItPrintsNamesTheTool(t *testing.T) {
	session := harness.Start(t)

	printed, failures, code := invoked(t, session.Binary, "kaboom")
	if code != 3 {
		t.Fatalf("an unknown command was refused with %d, want 3: %s", code, failures)
	}
	if !strings.Contains(failures, convention.Repository) {
		t.Errorf("what it printed does not name %s: %q", convention.Repository, failures)
	}

	helped, helpFailures, code := invoked(t, session.Binary, "--help")
	if code != 0 {
		t.Fatalf("--help exited %d: %s", code, helpFailures)
	}
	if !strings.Contains(helped, convention.Repository) {
		t.Errorf("what --help printed does not name %s: %q", convention.Repository, helped)
	}

	for _, said := range strings.Fields(printed + failures + helped + helpFailures) {
		if said == "socrates" || said == "scrat" {
			t.Errorf("what it printed names %q, and the command a person can type is %s", said, convention.Repository)
		}
	}
}

func installed(t *testing.T, root string) string {
	t.Helper()

	directory := t.TempDir()
	command := exec.Command("make", "install")
	command.Dir = root
	command.Env = append(append(os.Environ(), source.Toolchain()...), "GOBIN="+directory)

	var failures bytes.Buffer
	command.Stderr = &failures
	if err := command.Run(); err != nil {
		t.Fatalf("make install into %s: %v: %s", directory, err, failures.String())
	}
	return directory
}

func invoked(t *testing.T, binary string, arguments ...string) (string, string, int) {
	t.Helper()

	command := exec.Command(binary, arguments...)
	command.Dir = t.TempDir()

	var out, failures bytes.Buffer
	command.Stdout = &out
	command.Stderr = &failures

	code := 0
	if err := command.Run(); err != nil {
		var exited *exec.ExitError
		if !errors.As(err, &exited) {
			t.Fatalf("running %s: %v", binary, err)
		}
		code = exited.ExitCode()
	}
	return out.String(), failures.String(), code
}
