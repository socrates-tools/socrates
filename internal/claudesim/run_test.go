package claudesim_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ktei/socrates/internal/claudesim"
	"github.com/ktei/socrates/internal/findings"
	"github.com/ktei/socrates/internal/simulation"
)

var validArguments = []string{
	"-p",
	"--verbose",
	"--output-format", "stream-json",
	"--setting-sources", "user",
	"--tools", "Read,Grep,Glob",
	"--model", "opus",
	"--effort", "high",
	"--json-schema", findings.Schema,
	"--dangerously-skip-permissions",
}

func TestTheInvocationContractIsEnforced(t *testing.T) {
	cases := []struct {
		name      string
		arguments []string
		refused   string
	}{
		{name: "what socrates sends", arguments: validArguments},
		{name: "no print mode", arguments: []string{"--output-format", "stream-json"}, refused: "print mode"},
		{name: "an option the real CLI does not have", arguments: []string{"-p", "--nope"}, refused: "unknown option '--nope'"},
		{name: "a format nobody streams", arguments: []string{"-p", "--output-format", "text"}, refused: "stream-json"},
		{name: "a format the real CLI rejects", arguments: []string{"-p", "--output-format", "bogus"}, refused: "invalid"},
		{name: "an effort level the real CLI rejects", arguments: []string{"-p", "--output-format", "stream-json", "--effort", "colossal"}, refused: "invalid"},
		{name: "a setting source the real CLI rejects", arguments: []string{"-p", "--output-format", "stream-json", "--setting-sources", "everything"}, refused: "invalid"},
		{name: "a schema that is not JSON", arguments: []string{"-p", "--output-format", "stream-json", "--json-schema", "{not json"}, refused: "not JSON"},
		{name: "the prompt in the arguments", arguments: []string{"-p", "--output-format", "stream-json", "review this change"}, refused: "standard input"},
		{name: "a flag spelled with an equals sign", arguments: []string{"-p", "--output-format=stream-json"}},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			out, failures, code := simulate(t, testCase.arguments, "review the change")

			if testCase.refused == "" {
				if code != 0 {
					t.Fatalf("the invocation was refused with %q", failures)
				}
				if out == "" {
					t.Error("nothing was written to the stream")
				}
				return
			}

			if code == 0 {
				t.Fatalf("the invocation was accepted, answering %q", out)
			}
			if !strings.Contains(failures, testCase.refused) {
				t.Errorf("the refusal reads %q, want it to name %q", failures, testCase.refused)
			}
		})
	}
}

func TestTheStreamMatchesTheWireContract(t *testing.T) {
	out, _, code := simulate(t, validArguments, "review the change")
	if code != 0 {
		t.Fatalf("the simulator exited %d", code)
	}

	events := streamed(t, out)
	if len(events) < 3 {
		t.Fatalf("the stream holds %d events, want an init, an assistant and a result", len(events))
	}
	if events[0]["type"] != "system" || events[0]["subtype"] != "init" {
		t.Errorf("the stream opens with %v, want a system init event", events[0])
	}

	assistant := events[len(events)-2]
	if assistant["type"] != "assistant" {
		t.Errorf("the event before the result is %v, want an assistant event", assistant["type"])
	}

	result := events[len(events)-1]
	if result["type"] != "result" || result["subtype"] != "success" {
		t.Errorf("the stream ends with %v, want a successful result event", result)
	}
	if result["session_id"] == "" {
		t.Error("the result event carries no session identifier")
	}
	if result["usage"] == nil {
		t.Error("the result event carries no usage")
	}

	structured, marshalled := json.Marshal(result["structured_output"])
	if marshalled != nil {
		t.Fatalf("the structured output does not re-encode: %v", marshalled)
	}
	if _, err := findings.Parse(structured); err != nil {
		t.Errorf("the structured output is not a findings document: %v", err)
	}
}

func TestATestAssertsWhatWasAsked(t *testing.T) {
	log := filepath.Join(t.TempDir(), "invocations.jsonl")
	t.Setenv(simulation.LogVariable, log)
	t.Setenv("SOCRATES_RUN", "20260819T041530Z-ee66fe34")

	const prompt = "You are reviewing a change you did not write.\n\nThe diff follows.\n"
	if _, failures, code := simulate(t, validArguments, prompt); code != 0 {
		t.Fatalf("the simulator exited %d: %s", code, failures)
	}

	invocations, err := simulation.ReadLog(log)
	if err != nil {
		t.Fatalf("reading the log: %v", err)
	}
	if len(invocations) != 1 {
		t.Fatalf("the log holds %d invocations, want one", len(invocations))
	}

	invocation := invocations[0]
	if invocation.Prompt != prompt {
		t.Errorf("the log holds the prompt %q, want %q", invocation.Prompt, prompt)
	}
	if len(invocation.Arguments) != len(validArguments) {
		t.Errorf("the log holds %d arguments, want %d", len(invocation.Arguments), len(validArguments))
	}
	if !invocation.Has("--dangerously-skip-permissions") {
		t.Error("the log does not hold the arguments as they were passed")
	}
	if schema, found := invocation.ValueOf("--json-schema"); !found || schema != findings.Schema {
		t.Error("the log does not hold the schema that was passed")
	}
	working, err := os.Getwd()
	if err != nil {
		t.Fatalf("locating the working directory: %v", err)
	}
	if invocation.Directory != working {
		t.Errorf("the log holds the working directory %q, want %q", invocation.Directory, working)
	}
	if invocation.Environment["SOCRATES_RUN"] != "20260819T041530Z-ee66fe34" {
		t.Errorf("the log holds the environment %v, without what socrates set", invocation.Environment)
	}
	if invocation.Model != "opus" || invocation.Effort != "high" {
		t.Errorf("the log holds the model %q and the effort %q", invocation.Model, invocation.Effort)
	}
	if invocation.SettingSources != "user" || invocation.Tools != "Read,Grep,Glob" {
		t.Errorf("the log holds the setting sources %q and the tools %q", invocation.SettingSources, invocation.Tools)
	}
}

func simulate(t *testing.T, arguments []string, prompt string) (string, string, int) {
	t.Helper()

	var out, failures bytes.Buffer
	code := claudesim.Run(arguments, strings.NewReader(prompt), &out, &failures, os.Environ())
	return out.String(), failures.String(), code
}

func streamed(t *testing.T, out string) []map[string]any {
	t.Helper()

	var events []map[string]any
	for line := range strings.SplitSeq(strings.TrimSuffix(out, "\n"), "\n") {
		if line == "" {
			continue
		}
		var event map[string]any
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			t.Fatalf("a line of the stream does not parse: %v: %s", err, line)
		}
		events = append(events, event)
	}
	return events
}
