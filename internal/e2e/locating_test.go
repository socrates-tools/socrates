package e2e_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/socrates-tools/socrates/agora/surface"
	"github.com/socrates-tools/socrates/internal/exit"
)

func (s *surfaced) recordNamesAnAddressNothingAnswersAt() {
	s.t.Helper()

	open := s.record()
	open.Address = aFreeLoopbackAddress(s.t)
	written, err := json.Marshal(open)
	if err != nil {
		s.t.Fatalf("writing what the surface wrote about itself: %v", err)
	}
	if err := os.WriteFile(filepath.Join(s.surfaceHome(), "surface.json"), written, 0o600); err != nil {
		s.t.Fatalf("rewriting what the surface wrote about itself: %v", err)
	}
}

func TestAServeFindsTheSurfaceOnThePortWhateverTheRecordNames(t *testing.T) {
	// scenario: A record naming an address where nothing answers does not stop serve
	t.Parallel()
	held := standing(t)
	held.bindsWhatIsConfigured()

	held.serve(oneDecisionDocument)
	first := held.opened
	held.recordNamesAnAddressNothingAnswersAt()

	printed, failures, code := held.reading(nothingToDecideDocument, held.serving()...)
	if code != exit.NothingAsked {
		t.Fatalf("agora serve exited %d beside a record naming an address where nothing answers: %s", code, failures)
	}
	if !strings.HasPrefix(printed, "http://"+theSuitesDefault(t)+surface.SessionPrefix) {
		t.Errorf("agora serve printed %q, and the surface stands on %s", printed, theSuitesDefault(t))
	}

	held.opened = theSessionIn(t, strings.TrimSpace(printed))
	if held.opened == first {
		t.Error("serving a second document onto the surface standing reused the first session")
	}
}

func TestAStopStopsTheSurfaceOnThePortWhateverTheRecordNames(t *testing.T) {
	// scenario: socrates agora stop stops the surface standing on the port
	t.Parallel()
	held := standing(t)
	held.bindsWhatIsConfigured()

	held.serve(oneDecisionDocument)
	held.recordNamesAnAddressNothingAnswersAt()

	if _, failures, code := held.run("stop"); code != exit.NothingAsked {
		t.Fatalf("agora stop exited %d beside a record naming an address where nothing answers: %s", code, failures)
	}
	untilThePortIsFree(t, theSuitesDefault(t))
}

func TestAServeUnderOneRootRefusesAPortAnotherRootsSurfaceHolds(t *testing.T) {
	// scenario: A port held by an agora surface of another root is refused and named
	t.Parallel()
	owner := standing(t)
	owner.bindsWhatIsConfigured()
	owner.serve(oneDecisionDocument)

	another := standing(t)
	another.bindsWhatIsConfigured()

	printed, failures, code := another.reading(oneDecisionDocument, another.serving()...)
	if code == exit.NothingAsked {
		t.Fatalf("a serve under another root was served at %q beside a surface holding %s", printed, theSuitesDefault(t))
	}
	for _, named := range []string{owner.home, theSuitesDefault(t), surface.BindField, surface.AddressOption} {
		if !strings.Contains(failures, named) {
			t.Errorf("refusing a port another root's surface holds does not name %q: %s", named, failures)
		}
	}
	another.startedNoSurface("a port held by another root")

	if _, arrived := owner.awaitingOn(owner.opened, 50*time.Millisecond); arrived {
		t.Error("the other root's session was answered by a refused serve")
	}
}

func TestARefusedServeLeavesTheLiveSurfaceOfThisRootReachable(t *testing.T) {
	// scenario: A refusal for another root's port does not forget this root's live surface
	t.Parallel()
	owner := standing(t)
	owner.bindsWhatIsConfigured()
	owner.serve(oneDecisionDocument)

	mine := standing(t)
	mine.serve(oneDecisionDocument)
	recorded := mine.record()
	mine.bindsWhatIsConfigured()

	printed, failures, code := mine.reading(oneDecisionDocument, mine.serving()...)
	if code == exit.NothingAsked {
		t.Fatalf("a serve was answered at %q beside another root holding %s", printed, theSuitesDefault(t))
	}
	if !strings.Contains(failures, owner.home) {
		t.Errorf("the refusal does not name the root holding the port: %s", failures)
	}

	if after := mine.record(); after != recorded {
		t.Errorf("a refused serve changed what this root's live surface is reached with: %+v became %+v", recorded, after)
	}
	mine.answer(surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}, Words: "go on"}}})
	if printed, failures, code := mine.poll(); code != exit.NothingAsked || !strings.Contains(printed, "patch c-1") {
		t.Errorf("agora poll exited %d and printed %q, so it did not reach the live surface: %s", code, printed, failures)
	}
	mine.end(mine.opened)
	if _, failures, code := mine.run("stop"); code != exit.NothingAsked {
		t.Errorf("agora stop exited %d, so it did not reach the live surface: %s", code, failures)
	}
}
