package e2e_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-devtool/socrates/internal/document"
)

const callersPalette = `:root{--panel-face:#fbf9f4;--panel-ink:#2f2b23;--panel-rule:#d8d2c4}` +
	`html{color:var(--panel-ink);font-family:ui-monospace,monospace}` +
	`.panel{border:1px solid var(--panel-rule);background:var(--panel-face);padding:0.75rem}` +
	`.row{display:flex;justify-content:space-between;gap:1.5rem;padding:0.35rem 0}` +
	`.row+.row{border-top:1px solid var(--panel-rule)}` +
	`.label{color:#6f6a5e}.value{font-weight:600}#total{color:#8a2f2f}` +
	`.faint{color:#f2f0eb}` +
	`@media (min-width:20rem){.panel{border-radius:6px}}`

const overreachingPalette = `body{background:#ff00ff}main{background:#ff00ff}` +
	`.surface{background:#ff00ff}.decision{background:#ff00ff;display:none}` +
	`.question{color:#ff00ff}.commit{display:none}.send{display:none}` +
	`#decision-c-1{background:#ff00ff}` +
	`.panel{position:fixed;inset:0;width:400vw;height:400vh;background:#00ff00;z-index:2147483647}` +
	`@import url(https://elsewhere/x.css);` +
	`.row{background:url(https://elsewhere/p.png)}`

const paintedPanel = `<div class="panel">` +
	`<div class="row"><span class="label">Settlement date</span><span class="value">12 March</span></div>` +
	`<div class="row"><span class="label">Amount</span><span class="value" id="total">$1,240.00</span></div>` +
	`<p class="faint">Both dates are read from the same field.</p></div>`

const drawnPath = `<figure><svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 120 40" width="120" height="40" ` +
	`role="img" aria-label="the write path"><rect x="1" y="1" width="40" height="20" fill="none" stroke="#8a6d3b"/>` +
	`<text x="6" y="34" font-size="11">write</text></svg><figcaption>the write path</figcaption></figure>`

func illustrated(t *testing.T, stylesheet string) string {
	t.Helper()

	var written map[string]any
	if err := json.Unmarshal([]byte(twoDecisionDocument), &written); err != nil {
		t.Fatalf("reading the document: %v", err)
	}

	decisions := written["decisions"].([]any)
	decisions[0].(map[string]any)["fragment"] = map[string]any{"markup": paintedPanel}
	decisions[1].(map[string]any)["fragment"] = map[string]any{"markup": drawnPath}
	if stylesheet != "" {
		written["stylesheet"] = stylesheet
	}

	raw, err := json.Marshal(written)
	if err != nil {
		t.Fatalf("writing the document: %v", err)
	}
	return string(raw)
}

const whatThePaintingDid = `(() => {
  const panel = document.querySelector('.painted .drawn-panel');
  const total = document.querySelector('.painted #drawn-total');
  const box = document.querySelector('.decision[data-decision="c-1"] .illustration');
  if (!panel || !total || !box) {
    return { painted: false, why: 'the panel the stylesheet paints is not on the page' };
  }
  const held = getComputedStyle(panel);
  const rows = document.querySelectorAll('.painted .drawn-row');
  return {
    painted: true,
    face: held.backgroundColor,
    border: held.borderTopStyle + ' ' + held.borderTopColor,
    radius: held.borderTopLeftRadius,
    mono: held.fontFamily,
    rows: rows.length,
    laidOut: rows.length > 0 ? getComputedStyle(rows[0]).display : '',
    accent: getComputedStyle(total).color,
    within: panel.getBoundingClientRect().left >= box.getBoundingClientRect().left - 1,
  };
})()`

type wasPainted struct {
	Painted bool    `json:"painted"`
	Why     string  `json:"why"`
	Face    string  `json:"face"`
	Border  string  `json:"border"`
	Radius  string  `json:"radius"`
	Mono    string  `json:"mono"`
	Rows    int     `json:"rows"`
	LaidOut string  `json:"laidOut"`
	Accent  string  `json:"accent"`
	Within  bool    `json:"within"`
	Height  float64 `json:"height"`
}

func TestTheEmbeddedStylesheetPaints(t *testing.T) {
	held := standing(t)
	address := held.serve(illustrated(t, callersPalette))
	page := opening(t, address)

	var read wasPainted
	page.evaluate(whatThePaintingDid, &read)

	if !read.Painted {
		t.Fatalf("the illustration was not painted: %s", read.Why)
	}
	if read.Face != "rgb(251, 249, 244)" {
		t.Errorf("the panel is painted %q rather than the face the stylesheet stated", read.Face)
	}
	if read.Border != "solid rgb(216, 210, 196)" {
		t.Errorf("the panel is bordered %q rather than as the stylesheet stated", read.Border)
	}
	if read.Radius == "0px" {
		t.Error("the rule the stylesheet stated inside a media query did not paint")
	}
	if !strings.Contains(read.Mono, "monospace") {
		t.Errorf("the panel is set in %q rather than the family the stylesheet stated", read.Mono)
	}
	if read.Rows != 2 || read.LaidOut != "flex" {
		t.Errorf("the panel laid out %d rows as %q, and the stylesheet stated two rows as flex", read.Rows, read.LaidOut)
	}
	if read.Accent != "rgb(138, 47, 47)" {
		t.Errorf("the name the caller addressed by identifier is painted %q rather than as the stylesheet stated", read.Accent)
	}
	if !read.Within {
		t.Error("the panel the stylesheet painted is not inside the illustration's box")
	}

	if failures := page.violations(); len(failures) != 0 {
		t.Errorf("the browser refused what the page embedded: %v", failures)
	}
	nothingLeftTheMachine(t, address, page.requests.all())
}

const whatTheFrameCarries = `(() => {
  const chrome = [
    'body', 'main.surface', '.heading h1', '.preamble',
    '.decision', '.attribution .name', '.question', '.meaning',
    '.options legend', '.option-label', '.option-meaning', '.words-label',
    '.commit', '.reopen', '.message-label', '.send', '.finish', '.sends-meaning', '.state',
    '.evidence.code pre', '.evidence.illustration', '.hunk .line',
  ];
  const painted = {};
  for (const selector of chrome) {
    const held = document.querySelector(selector);
    if (!held) { painted[selector] = 'nothing'; continue; }
    const style = getComputedStyle(held);
    painted[selector] = [
      style.color, style.backgroundColor, style.borderTopColor, style.borderTopWidth,
      style.display, style.position, style.fontSize, style.fontFamily, style.fontWeight,
      style.padding, style.margin, style.width, style.textAlign, style.zIndex,
    ].join(' | ');
  }
  return painted;
})()`

func frameOf(t *testing.T, stylesheet string) map[string]string {
	t.Helper()

	held := standing(t)
	page := opening(t, held.serve(illustrated(t, stylesheet)))

	var read map[string]string
	page.evaluate(whatTheFrameCarries, &read)

	if len(read) < 20 {
		t.Fatalf("only %d parts of the frame were read, and a check that saw almost nothing would pass forever", len(read))
	}
	for selector, painted := range read {
		if painted == "nothing" {
			t.Fatalf("the page carries no %s, so comparing the frame proves nothing", selector)
		}
	}
	return read
}

func TestACallerCannotPaintTheFrame(t *testing.T) {
	// scenario: A selector naming the page is confined to the illustration
	unpainted := frameOf(t, "")
	overreached := frameOf(t, overreachingPalette)

	if diff := cmp.Diff(unpainted, overreached); diff != "" {
		t.Errorf("a caller's stylesheet moved the page's own frame (-with no stylesheet +with one that names it):\n%s", diff)
	}
}

const whatCoversTheControls = `(() => {
  const controls = [
    '.decision[data-decision="c-1"] .commit',
    '.decision[data-decision="c-1"] .choice',
    '.decision[data-decision="c-1"] .words-field',
    '.send', '.finish', '.message-field',
  ];
  const covered = [];
  for (const selector of controls) {
    const held = document.querySelector(selector);
    if (!held) { covered.push(selector + ' is not on the page'); continue; }
    held.scrollIntoView({ block: 'center' });
    const box = held.getBoundingClientRect();
    if (box.width < 1 || box.height < 1) { covered.push(selector + ' paints nothing'); continue; }
    const at = document.elementFromPoint(box.left + box.width / 2, box.top + box.height / 2);
    if (!at || !(at === held || held.contains(at) || at.contains(held))) {
      covered.push(selector + ' answers to ' + (at ? at.tagName + '.' + at.className : 'nothing'));
    }
  }

  const box = document.querySelector('.decision[data-decision="c-1"] .illustration');
  const style = getComputedStyle(box);
  const page = document.documentElement;
  return {
    covered,
    measured: controls.length,
    contained: style.contain,
    ceiling: parseFloat(style.maxHeight),
    height: box.getBoundingClientRect().height,
    scrolls: box.scrollHeight > box.clientHeight + 1,
    pageScroll: page.scrollWidth - page.clientWidth,
  };
})()`

type covering struct {
	Covered    []string `json:"covered"`
	Measured   int      `json:"measured"`
	Contained  string   `json:"contained"`
	Ceiling    float64  `json:"ceiling"`
	Height     float64  `json:"height"`
	Scrolls    bool     `json:"scrolls"`
	PageScroll float64  `json:"pageScroll"`
}

func TestAPaintedIllustrationCannotCoverThePage(t *testing.T) {
	held := standing(t)
	page := opening(t, held.serve(illustrated(t, overreachingPalette)))

	var read covering
	page.evaluate(whatCoversTheControls, &read)

	if read.Measured < 6 {
		t.Fatalf("only %d controls were looked at, and a check that saw almost nothing would pass forever", read.Measured)
	}
	if read.Contained == "none" || read.Contained == "" {
		t.Errorf("the illustration contains nothing, so a rule inside it decides where its result lands: %+v", read)
	}
	for _, covered := range read.Covered {
		t.Errorf("%s, and every control the page carries stays visible and operable", covered)
	}
	if read.Height > read.Ceiling+1 {
		t.Errorf("the illustration stands %vpx tall against a ceiling of %vpx", read.Height, read.Ceiling)
	}
	if !read.Scrolls {
		t.Errorf("the illustration does not scroll within its own bounds: %+v", read)
	}
	if read.PageScroll > 1 {
		t.Errorf("what the caller painted made the page scroll sideways by %vpx", read.PageScroll)
	}
}

const bothIllustrations = `(() => {
  const held = (decision, inside) =>
    document.querySelector('.decision[data-decision="' + decision + '"] .illustration ' + inside);
  const box = (decision) => document.querySelector('.decision[data-decision="' + decision + '"] .illustration');
  const inside = (decision, inner) => {
    const outer = box(decision);
    if (!outer || !inner) return false;
    const one = outer.getBoundingClientRect();
    const other = inner.getBoundingClientRect();
    return other.width > 0 && other.height > 0 &&
      other.left >= one.left - 1 && other.right <= one.right + 1 && other.top >= one.top - 1;
  };

  const panel = held('c-1', '.painted .drawn-panel');
  const drawn = held('c-2', '.painted svg');
  return {
    paintedBox: box('c-1') ? box('c-1').className : '',
    drawnBox: box('c-2') ? box('c-2').className : '',
    paintedWithin: inside('c-1', panel),
    drawnWithin: inside('c-2', drawn),
  };
})()`

type boxed struct {
	PaintedBox    string `json:"paintedBox"`
	DrawnBox      string `json:"drawnBox"`
	PaintedWithin bool   `json:"paintedWithin"`
	DrawnWithin   bool   `json:"drawnWithin"`
}

func TestAPaintedIllustrationRendersInTheSameBoxAsOneThatPaintsItself(t *testing.T) {
	held := standing(t)
	page := opening(t, held.serve(illustrated(t, callersPalette)))

	var read boxed
	page.evaluate(bothIllustrations, &read)

	if read.PaintedBox == "" || read.DrawnBox == "" {
		t.Fatalf("one of the two illustrations rendered no box at all: %+v", read)
	}
	if read.PaintedBox != read.DrawnBox {
		t.Errorf("a drawing renders in %q and a painted panel in %q, and the difference between them is what the caller supplied", read.DrawnBox, read.PaintedBox)
	}
	if !read.PaintedWithin || !read.DrawnWithin {
		t.Errorf("an illustration renders outside the box the page holds it in: %+v", read)
	}
}

func TestAPaintedIllustrationIsLegibleAtANarrowWidth(t *testing.T) {
	held := standing(t)
	page := opening(t, held.serve(illustrated(t, callersPalette)))

	page.narrow(360, 780)

	var read narrowed
	page.evaluate(narrowly, &read)
	if read.PageScroll > 1 {
		t.Errorf("at 360px the page scrolls sideways by %vpx", read.PageScroll)
	}
	if !read.DrawnVisible || !read.DrawnWithin {
		t.Errorf("at 360px the drawing is not held within the page and legible: %+v", read)
	}

	var painted wasPainted
	page.evaluate(narrowPainting, &painted)
	if !painted.Painted {
		t.Fatalf("at 360px the panel the stylesheet paints is not on the page: %s", painted.Why)
	}
	if painted.Height < 20 {
		t.Errorf("at 360px the panel the stylesheet paints stands %vpx tall, so it is not legible", painted.Height)
	}
	if !painted.Within {
		t.Errorf("at 360px the panel is not held within the illustration, and the illustration scrolls within its own bounds: %+v", painted)
	}
}

const narrowPainting = `(() => {
  const box = document.querySelector('.decision[data-decision="c-1"] .illustration');
  const panel = document.querySelector('.painted .drawn-panel');
  if (!box || !panel) return { painted: false, why: 'the panel is not on the page' };
  const style = getComputedStyle(box);
  return {
    painted: true,
    height: panel.getBoundingClientRect().height,
    within: box.clientWidth <= document.querySelector('main.surface').clientWidth + 1 &&
      (style.overflowX === 'auto' || style.overflowX === 'scroll'),
  };
})()`

const measuringInsideTheIllustration = `(() => {
  const parse = (colour) => {
    const held = colour.match(/[\d.]+/g);
    if (!held) return null;
    const [r, g, b, a] = held.map(Number);
    return { r, g, b, a: a === undefined ? 1 : a };
  };
  const luminance = ({ r, g, b }) => {
    const channel = (held) => {
      const scaled = held / 255;
      return scaled <= 0.03928 ? scaled / 12.92 : Math.pow((scaled + 0.055) / 1.055, 2.4);
    };
    return 0.2126 * channel(r) + 0.7152 * channel(g) + 0.0722 * channel(b);
  };
  const ratio = (one, other) => {
    const first = luminance(one);
    const second = luminance(other);
    return (Math.max(first, second) + 0.05) / (Math.min(first, second) + 0.05);
  };
  const behind = (element) => {
    for (let held = element; held; held = held.parentElement) {
      const colour = parse(getComputedStyle(held).backgroundColor);
      if (colour && colour.a === 1) return colour;
    }
    return { r: 255, g: 255, b: 255, a: 1 };
  };

  const failing = [];
  let measured = 0;
  for (const element of document.querySelectorAll('.illustration *')) {
    if (element.namespaceURI !== 'http://www.w3.org/1999/xhtml') continue;
    const owns = Array.from(element.childNodes).some((node) => node.nodeType === 3 && node.textContent.trim());
    if (!owns) continue;
    const painted = parse(getComputedStyle(element).color);
    if (!painted || painted.a < 1) continue;
    measured += 1;
    const held = ratio(painted, behind(element));
    if (held < 4.5) failing.push(element.tagName + '.' + String(element.className) + ' at ' + held.toFixed(2));
  }
  return { failing, measured };
})()`

func TestACallersColoursDoNotReachTheDecision(t *testing.T) {
	held := standing(t)
	page := opening(t, held.serve(illustrated(t, callersPalette)))

	page.evaluate(chose("c-1", 0), nil)
	page.click(within("c-1", ".commit"))

	for _, preference := range []string{"light", "dark"} {
		page.emulate(preference)

		var read weighed
		page.evaluate(measuringTheContrast, &read)

		if read.Measured < 30 {
			t.Fatalf("under a %s preference only %d texts the page paints were measured, and a check that saw almost nothing would pass forever", preference, read.Measured)
		}
		for _, failing := range read.Failing {
			t.Errorf("under a %s preference %s, and the floor holds every text the page itself paints", preference, failing)
		}
	}

	var inside weighed
	page.evaluate(measuringInsideTheIllustration, &inside)
	if inside.Measured == 0 {
		t.Fatal("no text inside the illustration was measured, so nothing shows that the sweep steps around one")
	}
	if len(inside.Failing) == 0 {
		t.Error("the caller painted nothing the floor would reject, so this says nothing about what the floor reaches")
	}
}

func TestAnOversizedStylesheetIsRefusedBeforeServing(t *testing.T) {
	held := standing(t)

	oversized := illustrated(t, strings.Repeat(callersPalette, 1+document.StylesheetCeiling/len(callersPalette)))
	printed, failures, code := held.reading(oversized, "serve", "--document", "-")

	if code == 0 {
		t.Fatalf("a document carrying an oversized stylesheet was served at %q", strings.TrimSpace(printed))
	}
	for _, named := range []string{"stylesheet", strconv.Itoa(document.StylesheetCeiling)} {
		if !strings.Contains(failures, named) {
			t.Errorf("refusing an oversized stylesheet does not say %q was exceeded: %s", named, failures)
		}
	}
	if _, err := os.Stat(filepath.Join(held.surfaceHome(), "surface.json")); err == nil {
		t.Error("a refused document left a session behind, and it starts no server")
	}
}
