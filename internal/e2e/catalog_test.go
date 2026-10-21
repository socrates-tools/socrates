package e2e_test

import (
	"fmt"
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
	return dialogOf(opener) + `.open === true`
}

func isClosed(opener string) string {
	return dialogOf(opener) + `.open === false`
}

func isModal(opener string) string {
	return dialogOf(opener) + `.matches(':modal')`
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
  const box = `+dialogOf(opener)+`.getBoundingClientRect();
  return { X: box.left + 4, Y: box.top + 4 };
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
