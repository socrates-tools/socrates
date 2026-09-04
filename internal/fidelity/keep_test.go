package fidelity_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/socrates-tools/socrates-review/internal/fidelity"
)

func TestARunThatAnsweredIsRecorded(t *testing.T) {
	root := t.TempDir()

	path, err := fidelity.Keep(root, "claude", "plain-answer", recordedStream, nil)
	if err != nil {
		t.Fatalf("keeping what the CLI answered: %v", err)
	}

	recordings, err := fidelity.Recordings(filepath.Dir(path))
	if err != nil {
		t.Fatalf("reading back what was recorded: %v", err)
	}
	if len(recordings) != 1 || recordings[0].Name != "plain-answer" {
		t.Errorf("an answered run was kept as %v, and one recording named plain-answer is what it is", recordings)
	}
}

func TestARunThatDidNotAnswerLeavesTheRecordingsAlone(t *testing.T) {
	root := t.TempDir()

	recorded, err := fidelity.Keep(root, "claude", "plain-answer", recordedStream, nil)
	if err != nil {
		t.Fatalf("recording an answer for a later run to overwrite: %v", err)
	}
	before, err := os.ReadFile(recorded)
	if err != nil {
		t.Fatalf("reading the recording: %v", err)
	}

	kept, err := fidelity.Keep(root, "claude", "plain-answer", refusedStream, errors.New("claude exited: exit status 1: Not logged in"))
	if err != nil {
		t.Fatalf("keeping what a run that did not answer wrote: %v", err)
	}
	if kept == recorded {
		t.Fatalf("a run that did not answer was kept at %s, which is where the recording lives", kept)
	}
	if _, err := os.Stat(kept); err != nil {
		t.Errorf("what a run that did not answer wrote was not kept for a person to read: %v", err)
	}

	after, err := os.ReadFile(recorded)
	if err != nil {
		t.Fatalf("reading the recording back: %v", err)
	}
	if string(after) != string(before) {
		t.Errorf("the recording was rewritten by a run that did not answer:\nbefore: %s\nafter: %s", before, after)
	}
}

const refusedStream = `{"type":"system","subtype":"init","session_id":"s","model":"claude","tools":["Read"],"cwd":"/tmp"}
{"type":"result","subtype":"success","session_id":"s","is_error":true,"result":"Not logged in · Please run /login","duration_ms":1,"num_turns":0}
`
