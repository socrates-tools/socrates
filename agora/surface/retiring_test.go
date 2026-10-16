package surface_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/socrates-tools/socrates/agora/surface"
	"github.com/socrates-tools/socrates/internal/state"
)

func (s *standing) retiring(t *testing.T, token string) int {
	t.Helper()
	return s.posting(t, s.at(surface.RetireRoute), token, nil)
}

func (s *standing) stoppedWithin(wait time.Duration) bool {
	select {
	case <-s.server.Stopping():
		return true
	case <-time.After(wait):
		return false
	}
}

func TestARequestWithNoTokenLearnsNothingOfWhetherTheSurfaceHoldsASession(t *testing.T) {
	// scenario: No request without a token learns whether the surface holds a session
	holding := serving(t, twoDecisions(), 0)
	idle := serving(t, twoDecisions(), 0)
	idle.exchange(t).End()

	withoutToken := map[string]int{}
	for named, held := range map[string]*standing{"holding": holding, "idle": idle} {
		withoutToken[named] = held.retiring(t, "")
		if held.stoppedWithin(200 * time.Millisecond) {
			t.Errorf("a request with no token stopped the surface that was %s", named)
		}
	}
	if withoutToken["holding"] != http.StatusForbidden || withoutToken["idle"] != http.StatusForbidden {
		t.Errorf("a request with no token was answered with %v, and both are refused alike", withoutToken)
	}

	for named, held := range map[string]*standing{"holding": holding, "idle": idle} {
		_, said := held.fetch(t, held.at(surface.IdentityRoute))
		var told map[string]any
		if err := json.Unmarshal([]byte(said), &told); err != nil {
			t.Fatalf("reading what the %s surface said about itself: %v", named, err)
		}
		if len(told) != 5 {
			t.Errorf("the %s surface tells a request with no token %v, and that is more than who it is", named, told)
		}
	}
}

func TestASurfaceHoldingASessionStillBeingAnsweredIsNotRetired(t *testing.T) {
	// scenario: A surface holding a session is reused and the session still answers
	held := serving(t, twoDecisions(), 0)

	if code := held.retiring(t, held.server.Control()); code != http.StatusConflict {
		t.Fatalf("retiring a surface holding a session was answered with %d, want %d", code, http.StatusConflict)
	}
	if held.stoppedWithin(200 * time.Millisecond) {
		t.Fatal("a surface holding a session was stopped")
	}
	if code := held.send(t, surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}}}}); code != http.StatusNoContent {
		t.Errorf("the session of a surface that refused to retire answered with %d", code)
	}
}

func TestASurfaceWhoseEverySessionWasEndedIsRetiredAndStops(t *testing.T) {
	// scenario: A surface holding no session is replaced
	held := serving(t, twoDecisions(), 0)
	held.exchange(t).End()

	if code := held.retiring(t, held.server.Control()); code != http.StatusNoContent {
		t.Fatalf("retiring a surface whose only session was ended was answered with %d, want %d", code, http.StatusNoContent)
	}
	if !held.stoppedWithin(5 * time.Second) {
		t.Error("a retired surface did not stop")
	}
}

func TestOneSessionLeftOpenBesideAnEndedOneKeepsTheSurface(t *testing.T) {
	// scenario: A surface holding a session is reused and the session still answers
	held := serving(t, twoDecisions(), 0)
	held.exchange(t).End()
	if _, err := held.server.Open(oneOtherDecision(), written(t, oneOtherDecision()), ""); err != nil {
		t.Fatalf("opening a second session: %v", err)
	}

	if code := held.retiring(t, held.server.Control()); code != http.StatusConflict {
		t.Errorf("retiring a surface with one session still open was answered with %d, want %d", code, http.StatusConflict)
	}
}

func TestARetiredSurfaceTakesNoSessionAndReopensNone(t *testing.T) {
	// scenario: A retired surface is not handed a session the replacement would lose
	held := serving(t, twoDecisions(), 0)
	held.exchange(t).End()
	if code := held.retiring(t, held.server.Control()); code != http.StatusNoContent {
		t.Fatalf("retiring was answered with %d", code)
	}

	if _, err := held.server.Open(oneOtherDecision(), written(t, oneOtherDecision()), ""); err == nil {
		t.Error("a retired surface took a new session")
	}
	if err := held.server.Revise(held.session, oneOtherDecision(), written(t, oneOtherDecision()), ""); err == nil {
		t.Error("a retired surface reopened a session that was ended")
	}
	if !held.exchange(t).Over() {
		t.Error("a retired surface's ended session was reopened")
	}
}

func TestACallerHoldingTheTokenRetiresAnIdleSurface(t *testing.T) {
	held := serving(t, twoDecisions(), 0)
	caller := surface.Reach(surface.Record{Address: held.server.Address(), Token: held.server.Control()})

	retired, err := caller.Retire(context.Background())
	if err != nil || retired {
		t.Fatalf("retiring a surface holding a session answered (%v, %v), want not retired and no error", retired, err)
	}

	held.exchange(t).End()
	retired, err = caller.Retire(context.Background())
	if err != nil || !retired {
		t.Fatalf("retiring a surface holding no session answered (%v, %v), want retired", retired, err)
	}
}

func TestBuildsDifferInAnythingTheBinarySaysAboutItself(t *testing.T) {
	// scenario: Another build is any difference in the version, the commit stamp or whether it was built uncommitted
	built := surface.Identity{Tool: "agora", Version: "v0.2.0", Commit: "abc1234", Uncommitted: false}

	for named, other := range map[string]surface.Identity{
		"the version":       {Version: "v0.3.0", Commit: "abc1234"},
		"the commit":        {Version: "v0.2.0", Commit: "def5678"},
		"being uncommitted": {Version: "v0.2.0", Commit: "abc1234", Uncommitted: true},
	} {
		if built.SameBuild(other) {
			t.Errorf("builds differing in %s were taken for one", named)
		}
	}
	if !built.SameBuild(surface.Identity{Tool: "agora", Root: "elsewhere", Version: "v0.2.0", Commit: "abc1234"}) {
		t.Error("two surfaces of one build were taken for two, and the root is not part of the build")
	}
}

func TestLocatingReportsASurfaceOfAnotherBuildAndASurfaceOfThisOne(t *testing.T) {
	// scenario: A surface of another build is told apart from one of this build
	held := serving(t, twoDecisions(), 0)

	same := locating(t, held.server.Address())
	if !same.Standing || same.OfAnotherBuild {
		t.Errorf("a surface of this build was found as %+v", same)
	}

	root, err := state.Root()
	if err != nil {
		t.Fatalf("resolving the root: %v", err)
	}
	older := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(writer).Encode(surface.Identity{Tool: surface.IdentityTool, Root: root, Version: "v0.0.1", Commit: "0000000"})
	}))
	t.Cleanup(older.Close)

	found := locating(t, strings.TrimPrefix(older.URL, "http://"))
	if !found.Standing || !found.OfAnotherBuild {
		t.Errorf("a surface of another build was found as %+v", found)
	}
}

func TestAnOldSurfaceLeavesTheRecordOfItsReplacementAlone(t *testing.T) {
	// scenario: A replaced surface does not forget what its replacement wrote
	t.Setenv(state.HomeVariable, t.TempDir())
	replacement := surface.Record{Address: "127.0.0.1:24672", Token: "replacement", Process: 2}
	writeRecord(t, replacement)

	if err := surface.ForgetRecordOf("old"); err != nil {
		t.Fatalf("forgetting what the old surface wrote: %v", err)
	}
	if kept, err := surface.ReadRecord(); err != nil || kept != replacement {
		t.Errorf("the record of the replacement is now %+v (%v), and it was %+v", kept, err, replacement)
	}

	if err := surface.ForgetRecordOf("replacement"); err != nil {
		t.Fatalf("forgetting what the surface wrote: %v", err)
	}
	if _, err := surface.ReadRecord(); err == nil {
		t.Error("a surface did not forget the record it wrote")
	}
}
