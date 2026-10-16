package surface_test

import (
	"io"
	"net/http"
	"regexp"
	"strings"
	"testing"
)

var fontReference = regexp.MustCompile(`url\(["']?\.?/?(fonts/[A-Za-z0-9-]+\.woff2)`)

type answeredWith struct {
	StatusCode int
	Header     http.Header
}

func answeredWithHeaders(t *testing.T, address string) (answeredWith, string) {
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
	return answeredWith{StatusCode: answered.StatusCode, Header: answered.Header}, string(said)
}

func TestTheTokensStylesheetIsServedToAFrameOfAnOpaqueOrigin(t *testing.T) {
	// scenario: A drawing linking the token stylesheet is painted with the page's own fonts
	held := serving(t, twoDecisions(), 0)

	answered, body := answeredWithHeaders(t, held.at("/assets/tokens.css"))

	if answered.StatusCode != http.StatusOK {
		t.Fatalf("the tokens stylesheet was answered with %d", answered.StatusCode)
	}
	if media := answered.Header.Get("content-type"); !strings.HasPrefix(media, "text/css") {
		t.Errorf("the tokens stylesheet was served as %q", media)
	}
	if allowed := answered.Header.Get("access-control-allow-origin"); allowed != "*" {
		t.Errorf("the tokens stylesheet allows origin %q, want * for a frame whose origin is opaque", allowed)
	}
	if !strings.Contains(body, "--font-body") || !strings.Contains(body, "--aegean-700") {
		t.Errorf("the tokens stylesheet states no type or colour token: %.120s", body)
	}
}

func TestEveryFontTheTokensStylesheetNamesIsServedToAFrameOfAnOpaqueOrigin(t *testing.T) {
	// scenario: A drawing linking the token stylesheet is painted with the page's own fonts
	held := serving(t, twoDecisions(), 0)

	_, body := answeredWithHeaders(t, held.at("/assets/tokens.css"))
	named := fontReference.FindAllStringSubmatch(body, -1)
	if len(named) != 3 {
		t.Fatalf("the tokens stylesheet names %d fonts, want the page's three: %.300s", len(named), body)
	}

	for _, face := range named {
		answered, font := answeredWithHeaders(t, held.at("/assets/"+face[1]))

		if answered.StatusCode != http.StatusOK || font == "" {
			t.Errorf("%s was answered with %d and %d bytes", face[1], answered.StatusCode, len(font))
		}
		if allowed := answered.Header.Get("access-control-allow-origin"); allowed != "*" {
			t.Errorf("%s allows origin %q, want * because a frame fetches a font in CORS mode", face[1], allowed)
		}
	}
}

func TestThePagesOwnScriptIsNotOfferedToEveryOrigin(t *testing.T) {
	// scenario: A drawing that does not link the tokens is painted only by its own styles
	held := serving(t, twoDecisions(), 0)

	answered, _ := answeredWithHeaders(t, held.at("/assets/surface.js"))

	if allowed := answered.Header.Get("access-control-allow-origin"); allowed != "" {
		t.Errorf("the page's script allows origin %q, and only the tokens and their fonts are for a frame to read", allowed)
	}
}
