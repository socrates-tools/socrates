package simulator

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rui-ktei/socrates/internal/claudesim"
	"github.com/rui-ktei/socrates/internal/testing/source"
)

const invokedAs = "claude"

type Simulator struct {
	t            *testing.T
	Directory    string
	Binary       string
	Link         string
	LogPath      string
	ScenarioPath string
}

func Install(t *testing.T) *Simulator {
	t.Helper()

	directory := t.TempDir()
	binary := filepath.Join(directory, "claudesim")

	source.Build(t, "./cmd/claudesim", binary)

	link := filepath.Join(directory, invokedAs)
	if err := os.Symlink(binary, link); err != nil {
		t.Fatalf("linking the simulator as %s: %v", invokedAs, err)
	}

	t.Setenv("PATH", directory+string(os.PathListSeparator)+os.Getenv("PATH"))

	log := filepath.Join(t.TempDir(), "invocations.jsonl")
	t.Setenv(claudesim.LogVariable, log)

	return &Simulator{
		t:         t,
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
	s.t.Setenv(claudesim.ScenarioVariable, path)
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
	s.t.Setenv(claudesim.RecordEnvVariable, joined)
}

func (s *Simulator) Invocations() []claudesim.Invocation {
	s.t.Helper()

	if _, err := os.Stat(s.LogPath); os.IsNotExist(err) {
		return nil
	}

	invocations, err := claudesim.ReadLog(s.LogPath)
	if err != nil {
		s.t.Fatalf("reading the invocation log: %v", err)
	}
	return invocations
}

func (s *Simulator) Processes() []int {
	invocations, err := claudesim.ReadLog(s.LogPath)
	if err != nil {
		return nil
	}

	var processes []int
	for _, invocation := range invocations {
		processes = append(processes, invocation.Process)
	}
	return processes
}

func (s *Simulator) WaitForChild(within time.Duration) (claudesim.Invocation, bool) {
	deadline := time.Now().Add(within)
	for time.Now().Before(deadline) {
		invocations, err := claudesim.ReadLog(s.LogPath)
		if err == nil {
			for _, invocation := range invocations {
				if invocation.Child != 0 {
					return invocation, true
				}
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	return claudesim.Invocation{}, false
}

func (s *Simulator) OnlyInvocation() claudesim.Invocation {
	s.t.Helper()

	invocations := s.Invocations()
	if len(invocations) != 1 {
		s.t.Fatalf("the simulator was invoked %d times, want once", len(invocations))
	}
	return invocations[0]
}
