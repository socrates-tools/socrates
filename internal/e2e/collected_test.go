package e2e_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/agora/surface"
	"github.com/socrates-tools/socrates/internal/exit"
)

const aPollThatWaitedWouldStillBeWaiting = 10 * time.Second

func TestAnAnswerOnePollPrintedIsPrintedAgainByPollLast(t *testing.T) {
	t.Parallel()
	const inTheirOwnWords = "take the lock, and say why in the commit"

	held := standing(t)
	held.serve(oneDecisionDocument)
	held.answer(surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}, Words: inTheirOwnWords}}})

	first, failures, code := held.poll()
	if code != exit.NothingAsked {
		t.Fatalf("agora poll exited %d: %s", code, failures)
	}
	if !strings.Contains(first, inTheirOwnWords) {
		t.Fatalf("agora poll did not print what the person wrote:\n%s", first)
	}

	started := time.Now()
	again, failures, code := held.poll("--last")
	if code != exit.NothingAsked {
		t.Fatalf("agora poll --last exited %d: %s", code, failures)
	}
	if took := time.Since(started); took > aPollThatWaitedWouldStillBeWaiting {
		t.Errorf("agora poll --last took %s, and it never waits", took)
	}
	if diff := cmp.Diff(first, again); diff != "" {
		t.Errorf("agora poll --last printed otherwise than the poll that collected the answer (-collected +again):\n%s", diff)
	}
}

func TestACollectedAnswerIsWrittenToTheSessionBeforePollPrintsIt(t *testing.T) {
	t.Parallel()
	const inTheirOwnWords = "leave it, the cache is rebuilt hourly"

	held := standing(t)
	held.serve(oneDecisionDocument)
	held.answer(surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{1}, Words: inTheirOwnWords}}})

	printed, failures, code := held.poll()
	if code != exit.NothingAsked {
		t.Fatalf("agora poll exited %d: %s", code, failures)
	}

	written, err := os.ReadFile(filepath.Join(held.surfaceHome(), "sessions", held.opened, "answers", "1", "answer.json"))
	if err != nil {
		t.Fatalf("the session holds no answer once agora poll printed one: %v", err)
	}
	var kept, handed surface.Answered
	if err := json.Unmarshal(written, &kept); err != nil {
		t.Fatalf("reading the answer the session holds: %v", err)
	}
	if err := json.Unmarshal([]byte(printed), &handed); err != nil {
		t.Fatalf("reading what agora poll printed: %v", err)
	}
	if diff := cmp.Diff(handed, kept); diff != "" {
		t.Errorf("the session holds otherwise than agora poll printed (-printed +held):\n%s", diff)
	}
}

func TestEveryAnswerCollectedInOneSessionStaysOnDiskInAFolderOfItsOwnInTheOrderCollected(t *testing.T) {
	t.Parallel()
	const inRoundOne = "the first round, in my own words"
	const inRoundTwo = "the second round, in my own words"

	held := standing(t)
	held.serve(oneDecisionDocument)

	for _, words := range []string{inRoundOne, inRoundTwo} {
		held.answer(surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}, Words: words}}})
		if _, failures, code := held.poll(); code != exit.NothingAsked {
			t.Fatalf("agora poll exited %d: %s", code, failures)
		}
	}

	for number, words := range map[string]string{"1": inRoundOne, "2": inRoundTwo} {
		written, err := os.ReadFile(filepath.Join(held.surfaceHome(), "sessions", held.opened, "answers", number, "answer.json"))
		if err != nil {
			t.Fatalf("the session holds no answer %s once two were collected: %v", number, err)
		}
		if !strings.Contains(string(written), words) {
			t.Errorf("answers/%s/answer.json holds %s, want the words %q", number, written, words)
		}
	}
}

func TestPollLastPrintsTheSecondAnswerCollectedAndNotTheFirstWithNoSurfaceRunning(t *testing.T) {
	t.Parallel()
	const inRoundOne = "the first round, in my own words"
	const inRoundTwo = "the second round, in my own words"

	held := standing(t)
	held.serve(oneDecisionDocument)
	for _, words := range []string{inRoundOne, inRoundTwo} {
		held.answer(surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}, Words: words}}})
		if _, failures, code := held.poll(); code != exit.NothingAsked {
			t.Fatalf("agora poll exited %d: %s", code, failures)
		}
	}
	if _, failures, code := held.run("stop"); code != exit.NothingAsked {
		t.Fatalf("agora stop exited %d: %s", code, failures)
	}

	printed, failures, code := held.poll("--last")
	if code != exit.NothingAsked {
		t.Fatalf("agora poll --last with no surface standing exited %d: %s", code, failures)
	}
	if !strings.Contains(printed, inRoundTwo) || strings.Contains(printed, inRoundOne) {
		t.Errorf("agora poll --last printed:\n%s\nwant the second answer alone", printed)
	}
}

func TestAnAnswerTheSessionCouldNotKeepStaysUncollectedAndTheNextPollCollectsIt(t *testing.T) {
	t.Parallel()
	const inTheirOwnWords = "the disk was full when I said this"

	held := standing(t)
	held.serve(oneDecisionDocument)
	held.answer(surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}, Words: inTheirOwnWords}}})

	blocking := filepath.Join(held.surfaceHome(), "sessions", held.opened, "answers")
	if err := os.WriteFile(blocking, []byte("a file where the folder of answers belongs"), 0o600); err != nil {
		t.Fatalf("making the session unable to keep an answer: %v", err)
	}
	if printed, _, code := held.poll(); code == exit.NothingAsked {
		t.Fatalf("agora poll succeeded though the session could not keep the answer: %s", printed)
	}

	if err := os.Remove(blocking); err != nil {
		t.Fatalf("letting the session keep answers again: %v", err)
	}
	printed, failures, code := held.poll()
	if code != exit.NothingAsked {
		t.Fatalf("the next agora poll exited %d: %s", code, failures)
	}
	if !strings.Contains(printed, inTheirOwnWords) {
		t.Errorf("the next agora poll did not collect the answer the session could not keep:\n%s", printed)
	}
}

func TestPollLastReadsTheAnswerOnceTheSurfaceHasStopped(t *testing.T) {
	t.Parallel()
	const inTheirOwnWords = "the surface is gone and this is not"

	held := standing(t)
	held.serve(oneDecisionDocument)
	held.answer(surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}, Words: inTheirOwnWords}}})
	if _, failures, code := held.poll(); code != exit.NothingAsked {
		t.Fatalf("agora poll exited %d: %s", code, failures)
	}
	if _, failures, code := held.run("stop"); code != exit.NothingAsked {
		t.Fatalf("agora stop exited %d: %s", code, failures)
	}

	printed, failures, code := held.poll("--last")
	if code != exit.NothingAsked {
		t.Fatalf("agora poll --last with no surface standing exited %d: %s", code, failures)
	}
	if !strings.Contains(printed, inTheirOwnWords) {
		t.Errorf("agora poll --last with no surface standing printed:\n%s", printed)
	}
}

func TestPollLastOnASessionNothingWasCollectedFromIsRefusedWithoutWaiting(t *testing.T) {
	t.Parallel()
	held := standing(t)
	held.serve(oneDecisionDocument)
	held.answer(surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}}}})

	started := time.Now()
	printed, failures, code := held.poll("--last")
	if code != exit.Unusable {
		t.Fatalf("agora poll --last before anything was collected exited %d, want %d: %s%s", code, exit.Unusable, printed, failures)
	}
	if took := time.Since(started); took > aPollThatWaitedWouldStillBeWaiting {
		t.Errorf("agora poll --last took %s, and it never waits", took)
	}
	if !strings.Contains(failures, "socrates agora poll --session "+held.opened) {
		t.Errorf("the refusal does not say how to wait for an answer: %s", failures)
	}

	if delivered, arrived := held.awaiting(shortWindow); !arrived || delivered.Decisions[0].ID != "c-1" {
		t.Errorf("agora poll --last collected the answer it was not handed: %+v", delivered)
	}
}

func TestPollLastNamingASessionOutsideTheSessionsIsRefused(t *testing.T) {
	t.Parallel()
	held := standing(t)
	held.serve(oneDecisionDocument)

	outside := filepath.Join(held.surfaceHome(), "answer.json")
	if err := os.WriteFile(outside, []byte(`{"decisions":[],"message":"read from outside the sessions","ended":false}`), 0o600); err != nil {
		t.Fatalf("planting an answer outside the sessions: %v", err)
	}

	printed, failures, code := held.run("poll", "--session", "..", "--last")
	if code == exit.NothingAsked {
		t.Fatalf("agora poll --last read a session outside the sessions: %s", printed)
	}
	if strings.Contains(printed, "read from outside the sessions") {
		t.Errorf("agora poll --last printed a file outside the sessions:\n%s", printed)
	}
	if !strings.Contains(failures, "..") {
		t.Errorf("the refusal does not name the session: %s", failures)
	}
}
