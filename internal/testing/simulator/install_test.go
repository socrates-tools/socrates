package simulator_test

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/testing/simulator"
	"github.com/socrates-tools/socrates/review/findings"
)

func TestTheSimulatorAnswersUnderTheNameSocratesInvokes(t *testing.T) {
	simulated := simulator.Install(t)
	simulated.Scenario("default:\n  text: I read the change\n")

	found, err := exec.LookPath("claude")
	if err != nil {
		t.Fatalf("looking for claude on PATH: %v", err)
	}
	if found != simulated.Link {
		t.Errorf("claude on PATH is %s, want the simulator at %s", found, simulated.Link)
	}

	const prompt = "review this change"
	command := exec.Command("claude",
		"-p",
		"--verbose",
		"--output-format", "stream-json",
		"--setting-sources", "user",
		"--json-schema", findings.Schema,
	)
	command.Dir = t.TempDir()
	command.Stdin = strings.NewReader(prompt)

	var out, failures bytes.Buffer
	command.Stdout = &out
	command.Stderr = &failures
	if err := command.Run(); err != nil {
		t.Fatalf("running the simulator through PATH: %v: %s", err, failures.String())
	}

	lines := strings.Split(strings.TrimSuffix(out.String(), "\n"), "\n")
	var result map[string]any
	if err := json.Unmarshal([]byte(lines[len(lines)-1]), &result); err != nil {
		t.Fatalf("the last line of the stream does not parse: %v", err)
	}
	structured, err := json.Marshal(result["structured_output"])
	if err != nil {
		t.Fatalf("re-encoding the structured output: %v", err)
	}
	if _, err := findings.Parse(structured); err != nil {
		t.Errorf("the simulator did not answer a findings document: %v", err)
	}

	invocation := simulated.OnlyInvocation()
	if invocation.Prompt != prompt {
		t.Errorf("the log holds the prompt %q, want %q", invocation.Prompt, prompt)
	}
	if invocation.Directory != command.Dir {
		t.Errorf("the log holds the working directory %q, want %q", invocation.Directory, command.Dir)
	}
	if !invocation.Has("--verbose") {
		t.Errorf("the log holds the arguments %v", invocation.Arguments)
	}
}

func TestTheSimulatorRecordsWhatItIsAskedTo(t *testing.T) {
	simulated := simulator.Install(t)
	simulated.Record("SCRAT_TEST_MARKER")

	command := exec.Command("claude", "-p", "--output-format", "stream-json")
	command.Dir = t.TempDir()
	command.Stdin = strings.NewReader("review this change")
	command.Env = append(command.Environ(), "SCRAT_TEST_MARKER=present")

	if out, err := command.CombinedOutput(); err != nil {
		t.Fatalf("running the simulator: %v: %s", err, out)
	}

	if marker := simulated.OnlyInvocation().Environment["SCRAT_TEST_MARKER"]; marker != "present" {
		t.Errorf("the log holds the marker %q, want present", marker)
	}
}

func TestTheSimulatorIsBuiltFromThisTree(t *testing.T) {
	simulated := simulator.Install(t)

	if filepath.Base(simulated.Link) != "claude" {
		t.Errorf("the simulator is linked as %s", filepath.Base(simulated.Link))
	}
	if filepath.Dir(simulated.Binary) != simulated.Directory {
		t.Errorf("the simulator was built outside its own directory: %s", simulated.Binary)
	}
}
