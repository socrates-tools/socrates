package e2e_test

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/chromedp/chromedp/kb"
)

func aDrawingOfHeight(rows int) string {
	return fmt.Sprintf(`<!doctype html>
<html><head><style>body { margin: 0; } .row { height: 100px; }</style></head>
<body>
<div id="rows">%s</div>
<script>
window.addEventListener('message', (event) => {
  if (event.data === 'grow') {
    const row = document.createElement('div');
    row.className = 'row';
    document.getElementById('rows').append(row);
  }
  if (event.data === 'lie') {
    for (const lie of ['640', {height: 640}, [640], null, true, NaN, Infinity, -Infinity, -1, 'tall']) parent.postMessage(lie, '*');
  }
  parent.postMessage({ type: 'report', inner: window.innerHeight, scroll: document.documentElement.scrollHeight }, '*');
});
</script>
</body></html>`, strings.Repeat(`<div class="row"></div>`, rows))
}

const framesOfBothCards = `[...document.querySelectorAll('iframe.drawing')]`

func (r *reader) heightOfTheFrame() float64 {
	r.t.Helper()

	var height float64
	r.evaluate(frameOfC1+`.getBoundingClientRect().height`, &height)
	return height
}

func (r *reader) waitForTheFrameToBe(height float64) {
	r.t.Helper()
	r.waitFor(fmt.Sprintf(`Math.abs(%s.getBoundingClientRect().height - %v) < 1`, frameOfC1, height))
}

func (r *reader) gapUnderTheFrame() float64 {
	r.t.Helper()

	var gap float64
	r.evaluate(`(() => {
  const frame = `+frameOfC1+`;
  const card = frame.closest('.evidence');
  const borderBelow = parseFloat(getComputedStyle(card).borderBottomWidth);
  return card.getBoundingClientRect().bottom - borderBelow - frame.getBoundingClientRect().bottom;
})()`, &gap)
	return gap
}

func (r *reader) askTheDrawingHowItStands(message string) map[string]any {
	r.t.Helper()

	r.evaluate(`window.__stood = null; if (!window.__watching) {
  window.__watching = true;
  window.addEventListener('message', (event) => {
    if (event.source === `+frameOfC1+`.contentWindow && event.data && event.data.type === 'report') window.__stood = event.data;
  });
} true`, nil)
	r.evaluate(frameOfC1+`.contentWindow.postMessage('`+message+`', '*'); true`, nil)
	r.waitFor(`window.__stood !== null`)

	var stood map[string]any
	r.evaluate(`window.__stood`, &stood)
	return stood
}

func aPageDrawing(t *testing.T, rows int) *reader {
	t.Helper()

	held := standing(t)
	aDrawing(t, held, "rows.html", aDrawingOfHeight(rows))
	return opening(t, held.serve(illustrating(t, oneDecisionDocument, "rows.html")))
}

func TestAShortDrawingsFrameIsTheDrawingsHeightWithNoGapUnderIt(t *testing.T) {
	// scenario: A short drawing's frame is the drawing's height, with no gap under it
	t.Parallel()
	page := aPageDrawing(t, 2)
	page.narrow(1280, 900)

	page.waitForTheFrameToBe(200)

	if gap := page.gapUnderTheFrame(); math.Abs(gap) > 1 {
		t.Errorf("the card leaves %vpx under the frame, want none", gap)
	}
}

func TestADrawingTallerThanTheWindowStopsAtTheWindowsHeightAndScrollsInsideItsFrame(t *testing.T) {
	// scenario: A drawing taller than the window stops at the window's height and scrolls inside its frame
	t.Parallel()
	page := aPageDrawing(t, 30)
	page.narrow(1280, 600)

	page.waitForTheFrameToBe(600)

	stood := page.askTheDrawingHowItStands("report")
	if stood["scroll"].(float64) <= stood["inner"].(float64) {
		t.Errorf("the drawing is %v tall in a frame %v tall, so nothing scrolls inside it", stood["scroll"], stood["inner"])
	}

	page.narrow(1280, 450)
	page.waitForTheFrameToBe(450)
}

func TestADrawingThatGrowsAfterLoadingGrowsItsFrame(t *testing.T) {
	// scenario: A drawing that grows after loading, such as by a script adding rows, grows its frame
	t.Parallel()
	page := aPageDrawing(t, 2)
	page.narrow(1280, 900)
	page.waitForTheFrameToBe(200)

	page.askTheDrawingHowItStands("grow")
	page.waitForTheFrameToBe(300)

	page.askTheDrawingHowItStands("grow")
	page.waitForTheFrameToBe(400)
}

func TestAFramePostingAnythingButANumberChangesNothingOnThePage(t *testing.T) {
	// scenario: A frame posting anything but a number changes nothing on the page
	t.Parallel()
	page := aPageDrawing(t, 2)
	page.narrow(1280, 900)
	page.waitForTheFrameToBe(200)

	page.askTheDrawingHowItStands("lie")

	if height := page.heightOfTheFrame(); math.Abs(height-200) > 1 {
		t.Errorf("the frame is %vpx after the drawing posted things that are not a height, want 200", height)
	}
}

func TestAMessageFromAnythingButTheCardsOwnFrameChangesNothingOnThePage(t *testing.T) {
	// scenario: A frame posting for another card changes nothing on the page
	t.Parallel()
	held := standing(t)
	aDrawing(t, held, "short.html", aDrawingOfHeight(2))
	aDrawing(t, held, "tall.html", aDrawingOfHeight(4))
	page := opening(t, held.serve(illustratingBoth(t, twoDecisionDocument, "short.html", "tall.html")))
	page.narrow(1280, 900)

	page.waitFor(`(() => { const [first, second] = ` + framesOfBothCards + `; return first && second && Math.abs(first.getBoundingClientRect().height - 200) < 1 && Math.abs(second.getBoundingClientRect().height - 400) < 1; })()`)

	page.evaluate(`window.postMessage(640, '*'); window.postMessage(0, '*'); true`, nil)
	page.evaluate(`new Promise((settled) => setTimeout(settled, 300))`, nil, awaitingThePromise)

	var heights []float64
	page.evaluate(framesOfBothCards+`.map((frame) => frame.getBoundingClientRect().height)`, &heights)
	if len(heights) != 2 || math.Abs(heights[0]-200) > 1 || math.Abs(heights[1]-400) > 1 {
		t.Errorf("the frames are %v after the page was posted numbers by something that is not their own frame, want [200 400]", heights)
	}
}

func TestADrawingThatGrowsWhileExpandedIsTheGrownHeightOnceClosed(t *testing.T) {
	// scenario: The height the page sets does not fight the expanded drawing and is right again on closing
	t.Parallel()
	page := aPageDrawing(t, 2)
	page.narrow(1280, 900)
	page.waitForTheFrameToBe(200)

	page.click(expandControl)
	expanded := page.whereTheFrameStands()
	if !expanded.fillsTheWindow() {
		t.Fatalf("the drawing did not expand: %+v", expanded)
	}

	page.askTheDrawingHowItStands("grow")
	page.evaluate(`new Promise((settled) => setTimeout(settled, 300))`, nil, awaitingThePromise)
	if still := page.whereTheFrameStands(); !still.fillsTheWindow() {
		t.Errorf("growing while expanded took the drawing off the window: %+v", still)
	}

	page.press(kb.Escape)
	page.waitForTheFrameToBe(300)
	if gap := page.gapUnderTheFrame(); math.Abs(gap) > 1 {
		t.Errorf("the card leaves %vpx under the frame after closing, want none", gap)
	}
}

func illustratingBoth(t *testing.T, base, first, second string) string {
	t.Helper()

	var written map[string]any
	if err := json.Unmarshal([]byte(illustrating(t, base, first)), &written); err != nil {
		t.Fatalf("reading the document: %v", err)
	}
	written["decisions"].([]any)[1].(map[string]any)["illustration"] = map[string]any{"path": second}

	raw, err := json.Marshal(written)
	if err != nil {
		t.Fatalf("writing the document: %v", err)
	}
	return string(raw)
}
