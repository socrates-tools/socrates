package e2e_test

import (
	"strings"
	"testing"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)

func TestAPageInABrowserThatRunsNoScriptSaysItNeedsJavaScriptAndShowsNothingOfTheDocument(t *testing.T) {
	// scenario: A browser that runs no script is shown one line and nothing of the document
	t.Parallel()
	held := standing(t)
	address := held.serve(oneDecisionDocument)

	page := browsing(t, ".send", ".decision")
	page.do(
		sending(emulation.SetScriptExecutionDisabled, emulation.SetScriptExecutionDisabledParams{Value: true}),
		chromedp.Navigate(address),
		chromedp.WaitReady(chromedp.CSS("body")),
	)

	var shown string
	page.evaluate(`document.body.innerText.trim()`, &shown)
	if !strings.Contains(shown, "needs JavaScript") {
		t.Errorf("a page that runs no script shows %q, and it does not say it needs JavaScript", shown)
	}
	if lines := strings.Split(shown, "\n"); len(lines) != 1 {
		t.Errorf("a page that runs no script shows %d lines, want one: %q", len(lines), shown)
	}

	for _, kept := range []string{"one thing to settle", "what follows is what was seen", "should the cache be read before the lock is taken", "take the lock first", "correctness first"} {
		if strings.Contains(shown, kept) {
			t.Errorf("a page that runs no script shows %q, which only the document says", kept)
		}
	}

	var drawn bool
	page.evaluate(`!!document.querySelector('.surface, .decision, ui-card, ui-choice')`, &drawn)
	if drawn {
		t.Error("a page that runs no script holds a piece of the surface, and it holds none")
	}
}
