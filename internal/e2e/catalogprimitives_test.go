package e2e_test

import (
	"bytes"
	"fmt"
	"image/png"
	"strings"
	"testing"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/chromedp"
)

const (
	badgeNeutralKey   = "a-neutral-tone-with-a-label"
	badgeAccentKey    = "an-accent-tone-with-a-label"
	badgeWarningKey   = "a-warning-tone-with-a-label"
	badgeErrorKey     = "an-error-tone-with-a-label"
	badgeIconAloneKey = "an-icon-alone"
	badgeExplainedKey = "a-label-that-says-what-it-means-on-hover-on-focus-and-to-a-screen-reader"

	overlineWordsAloneKey = "words-alone"
	overlineCountedKey    = "words-and-a-count"

	dotAroundKey   = "the-colour-around-it"
	dotNeutralKey  = "a-neutral-tone"
	dotAccentKey   = "an-accent-tone"
	dotResolvedKey = "a-resolved-tone"
	dotWarningKey  = "a-warning-tone"
	dotErrorKey    = "an-error-tone"

	cardSectionsKey = "a-body-and-a-footer"
	cardBodyKey     = "a-body-alone"
	cardPlainKey    = "no-sections"

	brandKey = "a-mark-and-a-name"

	catalogBadgeTitleSaid = "What this badge is"
	catalogBadgeHintSaid  = "A sentence saying what it means."
)

func (r *reader) paintedOf(selector string, properties ...string) map[string]string {
	r.t.Helper()

	var painted map[string]string
	r.evaluate(`(() => {
  const held = document.querySelector(`+jsonOf(r.t, selector)+`);
  if (!held) {
    return null;
  }
  const style = getComputedStyle(held);
  const painted = {};
  for (const property of `+jsonOfList(r.t, properties)+`) {
    painted[property] = style.getPropertyValue(property);
  }
  return painted;
})()`, &painted)
	if painted == nil {
		r.t.Fatalf("the catalog draws nothing at %s", selector)
	}
	return painted
}

func jsonOfList(t *testing.T, values []string) string {
	t.Helper()
	quoted := make([]string, len(values))
	for at, value := range values {
		quoted[at] = jsonOf(t, value)
	}
	return "[" + strings.Join(quoted, ",") + "]"
}

func (r *reader) pixelsOf(selector, property string) float64 {
	r.t.Helper()

	var pixels float64
	r.evaluate(`parseFloat(getComputedStyle(document.querySelector(`+jsonOf(r.t, selector)+`)).getPropertyValue(`+jsonOf(r.t, property)+`))`, &pixels)
	return pixels
}

func (r *reader) pixelsOfToken(token string) float64 {
	r.t.Helper()

	var pixels float64
	r.evaluate(`(() => {
  const probe = document.createElement('i');
  probe.style.cssText = 'position:absolute;visibility:hidden;width:var(`+token+`)';
  document.body.append(probe);
  const width = probe.getBoundingClientRect().width;
  probe.remove();
  return width;
})()`, &pixels)
	return pixels
}

func (r *reader) textOf(selector string) string {
	r.t.Helper()

	var text string
	r.evaluate(`document.querySelector(`+jsonOf(r.t, selector)+`).textContent.replace(/\s+/g, ' ').trim()`, &text)
	return text
}

func (r *reader) countOf(selector string) int {
	r.t.Helper()

	var counted int
	r.evaluate(`document.querySelectorAll(`+jsonOf(r.t, selector)+`).length`, &counted)
	return counted
}

func (r *reader) attributeOf(selector, attribute string) string {
	r.t.Helper()

	var value string
	r.evaluate(`document.querySelector(`+jsonOf(r.t, selector)+`).getAttribute(`+jsonOf(r.t, attribute)+`) || ''`, &value)
	return value
}

func badgeOf(key string) string {
	return stateOf("ui-badge", key) + " ui-badge"
}

func (r *reader) inkBetween(selector string, from, to float64) int {
	r.t.Helper()

	r.evaluate(`document.querySelector(`+jsonOf(r.t, selector)+`).scrollIntoView({ block: 'center' }); new Promise((settled) => requestAnimationFrame(() => requestAnimationFrame(settled)))`, nil, awaitingThePromise)
	var box struct{ Left, Top, Width, Height float64 }
	r.evaluate(`(() => { const box = document.querySelector(`+jsonOf(r.t, selector)+`).getBoundingClientRect(); return { Left: box.left, Top: box.top, Width: box.width, Height: box.height }; })()`, &box)
	painted := r.paintedOf(selector, "background-color")["background-color"]

	shot, err := png.Decode(bytes.NewReader(r.screenshot()))
	if err != nil {
		r.t.Fatalf("reading what the page painted: %v", err)
	}
	var wash [3]int
	if _, err := fmt.Sscanf(painted, "rgb(%d, %d, %d)", &wash[0], &wash[1], &wash[2]); err != nil {
		r.t.Fatalf("reading the colour %q the badge is painted: %v", painted, err)
	}

	inked := 0
	for x := int(box.Left + from); x < int(box.Left+to); x++ {
		for y := int(box.Top) + 5; y < int(box.Top+box.Height)-5; y++ {
			red, green, blue, _ := shot.At(x, y).RGBA()
			if abs(int(red>>8)-wash[0])+abs(int(green>>8)-wash[1])+abs(int(blue>>8)-wash[2]) > 60 {
				inked++
			}
		}
	}
	return inked
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

func TestTheCatalogPaintsEachBadgeToneInItsOwnWashAndInk(t *testing.T) {
	// scenario: A badge state with a tone paints its background and its text in that tone's tokens
	t.Parallel()
	page := catalogued(t)

	for _, tone := range []struct{ key, wash, ink, name string }{
		{badgeNeutralKey, "--surface-muted", "--fg-2", "neutral"},
		{badgeAccentKey, "--accent-wash", "--accent", "accent"},
		{badgeWarningKey, "--state-warning-wash", "--state-warning", "warning"},
		{badgeErrorKey, "--state-error-wash", "--state-error", "error"},
	} {
		t.Run(tone.name, func(t *testing.T) {
			badge := badgeOf(tone.key)

			painted := page.paintedOf(badge, "background-color", "color")

			if want := page.pageColourOf(tone.wash); painted["background-color"] != want {
				t.Errorf("the %s badge is painted %s, want %s", tone.name, painted["background-color"], want)
			}
			if want := page.pageColourOf(tone.ink); painted["color"] != want {
				t.Errorf("the %s badge's text is %s, want %s", tone.name, painted["color"], want)
			}
			if got := page.attributeOf(badge, "tone"); got != tone.name {
				t.Errorf("the badge states its tone as %q, want %q", got, tone.name)
			}
		})
	}
}

func TestTheCatalogDrawsABadgeWithALabelAsAnIconBesideItsWordInAPill(t *testing.T) {
	// scenario: A badge state with a label draws its icon and its word in a pill of the chip's height
	t.Parallel()
	page := catalogued(t)

	for _, key := range []string{badgeNeutralKey, badgeAccentKey, badgeWarningKey, badgeErrorKey} {
		t.Run(key, func(t *testing.T) {
			badge := badgeOf(key)

			if got := page.textOf(badge + ` [slot="label"]`); got != "A badge" {
				t.Errorf("the badge says %q, want its label", got)
			}
			if page.inkBetween(badge, 6, 18) == 0 {
				t.Error("the badge draws no icon beside its label")
			}
			if page.inkBetween(badge, 19, 21) != 0 {
				t.Error("the badge's label does not start after the icon and the gap that sets them apart")
			}
			if got, want := page.pixelsOf(badge, "height"), page.pixelsOfToken("--chip-height"); got != want {
				t.Errorf("the badge stands %vpx tall, want the chip height %vpx", got, want)
			}
			if got := page.pixelsOf(badge, "width"); got <= page.pixelsOfToken("--chip-height") {
				t.Errorf("the badge is %vpx wide, and a pill holding a word is wider than it is tall", got)
			}
		})
	}
}

func TestTheCatalogDrawsAnIconAloneBadgeAsACircleThatStillSaysWhatItMeansToAScreenReader(t *testing.T) {
	// scenario: A badge state holding an icon alone is a circle of the check's size, and is named to a screen reader and on hover
	t.Parallel()
	page := catalogued(t)
	badge := badgeOf(badgeIconAloneKey)

	if page.countOf(badge+` [slot="label"]`) != 0 {
		t.Error("the icon-alone badge draws a label")
	}
	check := page.pixelsOfToken("--check-size")
	if page.inkBetween(badge, 4, check-4) == 0 {
		t.Error("the icon-alone badge draws no icon")
	}
	painted := page.paintedOf(badge, "width", "height")
	if painted["width"] != painted["height"] || page.pixelsOf(badge, "width") != check {
		t.Errorf("the icon-alone badge is %s by %s, want a square of %vpx", painted["width"], painted["height"], check)
	}
	if got := page.pixelsOf(badge, "border-top-left-radius"); got < check/2 {
		t.Errorf("the icon-alone badge's corners are %vpx round, want a circle of at least %vpx", got, check/2)
	}

	if got := page.attributeOf(badge, "role"); got != "img" {
		t.Errorf("the icon-alone badge has the role %q, want img", got)
	}
	spoken := page.attributeOf(badge, "aria-label")
	if !strings.Contains(spoken, catalogBadgeTitleSaid) || !strings.Contains(spoken, catalogBadgeHintSaid) {
		t.Errorf("a screen reader hears %q from the icon-alone badge, want its title and its hint", spoken)
	}
	if page.attributeOf(badge, "tabindex") != "0" {
		t.Error("the icon-alone badge cannot take the keyboard cursor, so its meaning is out of reach without a pointer")
	}
}

func TestTheCatalogNamesAnExplainedBadgeToAScreenReaderAndShowsItsTipOnHoverAndOnFocus(t *testing.T) {
	// scenario: A badge state that says what it means is read aloud as its title and hint and shows both on hover and on focus
	t.Parallel()
	page := catalogued(t)
	page.do(sending(emulation.SetFocusEmulationEnabled, emulation.SetFocusEmulationEnabledParams{Enabled: true}))
	badge := badgeOf(badgeExplainedKey)

	if got := page.attributeOf(badge, "role"); got != "img" {
		t.Errorf("the explained badge has the role %q, want img", got)
	}
	if got, want := page.attributeOf(badge, "aria-label"), catalogBadgeTitleSaid+". "+catalogBadgeHintSaid; got != want {
		t.Errorf("a screen reader hears %q, want %q", got, want)
	}
	if page.attributeOf(badge, "tabindex") != "0" {
		t.Error("the explained badge cannot take the keyboard cursor")
	}
	if got := page.paintedOf(badge, "cursor")["cursor"]; got != "help" {
		t.Errorf("the pointer over the explained badge is %q, want help", got)
	}

	tipSays := `(() => {
  const tip = document.querySelector('.tooltip');
  return !tip.hidden && tip.getBoundingClientRect().height > 0 && tip.heading === ` + jsonOf(t, catalogBadgeTitleSaid) + ` && tip.hint === ` + jsonOf(t, catalogBadgeHintSaid) + `;
})()`
	tipIsHidden := `document.querySelector('.tooltip').hidden === true`

	page.evaluate(`document.querySelector(`+jsonOf(t, badge)+`).scrollIntoView({ block: 'center' }); new Promise((settled) => requestAnimationFrame(() => requestAnimationFrame(settled)))`, nil, awaitingThePromise)
	var centre struct{ X, Y float64 }
	page.evaluate(`(() => { const box = document.querySelector(`+jsonOf(t, badge)+`).getBoundingClientRect(); return { X: box.left + box.width / 2, Y: box.top + box.height / 2 }; })()`, &centre)
	page.do(sending(input.DispatchMouseEvent, input.DispatchMouseEventParams{Type: input.DispatchMouseEventTypeMouseMoved, X: centre.X, Y: centre.Y}))
	page.waitFor(tipSays)

	page.do(sending(input.DispatchMouseEvent, input.DispatchMouseEventParams{Type: input.DispatchMouseEventTypeMouseMoved, X: 1, Y: 1}))
	page.waitFor(tipIsHidden)

	page.do(chromedp.Focus(chromedp.CSS(badge)))
	page.waitFor(tipSays)
}

func TestAShadowRootHoldsEveryBadgeOfTheCatalogAndTheyStillCarryTheirTone(t *testing.T) {
	// scenario: Every ui-badge the catalog draws renders into a shadow root that adopts the shared stylesheet
	t.Parallel()
	page := catalogued(t)

	var drawn struct{ Badges, Shadowed, Adopting int }
	page.evaluate(`(() => {
  const badges = Array.from(document.querySelectorAll('ui-badge'));
  return {
    Badges: badges.length,
    Shadowed: badges.filter((badge) => badge.shadowRoot !== null).length,
    Adopting: badges.filter((badge) => badge.shadowRoot !== null && badge.shadowRoot.adoptedStyleSheets.length > 0).length,
  };
})()`, &drawn)

	if drawn.Badges == 0 {
		t.Fatal("the catalog draws no badge, so this proves nothing")
	}
	if drawn.Shadowed != drawn.Badges || drawn.Adopting != drawn.Badges {
		t.Errorf("of %d badges, %d have a shadow root and %d of those adopt the shared stylesheet, want all", drawn.Badges, drawn.Shadowed, drawn.Adopting)
	}
}

func TestABadgeWhoseLabelIsRemovedBecomesTheCircleOfTheChecksSize(t *testing.T) {
	// scenario: A badge that loses its slotted label is a circle of the check's size, and one that gains it is a pill of the chip's height
	t.Parallel()
	page := catalogued(t)
	badge := badgeOf(badgeNeutralKey)
	check := page.pixelsOfToken("--check-size")

	page.evaluate(`(() => { const label = document.querySelector(`+jsonOf(t, badge+` [slot="label"]`)+`); window.__label = label; label.remove(); })()`, nil)
	page.waitFor(`Math.round(document.querySelector(` + jsonOf(t, badge) + `).getBoundingClientRect().width) === ` + fmt.Sprint(check))
	if got := page.pixelsOf(badge, "height"); got != check {
		t.Errorf("a badge without its label is %vpx tall, want the check's %vpx", got, check)
	}

	page.evaluate(`document.querySelector(`+jsonOf(t, badge)+`).append(window.__label)`, nil)
	page.waitFor(`document.querySelector(` + jsonOf(t, badge) + `).getBoundingClientRect().width > ` + fmt.Sprint(check))
	if got, want := page.pixelsOf(badge, "height"), page.pixelsOfToken("--chip-height"); got != want {
		t.Errorf("a badge given its label back is %vpx tall, want the chip's %vpx", got, want)
	}
}

func (r *reader) renderedTextOf(selector string) string {
	r.t.Helper()

	var text string
	r.evaluate(`(() => {
  const drawn = (node) => {
    if (node.nodeType === Node.COMMENT_NODE) {
      return '';
    }
    if (node instanceof HTMLSlotElement) {
      return node.assignedNodes({ flatten: true }).map(drawn).join('');
    }
    return node.nodeType === Node.TEXT_NODE ? node.textContent : Array.from(node.childNodes).map(drawn).join('');
  };
  return drawn(document.querySelector(`+jsonOf(r.t, selector)+`).shadowRoot).replace(/\s+/g, ' ').trim();
})()`, &text)
	return text
}

func (r *reader) insideTheShadowRootOf(selector, within string, properties ...string) map[string]string {
	r.t.Helper()

	var painted map[string]string
	r.evaluate(`(() => {
  const held = document.querySelector(`+jsonOf(r.t, selector)+`).shadowRoot.querySelector(`+jsonOf(r.t, within)+`);
  if (!held) {
    return null;
  }
  const style = getComputedStyle(held);
  const painted = {};
  for (const property of `+jsonOfList(r.t, properties)+`) {
    painted[property] = style.getPropertyValue(property);
  }
  return painted;
})()`, &painted)
	if painted == nil {
		r.t.Fatalf("%s draws nothing at %s", selector, within)
	}
	return painted
}

func overlineOfState(key string) string {
	return stateOf("ui-overline", key) + " ui-overline"
}

func TestTheCatalogDrawsAnOverlineWithWordsAloneInCapitalsAndNoCount(t *testing.T) {
	// scenario: An overline state with words alone draws them in capitals in the quiet ink and draws no count
	t.Parallel()
	page := catalogued(t)
	overline := overlineOfState(overlineWordsAloneKey)

	painted := page.paintedOf(overline, "text-transform", "color", "font-weight")

	if got := page.textOf(overline); got != "An overline" {
		t.Errorf("the overline says %q, want its words alone", got)
	}
	if painted["text-transform"] != "uppercase" {
		t.Errorf("the overline's words are %q, want uppercase", painted["text-transform"])
	}
	if want := page.pageColourOf("--fg-3"); painted["color"] != want {
		t.Errorf("the overline is %s, want %s", painted["color"], want)
	}
	if painted["font-weight"] != "700" {
		t.Errorf("the overline's weight is %s, want bold", painted["font-weight"])
	}
	if got := page.renderedTextOf(overline); got != "An overline" {
		t.Errorf("the overline draws %q, want its words alone and no count it was not given", got)
	}
}

func TestTheCatalogDrawsAnOverlinesCountAfterItsWordsInALighterWeight(t *testing.T) {
	// scenario: An overline state with a count draws the number after the words, set apart and lighter than the words
	t.Parallel()
	page := catalogued(t)
	overline := overlineOfState(overlineCountedKey)

	if got := page.textOf(overline); got != "An overline" {
		t.Errorf("the overline's own words are %q, want its words with the count left to the element", got)
	}
	if got := page.attributeOf(overline, "count"); got != "3" {
		t.Errorf("the overline is given the count %q, want 3", got)
	}
	if got := page.renderedTextOf(overline); got != "An overline · 3" {
		t.Errorf("the overline draws %q, want its words and its count set apart by a dot", got)
	}
	words := page.paintedOf(overline, "font-weight")
	counted := page.insideTheShadowRootOf(overline, "span:last-child", "font-weight", "letter-spacing")
	if counted["font-weight"] != "500" || words["font-weight"] != "700" {
		t.Errorf("the count is weight %s beside words of weight %s, want 500 beside 700", counted["font-weight"], words["font-weight"])
	}
	if counted["letter-spacing"] != "normal" && counted["letter-spacing"] != "0px" {
		t.Errorf("the count is spaced %q, want it set tight", counted["letter-spacing"])
	}
	if spread := page.insideTheShadowRootOf(overline, "span:first-child", "letter-spacing")["letter-spacing"]; spread == counted["letter-spacing"] {
		t.Errorf("the count and the words are spaced alike (%q), want the words spread", spread)
	}
}

func dotOf(key string) string {
	return stateOf("ui-status-dot", key) + " ui-status-dot"
}

func TestTheCatalogPaintsEachStatusDotToneInItsToken(t *testing.T) {
	// scenario: A status dot state with a tone paints a six pixel circle in that tone's token and hides from a screen reader
	t.Parallel()
	page := catalogued(t)

	for _, tone := range []struct{ key, token, name string }{
		{dotNeutralKey, "--fg-2", "neutral"},
		{dotAccentKey, "--accent", "accent"},
		{dotResolvedKey, "--state-resolved", "resolved"},
		{dotWarningKey, "--state-warning", "warning"},
		{dotErrorKey, "--state-error", "error"},
	} {
		t.Run(tone.name, func(t *testing.T) {
			dot := dotOf(tone.key)

			painted := page.paintedOf(dot, "background-color", "width", "height")

			if want := page.pageColourOf(tone.token); painted["background-color"] != want {
				t.Errorf("the %s dot is painted %s, want %s", tone.name, painted["background-color"], want)
			}
			if painted["width"] != "6px" || painted["height"] != "6px" || page.pixelsOf(dot, "border-top-left-radius") < 3 {
				t.Errorf("the %s dot is %s by %s with corners %vpx round, want a six pixel circle", tone.name, painted["width"], painted["height"], page.pixelsOf(dot, "border-top-left-radius"))
			}
			if got := page.attributeOf(dot, "aria-hidden"); got != "true" {
				t.Errorf("the %s dot is hidden from a screen reader as %q, want true", tone.name, got)
			}
		})
	}
}

func TestTheCatalogPaintsAStatusDotWithNoToneInTheColourAroundIt(t *testing.T) {
	// scenario: A status dot state with no tone takes the colour of the text around it
	t.Parallel()
	page := catalogued(t)
	dot := dotOf(dotAroundKey)

	painted := page.paintedOf(dot, "background-color", "color")

	if page.attributeOf(dot, "tone") != "" {
		t.Error("the dot with no tone states a tone")
	}
	if painted["background-color"] != painted["color"] {
		t.Errorf("the dot is painted %s beside text of %s, want the colour around it", painted["background-color"], painted["color"])
	}
	if painted["background-color"] == "rgba(0, 0, 0, 0)" {
		t.Error("the dot paints nothing")
	}
}

func cardOf(key string) string {
	return stateOf("ui-card", key) + " ui-card"
}

func TestTheCatalogDrawsACardWithABodyAndAFooterWithTheFooterSunkUnderTheBody(t *testing.T) {
	// scenario: A card state with a body and a footer pads each section and sets the footer under the body on a sunk ground
	t.Parallel()
	page := catalogued(t)
	card := cardOf(cardSectionsKey)
	body := card + ` [data-part="body"]`
	footer := card + ` [data-part="footer"]`

	painted := page.paintedOf(card, "display", "background-color", "border-top-width", "border-top-left-radius")
	if painted["display"] != "block" || painted["border-top-width"] != "1px" {
		t.Errorf("the card is %q with a border of %s, want a block with a one pixel border", painted["display"], painted["border-top-width"])
	}
	if want := page.pageColourOf("--surface-card"); painted["background-color"] != want {
		t.Errorf("the card is painted %s, want %s", painted["background-color"], want)
	}
	if painted["border-top-left-radius"] == "0px" {
		t.Error("the card has square corners")
	}

	pad := page.pixelsOfToken("--card-pad")
	for name, section := range map[string]string{"body": body, "footer": footer} {
		if got := page.pixelsOf(section, "padding-left"); got != pad {
			t.Errorf("the card's %s is padded %vpx, want %vpx", name, got, pad)
		}
	}
	ground := page.paintedOf(footer, "background-color", "border-top-width")
	if want := page.pageColourOf("--surface-sunk"); ground["background-color"] != want {
		t.Errorf("the footer is painted %s, want %s", ground["background-color"], want)
	}
	if ground["border-top-width"] != "1px" {
		t.Errorf("the footer's top border is %s, want a rule one pixel thick", ground["border-top-width"])
	}

	var order struct {
		BodyBottom float64 `json:"bodyBottom"`
		FooterTop  float64 `json:"footerTop"`
	}
	page.evaluate(`({
  bodyBottom: document.querySelector(`+jsonOf(t, body)+`).getBoundingClientRect().bottom,
  footerTop: document.querySelector(`+jsonOf(t, footer)+`).getBoundingClientRect().top,
})`, &order)
	if order.FooterTop < order.BodyBottom-1 {
		t.Errorf("the footer begins at %v, above the body's end at %v", order.FooterTop, order.BodyBottom)
	}
}

func TestTheCatalogDrawsACardWithABodyAloneAndNoFooter(t *testing.T) {
	// scenario: A card state with a body alone pads its body and draws no footer or rule under it
	t.Parallel()
	page := catalogued(t)
	card := cardOf(cardBodyKey)

	if page.countOf(card+` [data-part="footer"]`) != 0 {
		t.Error("the body-alone card draws a footer")
	}
	if got, want := page.pixelsOf(card+` [data-part="body"]`, "padding-top"), page.pixelsOfToken("--card-pad"); got != want {
		t.Errorf("the body is padded %vpx, want %vpx", got, want)
	}
	if got := page.paintedOf(card, "border-top-width")["border-top-width"]; got != "1px" {
		t.Errorf("the body-alone card's border is %s, want one pixel", got)
	}
}

func TestTheCatalogDrawsACardWithNoSectionsAsTheCardAroundWhatItHolds(t *testing.T) {
	// scenario: A card state with no sections draws the border and ground around its content and pads nothing of its own
	t.Parallel()
	page := catalogued(t)
	card := cardOf(cardPlainKey)

	if page.countOf(card+` [data-part]`) != 0 {
		t.Error("the card with no sections draws a section")
	}
	if page.countOf(card+` p.catalog-card-plain`) != 1 {
		t.Error("the card with no sections holds no content")
	}
	painted := page.paintedOf(card, "background-color", "border-top-width", "padding-left")
	if want := page.pageColourOf("--surface-card"); painted["background-color"] != want {
		t.Errorf("the card is painted %s, want %s", painted["background-color"], want)
	}
	if painted["border-top-width"] != "1px" {
		t.Errorf("the card's border is %s, want one pixel", painted["border-top-width"])
	}
	if painted["padding-left"] != "0px" {
		t.Errorf("the card pads %s of its own, want none, since the page lays its content out", painted["padding-left"])
	}
}

func TestTheCatalogDrawsABrandAsADecorativeMarkBeforeItsNameInTheDisplayFace(t *testing.T) {
	// scenario: A brand state draws its mark, said to no one, before its name set in the display face and the strong ink
	t.Parallel()
	page := catalogued(t)
	brand := stateOf("ui-brand", brandKey) + " ui-brand"
	mark := brand + ` img[slot="mark"]`
	name := brand + ` span`

	if got := page.textOf(name); got != "A brand" {
		t.Errorf("the brand says %q, want its name", got)
	}
	page.waitFor(`document.querySelector(` + jsonOf(t, mark) + `).complete === true`)
	var markBox struct {
		Natural float64 `json:"natural"`
		Source  string  `json:"source"`
		Alt     string  `json:"alt"`
		Width   float64 `json:"width"`
		Height  float64 `json:"height"`
	}
	page.evaluate(`(() => {
  const held = document.querySelector(`+jsonOf(t, mark)+`);
  const box = held.getBoundingClientRect();
  return { natural: held.naturalWidth, source: held.getAttribute('src'), alt: held.getAttribute('alt'), width: box.width, height: box.height };
})()`, &markBox)
	if markBox.Natural == 0 {
		t.Errorf("the mark %q did not load", markBox.Source)
	}
	if markBox.Alt != "" {
		t.Errorf("the mark is named %q, want it left to the name beside it", markBox.Alt)
	}
	if markBox.Width != 28 || markBox.Height != 28 {
		t.Errorf("the mark is %vx%v, want 28x28", markBox.Width, markBox.Height)
	}

	var row struct {
		MarkRight, NameLeft, MarkMiddle, NameMiddle float64
	}
	page.evaluate(`(() => {
  const mark = document.querySelector(`+jsonOf(t, mark)+`).getBoundingClientRect();
  const name = document.querySelector(`+jsonOf(t, name)+`).getBoundingClientRect();
  return { MarkRight: mark.right, NameLeft: name.left, MarkMiddle: mark.top + mark.height / 2, NameMiddle: name.top + name.height / 2 };
})()`, &row)
	if row.NameLeft < row.MarkRight {
		t.Errorf("the name begins at %v, before the mark ends at %v", row.NameLeft, row.MarkRight)
	}
	if gap := row.MarkMiddle - row.NameMiddle; gap > 1 || gap < -1 {
		t.Errorf("the name sits %vpx off the mark's middle, want them on one line", gap)
	}

	painted := page.paintedOf(name, "font-family", "color", "font-weight")
	normalise := strings.NewReplacer(`"`, "", `'`, "", " ", "")
	if got, want := normalise.Replace(painted["font-family"]), normalise.Replace(page.pageToken("--font-display")); got != want {
		t.Errorf("the name is set in %s, want the display face %s", got, want)
	}
	if want := page.pageColourOf("--fg-1"); painted["color"] != want {
		t.Errorf("the name is %s, want %s", painted["color"], want)
	}
	if painted["font-weight"] != "700" {
		t.Errorf("the name's weight is %s, want bold", painted["font-weight"])
	}
}
