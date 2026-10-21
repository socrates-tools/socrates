package e2e_test

import (
	"slices"
	"testing"
)

const everyDocumentHasLoaded = `[document, ...Array.from(document.querySelectorAll('iframe')).map((frame) => frame.contentDocument)].every((held) => held !== null && held.readyState === 'complete')`

const primitivesWithoutAShadowRoot = `(() => {
  const documents = [document, ...Array.from(document.querySelectorAll('iframe')).map((frame) => frame.contentDocument)];
  const lacking = documents.flatMap((held) => Array.from(held.querySelectorAll('*'))).filter((element) => element.localName.startsWith('ui-') && element.shadowRoot === null);
  return Array.from(new Set(lacking.map((element) => element.localName))).sort();
})()`

const aPrimitiveDrawnWithNoShadowRoot = `(() => {
  customElements.define('ui-bare', class extends HTMLElement {});
  document.body.append(document.createElement('ui-bare'));
})()`

func (r *reader) primitivesLackingAShadowRoot() []string {
	r.t.Helper()

	r.waitFor(everyDocumentHasLoaded)
	var lacking []string
	r.evaluate(primitivesWithoutAShadowRoot, &lacking)
	return lacking
}

func TestEveryPrimitiveTheCatalogDrawsHasAShadowRoot(t *testing.T) {
	// scenario: Every ui- element the catalog draws, in the page and in each shell frame, renders into a shadow root
	t.Parallel()
	page := shellCatalogued(t)

	for _, name := range page.primitivesLackingAShadowRoot() {
		t.Errorf("the catalog draws <%s> without a shadow root, and every primitive draws its parts into one", name)
	}
}

func TestEveryPrimitiveThePageDrawsHasAShadowRoot(t *testing.T) {
	// scenario: Every ui- element a page of two decisions draws renders into a shadow root
	t.Parallel()
	page := openingAt(t, twoDecisionDocument, aWindowWideEnoughForThePanel, aShortWindow)

	for _, name := range page.primitivesLackingAShadowRoot() {
		t.Errorf("the page draws <%s> without a shadow root, and every primitive draws its parts into one", name)
	}
}

func TestAPrimitiveDrawnWithoutAShadowRootIsReportedNamingIt(t *testing.T) {
	// scenario: A ui- element that attaches no shadow root is named by the check the catalog and the page are held to
	t.Parallel()
	page := catalogued(t)
	page.evaluate(aPrimitiveDrawnWithNoShadowRoot, nil)

	if diff := page.primitivesLackingAShadowRoot(); !slices.Equal(diff, []string{"ui-bare"}) {
		t.Errorf("the primitives drawn without a shadow root are %v, want [ui-bare]", diff)
	}
}
