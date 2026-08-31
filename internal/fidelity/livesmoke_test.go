//go:build livesmoke

package fidelity_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/ktei/socrates/internal/agent"
	"github.com/ktei/socrates/internal/document"
	"github.com/ktei/socrates/internal/findings"
	"github.com/ktei/socrates/internal/state"
	"github.com/ktei/socrates/internal/surface"
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

	t.Run("what it answered reaches a page on a port of the window", func(t *testing.T) {
		t.Setenv(state.HomeVariable, t.TempDir())

		read := whatThePersonReads(document)
		raw, err := json.Marshal(read)
		if err != nil {
			t.Fatalf("writing the document the page is served from: %v", err)
		}

		window := surface.WindowAddresses()
		free, err := surface.FirstFreeIn(window)
		if err != nil {
			t.Fatalf("taking a port of the window: %v", err)
		}

		server, err := surface.Serve(free, surface.IdlePeriod)
		if err != nil {
			t.Fatalf("serving on %s: %v", free, err)
		}
		defer server.Stop()

		if !slices.Contains(window, server.Address()) {
			t.Fatalf("the surface bound %q, and the ports it takes where nothing states one run from %q to %q", server.Address(), window[len(window)-1], window[0])
		}

		opened, err := server.Open(read, raw)
		if err != nil {
			t.Fatalf("opening a session for what the real reviewer answered: %v", err)
		}

		address := "http://" + server.Address() + surface.SessionPrefix + opened.Session
		answered, err := http.Get(address)
		if err != nil {
			t.Fatalf("reaching %s: %v", address, err)
		}
		defer answered.Body.Close()

		page, err := io.ReadAll(answered.Body)
		if err != nil {
			t.Fatalf("reading the page: %v", err)
		}
		if answered.StatusCode != http.StatusOK {
			t.Fatalf("the page at %s answered with %d", address, answered.StatusCode)
		}
		for _, decision := range read.Decisions {
			if !strings.Contains(string(page), decision.ID) {
				t.Errorf("the page at %s does not hold %s, which the real reviewer reported", address, decision.ID)
			}
		}
		t.Logf("what the real reviewer answered was served as %d decisions at %s", len(read.Decisions), address)
	})
}

func whatThePersonReads(reported findings.Document) document.Document {
	read := document.Document{Title: "what a real reviewer answered"}

	for _, finding := range reported.Findings {
		read.Decisions = append(read.Decisions, document.Decision{
			ID:       finding.ID,
			Severity: string(finding.Severity),
			Question: finding.Description,
			Mode:     document.One,
			Options: []document.Option{
				{Label: "take it", Command: "socrates respond --finding " + finding.ID + " --fixed"},
				{Label: "leave it", Command: "socrates respond --finding " + finding.ID + " --rejected"},
			},
		})
	}

	if len(read.Decisions) == 0 {
		read.Decisions = append(read.Decisions, document.Decision{
			ID:       "nothing-found",
			Question: "the reviewer found nothing blocking, and there is nothing to settle",
			Mode:     document.One,
			Options:  []document.Option{{Label: "read it and stop", Command: "socrates status"}},
		})
	}
	return read
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
