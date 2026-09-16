package guidance_test

import (
	"context"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/testing/gitrepo"
	"github.com/socrates-tools/socrates/review/guidance"
)

func TestTheGuideIsReadAtTheBaseCommit(t *testing.T) {
	repository := gitrepo.New(t)
	repository.Write(guidance.Path, "what the base says\n")
	base := repository.Commit("a guide")
	repository.Write(guidance.Path, "what the change says\n")

	read := at(t, repository, base)

	if !strings.Contains(read.Text, "what the base says") {
		t.Errorf("the guide read at the base answered %q", read.Text)
	}
	if strings.Contains(read.Text, "what the change says") {
		t.Errorf("the guide read at the base carried what the change said: %q", read.Text)
	}
}

func TestAGuideHoldingNothing(t *testing.T) {
	repository := gitrepo.New(t)
	repository.Write(guidance.Path, "\n   \n\t\n")
	base := repository.Commit("a guide holding nothing")

	if read := at(t, repository, base); read.Text != "" {
		t.Errorf("a guide holding only whitespace answered %q, and there is nothing there to tell a reviewer", read.Text)
	}
}

func TestARepositoryWithNoGuide(t *testing.T) {
	repository := gitrepo.New(t)
	base := repository.Head()

	read, err := guidance.At(context.Background(), repository.Root, base, 0)
	if err != nil {
		t.Fatalf("a repository with no guide answered an error: %v", err)
	}
	if read.Text != "" {
		t.Errorf("a repository with no guide answered %q", read.Text)
	}
	if read.AtTheBase {
		t.Error("a repository with no guide reads as one whose base carries a guide")
	}
	if outcome := read.Reached(nil); outcome != guidance.None {
		t.Errorf("a repository with no guide reached %q, want %q", outcome, guidance.None)
	}
}

func TestAGuideOverTheCeilingIsNotCarried(t *testing.T) {
	repository := gitrepo.New(t)
	oversized := strings.Repeat("a convention nobody needs stated twice\n", 1+guidance.DefaultCeiling/38)
	repository.Write(guidance.Path, oversized)
	base := repository.Commit("a guide over the ceiling")

	if len(oversized) <= guidance.DefaultCeiling {
		t.Fatalf("the guide built for this test is %d bytes, which is not over the %d byte ceiling", len(oversized), guidance.DefaultCeiling)
	}

	read := at(t, repository, base)
	if read.Text != "" {
		t.Errorf("a guide over the ceiling answered %d bytes of text", len(read.Text))
	}
	if !read.TooLarge {
		t.Error("a guide over the ceiling does not read as one")
	}
	if read.Size != int64(len(oversized)) {
		t.Errorf("the guide reads as %d bytes, and it is %d on disk", read.Size, len(oversized))
	}
	if outcome := read.Reached(nil); outcome != guidance.OverTheCeiling {
		t.Errorf("a guide over the ceiling reached %q, want %q", outcome, guidance.OverTheCeiling)
	}
}

func TestAGuideAtTheCeilingIsCarriedWhole(t *testing.T) {
	repository := gitrepo.New(t)
	exact := strings.Repeat("x", guidance.DefaultCeiling)
	repository.Write(guidance.Path, exact)
	base := repository.Commit("a guide at the ceiling")

	read := at(t, repository, base)
	if len(read.Text) != guidance.DefaultCeiling {
		t.Errorf("a guide of exactly the ceiling was carried as %d bytes, want %d", len(read.Text), guidance.DefaultCeiling)
	}
}

func TestOneByteOverTheCeilingIsNotCarried(t *testing.T) {
	repository := gitrepo.New(t)
	repository.Write(guidance.Path, strings.Repeat("x", guidance.DefaultCeiling+1))
	base := repository.Commit("a guide one byte over the ceiling")

	if read := at(t, repository, base); read.Text != "" {
		t.Error("a guide one byte over the ceiling was carried")
	}
}

func TestTheCeilingIsMeasuredBeforeNeutralisation(t *testing.T) {
	repository := gitrepo.New(t)
	collapsible := strings.Repeat(strings.Repeat(" ", 64)+"a line neutralisation collapses\n", 1+guidance.DefaultCeiling/96)
	repository.Write(guidance.Path, collapsible)
	base := repository.Commit("a guide that only crosses the ceiling before it is neutralised")

	if len(collapsible) <= guidance.DefaultCeiling {
		t.Fatalf("the guide built for this test is %d bytes, which is not over the %d byte ceiling", len(collapsible), guidance.DefaultCeiling)
	}

	if read := at(t, repository, base); read.Text != "" {
		t.Error("a guide over the ceiling was carried because neutralisation would have brought it under one")
	}
}

func TestAGuideOnlyInTheChangeIsReachedAsSuch(t *testing.T) {
	repository := gitrepo.New(t)
	base := repository.Head()
	repository.Write(guidance.Path, "what the change says\n")

	read := at(t, repository, base)
	if read.Text != "" {
		t.Errorf("a guide the base does not carry answered %q", read.Text)
	}
	if outcome := read.Reached([]string{guidance.Path}); outcome != guidance.OnlyInTheChange {
		t.Errorf("a guide only in the change reached %q, want %q", outcome, guidance.OnlyInTheChange)
	}
}

func TestAGuideAtTheBaseIsReachedAsCarried(t *testing.T) {
	repository := gitrepo.New(t)
	repository.Write(guidance.Path, "one convention\n")
	base := repository.Commit("a guide")

	read := at(t, repository, base)
	if outcome := read.Reached([]string{guidance.Path}); outcome != guidance.Carried {
		t.Errorf("a guide at the base that the change also touches reached %q, want %q", outcome, guidance.Carried)
	}
}

func at(t *testing.T, repository *gitrepo.Repository, base string) guidance.Guide {
	t.Helper()

	read, err := guidance.At(context.Background(), repository.Root, base, 0)
	if err != nil {
		t.Fatalf("reading the guide at %s: %v", base, err)
	}
	return read
}

func TestAStatedGuidanceCeilingIsWhatTheGuideIsHeldTo(t *testing.T) {
	repository := gitrepo.New(t)
	guide := strings.Repeat("a convention\n", 200)
	repository.Write(guidance.Path, guide)
	base := repository.Commit("a guide between two ceilings")

	if len(guide) >= guidance.DefaultCeiling {
		t.Fatalf("the guide built for this test is %d bytes, and it has to be under the %d byte default", len(guide), guidance.DefaultCeiling)
	}

	under, err := guidance.At(context.Background(), repository.Root, base, int64(len(guide)-1))
	if err != nil {
		t.Fatalf("reading the guide under a stated ceiling: %v", err)
	}
	if !under.TooLarge {
		t.Error("a guide over the stated ceiling was carried, and the default is what held it")
	}

	over, err := guidance.At(context.Background(), repository.Root, base, int64(len(guide)))
	if err != nil {
		t.Fatalf("reading the guide at a stated ceiling: %v", err)
	}
	if over.TooLarge {
		t.Error("a guide exactly at the stated ceiling was not carried")
	}
}

func TestAnUnstatedCeilingIsTheDefault(t *testing.T) {
	repository := gitrepo.New(t)
	repository.Write(guidance.Path, strings.Repeat("x", guidance.DefaultCeiling+1))
	base := repository.Commit("a guide over the default")

	read := at(t, repository, base)
	if read.Ceiling != guidance.DefaultCeiling {
		t.Errorf("an unstated ceiling read as %d, want the default %d", read.Ceiling, guidance.DefaultCeiling)
	}
	if !read.TooLarge {
		t.Error("a guide over the default ceiling was carried when nothing stated a ceiling")
	}
}
