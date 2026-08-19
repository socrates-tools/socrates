//go:build record

package fidelity_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rui-ktei/socrates/internal/agent"
	"github.com/rui-ktei/socrates/internal/fidelity"
)

const recordingDeadline = 10 * time.Minute

func TestRecordWhatTheRealCLIAnswers(t *testing.T) {
	directory := fidelity.RecordingsDirectory

	for _, request := range fidelity.Requests() {
		t.Run(request.Name, func(t *testing.T) {
			asking, err := agent.NewClaude(agent.Options{Deadline: recordingDeadline})
			if err != nil {
				t.Fatalf("constructing the agent: %v", err)
			}

			var stream bytes.Buffer
			_, runErr := asking.Run(context.Background(), agent.Request{
				Role:       agent.Reviewer,
				Prompt:     request.Prompt,
				Dir:        subjectRepository(t),
				Schema:     request.Schema,
				Transcript: &stream,
			})

			if stream.Len() == 0 {
				t.Fatalf("the real CLI wrote nothing to record: %v", runErr)
			}

			shapes, err := fidelity.ShapesOf(stream.String())
			if err != nil {
				t.Fatalf("reading what the real CLI wrote: %v", err)
			}

			path, err := fidelity.Record(directory, request.Name, shapes)
			if err != nil {
				t.Fatalf("recording: %v", err)
			}
			t.Logf("recorded %s", path)

			if runErr != nil {
				t.Errorf("the real CLI answered %s, and the stream was recorded anyway: %v", request.Name, runErr)
			}
		})
	}
}

func subjectRepository(t *testing.T) string {
	t.Helper()

	directory := t.TempDir()
	path := filepath.Join(directory, "parse.go")
	if err := os.WriteFile(path, []byte(subject), 0o644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
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
