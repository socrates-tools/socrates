package e2e_test

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

const (
	aMinuteInMilliseconds = 60000
	shortenedTimeout      = 300
)

type heldAnswers struct {
	mu      sync.Mutex
	paused  []fetch.RequestID
	arrived chan struct{}
}

func (r *reader) holdingAnswers() *heldAnswers {
	r.t.Helper()

	held := &heldAnswers{arrived: make(chan struct{}, 16)}
	paused := chromedp.Events(r.context, fetch.RequestPaused)
	go func() {
		for stopped, err := range paused {
			if err != nil {
				return
			}
			held.mu.Lock()
			held.paused = append(held.paused, stopped.RequestID)
			held.mu.Unlock()
			held.arrived <- struct{}{}
		}
	}()
	r.do(sending(fetch.Enable, fetch.EnableParams{Patterns: []*fetch.RequestPattern{{URLPattern: "*/answer"}}}))
	return held
}

func (h *heldAnswers) count() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.paused)
}

func (h *heldAnswers) awaitTheFirst(t *testing.T) fetch.RequestID {
	t.Helper()

	select {
	case <-h.arrived:
	case <-time.After(10 * time.Second):
		t.Fatal("the press sent no request at all to hold")
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.paused[0]
}

func (r *reader) refuseTheRequest(held fetch.RequestID) {
	r.t.Helper()
	r.do(sending(fetch.FailRequest, fetch.FailRequestParams{RequestID: held, ErrorReason: network.ErrorReasonConnectionRefused}))
}

func (r *reader) shortenTheSendTimeout() {
	r.t.Helper()

	r.evaluate(fmt.Sprintf(`
		window.timeoutsAsked = []
		const original = AbortSignal.timeout.bind(AbortSignal)
		AbortSignal.timeout = (milliseconds) => {
			window.timeoutsAsked.push(milliseconds)
			return original(milliseconds >= %d ? %d : milliseconds)
		}
		true`, aMinuteInMilliseconds, shortenedTimeout), nil)
}

const bothPressesDisabled = `document.querySelector('.send').disabled === true && document.querySelector('.finish').disabled === true`

const bothPressesEnabled = `document.querySelector('.send').disabled === false && document.querySelector('.finish').disabled === false`

func (r *reader) stateSays() string {
	r.t.Helper()

	var said string
	r.evaluate(`document.querySelector('.state').textContent`, &said)
	return said
}

func (r *reader) stateTone() string {
	r.t.Helper()

	var tone string
	r.evaluate(`document.querySelector('.state').dataset.tone`, &tone)
	return tone
}

func writtenAndChosen(t *testing.T, page *reader) {
	t.Helper()

	page.evaluate(chose("c-1", 0), nil)
	page.write(".message-field", "still here")
}

func stillOnThePage(t *testing.T, page *reader) {
	t.Helper()

	var chosen int
	page.evaluate(`document.querySelectorAll('.decision[data-decision="c-1"] .choice[checked]').length`, &chosen)
	if chosen != 1 {
		t.Errorf("a failed send left %d options ticked on c-1, want the one that was chosen", chosen)
	}
	var written string
	page.evaluate(`document.querySelector('.message-field').value`, &written)
	if written != "still here" {
		t.Errorf("a failed send left the message as %q, want what was written", written)
	}
	var locked bool
	page.evaluate(`Array.from(document.querySelectorAll('.choice, .decline, .words-field, .message-field, .tray')).some((control) => control.disabled)`, &locked)
	if locked {
		t.Error("a failed send left a control read-only, and what the person wrote cannot be changed to send again")
	}
}

func TestAPressWhileASendIsInFlightDisablesBothPressesAndSendsNothingMore(t *testing.T) {
	t.Parallel()
	held, page := read(t, twoDecisionDocument)
	answers := page.holdingAnswers()

	writtenAndChosen(t, page)
	page.click(".send")
	answers.awaitTheFirst(t)

	page.waitFor(bothPressesDisabled)

	page.evaluate(`document.querySelector('form.exchange').requestSubmit()`, nil)
	page.evaluate(`document.querySelector('.send').click(); document.querySelector('.finish').click(); true`, nil)

	if requested := answers.count(); requested != 1 {
		t.Errorf("pressing again while a send was in flight left %d requests held, want the one", requested)
	}
	if delivered, arrived := held.awaiting(shortWindow); arrived {
		t.Errorf("an answer that was never let through reached the caller: %+v", delivered)
	}
}

func TestASendTheSurfaceHoldsUnansweredForAMinuteIsAbortedAndSaysNothingWasSent(t *testing.T) {
	t.Parallel()
	held, page := read(t, twoDecisionDocument)
	answers := page.holdingAnswers()
	page.shortenTheSendTimeout()

	writtenAndChosen(t, page)
	page.click(".send")
	answers.awaitTheFirst(t)
	page.waitFor(`document.querySelector('.state').textContent.includes('Nothing was sent')`)

	var asked []int
	page.evaluate(`window.timeoutsAsked`, &asked)
	if len(asked) != 1 || asked[0] != aMinuteInMilliseconds {
		t.Errorf("the page limited the send to %v milliseconds, want one limit of %d", asked, aMinuteInMilliseconds)
	}

	said := page.stateSays()
	if !strings.Contains(strings.ToLower(said), "minute") {
		t.Errorf("the page says %q, and it should name why nothing was sent", said)
	}
	if tone := page.stateTone(); tone != "error" {
		t.Errorf("the page said nothing was sent in the %q tone, want error", tone)
	}

	page.waitFor(bothPressesEnabled)
	stillOnThePage(t, page)
	if delivered, arrived := held.awaiting(shortWindow); arrived {
		t.Errorf("an aborted send reached the caller: %+v", delivered)
	}
}

func TestASendTheBrowserRejectsSaysNothingWasSentRatherThanStayingOnSending(t *testing.T) {
	t.Parallel()
	_, page := read(t, twoDecisionDocument)
	answers := page.holdingAnswers()

	writtenAndChosen(t, page)
	page.click(".send")
	page.waitFor(`document.querySelector('.state').textContent === 'Sending.'`)
	page.refuseTheRequest(answers.awaitTheFirst(t))
	page.waitFor(`document.querySelector('.state').textContent.includes('Nothing was sent')`)

	if tone := page.stateTone(); tone != "error" {
		t.Errorf("the page said nothing was sent in the %q tone, want error", tone)
	}

	page.waitFor(bothPressesEnabled)
	stillOnThePage(t, page)
}

func TestAnAnswerSentAfterAFailedSendIsDelivered(t *testing.T) {
	t.Parallel()
	held, page := read(t, twoDecisionDocument)
	answers := page.holdingAnswers()

	writtenAndChosen(t, page)
	page.click(".send")
	page.refuseTheRequest(answers.awaitTheFirst(t))
	page.waitFor(`document.querySelector('.state').textContent.includes('Nothing was sent')`)
	page.waitFor(bothPressesEnabled)

	page.do(sending(fetch.Disable, cdp.Empty{}))
	page.click(".send")

	delivered, arrived := held.awaiting(5 * time.Second)
	if !arrived {
		t.Fatal("sending again after a failed send delivered nothing")
	}
	if delivered.Message != "still here" {
		t.Errorf("the message came back as %q, want what was written before the failed send", delivered.Message)
	}
}

const theSendWentThrough = `document.querySelector('.send').disabled === true && !document.querySelector('.receipt').hidden`
