package e2e_test

import (
	"bytes"
	"fmt"
	"image/png"
	"testing"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/chromedp"
)

const (
	tooltipAnchors  = `Array.from(document.querySelectorAll('button.catalog-anchor')).map((held) => held.id)`
	scrollableFloor = `const floor = document.createElement('div'); floor.style.height = '4000px'; document.body.appendChild(floor); window.scrollTo(0, 0);`
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
  return tip.heading === held.dataset.title && tip.hint === (held.dataset.hint || '');
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

func TestTheCatalogShowsTheTipOfAPrimitiveCarryingADataTitleOnHoverAndOnFocusAndHidesItOnLeavingAndOnScroll(t *testing.T) {
	// scenario: A data-title on a primitive's own element shows the tooltip as it does on a plain element, on hover and on focus, and hides it on leaving and on scroll
	t.Parallel()
	page := catalogued(t)
	page.do(sending(emulation.SetFocusEmulationEnabled, emulation.SetFocusEmulationEnabledParams{Enabled: true}))
	page.evaluate(scrollableFloor, nil)
	badge := badgeOf(badgeExplainedKey)
	tip := `document.querySelector('.tooltip')`
	shownWith := `(() => { const tip = ` + tip + `; return !tip.hidden && tip.getBoundingClientRect().height > 0 && tip.heading === ` + jsonOf(t, catalogBadgeTitleSaid) + ` && tip.hint === ` + jsonOf(t, catalogBadgeHintSaid) + `; })()`
	hidden := tip + `.hidden === true`
	hover := func() {
		page.evaluate(`document.querySelector(`+jsonOf(t, badge)+`).scrollIntoView({ block: 'center' }); new Promise((settled) => requestAnimationFrame(() => requestAnimationFrame(settled)))`, nil, awaitingThePromise)
		var centre struct{ X, Y float64 }
		page.evaluate(`(() => { const box = document.querySelector(`+jsonOf(t, badge)+`).getBoundingClientRect(); return { X: box.left + box.width / 2, Y: box.top + box.height / 2 }; })()`, &centre)
		page.do(sending(input.DispatchMouseEvent, input.DispatchMouseEventParams{Type: input.DispatchMouseEventTypeMouseMoved, X: centre.X, Y: centre.Y}))
	}

	hover()
	page.waitFor(shownWith)
	page.evaluate(`window.scrollBy(0, 1)`, nil)
	page.waitFor(hidden)
	page.evaluate(`window.scrollTo(0, 0)`, nil)
	leaveTheAnchors(page)

	page.do(chromedp.Focus(chromedp.CSS(badge)))
	page.waitFor(shownWith)
	page.waitFor(`(() => {
  const box = document.querySelector(` + jsonOf(t, badge) + `).getBoundingClientRect();
  return Math.abs(parseFloat(` + tip + `.style.top) - (box.top - 8)) < 1;
})()`)
	page.evaluate(`window.scrollBy(0, 40)`, nil)
	page.waitFor(shownWith)
	page.waitFor(`(() => {
  const box = document.querySelector(` + jsonOf(t, badge) + `).getBoundingClientRect();
  return Math.abs(parseFloat(` + tip + `.style.top) - (box.top - 8)) < 1;
})()`)

	page.evaluate(`document.activeElement.blur()`, nil)
	page.waitFor(hidden)
}

func TestAShadowRootHoldsEveryIndicatorTheBrandAndTheTooltipOfTheCatalogAndThePage(t *testing.T) {
	// scenario: Every status dot, overline, progress bar, brand and tooltip renders into a shadow root that adopts the shared stylesheet, on the catalog and on the page
	t.Parallel()
	const everyOneAdoptsTheSharedStylesheet = `(() => {
  const found = {};
  for (const name of ['ui-status-dot', 'ui-overline', 'ui-progress-bar', 'ui-brand', 'ui-tooltip']) {
    const drawn = Array.from(document.querySelectorAll(name));
    found[name] = drawn.length > 0 && drawn.every((held) => held.shadowRoot !== null && held.shadowRoot.adoptedStyleSheets.length > 0);
  }
  return found;
})()`
	catalog := catalogued(t)
	_, page := read(t, glancedDocument)

	for place, reader := range map[string]*reader{"catalog": catalog, "page": page} {
		var adopted map[string]bool
		reader.evaluate(everyOneAdoptsTheSharedStylesheet, &adopted)
		for name, yes := range adopted {
			if !yes {
				t.Errorf("on the %s, not every <%s> has a shadow root adopting the shared stylesheet", place, name)
			}
		}
		if violations := reader.violations(); len(violations) != 0 {
			t.Errorf("the %s's policy blocked %v", place, violations)
		}
	}
}

func (r *reader) progressColourAt(selector string, fraction float64) string {
	r.t.Helper()

	r.evaluate(`document.querySelector(`+jsonOf(r.t, selector)+`).scrollIntoView({ block: 'center' }); new Promise((settled) => requestAnimationFrame(() => requestAnimationFrame(settled)))`, nil, awaitingThePromise)
	var box struct{ Left, Top, Width, Height float64 }
	r.evaluate(`(() => { const box = document.querySelector(`+jsonOf(r.t, selector)+`).getBoundingClientRect(); return { Left: box.left, Top: box.top, Width: box.width, Height: box.height }; })()`, &box)
	shot, err := png.Decode(bytes.NewReader(r.screenshot()))
	if err != nil {
		r.t.Fatalf("reading what the page painted: %v", err)
	}
	red, green, blue, _ := shot.At(int(box.Left+box.Width*fraction), int(box.Top+box.Height/2)).RGBA()
	return fmt.Sprintf("rgb(%d, %d, %d)", red>>8, green>>8, blue>>8)
}

func TestTheCatalogPaintsAProgressBarsDoneShareInTheAccentAndTheRestInTheTrackColour(t *testing.T) {
	// scenario: A progress bar state paints the share done in the accent and the share left in the track's colour
	t.Parallel()
	page := catalogued(t)
	accent := page.pageColourOf("--accent")
	track := page.pageColourOf("--surface-track")

	for _, state := range []struct {
		key        string
		done, left []float64
	}{
		{"nothing-done", nil, []float64{0.1, 0.5, 0.9}},
		{"partly-done", []float64{0.1, 0.4}, []float64{0.6, 0.9}},
		{"all-done", []float64{0.1, 0.5, 0.9}, nil},
	} {
		t.Run(state.key, func(t *testing.T) {
			bar := stateOf("ui-progress-bar", state.key) + " ui-progress-bar"
			for _, at := range state.done {
				if got := page.progressColourAt(bar, at); got != accent {
					t.Errorf("at %v of the bar it is %s, want the accent %s", at, got, accent)
				}
			}
			for _, at := range state.left {
				if got := page.progressColourAt(bar, at); got != track {
					t.Errorf("at %v of the bar it is %s, want the track's %s", at, got, track)
				}
			}
		})
	}
}

func TestTheCatalogMovesTheProgressBarWhenItsValueIsSet(t *testing.T) {
	// scenario: Setting the value of a progress bar in the catalog moves its bar
	t.Parallel()
	page := catalogued(t)
	bar := `[data-primitive="ui-progress-bar"] ui-progress-bar`
	accent := page.pageColourOf("--accent")

	var before struct{ Value, Max float64 }
	page.evaluate(`(() => { const bar = document.querySelector(`+jsonOf(t, bar)+`); return { Value: bar.value, Max: bar.max }; })()`, &before)
	if page.progressColourAt(bar, 0.9) == accent {
		t.Fatalf("the first bar the catalog draws is already done to %v of %v, so this proves nothing", before.Value, before.Max)
	}

	page.evaluate(`(async () => { const bar = document.querySelector(`+jsonOf(t, bar)+`); bar.value = bar.max; await bar.updateComplete; })()`, nil, awaitingThePromise)
	page.waitFor(`document.querySelector(` + jsonOf(t, bar) + `).value === document.querySelector(` + jsonOf(t, bar) + `).max`)

	page.evaluate(`new Promise((settled) => setTimeout(settled, 400))`, nil, awaitingThePromise)
	if got := page.progressColourAt(bar, 0.9); got != accent {
		t.Errorf("setting a progress bar to its maximum leaves %s at nine tenths of it, want the accent %s", got, accent)
	}
}
