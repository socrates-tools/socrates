package snapshot_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/socrates-tools/socrates/internal/testing/gitrepo"
	"github.com/socrates-tools/socrates/review/snapshot"
)

func TestSealingAWorktreeLeavesWhatASymlinkInItPointsAtWritable(t *testing.T) {
	somebodyElses := filepath.Join(t.TempDir(), "elsewhere.txt")
	if err := os.WriteFile(somebodyElses, []byte("not the change's\n"), 0o644); err != nil {
		t.Fatalf("writing a file outside the repository: %v", err)
	}

	repository := gitrepo.New(t)
	repository.Write("internal/reviewed.go", "package internal\n")
	if err := os.Symlink(somebodyElses, filepath.Join(repository.Root, "outside")); err != nil {
		t.Fatalf("committing a symlink that points outside the repository: %v", err)
	}

	worktree := checkout(t, take(t, repository), t.TempDir())
	if pointed, err := os.Readlink(filepath.Join(worktree.Path, "outside")); err != nil || pointed != somebodyElses {
		t.Fatalf("the sealed worktree holds no symlink to %s (read %q, %v), so this test is not reading the case it names", somebodyElses, pointed, err)
	}

	held, err := os.Stat(somebodyElses)
	if err != nil {
		t.Fatalf("reading the mode of what the symlink points at: %v", err)
	}
	if held.Mode().Perm() != 0o644 {
		t.Errorf("what the symlink points at is %s after sealing, so the seal reached outside the worktree through it", held.Mode().Perm())
	}
}

func TestSealingClearsTheWriteBitsAndKeepsEveryOther(t *testing.T) {
	for _, sealed := range []struct {
		named  string
		before fs.FileMode
		after  fs.FileMode
	}{
		{named: "a file only its owner may read stays unreadable to everybody else", before: 0o600, after: 0o400},
		{named: "a script keeps its executable bits", before: 0o755, after: 0o555},
		{named: "a file everybody may read stays readable", before: 0o644, after: 0o444},
	} {
		t.Run(sealed.named, func(t *testing.T) {
			tree := t.TempDir()
			held := filepath.Join(tree, "held")
			if err := os.WriteFile(held, []byte("content\n"), 0o600); err != nil {
				t.Fatalf("writing %s: %v", held, err)
			}
			if err := os.Chmod(held, sealed.before); err != nil {
				t.Fatalf("giving %s the mode %s: %v", held, sealed.before, err)
			}

			if err := snapshot.Seal(tree); err != nil {
				t.Fatalf("sealing %s: %v", tree, err)
			}
			t.Cleanup(func() { _ = snapshot.Unseal(tree) })

			read, err := os.Stat(held)
			if err != nil {
				t.Fatalf("reading the mode of %s: %v", held, err)
			}
			if read.Mode().Perm() != sealed.after {
				t.Errorf("a %s file seals to %s, want %s", sealed.before, read.Mode().Perm(), sealed.after)
			}
		})
	}
}
