package e2e_test

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/socrates-tools/socrates/agora/surface"
	"github.com/socrates-tools/socrates/internal/exit"
	"github.com/socrates-tools/socrates/internal/process"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

const (
	olderCommit    = "0000000"
	olderBuildFlag = "-X github.com/socrates-tools/socrates/internal/version.commit=" + olderCommit
)

func anOlderBuild(t *testing.T) string {
	t.Helper()
	return source.SharedLinkedWith(t, "./cmd/socrates", olderBuildFlag)
}

func (s *surfaced) servesFromTheOlderBuild() {
	s.t.Helper()
	s.binary = anOlderBuild(s.t)
}

func (s *surfaced) servesFromThisBuild() {
	s.t.Helper()
	s.binary = source.Shared(s.t, "./cmd/socrates")
}

func (s *surfaced) whoStandsAt(address string) surface.Identity {
	s.t.Helper()

	identity, err := surface.Ask(address).Identify(context.Background())
	if err != nil {
		s.t.Fatalf("asking %s who stands there: %v", address, err)
	}
	return identity
}

func untilTheProcessHasExited(t *testing.T, pid int) {
	t.Helper()

	deadline := time.Now().Add(listeningWait)
	for process.Alive(pid) {
		if time.Now().After(deadline) {
			t.Fatalf("the process %d of the replaced surface is still running", pid)
		}
		time.Sleep(listeningPoll)
	}
}

func TestASurfaceOfAnotherBuildHoldingNoSessionIsReplacedByTheBuildThatRanServe(t *testing.T) {
	t.Parallel()
	held := standing(t)
	held.bindsAPortNothingElseWillTake()

	held.servesFromTheOlderBuild()
	held.serve(oneDecisionDocument)
	replaced := held.record()
	held.end(held.opened)
	if told := held.whoStandsAt(held.bind); told.Commit != olderCommit {
		t.Fatalf("the surface standing was linked from the commit %q, and the test stands one from %q", told.Commit, olderCommit)
	}

	held.servesFromThisBuild()
	printed := held.serve(nothingToDecideDocument)

	if want := "http://" + held.bind + surface.SessionPrefix; !strings.HasPrefix(printed, want) {
		t.Errorf("agora serve printed %q, and the page is served from %s", printed, held.bind)
	}
	if told := held.whoStandsAt(held.bind); told.Commit == olderCommit {
		t.Errorf("the page is still drawn by the build linked from %q", told.Commit)
	}
	untilTheProcessHasExited(t, replaced.Process)
	if now := held.record(); now.Process == replaced.Process || now.Token == replaced.Token {
		t.Errorf("what the surface is reached with is still %+v, so nothing was started in place of the replaced surface", now)
	}

	page, err := http.Get(printed)
	if err != nil {
		t.Fatalf("opening the page agora serve printed: %v", err)
	}
	defer func() { _ = page.Body.Close() }()
	body, _ := io.ReadAll(page.Body)
	if page.StatusCode != http.StatusOK || !strings.Contains(string(body), "nothing to settle") {
		t.Errorf("the page agora serve printed was answered with %d and does not hold the document served: %s", page.StatusCode, body)
	}
}

func TestASurfaceOfAnotherBuildHoldingASessionIsReusedAndTheSessionStillAnswers(t *testing.T) {
	t.Parallel()
	held := standing(t)
	held.bindsAPortNothingElseWillTake()

	held.servesFromTheOlderBuild()
	held.serve(oneDecisionDocument)
	standingSurface := held.record()
	first := held.opened

	held.servesFromThisBuild()
	printed := held.serve(nothingToDecideDocument)

	if want := "http://" + held.bind + surface.SessionPrefix; !strings.HasPrefix(printed, want) {
		t.Errorf("agora serve printed %q, and the surface holding the session stands at %s", printed, held.bind)
	}
	if now := held.record(); now != standingSurface {
		t.Errorf("a surface holding a session was replaced: it was %+v and is %+v", standingSurface, now)
	}
	if told := held.whoStandsAt(held.bind); told.Commit != olderCommit {
		t.Errorf("the surface holding a session was linked from %q and is now linked from %q", olderCommit, told.Commit)
	}
	if !process.Alive(standingSurface.Process) {
		t.Error("the process of a surface holding a session has exited")
	}

	held.answerOn(first, surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}, Words: "go on"}}})
	held.opened = first
	if printed, failures, code := held.poll(); code != exit.NothingAsked || !strings.Contains(printed, "patch c-1") {
		t.Errorf("the session already there answered with exit %d and %q: %s", code, printed, failures)
	}
}

func TestASurfaceOfTheSameBuildIsReusedWhetherOrNotItHoldsASession(t *testing.T) {
	t.Parallel()
	held := standing(t)
	held.bindsAPortNothingElseWillTake()

	held.serve(oneDecisionDocument)
	standingSurface := held.record()
	held.end(held.opened)

	held.serve(nothingToDecideDocument)
	if now := held.record(); now != standingSurface {
		t.Errorf("a surface of this build was replaced: it was %+v and is %+v", standingSurface, now)
	}
}

func aSurfaceOfAnotherBuildThatNeverLetsGo(t *testing.T, held *surfaced) {
	t.Helper()

	identity := surface.Identity{Tool: surface.IdentityTool, Root: held.home, Version: "v0.0.1", Commit: olderCommit}
	standing := httptest.NewUnstartedServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case surface.IdentityRoute:
			_ = json.NewEncoder(writer).Encode(identity)
		case surface.RetireRoute:
			writer.WriteHeader(http.StatusNoContent)
		default:
			http.NotFound(writer, request)
		}
	}))
	listener, err := net.Listen("tcp", aPortTheOperatingSystemChooses)
	if err != nil {
		t.Fatalf("standing in for a surface on a port the operating system chooses: %v", err)
	}
	held.bind = listener.Addr().String()
	standing.Listener = listener
	standing.Start()
	t.Cleanup(standing.Close)

	written, err := json.Marshal(surface.Record{Address: held.bind, Token: "the-control-token", Process: os.Getpid()})
	if err != nil {
		t.Fatalf("writing what the surface wrote about itself: %v", err)
	}
	if err := os.MkdirAll(held.surfaceHome(), 0o700); err != nil {
		t.Fatalf("making the surface's directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(held.surfaceHome(), "surface.json"), written, 0o600); err != nil {
		t.Fatalf("writing what the surface wrote about itself: %v", err)
	}
}

func TestASurfaceOfAnotherBuildThatDoesNotLetGoOfThePortIsRefusedAsAPortHeldBySomethingElse(t *testing.T) {
	t.Parallel()
	held := standing(t)
	aSurfaceOfAnotherBuildThatNeverLetsGo(t, held)

	started := time.Now()
	printed, failures, code := held.reading(oneDecisionDocument, held.serving()...)
	if code == exit.NothingAsked {
		t.Fatalf("a port that was never let go was served at %q", printed)
	}
	for _, named := range []string{held.bind, surface.BindField, surface.AddressOption} {
		if !strings.Contains(failures, named) {
			t.Errorf("refusing a port that is never let go does not name %q: %s", named, failures)
		}
	}
	if waited := time.Since(started); waited > 30*time.Second {
		t.Errorf("serve waited %s for a port that was never let go, and the wait is bounded", waited)
	}
}

func TestAServeNamingASessionNeverReplacesTheSurfaceHoldingIt(t *testing.T) {
	t.Parallel()
	held := standing(t)
	held.bindsAPortNothingElseWillTake()

	held.servesFromTheOlderBuild()
	held.serve(oneDecisionDocument)
	standingSurface := held.record()
	held.end(held.opened)

	held.servesFromThisBuild()
	held.serveAgain(nothingToDecideDocument)

	if now := held.record(); now != standingSurface {
		t.Errorf("a serve naming a session on a surface of another build replaced it: it was %+v and is %+v", standingSurface, now)
	}
}
