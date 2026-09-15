//go:build unix

package agent_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/socrates-tools/socrates/internal/testing/simulator"
	"github.com/socrates-tools/socrates/socrates-review/agent"
	"github.com/socrates-tools/socrates/socrates-review/simulation"
)

func TestASilentAgentIsKilled(t *testing.T) {
	simulated := simulator.Install(t)
	simulated.Scenario("default:\n  behaviour: silence\n  hold: 5m\n  child_command: [sleep, \"60\"]\n")

	recorded := &collector{}
	claude := build(t, agent.Options{Deadline: 400 * time.Millisecond, Recorder: recorded})

	started := time.Now()
	if _, err := claude.Run(context.Background(), request()); err == nil {
		t.Fatal("a silent agent was waited on forever")
	} else if !strings.Contains(err.Error(), "400ms") && !strings.Contains(err.Error(), "deadline") {
		t.Errorf("the failure reads %q, without naming the deadline", err)
	}

	if waited := time.Since(started); waited > 10*time.Second {
		t.Errorf("the call took %s, so the deadline is not what ended it", waited)
	}

	reviewer := recorded.field(t, agent.EventStarted, "process")
	pid, ok := reviewer.(int)
	if !ok {
		t.Fatalf("the started event names the process as %T", reviewer)
	}
	waitUntilGone(t, pid, "the agent")

	child := simulated.OnlyInvocation().Child
	if child == 0 {
		t.Fatal("the simulator started no child, so this test does not reach the process group")
	}
	waitUntilGone(t, child, "the child the agent spawned")
}

func TestASlowButLiveAgentIsNotKilledEarly(t *testing.T) {
	simulated := simulator.Install(t)
	simulated.Scenario("default:\n  events: 5\n  interval: 40ms\n")

	claude := build(t, agent.Options{Deadline: 10 * time.Second})

	answer, err := claude.Run(context.Background(), request())
	if err != nil {
		t.Fatalf("a slow but live agent was killed: %v", err)
	}
	if len(answer.Structured) == 0 {
		t.Error("the slow agent's structured output was lost")
	}
	if simulated.OnlyInvocation().Prompt == "" {
		t.Error("the agent was never invoked")
	}
}

func TestACancelledCallStopsTheAgent(t *testing.T) {
	simulated := simulator.Install(t)
	simulated.Scenario("default:\n  behaviour: silence\n  hold: 5m\n  child_command: [sleep, \"60\"]\n")

	recorded := &collector{}
	claude := build(t, agent.Options{Deadline: time.Minute, Recorder: recorded})

	ctx, stop := context.WithCancel(context.Background())
	spawned := make(chan simulation.Invocation, 1)
	go func() {
		defer stop()
		if invocation, found := simulated.WaitForChild(30 * time.Second); found {
			spawned <- invocation
		}
	}()

	if _, err := claude.Run(ctx, request()); err == nil {
		t.Fatal("a cancelled call succeeded")
	}

	select {
	case invocation := <-spawned:
		waitUntilGone(t, invocation.Child, "the child the agent spawned")
	default:
		t.Fatal("the simulator started no child, so this test does not reach the process group")
	}
}

func waitUntilGone(t *testing.T, pid int, what string) {
	t.Helper()

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if gone(pid) {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Errorf("%s (process %d) is still running", what, pid)
}

func gone(pid int) bool {
	if err := syscall.Kill(pid, 0); err != nil {
		return true
	}

	stat, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return true
	}
	closing := strings.LastIndexByte(string(stat), ')')
	if closing < 0 {
		return false
	}
	fields := strings.Fields(string(stat)[closing+1:])
	return len(fields) > 0 && fields[0] == "Z"
}
