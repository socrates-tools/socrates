package surface_test

import (
	"bufio"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/socrates-tools/socrates-review/internal/surface"
)

type stated struct {
	name string
	data string
}

func events(t *testing.T, held *standing, within time.Duration) <-chan stated {
	t.Helper()

	request, err := http.NewRequest(http.MethodGet, held.under("/live"), nil)
	if err != nil {
		t.Fatalf("addressing the surface: %v", err)
	}

	answered, err := (&http.Client{Timeout: within}).Do(request)
	if err != nil {
		t.Fatalf("opening a stream onto the surface: %v", err)
	}
	t.Cleanup(func() { _ = answered.Body.Close() })

	if answered.StatusCode != http.StatusOK {
		t.Fatalf("a stream onto the surface was answered with %d", answered.StatusCode)
	}

	arriving := make(chan stated, 8)
	go func() {
		defer close(arriving)

		read := bufio.NewReader(answered.Body)
		holding := stated{}
		for {
			line, err := read.ReadString('\n')
			if err != nil {
				return
			}

			line = strings.TrimSpace(line)
			switch {
			case strings.HasPrefix(line, "event: "):
				holding.name = strings.TrimPrefix(line, "event: ")
			case strings.HasPrefix(line, "data: "):
				holding.data = strings.TrimPrefix(line, "data: ")
			case line == "" && holding.name != "":
				arriving <- holding
				holding = stated{}
			}
		}
	}()
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

	held.send(t, surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}}}})
	if _, arrived := held.await(t, 5*time.Second); !arrived {
		t.Fatal("an answer that leaves the exchange open was sent and nothing was delivered")
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
