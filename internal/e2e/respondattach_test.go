package e2e_test

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/testing/harness"
	"github.com/socrates-tools/socrates/review/prompt"
	"github.com/socrates-tools/socrates/review/runstate"
)

const screenshotBytes = "\x89PNG\r\n\x1a\nwhat the person saw on their screen"

func attachable(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
	return path
}

func TestAnAttachedFileReachesTheNextRoundsReviewer(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	name := parkedRun(t, session, dialogueScenario(patchableFindings, nothingFoundAgain))

	screenshot := attachable(t, "retry.png", screenshotBytes)
	trace := attachable(t, "trace.txt", "the trace the person captured")

	respond(t, session, name, "--finding", "f1", "--fixed", "--attach", screenshot, "--attach", trace)
	respond(t, session, name, "--finding", "f2", "--rejected", "--reason", "there is no retry loop in this function at all", "--attach", screenshot)

	recorded, err := session.Round(name, 1).Responses()
	if err != nil {
		t.Fatalf("reading the responses back: %v", err)
	}
	rejection, _ := recorded.For("f2")
	fix, _ := recorded.For("f1")

	t.Run("The copy survives the original being deleted", func(t *testing.T) {
		if err := os.Remove(screenshot); err != nil {
			t.Fatalf("deleting the original: %v", err)
		}
		if len(rejection.Attachments) != 1 {
			t.Fatalf("the rejection records %+v", rejection.Attachments)
		}
		if held := session.File(rejection.Attachments[0].Path); held != screenshotBytes {
			t.Errorf("the copy holds %q", held)
		}
		if len(fix.Attachments) != 2 || session.File(fix.Attachments[1].Path) != "the trace the person captured" {
			t.Errorf("the fix records %+v", fix.Attachments)
		}
	})

	out, failures, code := session.Socrates("continue", name, "--wait")
	if code != 0 {
		t.Fatalf("socrates review continue --wait exited %d: %s\n%s", code, failures, out)
	}
	if run := session.Run(name); run.State != runstate.Clean {
		t.Fatalf("the run ended as %s: %s", run.State, reason(run))
	}

	t.Run("The next round's prompt names each file's path and type for its finding", func(t *testing.T) {
		invocations := session.Simulator.Invocations()
		if len(invocations) != 2 {
			t.Fatalf("the reviewer was invoked %d times, want twice", len(invocations))
		}
		block, present := prompt.BodyOf(prompt.Sections(invocations[1].Prompt), prompt.PriorHeading)
		if !present {
			t.Fatalf("the second round was given no prior-round block:\n%s", invocations[1].Prompt)
		}
		for _, carried := range []string{
			rejection.Attachments[0].Path, "image/png",
			fix.Attachments[0].Path, fix.Attachments[1].Path, "text/plain",
		} {
			if !strings.Contains(block, carried) {
				t.Errorf("the prior-round block does not carry %q:\n%s", carried, block)
			}
		}
	})
}

func TestAnAttachmentIsAcceptedWithEveryResolutionButForwarded(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	name := theirChangeReviewed(t, session)
	picture := attachable(t, "shown.png", screenshotBytes)

	for finding, resolution := range map[string][]string{
		"f1": {"--rejected", "--reason", "the reviewer misread it"},
		"f2": {"--dropped", "--reason", "true and not worth raising"},
	} {
		printed, failures, code := session.Socrates(append([]string{"respond", name, "--finding", finding, "--attach", picture}, resolution...)...)
		if code != 4 {
			t.Errorf("respond %s %s --attach exited %d: %s\n%s", finding, resolution[0], code, failures, printed)
		}
	}

	t.Run("A forward with an attachment is refused naming why and records nothing", func(t *testing.T) {
		before := session.File(session.Round(name, 1).ResponsesPath())
		refused := refuse(t, session, "respond", name, "--finding", "f3", "--forwarded", "--attach", picture)
		for _, named := range []string{"forward", "pull request", "file"} {
			if !strings.Contains(refused, named) {
				t.Errorf("the refusal does not name %q: %s", named, refused)
			}
		}
		if after := session.File(session.Round(name, 1).ResponsesPath()); after != before {
			t.Errorf("a refused forward recorded something:\n%s", after)
		}
		copies, _ := filepath.Glob(filepath.Join(session.Round(name, 1).Directory, "attachments", "*", "*"))
		if len(copies) != 2 {
			t.Errorf("the refused forward left %d copies where two were recorded: %v", len(copies), copies)
		}
	})
}

func TestAnAttachmentThatCannotBeRecordedIsRefused(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	name := parkedRun(t, session, dialogueScenario(patchableFindings, nothingFoundAgain))

	t.Run("An eleventh --attach is refused naming the ceiling", func(t *testing.T) {
		arguments := []string{"respond", name, "--finding", "f2", "--rejected", "--reason", "no retry loop here"}
		for position := range 11 {
			arguments = append(arguments, "--attach", attachable(t, "f"+strconv.Itoa(position)+".txt", "x"))
		}
		refused := refuse(t, session, arguments...)
		if !strings.Contains(refused, "10") {
			t.Errorf("the refusal does not name the ceiling of ten: %s", refused)
		}
		if _, err := os.Stat(session.Round(name, 1).ResponsesPath()); !os.IsNotExist(err) {
			t.Error("a refused response was recorded")
		}
	})

	t.Run("A path that does not exist is refused naming it", func(t *testing.T) {
		missing := filepath.Join(t.TempDir(), "never-written.png")
		refused := refuse(t, session, "respond", name, "--finding", "f2", "--rejected", "--reason", "no retry loop here", "--attach", missing)
		if !strings.Contains(refused, missing) {
			t.Errorf("the refusal does not name %s: %s", missing, refused)
		}
		if _, err := os.Stat(session.Round(name, 1).ResponsesPath()); !os.IsNotExist(err) {
			t.Error("a refused response was recorded")
		}
	})

}
