package surface_test

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/socrates-tools/socrates/agora/surface"
)

type stated struct {
	name string
	data string
}

func theSurfacesOwnOrigin(held *standing) http.Header {
	return http.Header{"Origin": {"http://" + held.server.Address()}}
}

func dialing(held *standing, within time.Duration, origin http.Header) (*websocket.Conn, int, error) {
	opening, cancel := context.WithTimeout(context.Background(), within)
	defer cancel()

	conn, answered, err := websocket.Dial(opening, held.under("/live"), &websocket.DialOptions{HTTPHeader: origin})
	if answered == nil {
		return conn, 0, err
	}
	if answered.Body != nil {
		_ = answered.Body.Close()
	}
	return conn, answered.StatusCode, err
}

func listening(t *testing.T, held *standing, origin http.Header) (*websocket.Conn, <-chan stated) {
	t.Helper()

	conn, _, err := dialing(held, 5*time.Second, origin)
	if err != nil {
		t.Fatalf("opening a stream onto the surface: %v", err)
	}
	t.Cleanup(func() { _ = conn.CloseNow() })

	arriving := make(chan stated, 8)
	go func() {
		defer close(arriving)

		for {
			_, said, err := conn.Read(context.Background())
			if err != nil {
				return
			}
			name, data, _ := strings.Cut(string(said), " ")
			arriving <- stated{name: name, data: data}
		}
	}()
	return conn, arriving
}

func events(t *testing.T, held *standing, within time.Duration) <-chan stated {
	t.Helper()

	_, arriving := listening(t, held, theSurfacesOwnOrigin(held))
	return arriving
}

func next(t *testing.T, arriving <-chan stated, within time.Duration) stated {
	t.Helper()

	select {
	case event, open := <-arriving:
		if !open {
			return stated{}
		}
		return event
	case <-time.After(within):
		return stated{}
	}
}

func eventually(t *testing.T, within time.Duration, holds func() bool) bool {
	t.Helper()

	deadline := time.Now().Add(within)
	for time.Now().Before(deadline) {
		if holds() {
			return true
		}
		time.Sleep(10 * time.Millisecond)
	}
	return holds()
}

func nobodyIsReading(held *standing, t *testing.T) func() bool {
	return func() bool {
		return held.exchange(t).IdleSince(time.Now().Add(time.Hour), time.Minute)
	}
}

func TestAStreamIsToldWhereTheExchangeStandsWhenItConnects(t *testing.T) {
	held := serving(t, twoDecisions(), 0)

	if event := next(t, events(t, held, 5*time.Second), 2*time.Second); event.name != "revision" || event.data != "1" {
		t.Errorf("a stream opened onto an exchange that is not over stated %+v with nothing further changing, want revision 1", event)
	}
}

func TestAStreamThatMissedARevisionCatchesUpOnConnectingAgain(t *testing.T) {
	held := serving(t, twoDecisions(), 0)

	revision := held.exchange(t).Revise(twoDecisions())

	event := next(t, events(t, held, 5*time.Second), 2*time.Second)
	if event.name != "revision" || event.data != strconv.Itoa(revision) {
		t.Errorf("a stream opened after a document nothing was streaming for stated %+v, want revision %d", event, revision)
	}
}

func TestAStreamOpenedAfterTheEndIsToldAtOnce(t *testing.T) {
	held := serving(t, twoDecisions(), 0)

	held.send(t, surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}}}, Ends: true})
	if _, arrived := held.await(t, 5*time.Second); !arrived {
		t.Fatal("an answer that ends the exchange was sent and nothing was delivered")
	}

	if event := next(t, events(t, held, 5*time.Second), 2*time.Second); event.name != "closed" {
		t.Errorf("a stream opened onto an exchange that is already over stated %+v with nothing further changing, want closed", event)
	}
}

func TestAnOpenStreamReportsEachChangeOnceAndThenTheClose(t *testing.T) {
	// scenario: A stream is told where the exchange stands when it connects
	held := serving(t, twoDecisions(), 0)
	arriving := events(t, held, 20*time.Second)

	if event := next(t, arriving, 5*time.Second); event.name != "revision" || event.data != "1" {
		t.Fatalf("a stream onto an open exchange stated %+v on connection, want revision 1", event)
	}
	if event := next(t, arriving, 5*time.Second); event.name != "uncollected" || event.data != "no" {
		t.Fatalf("a stream onto an open exchange stated %+v on connection, want no answer uncollected", event)
	}

	held.send(t, surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}}}})
	if event := next(t, arriving, 5*time.Second); event.name != "uncollected" || event.data != "yes" {
		t.Errorf("a stream onto an open exchange stated %+v when an answer arrived, want an answer uncollected", event)
	}

	if _, arrived := held.await(t, 5*time.Second); !arrived {
		t.Fatal("an answer that leaves the exchange open was sent and nothing was delivered")
	}
	if event := next(t, arriving, 5*time.Second); event.name != "uncollected" || event.data != "no" {
		t.Errorf("a stream onto an open exchange stated %+v once the answer was collected, want no answer uncollected", event)
	}

	revision := held.exchange(t).Revise(twoDecisions())
	if event := next(t, arriving, 5*time.Second); event.name != "revision" || event.data != strconv.Itoa(revision) {
		t.Errorf("a stream onto an open exchange stated %+v when the caller wrote a new document, want revision %d", event, revision)
	}

	held.exchange(t).End()
	if event := next(t, arriving, 5*time.Second); event.name != "closed" {
		t.Errorf("a stream onto an open exchange stated %+v when the exchange ended, want closed", event)
	}
}

func TestAStreamConnectingWhileAnAnswerIsUncollectedIsToldAtOnce(t *testing.T) {
	held := serving(t, twoDecisions(), 0)

	held.send(t, surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}}}})

	arriving := events(t, held, 5*time.Second)
	if event := next(t, arriving, 2*time.Second); event.name != "revision" || event.data != "1" {
		t.Fatalf("a stream opened while an answer is uncollected stated %+v first, want revision 1", event)
	}
	if event := next(t, arriving, 2*time.Second); event.name != "uncollected" || event.data != "yes" {
		t.Errorf("a stream opened while an answer is uncollected stated %+v, want it told at once that an answer is uncollected", event)
	}
}

func TestAnUpgradeFromAnotherOriginIsRefusedAndOneFromTheSurfacesOwnIsAccepted(t *testing.T) {
	held := serving(t, twoDecisions(), 0)

	conn, answered, err := dialing(held, 5*time.Second, http.Header{"Origin": {"http://elsewhere.example"}})
	if err == nil {
		_ = conn.CloseNow()
		t.Fatal("an upgrade from another origin was accepted")
	}
	if answered != http.StatusForbidden {
		t.Errorf("an upgrade from another origin was answered %d, want %d", answered, http.StatusForbidden)
	}

	if event := next(t, events(t, held, 5*time.Second), 2*time.Second); event.name != "revision" {
		t.Errorf("an upgrade from the surface's own origin stated %+v, want it accepted and told where the exchange stands", event)
	}
}

func TestAnUpgradeFromTheOriginASurfaceIsReachedAtIsAccepted(t *testing.T) {
	held := servingWith(t, twoDecisions(), 0, surface.ReachedAt("https://socrates.localhost/"))

	_, arriving := listening(t, held, http.Header{"Origin": {"https://socrates.localhost"}})
	if event := next(t, arriving, 2*time.Second); event.name != "revision" {
		t.Errorf("an upgrade from the origin the surface is reached at stated %+v, want it accepted", event)
	}

	conn, _, err := dialing(held, 5*time.Second, http.Header{"Origin": {"https://elsewhere.example"}})
	if err == nil {
		_ = conn.CloseNow()
		t.Error("an upgrade from an origin that is neither the surface's own nor the one it is reached at was accepted")
	}
}

func TestAnUpgradeNamingASessionTheSurfaceDoesNotHoldIsAnsweredNotFound(t *testing.T) {
	held := serving(t, twoDecisions(), 0)
	held.session = "nosuchsession"

	conn, answered, err := dialing(held, 5*time.Second, theSurfacesOwnOrigin(held))
	if err == nil {
		_ = conn.CloseNow()
		t.Fatal("an upgrade naming a session the surface does not hold was accepted")
	}
	if answered != http.StatusNotFound {
		t.Errorf("an upgrade naming a session the surface does not hold was answered %d, want %d", answered, http.StatusNotFound)
	}
}

func TestTheLiveRouteNoLongerAnswersAsAnEventStream(t *testing.T) {
	held := serving(t, twoDecisions(), 0)

	request, err := http.NewRequest(http.MethodGet, held.under("/live"), nil)
	if err != nil {
		t.Fatalf("addressing the surface: %v", err)
	}
	request.Header.Set("accept", "text/event-stream")
	answered, err := (&http.Client{Timeout: 5 * time.Second}).Do(request)
	if err != nil {
		t.Fatalf("asking the live route for an event stream: %v", err)
	}
	defer func() { _ = answered.Body.Close() }()

	if answered.StatusCode == http.StatusOK || strings.Contains(answered.Header.Get("content-type"), "text/event-stream") {
		t.Errorf("the live route answered %d as %q to a request for an event stream", answered.StatusCode, answered.Header.Get("content-type"))
	}
}

func TestAStreamCountsAsAReaderForExactlyAsLongAsItIsOpen(t *testing.T) {
	held := serving(t, twoDecisions(), 0)
	conn, arriving := listening(t, held, theSurfacesOwnOrigin(held))
	next(t, arriving, 2*time.Second)

	if nobodyIsReading(held, t)() {
		t.Fatal("an open stream was not counted as a reader of its session")
	}

	if err := conn.Close(websocket.StatusNormalClosure, ""); err != nil {
		t.Fatalf("closing the stream: %v", err)
	}
	if !eventually(t, 2*time.Second, nobodyIsReading(held, t)) {
		t.Error("a stream that closed still counted as a reader of its session")
	}
}

func TestASessionWithAStreamOpenIsNotCollectedForIdlenessAndOneWhoseLastStreamClosedIs(t *testing.T) {
	held := serving(t, twoDecisions(), 200*time.Millisecond)
	conn, arriving := listening(t, held, theSurfacesOwnOrigin(held))
	next(t, arriving, 2*time.Second)

	time.Sleep(time.Second)
	if len(held.server.Holds()) != 1 {
		t.Fatalf("a session with a stream open was collected for idleness: it holds %v", held.server.Holds())
	}

	_ = conn.Close(websocket.StatusNormalClosure, "")
	if !eventually(t, 3*time.Second, func() bool { return len(held.server.Holds()) == 0 }) {
		t.Errorf("a session whose last stream closed was not collected: it holds %v", held.server.Holds())
	}
}

func TestAStreamThatDoesNotAnswerAPingIsClosedAndStopsCountingAsAReader(t *testing.T) {
	held := servingWith(t, twoDecisions(), 0, surface.Pinging(50*time.Millisecond, 100*time.Millisecond))

	conn, _, err := dialing(held, 5*time.Second, theSurfacesOwnOrigin(held))
	if err != nil {
		t.Fatalf("opening a stream onto the surface: %v", err)
	}
	t.Cleanup(func() { _ = conn.CloseNow() })

	if !eventually(t, 5*time.Second, func() bool { return held.exchange(t).IdleSince(time.Now().Add(time.Hour), time.Minute) }) {
		t.Error("a stream that never answered a ping still counted as a reader of its session")
	}
}

func TestAStreamThatDoesNotAnswerAPingIsClosedByTheSurface(t *testing.T) {
	held := servingWith(t, twoDecisions(), 0, surface.Pinging(50*time.Millisecond, 100*time.Millisecond))

	conn, _, err := dialing(held, 5*time.Second, theSurfacesOwnOrigin(held))
	if err != nil {
		t.Fatalf("opening a stream onto the surface: %v", err)
	}
	t.Cleanup(func() { _ = conn.CloseNow() })

	if !eventually(t, 5*time.Second, nobodyIsReading(held, t)) {
		t.Fatal("a stream that never answered a ping still counted as a reader of its session")
	}

	reading, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	for {
		if _, _, err := conn.Read(reading); err != nil {
			if reading.Err() != nil {
				t.Error("a stream that never answered a ping was still open after the surface stopped counting it")
			}
			return
		}
	}
}

func TestASessionWhoseOnlyStreamStoppedAnsweringPingsIsCollectedAfterItsIdlePeriod(t *testing.T) {
	held := servingWith(t, twoDecisions(), 200*time.Millisecond, surface.Pinging(50*time.Millisecond, 100*time.Millisecond))

	conn, _, err := dialing(held, 5*time.Second, theSurfacesOwnOrigin(held))
	if err != nil {
		t.Fatalf("opening a stream onto the surface: %v", err)
	}
	t.Cleanup(func() { _ = conn.CloseNow() })

	if len(held.server.Holds()) != 1 {
		t.Fatalf("a session with a stream open was not held: it holds %v", held.server.Holds())
	}
	if !eventually(t, 5*time.Second, func() bool { return len(held.server.Holds()) == 0 }) {
		t.Errorf("a session whose only stream stopped answering pings was not collected: it holds %v", held.server.Holds())
	}
}

func TestASessionWhoseStreamKeepsAnsweringPingsIsNotCollectedAcrossSeveralIdlePeriods(t *testing.T) {
	held := servingWith(t, twoDecisions(), 200*time.Millisecond, surface.Pinging(50*time.Millisecond, 100*time.Millisecond))
	_, arriving := listening(t, held, theSurfacesOwnOrigin(held))
	next(t, arriving, 2*time.Second)

	time.Sleep(time.Second)
	if len(held.server.Holds()) != 1 {
		t.Errorf("a session whose stream answered every ping was collected: it holds %v", held.server.Holds())
	}
}

func TestAStreamThatAnswersPingsKeepsCountingAsAReader(t *testing.T) {
	held := servingWith(t, twoDecisions(), 0, surface.Pinging(50*time.Millisecond, 100*time.Millisecond))
	_, arriving := listening(t, held, theSurfacesOwnOrigin(held))
	next(t, arriving, 2*time.Second)

	time.Sleep(time.Second)
	if nobodyIsReading(held, t)() {
		t.Error("a stream that answered every ping stopped counting as a reader of its session")
	}
}

func TestStoppingTheSurfaceClosesTheStreamsItHolds(t *testing.T) {
	held := serving(t, twoDecisions(), 0)
	_, arriving := listening(t, held, theSurfacesOwnOrigin(held))
	next(t, arriving, 2*time.Second)
	next(t, arriving, 2*time.Second)

	held.server.Stop()

	select {
	case _, open := <-arriving:
		if open {
			t.Fatal("a stream said something while the surface stopped")
		}
	case <-time.After(5 * time.Second):
		t.Error("a stream was still open after the surface stopped")
	}
}

func TestAnythingAPageSendsOnTheStreamIsIgnored(t *testing.T) {
	held := serving(t, twoDecisions(), 0)
	conn, arriving := listening(t, held, theSurfacesOwnOrigin(held))
	next(t, arriving, 2*time.Second)
	next(t, arriving, 2*time.Second)

	if err := conn.Write(context.Background(), websocket.MessageText, []byte("hello")); err != nil {
		t.Fatalf("writing to the stream: %v", err)
	}

	revision := held.exchange(t).Revise(twoDecisions())
	if event := next(t, arriving, 2*time.Second); event.name != "revision" || event.data != strconv.Itoa(revision) {
		t.Errorf("a stream whose page sent something stated %+v when the caller wrote a new document, want revision %d", event, revision)
	}
}
