package e2e_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates-review/internal/prompt"
	"github.com/socrates-tools/socrates-review/internal/runstate"
	"github.com/socrates-tools/socrates-review/internal/testing/harness"
)

const theBranchesOwnWork = "package internal\n\nfunc Added(document string) int {\n\treturn len(document)\n}\n"

func aBranchAndATrunkThatMoved(t *testing.T, session *harness.Session) string {
	t.Helper()

	repository := session.Repository
	repository.Write("shared/trunk.go", "package shared\n\nconst What = \"the trunk when the branch left it\"\n")
	forked := repository.Commit("what the trunk had when the branch left it")

	repository.Git("switch", "--quiet", "--create", "feature")
	repository.Write("internal/added.go", theBranchesOwnWork)
	repository.Commit("the branch's own committed work")

	repository.Git("switch", "--quiet", "main")
	repository.Write("shared/trunk.go", "package shared\n\nconst What = \"what the trunk did afterwards\"\n")
	repository.Write("shared/moved.go", "package shared\n\nfunc Moved() {}\n")
	repository.Commit("the trunk moved on")

	repository.Git("switch", "--quiet", "feature")
	return forked
}

func TestANamedBaseReviewsTheBranchesOwnWork(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(cleanScenario)

	forked := aBranchAndATrunkThatMoved(t, session)
	session.Repository.Write("internal/uncommitted.go", "package internal\n\nfunc Uncommitted() {}\n")

	name := session.Review("--base", "main", "--intent", "count the characters of a document")
	run := session.WaitFor(name, runstate.Clean)

	if run.Base != forked {
		t.Errorf("the run is based on %s, want the fork point %s: the diff would report every file the trunk changed as one the branch deleted", run.Base, forked)
	}

	delivered := session.Simulator.OnlyInvocation().Prompt
	sections := prompt.Sections(delivered)

	touched, present := prompt.BodyOf(sections, prompt.ChangedHeading)
	if !present {
		t.Fatalf("the reviewer was told nothing about which files the change touches:\n%s", delivered)
	}
	want := []string{"- internal/added.go", "- internal/uncommitted.go"}
	if diff := cmp.Diff(want, strings.Split(strings.TrimSpace(touched), "\n")); diff != "" {
		t.Errorf("the files the reviewer was told the change touches differ (-want +got):\n%s", diff)
	}

	change, present := prompt.BodyOf(sections, prompt.ChangeHeading)
	if !present {
		t.Fatalf("the reviewer was given no diff:\n%s", delivered)
	}
	for _, own := range []string{"func Added(document string) int", "func Uncommitted()"} {
		if !strings.Contains(change, own) {
			t.Errorf("the diff the reviewer was given does not hold %q, which the branch wrote:\n%s", own, change)
		}
	}
	for _, theTrunks := range []string{"shared/moved.go", "the trunk when the branch left it"} {
		if strings.Contains(change, theTrunks) {
			t.Errorf("the diff the reviewer was given holds %q, which is the trunk's work rather than the branch's:\n%s", theTrunks, change)
		}
	}
}

func TestTheStatusPrintsTheBaseItResolved(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(cleanScenario)

	forked := aBranchAndATrunkThatMoved(t, session)
	session.Repository.Write("internal/uncommitted.go", "package internal\n\nfunc Uncommitted() {}\n")

	name := session.Review("--base", "main")
	session.WaitFor(name, runstate.Clean)

	printed, failures, code := session.Socrates("status", name)
	if code != 0 {
		t.Fatalf("socrates-review status exited %d: %s", code, failures)
	}
	if !strings.Contains(printed, "base: "+forked) {
		t.Errorf("the status does not say the change was reviewed against %s:\n%s", forked, printed)
	}
}

func TestANamedBaseIsPinnedForTheWholeRun(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(dialogueScenario(patchableFindings, nothingFoundAgain))

	forked := aBranchAndATrunkThatMoved(t, session)
	session.Repository.Write("internal/added.go", answeredWrongly)

	name := session.Review("--base", "main", "--intent", "count the characters of a document")
	session.WaitFor(name, runstate.AwaitingResponse)

	session.Repository.Write("internal/added.go", answeredBetter)
	respond(t, session, name, "--finding", "f1", "--fixed", "--note", "an empty document now answers minus one")
	respond(t, session, name, "--finding", "f2", "--rejected", "--reason", "there is no retry loop in this function at all")
	session.Repository.Commit("the author committed the patches while the run waited")

	if _, failures, code := session.Socrates("continue", name, "--wait"); code != 0 {
		t.Fatalf("socrates-review continue --wait exited %d: %s", code, failures)
	}

	second, present := session.Recorded(name).Reviewed(2)
	if !present {
		t.Fatal("the run records no second round")
	}
	if second.Base != forked {
		t.Errorf("the second round is diffed against %s, want the fork point the run began with, %s", second.Base, forked)
	}

	delivered := session.Simulator.Invocations()[1].Prompt
	if !strings.Contains(delivered, "+func Answer(document string) int {") {
		t.Errorf("the second round was diffed against the commit the author made rather than against the run's base:\n%s", delivered)
	}
}

func TestABranchWithNothingAheadOfItsBaseIsNotCalledClean(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(cleanScenario)

	session.Repository.Git("switch", "--quiet", "--create", "feature")

	printed, failures, code := session.Socrates("review", "--base", "main")
	if code != 3 {
		t.Fatalf("reviewing a branch with nothing ahead of its base exited %d, want 3: %s\n%s", code, failures, printed)
	}
	for _, said := range []string{"nothing", "--base"} {
		if !strings.Contains(failures, said) {
			t.Errorf("the refusal does not say %q, and a person reading it has to work out why nothing was reviewed: %s", said, failures)
		}
	}
	if strings.Contains(failures, string(runstate.Clean)) {
		t.Errorf("the refusal reads as a clean verdict: %s", failures)
	}
	if names := session.Names(); len(names) != 0 {
		t.Errorf("a refused review left the runs %v behind, and a session acting on the state of one would read a pass", names)
	}
	if invocations := session.Simulator.Invocations(); len(invocations) != 0 {
		t.Errorf("a refused review spawned %d reviewers", len(invocations))
	}
}

func TestAWorkingTreeWithNothingInItIsNotCalledCleanEither(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(cleanScenario)

	printed, failures, code := session.Socrates("review")
	if code != 3 {
		t.Fatalf("reviewing a working tree identical to its base exited %d, want 3: %s\n%s", code, failures, printed)
	}
	if names := session.Names(); len(names) != 0 {
		t.Errorf("a refused review left the runs %v behind", names)
	}
}

func TestABaseNothingResolvesIsRefusedBeforeAnythingHappens(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(cleanScenario)
	session.Repository.Write("internal/added.go", theBranchesOwnWork)

	printed, failures, code := session.Socrates("review", "--base", "origin/develop")
	if code != 3 {
		t.Fatalf("reviewing against a base that does not resolve exited %d, want 3: %s\n%s", code, failures, printed)
	}
	if !strings.Contains(failures, "origin/develop") {
		t.Errorf("the refusal does not name what was asked for: %s", failures)
	}
	if names := session.Names(); len(names) != 0 {
		t.Errorf("a refused base left the runs %v behind", names)
	}
	if invocations := session.Simulator.Invocations(); len(invocations) != 0 {
		t.Errorf("a refused base spawned %d reviewers", len(invocations))
	}
}

func TestABaseWithNoCommonAncestorIsRefusedAndNamed(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(cleanScenario)

	session.Repository.Git("checkout", "--quiet", "--orphan", "unrelated")
	session.Repository.Write("unrelated.go", "package unrelated\n")
	session.Repository.Commit("a root commit sharing nothing with main")
	session.Repository.Write("internal/added.go", theBranchesOwnWork)

	printed, failures, code := session.Socrates("review", "--base", "main")
	if code != 3 {
		t.Fatalf("reviewing against a base with no common ancestor exited %d, want 3: %s\n%s", code, failures, printed)
	}
	if !strings.Contains(failures, "main") {
		t.Errorf("the refusal does not name the base: %s", failures)
	}
	if names := session.Names(); len(names) != 0 {
		t.Errorf("a refused base left the runs %v behind", names)
	}
}

func TestReviewingUncommittedWorkNamesNoBase(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(cleanScenario)
	session.Repository.Write("internal/added.go", theBranchesOwnWork)

	head := session.Repository.Head()
	name := session.Review("--intent", "count the characters of a document")
	run := session.WaitFor(name, runstate.Clean)

	if run.Base != head {
		t.Errorf("a review that named no base is based on %s, want the commit HEAD resolves to, %s", run.Base, head)
	}
}

func TestSocratesReachesNoRemoteToResolveABase(t *testing.T) {
	// scenario: socrates-review neither fetches nor rebases
	session := harness.Start(t)
	session.Simulator.Scenario(cleanScenario)

	forked := aBranchAndATrunkThatMoved(t, session)
	session.Repository.Write("internal/uncommitted.go", "package internal\n\nfunc Uncommitted() {}\n")

	branch := session.Repository.Git("rev-parse", "--abbrev-ref", "HEAD")
	before := session.Repository.Git("rev-list", "--all")
	reflog := session.Repository.Git("reflog", "--all")

	name := session.Review("--base", "main")
	session.WaitFor(name, runstate.Clean)

	if after := session.Repository.Git("rev-list", "--all"); after != before {
		t.Errorf("the commits reachable in the repository changed while the base was resolved:\nbefore:\n%s\nafter:\n%s", before, after)
	}
	if after := session.Repository.Git("reflog", "--all"); after != reflog {
		t.Errorf("resolving the base wrote to the reflog:\nbefore:\n%s\nafter:\n%s", reflog, after)
	}
	if after := session.Repository.Git("rev-parse", "--abbrev-ref", "HEAD"); after != branch {
		t.Errorf("resolving the base left the repository on %s rather than on %s", after, branch)
	}
	if run := session.Recorded(name); run.Base != forked {
		t.Errorf("the run is based on %s, want %s", run.Base, forked)
	}
}

func TestAReviewFromAWorktreeOfAnotherBranchIsOfThatBranchAlone(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(cleanScenario)

	forked := aBranchAndATrunkThatMoved(t, session)
	session.Repository.Git("switch", "--quiet", "main")
	session.Repository.Write("mine/unrelated.go", "package mine\n\nfunc Mine() {}\n")

	elsewhere := filepath.Join(t.TempDir(), "reviewing-their-change")
	session.Repository.Git("worktree", "add", "--quiet", "--detach", elsewhere, "feature")

	out, failures, code := session.SocratesIn(elsewhere, "review", "--base", "main", "--not-my-change", "--wait")
	if code != 0 {
		t.Fatalf("reviewing from a worktree of another branch exited %d: %s\n%s", code, failures, out)
	}

	name := firstLine(out)
	run := session.Run(name)
	if run.Base != forked {
		t.Errorf("the run is based on %s, want the fork point %s", run.Base, forked)
	}
	if run.Repository != elsewhere {
		t.Errorf("the run reads %s as the repository it reviewed, want the worktree it was started from, %s", run.Repository, elsewhere)
	}

	touched, present := prompt.BodyOf(prompt.Sections(session.Simulator.OnlyInvocation().Prompt), prompt.ChangedHeading)
	if !present {
		t.Fatal("the reviewer was told nothing about which files the change touches")
	}
	if !strings.Contains(touched, "internal/added.go") {
		t.Errorf("the reviewer was not given the branch's own work:\n%s", touched)
	}
	if strings.Contains(touched, "mine/unrelated.go") {
		t.Errorf("the reviewer was given the edits sitting in the other working tree, which are nobody's business here:\n%s", touched)
	}
}
