package promotion_test

import (
	"testing"

	"github.com/socrates-devtool/socrates/internal/promotion"
	"github.com/socrates-devtool/socrates/internal/testing/source"
)

func TestArchivingAChangeWouldDropNoScenario(t *testing.T) {
	checked, err := promotion.Check(source.Root(t))
	if err != nil {
		t.Fatalf("checking what archiving would drop: %v", err)
	}
	if checked.Changes == 0 {
		t.Skip("no change is under way, so there is nothing to archive: what the rule does is held by the tests over a fabricated tree")
	}
	if checked.Substituted == 0 && checked.Added == 0 {
		t.Fatal("a change is under way and proposes no requirement this check can read, so this test proves nothing")
	}

	for _, drop := range checked.Drops {
		t.Errorf("%s substitutes %q and %s: %v", drop.Spec, drop.Requirement, drop.Condition, drop.Scenarios)
	}
}
