package guidance_test

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"github.com/ktei/socrates/internal/guidance"
	"github.com/ktei/socrates/internal/refusal"
	"github.com/ktei/socrates/internal/testing/gitrepo"
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

	read, err := guidance.At(context.Background(), repository.Root, base)
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

func TestAGuideOverTheCeilingIsRefused(t *testing.T) {
	repository := gitrepo.New(t)
	oversized := strings.Repeat("a convention nobody needs stated twice\n", 1+guidance.Ceiling/38)
	repository.Write(guidance.Path, oversized)
	base := repository.Commit("a guide over the ceiling")

	if len(oversized) <= guidance.Ceiling {
		t.Fatalf("the guide built for this test is %d bytes, which is not over the %d byte ceiling", len(oversized), guidance.Ceiling)
	}

	_, err := guidance.At(context.Background(), repository.Root, base)
	if err == nil {
		t.Fatal("a guide over the ceiling was carried")
	}
	if !refusal.Is(err) {
		t.Errorf("a guide over the ceiling failed rather than refusing: %v", err)
	}
	for _, named := range []string{guidance.Path, strconv.Itoa(len(oversized)), strconv.Itoa(guidance.Ceiling)} {
		if !strings.Contains(err.Error(), named) {
			t.Errorf("the refusal does not name %q: %s", named, err)
		}
	}
}

func TestAGuideAtTheCeilingIsCarriedWhole(t *testing.T) {
	repository := gitrepo.New(t)
	exact := strings.Repeat("x", guidance.Ceiling)
	repository.Write(guidance.Path, exact)
	base := repository.Commit("a guide at the ceiling")

	read := at(t, repository, base)
	if len(read.Text) != guidance.Ceiling {
		t.Errorf("a guide of exactly the ceiling was carried as %d bytes, want %d", len(read.Text), guidance.Ceiling)
	}
}

func TestOneByteOverTheCeilingIsRefused(t *testing.T) {
	repository := gitrepo.New(t)
	repository.Write(guidance.Path, strings.Repeat("x", guidance.Ceiling+1))
	base := repository.Commit("a guide one byte over the ceiling")

	if _, err := guidance.At(context.Background(), repository.Root, base); err == nil {
		t.Error("a guide one byte over the ceiling was carried")
	}
}

func TestTheCeilingIsMeasuredBeforeNeutralisation(t *testing.T) {
	repository := gitrepo.New(t)
	collapsible := strings.Repeat(strings.Repeat(" ", 64)+"a line neutralisation collapses\n", 1+guidance.Ceiling/96)
	repository.Write(guidance.Path, collapsible)
	base := repository.Commit("a guide that only crosses the ceiling before it is neutralised")

	if len(collapsible) <= guidance.Ceiling {
		t.Fatalf("the guide built for this test is %d bytes, which is not over the %d byte ceiling", len(collapsible), guidance.Ceiling)
	}

	if _, err := guidance.At(context.Background(), repository.Root, base); err == nil {
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

	read, err := guidance.At(context.Background(), repository.Root, base)
	if err != nil {
		t.Fatalf("reading the guide at %s: %v", base, err)
	}
	return read
}
