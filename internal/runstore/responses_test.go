package runstore

import (
	"bytes"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/rui-ktei/socrates/internal/dialogue"
)

func TestTwoAnswersDoNotOverwriteEachOther(t *testing.T) {
	round := roundFor(t)

	var waiting sync.WaitGroup
	recorded := make([]bool, 8)

	for position := range recorded {
		waiting.Add(1)
		go func() {
			defer waiting.Done()
			_, _, err := round.RecordResponse(dialogue.Response{
				Finding:    "f" + strconv.Itoa(position),
				Resolution: dialogue.Rejected,
				Reason:     "the claim of f" + strconv.Itoa(position),
			})
			recorded[position] = err == nil
		}()
	}
	waiting.Wait()

	stored, err := round.Responses()
	if err != nil {
		t.Fatalf("reading the responses back: %v", err)
	}

	for position, succeeded := range recorded {
		if !succeeded {
			continue
		}
		identifier := "f" + strconv.Itoa(position)
		if _, present := stored.For(identifier); !present {
			t.Errorf("%s was recorded and is not in the document: %+v", identifier, stored.Recorded)
		}
	}
}

func TestAnAnswerIsAppendedRatherThanRewritten(t *testing.T) {
	round := roundFor(t)

	if _, _, err := round.RecordResponse(dialogue.Response{Finding: "f1", Resolution: dialogue.Fixed, Reason: "the first claim"}); err != nil {
		t.Fatalf("recording the first answer: %v", err)
	}

	before, err := os.ReadFile(round.ResponsesPath())
	if err != nil {
		t.Fatalf("reading the document: %v", err)
	}

	if _, _, err := round.RecordResponse(dialogue.Response{Finding: "f1", Resolution: dialogue.Rejected, Reason: "the second claim"}); err != nil {
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
