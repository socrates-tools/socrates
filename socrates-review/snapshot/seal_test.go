package snapshot_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/refusal"
	"github.com/socrates-tools/socrates/internal/testing/gitrepo"
	"github.com/socrates-tools/socrates/socrates-review/snapshot"
)

func TestAReviewerCannotCreateAFile(t *testing.T) {
	repository := gitrepo.New(t)
	repository.Write("internal/reviewed.go", "package internal\n")

	worktree := checkout(t, take(t, repository), t.TempDir())

	for _, path := range []string{"reviewer-notes.md", filepath.Join("internal", "reviewer-notes.md")} {
		created := filepath.Join(worktree.Path, path)
		if err := os.WriteFile(created, []byte("some thoughts\n"), 0o644); !errors.Is(err, fs.ErrPermission) {
			t.Errorf("creating %s returned %v, want a permission error", created, err)
		}
		if _, err := os.Stat(created); err == nil {
			t.Errorf("%s was created in the sealed worktree", created)
		}
	}

	if made := os.Mkdir(filepath.Join(worktree.Path, "notes"), 0o755); !errors.Is(made, fs.ErrPermission) {
		t.Errorf("creating a directory in the sealed worktree returned %v, want a permission error", made)
	}

	if dirty := gitrepo.At(t, worktree.Path).Git("status", "--porcelain=v2", "--untracked-files=all"); dirty != "" {
		t.Errorf("the worktree no longer matches the snapshot commit:\n%s", dirty)
	}
}

func TestAReviewerCannotModifyAFile(t *testing.T) {
	const reviewed = "package internal\n\nfunc Answer() int {\n\treturn 1\n}\n"

	repository := gitrepo.New(t)
	repository.Write("internal/reviewed.go", reviewed)

	worktree := checkout(t, take(t, repository), t.TempDir())
	held := filepath.Join(worktree.Path, "internal", "reviewed.go")

	if err := os.WriteFile(held, []byte("package internal\n"), 0o644); !errors.Is(err, fs.ErrPermission) {
		t.Errorf("overwriting %s returned %v, want a permission error", held, err)
	}

	appended, err := os.OpenFile(held, os.O_APPEND|os.O_WRONLY, 0o644)
	if !errors.Is(err, fs.ErrPermission) {
		t.Errorf("opening %s to append returned %v, want a permission error", held, err)
		if err == nil {
			_ = appended.Close()
		}
	}

	if removed := os.Remove(held); !errors.Is(removed, fs.ErrPermission) {
		t.Errorf("removing %s returned %v, want a permission error", held, removed)
	}

	if content := gitrepo.At(t, worktree.Path).Read("internal/reviewed.go"); content != reviewed {
		t.Errorf("the file holds %q, so the snapshot's content did not survive", content)
	}
}

func TestTheReviewerCanStillReadEverything(t *testing.T) {
	repository := gitrepo.New(t)
	repository.Write("internal/reviewed.go", "package internal\n\nfunc Answer() int {\n\treturn 1\n}\n")
	repository.Write("docs/nested/deeply/notes.md", "what the change is for\n")
	repository.Write("script.sh", "#!/bin/sh\necho ran\n")
	if err := os.Chmod(filepath.Join(repository.Root, "script.sh"), 0o755); err != nil {
		t.Fatalf("making the script executable: %v", err)
	}

	worktree := checkout(t, take(t, repository), t.TempDir())
	reviewing := gitrepo.At(t, worktree.Path)

	var read []string
	walked := filepath.WalkDir(worktree.Path, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || !entry.Type().IsRegular() {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(worktree.Path, path)
		if err != nil {
			return err
		}
		read = append(read, relative+" "+string(content))
		return nil
	})
	if walked != nil {
		t.Fatalf("listing and reading the sealed worktree: %v", walked)
	}

	for _, wanted := range []string{"internal/reviewed.go", "docs/nested/deeply/notes.md", "script.sh"} {
		found := false
		for _, line := range read {
			if strings.HasPrefix(line, wanted+" ") {
				found = true
			}
		}
		if !found {
			t.Errorf("%s was not readable in the sealed worktree, read %v", wanted, read)
		}
	}

	if content := reviewing.Read("docs/nested/deeply/notes.md"); content != "what the change is for\n" {
		t.Errorf("a nested file reads as %q", content)
	}
	if grepped := reviewing.Git("grep", "--name-only", "Answer"); grepped != "internal/reviewed.go" {
		t.Errorf("git grep in the sealed worktree found %q", grepped)
	}
	if dirty := reviewing.Git("status", "--porcelain=v2", "--untracked-files=all"); dirty != "" {
		t.Errorf("git status in the sealed worktree reports:\n%s", dirty)
	}

	executable, err := os.Stat(filepath.Join(worktree.Path, "script.sh"))
	if err != nil {
		t.Fatalf("reading the mode of the sealed script: %v", err)
	}
	if executable.Mode().Perm()&0o111 == 0 {
		t.Errorf("the sealed script reads as %s, so sealing took the executable bit with it", executable.Mode())
	}
}

func TestASealThatDoesNotHoldIsARefusal(t *testing.T) {
	repository := gitrepo.New(t)
	repository.Write("internal/reviewed.go", "package internal\n")

	worktree := checkout(t, take(t, repository), t.TempDir())

	if err := worktree.Unseal(); err != nil {
		t.Fatalf("standing in for a filesystem that takes the seal and enforces nothing: %v", err)
	}

	err := worktree.ProveTheSeal()
	if err == nil {
		t.Fatal("a worktree that can be written to proved the seal")
	}
	if !refusal.Is(err) {
		t.Errorf("the probe returned %v, which is not a refusal, so it does not exit with the usage code", err)
	}
	if !snapshot.SealNotEnforced(err) {
		t.Errorf("the probe returned %v, which no caller can recognise as an unenforced seal", err)
	}
	if !strings.Contains(err.Error(), worktree.Path) {
		t.Errorf("the refusal reads %q, without naming the worktree", err)
	}
	if !strings.Contains(err.Error(), "filesystem") {
		t.Errorf("the refusal reads %q, without naming the filesystem as the reason", err)
	}

	if dirty := gitrepo.At(t, worktree.Path).Git("status", "--porcelain=v2", "--untracked-files=all"); dirty != "" {
		t.Errorf("the probe left something behind:\n%s", dirty)
	}
}

func TestASealedWorktreeProvesTheSeal(t *testing.T) {
	repository := gitrepo.New(t)
	repository.Write("internal/reviewed.go", "package internal\n")

	worktree := checkout(t, take(t, repository), t.TempDir())

	if err := worktree.ProveTheSeal(); err != nil {
		t.Fatalf("a worktree the checkout sealed did not prove the seal: %v", err)
	}
}
