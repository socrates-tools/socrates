package e2e_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

const (
	libraryTokens   = "agora/web/packages/ui/tokens"
	landingGround   = "rgb(246, 242, 234)"
	changedGround   = "rgb(17, 34, 51)"
	groundTokenLine = "--limestone-100: #f6f2ea;"
	changedToken    = "--limestone-100: #112233;"
)

type landingPaint struct {
	Ground       string   `json:"ground"`
	BodyFont     string   `json:"bodyFont"`
	LoadedFaces  []string `json:"loadedFaces"`
	WordmarkFont string   `json:"wordmarkFont"`
}

const readingTheLandingPaint = `(async () => {
  await Promise.allSettled(['Cinzel', 'Manrope', 'JetBrains Mono'].map((face) => document.fonts.load('1em "' + face + '"')));
  return {
    ground: getComputedStyle(document.body).backgroundColor,
    bodyFont: getComputedStyle(document.body).fontFamily,
    wordmarkFont: getComputedStyle(document.querySelector('header a span')).fontFamily,
    loadedFaces: [...document.fonts].filter((face) => face.status === 'loaded').map((face) => face.family.replace(/["']/g, '')),
  };
})()`

func servedLandingPage(t *testing.T, site string) string {
	t.Helper()

	served := httptest.NewServer(http.FileServer(http.Dir(site)))
	t.Cleanup(served.Close)
	return served.URL
}

func paintOfTheLandingPage(t *testing.T, site string) landingPaint {
	t.Helper()

	page := openingFor(t, servedLandingPage(t, site), "header", "header")
	var painted landingPaint
	page.evaluate(readingTheLandingPaint, &painted, awaitingThePromise)
	return painted
}

func TestTheAssembledLandingPageHoldsTheLibrarysTokensFontsAndMark(t *testing.T) {
	root := requireMake(t)
	site := assembledSite(t, root)

	for _, stated := range []string{"tokens.css", "fonts.css"} {
		want, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(libraryTokens), stated))
		if err != nil {
			t.Fatalf("reading the library's %s: %v", stated, err)
		}
		got, err := os.ReadFile(filepath.Join(site, stated))
		if err != nil {
			t.Fatalf("the assembled site holds no %s: %v", stated, err)
		}
		if string(got) != string(want) {
			t.Errorf("the assembled site's %s is not the library's", stated)
		}
	}

	held, err := os.ReadDir(filepath.Join(root, filepath.FromSlash(libraryTokens), "fonts"))
	if err != nil {
		t.Fatalf("reading the library's fonts: %v", err)
	}
	for _, font := range held {
		if _, err := os.Stat(filepath.Join(site, "fonts", font.Name())); err != nil {
			t.Errorf("the assembled site holds no fonts/%s: %v", font.Name(), err)
		}
	}
}

func TestTheLandingPageLoadsTheTokensAheadOfItsOwnStylesheet(t *testing.T) {
	site := assembledSite(t, requireMake(t))

	page, err := os.ReadFile(filepath.Join(site, "index.html"))
	if err != nil {
		t.Fatalf("reading the assembled page: %v", err)
	}

	tokens := strings.Index(string(page), `href="tokens.css`)
	own := strings.Index(string(page), `href="site.generated.css`)
	if tokens < 0 || own < 0 || tokens > own {
		t.Errorf("the assembled page loads tokens.css at %d and site.generated.css at %d, and the tokens must come first", tokens, own)
	}
}

func TestTheLandingPageIsPaintedInThePalettesGroundAndTheLibrarysFonts(t *testing.T) {
	t.Parallel()
	site := assembledSite(t, requireMake(t))

	painted := paintOfTheLandingPage(t, site)

	if painted.Ground != landingGround {
		t.Errorf("the landing page's ground is %s, want the palette's limestone %s", painted.Ground, landingGround)
	}
	if !strings.HasPrefix(painted.BodyFont, "Manrope") {
		t.Errorf("the landing page's body is set in %s, want Manrope first", painted.BodyFont)
	}
	if !strings.HasPrefix(painted.WordmarkFont, "Cinzel") {
		t.Errorf("the landing page's wordmark is set in %s, want Cinzel first", painted.WordmarkFont)
	}
	for _, face := range []string{"Cinzel", "Manrope", "JetBrains Mono"} {
		if !slices.Contains(painted.LoadedFaces, face) {
			t.Errorf("the landing page loaded %v, and %s is not among them", painted.LoadedFaces, face)
		}
	}
}

func TestChangingAColourInTheTokensChangesItOnTheLandingPage(t *testing.T) {
	t.Parallel()
	site := assembledSite(t, requireMake(t))

	tokens := filepath.Join(site, "tokens.css")
	held, err := os.ReadFile(tokens)
	if err != nil {
		t.Fatalf("reading the assembled tokens: %v", err)
	}
	if !strings.Contains(string(held), groundTokenLine) {
		t.Fatalf("the assembled tokens state no %q to change", groundTokenLine)
	}

	if before := paintOfTheLandingPage(t, site).Ground; before != landingGround {
		t.Fatalf("the landing page's ground is %s before the change, want %s", before, landingGround)
	}

	changed := strings.Replace(string(held), groundTokenLine, changedToken, 1)
	if err := os.WriteFile(tokens, []byte(changed), 0o644); err != nil {
		t.Fatalf("changing the assembled tokens: %v", err)
	}

	if after := paintOfTheLandingPage(t, site).Ground; after != changedGround {
		t.Errorf("the landing page's ground is %s after the library's limestone changed, want %s", after, changedGround)
	}
}
