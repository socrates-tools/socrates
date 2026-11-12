package follow

import (
	"errors"
	"testing"
)

const (
	whole = `{"at":"2026-11-05T04:02:42Z","kind":"run-started","subject":"r"}` + "\n" +
		`{"at":"2026-11-05T04:02:43Z","kind":"agent-started","subject":"reviewer"}` + "\n"
	tornAcrossARotation = `{"at":"2026-11-05T04:02:42Z","kind":"run-started","subject":"r"}` + "\n" +
		`subject":"reviewer"}` + "\n"
)

func readsInTurn(contents ...string) func(string) ([]byte, error) {
	return func(string) ([]byte, error) {
		if len(contents) == 0 {
			return nil, errors.New("read more often than the test expected")
		}
		next := contents[0]
		if len(contents) > 1 {
			contents = contents[1:]
		}
		return []byte(next), nil
	}
}

func TestALogReadTornAcrossARotationIsReadAgainRatherThanEndingTheFollow(t *testing.T) {
	reading := tail{path: "events.jsonl", readFile: readsInTurn(tornAcrossARotation, whole)}

	events, exists, err := reading.more()
	if err != nil {
		t.Fatalf("a read torn across a rotation ended the follow: %v", err)
	}
	if !exists || len(events) != 2 {
		t.Errorf("the follower delivered %d events from a log that holds 2", len(events))
	}
}

func TestALineThatStaysUnreadableAcrossTwoReadsEndsTheFollow(t *testing.T) {
	reading := tail{path: "events.jsonl", readFile: readsInTurn(tornAcrossARotation)}

	if _, _, err := reading.more(); err == nil {
		t.Error("a log holding a line that is no event on every read was followed as if it were whole")
	}
}
