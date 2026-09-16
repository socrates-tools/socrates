package runstate_test

import (
	"testing"

	"github.com/socrates-tools/socrates/internal/exit"
	"github.com/socrates-tools/socrates/review/runstate"
)

func TestEachStateExitsTheCodeTheTableStates(t *testing.T) {
	for state, want := range map[runstate.State]int{
		runstate.Clean:     exit.NothingAsked,
		runstate.Escalated: exit.SomethingAsked,
		runstate.Failed:    exit.Failed,
		runstate.Starting:  exit.Working,
		runstate.Reviewing: exit.Working,
	} {
		if got := state.Code(); got != want {
			t.Errorf("%s exits %d, want %d", state, got, want)
		}
	}
}

func TestFailureIsNotEscalation(t *testing.T) {
	if runstate.Failed.Code() == runstate.Escalated.Code() {
		t.Fatalf("a failed run and an escalated run share the code %d, and they call for different actions", runstate.Failed.Code())
	}
}

func TestNoStateExitsWithTheCodeForAnUnusableCommand(t *testing.T) {
	for _, state := range runstate.States() {
		if state.Code() == exit.Unusable {
			t.Errorf("%s exits %d, which is the code reserved for a command that was not usable", state, exit.Unusable)
		}
	}
}

func TestACancelledRunIsNotReportedAsAReview(t *testing.T) {
	if code := runstate.Cancelled.Code(); code == exit.NothingAsked || code == exit.SomethingAsked {
		t.Fatalf("a cancelled run exits %d, which reports a review that never finished", code)
	}
}

func TestEveryStateHasACode(t *testing.T) {
	for _, state := range runstate.States() {
		code := state.Code()
		if code < exit.NothingAsked || code > exit.Working {
			t.Errorf("%s exits %d, which is outside the table", state, code)
		}
		if state.Terminal() == (code == exit.Working) {
			t.Errorf("%s is terminal=%v and exits %d", state, state.Terminal(), code)
		}
	}
}
