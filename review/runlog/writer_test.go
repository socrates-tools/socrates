package runlog_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/review/runlog"
)

func TestEventsAreOneJSONObjectPerLine(t *testing.T) {
	path := filepath.Join(t.TempDir(), runlog.FileName)

	writer, err := runlog.Open(path, 0)
	if err != nil {
		t.Fatalf("opening %s: %v", path, err)
	}
	for _, kind := range []string{"run-started", "agent-started", "run-ended"} {
		if err := writer.Record(runlog.Event{Kind: kind, Subject: "a-run"}); err != nil {
			t.Fatalf("recording a %s event: %v", kind, err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("closing %s: %v", path, err)
	}

	events, err := runlog.Read(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	if len(events) != 3 {
		t.Fatalf("%d events were read, want 3", len(events))
	}
	for _, event := range events {
		if event.At.IsZero() || event.Kind == "" || event.Subject == "" {
			t.Errorf("the event %+v carries no instant, kind or subject", event)
		}
	}
}

func TestAChattyReviewerDoesNotFillTheDisk(t *testing.T) {
	const ceiling = 4 << 10
	path := filepath.Join(t.TempDir(), runlog.FileName)

	writer, err := runlog.Open(path, ceiling)
	if err != nil {
		t.Fatalf("opening %s: %v", path, err)
	}
	defer func() { _ = writer.Close() }()

	for count := range 400 {
		if err := writer.Record(runlog.Event{Kind: "agent-output", Subject: "reviewer", Fields: map[string]any{"count": count}}); err != nil {
			t.Fatalf("recording event %d: %v", count, err)
		}
	}

	written, err := os.Stat(path)
	if err != nil {
		t.Fatalf("measuring %s: %v", path, err)
	}
	if written.Size() > ceiling {
		t.Errorf("the log is %d bytes, over its %d byte ceiling", written.Size(), ceiling)
	}

	events, err := runlog.Read(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}

	rotations := 0
	for _, event := range events {
		if event.Kind == runlog.EventRotated {
			rotations++
		}
	}
	if rotations == 0 {
		t.Error("the log holds no record that it rotated, so a reader would read it as the whole story")
	}

	last := events[len(events)-1]
	if last.Fields["count"] != float64(399) {
		t.Errorf("the last event is %+v, and the most recent events are what a rotation keeps", last)
	}
	if first := events[0]; first.Fields["count"] == float64(0) {
		t.Error("the log still holds the first event, so nothing was dropped and the ceiling did nothing")
	}
}

func TestTheWriterKeepsWritingAcrossARotation(t *testing.T) {
	const ceiling = 2 << 10
	path := filepath.Join(t.TempDir(), runlog.FileName)

	writer, err := runlog.Open(path, ceiling)
	if err != nil {
		t.Fatalf("opening %s: %v", path, err)
	}
	defer func() { _ = writer.Close() }()

	before, err := os.Stat(path)
	if err != nil {
		t.Fatalf("measuring %s: %v", path, err)
	}

	for count := range 200 {
		if err := writer.Record(runlog.Event{Kind: "agent-output", Subject: "reviewer", Fields: map[string]any{"count": count}}); err != nil {
			t.Fatalf("recording event %d: %v", count, err)
		}
	}
	if err := writer.Record(runlog.Event{Kind: "run-ended", Subject: "a-run"}); err != nil {
		t.Fatalf("recording after the rotation: %v", err)
	}

	after, err := os.Stat(path)
	if err != nil {
		t.Fatalf("measuring %s: %v", path, err)
	}
	if !os.SameFile(before, after) {
		t.Error("the log is a different file after rotating, so the writer's own descriptor points at a file nobody reads")
	}

	events, err := runlog.Read(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	if last := events[len(events)-1]; last.Kind != "run-ended" {
		t.Errorf("the last event is %+v, want the one written after the rotation", last)
	}

	entries, err := os.ReadDir(filepath.Dir(path))
	if err != nil {
		t.Fatalf("reading the directory: %v", err)
	}
	for _, entry := range entries {
		if entry.Name() != runlog.FileName {
			t.Errorf("rotating left %s behind, and a rotation that renames leaves the writer writing where nobody reads", entry.Name())
		}
	}
}

func TestALineOfSomethingElseIsNotAnEvent(t *testing.T) {
	path := filepath.Join(t.TempDir(), runlog.FileName)
	if err := os.WriteFile(path, []byte("not an event\n"), 0o644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}

	_, err := runlog.Read(path)
	if err == nil {
		t.Fatal("a line that is not JSON was read as an event")
	}
	if !strings.Contains(err.Error(), path) {
		t.Errorf("the failure does not name the file: %v", err)
	}
}
