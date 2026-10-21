package e2e_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/chromedp/kb"
)

const drawingWithACounter = `<!doctype html>
<html><body>
<p id="count">0</p>
<script>
let count = 0;
window.addEventListener('message', (event) => {
  if (typeof event.data === 'string' && event.data.startsWith('bump')) {
    count += 1;
    document.getElementById('count').textContent = String(count);
  }
  parent.postMessage({ count, shown: document.getElementById('count').textContent, answering: event.data }, '*');
});
</script>
</body></html>`

const expandControl = `.decision[data-decision="c-1"] ui-expand.expand`

const whereTheFrameStands = `(() => {
  const frame = ` + frameOfC1 + `;
  const box = frame.getBoundingClientRect();
  const backdrop = document.querySelector('.decision[data-decision="c-1"] .drawing-backdrop');
  const shade = backdrop ? backdrop.getBoundingClientRect() : { width: 0, height: 0 };
  return {
    position: getComputedStyle(frame).position,
    left: box.left, top: box.top, right: box.right, bottom: box.bottom,
    width: window.innerWidth, height: window.innerHeight,
    shaded: shade.width >= window.innerWidth && shade.height >= window.innerHeight && getComputedStyle(backdrop).display !== 'none',
    heldFocus: document.activeElement === document.querySelector('` + expandControl + `'),
    closeShown: (() => { const close = document.querySelector('.decision[data-decision="c-1"] .drawing-close'); return !!close && close.getBoundingClientRect().width > 0; })(),
  };
})()`

type frameStanding struct {
	Position   string  `json:"position"`
	Left       float64 `json:"left"`
	Top        float64 `json:"top"`
	Right      float64 `json:"right"`
	Bottom     float64 `json:"bottom"`
	Width      float64 `json:"width"`
	Height     float64 `json:"height"`
	Shaded     bool    `json:"shaded"`
	HeldFocus  bool    `json:"heldFocus"`
	CloseShown bool    `json:"closeShown"`
}

func (f frameStanding) fillsTheWindow() bool {
	return f.Position == "fixed" && f.Left <= 32 && f.Top <= 64 && f.Width-f.Right <= 32 && f.Height-f.Bottom <= 32 &&
		f.Right-f.Left >= f.Width-64 && f.Bottom-f.Top >= f.Height-96
}

func (r *reader) whereTheFrameStands() frameStanding {
	r.t.Helper()

	var stands frameStanding
	r.evaluate(whereTheFrameStands, &stands)
	return stands
}

func (r *reader) tellTheDrawing(message string) map[string]any {
	r.t.Helper()

	r.evaluate(`window.__heard = {}; if (!window.__listening) {
  window.__listening = true;
  window.__asked = 0;
  window.addEventListener('message', (event) => {
    if (event.source === `+frameOfC1+`.contentWindow && typeof event.data === 'object') window.__heard[event.data.answering] = event.data;
  });
} window.__telling = (told) => { window.__asked += 1; const tag = told + '#' + window.__asked; `+frameOfC1+`.contentWindow.postMessage(tag, '*'); return tag; }; true`, nil)

	r.waitFor(`(() => { window.__telling('ask'); return Object.keys(window.__heard).length > 0; })()`)

	var tag string
	r.evaluate(`window.__heard = {}; window.__telling(`+jsonOf(r.t, message)+`)`, &tag)
	r.waitFor(`window.__heard[` + jsonOf(r.t, tag) + `] !== undefined`)

	var heard map[string]any
	r.evaluate(`window.__heard[`+jsonOf(r.t, tag)+`]`, &heard)
	return heard
}

func aPageWithACountingDrawing(t *testing.T) *reader {
	t.Helper()

	held := standing(t)
	aDrawing(t, held, "counting.html", drawingWithACounter)
	page := opening(t, held.serve(illustrating(t, oneDecisionDocument, "counting.html")))
	page.waitFor(frameOfC1 + `.getBoundingClientRect().height > 0`)
	return page
}

func TestAnExpandedDrawingFillsTheWindowAtEveryWidthThePageSupports(t *testing.T) {
	// scenario: Pressing the expand control shows the drawing filling the window at every width, including below 960px
	t.Parallel()

	for _, width := range []int64{1280, 960, 800, 420} {
		t.Run(fmt.Sprintf("%d wide", width), func(t *testing.T) {
			page := aPageWithACountingDrawing(t)
			page.narrow(width, 700)

			if before := page.whereTheFrameStands(); before.fillsTheWindow() || before.CloseShown {
				t.Fatalf("the card already shows its drawing over the window before expanding: %+v", before)
			}

			page.click(expandControl)

			after := page.whereTheFrameStands()
			if !after.fillsTheWindow() {
				t.Errorf("the expanded drawing stands at %+v, want it laid over the whole window", after)
			}
			if !after.Shaded {
				t.Errorf("the page is not dimmed behind the expanded drawing: %+v", after)
			}
			if !after.CloseShown {
				t.Errorf("the expanded drawing offers no close button: %+v", after)
			}
		})
	}
}

func TestEscapeTheCloseButtonAndTheBackdropEachReturnToTheCardWithTheCursorOnTheExpandControl(t *testing.T) {
	// scenario: Escape, the close button and a press on the backdrop each return to the card with the cursor on the expand control
	t.Parallel()

	for _, closing := range []struct {
		by string
		do func(page *reader)
	}{
		{"escape", func(page *reader) { page.press(kb.Escape) }},
		{"the close button", func(page *reader) { page.click(".decision[data-decision=\"c-1\"] .drawing-close") }},
		{"a press on the backdrop", func(page *reader) { page.clickAt(4, 4) }},
	} {
		t.Run(closing.by, func(t *testing.T) {
			page := aPageWithACountingDrawing(t)
			page.click(expandControl)
			if !page.whereTheFrameStands().fillsTheWindow() {
				t.Fatal("the drawing did not expand")
			}

			closing.do(page)

			after := page.whereTheFrameStands()
			if after.fillsTheWindow() || after.CloseShown || after.Position == "fixed" {
				t.Errorf("closing by %s leaves the drawing over the window: %+v", closing.by, after)
			}
			if !after.HeldFocus {
				t.Errorf("closing by %s leaves the cursor off the expand control", closing.by)
			}
		})
	}
}

func TestADrawingKeepsTheStateItsScriptBuiltAcrossExpandingAndClosing(t *testing.T) {
	// scenario: A drawing whose script changed its state before expanding shows the same state expanded and after closing
	t.Parallel()
	page := aPageWithACountingDrawing(t)

	page.tellTheDrawing("bump")
	if before := page.tellTheDrawing("bump"); before["count"] != float64(2) {
		t.Fatalf("the drawing counted %v before expanding, want 2", before)
	}

	page.click(expandControl)
	if expanded := page.tellTheDrawing("ask"); expanded["count"] != float64(2) || expanded["shown"] != "2" {
		t.Errorf("the expanded drawing shows %v, want the state it built before expanding", expanded)
	}

	page.press(kb.Escape)
	if closed := page.tellTheDrawing("ask"); closed["count"] != float64(2) || closed["shown"] != "2" {
		t.Errorf("the closed drawing shows %v, want the state it built before expanding", closed)
	}
}

const whatTheExpandControlSays = `(() => {
  const control = document.querySelector('` + expandControl + `');
  const tip = document.querySelector('.tooltip');
  return {
    tipShown: !tip.hidden && tip.getBoundingClientRect().height > 0,
    tipSaid: (tip.heading + ' ' + tip.hint).trim(),
  };
})()`

type controlSaid struct {
	TipShown bool   `json:"tipShown"`
	TipSaid  string `json:"tipSaid"`
}

func TestTheExpandControlSaysWhatItDoesOnHoverOnFocusAndToAScreenReader(t *testing.T) {
	// scenario: The control says what it does on hover, on focus and to a screen reader
	t.Parallel()
	page := aPageWithACountingDrawing(t)
	page.do(sending(emulation.SetFocusEmulationEnabled, emulation.SetFocusEmulationEnabledParams{Enabled: true}))

	var resting controlSaid
	page.evaluate(whatTheExpandControlSays, &resting)
	if heard := page.buttonsNamed("Expand illustration"); len(heard) != 1 || heard[0].Expanded != "false" {
		t.Errorf("a screen reader hears %+v for the expand control, want one button named Expand illustration that says it is not expanded", heard)
	}
	if resting.TipShown {
		t.Errorf("the control shows %q before anything has its cursor", resting.TipSaid)
	}

	page.evaluate(`document.querySelector('`+expandControl+`').scrollIntoView({ block: 'center' }); new Promise((settled) => requestAnimationFrame(() => requestAnimationFrame(settled)))`, nil, awaitingThePromise)

	var centre struct{ X, Y float64 }
	page.evaluate(`(() => { const box = document.querySelector('`+expandControl+`').getBoundingClientRect(); return { X: box.left + box.width / 2, Y: box.top + box.height / 2 }; })()`, &centre)
	page.do(sending(input.DispatchMouseEvent, input.DispatchMouseEventParams{Type: input.DispatchMouseEventTypeMouseMoved, X: centre.X, Y: centre.Y}))

	var hovered controlSaid
	page.evaluate(whatTheExpandControlSays, &hovered)
	if !hovered.TipShown || !strings.Contains(strings.ToLower(hovered.TipSaid), "expand") {
		t.Errorf("hovering the expand control shows %+v, want it to say what it does", hovered)
	}

	page.do(sending(input.DispatchMouseEvent, input.DispatchMouseEventParams{Type: input.DispatchMouseEventTypeMouseMoved, X: 1, Y: 1}))
	page.focus(expandControl)

	var focused controlSaid
	page.evaluate(whatTheExpandControlSays, &focused)
	if !focused.TipShown || !strings.Contains(strings.ToLower(focused.TipSaid), "expand") {
		t.Errorf("focusing the expand control shows %+v, want it to say what it does", focused)
	}
}

const whereTheFocusIs = `(() => {
  const held = document.activeElement;
  return { behind: held !== document.body && !held.closest('.drawing-stage'), named: held.className || held.tagName };
})()`

type focusAt struct {
	Behind bool   `json:"behind"`
	Named  string `json:"named"`
}

func TestTabFromTheCloseButtonNeverLandsOnAControlBehindTheExpandedDrawing(t *testing.T) {
	// scenario: While a drawing is expanded, nothing behind it can be reached by Tab, and closing makes it reachable again
	t.Parallel()
	page := aPageWithACountingDrawing(t)

	page.click(expandControl)
	for press := 1; press <= 6; press++ {
		page.press(kb.Tab)
		var at focusAt
		page.evaluate(whereTheFocusIs, &at)
		if at.Behind {
			t.Fatalf("tab press %d from the close button put the cursor on %q, a control behind the dimmed page", press, at.Named)
		}
	}

	page.press(kb.Escape)
	page.evaluate(`document.querySelector('.decision[data-decision="c-1"] .choice').focus()`, nil)
	var back focusAt
	page.evaluate(whereTheFocusIs, &back)
	if back.Named == "BODY" || !back.Behind {
		t.Errorf("after closing, the controls behind the drawing cannot take the cursor again: %+v", back)
	}
}

func TestARevisionTakenWhileADrawingIsExpandedLeavesNothingOnThePageUnreachable(t *testing.T) {
	// scenario: A revision taken while a drawing is expanded leaves nothing on the page unreachable
	t.Parallel()
	held := standing(t)
	aDrawing(t, held, "counting.html", drawingWithACounter)
	written := illustrating(t, oneDecisionDocument, "counting.html")
	page := opening(t, held.serve(written))
	page.waitFor(frameOfC1 + `.getBoundingClientRect().height > 0`)

	page.click(expandControl)
	held.serveAgain(strings.Replace(written, "one thing to settle", "one thing, revised", 1))
	page.waitFor(`document.body.dataset.revision === '2'`)

	var unreachable int
	page.evaluate(`document.querySelectorAll('[inert]').length`, &unreachable)
	if unreachable != 0 {
		t.Errorf("%d elements are still unreachable after a revision replaced the expanded drawing", unreachable)
	}
}

func TestAnIllustrationsHeadingNamesNothingAndHoldsItsExpandControlAtTheRight(t *testing.T) {
	// scenario: The bar above a drawing carries no word, only the expand control at its right edge
	t.Parallel()
	page := aPageWithACountingDrawing(t)

	var heading struct {
		Text           string  `json:"text"`
		HeadingRight   float64 `json:"headingRight"`
		ControlRight   float64 `json:"controlRight"`
		HeadingPadding float64 `json:"headingPadding"`
	}
	page.evaluate(`(() => {
  const heading = document.querySelector('.decision[data-decision="c-1"] .evidence.illustrated .evidence-heading');
  const control = heading.querySelector('.expand');
  return {
    text: heading.innerText.trim(),
    headingRight: heading.getBoundingClientRect().right,
    controlRight: control.getBoundingClientRect().right,
    headingPadding: parseFloat(getComputedStyle(heading).paddingRight),
  };
})()`, &heading)

	if heading.Text != "" {
		t.Errorf("the bar above the drawing reads %q, want no word", heading.Text)
	}
	if gap := heading.HeadingRight - heading.HeadingPadding - heading.ControlRight; gap > 1 || gap < -1 {
		t.Errorf("the expand control stands %vpx short of the bar's right edge, want it at the right", gap)
	}
}
