package runstate

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTheSixStatesAndWhichAreTerminal(t *testing.T) {
	want := map[State]bool{
		Starting:  false,
		Reviewing: false,
		Clean:     true,
		Escalated: true,
		Failed:    true,
		Cancelled: true,
	}

	got := map[State]bool{}
	for _, state := range States() {
		got[state] = state.Terminal()
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("the states and their terminality differ (-want +got):\n%s", diff)
	}
	if State("reviewed").Valid() {
		t.Error("a state socrates does not have reads as valid")
	}
}
