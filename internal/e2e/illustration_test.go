package e2e_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/chromedp/chromedp"
	"github.com/socrates-tools/socrates/internal/exit"
)

const probingDrawing = `<!doctype html>
<html><body>
<p id="total">$1,240.00</p>
<script>
document.getElementById('total').dataset.ran = 'yes';
const probe = 'PROBEURL';
const attempt = async (result, name, work) => {
  try { result[name] = await work(); } catch (error) { result[name] = 'threw ' + error.name; }
};
const report = async () => {
  const result = { ran: document.getElementById('total').dataset.ran, text: document.getElementById('total').textContent };
  await attempt(result, 'parentDocument', () => { parent.document.title; return 'read'; });
  await attempt(result, 'token', () => parent.document.body.dataset.token);
  await attempt(result, 'cookie', () => document.cookie);
  await attempt(result, 'localStorage', () => localStorage.getItem('x'));
  await attempt(result, 'sessionStorage', () => sessionStorage.getItem('x'));
  await attempt(result, 'fetch', async () => { await fetch(probe + '/probe/fetch'); return 'fetched'; });
  await attempt(result, 'socket', () => new Promise((resolve) => {
    const socket = new WebSocket(probe.replace('http', 'ws') + '/probe/socket');
    socket.onopen = () => resolve('opened');
    socket.onerror = () => resolve('errored');
    setTimeout(() => resolve('silent'), 1500);
  }));
  await attempt(result, 'beacon', () => String(navigator.sendBeacon(probe + '/probe/beacon', 'x')));
  await attempt(result, 'form', () => {
    const form = document.createElement('form');
    form.method = 'post';
    form.action = probe + '/probe/form';
    document.body.append(form);
    form.submit();
    return 'submitted';
  });
  await attempt(result, 'navigate', () => { top.location = probe + '/probe/navigated'; return 'navigated'; });
  await attempt(result, 'popup', () => (window.open(probe + '/probe/popup') === null ? 'blocked' : 'opened'));
  await new Promise((resolve) => setTimeout(resolve, 400));
  result.admitted = String(window.fromAdmitted);
  result.elsewhere = String(window.fromElsewhere);
  return result;
};
let reporting = null;
window.addEventListener('message', async () => { reporting ||= report(); parent.postMessage(await reporting, '*'); });
</script>
</body></html>`

var drawingThatReports = strings.Replace(probingDrawing, "PROBEURL", "http://127.0.0.1:1", 1)

const frameOfC1 = `document.querySelector('.decision[data-decision="c-1"] iframe.drawing')`

func (s *standIn) record(request *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hosts = append(s.hosts, request.Host+request.URL.Path)
}

func aProbedDrawing(t *testing.T, held *surfaced, name string) *standIn {
	t.Helper()

	reached := &standIn{}
	elsewhere := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		reached.record(request)
		writer.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(elsewhere.Close)

	aDrawing(t, held, name, strings.Replace(probingDrawing, "PROBEURL", elsewhere.URL, 1))
	return reached
}

func aDrawing(t *testing.T, held *surfaced, name, contents string) {
	t.Helper()

	path := filepath.Join(held.where, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("making the folder of %s: %v", name, err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("writing the drawing %s: %v", name, err)
	}
}

func illustrating(t *testing.T, base, path string) string {
	t.Helper()

	var written map[string]any
	if err := json.Unmarshal([]byte(base), &written); err != nil {
		t.Fatalf("reading the document: %v", err)
	}
	written["decisions"].([]any)[0].(map[string]any)["illustration"] = map[string]any{"path": path}

	raw, err := json.Marshal(written)
	if err != nil {
		t.Fatalf("writing the document: %v", err)
	}
	return string(raw)
}

func (r *reader) askTheDrawing() map[string]any {
	r.t.Helper()

	r.evaluate(`window.__reports = []; window.addEventListener('message', (event) => {
  if (event.source === `+frameOfC1+`.contentWindow) window.__reports.push(event.data);
}); true`, nil)
	r.waitFor(`(() => { ` + frameOfC1 + `.contentWindow.postMessage('report', '*'); return window.__reports.length > 0; })()`)

	var reports []map[string]any
	r.evaluate(`window.__reports`, &reports)
	return reports[0]
}

func TestADrawingIsShownOnItsCardWithItsOwnScriptHavingRun(t *testing.T) {
	// scenario: A document naming an illustration by a path relative to its own folder is served
	t.Parallel()
	held := standing(t)
	aDrawing(t, held, "drawings/q3.html", drawingThatReports)

	page := opening(t, held.serve(illustrating(t, oneDecisionDocument, "drawings/q3.html")))

	var framed struct {
		Sandbox string  `json:"sandbox"`
		Width   float64 `json:"width"`
		Height  float64 `json:"height"`
	}
	page.evaluate(`({ sandbox: `+frameOfC1+`.getAttribute('sandbox'), width: `+frameOfC1+`.getBoundingClientRect().width, height: `+frameOfC1+`.getBoundingClientRect().height })`, &framed)
	if framed.Sandbox != "allow-scripts" {
		t.Errorf("the frame is sandboxed with %q, want allow-scripts alone", framed.Sandbox)
	}
	if framed.Width < 100 || framed.Height < 100 {
		t.Errorf("the frame is drawn %vx%v, so the card shows no drawing", framed.Width, framed.Height)
	}

	reported := page.askTheDrawing()
	if reported["ran"] != "yes" || reported["text"] != "$1,240.00" {
		t.Errorf("the drawing reports %v, want its own script to have run over its own markup", reported)
	}
}

func TestADocumentReadFromAFileResolvesTheIllustrationBesideIt(t *testing.T) {
	// scenario: A document naming an illustration by a path relative to its own folder is served
	t.Parallel()
	held := standing(t)

	beside := t.TempDir()
	if err := os.WriteFile(filepath.Join(beside, "q3.html"), []byte(drawingThatReports), 0o644); err != nil {
		t.Fatalf("writing the drawing: %v", err)
	}
	if err := os.WriteFile(filepath.Join(beside, "document.json"), []byte(illustrating(t, oneDecisionDocument, "q3.html")), 0o644); err != nil {
		t.Fatalf("writing the document: %v", err)
	}

	printed, failures, code := held.run("serve", "--document", filepath.Join(beside, "document.json"), "--address", held.bind)
	if code != exit.NothingAsked {
		t.Fatalf("agora serve exited %d: %s", code, failures)
	}

	page := opening(t, strings.TrimSpace(printed))
	if reported := page.askTheDrawing(); reported["ran"] != "yes" {
		t.Errorf("the drawing named beside the document file did not run: %v", reported)
	}
}

func TestADocumentOnStandardInputResolvesTheIllustrationAgainstTheWorkingDirectory(t *testing.T) {
	// scenario: A document piped on standard input resolves the path against the working directory
	t.Parallel()
	held := standing(t)
	aDrawing(t, held, "q3.html", drawingThatReports)

	address := held.serve(illustrating(t, oneDecisionDocument, "q3.html"))

	answered, err := http.Get(address + "/illustrations/c-1/q3.html")
	if err != nil {
		t.Fatalf("reaching the surface: %v", err)
	}
	defer func() { _ = answered.Body.Close() }()

	served, err := io.ReadAll(answered.Body)
	if err != nil {
		t.Fatalf("reading the drawing: %v", err)
	}
	if answered.StatusCode != http.StatusOK || string(served) != drawingThatReports {
		t.Errorf("the drawing in the working directory was answered with %d and %d bytes, want the file as written", answered.StatusCode, len(served))
	}
}

func TestAnIllustrationThatIsMissingOrNotHTMLIsRefusedBeforeAnythingIsServed(t *testing.T) {
	// scenario: A path that does not exist, or names a file that is not HTML, is refused
	t.Parallel()

	for _, refused := range []struct {
		named string
		path  string
	}{
		{"a path that does not exist", "missing/q3.html"},
		{"a file that is not HTML", "notes.txt"},
	} {
		t.Run(refused.named, func(t *testing.T) {
			held := standing(t)
			aDrawing(t, held, "notes.txt", "words")

			printed, failures, code := held.reading(illustrating(t, oneDecisionDocument, refused.path), "serve", "--document", "-")
			if code == exit.NothingAsked {
				t.Fatalf("a document naming %s was served at %q", refused.named, strings.TrimSpace(printed))
			}
			for _, named := range []string{"c-1", refused.path} {
				if !strings.Contains(failures, named) {
					t.Errorf("refusing %s does not name %q: %s", refused.named, named, failures)
				}
			}
			if _, err := os.Stat(filepath.Join(held.surfaceHome(), "surface.json")); err == nil {
				t.Error("a refused document left a surface behind")
			}
		})
	}
}

func TestADocumentCarryingAFragmentOrAStylesheetIsRefusedBeforeAnythingIsServed(t *testing.T) {
	// scenario: A document carrying fragment or stylesheet is refused naming the field
	t.Parallel()

	for field, put := range map[string]func(map[string]any){
		"fragment": func(written map[string]any) {
			written["decisions"].([]any)[0].(map[string]any)["fragment"] = map[string]any{"markup": "<p>drawn</p>"}
		},
		"stylesheet": func(written map[string]any) { written["stylesheet"] = ".panel{color:#222}" },
	} {
		t.Run(field, func(t *testing.T) {
			held := standing(t)

			var written map[string]any
			if err := json.Unmarshal([]byte(oneDecisionDocument), &written); err != nil {
				t.Fatalf("reading the document: %v", err)
			}
			put(written)
			raw, err := json.Marshal(written)
			if err != nil {
				t.Fatalf("writing the document: %v", err)
			}

			printed, failures, code := held.reading(string(raw), "serve", "--document", "-")
			if code == exit.NothingAsked {
				t.Fatalf("a document carrying %s was served at %q", field, strings.TrimSpace(printed))
			}
			if !strings.Contains(failures, field) {
				t.Errorf("refusing a document carrying %s does not name the field: %s", field, failures)
			}
			if _, err := os.Stat(filepath.Join(held.surfaceHome(), "surface.json")); err == nil {
				t.Error("a refused document left a surface behind")
			}
		})
	}
}

func TestADrawingsScriptCannotReachThePageOrAnythingOutside(t *testing.T) {
	// scenario: A drawing's script cannot read the page, send anything out or move the page
	t.Parallel()
	held := standing(t)
	probed := aProbedDrawing(t, held, "q3.html")
	address := held.serve(illustrating(t, oneDecisionDocument, "q3.html"))
	page := opening(t, address)

	reported := page.askTheDrawing()

	for _, reads := range []string{"parentDocument", "token", "cookie", "localStorage", "sessionStorage"} {
		if reported[reads] != "threw SecurityError" {
			t.Errorf("the drawing's attempt to read %s came back %v, want a SecurityError", reads, reported[reads])
		}
	}
	if fetched, _ := reported["fetch"].(string); !strings.HasPrefix(fetched, "threw") {
		t.Errorf("the drawing's fetch came back %v, want a refusal", reported["fetch"])
	}
	if reported["socket"] == "opened" {
		t.Error("the drawing opened a socket")
	}
	if reported["popup"] != "blocked" {
		t.Errorf("the drawing's popup came back %v, want blocked", reported["popup"])
	}

	if asked := probed.reached(); len(asked) != 0 {
		t.Errorf("something reached %v on the drawing's behalf, and nothing a drawing holds may leave it", asked)
	}

	var where string
	page.evaluate(`location.href`, &where)
	if where != address {
		t.Errorf("the drawing moved the page to %q from %q", where, address)
	}
	var kept bool
	page.evaluate(`!!`+frameOfC1, &kept)
	if !kept {
		t.Error("the drawing took its own card away")
	}
}

type standIn struct {
	mu    sync.Mutex
	hosts []string
}

func (s *standIn) reached() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return slices.Clone(s.hosts)
}

func TestADrawingLoadsAScriptFromAnAdmittedOriginAndNotFromAnyOther(t *testing.T) {
	// scenario: A drawing loads a script from an admitted CDN origin and is refused one from another
	t.Parallel()

	reached := &standIn{}
	cdn := httptest.NewTLSServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		reached.mu.Lock()
		reached.hosts = append(reached.hosts, request.Host)
		reached.mu.Unlock()

		writer.Header().Set("content-type", "text/javascript")
		if strings.HasPrefix(request.Host, "cdnjs.cloudflare.com") {
			_, _ = io.WriteString(writer, `window.fromAdmitted = 'loaded';`)
			return
		}
		_, _ = io.WriteString(writer, `window.fromElsewhere = 'loaded';`)
	}))
	t.Cleanup(cdn.Close)
	port := cdn.Listener.Addr().(*net.TCPAddr).Port

	held := standing(t)
	aDrawing(t, held, "q3.html", strings.Replace(drawingThatReports, "<script>\ndocument.getElementById",
		`<script src="https://cdnjs.cloudflare.com/lib.js"></script><script src="https://cdn.not-admitted.test/lib.js"></script><script>
document.getElementById`, 1))
	address := held.serve(illustrating(t, oneDecisionDocument, "q3.html"))

	page := opening(t, address,
		chromedp.IgnoreCertErrors,
		chromedp.Flag("host-resolver-rules", fmt.Sprintf("MAP cdnjs.cloudflare.com 127.0.0.1:%[1]d, MAP cdn.not-admitted.test 127.0.0.1:%[1]d, MAP * ~NOTFOUND, EXCLUDE 127.0.0.1", port)),
	)

	reported := page.askTheDrawing()

	if reported["admitted"] != "loaded" {
		t.Errorf("the script from an admitted origin did not load: %v", reported)
	}
	if reported["elsewhere"] != "undefined" {
		t.Errorf("the script from an origin that is not admitted loaded: %v", reported)
	}
	for _, host := range reached.reached() {
		if strings.HasPrefix(host, "cdn.not-admitted.test") {
			t.Errorf("the browser asked %s for a script, and a policy that refuses it asks nothing", host)
		}
	}
	if !slices.ContainsFunc(reached.reached(), func(host string) bool { return strings.HasPrefix(host, "cdnjs.cloudflare.com") }) {
		t.Errorf("the stand-in for the admitted origin was never asked: %v", reached.reached())
	}
}

func TestTheIllustrationRouteOpenedOnItsOwnStillRunsTheDrawingSandboxed(t *testing.T) {
	// scenario: Opening the illustration route on its own still runs the drawing sandboxed
	t.Parallel()
	held := standing(t)
	aDrawing(t, held, "q3.html", `<!doctype html><body><p id="total">on its own</p><script>
let kept = 'readable';
try { localStorage.getItem('x'); } catch (error) { kept = error.name; }
document.getElementById('total').dataset.storage = kept;
document.getElementById('total').dataset.origin = window.origin;
</script></body>`)
	address := held.serve(illustrating(t, oneDecisionDocument, "q3.html"))

	page := openingFor(t, address+"/illustrations/c-1/q3.html", "#total", "#total")

	var read struct {
		Storage string `json:"storage"`
		Origin  string `json:"origin"`
	}
	page.evaluate(`({ storage: document.getElementById('total').dataset.storage, origin: document.getElementById('total').dataset.origin })`, &read)

	if read.Storage != "SecurityError" || read.Origin != "null" {
		t.Errorf("the drawing opened on its own ran with storage %q and origin %q, want the sandbox's SecurityError and the opaque origin", read.Storage, read.Origin)
	}
}

func TestServingTheSameDocumentAgainShowsTheEditedDrawing(t *testing.T) {
	// scenario: Editing the HTML file and serving the same document again shows the edited drawing
	t.Parallel()
	held := standing(t)
	aDrawing(t, held, "q3.html", drawingThatReports)
	document := illustrating(t, oneDecisionDocument, "q3.html")

	page := opening(t, held.serve(document))
	if reported := page.askTheDrawing(); reported["text"] != "$1,240.00" {
		t.Fatalf("the first drawing reports %v", reported)
	}

	aDrawing(t, held, "q3.html", strings.ReplaceAll(drawingThatReports, "$1,240.00", "$9,999.00"))
	held.serveAgain(document)
	page.waitFor(`document.body.dataset.revision === '2'`)

	if reported := page.askTheDrawing(); reported["text"] != "$9,999.00" {
		t.Errorf("after serving again the drawing reports %v, want the edited file", reported)
	}
}
