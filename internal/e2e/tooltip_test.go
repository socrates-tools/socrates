package e2e_test

import (
	"fmt"
	"testing"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/chromedp"
)

const (
	tooltipAnchors  = `Array.from(document.querySelectorAll('button.catalog-anchor')).map((held) => held.id)`
	scrollableFloor = `document.body.style.paddingBottom = '4000px'; window.scrollTo(0, 0);`
)

func tooltipOf(anchor string) string {
	return fmt.Sprintf(`document.getElementById(document.getElementById(%q).dataset.tooltip)`, anchor)
}

func isShown(anchor string) string {
	return fmt.Sprintf(`(() => { const tip = %s; return !tip.hidden && tip.getBoundingClientRect().height > 0; })()`, tooltipOf(anchor))
}

func isHidden(anchor string) string {
	return tooltipOf(anchor) + `.hidden === true`
}

func saysTheTitleAndHintOf(anchor string) string {
	return fmt.Sprintf(`(() => {
  const held = document.getElementById(%q);
  const tip = %s;
  return tip.querySelector('[data-part="title"]').textContent === held.dataset.title && tip.querySelector('[data-part="hint"]').textContent === (held.dataset.hint || '');
})()`, anchor, tooltipOf(anchor))
}

func hoverOver(page *reader, anchor string) {
	page.t.Helper()
	page.evaluate(fmt.Sprintf(`document.getElementById(%q).scrollIntoView({ block: 'center' }); new Promise((settled) => requestAnimationFrame(() => requestAnimationFrame(settled)))`, anchor), nil, awaitingThePromise)
	var centre struct{ X, Y float64 }
	page.evaluate(fmt.Sprintf(`(() => { const box = document.getElementById(%q).getBoundingClientRect(); return { X: box.left + box.width / 2, Y: box.top + box.height / 2 }; })()`, anchor), &centre)
	page.do(sending(input.DispatchMouseEvent, input.DispatchMouseEventParams{Type: input.DispatchMouseEventTypeMouseMoved, X: centre.X, Y: centre.Y}))
}

func leaveTheAnchors(page *reader) {
	page.t.Helper()
	page.do(sending(input.DispatchMouseEvent, input.DispatchMouseEventParams{Type: input.DispatchMouseEventTypeMouseMoved, X: 1, Y: 1}))
}

func eachTooltipState(t *testing.T, page *reader, prove func(t *testing.T, anchor string)) {
	t.Helper()

	var anchors []string
	page.evaluate(tooltipAnchors, &anchors)
	if len(anchors) < 2 {
		t.Fatalf("the catalog offers %d anchors for a tooltip, and the tooltip has more than one state to draw: %v", len(anchors), anchors)
	}
	for _, anchor := range anchors {
		t.Run(anchor, func(t *testing.T) {
			prove(t, anchor)
		})
	}
}

func TestTheCatalogShowsEachTooltipStateOnHoverAndHidesItWhenTheCursorLeaves(t *testing.T) {
	// scenario: A tooltip in the catalog shows what its anchor says while the cursor is over the anchor
	t.Parallel()
	page := catalogued(t)

	eachTooltipState(t, page, func(t *testing.T, anchor string) {
		page.waitFor(isHidden(anchor))

		hoverOver(page, anchor)

		page.waitFor(isShown(anchor))
		page.waitFor(saysTheTitleAndHintOf(anchor))

		leaveTheAnchors(page)
		page.waitFor(isHidden(anchor))
	})
}

func TestTheCatalogShowsEachTooltipStateOnFocusAndHidesItWhenTheCursorLeaves(t *testing.T) {
	// scenario: A tooltip in the catalog shows what its anchor says while the anchor has the keyboard cursor
	t.Parallel()
	page := catalogued(t)
	page.do(sending(emulation.SetFocusEmulationEnabled, emulation.SetFocusEmulationEnabledParams{Enabled: true}))

	eachTooltipState(t, page, func(t *testing.T, anchor string) {
		page.do(chromedp.Focus(chromedp.CSS("#" + anchor)))

		page.waitFor(isShown(anchor))
		page.waitFor(saysTheTitleAndHintOf(anchor))

		page.evaluate(`document.activeElement.blur()`, nil)
		page.waitFor(isHidden(anchor))
	})
}

func TestTheCatalogHidesEachTooltipStateOnScrollWhileTheCursorIsOnlyOverItsAnchor(t *testing.T) {
	// scenario: A tooltip in the catalog shown by hover hides as soon as the page scrolls
	t.Parallel()
	page := catalogued(t)
	page.evaluate(scrollableFloor, nil)

	eachTooltipState(t, page, func(t *testing.T, anchor string) {
		hoverOver(page, anchor)
		page.waitFor(isShown(anchor))

		page.evaluate(`window.scrollBy(0, 1)`, nil)

		page.waitFor(isHidden(anchor))

		page.evaluate(`window.scrollTo(0, 0)`, nil)
		leaveTheAnchors(page)
	})
}

func TestTheCatalogKeepsEachTooltipStateWithTheFocusedAnchorWhenThePageScrolls(t *testing.T) {
	// scenario: A tooltip in the catalog shown by focus follows its anchor when the page scrolls
	t.Parallel()
	page := catalogued(t)
	page.do(sending(emulation.SetFocusEmulationEnabled, emulation.SetFocusEmulationEnabledParams{Enabled: true}))
	page.evaluate(scrollableFloor, nil)

	eachTooltipState(t, page, func(t *testing.T, anchor string) {
		page.do(chromedp.Focus(chromedp.CSS("#" + anchor)))
		page.waitFor(isShown(anchor))

		page.evaluate(`window.scrollBy(0, 40)`, nil)

		page.waitFor(isShown(anchor))
		page.waitFor(fmt.Sprintf(`(() => {
  const tip = %s;
  const box = document.getElementById(%q).getBoundingClientRect();
  return Math.abs(parseFloat(tip.style.top) - (box.top - 8)) < 1;
})()`, tooltipOf(anchor), anchor))

		page.evaluate(`document.activeElement.blur(); window.scrollTo(0, 0)`, nil)
	})
}

func TestTheCatalogMovesTheProgressBarWhenItsValueIsSet(t *testing.T) {
	// scenario: Setting the value of a progress bar in the catalog moves its bar
	t.Parallel()
	page := catalogued(t)

	var moved struct {
		Before float64 `json:"before"`
		After  float64 `json:"after"`
		Track  float64 `json:"track"`
		Max    float64 `json:"max"`
	}
	page.evaluate(`(() => {
  const bar = document.querySelector('[data-primitive="ui-progress-bar"] ui-progress-bar');
  const before = bar.value;
  bar.value = bar.max;
  return { before, after: bar.value, track: bar.querySelector('[data-part="track"]').value, max: bar.max };
})()`, &moved)

	if moved.After != moved.Max || moved.Track != moved.Max || moved.Before == moved.After {
		t.Errorf("setting a progress bar to its maximum leaves %+v, and the bar the person sees must follow its value", moved)
	}
}
