package coverage_test

import (
	"testing"

	"github.com/socrates-devtool/socrates/internal/coverage"
	"github.com/socrates-devtool/socrates/internal/testing/source"
)

func TestEveryScenarioUnderWayIsNamedByATest(t *testing.T) {
	root := source.Root(t)

	scenarios, err := coverage.Scenarios(root)
	if err != nil {
		t.Fatalf("reading the specifications: %v", err)
	}
	if len(scenarios) == 0 {
		t.Fatal("no scenario was found under openspec, so this test proves nothing")
	}

	claims, err := coverage.Claims(root)
	if err != nil {
		t.Fatalf("reading the tests: %v", err)
	}

	for _, gap := range coverage.Gaps(scenarios, claims) {
		if !gap.Underway {
			continue
		}
		t.Errorf("%s:%d promises %q and no test names it", gap.Spec, gap.Line, gap.Text)
	}
}
