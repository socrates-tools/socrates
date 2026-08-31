package stylesheet_test

import (
	"strings"
	"testing"

	"github.com/ktei/socrates/internal/stylesheet"
)

var illustration = stylesheet.Scope{Selector: ".painted", Prefix: "drawn-"}

func TestAStylesheetOutsideTheAllowlistIsReduced(t *testing.T) {
	for _, reduced := range []struct {
		named    string
		supplied string
		gone     []string
		kept     []string
	}{
		{
			named:    "an import",
			supplied: `@import url(https://elsewhere/x.css);.panel{color:#222}`,
			gone:     []string{"import", "elsewhere"},
			kept:     []string{".drawn-panel", "color:#222"},
		},
		{
			named:    "a font face",
			supplied: `@font-face{font-family:away;src:url(https://elsewhere/x.woff2)}.panel{color:#222}`,
			gone:     []string{"font-face", "elsewhere", "away"},
			kept:     []string{".drawn-panel", "color:#222"},
		},
		{
			named:    "a property outside the allowlist",
			supplied: `.panel{color:#222;behavior:evil;contain:none;-moz-binding:evil}`,
			gone:     []string{"behavior", "contain", "binding"},
			kept:     []string{"color:#222"},
		},
		{
			named:    "a value naming a remote resource",
			supplied: `.panel{color:#222;background:#fff url(https://elsewhere/p.png) no-repeat}`,
			gone:     []string{"elsewhere", "background"},
			kept:     []string{"color:#222"},
		},
		{
			named:    "a value naming a resource in escapes",
			supplied: `.panel{color:#222;background:\75 rl(https://elsewhere/p.png)}`,
			gone:     []string{"elsewhere", "background"},
			kept:     []string{"color:#222"},
		},
		{
			named:    "a value naming a resource split by a comment",
			supplied: `.panel{color:#222;background:u/**/rl(https://elsewhere/p.png)}`,
			gone:     []string{"elsewhere", "background"},
			kept:     []string{"color:#222"},
		},
		{
			named:    "a value naming a function that fetches",
			supplied: `.panel{color:#222;background-image:image-set("https://elsewhere/p.png" 1x)}`,
			gone:     []string{"elsewhere", "image-set"},
			kept:     []string{"color:#222"},
		},
		{
			named:    "a rule that closes the element it is embedded in",
			supplied: `.panel{color:#222;content:"</style><script>steal()</script>"}`,
			gone:     []string{"<", "steal"},
			kept:     []string{"color:#222"},
		},
		{
			named:    "an at-rule the allowlist does not name",
			supplied: `@keyframes away{from{opacity:0}to{opacity:1}}.panel{color:#222}`,
			gone:     []string{"keyframes", "away", "opacity"},
			kept:     []string{"color:#222"},
		},
		{
			named:    "a media query the allowlist does name",
			supplied: `@media (min-width:20rem){.panel{color:#222}}`,
			gone:     []string{},
			kept:     []string{"@media (min-width:20rem)", ".painted .drawn-panel", "color:#222"},
		},
	} {
		t.Run(reduced.named, func(t *testing.T) {
			embedded := stylesheet.Reduce(reduced.supplied, illustration)

			for _, absent := range reduced.gone {
				if strings.Contains(embedded, absent) {
					t.Errorf("%s survives the reduction as %q:\n%s", reduced.named, absent, embedded)
				}
			}
			for _, present := range reduced.kept {
				if !strings.Contains(embedded, present) {
					t.Errorf("reducing %s took %q with it:\n%s", reduced.named, present, embedded)
				}
			}
		})
	}
}

func TestAStylesheetThatDoesNotParseEmbedsNothing(t *testing.T) {
	for _, unparseable := range []struct {
		named    string
		supplied string
	}{
		{named: "a stray closing brace", supplied: `.panel{color:#222}}}`},
		{named: "a declaration outside any rule", supplied: `.panel{color:#222}color:#333;`},
	} {
		t.Run(unparseable.named, func(t *testing.T) {
			if embedded := stylesheet.Reduce(unparseable.supplied, illustration); embedded != "" {
				t.Errorf("%s embedded what the reduction managed to read:\n%s", unparseable.named, embedded)
			}
		})
	}
}

func TestACallerMayStateItsPaletteOnce(t *testing.T) {
	embedded := stylesheet.Reduce(`:root{--panel-ink:#1d1d1b;--panel-face:#fbf9f4}.panel{color:var(--panel-ink);background:var(--panel-face)}`, illustration)

	for _, stated := range []string{
		"--panel-ink:#1d1d1b",
		"--panel-face:#fbf9f4",
		"color:var(--panel-ink)",
		"background:var(--panel-face)",
	} {
		if !strings.Contains(embedded, stated) {
			t.Errorf("the reduction dropped %q, so a caller cannot state its palette once:\n%s", stated, embedded)
		}
	}
	if strings.Contains(embedded, ":root") {
		t.Errorf("the palette is stated on the page's own root rather than on the illustration:\n%s", embedded)
	}
}

func TestASelectorNamingThePageIsConfinedToTheIllustration(t *testing.T) {
	for _, confined := range []struct {
		named    string
		supplied string
		kept     string
	}{
		{named: "the document root", supplied: `:root{color:#222}`, kept: ".painted{"},
		{named: "the body", supplied: `body{color:#222}`, kept: ".painted{"},
		{named: "a name the page uses for its own chrome", supplied: `.decision .commit{color:#222}`, kept: ".painted .drawn-decision .drawn-commit{"},
		{named: "a decision of the page's own", supplied: `#decision-f3{color:#222}`, kept: ".painted #drawn-decision-f3{"},
		{named: "everything", supplied: `*{color:#222}`, kept: ".painted *{"},
	} {
		t.Run(confined.named, func(t *testing.T) {
			embedded := stylesheet.Reduce(confined.supplied, illustration)

			if !strings.Contains(embedded, confined.kept) {
				t.Errorf("a rule for %s was not confined to the illustration as %q:\n%s", confined.named, confined.kept, embedded)
			}
			for unscoped := range strings.SplitSeq(embedded, "}") {
				if selector, _, found := strings.Cut(unscoped, "{"); found && !strings.HasPrefix(strings.TrimSpace(selector), illustration.Selector) {
					t.Errorf("the reduction embedded %q, which names something outside the illustration", selector)
				}
			}
		})
	}
}

func TestASelectorThatCannotBeScopedIsDropped(t *testing.T) {
	for _, dropped := range []string{
		`.panel:has(.row){color:#222}`,
		`:host{color:#222}`,
		`.panel::part(row){color:#222}`,
		`.panel::slotted(span){color:#222}`,
	} {
		if embedded := stylesheet.Reduce(dropped, illustration); embedded != "" {
			t.Errorf("%q was embedded rather than dropped:\n%s", dropped, embedded)
		}
	}
}

func TestWhatTheCallerPaintedIsStillPainted(t *testing.T) {
	embedded := stylesheet.Reduce(
		`.panel{border:1px solid #d8d2c4;padding:0.75rem}`+
			`.panel .row+.row{border-top:1px solid #d8d2c4}`+
			`#total{font-weight:600}`+
			`.chart rect{fill:url(#band)}`,
		illustration)

	for _, painted := range []string{
		".painted .drawn-panel{border:1px solid #d8d2c4;padding:0.75rem;}",
		".painted .drawn-panel .drawn-row+.drawn-row{border-top:1px solid #d8d2c4;}",
		".painted #drawn-total{font-weight:600;}",
		".painted .drawn-chart rect{fill:url(#drawn-band);}",
	} {
		if !strings.Contains(embedded, painted) {
			t.Errorf("scoping took %q with it:\n%s", painted, embedded)
		}
	}
}

func TestNothingSuppliedIsNothingEmbedded(t *testing.T) {
	if embedded := stylesheet.Reduce("  \n ", illustration); embedded != "" {
		t.Errorf("a stylesheet of nothing but space reduces to %q", embedded)
	}
}
