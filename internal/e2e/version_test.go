package e2e_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-devtool/socrates/internal/testing/harness"
	"github.com/socrates-devtool/socrates/internal/testing/source"
)

const saysUncommitted = "uncommitted changes"

func TestTheBuildSaysWhichCommitItCameFrom(t *testing.T) {
	root := requireMake(t)

	commit := gitSaid(t, root, "rev-parse", "--short=7", "HEAD")
	uncommitted := gitSaid(t, root, "status", "--porcelain") != ""

	printed := versionOf(t, made(t, root))
	if !strings.Contains(printed, commit) {
		t.Errorf("socrates version says %q, which does not name the commit %s it was built from", printed, commit)
	}
	if said := strings.Contains(printed, saysUncommitted); said != uncommitted {
		t.Errorf("socrates version says the tree held uncommitted changes: %t, want %t: %q", said, uncommitted, printed)
	}
	if strings.Count(printed, "\n") != 0 {
		t.Errorf("socrates version printed more than one line:\n%s", printed)
	}
}

func TestABuildFromATreeWithUncommittedChangesSaysSo(t *testing.T) {
	root := requireMake(t)

	commit := gitSaid(t, root, "rev-parse", "--short=7", "HEAD")
	printed := versionOf(t, made(t, root, "UNCOMMITTED=yes"))

	if !strings.Contains(printed, saysUncommitted) {
		t.Errorf("a binary linked from a tree that held uncommitted changes says %q", printed)
	}
	if !strings.Contains(printed, commit) {
		t.Errorf("a binary linked from a tree that held uncommitted changes stopped naming the commit %s: %q", commit, printed)
	}
}

func TestABuildThatRecordedNoCommitSaysSo(t *testing.T) {
	root := requireMake(t)

	printed := versionOf(t, made(t, root, "COMMIT="))
	if !strings.Contains(printed, "no commit") {
		t.Errorf("a binary linked without the stamp says %q", printed)
	}
	if !strings.Contains(printed, "make install") {
		t.Errorf("a binary linked without the stamp does not say what to build it with: %q", printed)
	}
	if strings.Contains(printed, gitSaid(t, root, "rev-parse", "--short=7", "HEAD")) {
		t.Errorf("a binary linked without the stamp named a commit anyway: %q", printed)
	}
}

func TestScratVersionTakesNoArguments(t *testing.T) {
	session := harness.Start(t)

	printed, failures, code := session.Socrates("version", "the one I installed on Tuesday")
	if code != 3 {
		t.Errorf("socrates version with an argument exited %d, want 3: %s", code, failures)
	}
	if !strings.Contains(failures, "the one I installed on Tuesday") {
		t.Errorf("the refusal does not name the argument: %s", failures)
	}
	if strings.TrimSpace(printed) != "" {
		t.Errorf("socrates version printed a version for a usage error:\n%s", printed)
	}
}

func requireMake(t *testing.T) string {
	t.Helper()

	if _, err := exec.LookPath("make"); err != nil {
		t.Skip("make is not on PATH, and the stamp is the Makefile's")
	}
	return source.Root(t)
}

func made(t *testing.T, root string, overrides ...string) string {
	t.Helper()

	binary := filepath.Join(t.TempDir(), "socrates")
	arguments := append([]string{"build", "BINARY=" + binary}, overrides...)

	command := exec.Command("make", arguments...)
	command.Dir = root
	command.Env = append(os.Environ(), source.Toolchain()...)

	var failures bytes.Buffer
	command.Stderr = &failures
	if err := command.Run(); err != nil {
		t.Fatalf("make %s: %v: %s", strings.Join(arguments, " "), err, failures.String())
	}
	return binary
}

func versionOf(t *testing.T, binary string) string {
	t.Helper()

	command := exec.Command(binary, "version")
	command.Dir = t.TempDir()

	var out, failures bytes.Buffer
	command.Stdout = &out
	command.Stderr = &failures
	if err := command.Run(); err != nil {
		t.Fatalf("%s version: %v: %s", binary, err, failures.String())
	}
	return strings.TrimSpace(out.String())
}

func gitSaid(t *testing.T, root string, arguments ...string) string {
	t.Helper()

	command := exec.Command("git", arguments...)
	command.Dir = root

	var out bytes.Buffer
	command.Stdout = &out
	if err := command.Run(); err != nil {
		t.Fatalf("git %s in %s: %v", strings.Join(arguments, " "), root, err)
	}
	return strings.TrimSpace(out.String())
}
