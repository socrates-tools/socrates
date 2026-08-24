package e2e_test

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rui-ktei/socrates/internal/convention"
	"github.com/rui-ktei/socrates/internal/testing/harness"
	"github.com/rui-ktei/socrates/internal/testing/source"
)

func TestBothNamesAreInstalledAndReachTheSameBinary(t *testing.T) {
	root := requireMake(t)
	directory := installed(t, root)

	full := versionOf(t, filepath.Join(directory, "socrates"))
	short := versionOf(t, filepath.Join(directory, convention.ShortName))

	if full != short {
		t.Errorf("the two names report different builds: %q and %q", full, short)
	}
	if commit := gitSaid(t, root, "rev-parse", "--short=7", "HEAD"); !strings.Contains(full, commit) {
		t.Errorf("the installed binary reports %q, which does not name the commit %s", full, commit)
	}
}

func TestTheShortNameIsASecondNameForOneBinary(t *testing.T) {
	directory := installed(t, requireMake(t))

	entry, err := os.Lstat(filepath.Join(directory, convention.ShortName))
	if err != nil {
		t.Fatalf("the short name was not installed: %v", err)
	}
	if entry.Mode()&os.ModeSymlink == 0 {
		t.Errorf("the short name is a %s of its own rather than a second name for the one binary", entry.Mode())
	}

	full, err := os.Stat(filepath.Join(directory, "socrates"))
	if err != nil {
		t.Fatalf("the full name was not installed: %v", err)
	}
	reached, err := os.Stat(filepath.Join(directory, convention.ShortName))
	if err != nil {
		t.Fatalf("the short name reaches nothing: %v", err)
	}
	if !os.SameFile(full, reached) {
		t.Errorf("the short name reaches something other than the binary the full name is")
	}
}

func TestTheShortNameReportsTheFullOne(t *testing.T) {
	session := harness.Start(t)

	directory := t.TempDir()
	shortName := filepath.Join(directory, convention.ShortName)
	if err := os.Symlink(session.Binary, shortName); err != nil {
		t.Fatalf("naming the binary %s: %v", convention.ShortName, err)
	}

	printed, failures, code := invoked(t, shortName, "kaboom")
	if code != 3 {
		t.Fatalf("the short name refused an unknown command with %d, want 3: %s", code, failures)
	}
	if !strings.Contains(failures, "socrates") {
		t.Errorf("what the short name printed does not name socrates: %q", failures)
	}
	if strings.Contains(failures+printed, convention.ShortName) {
		t.Errorf("what the short name printed names the short name back: %q", failures+printed)
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
