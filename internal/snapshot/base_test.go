package snapshot_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/refusal"
	"github.com/socrates-tools/socrates/internal/snapshot"
	"github.com/socrates-tools/socrates/internal/testing/gitrepo"
)

func branchedFrom(t *testing.T) (*gitrepo.Repository, string) {
	t.Helper()

	repository := gitrepo.New(t)
	repository.Write("trunk.txt", "what the trunk had when the branch left it\n")
	forked := repository.Commit("the commit the branch grew out of")

	repository.Git("switch", "--quiet", "--create", "feature")
	repository.Write("branch.txt", "what the branch added\n")
	repository.Commit("the branch's own commit")

	repository.Git("switch", "--quiet", "main")
	repository.Write("trunk.txt", "what the trunk did afterwards\n")
	repository.Write("moved.txt", "a file the trunk added and the branch never saw\n")
	repository.Commit("the trunk moved on")

	repository.Git("switch", "--quiet", "feature")
	return repository, forked
}

func TestANamedBaseResolvesToTheForkPoint(t *testing.T) {
	repository, forked := branchedFrom(t)

	resolved, err := snapshot.ForkPoint(context.Background(), repository.Root, "main")
	if err != nil {
		t.Fatalf("resolving main as a base: %v", err)
	}
	if resolved != forked {
		t.Errorf("main resolves to %s, want the fork point %s: the diff would report every file the trunk changed as one the branch deleted", resolved, forked)
	}
}

func TestASnapshotAgainstANamedBaseHoldsTheBranchesOwnWork(t *testing.T) {
	repository, forked := branchedFrom(t)
	repository.Write("uncommitted.txt", "what the branch has not committed yet\n")

	taken, err := snapshot.Take(context.Background(), repository.Root, "main")
	if err != nil {
		t.Fatalf("taking a snapshot against main: %v", err)
	}
	if taken.Base != forked {
		t.Errorf("the snapshot is based on %s, want the fork point %s", taken.Base, forked)
	}

	diff, err := snapshot.Diff(context.Background(), taken.RepositoryRoot, taken.Base, taken.Commit)
	if err != nil {
		t.Fatalf("reading the diff of the snapshot: %v", err)
	}
	for _, own := range []string{"branch.txt", "uncommitted.txt"} {
		if !strings.Contains(diff, own) {
			t.Errorf("the diff does not hold %s, which is the branch's own work:\n%s", own, diff)
		}
	}
	if strings.Contains(diff, "moved.txt") {
		t.Errorf("the diff holds moved.txt, which the trunk added and the branch never touched:\n%s", diff)
	}
	if strings.Contains(diff, "what the trunk did afterwards") {
		t.Errorf("the diff holds what the trunk did after the fork point:\n%s", diff)
	}
}

func TestNoNamedBaseIsTheCommitHeadResolvesTo(t *testing.T) {
	repository, _ := branchedFrom(t)
	repository.Write("uncommitted.txt", "what the branch has not committed yet\n")

	taken, err := snapshot.Take(context.Background(), repository.Root, "")
	if err != nil {
		t.Fatalf("taking a snapshot against nothing named: %v", err)
	}
	if taken.Base != repository.Head() {
		t.Errorf("the snapshot is based on %s, want the commit HEAD resolves to, %s", taken.Base, repository.Head())
	}
	if diff := len(taken.Changed); diff != 1 {
		t.Errorf("the snapshot holds %v as changed, want the one uncommitted file", taken.Changed)
	}
}

func TestABaseThatDoesNotResolve(t *testing.T) {
	repository, _ := branchedFrom(t)

	refused := refuseBase(t, repository.Root, "origin/develop")

	if !strings.Contains(refused.Reason, "origin/develop") {
		t.Errorf("the refusal reads %q, without naming what was asked for", refused.Reason)
	}
	if !strings.Contains(refused.Reason, repository.Root) {
		t.Errorf("the refusal reads %q, without naming the repository it looked in", refused.Reason)
	}
}

func TestABaseWithNoCommonAncestor(t *testing.T) {
	repository, _ := branchedFrom(t)
	repository.Git("checkout", "--quiet", "--orphan", "unrelated")
	repository.Write("unrelated.txt", "a history of its own\n")
	repository.Commit("a root commit sharing nothing")
	repository.Git("switch", "--quiet", "feature")

	refused := refuseBase(t, repository.Root, "unrelated")

	if !strings.Contains(refused.Reason, "unrelated") {
		t.Errorf("the refusal reads %q, without naming the base", refused.Reason)
	}
	for _, said := range []string{"share", "ancestor", "common"} {
		if strings.Contains(refused.Reason, said) {
			return
		}
	}
	t.Errorf("the refusal reads %q, without saying the two histories have nothing in common", refused.Reason)
}

func refuseBase(t *testing.T, directory, ref string) *refusal.Refusal {
	t.Helper()

	taken, err := snapshot.Take(context.Background(), directory, ref)
	if err == nil {
		t.Fatalf("taking a snapshot against %s succeeded, returning %+v", ref, taken)
	}

	var refused *refusal.Refusal
	if !errors.As(err, &refused) {
		t.Fatalf("taking a snapshot against %s failed rather than refusing: %v", ref, err)
	}
	return refused
}
