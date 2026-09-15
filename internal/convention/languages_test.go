package convention_test

import (
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/internal/convention"
)

const registering = `var t=r();t.registerLanguage("go",e),t.registerLanguage("diff",n),t.registerLanguage("csharp",a);`

func TestTheLanguagesABundleRegistersAreReadFromIt(t *testing.T) {
	want := []string{"csharp", "diff", "go"}
	if diff := cmp.Diff(want, convention.LanguagesTheBundleRegisters(registering)); diff != "" {
		t.Errorf("what a minified bundle is read as registering differs (-want +got):\n%s", diff)
	}
}

func TestALanguageOnOneSideAndNotTheOtherIsReported(t *testing.T) {
	registered := convention.LanguagesTheBundleRegisters(registering)

	for _, held := range []struct {
		named          string
		stated         []string
		statedOnly     []string
		registeredOnly []string
	}{
		{
			named:  "the same set on both sides",
			stated: []string{"csharp", "diff", "go"},
		},
		{
			named:      "one stated that the bundle does not register",
			stated:     []string{"csharp", "diff", "go", "prolog"},
			statedOnly: []string{"prolog"},
		},
		{
			named:          "one registered that nothing states",
			stated:         []string{"csharp", "go"},
			registeredOnly: []string{"diff"},
		},
	} {
		t.Run(held.named, func(t *testing.T) {
			statedOnly, registeredOnly := convention.LanguagesOnOneSideOnly(held.stated, registered)

			if diff := cmp.Diff(held.statedOnly, statedOnly); diff != "" {
				t.Errorf("what is stated and not registered differs (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(held.registeredOnly, registeredOnly); diff != "" {
				t.Errorf("what is registered and not stated differs (-want +got):\n%s", diff)
			}
		})
	}
}

func TestABundleRegisteringNothingIsReadAsRegisteringNothing(t *testing.T) {
	if named := convention.LanguagesTheBundleRegisters("no library here at all"); len(named) != 0 {
		t.Errorf("a bundle registering nothing is read as registering %v", named)
	}
	if !slices.IsSorted(convention.LanguagesTheBundleRegisters(registering)) {
		t.Error("what a bundle registers comes back unsorted, and two sets are compared as sequences")
	}
}
