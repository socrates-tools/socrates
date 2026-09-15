package simulator

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/socrates-tools/socrates/internal/testing/source"
	"github.com/socrates-tools/socrates/socrates-review/agent"
	"github.com/socrates-tools/socrates/socrates-review/simulation"
)

type Simulator struct {
	t            *testing.T
	Vendor       agent.Vendor
	Directory    string
	Binary       string
	Link         string
	LogPath      string
	ScenarioPath string
}

func SimulatorFor(vendor agent.Vendor) string {
	return vendor.Name + "sim"
}

func Install(t *testing.T) *Simulator {
	t.Helper()
	return InstallFor(t, agent.TheDefault())
}

func InstallFor(t *testing.T, vendor agent.Vendor) *Simulator {
	t.Helper()

	directory := t.TempDir()
	named := SimulatorFor(vendor)
	binary := filepath.Join(directory, named)

	source.Build(t, "./socrates-review/cmd/"+named, binary)

	invokedAs := vendor.Binary
	link := filepath.Join(directory, invokedAs)
	if err := os.Symlink(binary, link); err != nil {
		t.Fatalf("linking the simulator as %s: %v", invokedAs, err)
	}

	t.Setenv("PATH", directory+string(os.PathListSeparator)+os.Getenv("PATH"))

	log := filepath.Join(t.TempDir(), "invocations.jsonl")
	t.Setenv(simulation.LogVariable, log)

	return &Simulator{
		t:         t,
		Vendor:    vendor,
		Directory: directory,
		Binary:    binary,
		Link:      link,
		LogPath:   log,
	}
}

func (s *Simulator) Scenario(scenario string) {
	s.t.Helper()

	path := filepath.Join(s.Directory, "scenario.yaml")
	if err := os.WriteFile(path, []byte(scenario), 0o644); err != nil {
		s.t.Fatalf("writing the scenario: %v", err)
	}
	s.t.Setenv(simulation.ScenarioVariable, path)
	s.ScenarioPath = path
}

func (s *Simulator) Record(names ...string) {
	s.t.Helper()

	joined := ""
	for position, name := range names {
		if position > 0 {
			joined += ","
		}
		joined += name
	}
	s.t.Setenv(simulation.RecordEnvVariable, joined)
}

func (s *Simulator) Invocations() []simulation.Invocation {
	s.t.Helper()

	if _, err := os.Stat(s.LogPath); os.IsNotExist(err) {
		return nil
	}

	invocations, err := simulation.ReadLog(s.LogPath)
	if err != nil {
		s.t.Fatalf("reading the invocation log: %v", err)
	}
	return invocations
}

func (s *Simulator) Processes() []int {
	invocations, err := simulation.ReadLog(s.LogPath)
	if err != nil {
		return nil
	}

	var processes []int
	for _, invocation := range invocations {
		processes = append(processes, invocation.Process)
	}
	return processes
}

func (s *Simulator) WaitForChild(within time.Duration) (simulation.Invocation, bool) {
	deadline := time.Now().Add(within)
	for time.Now().Before(deadline) {
		invocations, err := simulation.ReadLog(s.LogPath)
		if err == nil {
			for _, invocation := range invocations {
				if invocation.Child != 0 {
					return invocation, true
				}
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	return simulation.Invocation{}, false
}

func (s *Simulator) OnlyInvocation() simulation.Invocation {
	s.t.Helper()

	invocations := s.Invocations()
	if len(invocations) != 1 {
		s.t.Fatalf("the simulator was invoked %d times, want once", len(invocations))
	}
	return invocations[0]
}
