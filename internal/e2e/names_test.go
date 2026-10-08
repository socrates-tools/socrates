package e2e_test

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/convention"
	"github.com/socrates-tools/socrates/internal/exit"
	skillkit "github.com/socrates-tools/socrates/internal/skill"
	"github.com/socrates-tools/socrates/internal/testing/harness"
	"github.com/socrates-tools/socrates/internal/testing/source"
	"github.com/socrates-tools/socrates/review/skill"
)

func TestOneNameIsInstalledAndNothingIsPutBesideIt(t *testing.T) {
	t.Parallel()
	root := requireMake(t)
	directory, home := installed(t, root)

	tool := filepath.Join(directory, convention.Suite)
	reported := versionOf(t, tool)
	if commit := gitSaid(t, root, "rev-parse", "--short=7", "HEAD"); !strings.Contains(reported, commit) {
		t.Errorf("the installed binary reports %q, which does not name the commit %s", reported, commit)
	}

	if _, err := os.Stat(tool); err != nil {
		t.Fatalf("the tool was not installed under %s: %v", convention.Suite, err)
	}

	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatalf("reading what installing put on the path: %v", err)
	}
	for _, entry := range entries {
		if entry.Name() == convention.Suite {
			continue
		}
		if info, err := entry.Info(); err == nil && info.Mode().Perm()&0o111 != 0 && !info.IsDir() {
			t.Errorf("installing put %s beside %s, and the suite publishes one binary", entry.Name(), convention.Suite)
		}
	}

	// scenario: make install also installs the skill, under the home it was given rather than the person's own
	for _, one := range skill.All() {
		path := filepath.Join(skillkit.UserDirectoryUnder(home), one.Name, skillkit.FileName)
		content, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("make install did not install %s under the disposable home: %v", one.Name, err)
			continue
		}
		if string(content) != one.Document() {
			t.Errorf("what make install wrote for %s under %s is not what the skill's source generates", one.Name, path)
		}
	}
}

func TestWhatItPrintsNamesTheTool(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)

	printed, failures, code := invoked(t, session.Binary, "kaboom")
	if code != exit.Unusable {
		t.Fatalf("an unknown command was refused with %d, want %d: %s", code, exit.Unusable, failures)
	}
	if !strings.Contains(failures, convention.Suite) {
		t.Errorf("what it printed does not name %s: %q", convention.Suite, failures)
	}

	helped, helpFailures, code := invoked(t, session.Binary, "--help")
	if code != 0 {
		t.Fatalf("--help exited %d: %s", code, helpFailures)
	}
	if !strings.Contains(helped, convention.Suite) {
		t.Errorf("what --help printed does not name %s: %q", convention.Suite, helped)
	}

	for _, said := range strings.Fields(printed + failures + helped + helpFailures) {
		if said == convention.Repository || said == "scrat" {
			t.Errorf("what it printed offers %q, and that name was retired: run %s instead", said, convention.Suite)
		}
	}
}

func installed(t *testing.T, root string) (string, string) {
	t.Helper()

	directory := t.TempDir()
	home := t.TempDir()
	command := exec.Command("make", "install")
	command.Dir = root
	command.Env = append(append(os.Environ(), source.Toolchain()...), "GOBIN="+directory, "HOME="+home)

	var failures bytes.Buffer
	command.Stderr = &failures
	if err := command.Run(); err != nil {
		t.Fatalf("make install into %s: %v: %s", directory, err, failures.String())
	}
	return directory, home
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
