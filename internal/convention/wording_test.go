package convention

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-devtool/socrates/internal/testing/source"
)

func TestARefusalNamesTheNextAction(t *testing.T) {
	wording, err := UserFacingStrings(source.Root(t))
	if err != nil {
		t.Fatalf("reading what socrates says: %v", err)
	}

	for _, actionless := range wording.Actionless {
		t.Errorf("%s refuses without naming what the person can do about it", actionless)
	}
}

func TestNoEmDashReachesAPerson(t *testing.T) {
	wording, err := UserFacingStrings(source.Root(t))
	if err != nil {
		t.Fatalf("reading what socrates says: %v", err)
	}

	for _, dashed := range wording.EmDashes {
		t.Errorf("%s carries an em dash", dashed)
	}
}

func TestNothingSocratesPrintsNamesTheShortName(t *testing.T) {
	wording, err := UserFacingStrings(source.Root(t))
	if err != nil {
		t.Fatalf("reading what socrates says: %v", err)
	}

	for _, named := range wording.ShortNamed {
		t.Errorf("%s names the short name, and what socrates prints names socrates", named)
	}
}

func TestAStringNamingTheShortNameIsReported(t *testing.T) {
	root := t.TempDir()
	write(t, root, filepath.Join("internal", "thing", "messages.go"), `package thing

const named = "run `+ShortName+` --help to see what it does"
`)

	wording, err := UserFacingStrings(root)
	if err != nil {
		t.Fatalf("reading the fabricated package: %v", err)
	}

	want := []string{filepath.Join("internal", "thing", "messages.go") + ":3"}
	if diff := cmp.Diff(want, wording.ShortNamed); diff != "" {
		t.Errorf("the reported uses of the short name differ (-want +got):\n%s", diff)
	}
}

func TestEveryUserFacingStringIsInItsPackagesRegistry(t *testing.T) {
	wording, err := UserFacingStrings(source.Root(t))
	if err != nil {
		t.Fatalf("reading what socrates says: %v", err)
	}

	for _, inline := range wording.Inline {
		t.Errorf("%s states a message where it is used: move it into the package's %s", inline, RegistryFile)
	}
}

func TestAMessageStatedAtTheCallSiteIsReported(t *testing.T) {
	root := t.TempDir()
	write(t, root, filepath.Join("internal", "thing", "thing.go"), `package thing

import "fmt"

func Do() error {
	return fmt.Errorf("the thing could not be done: %w", nil)
}
`)

	wording, err := UserFacingStrings(root)
	if err != nil {
		t.Fatalf("reading the fabricated package: %v", err)
	}

	want := []string{filepath.Join("internal", "thing", "thing.go") + ":6"}
	if diff := cmp.Diff(want, wording.Inline); diff != "" {
		t.Errorf("the reported messages differ (-want +got):\n%s", diff)
	}
}

func TestAMessageInTheRegistryIsNotReported(t *testing.T) {
	root := t.TempDir()
	write(t, root, filepath.Join("internal", "thing", "thing.go"), `package thing

import "fmt"

func Do() error {
	return fmt.Errorf(couldNotBeDone, nil)
}
`)
	write(t, root, filepath.Join("internal", "thing", "messages.go"), `package thing

const couldNotBeDone = "the thing could not be done: %w"
`)

	wording, err := UserFacingStrings(root)
	if err != nil {
		t.Fatalf("reading the fabricated package: %v", err)
	}
	if len(wording.Inline) != 0 {
		t.Errorf("a message held in the registry was reported: %v", wording.Inline)
	}
}

func TestAFormatWithNoProseIsNotReported(t *testing.T) {
	root := t.TempDir()
	write(t, root, filepath.Join("internal", "thing", "thing.go"), `package thing

import "fmt"

func Do() string {
	return fmt.Sprintf("%s/%d", "a", 1)
}
`)

	wording, err := UserFacingStrings(root)
	if err != nil {
		t.Fatalf("reading the fabricated package: %v", err)
	}
	if len(wording.Inline) != 0 {
		t.Errorf("a format string carrying no prose was reported: %v", wording.Inline)
	}
}

func TestARefusalWithNoNextActionIsReported(t *testing.T) {
	root := t.TempDir()
	write(t, root, filepath.Join("internal", "thing", "thing.go"), `package thing

import "github.com/socrates-devtool/socrates/internal/refusal"

func Do() error {
	return refusal.New(somethingWentWrong)
}
`)
	write(t, root, filepath.Join("internal", "thing", "messages.go"), `package thing

const somethingWentWrong = "the thing is in a state socrates cannot use"
`)

	wording, err := UserFacingStrings(root)
	if err != nil {
		t.Fatalf("reading the fabricated package: %v", err)
	}

	want := []string{filepath.Join("internal", "thing", "thing.go") + ":6"}
	if diff := cmp.Diff(want, wording.Actionless); diff != "" {
		t.Errorf("the reported refusals differ (-want +got):\n%s", diff)
	}
}

func TestARefusalThatNamesTheNextActionIsNotReported(t *testing.T) {
	root := t.TempDir()
	write(t, root, filepath.Join("internal", "thing", "thing.go"), `package thing

import "github.com/socrates-devtool/socrates/internal/refusal"

func Do() error {
	return refusal.New(somethingWentWrong)
}
`)
	write(t, root, filepath.Join("internal", "thing", "messages.go"), `package thing

const somethingWentWrong = "the thing is in a state socrates cannot use: fix the thing and run the command again"
`)

	wording, err := UserFacingStrings(root)
	if err != nil {
		t.Fatalf("reading the fabricated package: %v", err)
	}
	if len(wording.Actionless) != 0 {
		t.Errorf("a refusal naming what to do next was reported: %v", wording.Actionless)
	}
}

func TestAnEmDashIsReported(t *testing.T) {
	root := t.TempDir()
	write(t, root, filepath.Join("internal", "thing", "messages.go"), `package thing

const dashed = "socrates refuses — and here is why"
`)

	wording, err := UserFacingStrings(root)
	if err != nil {
		t.Fatalf("reading the fabricated package: %v", err)
	}

	want := []string{filepath.Join("internal", "thing", "messages.go") + ":3"}
	if diff := cmp.Diff(want, wording.EmDashes); diff != "" {
		t.Errorf("the reported em dashes differ (-want +got):\n%s", diff)
	}
}
