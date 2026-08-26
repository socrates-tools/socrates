package render_test

import (
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/rui-ktei/socrates/internal/document"
	"github.com/rui-ktei/socrates/internal/render"
	"golang.org/x/net/html"
)

const answering, sending = "answering", "sending"

func oneDecision(mode document.Mode) document.Document {
	return document.Document{
		Title:    "one thing to settle",
		Preamble: "what follows is what was seen",
		Decisions: []document.Decision{{
			ID:       "c-1",
			Severity: "error",
			Audience: "a person",
			Location: "internal/cache/read.go:41",
			Question: "should the cache be read before the lock is taken",
			Meaning:  "two callers arriving together read the same stale entry",
			Mode:     mode,
			Options: []document.Option{
				{Label: "take the lock first", Meaning: "correctness first", Command: "patch c-1"},
				{Label: "leave it", Meaning: "the staleness is acceptable", Command: `agree c-1 --reason "acceptable"`},
			},
		}},
	}
}

func rendered(t *testing.T, read document.Document) *html.Node {
	t.Helper()

	written, err := render.Page(read, render.Frame{Token: "a-token", Revision: 1})
	if err != nil {
		t.Fatalf("rendering the page: %v", err)
	}
	parsed, err := html.Parse(strings.NewReader(written))
	if err != nil {
		t.Fatalf("reading back the rendered page: %v", err)
	}
	return parsed
}

func elements(node *html.Node, named ...string) []*html.Node {
	var found []*html.Node
	if node.Type == html.ElementNode && slices.Contains(named, node.Data) {
		found = append(found, node)
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		found = append(found, elements(child, named...)...)
	}
	return found
}

func attribute(node *html.Node, named string) string {
	for _, held := range node.Attr {
		if held.Key == named {
			return held.Val
		}
	}
	return ""
}

func hasClass(node *html.Node, named string) bool {
	return slices.Contains(strings.Fields(attribute(node, "class")), named)
}

func text(node *html.Node) string {
	var written strings.Builder
	if node.Type == html.TextNode {
		written.WriteString(node.Data)
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		written.WriteString(text(child))
	}
	return written.String()
}

func TestADocumentRendersWithoutTheCallerWritingMarkup(t *testing.T) {
	read := oneDecision(document.One)
	read.Decisions[0].Fragment = ""

	page := rendered(t, read)

	cards := elements(page, "li")
	if len(cards) != 1 || !hasClass(cards[0], "decision") {
		t.Fatalf("one decision rendered %d cards", len(cards))
	}

	card := cards[0]
	if named := attribute(card, "data-decision"); named != "c-1" {
		t.Errorf("the card names the decision %q, want %q", named, "c-1")
	}
	for _, expected := range []string{
		"should the cache be read before the lock is taken",
		"two callers arriving together read the same stale entry",
		"take the lock first",
		"leave it",
		"internal/cache/read.go:41",
		"a person",
		"error",
	} {
		if !strings.Contains(text(card), expected) {
			t.Errorf("the card does not carry %q:\n%s", expected, text(card))
		}
	}

	for _, choice := range elements(card, "input") {
		if attribute(choice, "type") != "radio" {
			t.Errorf("a single-choice decision renders a %q", attribute(choice, "type"))
		}
	}
}

func TestTheFreeTextFieldCannotBeSuppressed(t *testing.T) {
	for _, mode := range document.Modes() {
		t.Run(string(mode), func(t *testing.T) {
			read := oneDecision(mode)
			fields := reflect.ValueOf(&read.Decisions[0]).Elem()

			for position := range fields.NumField() {
				named := fields.Type().Field(position).Name
				t.Run(named, func(t *testing.T) {
					emptied := oneDecision(mode)
					held := reflect.ValueOf(&emptied.Decisions[0]).Elem().Field(position)
					held.Set(reflect.Zero(held.Type()))

					if !everyCardCarriesTheField(rendered(t, emptied)) {
						t.Errorf("emptying %s removed the field a person writes in", named)
					}
				})
			}

			if !everyCardCarriesTheField(rendered(t, read)) {
				t.Error("a whole decision renders no field for a person to write in")
			}
		})
	}
}

func everyCardCarriesTheField(page *html.Node) bool {
	cards := 0
	for _, card := range elements(page, "li") {
		if !hasClass(card, "decision") {
			continue
		}
		cards++
		fields := 0
		for _, written := range elements(card, "textarea") {
			if hasClass(written, "words-field") {
				fields++
			}
		}
		if fields != 1 {
			return false
		}
	}
	return cards > 0
}

func TestTheSurfaceOffersNoThirdInteraction(t *testing.T) {
	read := oneDecision(document.Many)
	read.Decisions[0].Fragment = `<form><input name="smuggled"><button>publish</button></form>`
	read.Decisions[0].Evidence = document.Evidence{Patch: "--- a/x\n+++ b/x\n@@ -1,2 +1,2 @@\n-old\n+new\n"}

	page := rendered(t, read)

	var parts []string
	for _, control := range elements(page, "input", "textarea", "button", "select", "form", "output", "dialog") {
		part := attribute(control, "data-part")
		if part != answering && part != sending {
			t.Errorf("the page offers a <%s class=%q> that is part of neither answering nor sending", control.Data, attribute(control, "class"))
			continue
		}
		parts = append(parts, part+" "+attribute(control, "class"))
	}

	slices.Sort(parts)
	want := []string{
		"answering choice",
		"answering choice",
		"answering commit",
		"answering reopen",
		"answering words-field",
		"sending discard",
		"sending finish",
		"sending keep",
		"sending message-field",
		"sending send",
	}
	if !slices.Equal(want, parts) {
		t.Errorf("the controls the page offers are %v, want %v", parts, want)
	}
}
