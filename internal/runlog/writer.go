package runlog

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

const FileName = "events.jsonl"

type Writer struct {
	mutex sync.Mutex
	file  *os.File
	now   func() time.Time
}

func Open(path string) (*Writer, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", path, err)
	}
	return &Writer{file: file, now: time.Now}, nil
}

func (w *Writer) Record(event Event) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if event.At.IsZero() {
		event.At = w.now().UTC()
	}

	encoded, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("encoding a %s event: %w", event.Kind, err)
	}
	if _, err := w.file.Write(append(encoded, '\n')); err != nil {
		return fmt.Errorf("writing a %s event: %w", event.Kind, err)
	}
	return nil
}

func (w *Writer) Close() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.file.Close()
}

func Read(path string) ([]Event, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	var events []Event
	for line := range strings.SplitSeq(string(content), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var event Event
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			return nil, fmt.Errorf("reading a line of %s: %w", path, err)
		}
		events = append(events, event)
	}
	return events, nil
}
