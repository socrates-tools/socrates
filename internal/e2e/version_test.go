package e2e_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/socrates-devtool/socrates/internal/testing/harness"
	"github.com/socrates-devtool/socrates/internal/testing/source"
)

const (
	saysUncommitted = "uncommitted changes"
	aReleaseTag     = "v9.9.9"
)

var aSemanticVersion = regexp.MustCompile(`v[0-9]+\.[0-9]+\.[0-9]+`)

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

func TestATaggedBuildReportsItsVersion(t *testing.T) {
	root := requireMake(t)
	tagged := aCopyOfTheSource(t, root)
	gitSaid(t, tagged, "tag", aReleaseTag)

	printed := versionOf(t, made(t, tagged))
	if !strings.Contains(printed, aReleaseTag) {
		t.Errorf("a binary built from a commit tagged %s reports %q, which names no release", aReleaseTag, printed)
	}
	if commit := gitSaid(t, tagged, "rev-parse", "--short=7", "HEAD"); !strings.Contains(printed, commit) {
		t.Errorf("a binary built from a tagged commit stopped naming the commit %s: %q", commit, printed)
	}
	if strings.Contains(printed, saysUncommitted) {
		t.Errorf("a binary built from a clean tagged tree says the tree held uncommitted changes: %q", printed)
	}
}

func TestAnUntaggedBuildDoesNotClaimAVersion(t *testing.T) {
	root := requireMake(t)
	untagged := aCopyOfTheSource(t, root)
	gitSaid(t, untagged, "tag", aReleaseTag)
	gitSaid(t, untagged, "-c", "user.name=A Person", "-c", "user.email=person@example.com", "commit", "--quiet", "--allow-empty", "-m", "a commit past the tag")

	printed := versionOf(t, made(t, untagged))
	if !strings.Contains(printed, "development build") {
		t.Errorf("a binary built one commit past %s reports %q, which does not say it is a development build", aReleaseTag, printed)
	}
	if strings.Contains(printed, aReleaseTag) {
		t.Errorf("a binary built one commit past %s borrowed that tag: %q", aReleaseTag, printed)
	}
	if named := aSemanticVersion.FindString(printed); named != "" {
		t.Errorf("a binary built from no tag names %s as its version: %q", named, printed)
	}
	if commit := gitSaid(t, untagged, "rev-parse", "--short=7", "HEAD"); !strings.Contains(printed, commit) {
		t.Errorf("a binary built from no tag stopped naming the commit %s: %q", commit, printed)
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

	var out, failures bytes.Buffer
	command.Stdout = &out
	command.Stderr = &failures
	if err := command.Run(); err != nil {
		t.Fatalf("git %s in %s: %v: %s", strings.Join(arguments, " "), root, err, failures.String())
	}
	return strings.TrimSpace(out.String())
}

func aCopyOfTheSource(t *testing.T, root string) string {
	t.Helper()

	directory := t.TempDir()
	listed := gitSaid(t, root, "ls-files", "-z", "--cached", "--others", "--exclude-standard")

	copied := 0
	for relative := range strings.SplitSeq(listed, "\x00") {
		if relative == "" {
			continue
		}

		from := filepath.Join(root, filepath.FromSlash(relative))
		held, err := os.ReadFile(from)
		if err != nil {
			continue
		}
		mode, err := os.Stat(from)
		if err != nil {
			t.Fatalf("reading the mode of %s: %v", relative, err)
		}

		to := filepath.Join(directory, filepath.FromSlash(relative))
		if err := os.MkdirAll(filepath.Dir(to), 0o755); err != nil {
			t.Fatalf("making room for %s: %v", relative, err)
		}
		if err := os.WriteFile(to, held, mode.Mode().Perm()); err != nil {
			t.Fatalf("writing %s: %v", relative, err)
		}
		copied++
	}
	if copied == 0 {
		t.Fatalf("nothing was copied out of %s, so what is built there is not this source", root)
	}

	gitSaid(t, directory, "-c", "init.defaultBranch=main", "init", "--quiet")
	gitSaid(t, directory, "add", "--all")
	gitSaid(t, directory, "-c", "user.name=A Person", "-c", "user.email=person@example.com", "commit", "--quiet", "-m", "the committed tree")
	return directory
}
