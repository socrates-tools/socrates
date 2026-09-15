package snapshot_test

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/socrates-tools/socrates/internal/testing/gitrepo"
	"github.com/socrates-tools/socrates/socrates-review/runstate"
	"github.com/socrates-tools/socrates/socrates-review/runstore"
	"github.com/socrates-tools/socrates/socrates-review/snapshot"
)

func TestAnUntrackedFileIsInTheSnapshot(t *testing.T) {
	repository := gitrepo.New(t)
	repository.Write("internal/added.go", "package internal\n")

	taken := take(t, repository)

	if !slices.Contains(repository.Tracked(taken.Commit), "internal/added.go") {
		t.Fatalf("the snapshot holds %v, without the untracked file", repository.Tracked(taken.Commit))
	}
	if content := repository.Contents(taken.Commit, "internal/added.go"); content != "package internal" {
		t.Errorf("the snapshot holds %q for the untracked file", content)
	}
}

func TestAModifiedTrackedFileIsInTheSnapshot(t *testing.T) {
	repository := gitrepo.New(t)
	repository.Write("tracked.txt", "committed\n")
	repository.Commit("a tracked file")
	repository.Write("tracked.txt", "modified\n")

	taken := take(t, repository)

	if content := repository.Contents(taken.Commit, "tracked.txt"); content != "modified" {
		t.Errorf("the snapshot holds %q, want the modified contents", content)
	}
	if content := repository.Contents(taken.Base, "tracked.txt"); content != "committed" {
		t.Errorf("the base holds %q, so the fixture is not what this test needs", content)
	}
}

func TestAStagedChangeAndAnUnstagedChangeAreBothInTheSnapshot(t *testing.T) {
	repository := gitrepo.New(t)
	repository.Write("tracked.txt", "committed\n")
	repository.Commit("a tracked file")

	repository.Write("staged.txt", "staged\n")
	repository.Git("add", "staged.txt")
	repository.Write("tracked.txt", "unstaged\n")

	taken := take(t, repository)

	if content := repository.Contents(taken.Commit, "staged.txt"); content != "staged" {
		t.Errorf("the snapshot holds %q for the staged file", content)
	}
	if content := repository.Contents(taken.Commit, "tracked.txt"); content != "unstaged" {
		t.Errorf("the snapshot holds %q for the unstaged change", content)
	}
}

func TestADeletedFileIsDeletedInTheSnapshot(t *testing.T) {
	repository := gitrepo.New(t)
	repository.Write("doomed.txt", "here for now\n")
	repository.Commit("a file to delete")
	repository.Delete("doomed.txt")

	taken := take(t, repository)

	if slices.Contains(repository.Tracked(taken.Commit), "doomed.txt") {
		t.Errorf("the snapshot holds %v, still naming the deleted file", repository.Tracked(taken.Commit))
	}
}

func TestAnIgnoredFileIsNotInTheSnapshot(t *testing.T) {
	repository := gitrepo.New(t)
	repository.Write(".gitignore", "ignored.txt\nbuilt/\n")
	repository.Write("ignored.txt", "noise\n")
	repository.Write("built/artifact", "noise\n")
	repository.Write("kept.txt", "signal\n")

	taken := take(t, repository)
	tracked := repository.Tracked(taken.Commit)

	for _, ignored := range []string{"ignored.txt", "built/artifact"} {
		if slices.Contains(tracked, ignored) {
			t.Errorf("the snapshot holds %s, which .gitignore excludes", ignored)
		}
	}
	if !slices.Contains(tracked, "kept.txt") {
		t.Errorf("the snapshot holds %v, without the file nothing ignores", tracked)
	}
}

func TestAChangeWithNothingToReview(t *testing.T) {
	repository := gitrepo.New(t)

	taken, err := snapshot.Take(context.Background(), repository.Root, "")
	if err != nil {
		t.Fatalf("taking the snapshot: %v", err)
	}

	if !taken.NothingToReview() {
		t.Errorf("a working tree identical to HEAD reports %v as changed", taken.Changed)
	}
	if taken.Commit != "" {
		t.Errorf("a working tree identical to HEAD produced the snapshot commit %s", taken.Commit)
	}
	if taken.Base != repository.Head() {
		t.Errorf("the snapshot's base is %s, want %s", taken.Base, repository.Head())
	}

	if _, err := taken.Checkout(context.Background(), filepath.Join(t.TempDir(), "worktree")); err == nil {
		t.Error("a snapshot with nothing to review was checked out, so a reviewer could be spawned on it")
	}

	store := runstore.New(t.TempDir())
	run, err := store.Create(runstore.Start{
		StartedAt:      time.Date(2026, 8, 19, 4, 15, 30, 0, time.UTC),
		RepositoryRoot: repository.Root,
		Base:           taken.Base,
	})
	if err != nil {
		t.Fatalf("creating the run: %v", err)
	}
	if err := run.End(runstate.Failed, "the working tree held nothing the base commit does not"); err != nil {
		t.Fatalf("ending the run: %v", err)
	}
	if err := store.Save(run); err != nil {
		t.Fatalf("saving the run: %v", err)
	}

	loaded, err := store.Load(run.Name)
	if err != nil {
		t.Fatalf("reading the run: %v", err)
	}
	if loaded.State != runstate.Failed {
		t.Errorf("a run with nothing to review reads as %q, want %q", loaded.State, runstate.Failed)
	}
	if loaded.Verdict == nil || !strings.Contains(loaded.Verdict.Reason, "nothing the base commit does not") {
		t.Errorf("the verdict is %+v, without recording that there was nothing to review", loaded.Verdict)
	}
	if loaded.Snapshot != "" {
		t.Errorf("a run with nothing to review names %s as its snapshot", loaded.Snapshot)
	}
}

func TestTheChangedPathsAreWhatChanged(t *testing.T) {
	repository := gitrepo.New(t)
	repository.Write("added.txt", "new\n")
	repository.Write("README.md", "changed\n")

	taken := take(t, repository)

	want := []string{"README.md", "added.txt"}
	got := slices.Clone(taken.Changed)
	slices.Sort(got)

	if !slices.Equal(want, got) {
		t.Errorf("the snapshot reports %v as changed, want %v", got, want)
	}
}

func TestGitStashCreateIsNotUsed(t *testing.T) {
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("reading the package directory: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		content, err := os.ReadFile(entry.Name())
		if err != nil {
			t.Fatalf("reading %s: %v", entry.Name(), err)
		}
		if strings.Contains(string(content), "stash") {
			t.Errorf("%s reaches for git stash, which silently ignores untracked files", entry.Name())
		}
	}
}

func take(t *testing.T, repository *gitrepo.Repository) snapshot.Snapshot {
	t.Helper()

	taken, err := snapshot.Take(context.Background(), repository.Root, "")
	if err != nil {
		t.Fatalf("taking the snapshot: %v", err)
	}
	if taken.NothingToReview() {
		t.Fatal("the snapshot found nothing to review, so this test asserts nothing")
	}
	return taken
}
