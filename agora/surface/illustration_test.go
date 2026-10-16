package surface_test

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/agora/document"
	"github.com/socrates-tools/socrates/agora/surface"
)

const theDrawing = `<!doctype html><title>q3</title><p id="total">$1,240.00</p><script>document.getElementById("total").dataset.ran="yes"</script>`

func illustrated(path string) document.Document {
	read := twoDecisions()
	read.Decisions[1].Illustration = document.Illustration{Path: path}
	return read
}

func aFolderHolding(t *testing.T, files map[string]string) string {
	t.Helper()

	folder := t.TempDir()
	for name, contents := range files {
		path := filepath.Join(folder, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("making the folder of %s: %v", name, err)
		}
		if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
			t.Fatalf("writing %s: %v", name, err)
		}
	}
	return folder
}

func servingFrom(t *testing.T, read document.Document, base string) *standing {
	t.Helper()

	held := serving(t, oneOtherDecision(), 0)
	opened, err := held.server.Open(read, written(t, read), base)
	if err != nil {
		t.Fatalf("opening a session on a document carrying an illustration: %v", err)
	}
	held.session, held.token = opened.Session, opened.Token
	return held
}

type response struct {
	code   int
	header http.Header
	body   string
}

func (s *standing) responseTo(t *testing.T, address string) response {
	t.Helper()

	answered, err := http.Get(address)
	if err != nil {
		t.Fatalf("reaching the surface: %v", err)
	}
	defer func() { _ = answered.Body.Close() }()

	said, err := io.ReadAll(answered.Body)
	if err != nil {
		t.Fatalf("reading what the surface said: %v", err)
	}
	return response{code: answered.StatusCode, header: answered.Header, body: string(said)}
}

func (s *standing) drawing(decision, file string) string {
	return s.under("/illustrations/" + url.PathEscape(decision) + "/" + url.PathEscape(file))
}

func TestAnIllustrationNamedRelativeToTheDocumentsFolderIsServedFromTheSession(t *testing.T) {
	// scenario: A document naming an illustration by a path relative to its own folder is served
	folder := aFolderHolding(t, map[string]string{"drawings/q3.html": theDrawing})
	held := servingFrom(t, illustrated("drawings/q3.html"), folder)

	answered := held.responseTo(t, held.drawing("c-2", "q3.html"))

	if answered.code != http.StatusOK {
		t.Fatalf("the illustration was answered with %d", answered.code)
	}
	if answered.body != theDrawing {
		t.Errorf("the illustration came back as %q, want the file as written", answered.body)
	}
	if media := answered.header.Get("content-type"); media != "text/html; charset=utf-8" {
		t.Errorf("the illustration was served as %q", media)
	}
}

func TestAnIllustrationNamedByAnAbsolutePathIsServed(t *testing.T) {
	folder := aFolderHolding(t, map[string]string{"q3.html": theDrawing})
	held := servingFrom(t, illustrated(filepath.Join(folder, "q3.html")), "")

	if answered := held.responseTo(t, held.drawing("c-2", "q3.html")); answered.code != http.StatusOK {
		t.Errorf("an illustration named by an absolute path was answered with %d", answered.code)
	}
}

func TestEveryResponseOfTheIllustrationRouteCarriesTheSandboxAndThePolicy(t *testing.T) {
	// scenario: Opening the illustration route on its own still runs the drawing sandboxed
	folder := aFolderHolding(t, map[string]string{"q3.html": theDrawing})
	held := servingFrom(t, illustrated("q3.html"), folder)

	for named, address := range map[string]string{
		"the drawing":          held.drawing("c-2", "q3.html"),
		"a file it never held": held.drawing("c-2", "other.html"),
		"a decision with none": held.drawing("c-1", "q3.html"),
		"a session never held": held.at(surface.SessionPrefix + "fedcba9876543210/illustrations/c-2/q3.html"),
	} {
		t.Run(named, func(t *testing.T) {
			answered := held.responseTo(t, address)
			policy := answered.header.Get("content-security-policy")

			if !strings.HasPrefix(policy, "sandbox allow-scripts;") {
				t.Errorf("the policy does not open with the sandbox:\n%s", policy)
			}
			if strings.Contains(policy, "allow-same-origin") {
				t.Errorf("the policy lets the drawing keep an origin of its own:\n%s", policy)
			}
			for _, admitted := range []string{"https://cdnjs.cloudflare.com", "https://cdn.jsdelivr.net", "https://unpkg.com", "http://" + held.server.Address()} {
				for _, directive := range []string{"script-src", "style-src", "font-src"} {
					if !strings.Contains(directiveOf(policy, directive), admitted) {
						t.Errorf("%s does not admit %s:\n%s", directive, admitted, policy)
					}
				}
			}
			for _, nothing := range []string{"connect-src", "frame-src", "form-action", "base-uri"} {
				if directiveOf(policy, nothing) != nothing+" 'none'" {
					t.Errorf("%s is not 'none':\n%s", nothing, policy)
				}
			}
			for _, directive := range []string{"img-src", "media-src"} {
				if strings.Contains(directiveOf(policy, directive), "https://") {
					t.Errorf("%s names a remote host:\n%s", directive, policy)
				}
			}
			if sniffing := answered.header.Get("x-content-type-options"); sniffing != "nosniff" {
				t.Errorf("x-content-type-options is %q, want nosniff", sniffing)
			}
			if caching := answered.header.Get("cache-control"); caching != "no-store" {
				t.Errorf("cache-control is %q, want no-store", caching)
			}
		})
	}
}

func directiveOf(policy, named string) string {
	for _, directive := range strings.Split(policy, "; ") {
		if strings.HasPrefix(directive, named+" ") || directive == named {
			return directive
		}
	}
	return ""
}

func TestTheIllustrationRouteServesOnlyWhatTheCopyHolds(t *testing.T) {
	folder := aFolderHolding(t, map[string]string{"q3.html": theDrawing, "secret.html": "not copied"})
	held := servingFrom(t, illustrated("q3.html"), folder)

	for named, address := range map[string]string{
		"a sibling the document never named":    held.drawing("c-2", "secret.html"),
		"a file that climbs out of the folder":  held.under("/illustrations/c-2/..%2F..%2F..%2Fdocument.json"),
		"a decision the document does not hold": held.drawing("c-9", "q3.html"),
		"a type the surface does not serve":     held.drawing("c-2", "q3.exe"),
	} {
		t.Run(named, func(t *testing.T) {
			if answered := held.responseTo(t, address); answered.code != http.StatusNotFound {
				t.Errorf("%s was answered with %d, want %d", named, answered.code, http.StatusNotFound)
			}
		})
	}
}

func TestADecisionIdentifierWithASlashStillNamesItsIllustration(t *testing.T) {
	folder := aFolderHolding(t, map[string]string{"q3.html": theDrawing})
	read := illustrated("q3.html")
	read.Decisions[1].ID = "c/2 ../x"
	held := servingFrom(t, read, folder)

	if answered := held.responseTo(t, held.drawing("c/2 ../x", "q3.html")); answered.code != http.StatusOK || answered.body != theDrawing {
		t.Errorf("an illustration of a decision whose identifier holds a slash was answered with %d", answered.code)
	}
}

func TestAnIllustrationThatCannotBeServedIsRefusedBeforeAnythingIsServed(t *testing.T) {
	// scenario: A path that does not exist, or a file that is not HTML, is refused naming the decision and the path
	folder := aFolderHolding(t, map[string]string{"notes.txt": "words", "q3.html": theDrawing, "drawings/q3.png": "png"})
	held := serving(t, oneOtherDecision(), 0)

	for _, refused := range []struct {
		named string
		path  string
		base  string
	}{
		{"a path that does not exist", "missing.html", folder},
		{"a file that is not HTML", "notes.txt", folder},
		{"an image", "drawings/q3.png", folder},
		{"a folder", "drawings", folder},
		{"a relative path with no folder to read it against", "q3.html", ""},
	} {
		t.Run(refused.named, func(t *testing.T) {
			before := len(held.server.Holds())
			read := illustrated(refused.path)

			_, err := held.server.Open(read, written(t, read), refused.base)
			if err == nil {
				t.Fatalf("a document whose illustration is %s was opened", refused.named)
			}
			for _, named := range []string{"c-2", refused.path} {
				if !strings.Contains(err.Error(), named) {
					t.Errorf("the refusal does not name %q: %v", named, err)
				}
			}
			if holds := len(held.server.Holds()); holds != before {
				t.Errorf("a refused document left the surface holding %d sessions, want %d", holds, before)
			}
			if sessions, err := os.ReadDir(filepath.Join(held.home, "sessions")); err != nil || len(sessions) != before {
				t.Errorf("a refused document left %d session directories on disk (%v), want %d", len(sessions), err, before)
			}
		})
	}
}

func TestAnIllustrationOverTheCeilingIsRefused(t *testing.T) {
	folder := aFolderHolding(t, map[string]string{"q3.html": strings.Repeat("x", surface.IllustrationCeiling+1)})
	held := serving(t, oneOtherDecision(), 0)
	read := illustrated("q3.html")

	if _, err := held.server.Open(read, written(t, read), folder); err == nil {
		t.Fatal("an illustration over the ceiling was opened")
	}
}

func TestARevisionCopiesTheDrawingAfreshAndTheOldCopyGoes(t *testing.T) {
	// scenario: Editing the HTML file and serving the same document again shows the edited drawing
	folder := aFolderHolding(t, map[string]string{"q3.html": theDrawing})
	read := illustrated("q3.html")
	held := servingFrom(t, read, folder)

	copies := func() []string {
		entries, err := os.ReadDir(filepath.Join(held.home, "sessions", held.session, "illustrations"))
		if err != nil {
			t.Fatalf("reading the copies of the session: %v", err)
		}
		var names []string
		for _, entry := range entries {
			names = append(names, entry.Name())
		}
		return names
	}
	first := copies()
	if len(first) != 1 {
		t.Fatalf("one illustration left %d copies", len(first))
	}

	edited := `<!doctype html><p>edited</p>`
	if err := os.WriteFile(filepath.Join(folder, "q3.html"), []byte(edited), 0o644); err != nil {
		t.Fatalf("editing the drawing: %v", err)
	}
	if err := held.server.Revise(held.session, read, written(t, read), folder); err != nil {
		t.Fatalf("revising the document: %v", err)
	}

	if answered := held.responseTo(t, held.drawing("c-2", "q3.html")); answered.body != edited {
		t.Errorf("a revision served %q, want the drawing as it now is on disk", answered.body)
	}
	second := copies()
	if len(second) != 1 || second[0] == first[0] {
		t.Errorf("the copies after a revision are %v, want the new one alone and not %v", second, first)
	}
}

func TestARevisionNamingNoIllustrationRemovesTheCopy(t *testing.T) {
	folder := aFolderHolding(t, map[string]string{"q3.html": theDrawing})
	held := servingFrom(t, illustrated("q3.html"), folder)

	plain := twoDecisions()
	if err := held.server.Revise(held.session, plain, written(t, plain), folder); err != nil {
		t.Fatalf("revising the document: %v", err)
	}

	if answered := held.responseTo(t, held.drawing("c-2", "q3.html")); answered.code != http.StatusNotFound {
		t.Errorf("a drawing the revision no longer names was answered with %d", answered.code)
	}
}

func TestARevisionWhoseIllustrationCannotBeCopiedChangesNothing(t *testing.T) {
	folder := aFolderHolding(t, map[string]string{"q3.html": theDrawing})
	read := illustrated("q3.html")
	held := servingFrom(t, read, folder)

	broken := illustrated("gone.html")
	if err := held.server.Revise(held.session, broken, written(t, broken), folder); err == nil {
		t.Fatal("a revision naming a drawing that does not exist was taken")
	}

	if answered := held.responseTo(t, held.drawing("c-2", "q3.html")); answered.code != http.StatusOK || answered.body != theDrawing {
		t.Errorf("a refused revision disturbed the drawing being served: %d", answered.code)
	}
}
