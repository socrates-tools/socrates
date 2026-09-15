package sanitise_test

import (
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/sanitise"
)

const prefix = "drawn-"

func TestAFragmentOutsideTheAllowlistIsStripped(t *testing.T) {
	for _, stripped := range []struct {
		named    string
		supplied string
		gone     []string
		kept     []string
	}{
		{
			named:    "a script",
			supplied: `<p>before</p><script>fetch("https://elsewhere")</script><p>after</p>`,
			gone:     []string{"script", "fetch", "elsewhere"},
			kept:     []string{"before", "after"},
		},
		{
			named:    "an event-handler attribute",
			supplied: `<p onclick="steal()" onmouseover="steal()" class="note">words</p>`,
			gone:     []string{"onclick", "onmouseover", "steal"},
			kept:     []string{"words", `class="drawn-note"`},
		},
		{
			named:    "a remote reference",
			supplied: `<p>before<img src="https://elsewhere/pixel.png"><a href="https://elsewhere">away</a></p>`,
			gone:     []string{"elsewhere", "src", "href", "<img"},
			kept:     []string{"before", "away"},
		},
		{
			named:    "a remote reference from an allowed attribute",
			supplied: `<svg viewBox="0 0 4 4"><rect fill="url(https://elsewhere/g.svg#x)" width="4" height="4"/></svg>`,
			gone:     []string{"elsewhere", "fill"},
			kept:     []string{`viewBox="0 0 4 4"`, `width="4"`},
		},
		{
			named:    "a form control",
			supplied: `<form action="https://elsewhere"><label>name<input name="who"></label><button>send</button></form><p>after</p>`,
			gone:     []string{"<form", "<input", "<button", "elsewhere"},
			kept:     []string{"after"},
		},
		{
			named:    "a style element and a style attribute",
			supplied: `<style>p{background:url(https://elsewhere)}</style><p style="background:url(https://elsewhere)">words</p>`,
			gone:     []string{"style", "elsewhere"},
			kept:     []string{"words"},
		},
	} {
		t.Run(stripped.named, func(t *testing.T) {
			reduced := sanitise.Fragment(stripped.supplied, prefix)

			for _, absent := range stripped.gone {
				if strings.Contains(reduced, absent) {
					t.Errorf("%s survives neutralising as %q:\n%s", stripped.named, absent, reduced)
				}
			}
			for _, present := range stripped.kept {
				if !strings.Contains(reduced, present) {
					t.Errorf("neutralising %s took %q with it:\n%s", stripped.named, present, reduced)
				}
			}
		})
	}
}

func TestAFragmentCannotClaimThePagesOwnNames(t *testing.T) {
	for _, claimed := range []struct {
		named    string
		supplied string
		gone     []string
		kept     []string
	}{
		{
			named:    "a class the page uses for its own chrome",
			supplied: `<div class="decision surface"><span class="commit">go on</span></div>`,
			gone:     []string{`class="decision`, `class="commit"`},
			kept:     []string{`class="drawn-decision drawn-surface"`, `class="drawn-commit"`},
		},
		{
			named:    "an identifier the page addresses",
			supplied: `<div id="decision-f3"><span aria-labelledby="decision-f3">go on</span></div>`,
			gone:     []string{`"decision-f3"`},
			kept:     []string{`id="drawn-decision-f3"`, `aria-labelledby="drawn-decision-f3"`},
		},
	} {
		t.Run(claimed.named, func(t *testing.T) {
			reduced := sanitise.Fragment(claimed.supplied, prefix)

			for _, absent := range claimed.gone {
				if strings.Contains(reduced, absent) {
					t.Errorf("%s reaches the page as %q:\n%s", claimed.named, absent, reduced)
				}
			}
			for _, present := range claimed.kept {
				if !strings.Contains(reduced, present) {
					t.Errorf("renaming %s did not leave %q:\n%s", claimed.named, present, reduced)
				}
			}
		})
	}
}

func TestAnIllustrationSurvivesNeutralising(t *testing.T) {
	drawn := `<figure class="drawn">` +
		`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 120 40" width="120" height="40" role="img" aria-label="two boxes">` +
		`<defs>` +
		`<marker id="arrow" markerWidth="6" markerHeight="6" refX="5" refY="3" orient="auto">` +
		`<path d="M0 0 L6 3 L0 6 z" fill="#8a6d3b"/></marker>` +
		`<linearGradient id="band"><stop offset="0" stop-color="#8a6d3b"/></linearGradient>` +
		`<clipPath id="frame"><rect x="0" y="0" width="120" height="40"/></clipPath>` +
		`<mask id="fade"><rect x="0" y="0" width="120" height="40" fill="#fff"/></mask>` +
		`</defs>` +
		`<rect x="1" y="1" width="40" height="20" rx="3" fill="url(#band)" clip-path="url(#frame)" mask="url(#fade)"/>` +
		`<line x1="41" y1="11" x2="78" y2="11" stroke="#8a6d3b" marker-end="url(#arrow)"/>` +
		`<text x="8" y="34" font-size="11" font-family="ui-monospace" text-anchor="start">read</text>` +
		`</svg><figcaption>the read path</figcaption></figure>`

	reduced := sanitise.Fragment(drawn, prefix)

	for _, renamed := range []string{
		`class="drawn-drawn"`,
		`id="drawn-arrow"`, `marker-end="url(#drawn-arrow)"`,
		`id="drawn-band"`, `fill="url(#drawn-band)"`,
		`id="drawn-frame"`, `clip-path="url(#drawn-frame)"`,
		`id="drawn-fade"`, `mask="url(#drawn-fade)"`,
	} {
		if !strings.Contains(reduced, renamed) {
			t.Errorf("renaming left the drawing without %q, so what it references no longer resolves:\n%s", renamed, reduced)
		}
	}
	for _, drawing := range []string{
		`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 120 40" width="120" height="40" role="img" aria-label="two boxes">`,
		`<path d="M0 0 L6 3 L0 6 z" fill="#8a6d3b"/>`,
		`<text x="8" y="34" font-size="11" font-family="ui-monospace" text-anchor="start">read</text>`,
		`<figcaption>the read path</figcaption>`,
	} {
		if !strings.Contains(reduced, drawing) {
			t.Errorf("an illustration within the allowlist did not survive neutralising with %q:\n%s", drawing, reduced)
		}
	}
	if strings.Count(reduced, prefix) != 9 {
		t.Errorf("the drawing carries %d prefixed names, want the 9 it wrote:\n%s", strings.Count(reduced, prefix), reduced)
	}
}

func TestAnIllustrationTheStylesheetPaintsSurvivesNeutralising(t *testing.T) {
	painted := `<div class="panel">` +
		`<div class="row"><span class="label">Settlement date</span><span class="value">12 March</span></div>` +
		`<div class="row"><span class="label">Amount</span><span class="value" id="total">$1,240.00</span></div>` +
		`<p class="note">Both dates are read from the same field.</p></div>`

	reduced := sanitise.Fragment(painted, prefix)

	want := `<div class="drawn-panel">` +
		`<div class="drawn-row"><span class="drawn-label">Settlement date</span><span class="drawn-value">12 March</span></div>` +
		`<div class="drawn-row"><span class="drawn-label">Amount</span><span class="drawn-value" id="drawn-total">$1,240.00</span></div>` +
		`<p class="drawn-note">Both dates are read from the same field.</p></div>`
	if reduced != want {
		t.Errorf("markup a stylesheet paints did not survive neutralising:\nwant %s\ngot  %s", want, reduced)
	}
}

func TestNothingSuppliedIsNothingEmbedded(t *testing.T) {
	if reduced := sanitise.Fragment("   \n ", prefix); reduced != "" {
		t.Errorf("a fragment of nothing but space reduces to %q", reduced)
	}
}
