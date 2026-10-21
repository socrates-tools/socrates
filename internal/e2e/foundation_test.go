package e2e_test

import "testing"

const everyBadgeAdoptsTheSharedStylesheet = `(() => {
  const badges = Array.from(document.querySelectorAll('ui-badge'));
  return badges.length > 0 && badges.every((badge) => badge.shadowRoot !== null && badge.shadowRoot.adoptedStyleSheets.length > 0);
})()`

const afterTwoFrames = `new Promise((settled) => requestAnimationFrame(() => requestAnimationFrame(settled)))`

func TestThePageAdoptsTheSharedStylesheetInEveryBadgeWithoutAPolicyViolation(t *testing.T) {
	// scenario: The page's badges adopt the one shared sheet and the page's policy blocks nothing
	t.Parallel()
	_, page := read(t, glancedDocument)

	page.waitFor(everyBadgeAdoptsTheSharedStylesheet)
	page.evaluate(afterTwoFrames, nil, awaitingThePromise)

	if violations := page.violations(); len(violations) != 0 {
		t.Errorf("the page's policy blocked %v while the badges adopted the shared stylesheet", violations)
	}
}

func TestTheCatalogAdoptsTheSharedStylesheetInEveryBadgeWithoutAPolicyViolation(t *testing.T) {
	// scenario: The catalog's badges adopt the one shared sheet and the catalog's policy blocks nothing
	t.Parallel()
	page := catalogued(t)

	page.waitFor(everyBadgeAdoptsTheSharedStylesheet)
	page.evaluate(afterTwoFrames, nil, awaitingThePromise)

	if violations := page.violations(); len(violations) != 0 {
		t.Errorf("the catalog's policy blocked %v while the badges adopted the shared stylesheet", violations)
	}
}

func TestAPrimitiveNotYetDefinedPaintsNothingUntilTheBundleDefinesIt(t *testing.T) {
	// scenario: A ui- element written to the page before its definition is hidden, and is shown the moment it is defined
	t.Parallel()
	page := catalogued(t)

	page.evaluate(`(() => {
  const pending = document.createElement('ui-pending-primitive');
  pending.textContent = 'unstyled words';
  document.body.append(pending);
})()`, nil)

	var before string
	page.evaluate(`getComputedStyle(document.querySelector('ui-pending-primitive')).display`, &before)
	if before != "none" {
		t.Errorf("an undefined primitive is displayed as %q, want none so that no unstyled slotted text flashes", before)
	}

	page.evaluate(`customElements.define('ui-pending-primitive', class extends HTMLElement {})`, nil)

	var after string
	page.evaluate(`getComputedStyle(document.querySelector('ui-pending-primitive')).display`, &after)
	if after == "none" {
		t.Error("a primitive stays hidden after the bundle defines it")
	}
}
