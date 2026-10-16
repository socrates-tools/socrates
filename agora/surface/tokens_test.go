package surface_test

import (
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/socrates-tools/socrates/agora/surface"
)

func TestOneSessionsTokenDoesNotAnswerForAnother(t *testing.T) {
	held := serving(t, twoDecisions(), 0)

	other, err := held.server.Open(oneOtherDecision(), written(t, oneOtherDecision()), "")
	if err != nil {
		t.Fatalf("opening a second session on a running surface: %v", err)
	}

	if code := held.sendingAs(t, held.session, other.Token, surface.Sent{Message: "let me in"}); code != http.StatusForbidden {
		t.Errorf("an answer to one session carrying another session's token was answered with %d", code)
	}

	read, arrived := held.await(t, 200*time.Millisecond)
	if arrived {
		t.Errorf("an answer carrying another session's token was delivered: %+v", read)
	}
}

func TestTheControlTokenIsAcceptedOnEveryRoute(t *testing.T) {
	held := serving(t, twoDecisions(), 0)

	if code := held.sendingAs(t, held.session, held.server.Control(), surface.Sent{Message: "from whoever started it"}); code != http.StatusNoContent {
		t.Errorf("an answer carrying the control token was answered with %d", code)
	}
	if _, arrived := held.await(t, 5*time.Second); !arrived {
		t.Error("an answer carrying the control token was not delivered")
	}
}

func TestOpeningASessionTakesTheControlTokenAlone(t *testing.T) {
	held := serving(t, twoDecisions(), 0)

	if code := held.posting(t, held.at(surface.SessionPrefix), held.token, written(t, oneOtherDecision())); code != http.StatusForbidden {
		t.Fatalf("opening a session carrying a page's own token was answered with %d", code)
	}
	if holds := held.server.Holds(); len(holds) != 1 {
		t.Fatalf("a refused opening left the surface holding %v", holds)
	}

	if code := held.posting(t, held.at(surface.SessionPrefix), held.server.Control(), written(t, oneOtherDecision())); code != http.StatusOK {
		t.Fatalf("opening a session carrying the control token was answered with %d", code)
	}
	if holds := held.server.Holds(); len(holds) != 2 {
		t.Errorf("the surface holds %v after a session was opened on it", holds)
	}
}

func TestStoppingTheServerTakesTheControlTokenAlone(t *testing.T) {
	held := serving(t, twoDecisions(), 0)

	if code := held.posting(t, held.at("/stop"), held.token, nil); code != http.StatusForbidden {
		t.Fatalf("stopping the server carrying a page's own token was answered with %d", code)
	}
	select {
	case <-held.server.Stopping():
		t.Fatal("a refused stop took the server down anyway")
	case <-time.After(200 * time.Millisecond):
	}

	if code := held.posting(t, held.at("/stop"), held.server.Control(), nil); code != http.StatusNoContent {
		t.Fatalf("stopping the server carrying the control token was answered with %d", code)
	}
	select {
	case <-held.server.Stopping():
	case <-time.After(5 * time.Second):
		t.Error("the server was stopped by whoever started it and is still running")
	}
}

func TestTwoCallersStoppingAtOnceTakeTheServerDownOnce(t *testing.T) {
	// scenario: An idle server stops
	held := serving(t, twoDecisions(), 0)

	together := make(chan struct{})
	var stopping sync.WaitGroup
	for range 64 {
		stopping.Add(1)
		go func() {
			defer stopping.Done()
			<-together
			held.server.Stop()
		}()
	}
	close(together)
	stopping.Wait()

	select {
	case <-held.server.Stopping():
	case <-time.After(5 * time.Second):
		t.Fatal("the server was stopped and is still running")
	}
}
