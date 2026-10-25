package surface_test

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/socrates-tools/socrates/agora/surface"
	"github.com/socrates-tools/socrates/internal/refusal"
)

func TestASessionTheSurfaceNeverMintedIsRefused(t *testing.T) {
	const neverMinted = "0123456789abcdef"

	held := serving(t, twoDecisions(), 0)
	before := held.server.Holds()

	err := held.server.Revise(neverMinted, oneOtherDecision(), written(t, oneOtherDecision()), "")
	if err == nil {
		t.Fatalf("naming the session %q brought one into existence", neverMinted)
	}
	if !refusal.Is(err) {
		t.Errorf("naming a session the surface does not hold reads as a failure rather than a refusal: %v", err)
	}
	if !strings.Contains(err.Error(), neverMinted) {
		t.Errorf("the refusal reads %q, without naming the session", err)
	}
	if after := held.server.Holds(); !slices.Equal(before, after) {
		t.Errorf("naming a session the surface does not hold left it holding %v, and it held %v", after, before)
	}
}

func TestASecondSessionHoldsItsOwnDocumentAndItsOwnToken(t *testing.T) {
	// scenario: Two sessions are served at once and neither disturbs the other
	held := serving(t, twoDecisions(), 0)

	opened, err := held.server.Open(oneOtherDecision(), written(t, oneOtherDecision()), "")
	if err != nil {
		t.Fatalf("opening a second session on a running surface: %v", err)
	}
	if opened.Session == held.session {
		t.Fatalf("a second session was opened under the identifier the first holds, %q", opened.Session)
	}

	first, standing := held.server.Token(held.session)
	if !standing {
		t.Fatalf("the surface no longer holds the session it was started with, %q", held.session)
	}
	if opened.Token == first {
		t.Error("two sessions carry the same token, and a token answers for the session it was minted for and no other")
	}

	held.exchange(t).Revise(twoDecisions())

	second, standing := held.server.Exchange(opened.Session)
	if !standing {
		t.Fatalf("the surface holds no session %q, and it was just opened", opened.Session)
	}
	read, _ := second.Document()
	if read.Title != oneOtherDecision().Title {
		t.Errorf("the second session holds %q, and it was opened on %q", read.Title, oneOtherDecision().Title)
	}

	holds := held.server.Holds()
	if !slices.Contains(holds, held.session) || !slices.Contains(holds, opened.Session) {
		t.Errorf("the surface holds %v, and two sessions were opened on it", holds)
	}
}

func TestACollectedSessionLeavesNothingBehind(t *testing.T) {
	// scenario: An idle session is collected and the ones being read are not
	held := serving(t, twoDecisions(), 200*time.Millisecond)

	within, err := surface.SessionPath(held.session)
	if err != nil {
		t.Fatalf("resolving where the surface wrote the session: %v", err)
	}
	for _, named := range []string{"session.json", "document.json"} {
		if _, err := os.Stat(filepath.Join(within, named)); err != nil {
			t.Fatalf("a session just opened wrote no %s: %v", named, err)
		}
	}

	select {
	case <-held.server.Stopping():
	case <-time.After(10 * time.Second):
		t.Fatal("no session was read for the idle period and the server is still running")
	}

	if _, err := os.Stat(within); !os.IsNotExist(err) {
		t.Errorf("a collected session left %s behind: %v", within, err)
	}
}

func TestACollectedSessionTakesEveryAnswerCollectedInItAway(t *testing.T) {
	// scenario: The session's directory, every collected answer included, is removed when the surface collects the session
	held := serving(t, twoDecisions(), 500*time.Millisecond)

	for _, words := range []string{"the first words", "the second words"} {
		held.send(t, surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}, Words: words}}})
		if _, arrived := held.await(t, 5*time.Second); !arrived {
			t.Fatal("an answer was sent and nothing was delivered")
		}
	}

	within, err := surface.SessionPath(held.session)
	if err != nil {
		t.Fatalf("resolving where the surface wrote the session: %v", err)
	}
	for _, number := range []string{"1", "2"} {
		if _, err := os.Stat(filepath.Join(within, "answers", number, "answer.json")); err != nil {
			t.Fatalf("answer %s was collected and the session does not hold it: %v", number, err)
		}
	}

	select {
	case <-held.server.Stopping():
	case <-time.After(10 * time.Second):
		t.Fatal("no session was read for the idle period and the server is still running")
	}
	if _, err := os.Stat(within); !os.IsNotExist(err) {
		t.Errorf("a collected session left %s behind: %v", within, err)
	}
}
