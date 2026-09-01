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

type recorded struct {
	line  []byte
	event runlog.Event
}

type tail struct {
	path      string
	delivered []byte
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

	fresh := held[t.resume(held):]
	if len(fresh) == 0 {
		return nil, true, nil
	}
	t.delivered = fresh[len(fresh)-1].line

	events := make([]runlog.Event, len(fresh))
	for position, held := range fresh {
		events[position] = held.event
	}
	return events, true, nil
}

func (t *tail) resume(held []recorded) int {
	if t.delivered == nil {
		return 0
	}
	for position := len(held) - 1; position >= 0; position-- {
		if bytes.Equal(held[position].line, t.delivered) {
			return position + 1
		}
	}
	return fromTheLastRotation(held)
}

func fromTheLastRotation(held []recorded) int {
	for position := len(held) - 1; position >= 0; position-- {
		if held[position].event.Kind == runlog.EventRotated {
			return position
		}
	}
	return 0
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

func (t *tail) parse(content []byte) ([]recorded, error) {
	var held []recorded

	for {
		line, rest, complete := bytes.Cut(content, []byte("\n"))
		if !complete {
			return held, nil
		}
		content = rest

		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		var event runlog.Event
		if err := json.Unmarshal(line, &event); err != nil {
			return held, fmt.Errorf(unreadableLine, t.path, err)
		}
		held = append(held, recorded{line: bytes.Clone(line), event: event})
	}
}
