package snapshot_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/rui-ktei/socrates/internal/snapshot"
	"github.com/rui-ktei/socrates/internal/testing/gitrepo"
)

func TestNotARepository(t *testing.T) {
	gitrepo.Hermetic(t)
	directory := t.TempDir()

	refusal := refuse(t, directory)

	if !strings.Contains(refusal.Reason, directory) {
		t.Errorf("the refusal reads %q, without naming the directory", refusal.Reason)
	}
	if !strings.Contains(refusal.Reason, "git") {
		t.Errorf("the refusal reads %q, without naming the condition", refusal.Reason)
	}
}

func TestARepositoryWithNoCommits(t *testing.T) {
	repository := gitrepo.Empty(t)
	repository.Write("first.txt", "not committed yet\n")

	refusal := refuse(t, repository.Root)

	if !strings.Contains(refusal.Reason, "base") {
		t.Errorf("the refusal reads %q, without saying there is no base to review against", refusal.Reason)
	}
}

func TestABareRepository(t *testing.T) {
	repository := gitrepo.Bare(t)

	refusal := refuse(t, repository.Root)

	if !strings.Contains(refusal.Reason, "working tree") {
		t.Errorf("the refusal reads %q, without saying there is no working tree to review", refusal.Reason)
	}
}

func TestARefusalIsNotAFailure(t *testing.T) {
	gitrepo.Hermetic(t)

	refusal := refuse(t, t.TempDir())
	if refusal.Error() != refusal.Reason {
		t.Errorf("the refusal reports %q rather than the reason it holds", refusal.Error())
	}
}

func refuse(t *testing.T, directory string) *snapshot.Refusal {
	t.Helper()

	taken, err := snapshot.Take(context.Background(), directory)
	if err == nil {
		t.Fatalf("taking a snapshot of %s succeeded, returning %+v", directory, taken)
	}

	var refusal *snapshot.Refusal
	if !errors.As(err, &refusal) {
		t.Fatalf("taking a snapshot of %s failed rather than refusing: %v", directory, err)
	}
	return refusal
}
