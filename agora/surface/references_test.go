package surface_test

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/agora/surface"
)

const theDrawingWithItsFiles = `<!doctype html>
<html><head>
<link rel="stylesheet" href="css/site.css">
<style>@import "inline-theme.css"; .inline { background: url(img/inline-bg.png) }</style>
<script src="app.js"></script>
</head><body>
<img src="shots/before.png" alt="before">
<img src="./shots/../shots/after.png?v=2#top" alt="after">
<img src="small.png" srcset="shots/one.webp 1x, shots/two.webp 2x" alt="sized">
<p style="background: url('img/attribute-bg.gif')">styled</p>
<a href="#top">top</a>
<a href="mailto:someone@example.com">mail</a>
<script src="https://cdn.jsdelivr.net/npm/mermaid/dist/mermaid.min.js"></script>
<script src="//cdnjs.cloudflare.com/lib.js"></script>
<img src="/absolute.png" alt="absolute">
<img src="data:image/png;base64,AAAA" alt="inline">
</body></html>`

func theFilesOfADrawing() map[string]string {
	return map[string]string{
		"q3.html":                theDrawingWithItsFiles,
		"css/site.css":           `/* url(commented-out.png) */ @import "theme.css"; @import url('more.css'); body { background: url("../img/bg.png?v=3#x") } @font-face { src: url(../fonts/face.woff2) format("woff2"), url(../fonts/face.woff) }`,
		"css/theme.css":          `.t { background: url(deeper/t.svg) }`,
		"css/more.css":           `@import "site.css"; .m { color: red }`,
		"css/deeper/t.svg":       `<svg xmlns="http://www.w3.org/2000/svg"/>`,
		"inline-theme.css":       `.i { color: blue }`,
		"img/inline-bg.png":      "inline-bg",
		"img/attribute-bg.gif":   "attribute-bg",
		"img/bg.png":             "bg",
		"fonts/face.woff2":       "woff2",
		"fonts/face.woff":        "woff",
		"app.js":                 `window.ran = true`,
		"shots/before.png":       "before",
		"shots/after.png":        "after",
		"small.png":              "small",
		"shots/one.webp":         "one",
		"shots/two.webp":         "two",
		"unrelated/untouched.js": "never referenced",
	}
}

func TestEveryRelativeReferenceOfTheHTMLAndOfTheCSSItReachesIsCopiedAndServedWithItsMediaType(t *testing.T) {
	// scenario: A drawing showing two PNG screenshots side by side, referenced relatively, renders both
	held := servingFrom(t, illustrated("q3.html"), aFolderHolding(t, theFilesOfADrawing()))

	for served, want := range map[string]struct{ contents, media string }{
		"css/site.css":         {theFilesOfADrawing()["css/site.css"], "text/css; charset=utf-8"},
		"css/theme.css":        {theFilesOfADrawing()["css/theme.css"], "text/css; charset=utf-8"},
		"css/more.css":         {theFilesOfADrawing()["css/more.css"], "text/css; charset=utf-8"},
		"css/deeper/t.svg":     {theFilesOfADrawing()["css/deeper/t.svg"], "image/svg+xml"},
		"inline-theme.css":     {"", "text/css; charset=utf-8"},
		"img/inline-bg.png":    {"", "image/png"},
		"img/attribute-bg.gif": {"", "image/gif"},
		"img/bg.png":           {"", "image/png"},
		"fonts/face.woff2":     {"", "font/woff2"},
		"fonts/face.woff":      {"", "font/woff"},
		"app.js":               {"", "text/javascript; charset=utf-8"},
		"shots/before.png":     {"before", "image/png"},
		"shots/after.png":      {"after", "image/png"},
		"small.png":            {"small", "image/png"},
		"shots/one.webp":       {"one", "image/webp"},
		"shots/two.webp":       {"two", "image/webp"},
	} {
		t.Run(served, func(t *testing.T) {
			answered := held.responseTo(t, held.drawing("c-2", served))

			if answered.code != http.StatusOK {
				t.Fatalf("%s was answered with %d", served, answered.code)
			}
			if want.contents != "" && answered.body != want.contents {
				t.Errorf("%s came back as %q, want the file as written", served, answered.body)
			}
			if media := answered.header.Get("content-type"); media != want.media {
				t.Errorf("%s was served as %q, want %q", served, media, want.media)
			}
			if sniffing := answered.header.Get("x-content-type-options"); sniffing != "nosniff" {
				t.Errorf("%s was served with x-content-type-options %q", served, sniffing)
			}
		})
	}
}

func TestAFileNothingReferencesIsNotCopied(t *testing.T) {
	held := servingFrom(t, illustrated("q3.html"), aFolderHolding(t, theFilesOfADrawing()))

	for _, left := range []string{"unrelated/untouched.js", "commented-out.png"} {
		if answered := held.responseTo(t, held.drawing("c-2", left)); answered.code != http.StatusNotFound {
			t.Errorf("%s, which nothing references, was answered with %d", left, answered.code)
		}
	}
}

func TestAnAbsoluteReferenceASchemeAndAFragmentAloneAreLeftAsTheyAreAndNothingIsLookedFor(t *testing.T) {
	// scenario: A reference to a CDN URL is not copied and still loads where ADR 77 admits it
	folder := aFolderHolding(t, theFilesOfADrawing())
	held := servingFrom(t, illustrated("q3.html"), folder)

	answered := held.responseTo(t, held.drawing("c-2", "q3.html"))

	if answered.code != http.StatusOK {
		t.Fatalf("a drawing referencing a CDN, an absolute path and a fragment was answered with %d", answered.code)
	}
	for _, left := range []string{
		`src="https://cdn.jsdelivr.net/npm/mermaid/dist/mermaid.min.js"`,
		`src="//cdnjs.cloudflare.com/lib.js"`,
		`src="/absolute.png"`,
		`href="#top"`,
		`href="mailto:someone@example.com"`,
		`src="data:image/png;base64,AAAA"`,
	} {
		if !strings.Contains(answered.body, left) {
			t.Errorf("the served drawing no longer holds %s", left)
		}
	}
}

func TestAReferenceThatCannotBeFollowedIsRefusedBeforeAnythingIsWrittenOrServed(t *testing.T) {
	// scenario: A reference such as ../secret.png is refused naming the decision and the reference
	// scenario: A referenced file that is missing, or of a type outside the list, is refused naming it
	for _, refused := range []struct {
		named     string
		reference string
		files     map[string]string
	}{
		{"a reference leaving the folder", "../secret.png", map[string]string{"drawings/q3.html": `<img src="../secret.png">`, "secret.png": "secret"}},
		{"a reference leaving the folder by a deeper climb", "sub/../../secret.png", map[string]string{"drawings/q3.html": `<img src="sub/../../secret.png">`, "secret.png": "secret"}},
		{"a stylesheet reaching out of the folder", "../../secret.png", map[string]string{"drawings/q3.html": `<link rel="stylesheet" href="s.css">`, "drawings/s.css": `p { background: url(../../secret.png) }`, "secret.png": "secret"}},
		{"a missing image", "missing.png", map[string]string{"drawings/q3.html": `<img src="missing.png">`}},
		{"a missing stylesheet import", "gone.css", map[string]string{"drawings/q3.html": `<style>@import "gone.css";</style>`}},
		{"a missing file a stylesheet names", "gone.png", map[string]string{"drawings/q3.html": `<link rel="stylesheet" href="s.css">`, "drawings/s.css": `p { background: url(gone.png) }`}},
		{"a missing file a srcset names", "gone@2x.png", map[string]string{"drawings/q3.html": `<img src="a.png" srcset="a.png 1x, gone@2x.png 2x">`, "drawings/a.png": "a"}},
		{"a file of a type outside the list", "notes.txt", map[string]string{"drawings/q3.html": `<a href="notes.txt">notes</a>`, "drawings/notes.txt": "words"}},
		{"an executable", "run.exe", map[string]string{"drawings/q3.html": `<script src="run.exe"></script>`, "drawings/run.exe": "MZ"}},
		{"a folder", "shots", map[string]string{"drawings/q3.html": `<img src="shots">`, "drawings/shots/a.png": "a"}},
	} {
		t.Run(refused.named, func(t *testing.T) {
			files := map[string]string{}
			for name, contents := range refused.files {
				files[name] = contents
			}
			folder := aFolderHolding(t, files)
			held := serving(t, oneOtherDecision(), 0)
			before := len(held.server.Holds())
			read := illustrated("drawings/q3.html")

			_, err := held.server.Open(read, written(t, read), folder)
			if err == nil {
				t.Fatalf("a drawing with %s was opened", refused.named)
			}
			for _, named := range []string{"c-2", refused.reference} {
				if !strings.Contains(err.Error(), named) {
					t.Errorf("the refusal does not name %q: %v", named, err)
				}
			}
			if leaving := strings.Contains(refused.named, "reaching out") || strings.Contains(refused.named, "leaving"); leaving && !strings.Contains(err.Error(), "leaves the folder") {
				t.Errorf("the refusal does not say the reference leaves the folder: %v", err)
			}
			if holds := len(held.server.Holds()); holds != before {
				t.Errorf("a refused drawing left the surface holding %d sessions, want %d", holds, before)
			}
			if sessions, err := os.ReadDir(filepath.Join(held.home, "sessions")); err != nil || len(sessions) != before {
				t.Errorf("a refused drawing left %d session directories on disk (%v), want %d", len(sessions), err, before)
			}
		})
	}
}

func TestAReferenceThroughASymbolicLinkOutOfTheFolderIsRefused(t *testing.T) {
	outside := aFolderHolding(t, map[string]string{"secret.png": "secret"})
	folder := aFolderHolding(t, map[string]string{"q3.html": `<img src="linked/secret.png">`})
	if err := os.Symlink(outside, filepath.Join(folder, "linked")); err != nil {
		t.Skipf("this platform makes no symbolic links: %v", err)
	}
	held := serving(t, oneOtherDecision(), 0)
	read := illustrated("q3.html")

	_, err := held.server.Open(read, written(t, read), folder)

	if err == nil {
		t.Fatal("a reference through a symbolic link out of the folder was followed")
	}
	if !strings.Contains(err.Error(), "linked/secret.png") {
		t.Errorf("the refusal does not name the reference: %v", err)
	}
}

func TestAnIllustrationWhoseFilesTotalOverTheCeilingIsRefusedNamingTheFileThatTippedIt(t *testing.T) {
	// scenario: An illustration whose files total over 16 MiB is refused naming the decision
	folder := aFolderHolding(t, map[string]string{"q3.html": `<img src="one.png"><img src="two.png">`, "one.png": "", "two.png": ""})
	for _, name := range []string{"one.png", "two.png"} {
		if err := os.Truncate(filepath.Join(folder, name), surfaceHalfTheCeiling()); err != nil {
			t.Fatalf("growing %s: %v", name, err)
		}
	}
	held := serving(t, oneOtherDecision(), 0)
	read := illustrated("q3.html")

	_, err := held.server.Open(read, written(t, read), folder)

	if err == nil {
		t.Fatal("an illustration whose files total over the ceiling was opened")
	}
	for _, named := range []string{"c-2", "two.png"} {
		if !strings.Contains(err.Error(), named) {
			t.Errorf("the refusal does not name %q: %v", named, err)
		}
	}
}

func TestAnIllustrationWhoseFilesTotalExactlyTheCeilingIsAccepted(t *testing.T) {
	folder := aFolderHolding(t, map[string]string{"q3.html": `<img src="one.png">`, "one.png": ""})
	if err := os.Truncate(filepath.Join(folder, "one.png"), surfaceCeiling()-int64(len(`<img src="one.png">`))); err != nil {
		t.Fatalf("growing the image: %v", err)
	}
	held := serving(t, oneOtherDecision(), 0)
	read := illustrated("q3.html")

	if _, err := held.server.Open(read, written(t, read), folder); err != nil {
		t.Errorf("an illustration totalling exactly the ceiling was refused: %v", err)
	}
}

func TestAFileReferencedTwiceAndStylesheetsThatImportEachOtherAreCopiedOnceAndEndTheWalk(t *testing.T) {
	folder := aFolderHolding(t, map[string]string{
		"q3.html": `<img src="a.png"><img src="a.png"><link rel="stylesheet" href="a.css">`,
		"a.png":   "a",
		"a.css":   `@import "b.css"; p { background: url(a.png) }`,
		"b.css":   `@import "a.css";`,
	})
	held := servingFrom(t, illustrated("q3.html"), folder)

	for _, served := range []string{"a.png", "a.css", "b.css"} {
		if answered := held.responseTo(t, held.drawing("c-2", served)); answered.code != http.StatusOK {
			t.Errorf("%s was answered with %d", served, answered.code)
		}
	}
}

func TestARevisionCopiesTheFilesAnIllustrationBringsAfreshToo(t *testing.T) {
	folder := aFolderHolding(t, map[string]string{"q3.html": `<img src="a.png">`, "a.png": "first"})
	read := illustrated("q3.html")
	held := servingFrom(t, read, folder)

	if err := os.WriteFile(filepath.Join(folder, "a.png"), []byte("second"), 0o644); err != nil {
		t.Fatalf("editing the image: %v", err)
	}
	if err := held.server.Revise(held.session, read, written(t, read), folder); err != nil {
		t.Fatalf("revising the document: %v", err)
	}

	if answered := held.responseTo(t, held.drawing("c-2", "a.png")); answered.body != "second" {
		t.Errorf("a revision served %q, want the image as it now is on disk", answered.body)
	}
}

func surfaceCeiling() int64 { return int64(surface.IllustrationCeiling) }

func surfaceHalfTheCeiling() int64 { return surfaceCeiling()/2 + 1 }
