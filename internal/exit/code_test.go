package exit_test

import (
	"testing"

	"github.com/socrates-tools/socrates/internal/exit"
)

func TestEachCodeMeansOneThing(t *testing.T) {
	meanings := map[int]string{
		0: "nothing is asked of you",
		1: "something is asked of you",
		2: "the command failed",
		3: "the command was not usable",
		4: "the work is still running",
	}

	named := map[int]string{
		exit.NothingAsked:   "nothing is asked of you",
		exit.SomethingAsked: "something is asked of you",
		exit.Failed:         "the command failed",
		exit.Unusable:       "the command was not usable",
		exit.Working:        "the work is still running",
	}

	if len(named) != len(meanings) {
		t.Fatalf("two of the five names share a code: %v", named)
	}
	for code, meaning := range meanings {
		if named[code] != meaning {
			t.Errorf("code %d means %q, want %q", code, named[code], meaning)
		}
	}
}
