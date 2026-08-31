package snapshot_test

import (
	"bytes"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-devtool/socrates/internal/testing/gitrepo"
)

func TestGitStatusIsUnchanged(t *testing.T) {
	repository := gitrepo.New(t)
	repository.Write("tracked.txt", "committed\n")
	repository.Commit("a tracked file")

	repository.Write("tracked.txt", "modified\n")
	repository.Write("untracked.txt", "new\n")
	repository.Write("staged.txt", "staged\n")
	repository.Git("add", "staged.txt")

	before := repository.Status()
	take(t, repository)
	after := repository.Status()

	if diff := cmp.Diff(before, after); diff != "" {
		t.Errorf("git status changed (-before +after):\n%s", diff)
	}
}

func TestTheIndexIsUnchanged(t *testing.T) {
	repository := gitrepo.New(t)
	repository.Write("partial.txt", "committed\n")
	repository.Commit("a file to stage half of")

	repository.Write("partial.txt", "staged\n")
	repository.Git("add", "partial.txt")
	repository.Write("partial.txt", "unstaged\n")

	index := repository.Bytes(".git/index")
	staged := repository.Git("diff", "--cached")
	unstaged := repository.Git("diff")

	take(t, repository)

	if !bytes.Equal(index, repository.Bytes(".git/index")) {
		t.Error("the person's index file is not byte-identical after the snapshot")
	}
	if diff := cmp.Diff(staged, repository.Git("diff", "--cached")); diff != "" {
		t.Errorf("the staged half changed (-before +after):\n%s", diff)
	}
	if diff := cmp.Diff(unstaged, repository.Git("diff")); diff != "" {
		t.Errorf("the unstaged half changed (-before +after):\n%s", diff)
	}
}

func TestNoRefNamesTheSnapshot(t *testing.T) {
	repository := gitrepo.New(t)
	repository.Write("changed.txt", "something to review\n")

	reflog := repository.Bytes(".git/logs/HEAD")
	logsBefore := administrativeFiles(repository, ".git/logs")

	taken := take(t, repository)

	named := repository.Git("for-each-ref", "--format=%(objectname) %(refname)")
	for line := range strings.SplitSeq(named, "\n") {
		if strings.HasPrefix(line, taken.Commit) {
			t.Errorf("%s names the snapshot commit", line)
		}
	}

	reachable := repository.Git("rev-list", "--all")
	if strings.Contains(reachable, taken.Commit) {
		t.Errorf("the snapshot commit %s is reachable from a ref", taken.Commit)
	}

	if !bytes.Equal(reflog, repository.Bytes(".git/logs/HEAD")) {
		t.Error("the reflog gained an entry")
	}
	if diff := cmp.Diff(logsBefore, administrativeFiles(repository, ".git/logs")); diff != "" {
		t.Errorf("the reflog files changed (-before +after):\n%s", diff)
	}
}

func TestTheSnapshotsParentIsTheCurrentHead(t *testing.T) {
	repository := gitrepo.New(t)
	repository.Write("changed.txt", "something to review\n")
	head := repository.Head()

	taken := take(t, repository)

	described := strings.Fields(repository.Git("rev-list", "--parents", "-n", "1", taken.Commit))
	if len(described) != 2 {
		t.Fatalf("the snapshot commit has the parents %v, want exactly one", described[1:])
	}
	if described[1] != head {
		t.Errorf("the snapshot's parent is %s, want %s", described[1], head)
	}
	if taken.Base != head {
		t.Errorf("the snapshot records %s as its base, want %s", taken.Base, head)
	}
}

func administrativeFiles(repository *gitrepo.Repository, prefix string) []string {
	within := strings.TrimPrefix(prefix, ".git/")

	var found []string
	for _, path := range repository.AdministrativeFiles() {
		if strings.HasPrefix(path, within) {
			found = append(found, path)
		}
	}
	slices.Sort(found)
	return found
}
