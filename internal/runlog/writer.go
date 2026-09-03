package runlog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	FileName       = "events.jsonl"
	DefaultCeiling = 1 << 20
	EventRotated   = "log-rotated"
)

type Writer struct {
	mutex   sync.Mutex
	file    *os.File
	path    string
	ceiling int64
	size    int64
	now     func() time.Time
}

func Open(path string, ceiling int64) (*Writer, error) {
	if ceiling <= 0 {
		ceiling = DefaultCeiling
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf(opening, path, err)
	}

	written, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return nil, fmt.Errorf(measuring, path, err)
	}

	return &Writer{file: file, path: path, ceiling: ceiling, size: written.Size(), now: time.Now}, nil
}

func (w *Writer) Record(event Event) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	line, err := w.line(event)
	if err != nil {
		return err
	}

	if w.size+int64(len(line)) > w.ceiling {
		if err := w.rotate(); err != nil {
			return err
		}
	}
	return w.write(line)
}

func (w *Writer) Close() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.file.Close()
}

func (w *Writer) line(event Event) ([]byte, error) {
	if event.At.IsZero() {
		event.At = w.now().UTC()
	}

	encoded, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf(encodingAnEvent, event.Kind, err)
	}
	return append(encoded, '\n'), nil
}

func (w *Writer) write(line []byte) error {
	if _, err := w.file.Write(line); err != nil {
		return fmt.Errorf(writingAnEvent, w.path, err)
	}
	w.size += int64(len(line))
	return nil
}

func (w *Writer) rotate() error {
	kept, err := w.tail()
	if err != nil {
		return err
	}

	if err := w.file.Truncate(0); err != nil {
		return fmt.Errorf(truncating, w.path, err)
	}
	w.size = 0

	if err := w.write(kept); err != nil {
		return err
	}

	line, err := w.line(Event{
		Kind:    EventRotated,
		Subject: FileName,
		Fields:  map[string]any{"ceiling": w.ceiling, "kept": len(kept)},
	})
	if err != nil {
		return err
	}
	return w.write(line)
}

func (w *Writer) tail() ([]byte, error) {
	content, err := os.ReadFile(w.path)
	if err != nil {
		return nil, fmt.Errorf(reading, w.path, err)
	}

	room := w.ceiling / 2
	if int64(len(content)) <= room {
		return content, nil
	}

	kept := content[int64(len(content))-room:]
	if position := bytes.IndexByte(kept, '\n'); position >= 0 {
		kept = kept[position+1:]
	}
	return kept, nil
}

func Read(path string) ([]Event, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf(reading, path, err)
	}

	var events []Event
	for line := range strings.SplitSeq(string(content), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var event Event
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			return nil, fmt.Errorf(readingALine, path, err)
		}
		events = append(events, event)
	}
	return events, nil
}
