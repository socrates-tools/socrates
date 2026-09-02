package convention_test

import (
	"testing"

	"github.com/ktei/socrates/internal/convention"
	"github.com/ktei/socrates/internal/testing/source"
)

func TestOnlyATestAboutReloadingReloadsThePage(t *testing.T) {
	reloading, err := convention.TestsReloadingThePage(source.Root(t))
	if err != nil {
		t.Fatalf("reading which tests reload the page: %v", err)
	}

	for _, unnamed := range reloading.Unnamed {
		t.Errorf("%s reloads the page before it asserts, and a reloaded page proves the weaker of the two claims: assert on the page the person is looking at, or name the test in convention.TestsThatReloadThePage", unnamed)
	}
	for _, stale := range reloading.Stale {
		t.Errorf("convention.TestsThatReloadThePage names %s, which reloads nothing: remove the entry", stale)
	}
	if len(reloading.Named) == 0 {
		t.Error("no test reloads the page at all, so nothing holds reloading to remaining sufficient")
	}
}
