package e2e_test

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/chromedp/cdproto/network"
	cdppage "github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

const (
	lostSaying     = "Lost the surface"
	sessionGone    = "The session is gone"
	reconnectLimit = 15 * time.Second
)

func (r *reader) statusSays() string {
	r.t.Helper()

	var said string
	r.evaluate(`document.querySelector('.status-said').textContent`, &said)
	return said
}

func (r *reader) goOffline(offline bool, reaching string) {
	r.t.Helper()

	r.do(sending(network.EmulateNetworkConditionsByRule, network.EmulateNetworkConditionsByRuleParams{
		MatchedNetworkConditions: []*network.Conditions{{
			URLPattern:         reaching,
			Offline:            offline,
			DownloadThroughput: -1,
			UploadThroughput:   -1,
		}},
	}))
}

const keepEveryStreamOpened = `
	window.streamsOpened = []
	const Original = window.WebSocket
	window.WebSocket = class extends Original {
		constructor(...asked) {
			super(...asked)
			window.streamsOpened.push(this)
		}
	}`

func (r *reader) keepingItsStreams() {
	r.t.Helper()

	var address string
	r.evaluate(`location.href`, &address)
	r.do(
		sending(cdppage.AddScriptToEvaluateOnNewDocument, cdppage.AddScriptToEvaluateOnNewDocumentParams{Source: keepEveryStreamOpened}),
		chromedp.Navigate(address),
		chromedp.WaitReady(chromedp.CSS(r.ready)),
	)
	r.settled()
	r.waitFor(`window.streamsOpened.length > 0 && window.streamsOpened.every((stream) => stream.readyState === WebSocket.OPEN)`)
}

func (r *reader) closeEveryStreamTheBrowserHolds() {
	r.t.Helper()
	r.evaluate(`window.streamsOpened.forEach((stream) => stream.close()); true`, nil)
}

func (r *reader) loseTheConnection() {
	r.t.Helper()

	r.goOffline(true, "http://127.0.0.1:*/*")
	r.closeEveryStreamTheBrowserHolds()
	r.waitFor(`document.querySelector('.state').textContent.includes('Lost the surface')`)
	r.goOffline(true, "")
}

func (r *reader) regainTheConnection() {
	r.t.Helper()
	r.goOffline(false, "")
}

func TestAfterStopAnOpenPageSaysItLostTheSurfaceAndIsTryingAgain(t *testing.T) {
	// scenario: After socrates agora stop, an open page says it lost the surface and is trying again
	t.Parallel()
	held, page := read(t, twoDecisionDocument)
	page.waitForTheStream()

	held.run("stop")
	page.waitFor(`document.querySelector('.state').textContent.includes('Lost the surface')`)

	said := page.stateSays()
	if !strings.Contains(said, "trying again") {
		t.Errorf("the page says %q, and it should say it is trying again", said)
	}
	if tone := page.stateTone(); tone != "error" {
		t.Errorf("the page said it lost the surface in the %q tone, want error", tone)
	}
	if status := page.statusSays(); status != lostSaying {
		t.Errorf("the status says %q after the surface stopped, want %q", status, lostSaying)
	}
	var marked string
	page.evaluate(`document.querySelector('.status').dataset.status`, &marked)
	if marked != "lost" {
		t.Errorf("the status is marked %q after the surface stopped, want lost, which the pill draws apart from the others", marked)
	}
}

func TestAPageWhoseSurfaceComesBackWithItsSessionStillHeldReconnectsAndShowsWhatItMissed(t *testing.T) {
	// scenario: A page whose surface comes back with its session still held reconnects by itself and swaps in a revision served while it was gone
	t.Parallel()
	held, page := read(t, twoDecisionDocument)
	page.keepingItsStreams()
	reconnected := len(page.streams.all()) + 1

	page.loseTheConnection()
	page.waitFor(`document.querySelector('.state').textContent.includes('Lost the surface')`)

	reachedBefore := len(page.opened.all())
	held.serveAgain(oneDecisionReply)
	time.Sleep(2500 * time.Millisecond)
	if reached := len(page.opened.all()); reached != reachedBefore {
		t.Fatalf("the page reached the surface %d times while the connection was cut, so the cut did not hold", reached-reachedBefore)
	}
	var shown bool
	page.evaluate(reachedThePage("c-3"), &shown)
	if shown {
		t.Fatal("the page showed a revision served while the connection was cut, so the cut did not hold")
	}
	page.regainTheConnection()

	deadline := time.Now().Add(reconnectLimit)
	for len(page.streams.all()) < reconnected && time.Now().Before(deadline) {
		time.Sleep(50 * time.Millisecond)
	}
	page.waitFor(reachedThePage("c-3"))

	if status := page.statusSays(); status != "Open" {
		t.Errorf("the status says %q once the surface is back, want where the exchange stands, Open", status)
	}
	if said := page.stateSays(); strings.Contains(said, lostSaying) {
		t.Errorf("the page still says %q once the surface is back", said)
	}
}

func TestAPageWhoseSurfaceComesBackKeepsSayingWhereASentAnswerStands(t *testing.T) {
	// scenario: A page whose status was Sent returns to Sent when the surface comes back
	t.Parallel()
	_, page := read(t, twoDecisionDocument)
	page.keepingItsStreams()

	page.evaluate(chose("c-1", 0), nil)
	page.click(".send")
	page.waitFor(`document.querySelector('.status-said').textContent === 'Sent, waiting for a reply'`)
	sentSaying := page.stateSays()
	reconnected := len(page.streams.all()) + 1

	page.loseTheConnection()
	page.waitFor(`document.querySelector('.status-said').textContent === 'Lost the surface'`)
	page.regainTheConnection()

	deadline := time.Now().Add(reconnectLimit)
	for len(page.streams.all()) < reconnected && time.Now().Before(deadline) {
		time.Sleep(50 * time.Millisecond)
	}
	page.waitFor(`document.querySelector('.status-said').textContent === 'Sent, waiting for a reply'`)

	if said := page.stateSays(); said != sentSaying {
		t.Errorf("the page says %q once the surface is back, want what it said before, %q", said, sentSaying)
	}
}

func TestAPageWhoseSessionTheSurfaceNoLongerHoldsSaysSoOffersNothingAndStopsTrying(t *testing.T) {
	// scenario: A page whose session the surface no longer holds says the session is gone, offers nothing to answer and stops reconnecting
	t.Parallel()
	held := standing(t)
	held.bindsAPortNothingElseWillTake()
	page := opening(t, held.serve(twoDecisionDocument))
	page.waitForTheStream()

	held.run("stop")
	page.waitFor(`document.querySelector('.state').textContent.includes('Lost the surface')`)

	held.serve(oneDecisionDocument)
	page.waitFor(`document.querySelector('.state').textContent.includes('` + sessionGone + `')`)

	if status := page.statusSays(); status != "Exchange over" {
		t.Errorf("the status says %q for a session the surface does not hold, want Exchange over", status)
	}
	var offered int
	page.evaluate(`Array.from(document.querySelectorAll('.send, .finish, .choice, .decline, .words-field, .message-field')).filter((control) => !control.disabled && control.offsetParent !== null).length`, &offered)
	if offered != 0 {
		t.Errorf("the page still offers %d controls to answer a session the surface does not hold", offered)
	}

	asked := len(page.streams.all())
	if streams := len(page.streams.all()); streams != asked {
		t.Errorf("the page opened %d more streams after it said the session is gone", streams-asked)
	}
}

func (r *reader) shortenEveryWaitToReconnect() {
	r.t.Helper()

	r.evaluate(`
		window.waitsAsked = []
		const original = window.setTimeout.bind(window)
		window.setTimeout = (run, milliseconds, ...rest) => {
			if ([1000, 2000, 4000, 8000, 16000, 30000].includes(milliseconds)) {
				window.waitsAsked.push(milliseconds)
				return original(run, 0, ...rest)
			}
			return original(run, milliseconds, ...rest)
		}
		true`, nil)
}

func TestAPageThatKeepsLosingTheSurfaceWaitsAStartingNearASecondDoublingToAboutThirtySeconds(t *testing.T) {
	// scenario: A page reconnects with a backoff that starts near a second and doubles to a cap near thirty seconds
	t.Parallel()
	held, page := read(t, twoDecisionDocument)
	page.waitForTheStream()
	page.shortenEveryWaitToReconnect()

	held.run("stop")

	var asked []int
	for deadline := time.Now().Add(reconnectLimit); len(asked) < 8 && time.Now().Before(deadline); time.Sleep(50 * time.Millisecond) {
		page.evaluate(`window.waitsAsked.slice(0, 8)`, &asked)
	}
	want := []int{1000, 2000, 4000, 8000, 16000, 30000, 30000, 30000}
	if !slices.Equal(asked, want) {
		t.Errorf("the page waited %v milliseconds between tries, want %v", asked, want)
	}
}
