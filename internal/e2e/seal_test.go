package e2e_test

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/rui-ktei/socrates/internal/exit"
	"github.com/rui-ktei/socrates/internal/state"
	"github.com/rui-ktei/socrates/internal/testing/harness"
)

const unenforcingDirectory = "SOCRATES_TEST_UNENFORCING_DIRECTORY"

func TestASealThatDoesNotHoldIsARefusal(t *testing.T) {
	stated := os.Getenv(unenforcingDirectory)
	if stated == "" {
		t.Skipf("set %s to a directory on a filesystem that ignores permission bits to run this", unenforcingDirectory)
	}
	if enforced(t, stated) {
		t.Fatalf("%s names %s, whose filesystem does enforce permissions, so this test would assert nothing", unenforcingDirectory, stated)
	}

	session := harness.Start(t)
	session.Repository.Write("internal/added.go", "package internal\n")

	home := filepath.Join(stated, "socrates-e2e-"+strconv.Itoa(os.Getpid()))
	t.Cleanup(func() { _ = os.RemoveAll(home) })

	out, failures, code := session.ScratWith(session.Repository.Root, []string{state.HomeVariable + "=" + home}, "review")

	if code != exit.Usage {
		t.Fatalf("scrat review exited %d on a filesystem that enforces nothing, want the usage code %d: %s%s", code, exit.Usage, out, failures)
	}
	if !strings.Contains(failures, "worktree") || !strings.Contains(failures, "filesystem") {
		t.Errorf("the refusal reads %q, without naming the worktree and the filesystem", failures)
	}
	if strings.TrimSpace(out) != "" {
		t.Errorf("scrat review printed %q, so a person is following a run that never started", out)
	}

	runs, err := os.ReadDir(filepath.Join(home, "runs"))
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("reading what the refused review left behind: %v", err)
	}
	if len(runs) != 0 {
		t.Errorf("the refused review left %d run directories behind", len(runs))
	}

	if listed := session.Repository.Worktrees(); strings.Contains(listed, home) {
		t.Errorf("git worktree list still names a worktree under %s:\n%s", home, listed)
	}
	if dirty := session.Repository.Git("status", "--porcelain=v1", "--untracked-files=all"); dirty != "?? internal/added.go" {
		t.Errorf("the person's repository reads:\n%s", dirty)
	}
}

func enforced(t *testing.T, directory string) bool {
	t.Helper()

	sealed, err := os.MkdirTemp(directory, "socrates-enforcement-")
	if err != nil {
		t.Fatalf("creating a directory under %s: %v", directory, err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(sealed) })

	if err := os.Chmod(sealed, 0o555); err != nil {
		t.Fatalf("clearing the write bits of %s: %v", sealed, err)
	}

	written, err := os.OpenFile(filepath.Join(sealed, "probe"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return true
	}
	written.Close()

	if err := os.Chmod(sealed, 0o755); err != nil {
		t.Fatalf("restoring the write bits of %s: %v", sealed, err)
	}
	return false
}
