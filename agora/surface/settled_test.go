package surface_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/socrates-tools/socrates/agora/document"
	"github.com/socrates-tools/socrates/agora/surface"
	"github.com/socrates-tools/socrates/internal/refusal"
)

const inThePersonsOwnWords = "the lock guards the read, so take it before the cache is touched"

func answeredAndCollected(t *testing.T, held *standing, sent surface.Sent) {
	t.Helper()

	if _, err := held.exchange(t).Receive(sent); err != nil {
		t.Fatalf("answering the page: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if _, arrived, err := held.exchange(t).Await(ctx, time.Second, keptAnywhere); err != nil || !arrived {
		t.Fatal("the answer the page sent was never collected")
	}
}

func theCacheQuestionAnswered() surface.Sent {
	return surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}, Words: inThePersonsOwnWords}}}
}

func settling(entry document.Settled) document.Document {
	return document.Document{
		Title:     "one thing left to settle",
		Settled:   []document.Settled{entry},
		Decisions: oneOtherDecision().Decisions,
	}
}

func asCollected() document.Settled {
	return document.Settled{
		ID:       "c-1",
		Question: "should the cache be read before the lock is taken",
		Chosen:   []string{"take the lock first"},
		Words:    inThePersonsOwnWords,
	}
}

func TestASettledAnswerIsStatedAsTheSessionCollectedIt(t *testing.T) {
	held := serving(t, twoDecisions(), 0)
	answeredAndCollected(t, held, theCacheQuestionAnswered())

	revised := settling(asCollected())
	if err := held.server.Revise(held.session, revised, written(t, revised)); err != nil {
		t.Fatalf("a settled answer stated exactly as it was collected was refused: %v", err)
	}
	if served, _ := held.exchange(t).Document(); len(served.Settled) != 1 || served.Settled[0].Words != inThePersonsOwnWords {
		t.Errorf("the page does not carry the settled answer it was served: %+v", served.Settled)
	}
}

func TestASettledAnswerRestatedOtherwiseThanCollectedIsRefused(t *testing.T) {
	for _, restated := range []struct {
		named  string
		change func(*document.Settled)
		names  string
	}{
		{named: "the person's words paraphrased", change: func(entry *document.Settled) {
			entry.Words = "take the lock first because it guards the read"
		}, names: "words"},
		{named: "the question reworded", change: func(entry *document.Settled) {
			entry.Question = "should the lock be taken at all"
		}, names: "question"},
		{named: "another option said to be chosen", change: func(entry *document.Settled) {
			entry.Chosen = []string{"leave it"}
		}, names: "chosen"},
		{named: "the chosen option left out", change: func(entry *document.Settled) {
			entry.Chosen = nil
		}, names: "chosen"},
	} {
		t.Run(restated.named, func(t *testing.T) {
			held := serving(t, twoDecisions(), 0)
			answeredAndCollected(t, held, theCacheQuestionAnswered())
			before, revision := held.exchange(t).Document()

			entry := asCollected()
			restated.change(&entry)
			revised := settling(entry)

			err := held.server.Revise(held.session, revised, written(t, revised))
			if err == nil {
				t.Fatal("a settled answer that is not what the session collected was served")
			}
			if !refusal.Is(err) {
				t.Errorf("a settled answer that is not what was collected reads as a failure rather than a refusal: %v", err)
			}
			for _, named := range []string{"c-1", restated.names, "agora poll"} {
				if !strings.Contains(err.Error(), named) {
					t.Errorf("the refusal reads %q, without naming %q", err, named)
				}
			}
			if after, now := held.exchange(t).Document(); now != revision || after.Title != before.Title {
				t.Errorf("a refused revision moved the page from revision %d to %d", revision, now)
			}
		})
	}
}

func TestOnlyWhatTheSessionCollectedBindsWhatIsSettled(t *testing.T) {
	t.Run("An identifier the session never collected an answer for", func(t *testing.T) {
		held := serving(t, twoDecisions(), 0)

		revised := settling(document.Settled{ID: "c-1", Question: "should the cache be read before the lock is taken", Words: "decided elsewhere"})
		if err := held.server.Revise(held.session, revised, written(t, revised)); err != nil {
			t.Errorf("a settled answer the surface holds nothing to compare with was refused: %v", err)
		}
	})

	t.Run("An answer sent but never collected", func(t *testing.T) {
		held := serving(t, twoDecisions(), 0)
		if _, err := held.exchange(t).Receive(theCacheQuestionAnswered()); err != nil {
			t.Fatalf("answering the page: %v", err)
		}

		entry := asCollected()
		entry.Words = "what the caller believes was answered"
		revised := settling(entry)
		if err := held.server.Revise(held.session, revised, written(t, revised)); err != nil {
			t.Errorf("an answer the caller never collected decided what it may settle: %v", err)
		}
	})

	t.Run("A question asked again binds what was answered last", func(t *testing.T) {
		held := serving(t, twoDecisions(), 0)
		answeredAndCollected(t, held, theCacheQuestionAnswered())

		held.exchange(t).Revise(twoDecisions())
		const answeredAgain = "on reflection, leave it"
		answeredAndCollected(t, held, surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{1}, Words: answeredAgain}}})

		stale := settling(asCollected())
		if err := held.server.Revise(held.session, stale, written(t, stale)); err == nil {
			t.Error("a settled answer stated as the first answer was served after the question was answered again")
		}

		latest := asCollected()
		latest.Chosen = []string{"leave it"}
		latest.Words = answeredAgain
		current := settling(latest)
		if err := held.server.Revise(held.session, current, written(t, current)); err != nil {
			t.Errorf("a settled answer stated as the last answer was refused: %v", err)
		}
	})
}
