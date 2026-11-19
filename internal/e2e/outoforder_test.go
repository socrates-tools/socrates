package e2e_test

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type heldAnswer struct {
	request fetch.RequestID
	network network.RequestID
}

func continuing(t *testing.T, page *reader, request fetch.RequestID) {
	t.Helper()
	if _, err := chromedp.Run(page.context, sending(fetch.ContinueResponse, fetch.ContinueResponseParams{RequestID: request})); err != nil {
		t.Errorf("letting a document answer through: %v", err)
	}
}

func holdingTheFirstDocumentAnswer(t *testing.T, page *reader) <-chan heldAnswer {
	t.Helper()

	page.do(sending(fetch.Enable, fetch.EnableParams{Patterns: []*fetch.RequestPattern{{
		URLPattern:   "*/document",
		RequestStage: fetch.RequestStageResponse,
	}}}))

	first := make(chan heldAnswer, 1)
	var once sync.Once
	paused := chromedp.Events(page.context, fetch.RequestPaused)
	go func() {
		for event, err := range paused {
			if err != nil {
				return
			}
			holding := false
			once.Do(func() {
				holding = true
				first <- heldAnswer{request: event.RequestID, network: event.NetworkID}
			})
			if !holding {
				continuing(t, page, event.RequestID)
			}
		}
	}()
	return first
}

func TestAnEarlierRevisionsAnswerArrivingAfterALaterOneReplacesNothing(t *testing.T) {
	// scenario: When the answer for revision 2 arrives after the page shows revision 3, the page keeps revision 3
	t.Parallel()
	held, page := read(t, twoDecisionDocument)
	page.waitForTheStream()

	first := holdingTheFirstDocumentAnswer(t, page)
	revise := func(revision int) {
		held.serveAgain(strings.Replace(twoDecisionDocument, "two things to settle", fmt.Sprintf("revision %d", revision), 1))
	}

	revise(2)
	var late heldAnswer
	select {
	case late = <-first:
	case <-time.After(10 * time.Second):
		t.Fatal("the page never fetched the document for revision 2")
	}

	revise(3)
	page.waitFor(`document.body.dataset.revision === '3'`)
	page.waitFor(`document.title === 'revision 3'`)

	finished := make(chan struct{})
	loaded := chromedp.Events(page.context, network.LoadingFinished)
	go func() {
		defer close(finished)
		for event, err := range loaded {
			if err != nil || event.RequestID == late.network {
				return
			}
		}
	}()
	continuing(t, page, late.request)
	select {
	case <-finished:
	case <-time.After(10 * time.Second):
		t.Fatal("the answer for revision 2 never finished loading")
	}
	page.evaluate(`new Promise((settled) => requestAnimationFrame(() => requestAnimationFrame(settled)))`, nil, awaitingThePromise)

	var shown struct {
		Revision string `json:"revision"`
		Title    string `json:"title"`
	}
	page.evaluate(`({ revision: document.body.dataset.revision, title: document.title })`, &shown)
	if shown.Revision != "3" || shown.Title != "revision 3" {
		t.Errorf("after revision 2's answer arrived late the page shows revision %s titled %q, want revision 3 titled %q", shown.Revision, shown.Title, "revision 3")
	}
}
