package agent

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/testing/source"
	"github.com/socrates-tools/socrates/review/findings"
	"github.com/socrates-tools/socrates/review/prompt"
	"github.com/socrates-tools/socrates/review/simulation"
)

func TestTheParsedTierIsExercisedAtTheWire(t *testing.T) {
	for _, answered := range []struct {
		named    string
		scenario string
		verdict  bool
	}{
		{named: "a reviewer that ends with a document", scenario: "", verdict: true},
		{named: "a reviewer that answered without a document", scenario: "default:\n  behaviour: no-structured-output\n"},
		{named: "a reviewer whose document violates the schema", scenario: "default:\n  behaviour: schema-violation\n"},
	} {
		t.Run(answered.named, func(t *testing.T) {
			reviewing := onTheParsedTier(t, answered.scenario)

			asked := prompt.Build(prompt.Request{
				Intent:                 "add a parser",
				Repository:             "/repo",
				Base:                   "1111111111111111111111111111111111111111",
				Snapshot:               "2222222222222222222222222222222222222222",
				Changed:                []string{"internal/added.go"},
				Diff:                   "diff --git a/internal/added.go b/internal/added.go\n",
				Ceiling:                1 << 20,
				DocumentEndsTheMessage: true,
			})

			answer, err := reviewing.Run(context.Background(), Request{
				Role:    Reviewer,
				Prompt:  asked.Text,
				Dir:     t.TempDir(),
				Schema:  json.RawMessage(findings.Schema),
				Honours: func(document json.RawMessage) error { _, err := findings.Parse(document); return err },
			})

			if !answered.verdict {
				if err == nil {
					t.Fatalf("%s answered and the call succeeded", answered.named)
				}
				return
			}

			if err != nil {
				t.Fatalf("%s answered and the call failed: %v", answered.named, err)
			}
			if len(answer.Structured) == 0 {
				t.Fatal("nothing was extracted from the reviewer's final message")
			}
			if _, err := findings.Parse(answer.Structured); err != nil {
				t.Errorf("what was extracted does not reach the same validation: %v", err)
			}
			if strings.Contains(string(answer.Structured), "```") {
				t.Errorf("what was extracted still carries its fence: %s", answer.Structured)
			}
		})
	}
}

func onTheParsedTier(t *testing.T, scenario string) Agent {
	t.Helper()

	directory := t.TempDir()
	binary := filepath.Join(directory, "claudesim")
	source.Build(t, "./review/cmd/claudesim", binary)

	if scenario != "" {
		path := filepath.Join(directory, "scenario.yaml")
		if err := os.WriteFile(path, []byte(scenario), 0o644); err != nil {
			t.Fatalf("writing the scenario: %v", err)
		}
		t.Setenv(simulation.ScenarioVariable, path)
	}
	t.Setenv(simulation.LogVariable, filepath.Join(directory, "invocations.jsonl"))

	extracting := Vendor{
		Name:                      TheDefault().Name,
		Binary:                    binary,
		Tier:                      Parsed,
		IsolatedFromTheRepository: true,
		ReportsWhyItFailed:        true,
		ModelSelectablePerRole:    true,
	}

	reviewing, err := newClaude(extracting, Options{Binary: binary})
	if err != nil {
		t.Fatalf("constructing an agent on the parsed tier: %v", err)
	}
	return reviewing
}
