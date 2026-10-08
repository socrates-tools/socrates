package simulator

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/socrates-tools/socrates/internal/testing/source"
	"github.com/socrates-tools/socrates/review/agent"
	"github.com/socrates-tools/socrates/review/simulation"
)

type Simulator struct {
	t            *testing.T
	Vendor       agent.Vendor
	Directory    string
	Binary       string
	Link         string
	LogPath      string
	ScenarioPath string
	isolated     bool
	variables    map[string]string
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

	built := filepath.Join(t.TempDir(), SimulatorFor(vendor))
	source.Build(t, "./review/cmd/"+SimulatorFor(vendor), built)

	installed := linked(t, vendor, built)
	installed.isolated = false
	t.Setenv("PATH", installed.Directory+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv(simulation.LogVariable, installed.LogPath)
	return installed
}

func Isolated(t *testing.T, vendor agent.Vendor) *Simulator {
	t.Helper()
	return linked(t, vendor, source.Shared(t, "./review/cmd/"+SimulatorFor(vendor)))
}

func linked(t *testing.T, vendor agent.Vendor, built string) *Simulator {
	t.Helper()

	directory := t.TempDir()
	named := SimulatorFor(vendor)
	binary := filepath.Join(directory, named)
	if err := os.Symlink(built, binary); err != nil {
		t.Fatalf("linking the simulator as %s: %v", named, err)
	}

	invokedAs := vendor.Binary
	link := filepath.Join(directory, invokedAs)
	if err := os.Symlink(binary, link); err != nil {
		t.Fatalf("linking the simulator as %s: %v", invokedAs, err)
	}

	log := filepath.Join(t.TempDir(), "invocations.jsonl")

	return &Simulator{
		t:         t,
		Vendor:    vendor,
		Directory: directory,
		Binary:    binary,
		Link:      link,
		LogPath:   log,
		isolated:  true,
		variables: map[string]string{simulation.LogVariable: log},
	}
}

func (s *Simulator) Environment() []string {
	environment := []string{"PATH=" + s.Directory + string(os.PathListSeparator) + os.Getenv("PATH")}
	for name, value := range s.variables {
		environment = append(environment, name+"="+value)
	}
	return environment
}

func (s *Simulator) set(name, value string) {
	s.t.Helper()

	s.variables[name] = value
	if !s.isolated {
		s.t.Setenv(name, value)
	}
}

func (s *Simulator) Scenario(scenario string) {
	s.t.Helper()

	path := filepath.Join(s.Directory, "scenario.yaml")
	if err := os.WriteFile(path, []byte(scenario), 0o644); err != nil {
		s.t.Fatalf("writing the scenario: %v", err)
	}
	s.set(simulation.ScenarioVariable, path)
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
	s.set(simulation.RecordEnvVariable, joined)
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
