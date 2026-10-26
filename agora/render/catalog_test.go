package render_test

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/agora/render"
	"golang.org/x/net/html"
)

func theCatalog(t *testing.T) *html.Node {
	t.Helper()

	written, err := render.Catalog()
	if err != nil {
		t.Fatalf("drawing the catalog: %v", err)
	}
	parsed, err := html.Parse(strings.NewReader(written))
	if err != nil {
		t.Fatalf("reading back the catalog: %v", err)
	}
	return parsed
}

func absentFromTheCatalog(declared contract, drawn writings) []string {
	var found []string
	var names []string
	for element := range declared.Elements {
		names = append(names, element)
	}
	for _, element := range names {
		if _, known := drawn[element]; !known {
			found = append(found, fmt.Sprintf("the contract declares <%s>, and the catalog draws no state of it", element))
		}
	}
	slices.Sort(found)
	return found
}

func TestEveryPrimitiveTheContractDeclaresIsDrawnInTheCatalog(t *testing.T) {
	for _, found := range absentFromTheCatalog(declaredByTheContract(t), writingsOf(theCatalog(t))) {
		t.Error(found)
	}
}

func TestTheCatalogWritesNoElementNoAttributeNoSlotAndNoPartTheContractDoesNotDeclare(t *testing.T) {
	for _, found := range undeclared(declaredByTheContract(t), writingsOf(theCatalog(t))) {
		t.Error(found)
	}
}

func TestTheCatalogDrawsEveryAttributeAndSlotOfEveryPrimitiveTheContractDeclares(t *testing.T) {
	for _, found := range unwritten(declaredByTheContract(t), writingsOf(theCatalog(t))) {
		t.Error(found)
	}
}

func TestAPrimitiveTheContractDeclaresAndTheCatalogDoesNotDrawIsReportedNamingIt(t *testing.T) {
	declared := contract{Elements: map[string]element{
		"ui-dialog":  {Slots: []string{"body"}},
		"ui-tooltip": {Attributes: []string{"tone"}},
	}}
	drawn := writings{"ui-dialog": {slots: []string{"body"}}}

	want := []string{`the contract declares <ui-tooltip>, and the catalog draws no state of it`}
	if diff := cmp.Diff(want, absentFromTheCatalog(declared, drawn)); diff != "" {
		t.Errorf("the primitives the catalog lacks differ (-want +got):\n%s", diff)
	}
}

func TestTheCatalogDrawsTheDialogInMoreThanOneState(t *testing.T) {
	drawn := 0
	var visit func(node *html.Node)
	visit = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "ui-dialog" {
			drawn++
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			visit(child)
		}
	}
	visit(theCatalog(t))

	if drawn < 2 {
		t.Errorf("the catalog draws %d dialogs, and the dialog has more than one state to draw", drawn)
	}
}

func framesOf(node *html.Node) []*html.Node {
	var found []*html.Node
	var visit func(node *html.Node)
	visit = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "iframe" && hasAttribute(node, "srcdoc") {
			found = append(found, node)
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			visit(child)
		}
	}
	visit(node)
	return found
}

func TestTheCatalogDrawsTheShellAtANarrowWidthAtNineSixtyAndAtTwelveEighty(t *testing.T) {
	var widths []string
	for _, frame := range framesOf(theCatalog(t)) {
		inside, err := html.Parse(strings.NewReader(attribute(frame, "srcdoc")))
		if err != nil {
			t.Fatalf("reading back a frame of the catalog: %v", err)
		}
		if _, drawn := writingsOf(inside)["ui-app-shell"]; drawn {
			widths = append(widths, attribute(frame, "width"))
		}
	}

	want := []string{"480", "960", "1280"}
	if diff := cmp.Diff(want, widths); diff != "" {
		t.Errorf("the widths the catalog draws the shell at differ (-want +got):\n%s", diff)
	}
}

func TestTheCatalogDrawsEveryControlInMoreThanOneState(t *testing.T) {
	// scenario: The button, field, choice, choice group, copy, expand and file tray each have more than one state drawn
	held := map[string]int{}
	var visit func(node *html.Node)
	visit = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "article" && hasClass(node, "catalog-state") {
			for child := node.FirstChild; child != nil; child = child.NextSibling {
				if child.Type == html.ElementNode && strings.HasPrefix(child.Data, "ui-") {
					held[child.Data]++
				}
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			visit(child)
		}
	}
	visit(theCatalog(t))

	for _, element := range []string{"ui-button", "ui-field", "ui-choice", "ui-choice-group", "ui-copy", "ui-expand", "ui-file-tray"} {
		if held[element] < 2 {
			t.Errorf("the catalog draws %d states of <%s>, and it has more than one to draw", held[element], element)
		}
	}
}
