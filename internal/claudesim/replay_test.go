package claudesim_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/corpus"
	"github.com/socrates-tools/socrates/internal/simulation"
)

const replayed = "dialogue-layer-first-look"

func TestAScenarioNamesARecordedAnswer(t *testing.T) {
	entries, err := corpus.Entries()
	if err != nil {
		t.Fatalf("reading the corpus: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("the corpus holds no recorded answer, so nothing here proves anything")
	}

	for _, entry := range entries {
		t.Run(entry.Name, func(t *testing.T) {
			scenario(t, "default:\n  recorded: "+entry.Name+"\n")

			out, failures, code := simulate(t, validArguments, "review the change")
			if code != 0 {
				t.Fatalf("the simulator exited %d: %s", code, failures)
			}
			if answered := structuredOutput(t, out); !bytes.Equal(answered, entry.Document) {
				t.Errorf("the structured result is not the recorded document, byte for byte: %s", difference(entry.Document, answered))
			}
		})
	}
}

func TestANamedAnswerThatDoesNotExistIsRefused(t *testing.T) {
	scenario(t, "default:\n  recorded: no-live-reviewer-ever-said-this\n")

	out, failures, code := simulate(t, validArguments, "review the change")
	if code == 0 {
		t.Fatalf("a scenario naming an answer the corpus does not hold was answered with %q", out)
	}
	if !strings.Contains(failures, "no-live-reviewer-ever-said-this") {
		t.Errorf("the refusal does not name the entry: %s", failures)
	}
	if strings.Contains(out, `"type":"result"`) {
		t.Errorf("the review was answered with a substitute: %s", out)
	}
}

func TestTheSimulatorDoesNotRepairARecordedAnswer(t *testing.T) {
	scenario(t, "default:\n  recorded: "+replayed+"\n")

	entry, err := corpus.Named(replayed)
	if err != nil {
		t.Fatalf("reading the recorded answer: %v", err)
	}

	elsewhere := t.TempDir()
	working, err := os.Getwd()
	if err != nil {
		t.Fatalf("locating the working directory: %v", err)
	}
	if err := os.Chdir(elsewhere); err != nil {
		t.Fatalf("moving to %s: %v", elsewhere, err)
	}
	t.Cleanup(func() { _ = os.Chdir(working) })

	first, failures, code := simulate(t, validArguments, "review the change against a9a3207, and report on internal/dialogue/dialogue.go")
	if code != 0 {
		t.Fatalf("the simulator exited %d: %s", code, failures)
	}
	second, failures, code := simulate(t, validArguments, "a prompt that names nothing the recorded answer mentions")
	if code != 0 {
		t.Fatalf("the simulator exited %d: %s", code, failures)
	}

	answered := structuredOutput(t, first)
	if !bytes.Equal(answered, structuredOutput(t, second)) {
		t.Error("the recorded answer was adapted to the prompt it received")
	}
	if !bytes.Equal(answered, entry.Document) {
		t.Errorf("the recorded answer was repaired on its way out: %s", difference(entry.Document, answered))
	}

	var document struct {
		Findings []struct {
			ID   string `json:"id"`
			File string `json:"file"`
		} `json:"findings"`
	}
	if err := json.Unmarshal(answered, &document); err != nil {
		t.Fatalf("reading the answer the simulator emitted: %v", err)
	}
	if len(document.Findings) == 0 {
		t.Fatal("the recorded answer holds no finding, so nothing here proves anything")
	}
	for _, finding := range document.Findings {
		if _, err := os.Stat(filepath.Join(elsewhere, finding.File)); !os.IsNotExist(err) {
			t.Fatalf("%s exists beside the simulator, so this proves nothing about adaptation", finding.File)
		}
		if finding.ID == "" {
			t.Error("a finding lost its identifier")
		}
	}
}

func scenario(t *testing.T, stated string) {
	t.Helper()

	path := filepath.Join(t.TempDir(), "scenario.yaml")
	if err := os.WriteFile(path, []byte(stated), 0o644); err != nil {
		t.Fatalf("writing the scenario: %v", err)
	}
	t.Setenv(simulation.ScenarioVariable, path)
}

func difference(recorded, answered []byte) string {
	at := 0
	for at < len(recorded) && at < len(answered) && recorded[at] == answered[at] {
		at++
	}

	from := max(at-40, 0)
	return fmt.Sprintf("%d bytes recorded, %d answered, differing at byte %d:\nrecorded %q\nanswered %q",
		len(recorded), len(answered), at, window(recorded, from), window(answered, from))
}

func window(document []byte, from int) string {
	if from >= len(document) {
		return ""
	}
	return string(document[from:min(from+120, len(document))])
}

func structuredOutput(t *testing.T, out string) []byte {
	t.Helper()

	var last []byte
	for line := range strings.SplitSeq(strings.TrimSuffix(out, "\n"), "\n") {
		if line == "" {
			continue
		}
		var event struct {
			Type             string          `json:"type"`
			StructuredOutput json.RawMessage `json:"structured_output"`
		}
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			t.Fatalf("a line of the stream does not parse: %v: %s", err, line)
		}
		if event.Type == "result" {
			last = event.StructuredOutput
		}
	}

	if last == nil {
		t.Fatalf("the stream carries no structured result:\n%s", out)
	}
	return last
}
