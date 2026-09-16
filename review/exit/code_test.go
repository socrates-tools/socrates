package exit_test

import (
	"testing"

	"github.com/socrates-tools/socrates/review/exit"
	"github.com/socrates-tools/socrates/review/runstate"
)

func TestEachVerdictHasItsOwnCode(t *testing.T) {
	for state, want := range map[runstate.State]int{
		runstate.Clean:     0,
		runstate.Escalated: 1,
		runstate.Failed:    2,
		runstate.Starting:  4,
		runstate.Reviewing: 4,
	} {
		if got := exit.For(state); got != want {
			t.Errorf("%s exits %d, want %d", state, got, want)
		}
	}
}

func TestFailureIsNotEscalation(t *testing.T) {
	if exit.For(runstate.Failed) == exit.For(runstate.Escalated) {
		t.Fatalf("a failed run and an escalated run share the code %d, and they call for different actions", exit.For(runstate.Failed))
	}
}

func TestAUsageErrorIsNotAVerdict(t *testing.T) {
	for _, state := range runstate.States() {
		if exit.For(state) == exit.Usage {
			t.Errorf("%s exits %d, which is the code reserved for a usage error", state, exit.Usage)
		}
	}
}

func TestACancelledRunIsNotReportedAsAReview(t *testing.T) {
	if code := exit.For(runstate.Cancelled); code == exit.Clean || code == exit.Escalated {
		t.Fatalf("a cancelled run exits %d, which reports a review that never finished", code)
	}
}

func TestEveryStateHasACode(t *testing.T) {
	for _, state := range runstate.States() {
		code := exit.For(state)
		if code < 0 || code > exit.Working {
			t.Errorf("%s exits %d, which is outside the table", state, code)
		}
		if state.Terminal() == (code == exit.Working) {
			t.Errorf("%s is terminal=%v and exits %d", state, state.Terminal(), code)
		}
	}
}
