package snapshot_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates-review/internal/refusal"
	"github.com/socrates-tools/socrates-review/internal/snapshot"
	"github.com/socrates-tools/socrates-review/internal/testing/gitrepo"
)

func TestNotARepository(t *testing.T) {
	gitrepo.Hermetic(t)
	directory := t.TempDir()

	refused := refuse(t, directory)

	if !strings.Contains(refused.Reason, directory) {
		t.Errorf("the refusal reads %q, without naming the directory", refused.Reason)
	}
	if !strings.Contains(refused.Reason, "git") {
		t.Errorf("the refusal reads %q, without naming the condition", refused.Reason)
	}
}

func TestARepositoryWithNoCommits(t *testing.T) {
	repository := gitrepo.Empty(t)
	repository.Write("first.txt", "not committed yet\n")

	refused := refuse(t, repository.Root)

	if !strings.Contains(refused.Reason, "base") {
		t.Errorf("the refusal reads %q, without saying there is no base to review against", refused.Reason)
	}
}

func TestABareRepository(t *testing.T) {
	repository := gitrepo.Bare(t)

	refused := refuse(t, repository.Root)

	if !strings.Contains(refused.Reason, "working tree") {
		t.Errorf("the refusal reads %q, without saying there is no working tree to review", refused.Reason)
	}
}

func TestARefusalIsNotAFailure(t *testing.T) {
	gitrepo.Hermetic(t)

	refused := refuse(t, t.TempDir())
	if refused.Error() != refused.Reason {
		t.Errorf("the refusal reports %q rather than the reason it holds", refused.Error())
	}
}

func refuse(t *testing.T, directory string) *refusal.Refusal {
	t.Helper()

	taken, err := snapshot.Take(context.Background(), directory, "")
	if err == nil {
		t.Fatalf("taking a snapshot of %s succeeded, returning %+v", directory, taken)
	}

	var refused *refusal.Refusal
	if !errors.As(err, &refused) {
		t.Fatalf("taking a snapshot of %s failed rather than refusing: %v", directory, err)
	}
	return refused
}
