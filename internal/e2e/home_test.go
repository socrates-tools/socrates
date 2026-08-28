package e2e_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rui-ktei/socrates/internal/exit"
	"github.com/rui-ktei/socrates/internal/runstate"
	"github.com/rui-ktei/socrates/internal/state"
	"github.com/rui-ktei/socrates/internal/surface"
	"github.com/rui-ktei/socrates/internal/testing/harness"
)

func TestOneVariableMovesEverything(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(cleanScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	stated := session.Home
	if err := os.MkdirAll(stated, 0o755); err != nil {
		t.Fatalf("creating the home the variable states: %v", err)
	}
	if err := os.WriteFile(filepath.Join(stated, "config.yaml"), []byte("roles:\n  reviewer:\n    model: sonnet\n"), 0o644); err != nil {
		t.Fatalf("writing the configuration under the stated home: %v", err)
	}

	name := session.Review()
	run := session.WaitForTerminal(name)
	if run.State != runstate.Clean {
		t.Fatalf("the run ended as %s: %s", run.State, reason(run))
	}

	if model := session.Simulator.OnlyInvocation().Model; model != "sonnet" {
		t.Errorf("the reviewer was invoked with the model %q, so the configuration under %s was not read", model, stated)
	}
	if _, err := os.Stat(filepath.Join(stated, "runs", name)); err != nil {
		t.Errorf("the run is not under the stated home: %v", err)
	}
	for _, elsewhere := range []string{".socrates", filepath.Join(".config", "socrates"), filepath.Join(".local", "state", "socrates")} {
		if _, err := os.Stat(filepath.Join(os.Getenv("HOME"), elsewhere)); err == nil {
			t.Errorf("socrates wrote %s under the person's own home, and one variable states where it writes", elsewhere)
		}
	}
}

func TestARelativeHomeIsRefused(t *testing.T) {
	session := harness.Start(t)
	session.Repository.Write("internal/added.go", "package internal\n")

	relative := filepath.Join("somewhere", "below")
	out, failures, code := session.SocratesWith(session.Repository.Root, []string{state.HomeVariable + "=" + relative}, "review")

	if code != exit.Usage {
		t.Fatalf("socrates review with a relative home exited %d, want the usage code %d: %s%s", code, exit.Usage, out, failures)
	}
	if !strings.Contains(failures, state.HomeVariable) || !strings.Contains(failures, relative) {
		t.Errorf("the refusal reads %q, without naming the variable and the path", failures)
	}
	if _, err := os.Stat(filepath.Join(session.Repository.Root, relative)); err == nil {
		t.Errorf("socrates read the relative home as a path below the repository")
	}
	if names := session.Names(); len(names) != 0 {
		t.Errorf("a run was written to the default home anyway: %v", names)
	}
}

func TestARunIsFoundBesideTheConfiguration(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(cleanScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	session.WaitForTerminal(name)

	if _, err := os.Stat(filepath.Join(session.Home, "runs", name)); err != nil {
		t.Fatalf("the run is not under socrates' home: %v", err)
	}

	beside := filepath.Join(session.Home, "config.yaml")
	if err := os.WriteFile(beside, []byte("roles: [not a mapping\n"), 0o644); err != nil {
		t.Fatalf("writing the configuration beside the runs: %v", err)
	}

	_, failures, code := session.Socrates("review")
	if code != exit.Usage {
		t.Fatalf("socrates review exited %d with a configuration file beside the runs that does not parse, want %d", code, exit.Usage)
	}
	if !strings.Contains(failures, beside) {
		t.Errorf("the refusal reads %q, and the configuration socrates reads lives beside the runs at %s", failures, beside)
	}
}

func TestOneVariableMovesEverythingAgoraHas(t *testing.T) {
	held := standing(t)
	stated := aFreeLoopbackAddress(t)
	held.configure("bind: " + stated + "\n")

	printed := held.serve(oneDecisionDocument)
	if !strings.Contains(printed, stated) {
		t.Errorf("agora serve printed %q, and the configuration under its home states the bind %s", printed, stated)
	}

	held.answer(surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}, Words: "go on"}}})

	own := whatIsUnder(t, held.home)
	if !own["config.yaml"] || !own["session.json"] {
		t.Errorf("agora's home holds %v, and the configuration it read and the session it wrote both live there", own)
	}
	for _, elsewhere := range []string{".agora", filepath.Join(".config", "agora"), filepath.Join(".local", "state", "agora")} {
		if _, err := os.Stat(filepath.Join(os.Getenv("HOME"), elsewhere)); err == nil {
			t.Errorf("agora wrote %s under the person's own home, and one variable states where it writes", elsewhere)
		}
	}
}

func TestNeitherProgramWritesInTheOthersHome(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(cleanScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	held := servingSurface(t)
	held.where = session.Repository.Root

	name := session.Review()
	session.WaitForTerminal(name)
	held.serve(oneDecisionDocument)

	socrates := whatIsUnder(t, session.Home)
	agora := whatIsUnder(t, held.home)

	if !socrates["runs"] {
		t.Fatalf("socrates' home holds %v, without the runs it writes", socrates)
	}
	if !agora["session.json"] {
		t.Fatalf("agora's home holds %v, without the session it writes", agora)
	}
	if socrates["session.json"] || socrates["document.json"] || socrates["config.yaml"] {
		t.Errorf("socrates' home holds %v, and what agora writes is agora's", socrates)
	}
	if agora["runs"] {
		t.Errorf("agora's home holds %v, and what socrates writes is socrates'", agora)
	}
}

func whatIsUnder(t *testing.T, directory string) map[string]bool {
	t.Helper()

	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatalf("reading %s: %v", directory, err)
	}

	found := map[string]bool{}
	for _, entry := range entries {
		found[entry.Name()] = true
	}
	return found
}
