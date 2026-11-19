package e2e_test

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	devtoolspage "github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/google/go-cmp/cmp"
)

const recordingEveryElementTheBundleRegisters = `(() => {
  const register = customElements.define.bind(customElements);
  window.__registered = [];
  customElements.define = (name, constructor, options) => {
    window.__registered.push(name);
    register(name, constructor, options);
  };
})()`

const primitivesTheCatalogDraws = `Array.from(document.querySelectorAll('[data-primitive]')).map((held) => held.dataset.primitive)`

func primitivesAmong(registered []string) []string {
	var primitives []string
	for _, name := range registered {
		if strings.HasPrefix(name, "ui-") {
			primitives = append(primitives, name)
		}
	}
	slices.Sort(primitives)
	return primitives
}

func registeredButNotDrawn(registered, drawn []string) []string {
	var found []string
	for _, name := range primitivesAmong(registered) {
		if !slices.Contains(drawn, name) {
			found = append(found, fmt.Sprintf("the library registers <%s>, and the catalog draws no state of it", name))
		}
	}
	return found
}

func drawnButNotRegistered(registered, drawn []string) []string {
	var found []string
	for _, name := range drawn {
		if !slices.Contains(registered, name) {
			found = append(found, fmt.Sprintf("the catalog draws <%s>, and the library registers no such element", name))
		}
	}
	return found
}

func catalogRegisteringWhatTheBundleDefines(t *testing.T) (registered, drawn []string) {
	t.Helper()

	first := catalogued(t)
	var address string
	first.evaluate(`location.href`, &address)

	tab, releaseTab := chromedp.NewContext(first.context)
	t.Cleanup(releaseTab)
	page := watching(t, tab, ".catalog-state", ".catalog-state")
	page.do(
		sending(devtoolspage.AddScriptToEvaluateOnNewDocument, devtoolspage.AddScriptToEvaluateOnNewDocumentParams{Source: recordingEveryElementTheBundleRegisters}),
		chromedp.Navigate(address),
		chromedp.WaitReady(chromedp.CSS(".catalog-state")),
	)
	page.waitFor(`Array.isArray(window.__registered) && window.__registered.includes('ui-dialog')`)
	page.evaluate(`window.__registered`, &registered)
	page.evaluate(primitivesTheCatalogDraws, &drawn)
	return registered, drawn
}

func TestEveryPrimitiveTheLibraryRegistersIsDrawnInTheCatalog(t *testing.T) {
	// scenario: A primitive the library registers has a section in the catalog the page draws
	t.Parallel()
	registered, drawn := catalogRegisteringWhatTheBundleDefines(t)

	if len(primitivesAmong(registered)) == 0 {
		t.Fatal("the bundle registered no primitive, and the catalog has nothing to be compared with")
	}
	for _, found := range registeredButNotDrawn(registered, drawn) {
		t.Error(found)
	}
}

func TestTheCatalogDrawsNoPrimitiveTheLibraryDoesNotRegister(t *testing.T) {
	// scenario: A section of the catalog names an element the library registers
	t.Parallel()
	registered, drawn := catalogRegisteringWhatTheBundleDefines(t)

	for _, found := range drawnButNotRegistered(registered, drawn) {
		t.Error(found)
	}
}

func TestAPrimitiveTheLibraryRegistersAndTheCatalogDoesNotDrawIsReportedNamingIt(t *testing.T) {
	// scenario: A registered primitive the catalog lacks is reported with its element name
	registered := []string{"ui-dialog", "ui-tooltip", "some-other-element"}
	drawn := []string{"ui-dialog"}

	want := []string{"the library registers <ui-tooltip>, and the catalog draws no state of it"}
	if diff := cmp.Diff(want, registeredButNotDrawn(registered, drawn)); diff != "" {
		t.Errorf("the primitives the catalog lacks differ (-want +got):\n%s", diff)
	}
}
