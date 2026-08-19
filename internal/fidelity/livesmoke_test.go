//go:build livesmoke

package fidelity_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rui-ktei/socrates/internal/agent"
	"github.com/rui-ktei/socrates/internal/findings"
)

const smokeDeadline = 10 * time.Minute

func TestTheRealCLIStillAcceptsWhatSocratesAsksOfIt(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, "parse.go")
	if err := os.WriteFile(path, []byte(smokeSubject), 0o644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}

	asking, err := agent.NewClaude(agent.Options{Deadline: smokeDeadline})
	if err != nil {
		t.Fatalf("constructing the agent: %v", err)
	}

	var stream bytes.Buffer
	answer, err := asking.Run(context.Background(), agent.Request{
		Role:       agent.Reviewer,
		Prompt:     "Read parse.go in your working directory and report what you find as the schema requires.",
		Dir:        directory,
		Schema:     []byte(findings.Schema),
		Honours:    func(structured json.RawMessage) error { _, err := findings.Parse(structured); return err },
		Transcript: &stream,
	})
	if err != nil {
		t.Fatalf("the real CLI refused what socrates asks of it: %v\n%s", err, stream.String())
	}

	document, err := findings.Parse(answer.Structured)
	if err != nil {
		t.Fatalf("the real CLI returned structured output socrates cannot read: %v\n%s", err, answer.Structured)
	}
	if answer.Model == "" {
		t.Error("the stream carried no model")
	}
	if answer.Usage.InputTokens == 0 && answer.Usage.OutputTokens == 0 {
		t.Error("the stream carried no token counts")
	}
	t.Logf("the real reviewer rated the change %s risk and reported %d findings", document.RiskLevel, len(document.Findings))
}

const smokeSubject = `package parser

func Parse(document string) (int, error) {
	total := 0
	for _, character := range document {
		total += int(character)
	}
	return total, nil
}
`
