//go:build drift

package fidelity_test

import (
	"bytes"
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/socrates-tools/socrates/internal/testing/simulator"
	"github.com/socrates-tools/socrates/internal/testing/source"
	"github.com/socrates-tools/socrates/socrates-review/agent"
	"github.com/socrates-tools/socrates/socrates-review/fidelity"
	"github.com/socrates-tools/socrates/socrates-review/simulation"
)

func TestTheSimulatorHasNotDriftedFromTheRealCLI(t *testing.T) {
	for _, vendor := range agent.Supported() {
		t.Run(vendor.Name, func(t *testing.T) {
			recordings, err := fidelity.Recordings(fidelity.RecordingsFor(vendor.Name))
			if err != nil {
				t.Fatalf("%v", err)
			}

			named := simulator.SimulatorFor(vendor)
			built := filepath.Join(t.TempDir(), named)
			source.Build(t, "./socrates-review/cmd/"+named, built)

			for _, request := range fidelity.Requests() {
				t.Run(request.Name, func(t *testing.T) {
					recorded, found := recordingNamed(recordings, request.Name)
					if !found {
						t.Fatalf("there is no recording named %s for %s: run make fixtures against the real CLI first", request.Name, vendor.Name)
					}

					produced, err := fidelity.ShapesOf(simulate(t, vendor, built, request))
					if err != nil {
						t.Fatalf("reading what the simulator produced: %v", err)
					}

					for _, difference := range fidelity.Compare(recorded.Shapes, produced) {
						t.Errorf("%s, so either %s changed or its simulator invented it (%s)", difference, vendor.Name, recorded.Path)
					}
				})
			}
		})
	}
}

func simulate(t *testing.T, vendor agent.Vendor, built string, request fidelity.Request) string {
	t.Helper()

	t.Setenv(simulation.LogVariable, filepath.Join(t.TempDir(), "invocations.jsonl"))

	asking, err := agent.New(vendor.Name, agent.Options{Binary: built, Deadline: time.Minute})
	if err != nil {
		t.Fatalf("constructing the simulated agent for %s: %v", vendor.Name, err)
	}

	var stream bytes.Buffer
	if _, err := asking.Run(context.Background(), agent.Request{
		Role:       agent.Reviewer,
		Prompt:     fidelity.Asked(vendor, request),
		Dir:        t.TempDir(),
		Scratch:    t.TempDir(),
		Schema:     request.Schema,
		Transcript: &stream,
	}); err != nil {
		t.Fatalf("the simulator for %s answered: %v", vendor.Name, err)
	}
	return stream.String()
}
