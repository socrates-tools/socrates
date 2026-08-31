package surface_test

import (
	"bufio"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/socrates-devtool/socrates/internal/surface"
)

func events(t *testing.T, held *standing, within time.Duration) <-chan string {
	t.Helper()

	request, err := http.NewRequest(http.MethodGet, held.under("/live"), nil)
	if err != nil {
		t.Fatalf("addressing the surface: %v", err)
	}

	answered, err := (&http.Client{Timeout: within}).Do(request)
	if err != nil {
		t.Fatalf("opening a stream onto the surface: %v", err)
	}
	t.Cleanup(func() { answered.Body.Close() })

	if answered.StatusCode != http.StatusOK {
		t.Fatalf("a stream onto the surface was answered with %d", answered.StatusCode)
	}

	arriving := make(chan string, 8)
	go func() {
		defer close(arriving)

		read := bufio.NewReader(answered.Body)
		for {
			line, err := read.ReadString('\n')
			if err != nil {
				return
			}
			if named, stated := strings.CutPrefix(line, "event: "); stated {
				arriving <- strings.TrimSpace(named)
			}
		}
	}()
	return arriving
}

func next(t *testing.T, arriving <-chan string, within time.Duration) string {
	t.Helper()

	select {
	case event, open := <-arriving:
		if !open {
			return ""
		}
		return event
	case <-time.After(within):
		return ""
	}
}

func TestAStreamOpenedAfterTheEndIsToldAtOnce(t *testing.T) {
	held := serving(t, twoDecisions(), 0)

	held.send(t, surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}}}, Ends: true})
	if _, arrived := held.await(t, 5*time.Second); !arrived {
		t.Fatal("an answer that ends the exchange was sent and nothing was delivered")
	}

	if event := next(t, events(t, held, 5*time.Second), 2*time.Second); event != "closed" {
		t.Errorf("a stream opened onto an exchange that is already over reported %q with nothing further changing, want %q", event, "closed")
	}
}

func TestAnOpenStreamReportsARevisionAndThenTheClose(t *testing.T) {
	held := serving(t, twoDecisions(), 0)
	arriving := events(t, held, 20*time.Second)

	held.exchange(t).Revise(twoDecisions())
	if event := next(t, arriving, 5*time.Second); event != "revision" {
		t.Errorf("a stream onto an open exchange reported %q when the caller wrote a new document, want %q", event, "revision")
	}

	held.exchange(t).End()
	if event := next(t, arriving, 5*time.Second); event != "closed" {
		t.Errorf("a stream onto an open exchange reported %q when the exchange ended, want %q", event, "closed")
	}
}
