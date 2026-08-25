package e2e_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
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
	_, page := read(t, twoDecisionDocument)

	var read measured
	page.evaluate(measuringTheText, &read)

	if read.Floor < 14 {
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
  return { background: body.backgroundColor, ink: body.color, tokens };
})()`

type painted struct {
	Background string            `json:"background"`
	Ink        string            `json:"ink"`
	Tokens     map[string]string `json:"tokens"`
}

func TestBothPreferencesArePainted(t *testing.T) {
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
	}

	if under["light"].Background == under["dark"].Background {
		t.Errorf("both preferences paint the same background %q, so the page does not follow the reader", under["light"].Background)
	}
	if under["light"].Ink == under["dark"].Ink {
		t.Errorf("both preferences paint the same text colour %q", under["light"].Ink)
	}
	for named, value := range under["light"].Tokens {
		if _, held := under["dark"].Tokens[named]; !held {
			t.Errorf("%s is defined under a light preference and not under a dark one: %q", named, value)
		}
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

	for _, wanted := range []string{"choice", "words-field", "commit", "message-field", "send", "finish"} {
		if !seen[wanted] {
			t.Errorf("the %s cannot be reached by keyboard alone", wanted)
		}
	}

	page.do(chromedp.Focus(within("c-1", ".choice"), chromedp.ByQuery), chromedp.KeyEvent(" "))

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
	held := standing(t)
	address := held.serve(twoDecisionDocument)
	page := opening(t, address)

	page.waitFor(`document.fonts.status === 'loaded'`)
	page.waitFor(`!!document.querySelector('code.hljs')`)

	nothingLeftTheMachine(t, address, page.requests.all())
}

func TestRenderingNeedsNoNetwork(t *testing.T) {
	held := standing(t)
	address := held.serve(twoDecisionDocument)

	page := opening(t, address, chromedp.Flag("host-resolver-rules", "MAP * ~NOTFOUND, EXCLUDE 127.0.0.1"))

	page.waitFor(`document.fonts.status === 'loaded'`)
	page.waitFor(`!!document.querySelector('code.hljs')`)

	var read struct {
		Face      bool `json:"face"`
		Mono      bool `json:"mono"`
		Diff      bool `json:"diff"`
		Drawn     bool `json:"drawn"`
		Highlight int  `json:"highlight"`
	}
	page.evaluate(`({
  face: document.fonts.check('16px "IBM Plex Sans"'),
  mono: document.fonts.check('16px "IBM Plex Mono"'),
  diff: !!document.querySelector('.hunk .line[data-kind="added"]'),
  drawn: !!document.querySelector('.illustration svg rect'),
  highlight: document.querySelectorAll('code.hljs .hljs-keyword, code.hljs .hljs-string, code.hljs .hljs-title').length,
})`, &read)

	if !read.Face || !read.Mono {
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
