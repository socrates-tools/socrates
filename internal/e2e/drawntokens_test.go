package e2e_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strings"
	"testing"
)

const tokensLink = `<link rel="stylesheet" href="/assets/tokens.css">`

const drawingReadingTokens = `<!doctype html>
<html><head>HEAD<style>
#line { font-family: var(--font-body); color: var(--aegean-700); }
</style></head><body>
<p id="line">painted</p>
<script>
window.addEventListener('message', async (event) => {
  const names = event.data.names;
  const faces = ['Cinzel', 'Manrope', 'JetBrains Mono'];
  await Promise.allSettled(faces.map((face) => document.fonts.load('1em "' + face + '"')));
  const line = document.getElementById('line');
  const root = getComputedStyle(document.documentElement);
  const tokens = {};
  for (const name of names) tokens[name] = root.getPropertyValue(name).trim();
  parent.postMessage({
    tokens,
    lineFont: getComputedStyle(line).fontFamily,
    lineColour: getComputedStyle(line).color,
    loadedFaces: [...document.fonts].filter((face) => face.status === 'loaded').map((face) => face.family.replace(/["']/g, '')),
  }, '*');
});
</script>
</body></html>`

func drawingReadingTokensWith(head string) string {
	return strings.Replace(drawingReadingTokens, "HEAD", head, 1)
}

var customProperty = regexp.MustCompile(`--[a-z0-9-]+:`)

func tokenNamesServedAt(t *testing.T, address string) []string {
	t.Helper()

	parsed, err := url.Parse(address)
	if err != nil {
		t.Fatalf("reading the address %q: %v", address, err)
	}
	answered, err := http.Get(parsed.Scheme + "://" + parsed.Host + "/assets/tokens.css")
	if err != nil {
		t.Fatalf("reaching the tokens stylesheet: %v", err)
	}
	defer func() { _ = answered.Body.Close() }()
	if answered.StatusCode != http.StatusOK {
		t.Fatalf("the tokens stylesheet was answered with %d", answered.StatusCode)
	}
	served, err := io.ReadAll(answered.Body)
	if err != nil {
		t.Fatalf("reading the tokens stylesheet: %v", err)
	}

	var names []string
	for _, found := range customProperty.FindAllString(string(served), -1) {
		name := strings.TrimSuffix(found, ":")
		if !slices.Contains(names, name) {
			names = append(names, name)
		}
	}
	return names
}

type asTheDrawingSawIt struct {
	Tokens      map[string]string `json:"tokens"`
	LineFont    string            `json:"lineFont"`
	LineColour  string            `json:"lineColour"`
	LoadedFaces []string          `json:"loadedFaces"`
}

func (r *reader) askTheDrawingAboutTokens(names []string) asTheDrawingSawIt {
	r.t.Helper()

	asked, err := json.Marshal(map[string]any{"names": names})
	if err != nil {
		r.t.Fatalf("writing the question: %v", err)
	}
	r.evaluate(`window.__reports = []; window.addEventListener('message', (event) => {
  if (event.source === `+frameOfC1+`.contentWindow && typeof event.data === 'object') window.__reports.push(event.data);
}); true`, nil)
	r.waitFor(`(() => { ` + frameOfC1 + `.contentWindow.postMessage(` + string(asked) + `, '*'); return window.__reports.length > 0; })()`)

	var reports []asTheDrawingSawIt
	r.evaluate(`window.__reports`, &reports)
	return reports[0]
}

func (r *reader) pageToken(name string) string {
	r.t.Helper()

	var value string
	r.evaluate(`getComputedStyle(document.documentElement).getPropertyValue('`+name+`').trim()`, &value)
	return value
}

func (r *reader) pageColourOf(token string) string {
	r.t.Helper()

	var colour string
	r.evaluate(`(() => {
  const probe = document.createElement('i');
  probe.style.color = 'var(`+token+`)';
  document.body.append(probe);
  const colour = getComputedStyle(probe).color;
  probe.remove();
  return colour;
})()`, &colour)
	return colour
}

func TestADrawingLinkingTheTokensIsPaintedWithThePagesFontsAndColours(t *testing.T) {
	// scenario: A drawing linking the token stylesheet is painted with the page's own fonts
	t.Parallel()
	held := standing(t)
	aDrawing(t, held, "q3.html", drawingReadingTokensWith(tokensLink))
	address := held.serve(illustrating(t, oneDecisionDocument, "q3.html"))
	page := opening(t, address)

	seen := page.askTheDrawingAboutTokens(tokenNamesServedAt(t, address))

	if !strings.HasPrefix(strings.ReplaceAll(seen.LineFont, `"`, ""), "Manrope") {
		t.Errorf("the linked drawing's paragraph is set in %q, want the page's body face, Manrope", seen.LineFont)
	}
	if want := page.pageColourOf("--aegean-700"); seen.LineColour != want {
		t.Errorf("the linked drawing's paragraph is coloured %q, want the page's own %q", seen.LineColour, want)
	}
	for _, face := range []string{"Cinzel", "Manrope", "JetBrains Mono"} {
		if !slices.Contains(seen.LoadedFaces, face) {
			t.Errorf("the linked drawing loaded %v, want %s among them, fetched from the surface by a frame with an opaque origin", seen.LoadedFaces, face)
		}
	}
}

func TestEveryTokenTheStylesheetStatesIsTheOneThePagePaintsWith(t *testing.T) {
	// scenario: The tokens in the stylesheet are the ones the page paints with
	t.Parallel()
	held := standing(t)
	aDrawing(t, held, "q3.html", drawingReadingTokensWith(tokensLink))
	address := held.serve(illustrating(t, oneDecisionDocument, "q3.html"))
	page := opening(t, address)

	names := tokenNamesServedAt(t, address)
	for _, known := range []string{"--aegean-700", "--font-body", "--size-body"} {
		if !slices.Contains(names, known) {
			t.Fatalf("the tokens stylesheet states %v, and none of it is %s", names, known)
		}
	}
	seen := page.askTheDrawingAboutTokens(names)

	for _, name := range names {
		if seen.Tokens[name] == "" {
			t.Errorf("the drawing reads no value for %s", name)
		}
		if onThePage := page.pageToken(name); seen.Tokens[name] != onThePage {
			t.Errorf("%s is %q in the drawing and %q on the page", name, seen.Tokens[name], onThePage)
		}
	}
}

func TestADrawingNotLinkingTheTokensIsPaintedOnlyByItsOwnStyles(t *testing.T) {
	// scenario: A drawing that does not link the tokens is painted only by its own styles
	t.Parallel()
	held := standing(t)
	aDrawing(t, held, "q3.html", drawingReadingTokensWith(""))
	address := held.serve(illustrating(t, oneDecisionDocument, "q3.html"))
	page := opening(t, address)

	names := tokenNamesServedAt(t, address)
	seen := page.askTheDrawingAboutTokens(names)

	for _, name := range names {
		if seen.Tokens[name] != "" {
			t.Errorf("the unlinked drawing reads %s as %q, want nothing from the page", name, seen.Tokens[name])
		}
	}
	if strings.Contains(seen.LineFont, "Manrope") {
		t.Errorf("the unlinked drawing's paragraph is set in %q, want the browser's own face", seen.LineFont)
	}
	if want := page.pageColourOf("--aegean-700"); seen.LineColour == want {
		t.Errorf("the unlinked drawing's paragraph is coloured %q, the page's own, though it linked nothing", seen.LineColour)
	}
	if len(seen.LoadedFaces) != 0 {
		t.Errorf("the unlinked drawing loaded %v, want no face", seen.LoadedFaces)
	}
}
