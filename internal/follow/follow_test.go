package follow_test

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
	"testing/synctest"
	"time"

	"github.com/rui-ktei/socrates/internal/follow"
	"github.com/rui-ktei/socrates/internal/runlog"
	"github.com/rui-ktei/socrates/internal/runstate"
	"github.com/rui-ktei/socrates/internal/runstore"
)

const repository = "/home/a/project"

func TestARunThatFailedBeforeWritingAnything(t *testing.T) {
	store, name := run(t, func(run *runstore.Run) {
		if err := run.End(runstate.Failed, "the reviewer went away before it wrote anything"); err != nil {
			t.Fatalf("ending the run: %v", err)
		}
	})

	type ending struct {
		state runstate.State
		err   error
	}
	stopped := make(chan ending, 1)
	go func() {
		state, err := following(store, name).Until(context.Background(), follow.Sink{})
		stopped <- ending{state: state, err: err}
	}()

	select {
	case reached := <-stopped:
		if reached.err != nil {
			t.Fatalf("following a run that failed before writing anything: %v", reached.err)
		}
		if reached.state != runstate.Failed {
			t.Errorf("the follower stopped at %s, want %s", reached.state, runstate.Failed)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("the follower is still waiting for an event a failed run will never write")
	}
}

func TestATerminalEventStopsAFollowerThatRunJSONWouldNot(t *testing.T) {
	store, name := run(t, nil)
	write(t, store, name, runlog.Event{
		At:      time.Now().UTC(),
		Kind:    "run-ended",
		Subject: name,
		Fields:  map[string]any{"state": string(runstate.Clean), "reason": "nothing blocking"},
	})

	if recorded := load(t, store, name); recorded.State != runstate.Reviewing {
		t.Fatalf("this test needs a run whose own record is not terminal, and it reads %s", recorded.State)
	}

	state, err := following(store, name).Until(context.Background(), follow.Sink{})
	if err != nil {
		t.Fatalf("following a run whose log ended before its record did: %v", err)
	}
	if state != runstate.Clean {
		t.Errorf("the follower stopped at %s, want the %s the log recorded", state, runstate.Clean)
	}
}

func TestEveryEventIsDeliveredOnce(t *testing.T) {
	store, name := run(t, nil)

	var delivered []runlog.Event
	sink := follow.Sink{Events: func(events []runlog.Event) error {
		delivered = append(delivered, events...)
		return nil
	}}

	write(t, store, name, runlog.Event{At: time.Now().UTC(), Kind: "run-started", Subject: name})
	write(t, store, name, runlog.Event{At: time.Now().UTC(), Kind: "agent-started", Subject: "reviewer"})

	go func() {
		time.Sleep(100 * time.Millisecond)
		write(t, store, name, runlog.Event{At: time.Now().UTC(), Kind: "agent-completed", Subject: "reviewer"})
		write(t, store, name, runlog.Event{
			At:      time.Now().UTC(),
			Kind:    "run-ended",
			Subject: name,
			Fields:  map[string]any{"state": string(runstate.Escalated), "reason": "one finding for the author"},
		})
	}()

	state, err := following(store, name).Until(context.Background(), follow.Sink{Events: sink.Events})
	if err != nil {
		t.Fatalf("following a run to its end: %v", err)
	}
	if state != runstate.Escalated {
		t.Errorf("the follower stopped at %s, want %s", state, runstate.Escalated)
	}

	var kinds []string
	for _, event := range delivered {
		kinds = append(kinds, event.Kind)
	}
	want := []string{"run-started", "agent-started", "agent-completed", "run-ended"}
	if len(kinds) != len(want) {
		t.Fatalf("the follower delivered %v, want %v", kinds, want)
	}
	for position, kind := range want {
		if kinds[position] != kind {
			t.Errorf("event %d was delivered as %s, want %s", position+1, kinds[position], kind)
		}
	}
}

func TestARotatedLogIsFollowedFromWhatItKept(t *testing.T) {
	store, name := run(t, nil)
	ceiling := int64(600)

	for position := range 3 {
		writeWith(t, store, name, ceiling, runlog.Event{At: time.Now().UTC(), Kind: "agent-started", Subject: "reviewer", Fields: map[string]any{"position": position}})
	}

	var delivered []runlog.Event
	sink := follow.Sink{Events: func(events []runlog.Event) error {
		delivered = append(delivered, events...)
		return nil
	}}

	go func() {
		time.Sleep(100 * time.Millisecond)
		for position := range 6 {
			writeWith(t, store, name, ceiling, runlog.Event{At: time.Now().UTC(), Kind: "agent-completed", Subject: "reviewer", Fields: map[string]any{"position": position}})
		}
		writeWith(t, store, name, ceiling, runlog.Event{
			At:      time.Now().UTC(),
			Kind:    "run-ended",
			Subject: name,
			Fields:  map[string]any{"state": string(runstate.Clean), "reason": "nothing blocking"},
		})
	}()

	state, err := following(store, name).Until(context.Background(), sink)
	if err != nil {
		t.Fatalf("following a log that rotated under the follower: %v", err)
	}
	if state != runstate.Clean {
		t.Errorf("the follower stopped at %s, want %s", state, runstate.Clean)
	}

	rotations := 0
	ended := 0
	for _, event := range delivered {
		switch event.Kind {
		case runlog.EventRotated:
			rotations++
		case "run-ended":
			ended++
		}
	}
	if rotations == 0 {
		t.Fatalf("this test needs the log to rotate under the follower, and it delivered %d events with no rotation", len(delivered))
	}
	if ended != 1 {
		t.Errorf("the run ended once and the follower delivered that %d times", ended)
	}
}

func TestWatchingImmediatelyAfterStarting(t *testing.T) {
	store, name := run(t, nil)

	waited := make(chan runstate.State, 1)
	sink := follow.Sink{Waiting: func(state runstate.State) error {
		waited <- state
		return nil
	}}

	go func() {
		time.Sleep(200 * time.Millisecond)
		write(t, store, name, runlog.Event{
			At:      time.Now().UTC(),
			Kind:    "run-ended",
			Subject: name,
			Fields:  map[string]any{"state": string(runstate.Clean), "reason": "nothing blocking"},
		})
	}()

	state, err := following(store, name).Until(context.Background(), sink)
	if err != nil {
		t.Fatalf("following a run whose log has not appeared yet: %v", err)
	}
	if state != runstate.Clean {
		t.Errorf("the follower stopped at %s, want %s", state, runstate.Clean)
	}

	select {
	case announced := <-waited:
		if announced != runstate.Reviewing {
			t.Errorf("the follower said the run was %s while it waited, want %s", announced, runstate.Reviewing)
		}
	default:
		t.Error("the follower waited for the log without saying that the run was still starting")
	}
}

func TestWaitingIsBounded(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		store, name := run(t, nil)

		started := time.Now()
		_, err := following(store, name).Until(context.Background(), follow.Sink{})
		if err == nil {
			t.Fatal("a run that wrote no event log and has no reviewer alive was followed without complaint")
		}
		if waited := time.Since(started); waited > 2*follow.DefaultPatience {
			t.Errorf("the follower waited %s, and its patience is %s", waited, follow.DefaultPatience)
		}
		if !strings.Contains(err.Error(), name) {
			t.Errorf("the refusal does not name the run: %v", err)
		}
	})
}

func TestARunIdentifierThatDoesNotExist(t *testing.T) {
	store := runstore.New(t.TempDir())

	_, err := following(store, "20200101T000000Z-deadbeef").Until(context.Background(), follow.Sink{})
	if err == nil {
		t.Fatal("following a run that does not exist was not refused")
	}
	if !strings.Contains(err.Error(), "20200101T000000Z-deadbeef") {
		t.Errorf("the refusal does not name the identifier asked for: %v", err)
	}
}

func TestHeldReadsWhatTheLogHoldsSoFar(t *testing.T) {
	store, name := run(t, nil)

	empty, err := follow.Held(store, name)
	if err != nil {
		t.Fatalf("reading a log that does not exist yet: %v", err)
	}
	if len(empty) != 0 {
		t.Errorf("a run with no log holds %d events", len(empty))
	}

	write(t, store, name, runlog.Event{At: time.Now().UTC(), Kind: "run-started", Subject: name})
	held, err := follow.Held(store, name)
	if err != nil {
		t.Fatalf("reading the log: %v", err)
	}
	if len(held) != 1 || held[0].Kind != "run-started" {
		t.Errorf("the log holds %+v, want the one event that was written", held)
	}
}

func following(store runstore.Store, name string) follow.Following {
	return follow.Following{Store: store, Name: name, Poll: 10 * time.Millisecond}
}

func run(t *testing.T, end func(*runstore.Run)) (runstore.Store, string) {
	t.Helper()

	store := runstore.New(filepath.Join(t.TempDir(), "runs"))
	created, err := store.Create(runstore.Start{
		StartedAt:      time.Now(),
		RepositoryRoot: repository,
		Base:           "a9a3207995f002632584efb33312c14123de4560",
		Snapshot:       "4445282d62930b67e699165cda6a212f8c2e47ab",
	})
	if err != nil {
		t.Fatalf("creating a run: %v", err)
	}

	if end == nil {
		if err := created.Advance(runstate.Reviewing); err != nil {
			t.Fatalf("advancing the run: %v", err)
		}
	} else {
		end(created)
	}
	if err := store.Save(created); err != nil {
		t.Fatalf("saving the run: %v", err)
	}
	return store, created.Name
}

func load(t *testing.T, store runstore.Store, name string) *runstore.Run {
	t.Helper()

	run, err := store.Load(name)
	if err != nil {
		t.Fatalf("reading the run %s: %v", name, err)
	}
	return run
}

func write(t *testing.T, store runstore.Store, name string, event runlog.Event) {
	t.Helper()
	writeWith(t, store, name, 0, event)
}

func writeWith(t *testing.T, store runstore.Store, name string, ceiling int64, event runlog.Event) {
	t.Helper()

	writer, err := runlog.Open(filepath.Join(store.Path(name), runlog.FileName), ceiling)
	if err != nil {
		t.Fatalf("opening the event log: %v", err)
	}
	defer writer.Close()

	if err := writer.Record(event); err != nil {
		t.Fatalf("recording an event: %v", err)
	}
}
