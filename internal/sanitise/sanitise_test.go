package sanitise_test

import (
	"strings"
	"testing"

	"github.com/rui-ktei/socrates/internal/sanitise"
)

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
			kept:     []string{"words", `class="note"`},
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
			reduced := sanitise.Fragment(stripped.supplied)

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

func TestAnIllustrationSurvivesNeutralising(t *testing.T) {
	drawn := `<figure class="drawn">` +
		`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 120 40" width="120" height="40" role="img" aria-label="two boxes">` +
		`<defs><marker id="arrow" markerWidth="6" markerHeight="6" refX="5" refY="3" orient="auto">` +
		`<path d="M0 0 L6 3 L0 6 z" fill="#8a6d3b"/></marker></defs>` +
		`<rect x="1" y="1" width="40" height="20" rx="3" fill="none" stroke="#8a6d3b" stroke-width="1.5"/>` +
		`<line x1="41" y1="11" x2="78" y2="11" stroke="#8a6d3b" marker-end="url(#arrow)"/>` +
		`<text x="8" y="34" font-size="11" font-family="ui-monospace" text-anchor="start">read</text>` +
		`</svg><figcaption>the read path</figcaption></figure>`

	reduced := sanitise.Fragment(drawn)

	if reduced != drawn {
		t.Errorf("an illustration within the allowlist did not survive neutralising unchanged:\nwant %s\ngot  %s", drawn, reduced)
	}
}

func TestNothingSuppliedIsNothingEmbedded(t *testing.T) {
	if reduced := sanitise.Fragment("   \n "); reduced != "" {
		t.Errorf("a fragment of nothing but space reduces to %q", reduced)
	}
}
