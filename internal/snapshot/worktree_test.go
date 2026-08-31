package snapshot_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-devtool/socrates/internal/snapshot"
	"github.com/socrates-devtool/socrates/internal/testing/gitrepo"
)

func TestTheWorktreeIsTheSnapshot(t *testing.T) {
	repository := gitrepo.New(t)
	repository.Write("reviewed.txt", "what the reviewer should see\n")

	taken := take(t, repository)
	worktree := checkout(t, taken, t.TempDir())

	checkedOut := gitrepo.At(t, worktree.Path)
	if head := checkedOut.Head(); head != taken.Commit {
		t.Errorf("the worktree is at %s, want the snapshot commit %s", head, taken.Commit)
	}
	if branch := checkedOut.Git("rev-parse", "--abbrev-ref", "HEAD"); branch != "HEAD" {
		t.Errorf("the worktree is on the branch %s, want a detached head", branch)
	}
	if dirty := checkedOut.Git("status", "--porcelain=v2", "--untracked-files=all"); dirty != "" {
		t.Errorf("the worktree does not match the snapshot tree:\n%s", dirty)
	}
	if content := checkedOut.Read("reviewed.txt"); content != "what the reviewer should see\n" {
		t.Errorf("the worktree holds %q for the reviewed file", content)
	}
}

func TestTheWorktreeGoesUnderTheRunDirectory(t *testing.T) {
	runDirectory := filepath.Join(t.TempDir(), "20260819T041530Z-ee66fe34")
	path := snapshot.WorktreePath(runDirectory)

	if filepath.Dir(path) != runDirectory {
		t.Errorf("the worktree path %s is not under the run directory %s", path, runDirectory)
	}
}

func TestEditingWhileTheReviewerReads(t *testing.T) {
	repository := gitrepo.New(t)
	repository.Write("reviewed.txt", "as it was when the review started\n")

	taken := take(t, repository)
	worktree := checkout(t, taken, t.TempDir())

	repository.Write("reviewed.txt", "as the person has since changed it\n")
	repository.Write("later.txt", "written after the reviewer started\n")

	checkedOut := gitrepo.At(t, worktree.Path)
	if content := checkedOut.Read("reviewed.txt"); content != "as it was when the review started\n" {
		t.Errorf("the reviewer's worktree holds %q, so it moved under the reviewer", content)
	}
	if _, err := os.Stat(filepath.Join(worktree.Path, "later.txt")); err == nil {
		t.Error("a file the person wrote after the snapshot appeared in the reviewer's worktree")
	}
	if dirty := checkedOut.Git("status", "--porcelain=v2", "--untracked-files=all"); dirty != "" {
		t.Errorf("the reviewer's worktree is no longer the snapshot:\n%s", dirty)
	}
}

func TestTheWorktreeIsCleanedUp(t *testing.T) {
	repository := gitrepo.New(t)
	repository.Write("reviewed.txt", "something to review\n")

	before := repository.Status()
	taken := take(t, repository)
	worktree := checkout(t, taken, t.TempDir())

	if !strings.Contains(repository.Worktrees(), worktree.Path) {
		t.Fatalf("git worktree list does not name %s, so this test asserts nothing", worktree.Path)
	}

	if err := worktree.Remove(context.Background()); err != nil {
		t.Fatalf("removing the worktree: %v", err)
	}

	if strings.Contains(repository.Worktrees(), worktree.Path) {
		t.Errorf("git worktree list still names %s:\n%s", worktree.Path, repository.Worktrees())
	}
	if _, err := os.Stat(worktree.Path); err == nil {
		t.Errorf("%s is still on disk", worktree.Path)
	}
	if after := repository.Status(); after != before {
		t.Errorf("the person's repository reads differently after the worktree came and went:\n%s", after)
	}
}

func checkout(t *testing.T, taken snapshot.Snapshot, into string) snapshot.Worktree {
	t.Helper()

	worktree, err := taken.Checkout(context.Background(), snapshot.WorktreePath(into))
	if err != nil {
		t.Fatalf("checking out the snapshot: %v", err)
	}
	t.Cleanup(func() {
		_ = worktree.Remove(context.Background())
	})
	return worktree
}
