package surface_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/agora/surface"
)

func TestEachSessionIsServedUnderAPathHoldingIt(t *testing.T) {
	// scenario: Two sessions are served at once and neither disturbs the other
	held := serving(t, twoDecisions(), 0)

	other, err := held.server.Open(oneOtherDecision(), written(t, oneOtherDecision()), "")
	if err != nil {
		t.Fatalf("opening a second session on a running surface: %v", err)
	}

	for _, each := range []struct {
		session, token, title, otherTitle string
	}{
		{held.session, held.token, twoDecisions().Title, oneOtherDecision().Title},
		{other.Session, other.Token, oneOtherDecision().Title, twoDecisions().Title},
	} {
		code, page := held.fetch(t, held.at(surface.SessionPrefix+each.session))
		if code != http.StatusOK {
			t.Fatalf("the page of session %s was answered with %d", each.session, code)
		}
		if !strings.Contains(page, surface.SessionPrefix+each.session) {
			t.Errorf("the page of session %s does not hold its own address", each.session)
		}

		served, said := held.servedTo(t, each.session, each.token)
		if served != http.StatusOK {
			t.Fatalf("the document of session %s was answered with %d", each.session, served)
		}
		if !strings.Contains(string(said), each.title) {
			t.Errorf("the document of session %s does not hold %q", each.session, each.title)
		}
		if strings.Contains(string(said), each.otherTitle) {
			t.Errorf("the document of session %s holds the other session's document", each.session)
		}
	}
}

func TestAPageNamingASessionTheSurfaceDoesNotHoldIsNotFound(t *testing.T) {
	// scenario: A session the surface never minted is refused
	const neverMinted = "fedcba9876543210"

	held := serving(t, twoDecisions(), 0)

	code, said := held.fetch(t, held.at(surface.SessionPrefix+neverMinted))
	if code != http.StatusNotFound {
		t.Fatalf("a page under a session the surface does not hold was answered with %d", code)
	}
	if !strings.Contains(said, neverMinted) {
		t.Errorf("what the surface said reads %q, without naming the session", said)
	}
}

func TestTheAssetsAreServedOutsideEverySessionsPath(t *testing.T) {
	// scenario: Two sessions are served at once and neither disturbs the other
	held := serving(t, twoDecisions(), 0)

	code, _ := held.fetch(t, held.at("/assets/surface.css"))
	if code != http.StatusOK {
		t.Errorf("the shared stylesheet was answered with %d", code)
	}

	under, _ := held.fetch(t, held.at(surface.SessionPrefix+held.session+"/assets/surface.css"))
	if under == http.StatusOK {
		t.Error("the assets are served under a session's own path as well, so a second page pays for them twice")
	}
}

func aDocumentLargerThan(ceiling int) []byte {
	padding := strings.Repeat("x", ceiling)
	return []byte(`{"title":"` + padding + `","decisions":[{"id":"c-1","question":"q","mode":"single-choice","options":[{"label":"yes","command":"agree c-1"}]}]}`)
}

func TestADocumentOverTheCeilingIsRefusedRatherThanTruncated(t *testing.T) {
	// scenario: A session is minted and the address holds it
	held := serving(t, twoDecisions(), 0)
	oversized := aDocumentLargerThan(surface.BodyCeiling)

	opening := held.posting(t, held.at(surface.SessionPrefix), held.server.Control(), oversized)
	if opening != http.StatusRequestEntityTooLarge {
		t.Errorf("opening a session on a document over the ceiling was answered with %d, want %d", opening, http.StatusRequestEntityTooLarge)
	}
	if holds := held.server.Holds(); len(holds) != 1 {
		t.Errorf("a document over the ceiling left the surface holding %v", holds)
	}

	revising := held.posting(t, held.at(surface.SessionPrefix+held.session+"/document"), held.token, oversized)
	if revising != http.StatusRequestEntityTooLarge {
		t.Errorf("writing a document over the ceiling was answered with %d, want %d", revising, http.StatusRequestEntityTooLarge)
	}
	read, _ := held.exchange(t).Document()
	if read.Title != twoDecisions().Title {
		t.Errorf("a document over the ceiling reached the page as %q", read.Title)
	}
}

func TestAPageIsReachableWithAndWithoutATrailingSlash(t *testing.T) {
	held := serving(t, twoDecisions(), 0)

	for _, typed := range []string{held.under(""), held.under("/")} {
		code, page := held.fetch(t, typed)
		if code != http.StatusOK {
			t.Errorf("the page at %s was answered with %d", typed, code)
			continue
		}
		if !strings.Contains(page, surface.SessionPrefix+held.session) {
			t.Errorf("what %s answered does not hold the address of its session %s", typed, held.session)
		}
	}
}

func TestOpeningASessionAnswersOnlyAtTheSessionsPathItself(t *testing.T) {
	held := serving(t, twoDecisions(), 0)
	opening := written(t, oneOtherDecision())

	if code := held.posting(t, held.under(""), held.server.Control(), opening); code != http.StatusMethodNotAllowed {
		t.Errorf("posting a document to the page of session %s was answered with %d, want %d, so something other than the page answers there", held.session, code, http.StatusMethodNotAllowed)
	}
	if holds := held.server.Holds(); len(holds) != 1 {
		t.Fatalf("posting a document to a session's page left the surface holding %v", holds)
	}

	if code := held.posting(t, held.at(surface.SessionPrefix), held.server.Control(), opening); code != http.StatusOK {
		t.Errorf("posting a document to %s was answered with %d, so no session was opened where one is", surface.SessionPrefix, code)
	}
	if holds := held.server.Holds(); len(holds) != 2 {
		t.Errorf("opening a session at %s left the surface holding %v", surface.SessionPrefix, holds)
	}
}
