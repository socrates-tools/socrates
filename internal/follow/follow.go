package follow

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/ktei/socrates/internal/runlog"
	"github.com/ktei/socrates/internal/runstate"
	"github.com/ktei/socrates/internal/runstore"
)

const (
	DefaultPoll     = 50 * time.Millisecond
	DefaultPatience = 10 * time.Second

	settlingPolls = 2

	stateField = "state"
)

type Sink struct {
	Events  func([]runlog.Event) error
	Waiting func(runstate.State) error
}

type Following struct {
	Store    runstore.Store
	Name     string
	Poll     time.Duration
	Patience time.Duration
}

func (f Following) Until(ctx context.Context, sink Sink) (runstate.State, error) {
	run, err := f.Store.Load(f.Name)
	if err != nil {
		return "", err
	}

	reading := &tail{path: logPath(f.Store, f.Name)}
	announced := false
	patience := f.patience()
	waitingSince := time.Now()

	for {
		events, exists, err := reading.more()
		if err != nil {
			return "", err
		}
		if err := deliver(sink, events); err != nil {
			return "", err
		}
		if ended, found := terminal(events); found {
			return ended, nil
		}

		if stopping(run.State) {
			return f.settle(ctx, sink, reading, run.State)
		}

		if !exists {
			if !announced {
				if err := waiting(sink, run.State); err != nil {
					return "", err
				}
				announced = true
			}
			if time.Since(waitingSince) >= patience {
				return runstate.Failed, fmt.Errorf(nothingEverWrote, f.Name, patience)
			}
		}

		if err := f.rest(ctx); err != nil {
			return "", err
		}

		run, err = f.Store.Load(f.Name)
		if err != nil {
			return "", err
		}
	}
}

func (f Following) settle(ctx context.Context, sink Sink, reading *tail, state runstate.State) (runstate.State, error) {
	for range settlingPolls {
		if err := f.rest(ctx); err != nil {
			return "", err
		}

		events, _, err := reading.more()
		if err != nil {
			return "", err
		}
		if err := deliver(sink, events); err != nil {
			return "", err
		}
		if ended, found := terminal(events); found {
			return ended, nil
		}
	}
	return state, nil
}

func (f Following) rest(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(f.poll()):
		return nil
	}
}

func (f Following) poll() time.Duration {
	if f.Poll <= 0 {
		return DefaultPoll
	}
	return f.Poll
}

func (f Following) patience() time.Duration {
	if f.Patience <= 0 {
		return DefaultPatience
	}
	return f.Patience
}

func deliver(sink Sink, events []runlog.Event) error {
	if len(events) == 0 || sink.Events == nil {
		return nil
	}
	return sink.Events(events)
}

func waiting(sink Sink, state runstate.State) error {
	if sink.Waiting == nil {
		return nil
	}
	return sink.Waiting(state)
}

func stopping(state runstate.State) bool {
	return state.Terminal() || state == runstate.AwaitingResponse
}

func terminal(events []runlog.Event) (runstate.State, bool) {
	for _, event := range events {
		stated, isText := event.Fields[stateField].(string)
		if !isText {
			continue
		}
		if state := runstate.State(stated); state.Terminal() {
			return state, true
		}
	}
	return "", false
}

func Held(store runstore.Store, name string) ([]runlog.Event, error) {
	reading := tail{path: logPath(store, name)}

	events, _, err := reading.more()
	if err != nil {
		return nil, err
	}
	return events, nil
}

func logPath(store runstore.Store, name string) string {
	return filepath.Join(store.Path(name), runlog.FileName)
}

type tail struct {
	path      string
	delivered int
	rotations int
}

func (t *tail) more() ([]runlog.Event, bool, error) {
	content, exists, err := t.read()
	if err != nil || !exists {
		return nil, exists, err
	}

	held, err := t.parse(content)
	if err != nil {
		return nil, true, err
	}

	rotations := 0
	kept := 0
	for position, event := range held {
		if event.Kind == runlog.EventRotated {
			rotations++
			kept = position
		}
	}

	if rotations > t.rotations {
		t.rotations = rotations
		t.delivered = len(held)
		return held[kept:], true, nil
	}

	if t.delivered > len(held) {
		t.delivered = len(held)
	}
	fresh := held[t.delivered:]
	t.delivered = len(held)
	return fresh, true, nil
}

func (t *tail) read() ([]byte, bool, error) {
	content, err := os.ReadFile(t.path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf(watchingTheLog, t.path, err)
	}
	return content, true, nil
}

func (t *tail) parse(content []byte) ([]runlog.Event, error) {
	var events []runlog.Event

	for {
		line, rest, complete := bytes.Cut(content, []byte("\n"))
		if !complete {
			return events, nil
		}
		content = rest

		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		var event runlog.Event
		if err := json.Unmarshal(line, &event); err != nil {
			return events, fmt.Errorf(unreadableLine, t.path, err)
		}
		events = append(events, event)
	}
}
