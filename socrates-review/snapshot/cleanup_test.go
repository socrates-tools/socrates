package snapshot_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/internal/testing/gitrepo"
)

func TestWorktreeCreationFails(t *testing.T) {
	repository := gitrepo.New(t)
	repository.Write("reviewed.txt", "something to review\n")

	occupied := filepath.Join(t.TempDir(), "worktree")
	if err := os.MkdirAll(occupied, 0o755); err != nil {
		t.Fatalf("creating the occupied directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(occupied, "in the way"), []byte("mine\n"), 0o644); err != nil {
		t.Fatalf("filling the occupied directory: %v", err)
	}

	administrativeBefore := repository.AdministrativeFiles()
	taken := take(t, repository)

	worktree, err := taken.Checkout(context.Background(), occupied)
	if err == nil {
		t.Fatalf("checking out over an occupied directory succeeded, returning %+v", worktree)
	}

	t.Run("A snapshot in progress leaves nothing behind when it fails", func(t *testing.T) {
		if strings.Contains(repository.Worktrees(), occupied) {
			t.Errorf("git worktree list names %s:\n%s", occupied, repository.Worktrees())
		}
		if leaked := leakedIndexes(t); len(leaked) != 0 {
			t.Errorf("the snapshot's own index is still on disk: %v", leaked)
		}
		if content, readErr := os.ReadFile(filepath.Join(occupied, "in the way")); readErr != nil || string(content) != "mine\n" {
			t.Errorf("the directory that was in the way was disturbed: %v", readErr)
		}

		added := addedAdministrativeFiles(administrativeBefore, repository.AdministrativeFiles())
		for _, path := range added {
			if !strings.HasPrefix(path, "objects/") {
				t.Errorf(".git/%s was left behind by a snapshot that failed", path)
			}
		}
	})
}

func TestTheTemporaryIndexNeverJoinsTheRepository(t *testing.T) {
	repository := gitrepo.New(t)
	repository.Write("reviewed.txt", "something to review\n")

	before := repository.AdministrativeFiles()

	taken := take(t, repository)
	worktree := checkout(t, taken, t.TempDir())

	for _, path := range addedAdministrativeFiles(before, repository.AdministrativeFiles()) {
		if strings.HasPrefix(path, "objects/") || strings.HasPrefix(path, "worktrees/") {
			continue
		}
		t.Errorf(".git/%s was written by taking a snapshot", path)
	}

	if leaked := leakedIndexes(t); len(leaked) != 0 {
		t.Errorf("the snapshot's own index is still on disk: %v", leaked)
	}

	if err := worktree.Remove(context.Background()); err != nil {
		t.Fatalf("removing the worktree: %v", err)
	}

	after := repository.AdministrativeFiles()
	for _, path := range addedAdministrativeFiles(before, after) {
		if strings.HasPrefix(path, "objects/") {
			continue
		}
		t.Errorf(".git/%s remains after the run", path)
	}
	if diff := cmp.Diff(removedAdministrativeFiles(before, after), []string(nil)); diff != "" {
		t.Errorf("taking a snapshot removed administrative files (-removed +none):\n%s", diff)
	}
}

func leakedIndexes(t *testing.T) []string {
	t.Helper()

	leaked, err := filepath.Glob(filepath.Join(os.TempDir(), "socrates-index-*"))
	if err != nil {
		t.Fatalf("looking for the snapshot's own index: %v", err)
	}
	return leaked
}

func addedAdministrativeFiles(before, after []string) []string {
	held := make(map[string]bool, len(before))
	for _, path := range before {
		held[path] = true
	}

	var added []string
	for _, path := range after {
		if !held[path] {
			added = append(added, path)
		}
	}
	return added
}

func removedAdministrativeFiles(before, after []string) []string {
	return addedAdministrativeFiles(after, before)
}
