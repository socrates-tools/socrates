package agent_test

import (
	"context"
	"strings"
	"sync"
	"testing"

	"github.com/socrates-tools/socrates/internal/agent"
	"github.com/socrates-tools/socrates/internal/runlog"
	"github.com/socrates-tools/socrates/internal/testing/simulator"
)

type collector struct {
	mutex    sync.Mutex
	recorded []runlog.Event
}

func (c *collector) Record(event runlog.Event) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.recorded = append(c.recorded, event)
	return nil
}

func (c *collector) events(kind string) []runlog.Event {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var found []runlog.Event
	for _, event := range c.recorded {
		if event.Kind == kind {
			found = append(found, event)
		}
	}
	return found
}

func (c *collector) field(t *testing.T, kind, name string) any {
	t.Helper()

	events := c.events(kind)
	if len(events) != 1 {
		t.Fatalf("%d %s events were recorded, want one", len(events), kind)
	}
	value, present := events[0].Fields[name]
	if !present {
		t.Fatalf("the %s event holds %v, without %s", kind, events[0].Fields, name)
	}
	return value
}

func TestWhatTheAgentCostIsRecorded(t *testing.T) {
	simulated := simulator.Install(t)
	simulated.Scenario(`
default:
  model: claude-opus-5
  input_tokens: 1200
  output_tokens: 340
`)

	recorded := &collector{}
	claude := build(t, agent.Options{
		Routing:  agent.Configuration{agent.Reviewer: {Effort: "xhigh"}},
		Recorder: recorded,
	})

	if _, err := claude.Run(context.Background(), request()); err != nil {
		t.Fatalf("running the reviewer: %v", err)
	}

	if model := recorded.field(t, agent.EventCompleted, "model"); model != "claude-opus-5" {
		t.Errorf("the completed event names the model %v, want the one that answered", model)
	}
	if effort := recorded.field(t, agent.EventCompleted, "effort"); effort != "xhigh" {
		t.Errorf("the completed event names the effort %v", effort)
	}
	if tokens := recorded.field(t, agent.EventCompleted, "input_tokens"); tokens != 1200 {
		t.Errorf("the completed event reports %v input tokens", tokens)
	}
	if tokens := recorded.field(t, agent.EventCompleted, "output_tokens"); tokens != 340 {
		t.Errorf("the completed event reports %v output tokens", tokens)
	}

	if started := recorded.events(agent.EventStarted); len(started) != 1 {
		t.Errorf("%d started events were recorded, want one", len(started))
	}
	if failed := recorded.events(agent.EventFailed); len(failed) != 0 {
		t.Errorf("a successful call recorded %d failures", len(failed))
	}
}

func TestAFailedCallIsRecordedAsAFailure(t *testing.T) {
	simulated := simulator.Install(t)
	simulated.Scenario("default:\n  behaviour: no-result\n")

	recorded := &collector{}
	claude := build(t, agent.Options{Recorder: recorded})

	if _, err := claude.Run(context.Background(), request()); err == nil {
		t.Fatal("a stream with no result event was accepted")
	}

	reason := recorded.field(t, agent.EventFailed, "reason")
	stated, isText := reason.(string)
	if !isText || !strings.Contains(stated, "result event") {
		t.Errorf("the failed event reads %v, without why it failed", reason)
	}
	if completed := recorded.events(agent.EventCompleted); len(completed) != 0 {
		t.Errorf("a failed call recorded %d completions", len(completed))
	}
}

func TestARunWithNoRecorderStillRuns(t *testing.T) {
	simulator.Install(t)

	if _, err := build(t, agent.Options{}).Run(context.Background(), request()); err != nil {
		t.Fatalf("running with no recorder: %v", err)
	}
}
