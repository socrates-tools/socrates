package e2e_test

import (
	"bytes"
	"fmt"
	"image/png"
	"math"
	"net/http"
	"strings"
	"testing"

	"github.com/chromedp/chromedp/kb"
)

const (
	catalogRoute         = "/catalog"
	developmentBuildTag  = "agoradev"
	dialogOpeners        = `Array.from(document.querySelectorAll('button.catalog-opener')).map((held) => held.id)`
	theCatalogIsDrawnFor = `ui-dialog`
)

func catalogued(t *testing.T) *reader {
	t.Helper()

	held := servingSurfaceBuiltWith(t, developmentBuildTag)
	address := held.serve(oneDecisionDocument)
	return openingFor(t, originOf(t, address)+catalogRoute, theCatalogIsDrawnFor, ".catalog-state")
}

func dialogOf(opener string) string {
	return fmt.Sprintf(`document.getElementById(document.getElementById(%q).dataset.opens)`, opener)
}

func isOpen(opener string) string {
	return dialogOf(opener) + `.isOpen === true`
}

func isClosed(opener string) string {
	return dialogOf(opener) + `.isOpen === false`
}

func isModal(opener string) string {
	return dialogOf(opener) + `.isModal`
}

func theCursorIsOn(opener string) string {
	return fmt.Sprintf(`document.activeElement === document.getElementById(%q)`, opener)
}

func eachDialogState(t *testing.T, page *reader, prove func(t *testing.T, opener string)) {
	t.Helper()

	var openers []string
	page.evaluate(dialogOpeners, &openers)
	if len(openers) < 2 {
		t.Fatalf("the catalog offers %d openers, and the dialog has more than one state to open: %v", len(openers), openers)
	}
	for _, opener := range openers {
		t.Run(opener, func(t *testing.T) {
			prove(t, opener)
		})
	}
}

func TestTheCatalogOpensEachDialogStateAsAModal(t *testing.T) {
	// scenario: Pressing the opener of a dialog state in the catalog opens it as a modal
	t.Parallel()
	page := catalogued(t)

	eachDialogState(t, page, func(t *testing.T, opener string) {
		page.waitFor(isClosed(opener))

		page.click("#" + opener)

		page.waitFor(isOpen(opener))
		page.waitFor(isModal(opener))

		page.press(kb.Escape)
		page.waitFor(isClosed(opener))
	})
}

func TestTheCatalogClosesEachDialogStateOnEscapeAndReturnsTheCursor(t *testing.T) {
	// scenario: Escape closes a dialog state in the catalog and returns the cursor to its opener
	t.Parallel()
	page := catalogued(t)

	eachDialogState(t, page, func(t *testing.T, opener string) {
		page.click("#" + opener)
		page.waitFor(isOpen(opener))

		page.press(kb.Escape)

		page.waitFor(isClosed(opener))
		page.waitFor(theCursorIsOn(opener))
	})
}

func TestTheCatalogClosesEachDialogStateOnAPressOnTheBackdropAndReturnsTheCursor(t *testing.T) {
	// scenario: A press on the backdrop closes a dialog state in the catalog and returns the cursor to its opener
	t.Parallel()
	page := catalogued(t)

	eachDialogState(t, page, func(t *testing.T, opener string) {
		page.click("#" + opener)
		page.waitFor(isOpen(opener))

		page.clickAt(4, 4)

		page.waitFor(isClosed(opener))
		page.waitFor(theCursorIsOn(opener))
	})
}

func TestTheCatalogKeepsADialogStateOpenOnAPressInsideIt(t *testing.T) {
	// scenario: A press on the dialog's own padding in the catalog does not close it
	t.Parallel()
	page := catalogued(t)

	eachDialogState(t, page, func(t *testing.T, opener string) {
		page.click("#" + opener)
		page.waitFor(isOpen(opener))

		var padding struct{ X, Y float64 }
		page.evaluate(`(() => {
  const box = `+dialogOf(opener)+`.querySelector('[slot="heading"]').getBoundingClientRect();
  return { X: box.left - 4, Y: box.top - 4 };
})()`, &padding)
		page.clickAt(padding.X, padding.Y)

		page.waitFor(isOpen(opener))
		page.press(kb.Escape)
		page.waitFor(isClosed(opener))
	})
}

func TestTheCatalogOpensADialogStateAgainAfterItWasClosed(t *testing.T) {
	// scenario: A dialog state in the catalog opens again after Escape closed it
	t.Parallel()
	page := catalogued(t)

	eachDialogState(t, page, func(t *testing.T, opener string) {
		page.click("#" + opener)
		page.waitFor(isOpen(opener))
		page.press(kb.Escape)
		page.waitFor(isClosed(opener))

		page.click("#" + opener)

		page.waitFor(isModal(opener))
		page.press(kb.Escape)
		page.waitFor(isClosed(opener))
	})
}

func TestTheCatalogPageStartsNoLiveStreamAndReachesOnlyItsOwnAssets(t *testing.T) {
	// scenario: The catalog draws the primitives with the page's own script and the page's features start only where there is a surface
	t.Parallel()
	page := catalogued(t)
	var openers []string
	page.evaluate(dialogOpeners, &openers)

	for _, asked := range page.requests.all() {
		if strings.HasSuffix(asked, "/live") {
			t.Errorf("the catalog requested %q, and a page with no surface follows no exchange", asked)
		}
		if !strings.Contains(asked, "/assets/") && !strings.HasSuffix(asked, catalogRoute) {
			t.Errorf("the catalog requested %q, and it asks for nothing but itself and the bundle", asked)
		}
	}
}

func TestTheDevelopmentBuildServesTheCatalogAtARouteUnderTheSurface(t *testing.T) {
	// scenario: A person opens the catalog on the development build
	t.Parallel()
	held := servingSurfaceBuiltWith(t, developmentBuildTag)
	address := held.serve(oneDecisionDocument)

	answered, said := fetched(t, originOf(t, address)+catalogRoute)

	if answered != http.StatusOK {
		t.Fatalf("the development build answered the catalog route with %d: %s", answered, said)
	}
	if !strings.Contains(said, "<ui-dialog") {
		t.Errorf("the catalog the development build serves draws no dialog: %s", said)
	}
}

func TestAReleasedBinaryAnswersTheCatalogRouteWithNotFound(t *testing.T) {
	// scenario: A released binary has no route to the catalog
	t.Parallel()
	held := standing(t)
	address := held.serve(oneDecisionDocument)

	answered, said := fetched(t, originOf(t, address)+catalogRoute)

	if answered != http.StatusNotFound {
		t.Errorf("a binary built without the %s tag answered the catalog route with %d: %s", developmentBuildTag, answered, said)
	}
}

const everyCardAndDialogAdoptsTheSharedStylesheet = `(() => {
  const held = Array.from(document.querySelectorAll('ui-card, ui-dialog'));
  return held.length > 0 && held.every((element) => element.shadowRoot !== null && element.shadowRoot.adoptedStyleSheets.length > 0);
})()`

func TestEveryCardAndDialogOfTheCatalogHasAShadowRootThatAdoptsTheSharedStylesheet(t *testing.T) {
	// scenario: The catalog's cards and dialogs render into shadow roots that adopt the one shared sheet and the catalog's policy blocks nothing
	t.Parallel()
	page := catalogued(t)

	page.waitFor(everyCardAndDialogAdoptsTheSharedStylesheet)
	page.evaluate(afterTwoFrames, nil, awaitingThePromise)

	if violations := page.violations(); len(violations) != 0 {
		t.Errorf("the catalog's policy blocked %v while the cards and the dialogs adopted the shared stylesheet", violations)
	}
}

func (r *reader) runOfTheColourAround(x, y int, colour string) (left, right int) {
	r.t.Helper()

	shot, err := png.Decode(bytes.NewReader(r.screenshot()))
	if err != nil {
		r.t.Fatalf("reading what the page painted: %v", err)
	}
	var ground [3]int
	if _, err := fmt.Sscanf(colour, "rgb(%d, %d, %d)", &ground[0], &ground[1], &ground[2]); err != nil {
		r.t.Fatalf("reading the colour %q: %v", colour, err)
	}
	same := func(at int) bool {
		red, green, blue, _ := shot.At(at, y).RGBA()
		return int(red>>8) == ground[0] && int(green>>8) == ground[1] && int(blue>>8) == ground[2]
	}
	left, right = x, x
	for same(left-1) && left > 0 {
		left--
	}
	for same(right+1) && right < shot.Bounds().Max.X-1 {
		right++
	}
	return left, right + 1
}

func TestTheCatalogDrawsAnOpenDialogAsACardGroundOfFourEightyPixelsPaddedByTheCardsPad(t *testing.T) {
	// scenario: An open dialog state is drawn on the card's ground, 480px wide, with the card's padding between its edge and its heading
	t.Parallel()
	page := catalogued(t)

	eachDialogState(t, page, func(t *testing.T, opener string) {
		page.click("#" + opener)
		page.waitFor(isOpen(opener))

		var heading struct{ Left, Top float64 }
		page.evaluate(`(() => {
  const box = `+dialogOf(opener)+`.querySelector('[slot="heading"]').getBoundingClientRect();
  return { Left: box.left, Top: box.top };
})()`, &heading)
		left, right := page.runOfTheColourAround(int(heading.Left), int(heading.Top)-4, page.pageColourOf("--surface-card"))

		if got := right - left; got < 479 || got > 481 {
			t.Errorf("the dialog is painted %dpx wide on the card's ground, want 480px", got)
		}
		if got, want := heading.Left-float64(left), page.pixelsOfToken("--card-pad"); math.Abs(got-want) > 1 {
			t.Errorf("the dialog pads its heading %vpx from its edge, want %vpx", got, want)
		}

		page.press(kb.Escape)
		page.waitFor(isClosed(opener))
	})
}

func TestTheCatalogsDialogHeadsItsBodyAndItsActionsInOrderThroughSlots(t *testing.T) {
	// scenario: The dialog state with everything draws its slotted heading above its body above its actions, and the state with a heading alone writes no body or actions
	t.Parallel()
	page := catalogued(t)
	headingAlone := stateOf("ui-dialog", "heading-alone")

	for _, slot := range []string{"body", "actions"} {
		if got := page.countOf(headingAlone + ` ui-dialog [slot="` + slot + `"]`); got != 0 {
			t.Errorf("the dialog state with a heading alone writes %d elements into its %s slot", got, slot)
		}
	}

	page.click("#open-catalog-dialog")
	page.waitFor(isOpen("open-catalog-dialog"))
	var order struct{ Heading, Body, Actions float64 }
	page.evaluate(`(() => {
  const host = `+dialogOf("open-catalog-dialog")+`;
  const top = (slot) => host.querySelector('[slot="' + slot + '"]').getBoundingClientRect().top;
  return { Heading: top('heading'), Body: top('body'), Actions: top('actions') };
})()`, &order)
	if !(order.Heading < order.Body && order.Body < order.Actions) {
		t.Errorf("the dialog draws its heading at %v, its body at %v and its actions at %v, want them in that order", order.Heading, order.Body, order.Actions)
	}
}
