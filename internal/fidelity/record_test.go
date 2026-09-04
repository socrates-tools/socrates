//go:build record

package fidelity_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ktei/socrates/internal/agent"
	"github.com/ktei/socrates/internal/fidelity"
	"github.com/ktei/socrates/internal/findings"
)

const (
	recordingDeadline = 10 * time.Minute

	attemptsOnTheParsedTier = 5

	anAccountIsAvailable = "SOCRATES_RECORD_AGAINST_AN_ACCOUNT"
	plainAnswer          = "ready"
	recordableDocument   = `{"findings":[{"id":"f1","severity":"warning","file":"parse.go","line":4,"description":"Parse sums the characters and can never fail, and its error return says otherwise to every caller.","action":"auto-fix"}],"risk_level":"low","risk_rationale":"one small function"}`
)

func TestRecordWhatTheRealCLIAnswers(t *testing.T) {
	for _, vendor := range agent.Supported() {
		t.Run(vendor.Name, func(t *testing.T) {
			for _, request := range fidelity.Requests() {
				t.Run(request.Name, func(t *testing.T) {
					stream, runErr := askTheRealCLI(t, vendor, request)
					if stream == "" {
						t.Fatalf("%s wrote nothing to record: %v", vendor.Name, runErr)
					}

					path, err := fidelity.Keep(fidelity.TestdataDirectory, vendor.Name, request.Name, stream, runErr)
					if err != nil {
						t.Fatalf("keeping what %s wrote: %v", vendor.Name, err)
					}

					if runErr != nil {
						t.Fatalf("%s did not answer %s: %v: what it wrote is under %s, and the recordings were left alone", vendor.Name, request.Name, runErr, path)
					}
					t.Logf("recorded %s", path)
				})
			}
		})
	}
}

func askTheRealCLI(t *testing.T, vendor agent.Vendor, request fidelity.Request) (string, error) {
	t.Helper()

	if answering := withoutAnAccount(t, vendor, request); answering != nil {
		return answering()
	}

	asking, err := agent.New(vendor.Name, agent.Options{Deadline: recordingDeadline})
	if err != nil {
		t.Fatalf("constructing the agent for %s: %v", vendor.Name, err)
	}

	var stream bytes.Buffer
	var runErr error

	for attempt := 1; attempt <= attemptsAllowed(vendor); attempt++ {
		stream.Reset()
		_, runErr = asking.Run(context.Background(), agent.Request{
			Role:       agent.Reviewer,
			Prompt:     fidelity.Asked(vendor, request),
			Dir:        subjectRepository(t),
			Scratch:    t.TempDir(),
			Schema:     request.Schema,
			Honours:    honours(request),
			Transcript: &stream,
		})
		if runErr == nil {
			recordAttempts(t, vendor, request, attempt, true)
			return stream.String(), nil
		}
		t.Logf("%s attempt %d of %d did not answer %s: %v", vendor.Name, attempt, attemptsAllowed(vendor), request.Name, runErr)
	}

	recordAttempts(t, vendor, request, attemptsAllowed(vendor), false)
	return stream.String(), runErr
}

func attemptsAllowed(vendor agent.Vendor) int {
	if vendor.Tier == agent.Parsed {
		return attemptsOnTheParsedTier
	}
	return 1
}

func recordAttempts(t *testing.T, vendor agent.Vendor, request fidelity.Request, attempts int, answered bool) {
	t.Helper()

	if vendor.Tier != agent.Parsed || len(request.Schema) == 0 {
		return
	}

	kept, err := fidelity.KeepObservation(fidelity.TestdataDirectory, fidelity.Observation{
		Vendor:    vendor.Name,
		Version:   "recorded by make fixtures",
		Criterion: "what-the-extracting-tier-costs",
		Question:  "how many attempts did one live review on the extracting tier take before a document arrived that parsed and satisfied the schema",
		Observed: map[string]string{
			"attempts":         strconv.Itoa(attempts),
			"attempts allowed": strconv.Itoa(attemptsOnTheParsedTier),
			"answered":         strconv.FormatBool(answered),
			"request":          request.Name,
		},
		Established: answered,
	})
	if err != nil {
		t.Fatalf("recording what the extracting tier cost: %v", err)
	}
	t.Logf("recorded %s", kept)
}

func honours(request fidelity.Request) func(json.RawMessage) error {
	if len(request.Schema) == 0 {
		return nil
	}
	return func(document json.RawMessage) error {
		_, err := findings.Parse(document)
		return err
	}
}

func withoutAnAccount(t *testing.T, vendor agent.Vendor, request fidelity.Request) func() (string, error) {
	t.Helper()

	if vendor.Name != "codex" || os.Getenv(anAccountIsAvailable) != "" {
		return nil
	}

	answer := plainAnswer
	if len(request.Schema) > 0 {
		answer = recordableDocument
	}
	address := answeringProvider(t, answer)
	directory := subjectRepository(t)

	return func() (string, error) {
		arguments := []string{"exec", "--json", "--sandbox", "read-only", "--ignore-user-config", "-c", "project_doc_max_bytes=0"}

		if len(request.Schema) > 0 {
			schema := filepath.Join(t.TempDir(), "schema.json")
			if err := os.WriteFile(schema, request.Schema, 0o644); err != nil {
				return "", err
			}
			arguments = append(arguments, "--output-schema", schema)
		}

		arguments = append(arguments,
			"-c", "model_provider=local",
			"-c", "model_providers.local.name=local",
			"-c", "model_providers.local.base_url="+address+"/v1",
			"-c", "model_providers.local.wire_api=responses",
			"-c", "model_providers.local.env_key=LOCAL_KEY",
			"-c", "model_providers.local.request_max_retries=0",
			"-c", "model_providers.local.stream_max_retries=0",
			"-",
		)

		command := exec.Command(vendor.Binary, arguments...)
		command.Dir = directory
		command.Stdin = strings.NewReader(request.Prompt)

		var stream, failures bytes.Buffer
		command.Stdout = &stream
		command.Stderr = &failures
		err := command.Run()
		return stream.String(), err
	}
}

func subjectRepository(t *testing.T) string {
	t.Helper()

	directory := t.TempDir()
	path := filepath.Join(directory, "parse.go")
	if err := os.WriteFile(path, []byte(subject), 0o644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
	for _, arguments := range [][]string{
		{"init", "-q", "."},
		{"add", "-A"},
		{"-c", "user.email=recording@socrates.invalid", "-c", "user.name=recording", "commit", "-qm", "the subject"},
	} {
		command := exec.Command("git", arguments...)
		command.Dir = directory
		if written, err := command.CombinedOutput(); err != nil {
			t.Fatalf("preparing the repository: %v: %s", err, written)
		}
	}
	return directory
}

const subject = `package parser

func Parse(document string) (int, error) {
	total := 0
	for _, character := range document {
		total += int(character)
	}
	return total, nil
}
`
