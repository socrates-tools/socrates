package surface_test

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

const theHeightScript = `<script src="/assets/height.js"></script>`

func TestTheHeightScriptEndsTheBodyOfTheCopiedDocumentAndNothingElseChanges(t *testing.T) {
	// scenario: The drawing's own markup and script are otherwise served exactly as written
	written := `<!doctype html>
<html><head><title>q3</title><script>window.ran = true</script></head>
<body>
<p id="total">$1,240.00</p>
<script>document.getElementById("total").dataset.ran = "yes"</script>
</body>
</html>
`
	want := `<!doctype html>
<html><head><title>q3</title><script>window.ran = true</script></head>
<body>
<p id="total">$1,240.00</p>
<script>document.getElementById("total").dataset.ran = "yes"</script>
` + theHeightScript + `</body>
</html>
`
	held := servingFrom(t, illustrated("q3.html"), aFolderHolding(t, map[string]string{"q3.html": written}))

	if answered := held.responseTo(t, held.drawing("c-2", "q3.html")); answered.body != want {
		t.Errorf("the drawing came back as %q, want %q", answered.body, want)
	}
}

func TestTheHeightScriptGoesBeforeTheLastClosingBodyTagWhateverItsCase(t *testing.T) {
	// scenario: The drawing's own markup and script are otherwise served exactly as written
	written := `<p>a string holding </body> in a script</p><BODY><p>x</p></BODY></html>`
	want := `<p>a string holding </body> in a script</p><BODY><p>x</p>` + theHeightScript + `</BODY></html>`
	held := servingFrom(t, illustrated("q3.html"), aFolderHolding(t, map[string]string{"q3.html": written}))

	if answered := held.responseTo(t, held.drawing("c-2", "q3.html")); answered.body != want {
		t.Errorf("the drawing came back as %q, want %q", answered.body, want)
	}
}

func TestADocumentWithNoClosingBodyTagStillEndsWithTheHeightScript(t *testing.T) {
	// scenario: The drawing's own markup and script are otherwise served exactly as written
	held := servingFrom(t, illustrated("q3.html"), aFolderHolding(t, map[string]string{"q3.html": `<p>short</p>`}))

	if answered := held.responseTo(t, held.drawing("c-2", "q3.html")); answered.body != `<p>short</p>`+theHeightScript {
		t.Errorf("the drawing came back as %q, want the script added at its end", answered.body)
	}
}

func TestTheHeightScriptIsAddedToTheEntryDocumentAlone(t *testing.T) {
	// scenario: The drawing's own markup and script are otherwise served exactly as written
	files := map[string]string{
		"q3.html":  `<body><img src="a.png"><script src="app.js"></script><link rel="stylesheet" href="site.css"></body>`,
		"a.png":    "png",
		"app.js":   `document.body.append("</body>")`,
		"site.css": `body { margin: 0 }`,
	}
	held := servingFrom(t, illustrated("q3.html"), aFolderHolding(t, files))

	for _, named := range []string{"a.png", "app.js", "site.css"} {
		if answered := held.responseTo(t, held.drawing("c-2", named)); answered.body != files[named] {
			t.Errorf("%s came back as %q, want it as written", named, answered.body)
		}
	}
}

func TestTheFileOnDiskIsLeftAsWritten(t *testing.T) {
	// scenario: The drawing's own markup and script are otherwise served exactly as written
	folder := aFolderHolding(t, map[string]string{"q3.html": theDrawing})
	servingFrom(t, illustrated("q3.html"), folder)

	held, err := os.ReadFile(filepath.Join(folder, "q3.html"))
	if err != nil {
		t.Fatalf("reading the drawing back: %v", err)
	}
	if string(held) != theDrawing {
		t.Errorf("the file on disk is %q after serving, want it untouched", held)
	}
}

func TestTheHeightScriptTheDrawingLoadsIsServedFromTheSurface(t *testing.T) {
	// scenario: A short drawing's frame is the drawing's height, with no gap under it
	held := serving(t, twoDecisions(), 0)

	answered, body := answeredWithHeaders(t, held.at("/assets/height.js"))

	if answered.StatusCode != http.StatusOK {
		t.Fatalf("the height script was answered with %d", answered.StatusCode)
	}
	if body == "" {
		t.Error("the height script is empty")
	}
}
