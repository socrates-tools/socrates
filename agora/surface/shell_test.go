package surface_test

import (
	"html"
	"net/http"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/agora/document"
	"github.com/socrates-tools/socrates/agora/surface"
)

func aDocumentWithEveryKindOfWord() document.Document {
	return document.Document{
		Title:    "a title nobody else would write",
		Preamble: "a preamble only this document holds",
		Source: []document.Entry{
			{ID: "c-1", Name: "the name of an entry", Severity: "error", Audience: "needs a person", Location: "internal/cache/read.go:41"},
		},
		Decisions: []document.Decision{
			{
				ID:       "c-1",
				Question: "a question only this document asks",
				Meaning:  "a meaning only this document gives",
				Mode:     document.One,
				Options: []document.Option{
					{Label: "an option label of its own", Meaning: "an option meaning of its own", Command: "patch c-1"},
				},
				Evidence: document.Evidence{{Caption: "a caption of its own", Language: "go", Code: "func OnlyHere() {}\n"}},
			},
		},
		Settled: []document.Settled{{ID: "s-1", Name: "a settled name", Question: "a settled question", Chosen: []string{"a settled choice"}, Words: "settled words"}},
	}
}

func TestThePageIsAShellThatHoldsNoWordOfTheDocument(t *testing.T) {
	// scenario: The page's route serves a shell and nothing the document says
	held := serving(t, aDocumentWithEveryKindOfWord(), 0)

	code, page := held.fetch(t, held.under(""))
	if code != http.StatusOK {
		t.Fatalf("the page was answered with %d", code)
	}

	for _, said := range []string{
		"a title nobody else would write",
		"a preamble only this document holds",
		"the name of an entry",
		"internal/cache/read.go",
		"a question only this document asks",
		"a meaning only this document gives",
		"an option label of its own",
		"an option meaning of its own",
		"a caption of its own",
		"OnlyHere",
		"a settled name",
		"a settled question",
		"a settled choice",
		"settled words",
		"patch c-1",
	} {
		if strings.Contains(page, said) {
			t.Errorf("the page holds %q, which only the document says", said)
		}
	}
}

func TestTheShellCarriesTheAddressTheTokenTheBundleAndOneLineForAPageWithoutScript(t *testing.T) {
	// scenario: The shell holds the address, the token, the bundle and the line for a browser with script turned off
	held := serving(t, twoDecisions(), 0)

	_, served := held.fetch(t, held.under(""))
	page := html.UnescapeString(served)

	for name, expected := range map[string]string{
		"its address":                 `data-base="` + surface.SessionPrefix + held.session + `"`,
		"its token":                   `data-token="` + held.token + `"`,
		"the bundle's script":         `<script src="/assets/surface.js">`,
		"the bundle's stylesheet":     `href="/assets/surface.css"`,
		"a policy that admits script": "script-src 'self'",
	} {
		if !strings.Contains(page, expected) {
			t.Errorf("the shell does not hold %s: %q is missing from\n%s", name, expected, page)
		}
	}

	if got := strings.Count(page, "<noscript>"); got != 1 {
		t.Errorf("the shell holds %d lines for a page with script turned off, want 1", got)
	}
	if !strings.Contains(page, "needs JavaScript") {
		t.Errorf("the line for a page with script turned off does not say it needs JavaScript:\n%s", page)
	}
}
