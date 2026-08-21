package e2e_test

import (
	"encoding/json"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/rui-ktei/socrates/internal/runstate"
	"github.com/rui-ktei/socrates/internal/testing/gitrepo"
	"github.com/rui-ktei/socrates/internal/testing/harness"
)

func TestTheCommandSurfaceIsWhatItSays(t *testing.T) {
	session := harness.Start(t)

	printed, failures, code := session.Scrat("--help")
	if code != 0 {
		t.Fatalf("scrat --help exited %d: %s", code, failures)
	}

	listed := commandsIn(printed)
	want := []string{"cancel", "continue", "findings", "help", "logs", "respond", "review", "status", "watch"}
	if !slices.Equal(listed, want) {
		t.Errorf("scrat lists %v, want %v", listed, want)
	}

	for _, forbidden := range []string{"commit", "push", "pr", "pull-request", "merge", "land", "ci", "test", "lint", "verify"} {
		if slices.Contains(listed, forbidden) {
			t.Errorf("scrat lists %q, which lands or verifies a change", forbidden)
		}
	}
}

func commandsIn(help string) []string {
	var listed []string
	inside := false

	for line := range strings.SplitSeq(help, "\n") {
		if strings.HasPrefix(line, "Available Commands:") {
			inside = true
			continue
		}
		if !inside {
			continue
		}
		if strings.TrimSpace(line) == "" {
			break
		}
		listed = append(listed, strings.Fields(line)[0])
	}

	slices.Sort(listed)
	return listed
}

func TestTheDefaultDoesNotBlock(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(slowScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	printed, failures, code := session.Scrat("review")
	if code != 0 {
		t.Fatalf("scrat review exited %d: %s", code, failures)
	}

	name := firstLine(printed)
	if !strings.Contains(printed, "scrat status "+name) {
		t.Errorf("scrat review does not say how to follow the run:\n%s", printed)
	}

	run := session.WaitFor(name, runstate.Reviewing)
	if run.State.Terminal() {
		t.Fatalf("the run was already %s when scrat review returned", run.State)
	}

	_, _, statusCode := session.Scrat("status", name)
	if statusCode != 4 {
		t.Errorf("scrat status of a running review exited %d, want 4", statusCode)
	}
}

func TestWaitingReturnsTheVerdict(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(cleanScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	printed, failures, code := session.Scrat("review", "--wait")
	if code != 0 {
		t.Fatalf("scrat review --wait exited %d: %s\n%s", code, failures, printed)
	}
	if !strings.Contains(printed, string(runstate.Clean)) {
		t.Errorf("scrat review --wait does not report the verdict:\n%s", printed)
	}

	run := session.Run(firstLine(printed))
	if run.State != runstate.Clean {
		t.Errorf("the run is %s after scrat review --wait returned zero", run.State)
	}
}

func TestStartingIsNotAVerdict(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(escalatingScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	printed, failures, code := session.Scrat("review")
	if code != 0 {
		t.Fatalf("scrat review exited %d for a run that later escalates: %s", code, failures)
	}

	run := session.WaitForTerminal(firstLine(printed))
	if run.State != runstate.Escalated {
		t.Fatalf("the run ended as %s, want escalated: %s", run.State, reason(run))
	}
}

func TestEachVerdictHasItsOwnCodeFromTheCommand(t *testing.T) {
	for _, verdict := range []struct {
		scenario string
		state    runstate.State
		code     int
	}{
		{scenario: cleanScenario, state: runstate.Clean, code: 0},
		{scenario: escalatingScenario, state: runstate.Escalated, code: 1},
		{scenario: "default:\n  behaviour: no-result\n", state: runstate.Failed, code: 2},
	} {
		t.Run(string(verdict.state), func(t *testing.T) {
			session := harness.Start(t)
			session.Simulator.Scenario(verdict.scenario)
			session.Repository.Write("internal/added.go", "package internal\n")

			name := session.Review()
			run := session.WaitForTerminal(name)
			if run.State != verdict.state {
				t.Fatalf("the run ended as %s, want %s: %s", run.State, verdict.state, reason(run))
			}

			if _, _, code := session.Scrat("status", name); code != verdict.code {
				t.Errorf("scrat status of a %s run exited %d, want %d", verdict.state, code, verdict.code)
			}
			if _, _, code := session.Scrat("findings", name); code != verdict.code {
				t.Errorf("scrat findings of a %s run exited %d, want %d", verdict.state, code, verdict.code)
			}
		})
	}
}

func TestAUsageErrorIsNotAVerdictFromTheCommand(t *testing.T) {
	session := harness.Start(t)

	printed, failures, code := session.Scrat("status", "20200101T000000Z-deadbeef")
	if code != 3 {
		t.Errorf("scrat status of a run that does not exist exited %d, want 3", code)
	}
	if !strings.Contains(failures, "20200101T000000Z-deadbeef") {
		t.Errorf("the refusal does not name the run asked for: %s", failures)
	}
	for _, verdict := range []string{string(runstate.Clean), string(runstate.Escalated), string(runstate.Failed)} {
		if strings.Contains(printed, verdict) {
			t.Errorf("a usage error printed the verdict %q:\n%s", verdict, printed)
		}
	}
}

func TestMachineReadableOnRequest(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(escalatingScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	session.WaitForTerminal(name)

	printed, failures, code := session.Scrat("status", "--json", name)
	if code != 1 {
		t.Fatalf("scrat status --json exited %d: %s", code, failures)
	}

	var status struct {
		Run      string `json:"run"`
		State    string `json:"state"`
		Round    int    `json:"round"`
		Findings struct {
			Total          int `json:"total"`
			NeedingAPerson int `json:"needing_a_person"`
		} `json:"findings"`
	}
	if err := json.Unmarshal([]byte(printed), &status); err != nil {
		t.Fatalf("scrat status --json printed something other than one JSON object: %v:\n%s", err, printed)
	}
	if status.Run != name || status.State != string(runstate.Escalated) || status.Round != 1 {
		t.Errorf("the object reads %+v", status)
	}
	if status.Findings.Total != 3 || status.Findings.NeedingAPerson != 1 {
		t.Errorf("the counts read %+v", status.Findings)
	}
	if strings.Count(strings.TrimSpace(printed), "\n") != 0 {
		t.Errorf("scrat status --json printed more than one line:\n%s", printed)
	}
}

func TestFindingsAreReadableFromTheCommand(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(escalatingScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	session.WaitForTerminal(name)

	printed, failures, code := session.Scrat("findings", name)
	if code != 1 {
		t.Fatalf("scrat findings exited %d: %s", code, failures)
	}
	for _, needed := range []string{
		"f1", "error", "internal/added.go:1", "The added file answers wrongly for an empty input.", "needs a person",
		"f2", "can be patched", "f3", "informational",
		"medium", "one changed path answers wrongly without erroring",
	} {
		if !strings.Contains(printed, needed) {
			t.Errorf("scrat findings does not print %q:\n%s", needed, printed)
		}
	}

	asJSON, _, _ := session.Scrat("findings", "--json", name)
	var reported struct {
		Findings []map[string]any `json:"findings"`
	}
	if err := json.Unmarshal([]byte(asJSON), &reported); err != nil {
		t.Fatalf("scrat findings --json printed something other than one JSON object: %v:\n%s", err, asJSON)
	}
	if reported.Findings[0]["action"] != "ask-user" {
		t.Errorf("the finding that stated no action is reported as %v", reported.Findings[0]["action"])
	}
}

func TestTheBareCommandFindsThisRepositorysRun(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(cleanScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	elsewhere := gitrepo.New(t)
	elsewhere.Write("internal/other.go", "package internal\n")

	here := session.Review()
	there := strings.TrimSpace(firstLine(mustReview(t, session, elsewhere.Root)))
	if here == there {
		t.Fatalf("both repositories produced the run %s", here)
	}

	printed, failures, code := session.ScratIn(elsewhere.Root, "status")
	if code != 0 && code != 4 {
		t.Fatalf("scrat status in the second repository exited %d: %s", code, failures)
	}
	if !strings.Contains(printed, there) {
		t.Errorf("scrat status in the second repository reported %q, want the run %s", printed, there)
	}
	if strings.Contains(printed, here) {
		t.Errorf("scrat status in the second repository reported the first repository's run:\n%s", printed)
	}
}

func mustReview(t *testing.T, session *harness.Session, directory string) string {
	t.Helper()

	printed, failures, code := session.ScratIn(directory, "review")
	if code != 0 {
		t.Fatalf("scrat review in %s exited %d: %s", directory, code, failures)
	}
	return printed
}

func TestARepositoryWithNoRuns(t *testing.T) {
	session := harness.Start(t)
	untouched := gitrepo.New(t)

	printed, failures, code := session.ScratIn(untouched.Root, "status")
	if code != 3 {
		t.Errorf("scrat status in a repository that has never been reviewed exited %d, want 3", code)
	}
	if !strings.Contains(failures, untouched.Root) {
		t.Errorf("the refusal does not name the repository: %s", failures)
	}
	if strings.TrimSpace(printed) != "" {
		t.Errorf("scrat status printed a report for a repository with no runs:\n%s", printed)
	}
}

func TestAnIdentifierWinsOverTheDefault(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(cleanScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	session.WaitForTerminal(name)

	elsewhere := gitrepo.New(t)
	printed, failures, code := session.ScratIn(elsewhere.Root, "status", name)
	if code != 0 {
		t.Fatalf("scrat status of a named run exited %d from another repository: %s", code, failures)
	}
	if !strings.Contains(printed, name) {
		t.Errorf("scrat status reported %q, want the run %s", printed, name)
	}
}

func TestCancellingStopsEverything(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(slowScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	session.WaitFor(name, runstate.Reviewing)
	started := session.Recorded(name)

	printed, failures, code := session.Scrat("cancel", name)
	if code != 0 {
		t.Fatalf("scrat cancel exited %d: %s", code, failures)
	}
	if !strings.Contains(printed, name) {
		t.Errorf("scrat cancel does not name the run it stopped:\n%s", printed)
	}

	run := session.Recorded(name)
	if run.State != runstate.Cancelled {
		t.Fatalf("the run is %s after being cancelled", run.State)
	}
	if _, err := os.Stat(started.Worktree); !os.IsNotExist(err) {
		t.Errorf("the worktree at %s survived the cancellation", started.Worktree)
	}
	waitForDeath(t, name, started.Reviewer, started.Agent)
}

func TestCancellingATerminalRunChangesNothing(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(cleanScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	before := session.WaitForTerminal(name)

	printed, failures, code := session.Scrat("cancel", name)
	if code != 0 {
		t.Fatalf("scrat cancel of a finished run exited %d: %s", code, failures)
	}
	if !strings.Contains(printed, string(before.State)) {
		t.Errorf("scrat cancel does not say the run had already finished as %s:\n%s", before.State, printed)
	}

	after := session.Recorded(name)
	if after.State != before.State {
		t.Errorf("cancelling a %s run made it %s", before.State, after.State)
	}
}

func TestNoConfigurationFileIsCreated(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(cleanScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	run := session.WaitForTerminal(name)
	if run.State != runstate.Clean {
		t.Fatalf("the run ended as %s with no configuration file present: %s", run.State, reason(run))
	}

	path := configurationPath(t)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("socrates created %s", path)
	}
}

func TestAnUnreadableConfigurationFileIsARefusalFromTheCommand(t *testing.T) {
	session := harness.Start(t)
	path := configurationPath(t)

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("creating the configuration directory: %v", err)
	}
	if err := os.WriteFile(path, []byte("roles: [not a mapping\n"), 0o644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}

	_, failures, code := session.Scrat("review")
	if code != 3 {
		t.Errorf("scrat review with a configuration file that does not parse exited %d, want 3", code)
	}
	if !strings.Contains(failures, path) {
		t.Errorf("the refusal does not name the file: %s", failures)
	}
	if names := session.Names(); len(names) != 0 {
		t.Errorf("a run was started anyway: %v", names)
	}
}

func TestTheReviewerModelIsWhatTheConfigurationSays(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(cleanScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	path := configurationPath(t)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("creating the configuration directory: %v", err)
	}
	if err := os.WriteFile(path, []byte("roles:\n  reviewer:\n    model: sonnet\n    effort: low\n"), 0o644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}

	name := session.Review()
	session.WaitForTerminal(name)

	invocation := session.Simulator.OnlyInvocation()
	if invocation.Model != "sonnet" {
		t.Errorf("the reviewer was invoked with the model %q, want the configured one", invocation.Model)
	}
	if invocation.Effort != "low" {
		t.Errorf("the reviewer was invoked with the effort %q, want the configured one", invocation.Effort)
	}
}

func TestTheRepositoryIsNeverWrittenTo(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(escalatingScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	before := session.Repository.Status()
	tracked := workingTree(t, session.Repository.Root)

	name := session.Review()
	run := session.WaitForTerminal(name)
	if run.State != runstate.Escalated {
		t.Fatalf("the run ended as %s: %s", run.State, reason(run))
	}

	if after := session.Repository.Status(); after != before {
		t.Errorf("git status changed:\nbefore:\n%s\nafter:\n%s", before, after)
	}
	if after := workingTree(t, session.Repository.Root); !maps.Equal(after, tracked) {
		t.Errorf("the working tree changed:\nbefore: %v\nafter: %v", tracked, after)
	}
	for path := range workingTree(t, session.Repository.Root) {
		if strings.Contains(strings.ToLower(path), "socrates") || strings.Contains(path, "scrat") {
			t.Errorf("socrates left %s inside the repository", path)
		}
	}
}

func workingTree(t *testing.T, root string) map[string]string {
	t.Helper()

	found := map[string]string{}
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if entry.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		relative, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return relErr
		}
		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		found[relative] = string(content)
		return nil
	})
	if err != nil {
		t.Fatalf("reading the working tree at %s: %v", root, err)
	}
	return found
}

func configurationPath(t *testing.T) string {
	t.Helper()

	home := os.Getenv("XDG_CONFIG_HOME")
	if home == "" {
		t.Fatal("the harness left no XDG_CONFIG_HOME, so this test would read the person's own configuration")
	}
	return filepath.Join(home, "socrates", "config.yaml")
}
