//go:build drift

package fidelity_test

import (
	"bytes"
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/rui-ktei/socrates/internal/agent"
	"github.com/rui-ktei/socrates/internal/fidelity"
	"github.com/rui-ktei/socrates/internal/testing/source"
)

func TestTheSimulatorHasNotDriftedFromTheRealCLI(t *testing.T) {
	recordings, err := fidelity.Recordings(fidelity.RecordingsDirectory)
	if err != nil {
		t.Fatalf("%v", err)
	}

	simulator := filepath.Join(t.TempDir(), "claudesim")
	source.Build(t, "./cmd/claudesim", simulator)

	for _, request := range fidelity.Requests() {
		t.Run(request.Name, func(t *testing.T) {
			recorded, found := recordingNamed(recordings, request.Name)
			if !found {
				t.Fatalf("there is no recording named %s: run make fixtures against the real CLI first", request.Name)
			}

			produced, err := fidelity.ShapesOf(simulate(t, simulator, request))
			if err != nil {
				t.Fatalf("reading what the simulator produced: %v", err)
			}

			for _, difference := range fidelity.Compare(recorded.Shapes, produced) {
				t.Errorf("%s, so either the CLI changed or the simulator invented it (%s)", difference, recorded.Path)
			}
		})
	}
}

func simulate(t *testing.T, simulator string, request fidelity.Request) string {
	t.Helper()

	asking, err := agent.NewClaude(agent.Options{Binary: simulator, Deadline: time.Minute})
	if err != nil {
		t.Fatalf("constructing the simulated agent: %v", err)
	}

	var stream bytes.Buffer
	if _, err := asking.Run(context.Background(), agent.Request{
		Role:       agent.Reviewer,
		Prompt:     request.Prompt,
		Dir:        t.TempDir(),
		Schema:     request.Schema,
		Transcript: &stream,
	}); err != nil {
		t.Fatalf("the simulator answered: %v", err)
	}
	return stream.String()
}

func recordingNamed(recordings []fidelity.Recording, name string) (fidelity.Recording, bool) {
	for _, recording := range recordings {
		if recording.Name == name {
			return recording, true
		}
	}
	return fidelity.Recording{}, false
}
