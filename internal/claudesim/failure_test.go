package claudesim_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/socrates-tools/socrates-review/internal/findings"
	"github.com/socrates-tools/socrates-review/internal/simulation"
)

func TestTheSimulatorProducesAFailureOnDemand(t *testing.T) {
	t.Run("a truncated stream", func(t *testing.T) {
		out, _, code := scripted(t, "behaviour: truncated-stream")

		if code != 0 {
			t.Fatalf("the simulator exited %d, and a truncated stream is not an exit code", code)
		}
		if strings.HasSuffix(out, "\n") {
			t.Fatal("the stream ends with a newline, so it is not truncated part-way through a line")
		}

		lines := strings.Split(out, "\n")
		last := lines[len(lines)-1]
		if json.Valid([]byte(last)) {
			t.Errorf("the last line %q parses, so nothing was truncated", last)
		}
		if !strings.Contains(last, `"type":"result"`) {
			t.Errorf("the truncated line is %q, want the result event cut part-way", last)
		}
	})

	t.Run("a missing result event", func(t *testing.T) {
		out, _, code := scripted(t, "behaviour: no-result")

		if code != 0 {
			t.Fatalf("the simulator exited %d", code)
		}
		for _, event := range streamed(t, out) {
			if event["type"] == "result" {
				t.Error("the stream carries a result event")
			}
		}
		if !strings.Contains(out, `"type":"assistant"`) {
			t.Error("the stream carries nothing at all, so it does not isolate the missing result")
		}
	})

	t.Run("no structured output", func(t *testing.T) {
		out, _, code := scripted(t, "behaviour: no-structured-output")

		if code != 0 {
			t.Fatalf("the simulator exited %d", code)
		}
		result := lastEvent(t, out)
		if result["type"] != "result" {
			t.Fatalf("the stream ends with %v, want a result event", result["type"])
		}
		if _, present := result["structured_output"]; present {
			t.Error("the result event carries structured output")
		}
	})

	t.Run("a schema violation", func(t *testing.T) {
		out, _, code := scripted(t, "behaviour: schema-violation")

		if code != 0 {
			t.Fatalf("the simulator exited %d", code)
		}
		structured, err := json.Marshal(lastEvent(t, out)["structured_output"])
		if err != nil {
			t.Fatalf("re-encoding the structured output: %v", err)
		}
		if !json.Valid(structured) {
			t.Fatal("the structured output is not JSON, so it does not isolate a schema violation")
		}
		if _, err := findings.Parse(structured); err == nil {
			t.Errorf("the structured output %s satisfies the schema", structured)
		}
	})

	t.Run("a non-zero exit", func(t *testing.T) {
		out, failures, code := scripted(t, "behaviour: non-zero-exit\nstderr: the provider gave up")

		if code == 0 {
			t.Fatal("the simulator exited zero")
		}
		if !strings.Contains(failures, "the provider gave up") {
			t.Errorf("the error stream reads %q, without what the scenario asked for", failures)
		}

		result := lastEvent(t, out)
		if result["type"] != "result" {
			t.Errorf("the stream ends with %v, and this failure is a complete stream followed by an exit", result["type"])
		}
		if _, present := result["structured_output"]; !present {
			t.Error("the complete stream carries no structured output, so the exit code is not what fails the call")
		}
	})

	t.Run("a result reporting an error", func(t *testing.T) {
		out, failures, code := scripted(t, "behaviour: error-result\ntext: \"API Error: 529 Overloaded\"")

		if code == 0 {
			t.Fatal("the simulator exited zero, and the recorded live shape exits non-zero")
		}
		if failures != "" {
			t.Errorf("the error stream reads %q, and the live shape said everything in the result event", failures)
		}

		result := lastEvent(t, out)
		if result["is_error"] != true {
			t.Errorf("the result event reports is_error %v, want the error the account belongs to", result["is_error"])
		}
		if result["subtype"] != "success" {
			t.Errorf("the result event carries the subtype %v, and the live shape carried success beside the error", result["subtype"])
		}
		if result["result"] != "API Error: 529 Overloaded" {
			t.Errorf("the result event carries the account %v, want what the scenario asked for", result["result"])
		}
		if _, present := result["structured_output"]; present {
			t.Error("the result event carries structured output, and a call that reported an error answered nothing")
		}
	})

	t.Run("a result reporting an error and exiting zero", func(t *testing.T) {
		out, _, code := scripted(t, "behaviour: error-result\nexit: 0")

		if code != 0 {
			t.Fatalf("the simulator exited %d, and a scenario asking for a zero exit drives the branch a non-zero exit masks", code)
		}
		if lastEvent(t, out)["is_error"] != true {
			t.Error("the result event does not report an error, so nothing distinguishes this from an answer")
		}
	})

	t.Run("silence", func(t *testing.T) {
		out, _, code := scripted(t, "behaviour: silence\nhold: 20ms")

		if code != 0 {
			t.Fatalf("the simulator exited %d", code)
		}
		if out != "" {
			t.Errorf("the silent agent wrote %q", out)
		}
	})

	t.Run("a child process that outlives the agent", func(t *testing.T) {
		log := filepath.Join(t.TempDir(), "invocations.jsonl")
		t.Setenv(simulation.LogVariable, log)

		_, _, code := scripted(t, "child_command: [sleep, \"30\"]")
		if code != 0 {
			t.Fatalf("the simulator exited %d", code)
		}

		invocations, err := simulation.ReadLog(log)
		if err != nil {
			t.Fatalf("reading the log: %v", err)
		}
		child := invocations[0].Child
		if child == 0 {
			t.Fatal("the log records no child, so a test cannot look for it")
		}
		t.Cleanup(func() {
			_ = syscall.Kill(child, syscall.SIGKILL)
		})

		if err := syscall.Kill(child, 0); err != nil {
			t.Errorf("the child %d is not running: %v", child, err)
		}
	})

	t.Run("a reviewer that writes", func(t *testing.T) {
		working := t.TempDir()
		t.Chdir(working)

		_, _, code := scripted(t, "writes:\n  - path: internal/helped.go\n    content: \"package internal\\n\"")
		if code != 0 {
			t.Fatalf("the simulator exited %d", code)
		}

		written, err := os.ReadFile(filepath.Join(working, "internal", "helped.go"))
		if err != nil {
			t.Fatalf("reading what the reviewer wrote: %v", err)
		}
		if string(written) != "package internal\n" {
			t.Errorf("the reviewer wrote %q", written)
		}
	})
}

func TestTheScenarioSelectsByPrompt(t *testing.T) {
	scenario := `
default:
  behaviour: findings
  text: the default answer
responses:
  - when: reviewing a change you did not write
    behaviour: non-zero-exit
    stderr: the reviewer prompt was matched
`
	writeScenario(t, scenario)

	_, failures, code := simulate(t, validArguments, "You are reviewing a change you did not write.")
	if code == 0 || !strings.Contains(failures, "the reviewer prompt was matched") {
		t.Errorf("the matching response was not chosen: exit %d, %q", code, failures)
	}

	out, _, code := simulate(t, validArguments, "something else entirely")
	if code != 0 {
		t.Fatalf("the default was not chosen: exit %d", code)
	}
	if !strings.Contains(out, "the default answer") {
		t.Errorf("the stream reads %q, without the default's text", out)
	}
}

func TestAScenarioNobodyCanHonourIsRefused(t *testing.T) {
	for _, scenario := range []string{
		"default:\n  behaviour: be-nice",
		"default:\n  hold: soon",
		"default:\n  behaviour: findings\n  unknown_field: 1",
	} {
		t.Run(scenario, func(t *testing.T) {
			writeScenario(t, scenario)

			if _, failures, code := simulate(t, validArguments, "review the change"); code == 0 {
				t.Errorf("the scenario was accepted, and the simulator answered anyway: %q", failures)
			}
		})
	}
}

func TestNoScenarioIsACleanReview(t *testing.T) {
	t.Setenv(simulation.ScenarioVariable, "")

	out, _, code := simulate(t, validArguments, "review the change")
	if code != 0 {
		t.Fatalf("the simulator exited %d", code)
	}

	structured, err := json.Marshal(lastEvent(t, out)["structured_output"])
	if err != nil {
		t.Fatalf("re-encoding the structured output: %v", err)
	}
	document, err := findings.Parse(structured)
	if err != nil {
		t.Fatalf("the default answer is not a findings document: %v", err)
	}
	if len(document.Findings) != 0 {
		t.Errorf("the default answer holds %d findings, want none", len(document.Findings))
	}
}

func scripted(t *testing.T, response string) (string, string, int) {
	t.Helper()

	indented := "  " + strings.ReplaceAll(strings.TrimSpace(response), "\n", "\n  ")
	writeScenario(t, "default:\n"+indented+"\n")
	return simulate(t, validArguments, "review the change")
}

func writeScenario(t *testing.T, scenario string) {
	t.Helper()

	path := filepath.Join(t.TempDir(), "scenario.yaml")
	if err := os.WriteFile(path, []byte(scenario), 0o644); err != nil {
		t.Fatalf("writing the scenario: %v", err)
	}
	t.Setenv(simulation.ScenarioVariable, path)
}

func lastEvent(t *testing.T, out string) map[string]any {
	t.Helper()

	events := streamed(t, out)
	if len(events) == 0 {
		t.Fatal("the stream is empty")
	}
	return events[len(events)-1]
}
