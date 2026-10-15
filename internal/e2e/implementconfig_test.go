package e2e_test

import (
	skillkit "github.com/socrates-tools/socrates/internal/skill"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/implement/config"
	"github.com/socrates-tools/socrates/implement/skill"
	"github.com/socrates-tools/socrates/internal/exit"
	"github.com/socrates-tools/socrates/internal/testing/harness"
)

func implementConfigurationPath(session *harness.Session) string {
	return filepath.Join(session.Root, "implement", "config.yaml")
}

func stateImplementConfiguration(t *testing.T, session *harness.Session, content string) {
	t.Helper()

	path := implementConfigurationPath(session)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("creating the configuration directory: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
}

func implementSettingsIn(t *testing.T, printed string) map[string][2]string {
	t.Helper()

	reported := map[string][2]string{}
	for line := range strings.SplitSeq(printed, "\n") {
		fields := strings.Fields(line)
		if len(fields) != 3 {
			continue
		}
		reported[fields[0]] = [2]string{fields[1], fields[2]}
	}
	return reported
}

func TestImplementConfigWithNoFileReportsTheThreeSettingsAtTheirDefaults(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)

	printed, failures, code := session.Top("implement", "config")
	if code != exit.NothingAsked {
		t.Fatalf("socrates implement config with no file exited %d: %s", code, failures)
	}

	want := map[string][2]string{
		"concurrency":              {"4", string(config.Default)},
		"roles.implementer.model":  {"sonnet", string(config.Default)},
		"roles.implementer.effort": {"high", string(config.Default)},
	}
	if diff := cmp.Diff(want, implementSettingsIn(t, printed)); diff != "" {
		t.Errorf("what socrates implement config reports differs (-want +got):\n%s", diff)
	}
	if first := firstLine(printed); first != implementConfigurationPath(session) {
		t.Errorf("socrates implement config names %q, and the file under the stated home is %s", first, implementConfigurationPath(session))
	}
	if _, err := os.Stat(implementConfigurationPath(session)); !os.IsNotExist(err) {
		t.Errorf("socrates implement config created %s", implementConfigurationPath(session))
	}
}

func TestImplementConfigReportsWhatAPersonWroteAsSet(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	stateImplementConfiguration(t, session, "concurrency: 2\nroles:\n  implementer:\n    effort: max\n")

	printed, failures, code := session.Top("implement", "config")
	if code != exit.NothingAsked {
		t.Fatalf("socrates implement config exited %d: %s", code, failures)
	}

	want := map[string][2]string{
		"concurrency":              {"2", string(config.Set)},
		"roles.implementer.model":  {"sonnet", string(config.Default)},
		"roles.implementer.effort": {"max", string(config.Set)},
	}
	if diff := cmp.Diff(want, implementSettingsIn(t, printed)); diff != "" {
		t.Errorf("what socrates implement config reports differs (-want +got):\n%s", diff)
	}
}

func TestImplementConfigRefusesAValueOutOfRangeAndAFieldItDoesNotHaveNamingTheFile(t *testing.T) {
	t.Parallel()
	for name, content := range map[string]string{
		"a concurrency of zero": "concurrency: 0\n",
		"an unknown effort":     "roles:\n  implementer:\n    effort: extreme\n",
		"an unknown field":      "parallelism: 3\n",
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			session := harness.Start(t)
			stateImplementConfiguration(t, session, content)

			printed, failures, code := session.Top("implement", "config")
			if code != exit.Unusable {
				t.Fatalf("socrates implement config with %s exited %d, want %d", name, code, exit.Unusable)
			}
			if !strings.Contains(failures, implementConfigurationPath(session)) {
				t.Errorf("the refusal does not name the file: %s", failures)
			}
			if printed != "" {
				t.Errorf("socrates implement config reported settings while the file is refused:\n%s", printed)
			}
		})
	}
}

func TestImplementReadsNothingFromReviewsFile(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	stateConfiguration(t, session, "diff_ceiling: 4096\n")

	printed, failures, code := session.Top("implement", "config")
	if code != exit.NothingAsked {
		t.Fatalf("socrates implement config exited %d: %s", code, failures)
	}
	for _, reported := range implementSettingsIn(t, printed) {
		if reported[1] != string(config.Default) {
			t.Errorf("a setting reads as %s while only review's file states anything:\n%s", reported[1], printed)
		}
	}
}

func TestTheTopLevelHelpListsImplementAndImplementHelpListsWhatItHolds(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)

	top, failures, code := session.Top("--help")
	if code != exit.NothingAsked || !strings.Contains(top, "implement") {
		t.Errorf("socrates --help exited %d and reads %q%s, without listing implement", code, top, failures)
	}

	own, failures, code := session.Top("implement", "--help")
	if code != exit.NothingAsked || !strings.Contains(own, "config") {
		t.Errorf("socrates implement --help exited %d and reads %q%s, without listing config", code, own, failures)
	}
}

func defaultBranchReported(t *testing.T, printed string) (string, bool) {
	t.Helper()

	for line := range strings.SplitSeq(printed, "\n") {
		fields := strings.Fields(line)
		if len(fields) > 0 && fields[0] == "default_branch" {
			return strings.Join(fields[1:len(fields)-1], " "), true
		}
	}
	return "", false
}

func TestImplementDefaultBranchSavesTheBranchForTheRepositoryAndConfigReportsItThere(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	session.Repository.Git("branch", "dev")

	printed, failures, code := session.Top("implement", "default-branch", "dev")
	if code != exit.NothingAsked {
		t.Fatalf("socrates implement default-branch dev exited %d: %s", code, failures)
	}
	if !strings.Contains(printed, session.Repository.Root) {
		t.Errorf("the confirmation does not name the repository: %q", printed)
	}

	written := session.File(implementConfigurationPath(session))
	if want := "default_branches:\n  " + session.Repository.Root + ": dev\n"; written != want {
		t.Errorf("the file holds:\n%s\nwant:\n%s", written, want)
	}

	reported, failures, code := session.Top("implement", "config")
	if code != exit.NothingAsked {
		t.Fatalf("socrates implement config exited %d: %s", code, failures)
	}
	if branch, present := defaultBranchReported(t, reported); !present || branch != "dev" {
		t.Errorf("config reports default_branch %q (present %v):\n%s", branch, present, reported)
	}
	if !strings.HasSuffix(strings.TrimSpace(lastLine(reported)), "set") {
		t.Errorf("the default branch is not reported as set:\n%s", reported)
	}
}

func lastLine(printed string) string {
	lines := strings.Split(strings.TrimSpace(printed), "\n")
	return lines[len(lines)-1]
}

func TestImplementDefaultBranchAcceptsARemoteTrackingBranch(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	session.Repository.Git("branch", "dev")
	clone := filepath.Join(t.TempDir(), "clone")
	session.Repository.Git("clone", "--quiet", session.Repository.Root, clone)

	_, failures, code := session.TopIn(clone, "implement", "default-branch", "origin/dev")
	if code != exit.NothingAsked {
		t.Fatalf("saving a remote-tracking branch exited %d: %s", code, failures)
	}
}

func TestImplementDefaultBranchRefusesABranchTheRepositoryDoesNotHoldAndLeavesTheFileAlone(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	session.Repository.Git("tag", "release")
	held := "# mine\nconcurrency: 2\n"
	stateImplementConfiguration(t, session, held)

	for _, branch := range []string{"nowhere", "release", "origin/dev"} {
		_, failures, code := session.Top("implement", "default-branch", branch)
		if code != exit.Unusable {
			t.Errorf("saving %q exited %d, want %d", branch, code, exit.Unusable)
		}
		if !strings.Contains(failures, branch) || !strings.Contains(failures, "local or remote-tracking") {
			t.Errorf("the refusal of %q does not say what a repository holds: %s", branch, failures)
		}
		if now := session.File(implementConfigurationPath(session)); now != held {
			t.Errorf("a refused save touched the file:\n%s", now)
		}
	}
}

func TestAWorktreeOfTheCloneReadsTheCloneEntryAndASecondCloneReadsNothing(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	session.Repository.Git("branch", "dev")
	worktree := filepath.Join(t.TempDir(), "worktree")
	session.Repository.Git("worktree", "add", "--quiet", "-b", "feature", worktree)
	secondClone := filepath.Join(t.TempDir(), "second.clone")
	session.Repository.Git("clone", "--quiet", session.Repository.Root, secondClone)

	if _, failures, code := session.Top("implement", "default-branch", "dev"); code != exit.NothingAsked {
		t.Fatalf("saving from the clone exited %d: %s", code, failures)
	}

	for name, directory := range map[string]string{"the clone": session.Repository.Root, "a worktree of it": worktree} {
		reported, failures, code := session.TopIn(directory, "implement", "config")
		if code != exit.NothingAsked {
			t.Fatalf("config from %s exited %d: %s", name, code, failures)
		}
		if branch, _ := defaultBranchReported(t, reported); branch != "dev" {
			t.Errorf("from %s default_branch reads %q, want dev:\n%s", name, branch, reported)
		}
	}

	reported, failures, code := session.TopIn(secondClone, "implement", "config")
	if code != exit.NothingAsked {
		t.Fatalf("config from the second clone exited %d: %s", code, failures)
	}
	if branch, present := defaultBranchReported(t, reported); !present || branch != "" {
		t.Errorf("a second clone at another path reports default_branch %q (present %v), want it present and empty:\n%s", branch, present, reported)
	}
	if !strings.HasSuffix(lastLine(reported), "default") {
		t.Errorf("an unremembered default branch is not reported as default:\n%s", reported)
	}
}

func TestSavingFromAWorktreeWritesTheEntryOfTheCloneAndAPathWithADotStaysOneKey(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	session.Repository.Git("branch", "dev")
	worktree := filepath.Join(t.TempDir(), "my.worktree")
	session.Repository.Git("worktree", "add", "--quiet", "-b", "feature", worktree)

	if _, failures, code := session.TopIn(worktree, "implement", "default-branch", "dev"); code != exit.NothingAsked {
		t.Fatalf("saving from a worktree exited %d: %s", code, failures)
	}

	if want := "default_branches:\n  " + session.Repository.Root + ": dev\n"; session.File(implementConfigurationPath(session)) != want {
		t.Errorf("the file holds:\n%s\nwant:\n%s", session.File(implementConfigurationPath(session)), want)
	}
}

func TestADefaultBranchIsKeyedByAPathHoldingADotAndASecondRepositoryAddsASecondEntry(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	session.Repository.Git("branch", "dev")
	dotted := filepath.Join(t.TempDir(), "my.repo")
	session.Repository.Git("clone", "--quiet", session.Repository.Root, dotted)
	resolved, err := filepath.EvalSymlinks(dotted)
	if err != nil {
		t.Fatal(err)
	}

	for directory, branch := range map[string]string{session.Repository.Root: "dev", dotted: "origin/dev"} {
		if _, failures, code := session.TopIn(directory, "implement", "default-branch", branch); code != exit.NothingAsked {
			t.Fatalf("saving %s exited %d: %s", branch, code, failures)
		}
	}

	loaded, err := config.LoadFrom(implementConfigurationPath(session))
	if err != nil {
		t.Fatalf("loading: %v", err)
	}
	want := map[string]string{session.Repository.Root: "dev", resolved: "origin/dev"}
	if diff := cmp.Diff(want, loaded.DefaultBranches); diff != "" {
		t.Errorf("the entries differ (-want +got):\n%s", diff)
	}
}

func TestOutsideARepositoryDefaultBranchIsRefusedNamingWhatItNeedsAndConfigReportsNoDefaultBranch(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	outside := t.TempDir()

	_, failures, code := session.TopIn(outside, "implement", "default-branch", "main")
	if code != exit.Unusable {
		t.Errorf("default-branch outside a repository exited %d, want %d", code, exit.Unusable)
	}
	if !strings.Contains(failures, "inside a git repository") {
		t.Errorf("the refusal does not name what it needs: %s", failures)
	}
	if _, err := os.Stat(implementConfigurationPath(session)); !os.IsNotExist(err) {
		t.Errorf("a refused save created the file")
	}

	reported, failures, code := session.TopIn(outside, "implement", "config")
	if code != exit.NothingAsked {
		t.Fatalf("config outside a repository exited %d: %s", code, failures)
	}
	if _, present := defaultBranchReported(t, reported); present {
		t.Errorf("config outside a repository reports a default_branch:\n%s", reported)
	}
}

func TestDefaultBranchKeepsEveryCommentKeyOrderAndBlankLineOfTheFile(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	session.Repository.Git("branch", "dev")
	held := "# how many at once\nconcurrency: 2\n\n# which model\nroles:\n  implementer:\n    effort: max\n    model: opus\n"
	stateImplementConfiguration(t, session, held)

	if _, failures, code := session.Top("implement", "default-branch", "dev"); code != exit.NothingAsked {
		t.Fatalf("saving exited %d: %s", code, failures)
	}

	want := held + "default_branches:\n  " + session.Repository.Root + ": dev\n"
	if now := session.File(implementConfigurationPath(session)); now != want {
		t.Errorf("the file holds:\n%s\nwant:\n%s", now, want)
	}

	if _, failures, code := session.Top("implement", "default-branch", "main"); code != exit.NothingAsked {
		t.Fatalf("saving again exited %d: %s", code, failures)
	}
	if now := session.File(implementConfigurationPath(session)); !strings.HasPrefix(now, held) || !strings.HasSuffix(now, ": main\n") {
		t.Errorf("changing the branch disturbed the rest:\n%s", now)
	}
}

func TestInABareRepositoryDefaultBranchIsRefusedNamingTheMainWorktreeItLacks(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	bare := filepath.Join(t.TempDir(), "bare.git")
	session.Repository.Git("clone", "--quiet", "--bare", session.Repository.Root, bare)

	_, failures, code := session.TopIn(bare, "implement", "default-branch", "main")
	if code != exit.Unusable || !strings.Contains(failures, "main worktree") {
		t.Errorf("default-branch in a bare repository exited %d with %q", code, failures)
	}
}

func TestImplementSkillWritesTheDefinitionUnderTheDefaultDirectoryAtTheDefaults(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)

	printed, failures, code := session.Top("implement", "skill")
	if code != exit.NothingAsked {
		t.Fatalf("socrates implement skill exited %d: %s", code, failures)
	}

	path := filepath.Join(skillkit.ClaudeCode.AgentsDirectoryUnder(session.UserHome), skill.AgentFileName())
	if !strings.Contains(printed, path) {
		t.Errorf("socrates implement skill does not say where it wrote, and %s is where it did:\n%s", path, printed)
	}
	if written := session.File(path); written != skill.DefaultDefinition() {
		t.Errorf("what was installed is not the definition at the defaults:\n%s", written)
	}
}

func TestImplementSkillWritesTheSettingsInForceUnderAgentsToAndRecordsTheLastPlace(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	stateImplementConfiguration(t, session, "roles:\n  implementer:\n    model: opus\n    effort: max\n")
	first, second := t.TempDir(), t.TempDir()

	for _, directory := range []string{first, second} {
		if _, failures, code := session.Top("implement", "skill", "--host", skillkit.ClaudeCode.Flag, "--agents-to", directory); code != exit.NothingAsked {
			t.Fatalf("socrates implement skill --host claude --agents-to %s exited %d: %s", directory, code, failures)
		}
	}

	written := session.File(filepath.Join(second, skill.AgentFileName()))
	for _, want := range []string{"model: opus\n", "effort: max\n"} {
		if !strings.Contains(written, want) {
			t.Errorf("the definition does not state %q:\n%s", want, written)
		}
	}
	record := session.File(filepath.Join(session.Root, "implement", skill.RecordFile))
	if !strings.Contains(record, second) || strings.Contains(record, first) {
		t.Errorf("the record does not name only the last directory %s:\n%s", second, record)
	}
	if strings.Contains(session.File(implementConfigurationPath(session)), second) {
		t.Error("the configuration file names where the definition was installed")
	}
}
