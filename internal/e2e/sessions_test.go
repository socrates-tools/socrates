package e2e_test

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates-review/internal/surface"
)

func (s *surfaced) twoPages() (string, string, string, string) {
	s.t.Helper()

	firstAddress := s.serve(twoDecisionDocument)
	first := s.opened
	secondAddress := s.serve(oneDecisionDocument)
	second := s.opened

	if first == second {
		s.t.Fatalf("two documents served without naming a session landed on the same page, %q", first)
	}
	return firstAddress, first, secondAddress, second
}

func fetched(t *testing.T, address string) (int, string) {
	t.Helper()

	answered, err := http.Get(address)
	if err != nil {
		t.Fatalf("reaching %s: %v", address, err)
	}
	defer func() { _ = answered.Body.Close() }()

	said, err := io.ReadAll(answered.Body)
	if err != nil {
		t.Fatalf("reading what %s said: %v", address, err)
	}
	return answered.StatusCode, string(said)
}

func TestASessionIsMintedAndTheAddressHoldsIt(t *testing.T) {
	held := standing(t)

	address := held.serve(oneDecisionDocument)
	bound := held.record().Address

	if want := "http://" + bound + surface.SessionPrefix; !strings.HasPrefix(address, want) {
		t.Fatalf("agora serve printed %q, and the address holding a page begins %q", address, want)
	}
	if held.opened == "" {
		t.Fatal("agora serve printed an address holding no session")
	}

	code, page := fetched(t, address)
	if code != http.StatusOK {
		t.Fatalf("the page at the address agora serve printed answered with %d", code)
	}
	if !strings.Contains(page, "should the cache be read before the lock is taken") {
		t.Errorf("the page at the address agora serve printed is not the document that was served:\n%s", page)
	}
}

func TestASessionIsNamedToPollIt(t *testing.T) {
	const onTheSecondPage = "this one is the second page"

	held := standing(t)
	_, first, _, second := held.twoPages()

	held.answerOn(second, surface.Sent{
		Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}, Words: onTheSecondPage}},
	})

	printed, failures, code := held.run("poll", "--session", second)
	if code != 0 {
		t.Fatalf("agora poll --session exited %d: %s", code, failures)
	}
	if !strings.Contains(printed, onTheSecondPage) {
		t.Errorf("agora poll on the session that was answered printed:\n%s", printed)
	}

	if delivered, arrived := held.awaitingOn(first, shortWindow); arrived {
		t.Errorf("polling the other session reported what was answered on the second: %+v", delivered)
	}
}

func TestAPollNamingNoSessionIsRefused(t *testing.T) {
	// scenario: A session is named to poll it
	held := standing(t)
	held.serve(oneDecisionDocument)

	for _, command := range []string{"poll", "end"} {
		t.Run(command, func(t *testing.T) {
			printed, failures, code := held.run(command)
			if code == 0 {
				t.Fatalf("agora %s naming no session succeeded: %s", command, printed)
			}
			if !strings.Contains(failures, "--session") {
				t.Errorf("refusing agora %s naming no session does not name the flag: %s", command, failures)
			}
		})
	}
}

func TestANewDocumentReachesTheSessionItNamedAndNoOther(t *testing.T) {
	held := standing(t)
	firstAddress, _, secondAddress, second := held.twoPages()

	held.serveOnto(second, oneDecisionReply)

	if _, page := fetched(t, secondAddress); !strings.Contains(page, "one thing left") {
		t.Errorf("the page a new document was written to did not take it up:\n%s", page)
	}
	if _, page := fetched(t, firstAddress); strings.Contains(page, "one thing left") {
		t.Errorf("a document written to one session reached the other:\n%s", page)
	}
}

func TestASessionTheSurfaceNeverMintedIsRefusedByTheBinary(t *testing.T) {
	// scenario: A session the surface never minted is refused
	const neverMinted = "fedcba9876543210"

	held := standing(t)
	held.serve(oneDecisionDocument)
	before := held.opened

	printed, failures, code := held.reading(oneDecisionDocument, held.serving("--session", neverMinted)...)
	if code == 0 {
		t.Fatalf("naming a session the surface does not hold brought one into existence at %q", printed)
	}
	if !strings.Contains(failures, neverMinted) {
		t.Errorf("the refusal reads %q, without naming the session", failures)
	}

	if _, err := os.Stat(filepath.Join(held.surfaceHome(), "sessions", neverMinted)); err == nil {
		t.Errorf("a refused session was written under %s", filepath.Join(held.surfaceHome(), "sessions"))
	}
	if held.opened != before {
		t.Errorf("a refused session replaced the one the surface holds")
	}
}

func TestEndingOneSessionLeavesTheOthersBeingRead(t *testing.T) {
	held := standing(t)
	_, first, secondAddress, second := held.twoPages()

	held.end(first)

	code, _ := fetched(t, secondAddress)
	if code != http.StatusOK {
		t.Fatalf("ending one session answered the other's page with %d", code)
	}

	held.answerOn(second, surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}}}})
	delivered, arrived := held.awaitingOn(second, 5*time.Second)
	if !arrived {
		t.Fatal("the session that was not ended no longer answers")
	}
	if diff := cmp.Diff([]string{"patch c-1"}, delivered.Decisions[0].Commands); diff != "" {
		t.Errorf("what the session that was not ended delivered differs (-want +got):\n%s", diff)
	}
}

func TestStoppingIsTheServerNotASession(t *testing.T) {
	held := standing(t)
	firstAddress, _, secondAddress, _ := held.twoPages()
	bound := held.record().Address

	if _, failures, code := held.run("stop"); code != 0 {
		t.Fatalf("agora stop exited %d: %s", code, failures)
	}
	untilThePortIsFree(t, bound)

	for _, address := range []string{firstAddress, secondAddress} {
		answered, err := http.Get(address)
		if err == nil {
			_ = answered.Body.Close()
			t.Errorf("%s still answers after the server was stopped", address)
		}
	}
}

func TestStopNamesNoSession(t *testing.T) {
	// scenario: Stopping is the server, not a session
	held := standing(t)
	held.serve(oneDecisionDocument)

	printed, failures, code := held.run("stop", "--session", held.opened)
	if code == 0 {
		t.Fatalf("agora stop took a session: %s", printed)
	}
	if !strings.Contains(failures, "session") {
		t.Errorf("refusing agora stop --session does not name the flag it does not offer: %s", failures)
	}
}

func TestAServeWhosePortSomethingElseHoldsIsRefused(t *testing.T) {
	// scenario: The default bind is a fixed loopback port
	held := standing(t)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("holding a loopback port the surface will be asked for: %v", err)
	}
	defer func() { _ = listener.Close() }()
	taken := listener.Addr().String()

	printed, failures, code := held.reading(oneDecisionDocument, "serve", "--document", "-", "--address", taken)
	if code == 0 {
		t.Fatalf("agora serve took a port something else holds and printed %q", printed)
	}
	if !strings.Contains(failures, taken) {
		t.Errorf("refusing a port something else holds does not name it: %s", failures)
	}
	if !strings.Contains(failures, "--address") {
		t.Errorf("refusing a port something else holds does not name the flag that moves it: %s", failures)
	}
	if _, err := os.Stat(filepath.Join(held.surfaceHome(), "surface.json")); err == nil {
		t.Error("a refused bind left a surface behind, and it starts no server")
	}
}

func TestEachSessionIsWrittenUnderADirectoryOfItsOwn(t *testing.T) {
	// scenario: Neither program writes in the other's home
	held := standing(t)
	_, first, _, second := held.twoPages()

	for _, session := range []string{first, second} {
		within := filepath.Join(held.surfaceHome(), "sessions", session)
		for _, named := range []string{"session.json", "document.json"} {
			if _, err := os.Stat(filepath.Join(within, named)); err != nil {
				t.Errorf("the session %s wrote no %s: %v", session, named, err)
			}
		}
	}
}

func TestTwoSessionsAreServedAtOnceAndNeitherDisturbsTheOther(t *testing.T) {
	held := standing(t)
	firstAddress, first, secondAddress, second := held.twoPages()

	firstPage := opening(t, firstAddress)
	secondPage := opening(t, secondAddress)

	secondPage.evaluate(chose("c-1", 0), nil)
	secondPage.click(within("c-1", ".commit"))
	secondPage.click(".finish")
	secondPage.waitFor(`document.querySelector('.send').disabled === true`)

	delivered, arrived := held.awaitingOn(second, 5*time.Second)
	if !arrived {
		t.Fatal("answering on the second page delivered nothing")
	}
	if diff := cmp.Diff([]string{"patch c-1"}, delivered.Decisions[0].Commands); diff != "" {
		t.Errorf("what the second page delivered differs (-want +got):\n%s", diff)
	}

	var live bool
	firstPage.evaluate(`!document.querySelector('.send').disabled`, &live)
	if !live {
		t.Fatal("answering on one page closed the other")
	}

	firstPage.evaluate(chose("c-1", 0), nil)
	firstPage.click(within("c-1", ".commit"))
	firstPage.click(".send")
	firstPage.waitFor(`document.querySelector('.state').textContent.length > 0`)

	stillOpen, arrived := held.awaitingOn(first, 5*time.Second)
	if !arrived {
		t.Fatal("the page nobody had touched no longer answers")
	}
	if stillOpen.Ended {
		t.Error("finishing on one page reported the other as over")
	}

	nothingLeftTheMachine(t, firstAddress, firstPage.requests.all())
	nothingLeftTheMachine(t, secondAddress, secondPage.requests.all())
}

func TestAStrangerOnThePortIsNotMistakenForTheSurface(t *testing.T) {
	// scenario: The default bind is a fixed loopback port
	held := standing(t)
	held.serve(oneDecisionDocument)

	bound := held.record().Address
	if _, failures, code := held.run("stop"); code != 0 {
		t.Fatalf("agora stop exited %d: %s", code, failures)
	}
	untilThePortIsFree(t, bound)

	if err := os.WriteFile(filepath.Join(held.surfaceHome(), "surface.json"), whatAKilledSurfaceLeaves(t, bound), 0o600); err != nil {
		t.Fatalf("leaving behind what a killed surface leaves: %v", err)
	}

	stranger := aStrangerAnsweringOn(t, bound)
	defer func() { _ = stranger.Close() }()

	printed, failures, code := held.reading(oneDecisionDocument, "serve", "--document", "-", "--address", bound)
	if code == 0 {
		t.Fatalf("agora served at %q with a stranger holding the port it remembered", printed)
	}
	if !strings.Contains(failures, bound) || !strings.Contains(failures, "--address") {
		t.Errorf("what agora said is %q, and it names neither the port nor the flag that moves it", failures)
	}
	if strings.Contains(failures, "asked by whoever started it") {
		t.Errorf("agora reported a stranger on its port as a question of permission: %s", failures)
	}
	if _, err := os.Stat(filepath.Join(held.surfaceHome(), "surface.json")); err == nil {
		t.Error("agora kept a record naming a surface that is not there, so every later command repeats this")
	}
}

func whatAKilledSurfaceLeaves(t *testing.T, address string) []byte {
	t.Helper()

	written, err := json.Marshal(surface.Record{Address: address, Token: "a-token-no-surface-was-started-with", Process: 1})
	if err != nil {
		t.Fatalf("writing the record a killed surface leaves: %v", err)
	}
	return written
}

func aStrangerAnsweringOn(t *testing.T, address string) *http.Server {
	t.Helper()

	listener, err := net.Listen("tcp", address)
	if err != nil {
		t.Fatalf("taking the port the killed surface left: %v", err)
	}

	standing := &http.Server{
		Handler:           http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) { writer.WriteHeader(http.StatusOK) }),
		ReadHeaderTimeout: 10 * time.Second,
	}
	go func() { _ = standing.Serve(listener) }()
	return standing
}

func TestAnotherAgorasSurfaceIsNotMistakenForThisOne(t *testing.T) {
	// scenario: The default bind is a fixed loopback port
	shared := aFreeLoopbackAddress(t)

	foreign := standing(t)
	foreign.bind = shared
	foreign.serve(oneDecisionDocument)
	if bound := foreign.record().Address; bound != shared {
		t.Fatalf("the other agora bound %q, and the port both were asked for is %q", bound, shared)
	}

	held := servingSurface(t)
	if err := os.MkdirAll(held.surfaceHome(), 0o700); err != nil {
		t.Fatalf("making the directory the surface writes under: %v", err)
	}
	if err := os.WriteFile(filepath.Join(held.surfaceHome(), "surface.json"), whatAKilledSurfaceLeaves(t, shared), 0o600); err != nil {
		t.Fatalf("leaving behind what a killed surface leaves: %v", err)
	}

	printed, failures, code := held.reading(oneDecisionDocument, "serve", "--document", "-", "--address", shared)
	if code == 0 {
		t.Fatalf("agora served at %q on a port another agora holds", printed)
	}
	if !strings.Contains(failures, shared) || !strings.Contains(failures, "--address") {
		t.Errorf("what agora said is %q, and it names neither the port nor the flag that moves it", failures)
	}
	if strings.Contains(failures, "asked by whoever started it") {
		t.Errorf("agora reported another agora on its port as a question of permission: %s", failures)
	}
	if _, err := os.Stat(filepath.Join(held.surfaceHome(), "surface.json")); err == nil {
		t.Error("agora kept a record naming a surface that is not its own, so every later command repeats this")
	}

	if holds := foreign.record().Address; holds != shared {
		t.Errorf("the other agora's own record moved to %q", holds)
	}
	if code, _ := fetched(t, thePageAt("http://"+shared, foreign.opened)); code != http.StatusOK {
		t.Errorf("the other agora's page answered with %d after this one tried to take its port", code)
	}
}
