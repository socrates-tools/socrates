package e2e_test

import (
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/google/go-cmp/cmp"
)

const measuringTheText = `(() => {
  const floor = parseFloat(getComputedStyle(document.documentElement).getPropertyValue('--text-floor'));
  const below = [];
  for (const held of document.querySelectorAll('body *')) {
    if (held.namespaceURI !== 'http://www.w3.org/1999/xhtml') continue;
    if (held.closest('.illustration')) continue;
    const owns = Array.from(held.childNodes).some((node) => node.nodeType === 3 && node.textContent.trim());
    if (!owns) continue;
    const style = getComputedStyle(held);
    if (style.display === 'none' || style.visibility === 'hidden') continue;
    const size = parseFloat(style.fontSize);
    if (size < floor) below.push(held.tagName + '.' + held.className + ' at ' + size);
  }
  return { floor, below, measured: document.querySelectorAll('body *').length };
})()`

type measured struct {
	Floor    float64  `json:"floor"`
	Below    []string `json:"below"`
	Measured int      `json:"measured"`
}

func TestNothingRendersBelowTheFloor(t *testing.T) {
	t.Parallel()
	_, page := read(t, twoDecisionDocument)

	var read measured
	page.evaluate(measuringTheText, &read)

	if read.Floor < 12 {
		t.Fatalf("the floor the page states is %vpx, and a floor that low is not a floor", read.Floor)
	}
	if read.Measured < 20 {
		t.Fatalf("only %d elements were measured, and a check that saw almost nothing would pass forever", read.Measured)
	}
	for _, small := range read.Below {
		t.Errorf("%s renders below the floor of %vpx the page states", small, read.Floor)
	}
}

const painting = `(() => {
  const root = getComputedStyle(document.documentElement);
  const body = getComputedStyle(document.body);
  const tokens = {};
  for (const named of Array.from(root)) {
    if (named.startsWith('--')) tokens[named] = root.getPropertyValue(named).trim();
  }
  return { background: body.backgroundColor, ink: body.color, scheme: root.colorScheme, tokens };
})()`

type painted struct {
	Background string            `json:"background"`
	Ink        string            `json:"ink"`
	Scheme     string            `json:"scheme"`
	Tokens     map[string]string `json:"tokens"`
}

func TestThePageIsPaintedLightWhateverThePreference(t *testing.T) {
	t.Parallel()
	_, page := read(t, twoDecisionDocument)

	under := map[string]painted{}
	for _, preference := range []string{"light", "dark"} {
		page.emulate(preference)

		var read painted
		page.evaluate(painting, &read)
		under[preference] = read

		if read.Background == "" || strings.Contains(read.Background, "rgba(0, 0, 0, 0)") {
			t.Errorf("under a %s preference the page paints no background of its own: %q", preference, read.Background)
		}
		if read.Ink == "" {
			t.Errorf("under a %s preference the page paints no text colour of its own", preference)
		}
		if len(read.Tokens) < 10 {
			t.Errorf("under a %s preference the page declares %d colours, and a check that saw almost none would pass forever", preference, len(read.Tokens))
		}
		for named, value := range read.Tokens {
			if value == "" {
				t.Errorf("under a %s preference %s is left undefined", preference, named)
			}
		}
		if read.Scheme != "light" {
			t.Errorf("under a %s preference the page states the colour scheme %q, and it is light whatever the reader prefers", preference, read.Scheme)
		}
	}

	if diff := cmp.Diff(under["light"], under["dark"]); diff != "" {
		t.Errorf("a dark preference paints the page differently from a light one, and the page has one palette (-light +dark):\n%s", diff)
	}
}

const focused = `(() => {
  const held = document.activeElement;
  if (!held || held === document.body) return { tag: '', name: '', className: '' };
  const label = held.closest('label');
  const name = held.getAttribute('aria-label') || (label ? label.textContent.trim() : '') || held.textContent.trim();
  return { tag: held.tagName, name, className: String(held.className) };
})()`

type reached struct {
	Tag       string `json:"tag"`
	Name      string `json:"name"`
	ClassName string `json:"className"`
}

func TestADecisionCanBeAnsweredWithoutAPointer(t *testing.T) {
	t.Parallel()
	_, page := read(t, twoDecisionDocument)

	seen := map[string]bool{}
	for range 40 {
		page.do(chromedp.KeyEvent(kb.Tab))

		var held reached
		page.evaluate(focused, &held)
		if held.Tag == "" {
			continue
		}
		for _, named := range strings.Fields(held.ClassName) {
			seen[named] = true
			if strings.TrimSpace(held.Name) == "" {
				t.Errorf("the %s reached by keyboard carries no name a reader could hear", named)
			}
		}
	}

	for _, wanted := range []string{"choice", "words-field", "message-field", "send", "finish"} {
		if !seen[wanted] {
			t.Errorf("the %s cannot be reached by keyboard alone", wanted)
		}
	}

	page.do(chromedp.Focus(chromedp.CSS(within("c-1", ".choice"))), chromedp.KeyEvent(" "))

	var ticked bool
	page.evaluate(`document.querySelector('.decision[data-decision="c-1"] .choice').checked`, &ticked)
	if !ticked {
		t.Error("an option focused by keyboard cannot be chosen by keyboard")
	}
}

const narrowly = `(() => {
  const page = document.documentElement;
  const hunk = document.querySelector('.hunk');
  const drawn = document.querySelector('.illustration svg');
  const surface = document.querySelector('main.surface');
  return {
    pageScroll: page.scrollWidth - page.clientWidth,
    hunkOverflows: hunk ? getComputedStyle(hunk).overflowX : '',
    hunkWithin: hunk ? hunk.clientWidth <= surface.clientWidth : false,
    hunkHasContent: hunk ? hunk.scrollWidth > 0 : false,
    drawnWithin: drawn ? drawn.getBoundingClientRect().width <= surface.clientWidth + 1 : false,
    drawnVisible: drawn ? drawn.getBoundingClientRect().width > 0 : false,
  };
})()`

type narrowed struct {
	PageScroll     float64 `json:"pageScroll"`
	HunkOverflows  string  `json:"hunkOverflows"`
	HunkWithin     bool    `json:"hunkWithin"`
	HunkHasContent bool    `json:"hunkHasContent"`
	DrawnWithin    bool    `json:"drawnWithin"`
	DrawnVisible   bool    `json:"drawnVisible"`
}

func TestTheNarrowCaseDoesNotShrinkTheEvidenceToNothing(t *testing.T) {
	t.Parallel()
	_, page := read(t, twoDecisionDocument)

	page.narrow(360, 780)

	var read narrowed
	page.evaluate(narrowly, &read)

	if read.PageScroll > 1 {
		t.Errorf("at 360px the page scrolls sideways by %vpx", read.PageScroll)
	}
	if !slices.Contains([]string{"auto", "scroll"}, read.HunkOverflows) {
		t.Errorf("at 360px the diff overflows as %q rather than scrolling within its own bounds", read.HunkOverflows)
	}
	if !read.HunkWithin || !read.HunkHasContent {
		t.Errorf("at 360px the diff is not held within the page and legible: %+v", read)
	}
	if !read.DrawnVisible || !read.DrawnWithin {
		t.Errorf("at 360px the illustration is not held within the page and legible: %+v", read)
	}
}

func TestThePageReachesNothing(t *testing.T) {
	t.Parallel()
	held := standing(t)
	address := held.serve(twoDecisionDocument)
	page := opening(t, address)

	page.waitFor(`document.fonts.status === 'loaded'`)
	page.waitFor(`!!document.querySelector('code.hljs')`)

	nothingLeftTheMachine(t, address, page.requests.all())
}

func TestRenderingNeedsNoNetwork(t *testing.T) {
	t.Parallel()
	held := standing(t)
	address := held.serve(twoDecisionDocument)

	page := opening(t, address, chromedp.Flag("host-resolver-rules", "MAP * ~NOTFOUND, EXCLUDE 127.0.0.1"))

	page.waitFor(`document.fonts.status === 'loaded'`)
	page.waitFor(`!!document.querySelector('code.hljs')`)

	var read struct {
		Face      bool `json:"face"`
		Display   bool `json:"display"`
		Mono      bool `json:"mono"`
		Diff      bool `json:"diff"`
		Drawn     bool `json:"drawn"`
		Highlight int  `json:"highlight"`
	}
	page.evaluate(`({
  face: document.fonts.check('16px "Manrope"'),
  display: document.fonts.check('40px "Cinzel"'),
  mono: document.fonts.check('16px "JetBrains Mono"'),
  diff: !!document.querySelector('.hunk .line[data-kind="added"]'),
  drawn: !!document.querySelector('.illustration svg rect'),
  highlight: document.querySelectorAll('code.hljs .hljs-keyword, code.hljs .hljs-string, code.hljs .hljs-title').length,
})`, &read)

	if !read.Face || !read.Display || !read.Mono {
		t.Errorf("with no name resolution the fonts did not render: %+v", read)
	}
	if !read.Diff {
		t.Error("with no name resolution the diff did not render")
	}
	if !read.Drawn {
		t.Error("with no name resolution the illustration did not render")
	}
	if read.Highlight == 0 {
		t.Error("with no name resolution the highlighting did not render")
	}

	nothingLeftTheMachine(t, address, page.requests.all())
}

func TestTheTabAndTheTopBarCarryTheSocratesMark(t *testing.T) {
	t.Parallel()
	held := standing(t)
	address := held.serve(twoDecisionDocument)
	page := opening(t, address)

	page.waitFor(`document.querySelector('.brand-mark') && document.querySelector('.brand-mark').complete`)

	var read struct {
		Icon     string `json:"icon"`
		Mark     string `json:"mark"`
		Drawn    int    `json:"drawn"`
		Shown    int    `json:"shown"`
		Silent   bool   `json:"silent"`
		Fetched  int    `json:"fetched"`
		Answered string `json:"answered"`
	}
	page.evaluate(`(async () => {
  const icon = document.querySelector('link[rel="icon"]');
  const mark = document.querySelector('.topbar .brand-mark');
  const fetched = icon ? await fetch(icon.href) : null;
  return {
    icon: icon ? new URL(icon.href).pathname : '',
    mark: mark ? new URL(mark.src).pathname : '',
    drawn: mark ? mark.naturalWidth : 0,
    shown: mark ? Math.round(mark.getBoundingClientRect().width) : 0,
    silent: mark ? mark.alt === '' : false,
    fetched: fetched ? fetched.status : 0,
    answered: fetched ? fetched.headers.get('content-type') : '',
  };
})()`, &read, awaitingThePromise)

	if read.Icon == "" || read.Icon != read.Mark {
		t.Errorf("the tab's icon is %q and the top bar's mark is %q, want both to be the one socrates mark", read.Icon, read.Mark)
	}
	if read.Fetched != 200 || !strings.HasPrefix(read.Answered, "image/svg+xml") {
		t.Errorf("the tab's icon answered %d as %q, want 200 as image/svg+xml", read.Fetched, read.Answered)
	}
	if read.Drawn == 0 || read.Shown == 0 {
		t.Errorf("the top bar's mark did not draw: %+v", read)
	}
	if !read.Silent {
		t.Error("the top bar's mark is announced beside the wordmark that already says socrates")
	}
}

const whatIsHidden = `(() => {
  const painted = [];
  const held = document.querySelectorAll('[hidden], dialog:not([open])');
  for (const element of held) {
    const style = getComputedStyle(element);
    const box = element.getBoundingClientRect();
    if (style.display !== 'none' || box.width > 0 || box.height > 0) {
      painted.push(element.tagName + '.' + element.className + ' displays as ' + style.display + ' at ' + Math.round(box.width) + 'x' + Math.round(box.height));
    }
  }
  return { painted, hidden: held.length };
})()`

type hiding struct {
	Painted []string `json:"painted"`
	Hidden  int      `json:"hidden"`
}

func TestNothingHiddenIsPainted(t *testing.T) {
	t.Parallel()
	_, page := read(t, twoDecisionDocument)

	var read hiding
	page.evaluate(whatIsHidden, &read)

	if read.Hidden == 0 {
		t.Fatal("the page hides nothing at all, and a check that saw nothing would pass forever")
	}
	for _, painted := range read.Painted {
		t.Errorf("%s, and a person clicks what they can see", painted)
	}
}

const readingTheMark = `(() => {
  const option = document.querySelector('.decision[data-decision="c-2"] .option[data-recommended="yes"]');
  if (!option) return { found: false };
  const mark = option.querySelector('.option-mark');
  const style = mark ? getComputedStyle(mark) : null;
  return {
    found: true,
    name: option.textContent.replace(/\s+/g, ' ').trim(),
    said: mark ? mark.textContent.trim() : '',
    size: style ? parseFloat(style.fontSize) : 0,
    shown: style ? style.display !== 'none' && style.visibility !== 'hidden' : false,
    first: option === document.querySelector('.decision[data-decision="c-2"] .option'),
    floor: parseFloat(getComputedStyle(document.documentElement).getPropertyValue('--text-floor')),
  };
})()`

type marked struct {
	Found bool    `json:"found"`
	Name  string  `json:"name"`
	Said  string  `json:"said"`
	Size  float64 `json:"size"`
	Shown bool    `json:"shown"`
	First bool    `json:"first"`
	Floor float64 `json:"floor"`
}

func TestTheMarkIsSaidRatherThanPainted(t *testing.T) {
	t.Parallel()
	_, page := read(t, twoDecisionDocument)

	page.evaluate(`(() => {
  const stripped = document.createElement('style');
  stripped.textContent = '* { color: black !important; background: white !important; border-color: black !important; }';
  document.head.appendChild(stripped);
})()`, nil)

	var read marked
	page.evaluate(readingTheMark, &read)

	if !read.Found {
		t.Fatal("the page marks no option, and the document marks one")
	}
	if !read.Shown || read.Said == "" {
		t.Errorf("with every colour discarded the mark says %q and is shown: %v", read.Said, read.Shown)
	}
	if !strings.Contains(strings.ToLower(read.Name), strings.ToLower(read.Said)) {
		t.Errorf("the option is announced as %q, and the mark %q is not part of it", read.Name, read.Said)
	}
	if !strings.Contains(strings.ToLower(read.Said), "recommend") {
		t.Errorf("the mark says %q, and it has to say that this is the option the caller would take", read.Said)
	}
	if !read.First {
		t.Error("the marked option is not the first a person meets")
	}
	if read.Size < read.Floor {
		t.Errorf("the mark renders at %vpx, below the floor of %vpx the page states", read.Size, read.Floor)
	}
}

const readingTheMarkedMeaning = `(() => {
  const option = document.querySelector('.decision[data-decision="c-2"] .option[data-recommended="yes"]');
  const mark = option ? option.querySelector('.option-mark') : null;
  const meaning = option ? option.querySelector('.option-meaning') : null;
  if (!mark || !meaning) return { found: false };
  return {
    found: true,
    meaningRight: meaning.getBoundingClientRect().right,
    markRight: mark.getBoundingClientRect().right,
  };
})()`

type markedMeaning struct {
	Found        bool    `json:"found"`
	MeaningRight float64 `json:"meaningRight"`
	MarkRight    float64 `json:"markRight"`
}

func TestAtANarrowWidthTheMeaningOfAMarkedOptionRunsUnderTheMark(t *testing.T) {
	t.Parallel()
	_, page := read(t, twoDecisionDocument)

	page.narrow(360, 780)

	var read markedMeaning
	page.evaluate(readingTheMarkedMeaning, &read)

	if !read.Found {
		t.Fatal("the page marks no option with a meaning, and the document marks one")
	}
	if read.MeaningRight < read.MarkRight-1 {
		t.Errorf("at 360px the meaning of the marked option ends at %vpx and the mark at %vpx, leaving the width under the mark empty", read.MeaningRight, read.MarkRight)
	}
}

const contrastFloor = 4.5

var measuringTheContrast = `(() => {
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
    const root = parse(getComputedStyle(document.documentElement).backgroundColor);
    return root && root.a === 1 ? root : { r: 255, g: 255, b: 255, a: 1 };
  };
  const named = (element) => element.tagName + '.' + String(element.className);

  const failing = [];
  let measured = 0;
  const weigh = (element, colour, what) => {
    const painted = parse(colour);
    if (!painted || painted.a < 1) return;
    measured += 1;
    const held = ratio(painted, behind(element));
    if (held < ` + strconv.FormatFloat(contrastFloor, 'f', -1, 64) + `) {
      failing.push(named(element) + what + ' at ' + held.toFixed(2));
    }
  };

  for (const element of document.querySelectorAll('body *')) {
    if (element.namespaceURI !== 'http://www.w3.org/1999/xhtml') continue;
    if (element.closest('.illustration')) continue;
    const style = getComputedStyle(element);
    if (style.display === 'none' || style.visibility === 'hidden') continue;

    const owns = Array.from(element.childNodes).some((node) => node.nodeType === 3 && node.textContent.trim());
    if (owns) weigh(element, style.color, '');
    if (element.matches('textarea, input:not([type="radio"]):not([type="checkbox"])')) {
      weigh(element, getComputedStyle(element, '::placeholder').color, ' placeholder');
    }
  }
  return { failing, measured };
})()`

type weighed struct {
	Failing  []string `json:"failing"`
	Measured int      `json:"measured"`
}

func TestEveryTextMeetsTheContrastFloor(t *testing.T) {
	t.Parallel()
	_, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)

	for _, preference := range []string{"light", "dark"} {
		page.emulate(preference)

		var read weighed
		page.evaluate(measuringTheContrast, &read)

		if read.Measured < 30 {
			t.Fatalf("under a %s preference only %d texts were measured, and a check that saw almost nothing would pass forever", preference, read.Measured)
		}
		for _, failing := range read.Failing {
			t.Errorf("under a %s preference %s, and the floor the page is held to is %v", preference, failing, contrastFloor)
		}
	}
}

func TestTheReceiptMeetsTheContrastFloor(t *testing.T) {
	t.Parallel()
	_, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-1", 0), nil)
	page.write(".message-field", "one more look please")
	page.click(".send")
	page.waitFor(`!document.querySelector('.receipt').hidden`)

	for _, preference := range []string{"light", "dark"} {
		page.emulate(preference)

		var read weighed
		page.evaluate(measuringTheContrast, &read)

		if read.Measured < 30 {
			t.Fatalf("under a %s preference only %d texts were measured with the receipt showing, and a check that saw almost nothing would pass forever", preference, read.Measured)
		}
		for _, failing := range read.Failing {
			t.Errorf("under a %s preference %s, and the floor the page is held to is %v", preference, failing, contrastFloor)
		}
	}
}

func TestTheConfirmationMeetsTheContrastFloor(t *testing.T) {
	t.Parallel()
	_, page := read(t, twoDecisionDocument)

	page.evaluate(chose("c-2", 2), nil)
	page.click(".send")
	page.waitFor(confirmationIsOpen)

	for _, preference := range []string{"light", "dark"} {
		page.emulate(preference)

		var read weighed
		page.evaluate(measuringTheContrast, &read)

		if read.Measured < 30 {
			t.Fatalf("under a %s preference only %d texts were measured with the confirmation open, and a check that saw almost nothing would pass forever", preference, read.Measured)
		}
		for _, failing := range read.Failing {
			t.Errorf("under a %s preference %s, and the floor the page is held to is %v", preference, failing, contrastFloor)
		}
	}
}
