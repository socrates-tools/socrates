package runstore

import (
	"bytes"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/socrates-tools/socrates/review/dialogue"
)

func TestTwoAnswersDoNotOverwriteEachOther(t *testing.T) {
	round := roundFor(t)

	const answers = 8

	var waiting sync.WaitGroup
	failures := make([]error, answers)

	for position := range failures {
		waiting.Add(1)
		go func() {
			defer waiting.Done()
			failures[position] = round.AppendResponse(dialogue.Response{
				Finding:    "f" + strconv.Itoa(position),
				Resolution: dialogue.Rejected,
				Reason:     "the claim of f" + strconv.Itoa(position),
			})
		}()
	}
	waiting.Wait()

	stored, err := round.Responses()
	if err != nil {
		t.Fatalf("reading the responses back: %v", err)
	}
	if len(stored.Recorded) != answers {
		t.Errorf("%d of %d answers are stored: %+v", len(stored.Recorded), answers, stored.Recorded)
	}

	for position, failure := range failures {
		identifier := "f" + strconv.Itoa(position)
		if failure != nil {
			t.Errorf("recording %s alongside the others failed: %v", identifier, failure)
			continue
		}
		if _, present := stored.For(identifier); !present {
			t.Errorf("%s was recorded and is not in the document: %+v", identifier, stored.Recorded)
		}
	}
}

func TestAnAnswerIsAppendedRatherThanRewritten(t *testing.T) {
	round := roundFor(t)

	if err := round.AppendResponse(dialogue.Response{Finding: "f1", Resolution: dialogue.Fixed, Reason: "the first claim"}); err != nil {
		t.Fatalf("recording the first answer: %v", err)
	}

	before, err := os.ReadFile(round.ResponsesPath())
	if err != nil {
		t.Fatalf("reading the document: %v", err)
	}

	if err := round.AppendResponse(dialogue.Response{Finding: "f1", Resolution: dialogue.Rejected, Reason: "the second claim"}); err != nil {
		t.Fatalf("recording the replacement: %v", err)
	}

	after, err := os.ReadFile(round.ResponsesPath())
	if err != nil {
		t.Fatalf("reading the document again: %v", err)
	}
	if !bytes.HasPrefix(after, before) {
		t.Errorf("the second answer rewrote what the first wrote:\nbefore: %s\nafter:  %s", before, after)
	}

	stored, err := round.Responses()
	if err != nil {
		t.Fatalf("reading the responses back: %v", err)
	}
	if len(stored.Recorded) != 2 {
		t.Errorf("%d answers are stored, want both: %+v", len(stored.Recorded), stored.Recorded)
	}
	answer, present := stored.For("f1")
	if !present || answer.Reason != "the second claim" {
		t.Errorf("the answer that counts is %+v, want the second claim", answer)
	}
	if answers := stored.Answers("f1"); answers != 2 {
		t.Errorf("answering f1 twice reads as %d answers, so the replacement is not visible", answers)
	}
}

func roundFor(t *testing.T) Round {
	t.Helper()

	store := New(t.TempDir())
	run := started(t, store)

	round, err := store.CreateRound(run.Name, 1)
	if err != nil {
		t.Fatalf("creating the round: %v", err)
	}
	return round
}
