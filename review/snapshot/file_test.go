package snapshot_test

import (
	"context"
	"testing"

	"github.com/socrates-tools/socrates/internal/testing/gitrepo"
	"github.com/socrates-tools/socrates/review/snapshot"
)

func TestAFileIsReadAtTheCommitRatherThanFromTheWorkingTree(t *testing.T) {
	repository := gitrepo.New(t)
	repository.Write("notes.md", "what the commit says\n")
	base := repository.Commit("a file to read")
	repository.Write("notes.md", "what the working tree says\n")

	content, present, err := snapshot.FileAt(context.Background(), repository.Root, base, "notes.md")
	if err != nil {
		t.Fatalf("reading notes.md at %s: %v", base, err)
	}
	if !present {
		t.Fatal("the commit carries notes.md and the read reported it absent")
	}
	if string(content) != "what the commit says\n" {
		t.Errorf("the read answered %q, which is not what the commit holds", content)
	}
}

func TestAFileTheCommitDoesNotCarryIsAbsentRatherThanAFailure(t *testing.T) {
	repository := gitrepo.New(t)
	base := repository.Head()

	content, present, err := snapshot.FileAt(context.Background(), repository.Root, base, "nothing/here.md")
	if err != nil {
		t.Fatalf("reading a path the commit does not carry answered an error: %v", err)
	}
	if present {
		t.Error("a path the commit does not carry was reported present")
	}
	if content != nil {
		t.Errorf("a path the commit does not carry answered %q", content)
	}
}

func TestADirectoryAtThePathIsNotAFile(t *testing.T) {
	repository := gitrepo.New(t)
	repository.Write("notes/inner.md", "inside\n")
	base := repository.Commit("a directory where a file was asked for")

	_, present, err := snapshot.FileAt(context.Background(), repository.Root, base, "notes")
	if err != nil {
		t.Fatalf("reading a directory answered an error: %v", err)
	}
	if present {
		t.Error("a directory was reported as a file the commit carries")
	}
}

func TestReadingAFileAtACommitWritesNothing(t *testing.T) {
	// scenario: Reading the guide writes nothing
	repository := gitrepo.New(t)
	repository.Write("notes.md", "committed\n")
	base := repository.Commit("a file to read")
	repository.Write("untracked.md", "not committed\n")

	before := repository.Status()
	if _, _, err := snapshot.FileAt(context.Background(), repository.Root, base, "notes.md"); err != nil {
		t.Fatalf("reading notes.md at %s: %v", base, err)
	}
	if after := repository.Status(); after != before {
		t.Errorf("reading a file at a commit changed what git status reports:\nbefore:\n%s\nafter:\n%s", before, after)
	}
}
