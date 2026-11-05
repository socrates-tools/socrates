package e2e_test

import (
	"strings"
	"testing"
	"time"
)

const pagesOpenAtOnce = 7

func TestASendFromTheLastOfSevenPagesOpenOnOneSurfaceIsCollected(t *testing.T) {
	// scenario: With seven pages open on one surface, a send from the last is collected
	t.Parallel()
	held := standing(t)

	first := opening(t, held.serve(twoDecisionDocument))
	pages := []*reader{first}
	for len(pages) < pagesOpenAtOnce {
		pages = append(pages, openingAlongside(t, first, held.serve(twoDecisionDocument)))
	}
	for _, page := range pages {
		page.waitForTheStream()
	}

	last := pages[len(pages)-1]
	last.evaluate(chose("c-1", 0), nil)
	last.click(".send")

	if _, arrived := held.awaiting(10 * time.Second); !arrived {
		var said string
		last.evaluate(`document.querySelector('.state').textContent`, &said)
		t.Fatalf("a send from the last of %d open pages was not collected, and the page says %q", pagesOpenAtOnce, said)
	}
}

func TestAPageHoldsItsStreamAsAWebSocketWithoutBreakingItsPolicy(t *testing.T) {
	// scenario: A page holds its stream on the surface's own address
	t.Parallel()
	held, page := read(t, twoDecisionDocument)
	page.waitForTheStream()

	held.serveAgain(oneDecisionReply)
	page.waitFor(reachedThePage("c-3"))

	streams := page.streams.all()
	if len(streams) != 1 || !strings.HasPrefix(streams[0], "ws://"+held.record().Address+"/sessions/"+held.opened+"/live") {
		t.Errorf("the page held the streams %v, want one WebSocket on its session's live route", streams)
	}
	if violations := page.violations(); len(violations) != 0 {
		t.Errorf("the page's policy blocked %v while the page held its stream", violations)
	}
}
