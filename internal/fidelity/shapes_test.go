package fidelity_test

import (
	"strings"
	"testing"

	"github.com/ktei/socrates/internal/fidelity"
)

const recordedStream = `{"type":"system","subtype":"init","session_id":"s","model":"claude","tools":["Read"],"cwd":"/tmp"}
{"type":"assistant","session_id":"s","message":{"id":"m","type":"message","role":"assistant","model":"claude","content":[{"type":"text","text":"hello"}],"usage":{"input_tokens":1,"output_tokens":2,"cache_read_input_tokens":0,"cache_creation_input_tokens":0}}}
{"type":"result","subtype":"success","session_id":"s","is_error":false,"result":"hello","structured_output":{"findings":[],"risk_level":"low","risk_rationale":"r"},"usage":{"input_tokens":1,"output_tokens":2,"cache_read_input_tokens":0,"cache_creation_input_tokens":0},"duration_ms":1,"num_turns":1,"total_cost_usd":0.01}
`

func TestDriftIsVisible(t *testing.T) {
	recorded, err := fidelity.ShapesOf(recordedStream)
	if err != nil {
		t.Fatalf("reading the recorded stream: %v", err)
	}

	drifted, err := fidelity.ShapesOf(strings.ReplaceAll(recordedStream, `"structured_output"`, `"structuredOutput"`))
	if err != nil {
		t.Fatalf("reading the drifted stream: %v", err)
	}

	differences := fidelity.Compare(recorded, drifted)
	if len(differences) == 0 {
		t.Fatal("a stream carrying a field the recording does not have compared equal")
	}

	reported := ""
	for _, difference := range differences {
		reported += difference.String() + "\n"
	}
	if !strings.Contains(reported, "structuredOutput") {
		t.Errorf("the drift does not name the field that differs:\n%s", reported)
	}
	if !strings.Contains(reported, "result") {
		t.Errorf("the drift does not name the event the field is in:\n%s", reported)
	}
}

func TestAnEventTypeNobodyRecordedIsDrift(t *testing.T) {
	recorded, err := fidelity.ShapesOf(recordedStream)
	if err != nil {
		t.Fatalf("reading the recorded stream: %v", err)
	}
	invented, err := fidelity.ShapesOf(recordedStream + `{"type":"summary","session_id":"s"}` + "\n")
	if err != nil {
		t.Fatalf("reading the invented stream: %v", err)
	}

	differences := fidelity.Compare(recorded, invented)
	if len(differences) != 1 || differences[0].Type != "summary" {
		t.Errorf("an event type the recording has none of gave %v", differences)
	}
}

func TestTheSameStreamHasNoDrift(t *testing.T) {
	recorded, err := fidelity.ShapesOf(recordedStream)
	if err != nil {
		t.Fatalf("reading the recorded stream: %v", err)
	}
	if differences := fidelity.Compare(recorded, recorded); len(differences) != 0 {
		t.Errorf("a stream compared against itself reported %v", differences)
	}
}

func TestAFieldTheRecordingHasAndTheSimulatorDoesNotIsNotDrift(t *testing.T) {
	recorded, err := fidelity.ShapesOf(recordedStream)
	if err != nil {
		t.Fatalf("reading the recorded stream: %v", err)
	}
	leaner, err := fidelity.ShapesOf(strings.ReplaceAll(recordedStream, `,"total_cost_usd":0.01`, ""))
	if err != nil {
		t.Fatalf("reading the leaner stream: %v", err)
	}

	if differences := fidelity.Compare(recorded, leaner); len(differences) != 0 {
		t.Errorf("a simulator that produces less than the CLI reported %v, and only invention is drift", differences)
	}
}

func TestALineThatIsNotAnEventIsRefused(t *testing.T) {
	if _, err := fidelity.ShapesOf("this is not an event\n"); err == nil {
		t.Fatal("a line that is not JSON was read as an event")
	}
}
