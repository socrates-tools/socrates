package e2e_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/agora/surface"
	"github.com/socrates-tools/socrates/internal/exit"
	"github.com/socrates-tools/socrates/internal/state"
	"github.com/socrates-tools/socrates/internal/testing/harness"
	"github.com/socrates-tools/socrates/review/runstate"
)

func TestOneVariableMovesEverything(t *testing.T) {
	t.Parallel()
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

	name := session.Begin()
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
		if _, err := os.Stat(filepath.Join(session.UserHome, elsewhere)); err == nil {
			t.Errorf("socrates wrote %s under the person's own home, and one variable states where it writes", elsewhere)
		}
	}
}

func TestARelativeHomeIsRefused(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	session.Repository.Write("internal/added.go", "package internal\n")

	relative := filepath.Join("somewhere", "below")
	out, failures, code := session.SocratesWith(session.Repository.Root, []string{state.HomeVariable + "=" + relative}, "begin")

	if code != exit.Unusable {
		t.Fatalf("socrates review begin with a relative home exited %d, want the usage code %d: %s%s", code, exit.Unusable, out, failures)
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
	t.Parallel()
	session := harness.Start(t)
	session.Simulator.Scenario(cleanScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Begin()
	session.WaitForTerminal(name)

	if _, err := os.Stat(filepath.Join(session.Home, "runs", name)); err != nil {
		t.Fatalf("the run is not under socrates' home: %v", err)
	}

	beside := filepath.Join(session.Home, "config.yaml")
	if err := os.WriteFile(beside, []byte("roles: [not a mapping\n"), 0o644); err != nil {
		t.Fatalf("writing the configuration beside the runs: %v", err)
	}

	_, failures, code := session.Socrates("begin")
	if code != exit.Unusable {
		t.Fatalf("socrates review begin exited %d with a configuration file beside the runs that does not parse, want %d", code, exit.Unusable)
	}
	if !strings.Contains(failures, beside) {
		t.Errorf("the refusal reads %q, and the configuration socrates reads lives beside the runs at %s", failures, beside)
	}
}

func TestOneVariableMovesEverythingBothToolsHave(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	session.Simulator.Scenario(cleanScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	if err := os.WriteFile(filepath.Join(session.Home, state.ConfigurationFile), []byte("roles:\n  reviewer:\n    model: sonnet\n"), 0o644); err != nil {
		t.Fatalf("writing review's configuration file under the stated root: %v", err)
	}

	held := servingSurface(t)
	held.home = session.Root
	held.where = session.Repository.Root
	stated := aFreeLoopbackAddress(t)
	held.configure("bind: " + stated + "\n")
	held.bindsWhatIsConfigured()

	name := session.Begin()
	session.WaitForTerminal(name)

	printed := held.serve(oneDecisionDocument)
	if !strings.Contains(printed, stated) {
		t.Errorf("agora serve printed %q, and agora's configuration file under the stated root states the bind %s", printed, stated)
	}
	if model := session.Simulator.OnlyInvocation().Model; model != "sonnet" {
		t.Errorf("the reviewer was invoked with the model %q, so review's configuration file under %s was not read", model, session.Home)
	}

	held.answer(surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}, Words: "go on"}}})

	root := whatIsUnder(t, session.Root)
	if !root[string(state.Review)] || !root[string(state.Agora)] {
		t.Errorf("the stated root holds %v, and each tool keeps a directory of its own under it", root)
	}
	review := whatIsUnder(t, session.Home)
	if !review["runs"] || !review[state.ConfigurationFile] {
		t.Errorf("review's directory holds %v, and the runs and review's configuration file live under it", review)
	}
	if review["surface"] || review["surface.json"] || review["sessions"] {
		t.Errorf("review's directory holds %v, and nothing the surface writes belongs in it", review)
	}
	agora := whatIsUnder(t, state.Agora.HomeUnder(session.Root))
	if !agora[state.ConfigurationFile] || !agora["surface.json"] || !agora["sessions"] {
		t.Errorf("agora's directory holds %v, and agora's configuration file and what the surface writes live under it", agora)
	}
	for _, elsewhere := range []string{".socrates", ".agora", filepath.Join(".config", "agora"), filepath.Join(".local", "state", "agora")} {
		if _, err := os.Stat(filepath.Join(session.UserHome, elsewhere)); err == nil {
			t.Errorf("%s was written under the person's own home, and one variable states where both tools write", elsewhere)
		}
	}
}

func TestTheSurfaceWritesInAgorasDirectoryAndNothingInReviews(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	session.Simulator.Scenario(cleanScenario)
	session.Repository.Write("internal/added.go", "package internal\n")

	held := servingSurface(t)
	held.home = session.Root
	held.where = session.Repository.Root

	name := session.Begin()
	session.WaitForTerminal(name)
	held.serve(oneDecisionDocument)

	own := whatIsUnder(t, held.surfaceHome())
	if !own["surface.json"] || !own["sessions"] {
		t.Fatalf("agora's directory holds %v, without what the surface wrote about itself and the session it opened", own)
	}

	beside := whatIsUnder(t, session.Home)
	if !beside["runs"] || beside["surface"] || beside["surface.json"] || beside["sessions"] {
		t.Fatalf("review's directory holds %v, and nothing the surface writes belongs beside the runs", beside)
	}

	runs := filepath.Join(session.Home, "runs")
	if err := filepath.WalkDir(runs, func(path string, entry fs.DirEntry, walked error) error {
		if walked != nil {
			return walked
		}
		if entry.Name() == "surface.json" || entry.Name() == "sessions" || entry.Name() == "session.json" {
			t.Errorf("a run directory holds %s, and what the surface writes is beside the runs rather than inside one", path)
		}
		return nil
	}); err != nil {
		t.Fatalf("reading what is under %s: %v", runs, err)
	}
}

func TestTheVariableThatMovedTheSurfacesOwnHomeIsRefused(t *testing.T) {
	t.Parallel()
	held := standing(t)
	moved := t.TempDir()
	held.also = []string{surface.RetiredHomeVariable + "=" + moved}

	printed, failures, code := held.reading(oneDecisionDocument, "serve", "--document", "-", "--address", aPortTheOperatingSystemChooses)
	if code == 0 {
		t.Fatalf("a surface was served at %q while %s was still set", printed, surface.RetiredHomeVariable)
	}
	if !strings.Contains(failures, surface.RetiredHomeVariable) || !strings.Contains(failures, state.HomeVariable) {
		t.Errorf("the refusal reads %q, and it names the variable that moved the surface's home and the one that replaced it", failures)
	}
	if instruction := "set " + state.HomeVariable + " instead"; !strings.Contains(failures, instruction) {
		t.Errorf("the refusal reads %q, and it tells the person to %q", failures, instruction)
	}
	if !strings.Contains(failures, "root") {
		t.Errorf("the refusal reads %q, and it names the brand's root %s moves", failures, state.HomeVariable)
	}
	if strings.Contains(failures, "both binaries") {
		t.Errorf("the refusal reads %q, and the suite is one binary", failures)
	}
	if _, err := os.Stat(filepath.Join(held.surfaceHome(), "surface.json")); err == nil {
		t.Error("a refused startup left a surface behind, and it starts no server")
	}
	if _, err := os.Stat(moved); err == nil {
		if entries := whatIsUnder(t, moved); len(entries) != 0 {
			t.Errorf("the home that variable states holds %v, and it is read by nothing", entries)
		}
	}
}

func TestARelativeHomeIsRefusedByTheSurfaceToo(t *testing.T) {
	// scenario: A relative home is refused
	t.Parallel()
	held := standing(t)
	relative := filepath.Join("somewhere", "below")
	held.home = relative

	printed, failures, code := held.reading(oneDecisionDocument, "serve", "--document", "-", "--address", aPortTheOperatingSystemChooses)
	if code == 0 {
		t.Fatalf("a surface was served at %q with a relative home", printed)
	}
	if !strings.Contains(failures, state.HomeVariable) || !strings.Contains(failures, relative) {
		t.Errorf("the refusal reads %q, without naming the variable and the path", failures)
	}
	if _, err := os.Stat(filepath.Join(held.where, relative)); err == nil {
		t.Error("the surface read the relative home as a path below where it was run")
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
