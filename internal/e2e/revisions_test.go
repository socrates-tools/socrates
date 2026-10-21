package e2e_test

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/domdebugger"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/google/go-cmp/cmp"
)

func (r *reader) listenersOn(expression string, subtree bool) []string {
	r.t.Helper()

	var held []string
	r.do(chromedp.Func(func(ctx context.Context, target *chromedp.Target) error {
		evaluated, err := cdp.Call(ctx, target, runtime.Evaluate, runtime.EvaluateParams{Expression: expression})
		if err != nil {
			return err
		}
		if evaluated.Result == nil || evaluated.Result.ObjectID == "" {
			return fmt.Errorf("%s names no object on the page", expression)
		}
		asked := domdebugger.GetEventListenersParams{ObjectID: evaluated.Result.ObjectID}
		if subtree {
			asked.Depth = new(int64(-1))
		}
		answered, err := cdp.Call(ctx, target, domdebugger.GetEventListeners, asked)
		if err != nil {
			return err
		}
		for _, listener := range answered.Listeners {
			held = append(held, fmt.Sprintf("%s capture=%t", listener.Type, listener.UseCapture))
		}
		return nil
	}))
	slices.Sort(held)
	return held
}

func (r *reader) listenersOfThePage() map[string][]string {
	r.t.Helper()

	return map[string][]string{
		"window":   r.listenersOn("window", false),
		"document": r.listenersOn("document", false),
		"body":     r.listenersOn("document.body", false),
	}
}

func TestSeveralRevisionsInARowLeaveNoListenerOfAnEarlierSurfaceOnThePage(t *testing.T) {
	// scenario: After several revisions in a row, the page holds no listener left from an earlier surface
	t.Parallel()
	held, page := read(t, twoDecisionDocument)

	page.evaluate(`window.firstSurface = document.querySelector('.surface')`, nil)
	standing := page.listenersOfThePage()
	if got := page.listenersOn("window.firstSurface", true); len(got) == 0 {
		t.Fatal("the first surface holds no listener at all, so a listener it left behind could not be seen")
	}

	for revision := 2; revision <= 6; revision++ {
		revised := strings.Replace(twoDecisionDocument, "two things to settle", fmt.Sprintf("revision %d", revision), 1)
		held.serveAgain(revised)
		page.waitFor(fmt.Sprintf(`document.body.dataset.revision === '%d'`, revision))
		page.waitFor(fmt.Sprintf(`document.title === 'revision %d'`, revision))
	}

	if diff := cmp.Diff(standing, page.listenersOfThePage()); diff != "" {
		t.Errorf("the window, the document and the body hold other listeners after six revisions than after the first page (-first +sixth):\n%s", diff)
	}
	if left := page.listenersOn("window.firstSurface", true); len(left) != 0 {
		t.Errorf("the first surface still holds %d listeners after it was replaced, and nothing released them: %v", len(left), left)
	}
}
